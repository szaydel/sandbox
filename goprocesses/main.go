package main

/*
#include <time.h>
#include <unistd.h>
extern int clock_gettime(clockid_t clock_id, struct timespec *tp);
extern long sysconf(int name);
*/
import "C"

import (
	// "encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func handleErr(e error, doPanic bool) {
	if e != nil && e != io.EOF {
		if doPanic {
			panic(e)
		}
		fmt.Fprintf(os.Stderr, "Error: %v\n", e)
	}
}

func ticksToNsecs(ticks int64) int64 {
	var hz_per_sec_c C.long
	var secs int64
	hz_per_sec_c = C.sysconf(C._SC_CLK_TCK)
	secs = (ticks / int64(hz_per_sec_c))
	return secs * 1e9
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

func nullByteToSpace(b []byte) []byte {
	for i, v := range b {
		if v == 0x0 {
			b[i] = 0x20 // ASCII space character
		}
	}
	return b
}

// type Count int
// type Role string

// const Proxy Role = "proxy"
// const Worker Role = "worker"
// const Manager Role = "manager"
// const Logger Role = "logger"

// func expectedCount(r Role) Count {
// 	m := map[Role]Count{
// 		Proxy:   2,
// 		Worker:  1,
// 		Manager: 1,
// 		Logger:  1,
// 	}
// 	return m[r]
// }

// CommandLine is a representation of a process' arguments and name
// separated into a slice of argument strings and a name string.
type CommandLine struct {
	args        []string
	programName string
}

// ProgramName is the name of the program for the given process, often
// known as argv[0].
func (c CommandLine) ProgramName() string {
	return c.programName
}

// Args is the set of arguments, sans what would be argv[0], i.e. ProgramName.
// It is effectively argv[1...N]
func (c CommandLine) Args() []string {
	return c.args
}

func cmdLineArgs(pid int) *CommandLine {
	var buf = make([]byte, 256)
	var f *os.File
	var err error
	var n int
	var programName string
	var restOfCmdline []string

	var cmdLinePath = fmt.Sprintf("/proc/%d/cmdline", pid)
	if f, err = os.Open(cmdLinePath); err != nil {
		handleErr(err, false)
	}
	defer f.Close()
	if n, err = f.Read(buf); err != nil {
		handleErr(err, false)
	}
	cmdLineSlice := strings.Split(string(nullByteToSpace(buf[:n])), " ")
	programName = cmdLineSlice[0]

	if len(cmdLineSlice) > 1 {
		restOfCmdline = cmdLineSlice[1:]
	} else {
		restOfCmdline = []string{}
	}
	return &CommandLine{
		programName: programName,
		args:        restOfCmdline,
	}
}

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
		// fmt.Printf("procfile: %s\n", procfile)
		pSlc := strings.Split(procfile, "/")
		// FIXME: Check length of slice
		pid, err := strconv.Atoi(pSlc[len(pSlc)-1])
		handleErr(err, true)
		args := cmdLineArgs(pid)
		if args.ProgramName() == name {
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

// IntervalReport is a point in time view of process' CPU usage with three
// figures, WindowRate, LifeTimeRate and CurrentRate.
// WindowRate - an average of samples over several intervals, which effectively
// makes the data smoother.
// LifeTimeRate - rate of time spent on CPU over total process' runtime,
// computed over the entire lifetime of process; least volatile.
// CurrentRate - derivative between two interval samples; most volatile.
type IntervalReport struct {
	Timestamp    time.Time
	WindowRate   float64
	LifetimeRate float64
	CurrentRate  float64
}

func startIntervalReport(c <-chan *IntervalReport) {
	for {
		select {
		case v := <-c:
			fmt.Printf("%+v\n", v)
		default:
			<-time.NewTimer(1 * time.Second).C
		}
	}
}

// startMonitors periodically scans the process table by reading through /proc
// and picks out only those processes that we are interested in. These processes
// are then added to a map of process roles to PIDs, where a role is something
// like logger, or manager or worker-XX, etc. For each entry in this map, a
// channel is created in a corresponding map. Also, a goroutine is started for
// every process, and the channel for the given process is passed to this
// goroutine, establishing a one-way communication mechanism. This channel is in
// essence an updates channel. Initially, information about each process that we
// intend to monitor is passed to the goroutine on the other end of the channel,
// and then subsequently, any time this process is restarted, new information
// about the process is passed to goroutine responsible for this process.
// As a side-effect of these periodic checks, if we detect at some point a
// process that is not already in the map, we begin to track this process and
// create a new monitor goroutine for it.
func startMonitors(stopChan chan struct{}, repChan chan *IntervalReport) {
	var processMap = make(map[string]int)
	var channelsMap = make(map[string]chan *ProcInfo)
	// var ivReportChan = make(chan *IntervalReport)
	for {
		processes := findProcsByName("bro")
		for _, p := range processes {
			if v, ok := processMap[p.Role]; ok {
				if v != p.PID { // This process' PID changed
					fmt.Println("PID changed, take action, save pid")
					processMap[p.Role] = p.PID
					channelsMap[p.Role] <- p
				}
			} else { // This process' Role is not already in the map
				processMap[p.Role] = p.PID
				channelsMap[p.Role] = make(chan *ProcInfo)
				go monitor(channelsMap[p.Role], repChan)
				channelsMap[p.Role] <- p
				fmt.Printf("Added %s => %d to map\n", p.Role, p.PID)
			}
		}
		// It may take this much time to detect that a process got restarted
		// or that a new process was added to system.
		time.Sleep(4 * time.Second)
	}
}

func monitor(p <-chan *ProcInfo, r chan<- *IntervalReport) {
	const window = 10
	var watching *ProcInfo
	var samples = make([]float64, window)
	var counter uint64
	var times CPUTimes
	var lifetimeRate float64
	for {
		select {
		case v := <-p:
			watching = v
			fmt.Printf("monitoring: %s with PID: %d %p\n", watching.Role, watching.PID, watching)
		default:
			if s, ok := watching.Stat(); ok {
				lifetimeRate = float64(s.OnCPUTimeTotal()) / float64(watching.ProcAgeAsTicks())
				samples[counter%window] = lifetimeRate

				times.PrevOnCPUTime = times.CurrentOnCPUTime
				times.PrevRunTime = times.CurrentRunTime

				// If we don't have any value for previous runtime, this is the
				// first time we gather stats. In this case we set both current
				// and previous values to the sample we just collected.
				// If previous value is larger than current value, we are no
				// longer looking at same process. We behave just like we would
				// on first run, setting both previous and current values to the
				// sample we just collected.
				if times.PrevRunTime > times.CurrentRunTime ||
					times.PrevRunTime == 0 {
					times.PrevOnCPUTime = s.OnCPUTimeTotal()
					times.CurrentOnCPUTime = s.OnCPUTimeTotal()
					times.PrevRunTime = watching.ProcAgeAsTicks()
					times.CurrentRunTime = watching.ProcAgeAsTicks()
				} else {
					times.CurrentOnCPUTime = s.OnCPUTimeTotal()
					times.CurrentRunTime = watching.ProcAgeAsTicks()
				}
			} else {
				samples[counter%10] = math.NaN()
				times.Reset()
			}
			counter++
			// fmt.Printf("counter: %d | %+v\n", counter, samples)
			if counter >= 10 {
				// fmt.Printf("DELTA: %f | Avg: %f Latest: %f\n", times.Delta(), avg(samples), lifetimeRate)
				r <- &IntervalReport{
					Timestamp: time.Now(),
					WindowRate: avg(samples),
					LifetimeRate: lifetimeRate,
					CurrentRate: times.Delta(),
				}
				// fmt.Printf("counter: %d | avg: %f\n", counter, avg(samples))
			}

			// times.Delta()
			time.Sleep(1 * time.Second)
		}
	}
}

func main() {
	// processes := findProcsByName("bro")
	// j, _ := json.Marshal(processes)
	// fmt.Printf("%v\n", string(j))

	shutdownChan := make(chan struct{})
	intervalReportChan := make(chan *IntervalReport)
	go startMonitors(shutdownChan, intervalReportChan)
	go startIntervalReport(intervalReportChan)
	time.Sleep(100 * time.Second)
	shutdownChan <- struct{}{}
}
