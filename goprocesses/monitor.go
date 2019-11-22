package main

import (
	"context"
	"fmt"
	"math"
	"os"
	"time"
)
// MaxNotSeenIntervals is the number of times we allow for a process to not be 
// seen in the process table before removing its traces and stopping associated
// goroutine.
const MaxNotSeenIntervals = 2

// ProcRefreshInterval is the amount of time between rescans of the process 
// table. This is the upper limit to amount of time it can take to detect 
// changes with the processes of interest.
const ProcRefreshInterval = 4 * time.Second
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
func startMonitors(
	ctx context.Context,
	repChan chan *IntervalReport,
	processes func() []*ProcInfo) {
	var processMap = struct {
		t map[string]struct{}
		current map[string]int
		notSeen map[string]int
	}{
		current: make(map[string]int),
		notSeen: make(map[string]int),
	}
	// channels is used for communication with goroutines which we start here
	// for every process of interest.
	var channels = struct {
		pi map[string]chan *ProcInfo
	}{
		pi: make(map[string]chan *ProcInfo),
	}
	for {
		select {
		case <-ctx.Done():
			return
			// shutdown all monitor processes
		default:
			// We want to create this map each time through this loop. This map
			// is expected to be transient and its contents are only good for a
			// single iteration of this loop.
			processMap.t = make(map[string]struct{})
			for _, p := range processes() {
				processMap.t[p.Role] = struct{}{}
				if v, ok := processMap.current[p.Role]; ok {
					if v != p.PID { // This process' PID changed
						fmt.Println("PID changed, take action, save pid")
						processMap.current[p.Role] = p.PID
						channels.pi[p.Role] <- p
					}
				} else { // This process' Role is not already in the map
					processMap.current[p.Role] = p.PID
					processMap.notSeen[p.Role] = 0
					// Create channels if this is a new process which we have
					// never seen before and do not already have its role in 
					// the processMap.current map. If this is a new process ID 
					// for a previously seen role, we will already have these,
					// and instead we just update the process ID above.
					channels.pi[p.Role] = make(chan *ProcInfo)
					go monitor(channels.pi[p.Role], 
						repChan)
					channels.pi[p.Role] <- p
					fmt.Printf("Added %s => %d to map\n", p.Role, p.PID)
				}
			}
			for role := range processMap.current {
				if _, ok := processMap.t[role]; !ok {
					if processMap.notSeen[role] < MaxNotSeenIntervals {
						processMap.notSeen[role]++
						continue
					}
					// We need to notify corresponding goroutine that it needs 
					// to shutdown! After telling relevant goroutine to stop, 
					// remove the no longer existing role from the current map.
					delete(processMap.current, role)
					// this signals associated goroutine to stop and return, 
					// otherwise we are going to have leaking goroutines.
					close(channels.pi[role])
					channels.pi[role] = nil
					delete(channels.pi, role)
					delete(processMap.notSeen, role)
				}
			}
			// It may take this much time to detect that a process got restarted
			// or that a new process was added to system.
			time.Sleep(ProcRefreshInterval)
		}
	}
}

func monitor(p <-chan *ProcInfo, r chan<- *IntervalReport) {
	const window = 10
	var watching *ProcInfo
	var samples = make([]float64, window)
	var counter uint64
	var times CPUTimes
	var lifetimeRate float64
	var osPageSize = os.Getpagesize()
	for {
		select {
		case v := <-p:
			watching = v
			if watching != nil {
			fmt.Printf("monitoring: %s with PID: %d %p\n", watching.Role, watching.PID, watching)
			} else {
				// This goroutine is expected to go away now because the 
				// monitored process was removed by system from process table.
				return
			}

		default:
			var s ProcStat
			var ok bool
			if s, ok = watching.Stat(); ok {
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
				samples[counter%window] = math.NaN()
				times.Reset()
			}
			counter++
			// fmt.Printf("counter: %d | %+v\n", counter, samples)
			if counter >= window {
				// fmt.Printf("DELTA: %f | Avg: %f Latest: %f\n", times.Delta(), avg(samples), lifetimeRate)
				r <- &IntervalReport{
					PID:             watching.PID,
					Role:            watching.Role,
					Timestamp:       time.Now(),
					WindowRate:      avg(samples),
					StandardDev:     stddev(samples),
					LifetimeRate:    lifetimeRate,
					CurrentRate:     times.Delta(),
					VirtMemoryBytes: s.VSize,
					RSSBytes:        s.RSS * osPageSize,
				}
				// fmt.Printf("counter: %d | avg: %f\n", counter, avg(samples))
			}

			// times.Delta()
			time.Sleep(1 * time.Second)
		}
	}
}
