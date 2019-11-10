package main

import (
	"bytes"
	"fmt"
	"os"
	"time"
)

// ProcStat is essentially parsed contents of /proc/<pid>/stat.
type ProcStat struct {
	// The process ID.
	PID int `json:"pid"`
	// The filename of the executable.
	Comm string `json:"comm"`
	// The process state.
	State string `json:"process_state"`
	// The PID of the parent of this process.
	PPID int `json:"parent_pid"`
	// The process group ID of the process.
	PGRP int `json:"process_group"`
	// The session ID of the process.
	Session int `json:"session"`
	// The controlling terminal of the process.
	TTY int `json:"tty"`
	// The ID of the foreground process group of the controlling terminal of
	// the process.
	TPGID int `json:"term_process_group_id"`
	// The kernel flags word of the process.
	Flags uint `json:"flags"`
	// The number of minor faults the process has made which have not required
	// loading a memory page from disk.
	MinFlt uint `json:"minor_faults"`
	// The number of minor faults that the process's waited-for children have
	// made.
	CMinFlt uint `json:"child_minor_faults"`
	// The number of major faults the process has made which have required
	// loading a memory page from disk.
	MajFlt uint `json:"major_faults"`
	// The number of major faults that the process's waited-for children have
	// made.
	CMajFlt uint `json:"child_major_faults"`
	// Amount of time that this process has been scheduled in user mode,
	// measured in clock ticks.
	UTime uint `json:"user_time_ticks"`
	// Amount of time that this process has been scheduled in kernel mode,
	// measured in clock ticks.
	STime uint `json:"system_time_ticks"`
	// Amount of time that this process's waited-for children have been
	// scheduled in user mode, measured in clock ticks.
	CUTime uint `json:"child_user_time_ticks"`
	// Amount of time that this process's waited-for children have been
	// scheduled in kernel mode, measured in clock ticks.
	CSTime uint `json:"child_system_time_ticks"`
	// For processes running a real-time scheduling policy, this is the negated
	// scheduling priority, minus one.
	Priority int `json:"priority"`
	// The nice value, a value in the range 19 (low priority) to -20 (high
	// priority).
	Nice int `json:"nice"`
	// Number of threads in this process.
	NumThreads int `json:"num_threads"`
	// The time the process started after system boot, the value is expressed
	// in clock ticks.
	Starttime uint64 `json:"starttime"`
	// Virtual memory size in bytes.
	VSize uint `json:"virt_memory_bytes"`
	// Resident set size in pages.
	RSS int `json:"rss_pages"`
}

// ProcInfo maintains information about a single process
type ProcInfo struct {
	Name        string        `json:"name"`
	Role        string        `json:"role"`
	Args        []string      `json:"args"`
	PID         int           `json:"pid"`
	AgeTicks    int64         `json:"age_ticks"`
	AgeDuration time.Duration `json:"age_nanoseconds"`
	S           *ProcStat     `json:"process_stats"`
}

func (ps ProcStat) OnCPUTimeTotal() int64 {
	return int64(ps.STime + ps.UTime + ps.CSTime + ps.CUTime)
}

// Stat returns the current status information of the process.
func (p ProcInfo) Stat() (ProcStat, bool) {
	data, err := ReadFileNoStat(p.path("stat"))
	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			return ProcStat{}, false
		}
		handleErr(err, true)
	}
	var (
		ignore int

		s = ProcStat{PID: p.PID}
		l = bytes.Index(data, []byte("("))
		r = bytes.LastIndex(data, []byte(")"))
	)

	if l < 0 || r < 0 {
		handleErr(fmt.Errorf(
			"unexpected format, couldn't extract comm: %s",
			data,
		), true)
	}

	s.Comm = string(data[l+1 : r])
	_, err = fmt.Fscan(
		bytes.NewBuffer(data[r+2:]),
		&s.State,
		&s.PPID,
		&s.PGRP,
		&s.Session,
		&s.TTY,
		&s.TPGID,
		&s.Flags,
		&s.MinFlt,
		&s.CMinFlt,
		&s.MajFlt,
		&s.CMajFlt,
		&s.UTime,
		&s.STime,
		&s.CUTime,
		&s.CSTime,
		&s.Priority,
		&s.Nice,
		&s.NumThreads,
		&ignore,
		&s.Starttime,
		&s.VSize,
		&s.RSS,
	)
	if err != nil {
		handleErr(err, true)
	}
	return s, true
}

func (p ProcInfo) path(name string) string {
	return fmt.Sprintf("/proc/%d/%s", p.PID, name)
}

func (p ProcInfo) ProcAgeAsTicks() int64 {
	return monotonicClockTicks() - int64(p.S.Starttime)
}

func (p ProcInfo) ProcAgeAsDuration() time.Duration {
	return time.Duration(
		monotonicSinceBoot().Nanoseconds() -
			ticksToNsecs(int64(p.S.Starttime)))
}

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
