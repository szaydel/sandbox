package main

import (
	"reflect"
	"testing"
	"time"
)

func TestReadFileNoStat(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadFileNoStat(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadFileNoStat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadFileNoStat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildProcInfo(t *testing.T) {
	type args struct {
		procfile string
	}
	tests := []struct {
		name string
		args args
		want *ProcInfo
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildProcInfo(tt.args.procfile); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildProcInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_findProcsByName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want []*ProcInfo
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findProcsByName(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findProcsByName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProcStat_OnCPUTimeTotal(t *testing.T) {
	type fields struct {
		PID        int
		Comm       string
		State      string
		PPID       int
		PGRP       int
		Session    int
		TTY        int
		TPGID      int
		Flags      uint
		MinFlt     uint
		CMinFlt    uint
		MajFlt     uint
		CMajFlt    uint
		UTime      uint
		STime      uint
		CUTime     uint
		CSTime     uint
		Priority   int
		Nice       int
		NumThreads int
		Starttime  uint64
		VSize      uint
		RSS        int
	}
	tests := []struct {
		name   string
		fields fields
		want   int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := ProcStat{
				PID:        tt.fields.PID,
				Comm:       tt.fields.Comm,
				State:      tt.fields.State,
				PPID:       tt.fields.PPID,
				PGRP:       tt.fields.PGRP,
				Session:    tt.fields.Session,
				TTY:        tt.fields.TTY,
				TPGID:      tt.fields.TPGID,
				Flags:      tt.fields.Flags,
				MinFlt:     tt.fields.MinFlt,
				CMinFlt:    tt.fields.CMinFlt,
				MajFlt:     tt.fields.MajFlt,
				CMajFlt:    tt.fields.CMajFlt,
				UTime:      tt.fields.UTime,
				STime:      tt.fields.STime,
				CUTime:     tt.fields.CUTime,
				CSTime:     tt.fields.CSTime,
				Priority:   tt.fields.Priority,
				Nice:       tt.fields.Nice,
				NumThreads: tt.fields.NumThreads,
				Starttime:  tt.fields.Starttime,
				VSize:      tt.fields.VSize,
				RSS:        tt.fields.RSS,
			}
			if got := ps.OnCPUTimeTotal(); got != tt.want {
				t.Errorf("ProcStat.OnCPUTimeTotal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProcInfo_Stat(t *testing.T) {
	type fields struct {
		Name        string
		Role        string
		Args        []string
		PID         int
		PIDChaged   bool
		AgeTicks    int64
		AgeDuration time.Duration
		S           *ProcStat
	}
	tests := []struct {
		name   string
		fields fields
		want   ProcStat
		want1  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := ProcInfo{
				Name:        tt.fields.Name,
				Role:        tt.fields.Role,
				Args:        tt.fields.Args,
				PID:         tt.fields.PID,
				PIDChaged:   tt.fields.PIDChaged,
				AgeTicks:    tt.fields.AgeTicks,
				AgeDuration: tt.fields.AgeDuration,
				S:           tt.fields.S,
			}
			got, got1 := p.Stat()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProcInfo.Stat() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ProcInfo.Stat() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestProcInfo_path(t *testing.T) {
	type fields struct {
		Name        string
		Role        string
		Args        []string
		PID         int
		PIDChaged   bool
		AgeTicks    int64
		AgeDuration time.Duration
		S           *ProcStat
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := ProcInfo{
				Name:        tt.fields.Name,
				Role:        tt.fields.Role,
				Args:        tt.fields.Args,
				PID:         tt.fields.PID,
				PIDChaged:   tt.fields.PIDChaged,
				AgeTicks:    tt.fields.AgeTicks,
				AgeDuration: tt.fields.AgeDuration,
				S:           tt.fields.S,
			}
			if got := p.path(tt.args.name); got != tt.want {
				t.Errorf("ProcInfo.path() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProcInfo_ProcAgeAsTicks(t *testing.T) {
	type fields struct {
		Name        string
		Role        string
		Args        []string
		PID         int
		PIDChaged   bool
		AgeTicks    int64
		AgeDuration time.Duration
		S           *ProcStat
	}
	tests := []struct {
		name   string
		fields fields
		want   int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := ProcInfo{
				Name:        tt.fields.Name,
				Role:        tt.fields.Role,
				Args:        tt.fields.Args,
				PID:         tt.fields.PID,
				PIDChaged:   tt.fields.PIDChaged,
				AgeTicks:    tt.fields.AgeTicks,
				AgeDuration: tt.fields.AgeDuration,
				S:           tt.fields.S,
			}
			if got := p.ProcAgeAsTicks(); got != tt.want {
				t.Errorf("ProcInfo.ProcAgeAsTicks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProcInfo_ProcAgeAsDuration(t *testing.T) {
	type fields struct {
		Name        string
		Role        string
		Args        []string
		PID         int
		PIDChaged   bool
		AgeTicks    int64
		AgeDuration time.Duration
		S           *ProcStat
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Duration
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := ProcInfo{
				Name:        tt.fields.Name,
				Role:        tt.fields.Role,
				Args:        tt.fields.Args,
				PID:         tt.fields.PID,
				PIDChaged:   tt.fields.PIDChaged,
				AgeTicks:    tt.fields.AgeTicks,
				AgeDuration: tt.fields.AgeDuration,
				S:           tt.fields.S,
			}
			if got := p.ProcAgeAsDuration(); got != tt.want {
				t.Errorf("ProcInfo.ProcAgeAsDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}
