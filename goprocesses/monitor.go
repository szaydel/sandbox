package main

import (
	"context"
	"errors"
	"log"
	"math"
	"os"
	"time"
)

// MaxNotSeenIntervals is the number of times we allow for a process to not be
// seen in the process table before removing its traces and stopping associated
// goroutine.
const MaxNotSeenIntervals = 5

// ProcRefreshInterval is the amount of time between rescans of the process
// table. This is the upper limit to amount of time it can take to detect
// changes with the processes of interest.
const ProcRefreshInterval = 4 * time.Second

type MonitoredProcesses struct {
	transient  map[string]struct{}
	persistent map[string]int
	notSeen    map[string]int
}

func (mp *MonitoredProcesses) NewTransient() map[string]struct{} {
	mp.transient = make(map[string]struct{})
	return mp.transient
}

func (mp *MonitoredProcesses) InsertIntoTransient(role string) {
	mp.transient[role] = struct{}{}
}

func (mp *MonitoredProcesses) RegisterNewProcess(role string, pid int) {
	mp.persistent[role] = pid
	mp.notSeen[role] = 0
}

func (mp *MonitoredProcesses) pidOf(role string) (int, bool) {
	if pid, ok := mp.persistent[role]; ok {
		return pid, true
	}
	return -1, false
}

func (mp *MonitoredProcesses) PidChanged(role string, pid int) (bool, error) {
	if lastKnownPid, ok := mp.pidOf(role); ok {
		return pid != lastKnownPid, nil
	}
	return false, errors.New("this role not seen previously")
}

func (mp *MonitoredProcesses) UpdatePid(role string, pid int) {
	mp.persistent[role] = pid
}

func (mp *MonitoredProcesses) ResetNotSeen(role string) {
	mp.notSeen[role] = 0
}

func (mp *MonitoredProcesses) IncrNotSeen(role string) {
	if _, ok := mp.notSeen[role]; !ok {
		mp.notSeen[role] = 1
		return
	}
	mp.notSeen[role]++
}

func (mp *MonitoredProcesses) RemoveMonitored(role string) {
	delete(mp.persistent, role)
	delete(mp.notSeen, role)
	// We do not do anything about the transient map because it is recreated 
	// frequently unlike the other two maps.
}

func (mp *MonitoredProcesses) NotSeenCount(role string) int {
	if count, ok := mp.notSeen[role]; ok {
		return count
	}
	mp.notSeen[role] = 0 // If for some reason it is not in map, add to map.
	return 0
}

func (mp *MonitoredProcesses) NotSeen() map[string]int {
	return mp.notSeen
}

func (mp *MonitoredProcesses) NotSeenFewerThan(role string, count int) bool {
	return mp.notSeen[role] < count
}

func (mp *MonitoredProcesses) Persistent() map[string]int {
	return mp.persistent
}

func (mp *MonitoredProcesses) Transient() map[string]struct{} {
	return mp.transient
}

func (mp *MonitoredProcesses) InTransient(role string) bool {
	_, ok := mp.transient[role]
	return ok
}

