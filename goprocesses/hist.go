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

func (h *Histogram) Insert(n float64) {
	slot := sort.SearchFloat64s(h.slots, n)
	for i := len(h.slots) - 1; i >= slot; i-- {
		h.counts[i]++
	}
}

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

// func main() {
// 	h := NewHist()
// 	for _, n := range [...]float64{
// 		0.,
// 		0.0005,
// 		0.001,
// 		0.01,
// 		0.1,
// 		0.11,
// 		0.12,
// 		0.199,
// 		0.2,
// 		1.0002,
// 		2.0,
// 		0.45,
// 	} {
// 		h.Insert(n)
// 	}
// 	fmt.Printf("%v\n", h)

// }
