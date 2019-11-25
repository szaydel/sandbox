package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
)

const Port = 8080

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

	http.HandleFunc("/info", allInfoHandler)
	http.HandleFunc("/info/", roleInfoHandler) // children of /info route

	log.Printf("Server starting on port %v\n", Port)
	go http.ListenAndServe(fmt.Sprintf(":%v", Port), nil)

	<-ctx.Done()
}

func allInfoHandler(w http.ResponseWriter, r *http.Request) {
	data, err := metricsReport.ToJSON()
	if err != nil {
		handleErr(err, false)
		http.Error(w,
			http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError,
		)
		return
	}
	fmt.Fprint(w, string(data))
}

func roleInfoHandler(w http.ResponseWriter, r *http.Request) {
	role := strings.TrimPrefix(r.URL.Path, "/info/")
	//fmt.Printf("URL.Path: %+v\n", role)
	//fmt.Printf("URL.RequestURI(): %+v\n", r.URL.RequestURI())
	data, err := metricsReport.RoleToJSON(role)
	if err != nil {
		handleErr(
			fmt.Errorf(
				"failed getting metrics for %s with: %s", role, err),
			false,
		)
		if err == errNoInfoForRole {
			http.NotFound(w, r)
		} else {
			http.Error(w,
				http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError,
			)
		}
		return
	}
	fmt.Fprint(w, string(data))
}
