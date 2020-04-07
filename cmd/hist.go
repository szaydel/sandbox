package main

import (
	"math"
	"sort"
)

// Histogram is a structure that implements a very basic histogram function,
// which acts as a cumulative histogram, meaning each bucket is a count of
// values which fell into all lower buckets in addition to values falling it.
// In other words, last bucket, which is +Inf is a count of every other bucket
// because there is no value larger than +Inf.
type Histogram struct {
	slots  []float64
	counts []int64
}

// NewHist returns a new histogram structure ready for use, initialized with
// slots that should cover our entire range of values.
// This histogram is cumulative, which means every bucket is a count of
// everything that is equal to the bucket's size or smaller.
// In other words every observation falls into the last bucket, because every
// observation is going to be less than +Inf.
func NewHist() *Histogram {
	return &Histogram{
		slots:  []float64{0.0001, 0.001, 0.01, 0.1, 0.2, 0.4, 0.8, math.Inf(+1)},
		counts: make([]int64, 8),
	}

}

// Insert is the primary method of the histogram used to insert new observations
// into the histogram.
func (h *Histogram) Insert(n float64) {
	slot := sort.SearchFloat64s(h.slots, n)
	for i := len(h.slots) - 1; i >= slot; i-- {
		h.counts[i]++
	}
}

// Map returns contents of a histogram as float64->int64 map.
func (h *Histogram) Map() map[float64]int64 {
	return map[float64]int64{
		0.0001:       h.counts[0],
		0.001:        h.counts[1],
		0.01:         h.counts[2],
		0.1:          h.counts[3],
		0.2:          h.counts[4],
		0.4:          h.counts[5],
		0.8:          h.counts[6],
		math.Inf(+1): h.counts[7],
	}
}

// JSONSafeMap returns contents of a histogram as string->int64 map.
// It is a representation of the histogram where keys are string values
// instead of floating point values, in order to properly serialize to JSON
// without creating a custom encoder.
func (h *Histogram) JSONSafeMap() map[string]int64 {
	return map[string]int64{
		"0.0001": h.counts[0],
		"0.001":  h.counts[1],
		"0.01":   h.counts[2],
		"0.1":    h.counts[3],
		"0.2":    h.counts[4],
		"0.4":    h.counts[5],
		"0.8":    h.counts[6],
		"+Inf":   h.counts[7],
	}
}
