package main
import (
	"sync"
)

var intervalReportMap = Summaries{
	m: make(map[string]*IntervalReport),
	mtx: sync.Mutex{},
}