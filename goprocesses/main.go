package main

import (
	"context"
	"os"
	"os/signal"
	"net/http"
	"log"
	"fmt"
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
	// trap Ctrl+C and call cancel on the context
	ctx, cancel := context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	defer func() {
		signal.Stop(sigChan)
		cancel()
	}()
	go func() {
		select {
		case <-sigChan:
			cancel()
		case <-ctx.Done():
		}
	}()

	intervalReportChan := make(chan *IntervalReport)
	go startMonitors(
		ctx,
		intervalReportChan,
		func() []*ProcInfo {
			return findProcsByName("/workspace/sandbox/bin/bro")
		})
	go startIntervalReport(intervalReportChan)


	port := 8080
	http.HandleFunc("/info", handler)

	log.Printf("Server starting on port %v\n", port)
	go http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
	
	<-ctx.Done()
}

func handler(w http.ResponseWriter, r *http.Request) {
	data, err := intervalReportMap.ToJSON()
	handleErr(err, true)
	fmt.Fprint(w, string(data))
}