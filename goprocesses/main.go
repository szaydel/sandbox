package main

import (
	"time"
)

// type Count int
// type Role string

// const Proxy Role = "proxy"
// const Worker Role = "worker"
// const Manager Role = "manager"
// const Logger Role = "logger"

// func expectedCount(r Role) Count {
// 	m := map[Role]Count{
// 		Proxy:   2,
// 		Worker:  1,
// 		Manager: 1,
// 		Logger:  1,
// 	}
// 	return m[r]
// }

func main() {
	// processes := findProcsByName("bro")
	// j, _ := json.Marshal(processes)
	// fmt.Printf("%v\n", string(j))

	shutdownChan := make(chan struct{})
	intervalReportChan := make(chan *IntervalReport)
	go startMonitors(shutdownChan, intervalReportChan)
	go startIntervalReport(intervalReportChan)
	time.Sleep(100 * time.Second)
	shutdownChan <- struct{}{}
}
