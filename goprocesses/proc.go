package main

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
