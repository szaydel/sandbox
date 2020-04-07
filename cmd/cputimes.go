package main

// CPUTimes tracks two observations of time for a process. There are
// two samples stored previous observation and current, for the purposes
// of computing a delta.
type CPUTimes struct {
	PrevRunTime      int64 // Total time spent running on or off CPU - last
	CurrentRunTime   int64 // Total time spent running on or off CPU - latest
	PrevOnCPUTime    int64 // Time spent on CPU - last
	CurrentOnCPUTime int64 // Time spent on CPU - latest
}

// Delta computes a derivative between current sample and previously taken
// sample. We first take the difference in time on CPU, and divide this value
// by effectively interval between collections. If interval is say 100 ticks,
// and the difference in on CPU time is 90 ticks, rate is 90/100, or 0.9, which
// is to say that process spent 90% of last interval on CPU.
func (t CPUTimes) Delta() float64 {
	// fmt.Printf("CurOnCPU: %d PrevOnCPU: %d CurWT: %d PrevWT: %d\n",
	// t.CurrentOnCPUTime, t.PrevOnCPUTime,
	// t.CurrentRunTime, t.PrevRunTime)
	return float64(t.CurrentOnCPUTime-t.PrevOnCPUTime) / float64(t.CurrentRunTime-t.PrevRunTime)
}

// Reset will zero-out all CPU times. This is useful for instances where
// we can no longer gather statistics, possibly because process was restarted.
func (t *CPUTimes) Reset() {
	t.PrevRunTime = 0
	t.CurrentRunTime = 0
	t.PrevOnCPUTime = 0
	t.CurrentOnCPUTime = 0
}
