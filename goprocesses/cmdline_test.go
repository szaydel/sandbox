package main

import (
	"reflect"
	"testing"
)

func TestCommandLine_ProgramName(t *testing.T) {
	type fields struct {
		args        []string
		programName string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := CommandLine{
				args:        tt.fields.args,
				programName: tt.fields.programName,
			}
			if got := c.ProgramName(); got != tt.want {
				t.Errorf("CommandLine.ProgramName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommandLine_Args(t *testing.T) {
	type fields struct {
		args        []string
		programName string
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := CommandLine{
				args:        tt.fields.args,
				programName: tt.fields.programName,
			}
			if got := c.Args(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CommandLine.Args() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cmdLineArgs(t *testing.T) {
	type args struct {
		pid int
	}
	tests := []struct {
		name string
		args args
		want *CommandLine
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cmdLineArgs(tt.args.pid); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cmdLineArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}
