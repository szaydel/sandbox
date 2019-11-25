package main

import (
	"sync"
)


var singleton sync.Once
// metricsReport is the only instance of Summaries struct used in the program.
// There is not really a concept of singletons in Golang, so this is next best
// approximation thereof.
var metricsReport  *Summaries

// NewSummaries returns a singleton instance of a Summaries variable, which
// is used everywhere else.
func NewSummaries() {
	singleton.Do(func() {
		metricsReport = &Summaries{
			m:   make(map[string]*IntervalReport),
			mtx: sync.RWMutex{},
		}
	})
}

func init() {
	NewSummaries()
}
