package main

import (
	"fmt"
	"math"
	"time"
)

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
					PID:          watching.PID,
					Role:         watching.Role,
					Timestamp:    time.Now(),
					WindowRate:   avg(samples),
					LifetimeRate: lifetimeRate,
					CurrentRate:  times.Delta(),
				}
				// fmt.Printf("counter: %d | avg: %f\n", counter, avg(samples))
			}

			// times.Delta()
			time.Sleep(1 * time.Second)
		}
	}
}
