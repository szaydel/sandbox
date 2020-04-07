package main

import (
	"reflect"
	"testing"
)

func TestCPUTimes_Delta(t *testing.T) {
	type fields struct {
		PrevRunTime      int64
		CurrentRunTime   int64
		PrevOnCPUTime    int64
		CurrentOnCPUTime int64
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{name: "Emulated CPU time Ratio",
			fields: fields{PrevRunTime: 1000, CurrentRunTime: 1100, PrevOnCPUTime: 100, CurrentOnCPUTime: 110},
			want:   10 / 100.0,
		},
		{name: "Emulated CPU time Ratio",
			fields: fields{PrevRunTime: 1000, CurrentRunTime: -1000, PrevOnCPUTime: 100, CurrentOnCPUTime: -100},
			want:   -10 / -100.0,
		},
		{name: "Emulated CPU time Ratio",
			fields: fields{PrevRunTime: 564800, CurrentRunTime: 564900, PrevOnCPUTime: 3450, CurrentOnCPUTime: 3500},
			want:   1 / 2.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := CPUTimes{
				PrevRunTime:      tt.fields.PrevRunTime,
				CurrentRunTime:   tt.fields.CurrentRunTime,
				PrevOnCPUTime:    tt.fields.PrevOnCPUTime,
				CurrentOnCPUTime: tt.fields.CurrentOnCPUTime,
			}
			if got := tr.Delta(); got != tt.want {
				t.Errorf("CPUTimes.Delta() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCPUTimes_Reset(t *testing.T) {
	type fields struct {
		PrevRunTime      int64
		CurrentRunTime   int64
		PrevOnCPUTime    int64
		CurrentOnCPUTime int64
	}
	tests := []struct {
		name   string
		fields fields
		want *CPUTimes
	}{
		{
			name: "All fields must be reset to zeros",
			fields: fields{
				PrevRunTime: 1234,
				CurrentRunTime: 1234,
				PrevOnCPUTime: 1234,
				CurrentOnCPUTime: 1234,
			},
			want: &CPUTimes{
				PrevRunTime: 0,
				CurrentRunTime: 0,
				PrevOnCPUTime: 0,
				CurrentOnCPUTime: 0,
			},
		},{
			name: "All fields must be reset to zeros",
			fields: fields{
				PrevRunTime: -1234,
				CurrentRunTime: -1234,
				PrevOnCPUTime: -1234,
				CurrentOnCPUTime: -1234,
			},
			want: &CPUTimes{
				PrevRunTime: 0,
				CurrentRunTime: 0,
				PrevOnCPUTime: 0,
				CurrentOnCPUTime: 0,
			},
		},{
			name: "All fields must be reset to zeros",
			fields: fields{
				PrevRunTime: 0,
				CurrentRunTime: 0,
				PrevOnCPUTime: 0,
				CurrentOnCPUTime: 0,
			},
			want: &CPUTimes{
				PrevRunTime: 0,
				CurrentRunTime: 0,
				PrevOnCPUTime: 0,
				CurrentOnCPUTime: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &CPUTimes{
				PrevRunTime:      tt.fields.PrevRunTime,
				CurrentRunTime:   tt.fields.CurrentRunTime,
				PrevOnCPUTime:    tt.fields.PrevOnCPUTime,
				CurrentOnCPUTime: tt.fields.CurrentOnCPUTime,
			}
			tr.Reset()
			if !reflect.DeepEqual(tr, tt.want) {
				t.Errorf("Reset() = %+v, want %+v", tr, tt.want)
			}

		})
	}
}
