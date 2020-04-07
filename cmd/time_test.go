package main

import (
	"time"
)
import "testing"

func Test_ticksToNsecs(t *testing.T) {
	type args struct {
		ticks int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{name: "zero seconds", args: args{ticks: 0}, want: 0},
		{name: "one-100th of a second", args: args{ticks: 1}, want: 1e7},
		{name: "one tenth of a second", args: args{ticks: 10}, want: 1e8},
		{name: "one second", args: args{ticks: 100}, want: 1e9},
		{name: "ten seconds", args: args{ticks: 1000}, want: 1e10},
		{name: "100 seconds", args: args{ticks: 10000}, want: 1e11},
		{name: "1000 seconds", args: args{ticks: 100000}, want: 1e12},
		{name: "10000 seconds", args: args{ticks: 1000000}, want: 1e13},
		{name: "100000 seconds", args: args{ticks: 10000000}, want: 1e14},
		{name: "1000000 seconds", args: args{ticks: 100000000}, want: 1e15},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ticksToNsecs(tt.args.ticks); got != tt.want {
				t.Errorf("ticksToNsecs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_monotonicClockTicks(t *testing.T) {
	const c = 100
	var samples = make([]int64, c)
	for i := 0; i < c; i++ {
		samples[i] = monotonicClockTicks()
	}
	if !timeIsIncreasingInt64(samples) {
		t.Errorf("monotonicClockTicks() expecting only increasing values")
	}
}

func Test_monotonicSinceBoot(t *testing.T) {
	const c = 100
	var samples = make([]time.Duration, c)
	for i := 0; i < c; i++ {
		samples[i] = monotonicSinceBoot()
	}
	if !timeIsIncreasing(samples) {
		t.Errorf("monotonicSinceBoot() expecting only increasing values")
	}
}

func timeIsIncreasingInt64(samples []int64) bool {
	for i := 0; i < len(samples)-1; i++ {
		if samples[i] > samples[i+1] { // next sample must be > previous sample
			return false
		}
	}
	return true
}

func timeIsIncreasing(samples []time.Duration) bool {
	for i := 0; i < len(samples)-1; i++ {
		if samples[i] > samples[i+1] { // next sample must be > previous sample
			return false
		}
	}
	return true
}
