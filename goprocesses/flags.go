package main

import "flag"

func setupCliFlags() {
	flag.StringVar(&exeLocation, "exeLocation", "/workspace/sandbox/bin/bro", "Path to executable to be monitored")
	flag.IntVar(&port, "port", defaultPort, "Listen on this port")
	flag.DurationVar(&reportInterval, "report-interval", defaultReportInterval, "Print summaries for all monitored processes with this interval")
	flag.StringVar(&hostname, "hostname", defaultHostname, "Address on which to listen")
	flag.Uint64Var(&windowSize, "window-size", defaultWindowSize, "Length of samples window over which statistics are calculated; the larger the number, the smoother the data")
	flag.Parse()
}