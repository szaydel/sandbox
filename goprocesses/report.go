package main

import (
	"fmt"
	"time"
)

// IntervalReport is a point in time view of process' CPU usage with three
// figures, WindowRate, LifeTimeRate and CurrentRate.
// WindowRate - an average of samples over several intervals, which effectively
// makes the data smoother.
// LifeTimeRate - rate of time spent on CPU over total process' runtime,
// computed over the entire lifetime of process; least volatile.
// CurrentRate - derivative between two interval samples; most volatile.
type IntervalReport struct {
	PID          int
	Role         string
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
