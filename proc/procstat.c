#define _GNU_SOURCE

#include <errno.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>
#include <unistd.h>

typedef unsigned long int u_num;

typedef struct {
  u_num user;
  u_num sys;
  u_num start;
  struct timespec ts;
} times_t;

// time_since_boot returns duration since kernel came online in ticks.
static inline int time_since_boot(void) {
  FILE *procuptime;
  int sec, ssec;

  procuptime = fopen("/proc/uptime", "r");
  fscanf(procuptime, "%d.%ds", &sec, &ssec);
  fclose(procuptime);
  return (sec * sysconf(_SC_CLK_TCK)) + ssec;
}

// parse_times extracts CPU times from buf, which is expected to be
// contents of /proc/<PID>/stat file.
static inline times_t parse_times(char *buf) {
#define UTIME_IDX 13
#define STIME_IDX 14
#define CUTIME_IDX 15
#define CSTIME_IDX 16
#define STARTTIME_IDX 21

  times_t t = {
      .user = 0,
      .sys = 0,
      .start = 0,
  };
  char *tok;

  for (size_t i = 0; (tok = strsep(&buf, " ")) != NULL; i++) {
    switch (i) {
    case UTIME_IDX:
    case CUTIME_IDX:
      t.user += strtoul(tok, NULL, 10);
      break;
    case CSTIME_IDX:
    case STIME_IDX:
      t.sys += strtoul(tok, NULL, 10);
      break;
    case STARTTIME_IDX:
      t.start = strtoul(tok, NULL, 10);
      break;
    }
  }
  return t;
}
// process_age returns duration in ticks, given a process' start time.
static inline u_num process_age(u_num start_time) {
  int since_boot = time_since_boot();
  int age_ticks = since_boot - start_time;
  return age_ticks;
}

// lifetime_cpu_ticks_per_tick returns effectively a ratio of CPU times over
// age of the process, i.e. an average load over process' lifetime.
static inline double lifetime_cpu_ticks_per_tick(times_t t) {
  double total_cpu_ticks = (double)(t.user + t.sys);
  return total_cpu_ticks / (double)process_age(t.start);
}

// cpu_ticks_per_tick returns a ratio of CPU time over walltime measured
// between two samples.
// Seconds per second should not be > 1.0, where 1.0 means process spent
// 100% of measured sample interval consuming CPU.
static inline double cpu_ticks_per_tick(times_t first, times_t second) {
  double cpu_times[2], times[2];
  cpu_times[0] = first.user + first.sys;
  cpu_times[1] = second.user + second.sys;
  times[0] = first.ts.tv_sec + ((double)first.ts.tv_nsec / 1e9);
  times[1] = second.ts.tv_sec + ((double)second.ts.tv_nsec / 1e9);
#ifdef DEBUG
  printf("cpu_times[0] = %f cpu_times[1] = %f\n", cpu_times[0], cpu_times[1]);
  printf("times[0] = %f times[1] = %f\n", times[0], times[1]);
#endif
  return (cpu_times[1] - cpu_times[0]) /
         ((times[1] - times[0]) * sysconf(_SC_CLK_TCK));
}

// read_stat reads process stat information into a buffer for given PID.
static inline char *read_stat(char *buf, size_t len, pid_t pid) {
  FILE *input = NULL;
  char *path = malloc(sizeof(char) * 1024);
  if (!path) {
    perror("malloc(...)");
    return NULL;
  }
  asprintf(&path, "/proc/%d/stat", pid);
  input = fopen(path, "r");
  if (!input) {
    free(path);
    return NULL;
  }
  free(path);

  fread(buf, len, 1, input);
  if (ferror(input)) {
    strerror(errno);
    fclose(input);
    return NULL;
  }
  fclose(input);
  return buf;
}

// collect_times repeatedly samples process' /proc/<PID>/stat file and
// populates an allocated times_t array with these observations.
// This array is returned after all samples are gathered. Delay between
// observations is fixed at 1-second.
static inline times_t *collect_times(size_t n, pid_t pid) {
  size_t bufsz = 1024;
  char *buf = malloc(sizeof(char) * bufsz);
  if (!buf) {
    perror("malloc(...)");
    return NULL;
  }
  times_t *tt = malloc(n * sizeof(times_t));
  if (!tt) {
    perror("malloc(...)");
    free(buf);
    return NULL;
  }
  for (size_t i = 0; i < n; i++) {
    if (read_stat(buf, bufsz, pid) == NULL) {
      free(tt);
      free(buf);
      return NULL;
    }
    tt[i] = parse_times(buf);
    if (clock_gettime(CLOCK_MONOTONIC, &tt[i].ts) == -1) {
      free(tt);
      free(buf);
      return NULL;
    }
#ifdef DEBUG
    printf("index[%zu] => %ld(sec) %ld(ns)\n", i, tt[i].ts.tv_sec,
           tt[i].ts.tv_nsec);
#endif
    sleep(1);
  }
  free(buf);
  return tt;
}

// avg_times returns an average value for times_t array by computing
// a diff between neighboring pairs, and then dividing that sum by
// number of observations.
static inline double avg_times(times_t *tt, size_t n) {
  double total = 0;
  for (size_t i = 0; i < n - 1; i++) {
    total += cpu_ticks_per_tick(tt[i], tt[i + 1]);
  }
  return total / (n - 1);
}

void usage(char *name) {
  fprintf(stderr, "usage: %s <PID> [number of 1-second samples]\n", name);
  exit(2);
}

int main(int argc, char *argv[]) {
  if (argc < 2) {
    usage(argv[0]);
  }

  if (strcmp(argv[1], "-h") == 0 || strcmp(argv[1], "--help") == 0) {
    usage(argv[0]);
  }

  pid_t pid = atoi(argv[1]);

  if (!pid) {
    fprintf(stderr, "Failed to parse PID from argument\n");
    return 1;
  }
  int intervals = 2;
  if (argc > 2) {
    int parsed_intervals = atoi(argv[2]);
    if (parsed_intervals) {
      if (parsed_intervals < 2) {
        parsed_intervals = 2;
      }
      intervals = parsed_intervals;
    }
  }

  times_t *points = collect_times(intervals, pid);
  if (points) {
    printf("%f %f %f\n", avg_times(points, intervals),
           cpu_ticks_per_tick(points[0], points[intervals - 1]),
           lifetime_cpu_ticks_per_tick(points[intervals - 1]));
    // pointless to free(points) here
  } else {
    fprintf(stderr, "Failed to collect information, PID not valid?\n");
    return 1;
  }
  return 0;
}
