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

func variance(nums []float64) float64 {
	var μ, σ2 float64
	if len(nums) == 0 {
		return math.NaN()
	}
	μ = sum(nums) / float64(len(nums))

	for _, n := range nums {
		σ2 += (n - μ) * (n - μ)
	}
	return σ2 / float64(len(nums)-1) // -1 for sample variance
}

func stddev(nums []float64) float64 {
	var σ2 float64
	if len(nums) == 0 {
		return math.NaN()
	}
	σ2 = variance(nums)
	return math.Pow(σ2, 0.5)
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

func avg(nums []float64) float64 {
	var skip = countNaNs(nums)
	// If all values are NaNs, cnt will be zero, and as a result, dividing by
	// zero will also be a NaN value.
	var cnt = len(nums) - skip
	return sum(nums) / float64(cnt)
}
