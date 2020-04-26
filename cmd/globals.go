package main

import (
	"sync"
	"time"
)

const defaultHostname = "localhost"
const defaultPort = 8080
const defaultReportInterval = time.Second * 5

// defaultWindowSize is the number of samples for statistical functions like
// average, standard deviation, etc. The larger the window the smoother the data
// is going to appear, becasuse extreme observations play a lesser role as the
// window size increases.
const defaultWindowSize = 10

var exeLocation string
var hostname string
var port int
var windowSize uint64
var reportInterval time.Duration

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
