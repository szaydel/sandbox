package main

import (
	"reflect"
	"testing"
)

func Test_handleErr(t *testing.T) {
	type args struct {
		e       error
		doPanic bool
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handleErr(tt.args.e, tt.args.doPanic)
		})
	}
}

func Test_nullByteToSpace(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := nullByteToSpace(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("nullByteToSpace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isTargetProcess(t *testing.T) {
	type args struct {
		pid    int
		target string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isTargetProcess(tt.args.pid, tt.args.target); got != tt.want {
				t.Errorf("isTargetProcess() = %v, want %v", got, tt.want)
			}
		})
	}
}
