package main

import "testing"

func TestNewSummaries(t *testing.T) {
	if metricsReport != nil {
		t.Fatal("NewSummariesSingleton(): metricsReport must be nil")
	}
	if !NewSummariesSingleton() {
		t.Errorf("NewSummariesSingleton(): expected true")
	}
	if metricsReport == nil {
		t.Errorf("NewSummariesSingleton(): metricsReport must not be nil")
	}
}