func NewMonitoredProcesses() *MonitoredProcesses {
	return &MonitoredProcesses{
		transient:  make(map[string]struct{}),
		persistent: make(map[string]int),
		notSeen:    make(map[string]int),
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
func startMonitors(
	ctx context.Context,
	repChan chan *IntervalReport,
	processes func() []*ProcInfo) {
	// var processMap = struct {
	// 	t       map[string]struct{}
	// 	current map[string]int
	// 	notSeen map[string]int
	// }{
	// 	current: make(map[string]int),
	// 	notSeen: make(map[string]int),
	// }
	var mp = NewMonitoredProcesses()
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
			// processMap.t = make(map[string]struct{})
			mp.NewTransient()
			for _, p := range processes() {
				// processMap.t[p.Role] = struct{}{}
				mp.InsertIntoTransient(p.Role)
				if changed, err := mp.PidChanged(p.Role, p.PID); err == nil {
					if changed {
						if pid, ok := mp.pidOf(p.Role); ok {
							log.Printf("PID for process %s changed from %d to %d",
								p.Role, pid, p.PID)
						}
						mp.UpdatePid(p.Role, p.PID)
						// if v, ok := processMap.current[p.Role]; ok {
						// if v != p.PID { // This process' PID changed
						// log.Printf("PID for process %s changed from %d to %d",
						// p.Role, v, p.PID)
						// processMap.current[p.Role] = p.PID
						p.PIDChaged = true
						channels.pi[p.Role] <- p
					}
				} else { // This process' Role is not already in the map
					mp.RegisterNewProcess(p.Role, p.PID)
					// processMap.current[p.Role] = p.PID
					// processMap.notSeen[p.Role] = 0
					// Create channels if this is a new process which we have
					// never seen before and do not already have its role in
					// the processMap.current map. If this is a new process ID
					// for a previously seen role, we will already have these,
					// and instead we just update the process ID above.
					channels.pi[p.Role] = make(chan *ProcInfo)
					go monitor(channels.pi[p.Role],
						repChan)
					channels.pi[p.Role] <- p
					log.Printf("Added %s with PID %d to map", p.Role, p.PID)
				}
			}
			// for role := range processMap.current {
			for role, pid := range mp.Persistent() {
				// if _, ok := processMap.t[role]; !ok {
				if !mp.InTransient(role) {
					if mp.NotSeenFewerThan(role, MaxNotSeenIntervals) {
						mp.IncrNotSeen(role)
						log.Printf("PID %d for process %s no longer seen", pid, role)
						continue
					}
					// if processMap.notSeen[role] < MaxNotSeenIntervals {
					// 	processMap.notSeen[role]++
					// 	log.Println("Not seen")
					// 	continue
					// }
					// We need to notify corresponding goroutine that it needs
					// to shutdown! After telling relevant goroutine to stop,
					// remove the no longer existing role from the current map.
					// delete(processMap.current, role)

					// this signals associated goroutine to stop and return,
					// otherwise we are going to have leaking goroutines.
					log.Printf("Removing %s from list of monitored processes", role)
					close(channels.pi[role])
					channels.pi[role] = nil
					delete(channels.pi, role)
					// delete(processMap.notSeen, role)
					mp.RemoveMonitored(role)
				} else {
					// If process was not seen for whatever reason and is now
					// seen again, reset the count to make sure that next time
					// process is not seen again, we again start counting from
					// a zero counter.
					if mp.NotSeenCount(role) > 0 {
						mp.ResetNotSeen(role)
					}
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
	var counter uint64
	var histogram = NewHist()
	var initTimestamp = time.Now()
	var lifetimeRate float64
	var osPageSize = os.Getpagesize()
	var newPIDCounter uint64
	var samples = make([]float64, window)
	var times CPUTimes
	var watching *ProcInfo

	for {
		select {
		case watching = <-p:
			if watching != nil {
				if counter > 0 && watching.PIDChaged {
					newPIDCounter++
					log.Printf(
						"Resume monitor for %s with new PID: %d *ProcInfo: %p",
						watching.Role, watching.PID, watching)
				} else {
					log.Printf("Start monitor for %s with PID: %d *ProcInfo: %p",
						watching.Role, watching.PID, watching)
				}
			} else {
				// When we receive a nil value on this channel, we know that
				// this goroutine is expected to go away now because the
				// monitored process was removed by system from process table.
				log.Println("Shutting down monitor goroutine")
				return
			}

		default:
			var s ProcStat
			var ok bool
			if s, ok = watching.Stat(); ok {
				lifetimeRate = float64(s.OnCPUTimeTotal()) / float64(watching.ProcAgeAsTicks())
				samples[counter%window] = lifetimeRate

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
					times.PrevOnCPUTime = times.CurrentOnCPUTime
					times.PrevRunTime = times.CurrentRunTime
					times.CurrentOnCPUTime = s.OnCPUTimeTotal()
					times.CurrentRunTime = watching.ProcAgeAsTicks()
					histogram.Insert(times.Delta())
				}
			} else {
				samples[counter%window] = math.NaN()
				times.Reset()
			}
			counter++
			if counter >= window {
				r <- &IntervalReport{
					PID:             watching.PID,
					Role:            watching.Role,
					InitTimestamp:   initTimestamp,
					Timestamp:       time.Now(),
					Age:             watching.ProcAgeAsDuration(),
					WindowRate:      avg(samples),
					StandardDev:     stddev(samples),
					LifetimeRate:    lifetimeRate,
					CurrentRate:     times.Delta(),
					RateHistogram:   histogram.JSONSafeMap(),
					TimesRestated:   newPIDCounter,
					VirtMemoryBytes: s.VSize,
					RSSBytes:        s.RSS * osPageSize,
				}
				// fmt.Printf("counter: %d | avg: %f\n", counter, avg(samples))
			}

			<-time.NewTimer(time.Second).C
		}
	}
}
