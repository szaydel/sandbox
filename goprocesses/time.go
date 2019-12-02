package main

/*
#include <time.h>
#include <unistd.h>
extern int clock_gettime(clockid_t clock_id, struct timespec *tp);
extern long sysconf(int name);
*/
import "C"

import (
	"time"
)

func ticksToNsecs(ticks int64) int64 {
	var hz_per_sec_c C.long
	var nsecs int64
	hz_per_sec_c = C.sysconf(C._SC_CLK_TCK)
	nsecs = (1e9 * ticks / int64(hz_per_sec_c))
	return nsecs
}

func monotonicClockTicks() int64 {
	var ts C.struct_timespec
	var hz_per_sec_c C.long
	var ns int64
	C.clock_gettime(C.CLOCK_MONOTONIC, &ts)
	hz_per_sec_c = C.sysconf(C._SC_CLK_TCK)

	ns = int64(ts.tv_sec) * 1e9
	ns += int64(ts.tv_nsec)
	return (ns * int64(hz_per_sec_c)) / 1e9
}

func monotonicSinceBoot() time.Duration {
	var ts C.struct_timespec
	C.clock_gettime(C.CLOCK_MONOTONIC, &ts)
	return time.Duration(int64(ts.tv_sec*1e9) + int64(ts.tv_nsec))
}
