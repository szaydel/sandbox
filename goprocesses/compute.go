package main

import (
	"math"
)

func sum(nums []float64) float64 {
	var total float64
	for _, n := range nums {
		if !math.IsNaN(n) {
			total += n
		}
	}
	return total
}

func avg(nums []float64) float64 {
	var skip = countNaNs(nums)
	// If all values are NaNs, cnt will be zero, and as a result, dividing by
	// zero will also be a NaN value.
	var cnt = len(nums) - skip
	return sum(nums) / float64(cnt)
}

func countNaNs(nums []float64) int {
	if len(nums) == 0 {
		return 0
	}
	if len(nums) == 1 {
		if math.IsNaN(nums[0]) {
			return 1
		}
		return 0
	}
	if math.IsNaN(nums[0]) {
		return 1 + countNaNs(nums[1:])
	}
	return countNaNs(nums[1:])
}
