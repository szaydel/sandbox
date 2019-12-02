package main

import (
	"sync"
)

var singleton sync.Once

// metricsReport is the only instance of Summaries struct used in the program.
// There is not really a concept of singletons in Golang, so this is next best
// approximation thereof.
var metricsReport *Summaries

// NewSummariesSingleton returns a singleton instance of a Summaries global
// variable, which is used everywhere else.
func NewSummariesSingleton() bool {
	var initialized bool
	singleton.Do(func() {
		metricsReport = &Summaries{
			m:   make(map[string]*IntervalReport),
			mtx: sync.RWMutex{},
		}
		initialized = true
	})
	return initialized
}
