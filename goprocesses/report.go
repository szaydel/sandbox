package main

import (
	"errors"
	"time"
	"sync"
	"encoding/json"
)

// IntervalReport is a point in time view of process' CPU usage with three
// figures, WindowRate, LifeTimeRate and CurrentRate.
// WindowRate - an average of samples over several intervals, which effectively
// makes the data smoother.
// StandardDev - standard deviation for samples in this window.
// LifeTimeRate - rate of time spent on CPU over total process' runtime,
// computed over the entire lifetime of process; least volatile.
// CurrentRate - derivative between two interval samples; most volatile.
type IntervalReport struct {
	PID             int `json:"pid"` 
	Role            string `json:"role"` 
	Timestamp       time.Time `json:"timestamp"` 
	WindowRate      float64 `json:"window_rate"` 
	StandardDev     float64 `json:"standard_dev"` 
	LifetimeRate    float64 `json:"lifetime_rate"` 
	CurrentRate     float64 `json:"current_rate"` 
	VirtMemoryBytes uint `json:"virtual_memory_bytes"` 
	RSSBytes        int `json:"rss_bytes"` 
}

func startIntervalReport(c <-chan *IntervalReport) {
	for {
		select {
		case v := <-c:
			intervalReportMap.Insert(v)
			//fmt.Printf("%+v | %p\n", v, v)
		default:
			<-time.NewTimer(1 * time.Second).C
		}
	}
}

// Summaries is used as a global singleton to keep track of running
// statistics for processes being monitored.
type Summaries struct {
	m map[string]*IntervalReport
	mtx sync.Mutex
}

// Insert updates the map with latest interval report.
func (s *Summaries) Insert(r *IntervalReport) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.m[r.Role] = r
}

// PairToJSON returns serialized version of a ringle entry from interval 
// summaries map, assuming entry is found in the map.
func (s *Summaries) PairToJSON(role string) ([]byte, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if v, ok := s.m[role] ; ok {
	return json.Marshal(v)
	}
	return []byte{}, errors.New("role does not exist in the map")
}

// ToJSON returns serialized version of the interval summaries map.
func (s *Summaries) ToJSON() ([]byte, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	return json.Marshal(s.m)
}