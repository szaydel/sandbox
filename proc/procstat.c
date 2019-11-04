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
  return (sec*sysconf(_SC_CLK_TCK))+ssec;
}

// parse_times extracts CPU times from buf, which is expected to be
// contents of /proc/<PID>/stat file.
static inline times_t parse_times(char *buf) {
  #define UTIME_IDX 13
  #define STIME_IDX 14
  #define STARTTIME_IDX 21

  times_t t = {};
  char *tok;

  for (size_t i = 0 ; (tok = strsep(&buf, " ")) != NULL ; i++) {
      if (i == UTIME_IDX) {
        t.user = strtoul(tok, NULL, 10);
      } else if (i == STIME_IDX) {
        t.sys = strtoul(tok, NULL, 10);
      } else if (i == STARTTIME_IDX) {
        t.start = strtoul(tok, NULL, 10);
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

// lifetime_avg_load returns effectively a ratio of CPU times over
// age of the process, i.e. an average load over process' lifetime.
static inline double lifetime_avg_load(times_t t) {
  double total_cpu_ticks = (double)(t.user + t.sys);
  return total_cpu_ticks / (double)process_age(t.start);
}

// load_diff returns a relative difference between two samples.
// This difference should not be greater than 1 assuming a single core.
static inline double load_diff(times_t first, times_t second) {
  double cpu_times[2], times[2];
  cpu_times[0] = first.user + first.sys;
  cpu_times[1] = second.user + second.sys;
  times[0] = first.ts.tv_sec + ((double)(first.ts.tv_sec) / 1e9);
  times[1] = second.ts.tv_sec + ((double)(second.ts.tv_sec) / 1e9);
  return (cpu_times[1] - cpu_times[0]) /
        ((times[1] - times[0]) * sysconf(_SC_CLK_TCK));
}

// read_stat reads process stat information into a buffer for given PID.
static inline char *read_stat(char *buf, size_t len, pid_t pid) {
    FILE *input = NULL;
    char *path = malloc(sizeof(char) * 1024);
    asprintf(&path, "/proc/%d/stat", pid);
    input = fopen(path, "r");
    if(!input) {
      free(path);
      return NULL;
    }
    free(path);

  fread(buf, len, 1, input);
  if (ferror(input)) {
    strerror(errno);
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
  times_t *tt = malloc(n * sizeof(times_t));
  for (size_t i = 0 ; i < n ; i++) {
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
    printf("index[%zu] => %ld(sec) %ld(ns)\n",
        i, tt[i].ts.tv_sec, tt[i].ts.tv_nsec);
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
    for (size_t i = 0; i < n-1; i++) {
        total += load_diff(tt[i], tt[i+1]);
    }
    return total / n;
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
    printf("%f %f %f\n",
        avg_times(points, intervals),
        load_diff(points[0], points[intervals-1]),
        lifetime_avg_load(points[intervals-1])
    );
  // pointless to free(points) here
  } else {
      fprintf(stderr, "Failed to collect information, PID not valid?\n");
      return 1;
  }
  return 0;
}
