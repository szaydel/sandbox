package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	setupCliFlags() // Setup and parse command line flags
	// trap Ctrl+C and call cancel on the context
	ctx, cancel := context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)
	signal.Notify(sigChan, syscall.SIGTERM)
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

	if !NewSummariesSingleton() {
		panic("failed to initialize Summaries structure")
	}
	intervalReportChan := make(chan *IntervalReport)
	go startMonitors(
		ctx,
		intervalReportChan,
		func() []*ProcInfo {
			return findProcsByName(exeLocation)
		})
	go startIntervalReport(intervalReportChan)

	http.HandleFunc("/info", allInfoHandler)
	http.HandleFunc("/info/", roleInfoHandler) // children of /info route
	http.HandleFunc("/metrics", prometheusMetricsHandler) // prometheus output

	log.Printf("Server starting on %s:%d\n", hostname, port)
	go http.ListenAndServe(fmt.Sprintf("%s:%d", hostname, port), nil)

	<-ctx.Done()
}
