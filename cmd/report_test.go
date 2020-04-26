package main

import (
	"reflect"
	"sync"
	"testing"
	"time"
)

func Test_startIntervalReport(t *testing.T) {
	type args struct {
		c chan *IntervalReport
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{name: "A", args: args{c: make(chan *IntervalReport)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reportInterval = 1 * time.Second
			NewSummariesSingleton()
			// TODO: Read stdout and make sure we get data we expect
			go startIntervalReport(tt.args.c)
			tt.args.c <- &IntervalReport{}
			tt.args.c <- nil
		})
	}
}

func TestSummaries_Insert(t *testing.T) {
	type fields struct {
		m   map[string]*IntervalReport
		mtx sync.RWMutex
	}
	type args struct {
		r *IntervalReport
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Summaries{
				m:   tt.fields.m,
				mtx: tt.fields.mtx,
			}
			s.Insert(tt.args.r)
		})
	}
}

func TestSummaries_Len(t *testing.T) {
	type fields struct {
		m   map[string]*IntervalReport
		mtx sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Summaries{
				m:   tt.fields.m,
				mtx: tt.fields.mtx,
			}
			if got := s.Len(); got != tt.want {
				t.Errorf("Summaries.Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSummaries_findRole(t *testing.T) {
	type fields struct {
		m   map[string]*IntervalReport
		mtx sync.RWMutex
	}
	type args struct {
		role string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *IntervalReport
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Summaries{
				m:   tt.fields.m,
				mtx: tt.fields.mtx,
			}
			if got := s.findRole(tt.args.role); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Summaries.findRole() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSummaries_safeIntervalReport(t *testing.T) {
	type fields struct {
		m   map[string]*IntervalReport
		mtx sync.RWMutex
	}
	type args struct {
		role string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *IntervalReport
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Summaries{
				m:   tt.fields.m,
				mtx: tt.fields.mtx,
			}
			if got := s.safeIntervalReport(tt.args.role); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Summaries.safeIntervalReport() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSummaries_RoleToJSON(t *testing.T) {
	type fields struct {
		m   map[string]*IntervalReport
		mtx sync.RWMutex
	}
	type args struct {
		role string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Summaries{
				m:   tt.fields.m,
				mtx: tt.fields.mtx,
			}
			got, err := s.RoleToJSON(tt.args.role)
			if (err != nil) != tt.wantErr {
				t.Errorf("Summaries.RoleToJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Summaries.RoleToJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSummaries_ToJSON(t *testing.T) {
	type fields struct {
		m   map[string]*IntervalReport
		mtx sync.RWMutex
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Summaries{
				m:   tt.fields.m,
				mtx: tt.fields.mtx,
			}
			got, err := s.ToJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Summaries.ToJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Summaries.ToJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}
