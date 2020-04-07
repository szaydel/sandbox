package main

import (
	"math"
	"math/rand"
	"reflect"
	"testing"
)

func TestNewHist(t *testing.T) {
	tests := []struct {
		name string
		want *Histogram
	}{
		{name: "New Empty Histogram", want: &Histogram{slots: []float64{0.0001, 0.001, 0.01, 0.1, 0.2, 0.4, 0.8, math.Inf(+1)}, counts: []int64{0, 0, 0, 0, 0, 0, 0, 0}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHist(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHistogram_Insert(t *testing.T) {
	type fields struct {
		slots  []float64
		counts []int64
	}
	type args struct {
		fn func(h *Histogram) *Histogram
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Histogram
	}{
		{name: "10K Random variates from seed 2^13 - 1",
			fields: fields{
				slots:  []float64{0.0001, 0.001, 0.01, 0.1, 0.2, 0.4, 0.8, math.Inf(+1)},
				counts: []int64{0, 0, 0, 0, 0, 0, 0, 0},
			},
			args: args{
				fn: func(h *Histogram) *Histogram {
					rand.Seed(981265)
					for i := 0; i < 10000; i++ {
						h.Insert(rand.Float64())
					}
					return h
				}},
			want: &Histogram{
				slots:  []float64{0.0001, 0.001, 0.01, 0.1, 0.2, 0.4, 0.8, math.Inf(+1)},
				counts: []int64{4, 19, 117, 1026, 2001, 3966, 8017, 10000},
			},
		},
		{name: "10K Random variates from seed 2^19 - 1",
			fields: fields{
				slots:  []float64{0.0001, 0.001, 0.01, 0.1, 0.2, 0.4, 0.8, math.Inf(+1)},
				counts: []int64{0, 0, 0, 0, 0, 0, 0, 0},
			},
			args: args{
				fn: func(h *Histogram) *Histogram {
					rand.Seed(524287)
					for i := 0; i < 10000; i++ {
						h.Insert(rand.Float64())
					}
					return h
				}},
			want: &Histogram{
				slots:  []float64{0.0001, 0.001, 0.01, 0.1, 0.2, 0.4, 0.8, math.Inf(+1)},
				counts: []int64{2, 12, 99, 1052, 2056, 4047, 8019, 10000},
			},
		},
		{name: "10K Random variates from seed 2^31 - 1",
			fields: fields{
				slots:  []float64{0.0001, 0.001, 0.01, 0.1, 0.2, 0.4, 0.8, math.Inf(+1)},
				counts: []int64{0, 7, 102, 958, 1899, 3882, 7994, 10000},
			},
			args: args{
				fn: func(h *Histogram) *Histogram {
					rand.Seed(2147483647)
					for i := 0; i < 10000; i++ {
						h.Insert(rand.Float64())
					}
					return h
				}},
			want: &Histogram{
				slots:  []float64{0.0001, 0.001, 0.01, 0.1, 0.2, 0.4, 0.8, math.Inf(+1)},
				counts: []int64{1, 15, 188, 1917, 3860, 7796, 15862, 20000},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Histogram{
				slots:  tt.fields.slots,
				counts: tt.fields.counts,
			}
			if got := tt.args.fn(h); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHist() = %v, want %v", got, tt.want)
			}

			t.Logf("%+v\n", h)
		})
	}
}

func TestHistogram_Map(t *testing.T) {
	type fields struct {
		slots  []float64
		counts []int64
	}
	tests := []struct {
		name   string
		fields fields
		want   map[float64]int64
	}{
		{name: "10K Random variates from seed 2^13 - 1",
			fields: fields{
				slots:  []float64{0.0001, 0.001, 0.01, 0.1, 0.2, 0.4, 0.8, math.Inf(+1)},
				counts: []int64{4, 19, 117, 1026, 2001, 3966, 8017, 10000},
			},
			want: map[float64]int64{0.0001: 4, 0.001: 19, 0.01: 117, 0.1: 1026, 0.2: 2001, 0.4: 3966, 0.8: 8017, math.Inf(+1): 10000},
		},
		{name: "10K Random variates from seed 2^19 - 1",
			fields: fields{
				slots:  []float64{0.0001, 0.001, 0.01, 0.1, 0.2, 0.4, 0.8, math.Inf(+1)},
				counts: []int64{2, 12, 99, 1052, 2056, 4047, 8019, 10000},
			},
			want: map[float64]int64{0.0001: 2, 0.001: 12, 0.01: 99, 0.1: 1052, 0.2: 2056, 0.4: 4047, 0.8: 8019, math.Inf(+1): 10000},
		},
		{name: "10K Random variates from seed 2^31 - 1",
			fields: fields{
				slots:  []float64{0.0001, 0.001, 0.01, 0.1, 0.2, 0.4, 0.8, math.Inf(+1)},
				counts: []int64{4, 19, 117, 1026, 2001, 3966, 8017, 10000},
			},
			want: map[float64]int64{0.0001: 4, 0.001: 19, 0.01: 117, 0.1: 1026, 0.2: 2001, 0.4: 3966, 0.8: 8017, math.Inf(+1): 10000},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Histogram{
				slots:  tt.fields.slots,
				counts: tt.fields.counts,
			}
			if got := h.Map(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Histogram.Map() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHistogram_JSONSafeMap(t *testing.T) {
	type fields struct {
		slots  []float64
		counts []int64
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]int64
	}{
		{name: "10K Random variates from seed 2^13 - 1",
			fields: fields{
				slots:  []float64{0.0001, 0.001, 0.01, 0.1, 0.2, 0.4, 0.8, math.Inf(+1)},
				counts: []int64{4, 19, 117, 1026, 2001, 3966, 8017, 10000},
			},
			want: map[string]int64{"0.0001": 4, "0.001": 19, "0.01": 117, "0.1": 1026, "0.2": 2001, "0.4": 3966, "0.8": 8017, "+Inf": 10000},
		},
		{name: "10K Random variates from seed 2^19 - 1",
			fields: fields{
				slots:  []float64{0.0001, 0.001, 0.01, 0.1, 0.2, 0.4, 0.8, math.Inf(+1)},
				counts: []int64{2, 12, 99, 1052, 2056, 4047, 8019, 10000},
			},
			want: map[string]int64{"0.0001": 2, "0.001": 12, "0.01": 99, "0.1": 1052, "0.2": 2056, "0.4": 4047, "0.8": 8019, "+Inf": 10000},
		},
		{name: "10K Random variates from seed 2^31 - 1",
			fields: fields{
				slots:  []float64{0.0001, 0.001, 0.01, 0.1, 0.2, 0.4, 0.8, math.Inf(+1)},
				counts: []int64{4, 19, 117, 1026, 2001, 3966, 8017, 10000},
			},
			want: map[string]int64{"0.0001": 4, "0.001": 19, "0.01": 117, "0.1": 1026, "0.2": 2001, "0.4": 3966, "0.8": 8017, "+Inf": 10000},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Histogram{
				slots:  tt.fields.slots,
				counts: tt.fields.counts,
			}
			if got := h.JSONSafeMap(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Histogram.JSONSafeMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
