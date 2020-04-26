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

// MonitoredProcesses is a structure which encapsulates data and associated
// operations in the startMonitors function. It is assumed to not be thread-safe
// because there is only ever one instance of this structure around, because
// there is only ever a single instance of the startMonitors in existence. If
// this design changes in the future, locking may become necessary to protect
// concurrent modifications to maps.
type MonitoredProcesses struct {
	transient  map[string]struct{}
	persistent map[string]int
	notSeen    map[string]int
	channels   map[string]chan *ProcInfo
}

// NewTransient creates a new transient map.
func (mp *MonitoredProcesses) NewTransient() map[string]struct{} {
	mp.transient = make(map[string]struct{})
	return mp.transient
}

// InsertIntoTransient registers a role by making sure a key with role's name
// is present. This is used for the purposes of comparison between what we
// believe is current and what is actually current in the process table on the
// system.
func (mp *MonitoredProcesses) InsertIntoTransient(role string) {
	mp.transient[role] = struct{}{}
}

// RegisterNewProcess creates and initializes all the necessary pieces before
// we can start a new monitor thread.
func (mp *MonitoredProcesses) RegisterNewProcess(role string, pid int) {
	mp.persistent[role] = pid
	mp.notSeen[role] = 0
	mp.channels[role] = make(chan *ProcInfo)
}

func (mp *MonitoredProcesses) pidOf(role string) (int, bool) {
	if pid, ok := mp.persistent[role]; ok {
		return pid, true
	}
	return -1, false
}

// PidChanged reports whether or not a PID changed for process being monitored.
// If a process if not already known an error is returned with false, otherwise
// nil error and true are returned.
func (mp *MonitoredProcesses) PidChanged(role string, pid int) (bool, error) {
	if lastKnownPid, ok := mp.pidOf(role); ok {
		return pid != lastKnownPid, nil
	}
	return false, errors.New("this role not seen previously")
}

// UpdatePid replaces PID for a given process with another value.
func (mp *MonitoredProcesses) UpdatePid(role string, pid int) {
	mp.persistent[role] = pid
}

// ResetNotSeen resets to zero a counter value tracking number of times
// a process known to have existed at some point is no longer seen.
func (mp *MonitoredProcesses) ResetNotSeen(role string) {
	mp.notSeen[role] = 0
}

// IncrNotSeen increments by one count of times a process known to have existed
// at some point, which no longer appears to exist.
func (mp *MonitoredProcesses) IncrNotSeen(role string) {
	if _, ok := mp.notSeen[role]; !ok {
		mp.notSeen[role] = 1
		return
	}
	mp.notSeen[role]++
}

func (mp *MonitoredProcesses) closeChannels(role string) {
	close(mp.channels[role])
	mp.channels[role] = nil
}

// RemoveMonitored is effectively a reverse of RegisterNewProcess. It removes
// structures in memory associated with a process which we no longer want to
// monitor.
func (mp *MonitoredProcesses) RemoveMonitored(role string) {
	mp.closeChannels(role)
	delete(mp.channels, role)
	delete(mp.persistent, role)
	delete(mp.notSeen, role)
	// We do not do anything about the transient map because it is recreated
	// frequently unlike the other two maps.
}

// NotSeenCount returns count of times a process known to have existed
// previously was not seen.
func (mp *MonitoredProcesses) NotSeenCount(role string) int {
	if count, ok := mp.notSeen[role]; ok {
		return count
	}
	mp.notSeen[role] = 0 // If for some reason it is not in map, add to map.
	return 0
}

// NotSeen returns the map of processes and their "not seen" counts.
func (mp *MonitoredProcesses) NotSeen() map[string]int {
	return mp.notSeen
}

// NotSeenFewerThan just returns a boolean value resulting from a comparison
// of number of times not seen and value passed in via the count argument.
func (mp *MonitoredProcesses) NotSeenFewerThan(role string, count int) bool {
	return mp.notSeen[role] < count
}

// Persistent retuns the persistent map of role->PIDs.
func (mp *MonitoredProcesses) Persistent() map[string]int {
	return mp.persistent
}

// Transient retuns the transient map of role->struct{}(s).
func (mp *MonitoredProcesses) Transient() map[string]struct{} {
	return mp.transient
}

// InTransient does a membership check, reporting whether or not a role is
// present in the transient map.
func (mp *MonitoredProcesses) InTransient(role string) bool {
	_, ok := mp.transient[role]
	return ok
}

// NewMonitoredProcesses returns an initialized and ready to go structure used
// by startMonitors function. It is a convenience mechanism mainly to simplify
// creating all required maps.
func NewMonitoredProcesses() *MonitoredProcesses {
	return &MonitoredProcesses{
		transient:  make(map[string]struct{}),
		persistent: make(map[string]int),
		notSeen:    make(map[string]int),
		channels:   make(map[string]chan *ProcInfo),
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
	var mp = NewMonitoredProcesses()

	for {
		select {
		case <-ctx.Done():
			return
			// shutdown all monitor processes
		default:
			// We want to create this map each time through this loop. This map
			// is expected to be transient and its contents are only good for a
			// single iteration of this loop.
			mp.NewTransient()
			for _, p := range processes() {
				mp.InsertIntoTransient(p.Role)
				if changed, err := mp.PidChanged(p.Role, p.PID); err == nil {
					if changed {
						if pid, ok := mp.pidOf(p.Role); ok {
							log.Printf("PID for process %s changed from %d to %d",
								p.Role, pid, p.PID)
						}
						mp.UpdatePid(p.Role, p.PID)
						p.PIDChaged = true
						mp.channels[p.Role] <- p
					}
				} else { // This process' Role is not already in the map
					// Register a new process if this is a process which we have
					// never seen before and do not already have a monitor
					// thread for its role. If this is a new process ID for a
					// previously seen role, we will already have these, and
					// instead we just update the process ID above.
					mp.RegisterNewProcess(p.Role, p.PID)

					// Start a new monitor thread for role we have not yet seen,
					// or have seen before but removed because it was not seen
					// for a number of intervals.
					go monitor(mp.channels[p.Role], repChan)
					mp.channels[p.Role] <- p
					log.Printf("Added %s with PID %d to map", p.Role, p.PID)
				}
			}

			for role, pid := range mp.Persistent() {
				if !mp.InTransient(role) {
					if mp.NotSeenFewerThan(role, MaxNotSeenIntervals) {
						mp.IncrNotSeen(role)
						log.Printf("PID %d for process %s no longer seen", pid, role)
						continue
					}
					// We need to notify corresponding goroutine that it needs
					// to shutdown! After telling relevant goroutine to stop,
					// remove the no longer existing role from the current map.
					log.Printf("Removing %s from list of monitored processes", role)
					// RemoveMonitored(...) signals associated goroutine to
					// stop and return, otherwise we are going to have leaking
					// goroutines.
					mp.RemoveMonitored(role)
					// Do not attempt to send on the channel for the process
					// after calling RemoveMonitored(...) here to prevent a
					// send on closed channel panic.
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
	var counter uint64
	var histogram = NewHist()
	var initTimestamp = time.Now()
	var lifetimeRate float64
	var osPageSize = os.Getpagesize()
	var newPIDCounter uint64
	var times CPUTimes
	var watching *ProcInfo
	var window = windowSize
	var samples = make([]float64, window)

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
