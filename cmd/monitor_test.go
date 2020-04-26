package main

import (
	"context"
	"os"
	"testing"
)

func Test_startMonitors(t *testing.T) {
	type args struct {
		ctx       context.Context
		repChan   chan *IntervalReport
		processes func() []*ProcInfo
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startMonitors(tt.args.ctx, tt.args.repChan, tt.args.processes)
		})
	}
}

// Test_monitor test needs more than 5 seconds to run, therefore, set timeout
// to at least 10 seconds. Otherwise, expect to see something like:
// panic: test timed out after 5s
func Test_monitor(t *testing.T) {
	p := make(chan *ProcInfo)
	r := make(chan *IntervalReport)

	pi := findProcsByName(os.Args[0])[0]
	go monitor(p, r)
	p <- pi

	// Expect to get two interval reports before closing the *ProcInfo channel.
	t.Logf("%+v\n", <-r)
	t.Logf("%+v\n", <-r)
	close(p)
	p = nil

}
