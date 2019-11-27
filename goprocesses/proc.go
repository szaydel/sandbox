package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// ReadFileNoStat uses ioutil.ReadAll to read contents of entire file.
// This is similar to ioutil.ReadFile but without the call to os.Stat, because
// many files in /proc and /sys report incorrect file sizes (either 0 or 4096).
// Reads a max file size of 512kB.  For files larger than this, a scanner
// should be used.
func ReadFileNoStat(filename string) ([]byte, error) {
	const maxBufferSize = 1024 * 512

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := io.LimitReader(f, maxBufferSize)
	return ioutil.ReadAll(reader)
}

func buildProcInfo(procfile string) *ProcInfo {
	pSlc := strings.Split(procfile, "/")
	// FIXME: Check length of slice
	pid, err := strconv.Atoi(pSlc[len(pSlc)-1])
	handleErr(err, true)
	args := cmdLineArgs(pid)
	proci := &ProcInfo{}
	proci.Name = args.ProgramName()
	proci.Args = args.Args()
	if len(proci.Args) > 1 {
		proci.Role = proci.Args[1]
	} else {
		proci.Role = "unknown"
	}
	proci.PID = pid
	var s ProcStat
	var ok bool
	if s, ok = proci.Stat(); !ok {
		return nil
	}
	proci.S = &s
	proci.AgeTicks = proci.ProcAgeAsTicks()
	proci.AgeDuration = proci.ProcAgeAsDuration()
	return proci
}

func findProcsByName(name string) []*ProcInfo {
	paths, err := filepath.Glob("/proc/[0-9]*")
	handleErr(err, true)
	var piSlc = make([]*ProcInfo, 0)
	for _, procfile := range paths {
		pSlc := strings.Split(procfile, "/")
		// FIXME: Check length of slice
		pid, err := strconv.Atoi(pSlc[len(pSlc)-1])
		handleErr(err, true)
		if isTargetProcess(pid, name) {
			//args := cmdLineArgs(pid)
			// if args.ProgramName() == name {
			// If buildProcInfo returns nil, a process is likely no longer valid
			// and instead of adding it to this slice, we skip it.
			// This check runs periodically and if the process that just went
			// away is restarted, it will get picked-up on next run.
			pi := buildProcInfo(procfile)
			if pi != nil {
				piSlc = append(piSlc, pi)
			}
		}
	}
	return piSlc
}

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
	PIDChaged   bool          `json:"pid_changed"`
	AgeTicks    int64         `json:"age_ticks"`
	AgeDuration time.Duration `json:"age_nanoseconds"`
	S           *ProcStat     `json:"process_stats"`
}

// OnCPUTimeTotal returns total amount of time process and its children spent
// on CPU.
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
