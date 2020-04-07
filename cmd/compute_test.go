package main

import (
	"math"
	"testing"
)

func Test_sum(t *testing.T) {
	type args struct {
		nums []float64
	}
	tests := []struct {
		name string
		args args
		want float64
		ε    float64
	}{
		// >>> random.seed(8120)
		// >>> nums = [random.randrange(0, 100) for n in range(100)]
		{name: "100 non-negative randomly generated integers", args: args{nums: []float64{
			1, 56, 35, 63, 96, 18, 48, 87, 9, 58, 0, 38, 77, 90, 88, 4, 81, 8,
			27, 92, 33, 35, 60, 32, 79, 7, 11, 82, 95, 79, 53, 69, 0, 27, 5, 62,
			98, 43, 61, 27, 24, 64, 2, 62, 83, 85, 44, 41, 6, 82, 58, 34, 17,
			24, 48, 26, 28, 5, 20, 22, 83, 88, 36, 77, 62, 28, 54, 21, 65, 84,
			5, 93, 99, 18, 17, 31, 42, 16, 77, 58, 4, 71, 41, 28, 75, 78, 7, 90,
			61, 0, 20, 67, 33, 50, 57, 30, 92, 68, 8, 11},
		}, want: 4654, ε: 0},
		// >>> random.seed(8120)
		// >>> nums = [random.randrange(-100, 100) for n in range(100)]
		{name: "100 randomly generated integers", args: args{nums: []float64{
			-98, 13, -29, 27, 93, -63, -4, 75, -82, 16, -100, -23, 55, 80, 77,
			-92, 63, -83, -46, 85, -34, -29, 21, -36, 58, -86, -78, 64, 90, 58,
			7, 39, -99, -45, -90, 25, 97, -13, 22, -46, -51, 28, -95, 25, 67,
			70, -11, -18, -88, 64, 16, -32, -65, -52, -4, -48, -44, -90, -60,
			-55, 66, 76, -28, 54, 25, -44, 9, -57, 30, 68, -89, 86, 98, -64,
			-65, -38, -15, -68, 54, 17, -91, 43, -18, -43, 50, 57, -86, 81, 22,
			-100, -59, 35, -33, 0, 14, -40, 85, 36, -83, -77},
		}, want: -646, ε: 0},
		// >>> random.seed(8120)
		// >>> nums = [random.gauss(0, 1) for n in range(10)]
		{name: "10 randomly generated gaussian", args: args{nums: []float64{
			0.8101616449919058, 0.0548988074717373, 0.07886885577556534,
			-2.3293260605060047, -0.2810299121940123, 0.2693364335404805,
			0.8479484467222266, 0.008913715915297816, -0.4237541806367912,
			-1.4756423117462945},
		}, want: -2.4396245606658895, ε: 0.000000001},
		// >>> random.seed(8120)
		// >>> nums = [random.gauss(0, 1) for n in range(10)]
		{name: "10 randomly generated gaussian", args: args{nums: []float64{
			10.008101616449919, 10.000548988074717, 10.000788688557755,
			9.97670673939494, 9.99718970087806, 10.002693364335405,
			10.008479484467223, 10.000089137159152, 9.995762458193632,
			9.985243576882537},
		}, want: 99.97560375439335, ε: 0.000000001},
		{name: "zero length nums array", args: args{nums: []float64{}}, want: 0, ε: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sum(tt.args.nums); !tolerance(got, tt.want, tt.ε) {
				t.Errorf("sum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_variance(t *testing.T) {
	type args struct {
		nums []float64
	}
	tests := []struct {
		name string
		args args
		want float64
		ε    float64
	}{
		// >>> random.seed(8120)
		// >>> nums = [random.randrange(0, 100) for n in range(100)]
		{name: "100 non-negative randomly generated integers", args: args{nums: []float64{
			1, 56, 35, 63, 96, 18, 48, 87, 9, 58, 0, 38, 77, 90, 88, 4, 81, 8,
			27, 92, 33, 35, 60, 32, 79, 7, 11, 82, 95, 79, 53, 69, 0, 27, 5, 62,
			98, 43, 61, 27, 24, 64, 2, 62, 83, 85, 44, 41, 6, 82, 58, 34, 17,
			24, 48, 26, 28, 5, 20, 22, 83, 88, 36, 77, 62, 28, 54, 21, 65, 84,
			5, 93, 99, 18, 17, 31, 42, 16, 77, 58, 4, 71, 41, 28, 75, 78, 7, 90,
			61, 0, 20, 67, 33, 50, 57, 30, 92, 68, 8, 11},
		}, want: 905.1397979797987, ε: 0.000000001},
		// >>> random.seed(8120)
		// >>> nums = [random.randrange(-100, 100) for n in range(100)]
		{name: "100 randomly generated integers", args: args{nums: []float64{
			-98, 13, -29, 27, 93, -63, -4, 75, -82, 16, -100, -23, 55, 80, 77,
			-92, 63, -83, -46, 85, -34, -29, 21, -36, 58, -86, -78, 64, 90, 58,
			7, 39, -99, -45, -90, 25, 97, -13, 22, -46, -51, 28, -95, 25, 67,
			70, -11, -18, -88, 64, 16, -32, -65, -52, -4, -48, -44, -90, -60,
			-55, 66, 76, -28, 54, 25, -44, 9, -57, 30, 68, -89, 86, 98, -64,
			-65, -38, -15, -68, 54, 17, -91, 43, -18, -43, 50, 57, -86, 81, 22,
			-100, -59, 35, -33, 0, 14, -40, 85, 36, -83, -77},
		}, want: 3623.240808080808, ε: 0.000000001},
		// >>> random.seed(8120)
		// >>> nums = [random.gauss(0, 1) for n in range(10)]
		{name: "10 randomly generated gaussian", args: args{nums: []float64{
			0.8101616449919058, 0.0548988074717373, 0.07886885577556534,
			-2.3293260605060047, -0.2810299121940123, 0.2693364335404805,
			0.8479484467222266, 0.008913715915297816, -0.4237541806367912,
			-1.4756423117462945},
		}, want: 0.9693203277082666, ε: 0.000000001},
		// >>> random.seed(8120)
		// >>> nums = [random.gauss(0, 1) for n in range(10)]
		{name: "10 randomly generated gaussian", args: args{nums: []float64{
			10.008101616449919, 10.000548988074717, 10.000788688557755,
			9.97670673939494, 9.99718970087806, 10.002693364335405,
			10.008479484467223, 10.000089137159152, 9.995762458193632,
			9.985243576882537},
		}, want: 9.693203277082747e-05, ε: 0.000000001},
		{name: "zero length nums array", args: args{nums: []float64{}}, want: math.NaN(), ε: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := variance(tt.args.nums); !tolerance(got, tt.want, tt.ε) {
				t.Errorf("variance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stddev(t *testing.T) {
	type args struct {
		nums []float64
	}
	tests := []struct {
		name string
		args args
		want float64
		ε    float64
	}{
		// >>> random.seed(8120)
		// >>> nums = [random.randrange(0, 100) for n in range(100)]
		{name: "100 non-negative randomly generated integers", args: args{nums: []float64{
			1, 56, 35, 63, 96, 18, 48, 87, 9, 58, 0, 38, 77, 90, 88, 4, 81, 8,
			27, 92, 33, 35, 60, 32, 79, 7, 11, 82, 95, 79, 53, 69, 0, 27, 5, 62,
			98, 43, 61, 27, 24, 64, 2, 62, 83, 85, 44, 41, 6, 82, 58, 34, 17,
			24, 48, 26, 28, 5, 20, 22, 83, 88, 36, 77, 62, 28, 54, 21, 65, 84,
			5, 93, 99, 18, 17, 31, 42, 16, 77, 58, 4, 71, 41, 28, 75, 78, 7, 90,
			61, 0, 20, 67, 33, 50, 57, 30, 92, 68, 8, 11},
		}, want: 30.085541344303557, ε: 0.000000001},
		// >>> random.seed(8120)
		// >>> nums = [random.randrange(-100, 100) for n in range(100)]
		{name: "100 randomly generated integers", args: args{nums: []float64{
			-98, 13, -29, 27, 93, -63, -4, 75, -82, 16, -100, -23, 55, 80, 77,
			-92, 63, -83, -46, 85, -34, -29, 21, -36, 58, -86, -78, 64, 90, 58,
			7, 39, -99, -45, -90, 25, 97, -13, 22, -46, -51, 28, -95, 25, 67,
			70, -11, -18, -88, 64, 16, -32, -65, -52, -4, -48, -44, -90, -60,
			-55, 66, 76, -28, 54, 25, -44, 9, -57, 30, 68, -89, 86, 98, -64,
			-65, -38, -15, -68, 54, 17, -91, 43, -18, -43, 50, 57, -86, 81, 22,
			-100, -59, 35, -33, 0, 14, -40, 85, 36, -83, -77},
		}, want: 60.193361827371035, ε: 0.000000001},
		// >>> random.seed(8120)
		// >>> nums = [random.gauss(0, 1) for n in range(10)]
		{name: "10 randomly generated gaussian", args: args{nums: []float64{
			0.8101616449919058, 0.0548988074717373, 0.07886885577556534,
			-2.3293260605060047, -0.2810299121940123, 0.2693364335404805,
			0.8479484467222266, 0.008913715915297816, -0.4237541806367912,
			-1.4756423117462945},
		}, want: 0.9845406683871757, ε: 0.000000001},
		// >>> random.seed(8120)
		// >>> nums = [random.gauss(0, 1) for n in range(10)]
		{name: "10 randomly generated gaussian", args: args{nums: []float64{
			10.008101616449919, 10.000548988074717, 10.000788688557755,
			9.97670673939494, 9.99718970087806, 10.002693364335405,
			10.008479484467223, 10.000089137159152, 9.995762458193632,
			9.985243576882537},
		}, want: 0.009845406683871799, ε: 0.000000001},
		{name: "zero length nums array", args: args{nums: []float64{}}, want: math.NaN(), ε: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stddev(tt.args.nums); !tolerance(got, tt.want, tt.ε) {
				t.Errorf("stddev() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_countNaNs(t *testing.T) {
	type args struct {
		nums []float64
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// >>> random.seed(8120)
		// >>> nums = [random.randrange(0, 100) for n in range(100)]
		{name: "100 non-negative randomly generated integers", args: args{nums: []float64{
			1, 56, 35, 63, 96, 18, 48, 87, 9, 58, 0, 38, 77, 90, 88, 4, 81, 8,
			27, 92, 33, 35, 60, 32, 79, 7, 11, 82, 95, 79, 53, 69, 0, 27, 5, 62,
			98, 43, 61, 27, 24, 64, 2, 62, 83, 85, 44, 41, 6, 82, 58, 34, 17,
			24, 48, 26, 28, 5, 20, 22, 83, 88, 36, 77, 62, 28, 54, 21, 65, 84,
			5, 93, 99, 18, 17, 31, 42, 16, 77, 58, 4, 71, 41, 28, 75, 78, 7, 90,
			61, 0, 20, 67, 33, 50, 57, 30, 92, 68, 8, 11},
		}, want: 0},
		// >>> random.seed(8120)
		// >>> nums = [random.randrange(-100, 100) for n in range(100)]
		{name: "100 randomly generated integers", args: args{nums: []float64{
			-98, 13, -29, 27, 93, -63, -4, 75, -82, 16, -100, -23, 55, 80, 77,
			-92, 63, -83, -46, 85, -34, -29, 21, -36, 58, -86, -78, 64, 90, 58,
			7, 39, -99, -45, -90, 25, 97, -13, 22, -46, -51, 28, -95, 25, 67,
			70, -11, -18, -88, 64, 16, -32, -65, -52, -4, -48, -44, -90, -60,
			-55, 66, 76, -28, 54, 25, -44, 9, -57, 30, 68, -89, 86, 98, -64,
			-65, -38, -15, -68, 54, 17, -91, 43, -18, -43, 50, 57, -86, 81, 22,
			-100, -59, 35, -33, 0, 14, -40, 85, 36, -83, -77},
		}, want: 0},
		// >>> random.seed(8120)
		// >>> nums = [random.gauss(0, 1) for n in range(10)]
		{name: "10 randomly generated gaussian", args: args{nums: []float64{
			0.8101616449919058, 0.0548988074717373, 0.07886885577556534,
			-2.3293260605060047, -0.2810299121940123, 0.2693364335404805,
			0.8479484467222266, 0.008913715915297816, -0.4237541806367912,
			-1.4756423117462945},
		}, want: 0},
		// >>> random.seed(8120)
		// >>> nums = [random.gauss(0, 1) for n in range(10)]
		{name: "10 randomly generated gaussian", args: args{nums: []float64{
			10.008101616449919, 10.000548988074717, 10.000788688557755,
			9.97670673939494, 9.99718970087806, 10.002693364335405,
			10.008479484467223, 10.000089137159152, 9.995762458193632,
			9.985243576882537},
		}, want: 0},
		{name: "zero length nums array",
			args: args{nums: []float64{}}, want: 0},
		{name: "only NaNs in nums array",
			args: args{nums: []float64{
				math.NaN(), math.NaN(), math.NaN(), math.NaN(), math.NaN()}},
			want: 5},
		{name: "some NaNs in nums array",
			args: args{nums: []float64{
				math.NaN(), -0.4705235438576944, math.NaN(), 0.5213549821939247,
				-1.5359989884260459, -0.22517298905321154, math.NaN()}},
			want: 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := countNaNs(tt.args.nums); got != tt.want {
				t.Errorf("countNaNs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_avg(t *testing.T) {
	type args struct {
		nums []float64
	}
	tests := []struct {
		name string
		args args
		want float64
		ε    float64
	}{
		// >>> random.seed(8120)
		// >>> nums = [random.randrange(0, 100) for n in range(100)]
		{name: "100 non-negative randomly generated integers", args: args{nums: []float64{
			1, 56, 35, 63, 96, 18, 48, 87, 9, 58, 0, 38, 77, 90, 88, 4, 81, 8,
			27, 92, 33, 35, 60, 32, 79, 7, 11, 82, 95, 79, 53, 69, 0, 27, 5, 62,
			98, 43, 61, 27, 24, 64, 2, 62, 83, 85, 44, 41, 6, 82, 58, 34, 17,
			24, 48, 26, 28, 5, 20, 22, 83, 88, 36, 77, 62, 28, 54, 21, 65, 84,
			5, 93, 99, 18, 17, 31, 42, 16, 77, 58, 4, 71, 41, 28, 75, 78, 7, 90,
			61, 0, 20, 67, 33, 50, 57, 30, 92, 68, 8, 11},
		}, want: 46.54, ε: 0.000000001},
		// >>> random.seed(8120)
		// >>> nums = [random.randrange(-100, 100) for n in range(100)]
		{name: "100 randomly generated integers", args: args{nums: []float64{
			-98, 13, -29, 27, 93, -63, -4, 75, -82, 16, -100, -23, 55, 80, 77,
			-92, 63, -83, -46, 85, -34, -29, 21, -36, 58, -86, -78, 64, 90, 58,
			7, 39, -99, -45, -90, 25, 97, -13, 22, -46, -51, 28, -95, 25, 67,
			70, -11, -18, -88, 64, 16, -32, -65, -52, -4, -48, -44, -90, -60,
			-55, 66, 76, -28, 54, 25, -44, 9, -57, 30, 68, -89, 86, 98, -64,
			-65, -38, -15, -68, 54, 17, -91, 43, -18, -43, 50, 57, -86, 81, 22,
			-100, -59, 35, -33, 0, 14, -40, 85, 36, -83, -77},
		}, want: -6.46, ε: 0.000000001},
		// >>> random.seed(8120)
		// >>> nums = [random.gauss(0, 1) for n in range(10)]
		{name: "10 randomly generated gaussian", args: args{nums: []float64{
			0.8101616449919058, 0.0548988074717373, 0.07886885577556534,
			-2.3293260605060047, -0.2810299121940123, 0.2693364335404805,
			0.8479484467222266, 0.008913715915297816, -0.4237541806367912,
			-1.4756423117462945},
		}, want: -0.24396245606658895, ε: 0.000000001},
		// >>> random.seed(8120)
		// >>> nums = [random.gauss(0, 1) for n in range(10)]
		{name: "10 randomly generated gaussian", args: args{nums: []float64{
			10.008101616449919, 10.000548988074717, 10.000788688557755,
			9.97670673939494, 9.99718970087806, 10.002693364335405,
			10.008479484467223, 10.000089137159152, 9.995762458193632,
			9.985243576882537},
		}, want: 9.997560375439335, ε: 0.000000001},
		{name: "zero length nums array", args: args{nums: []float64{}}, want: math.NaN(), ε: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := avg(tt.args.nums); !tolerance(got, tt.want, tt.ε) {
				t.Errorf("avg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func tolerance(actual, want, ε float64) bool {
	if (math.IsNaN(want) && math.IsNaN(actual)) || actual == want {
		return true
	}
	if actual > want {

		return actual-want <= ε
	}
	return want-actual <= ε
}
