package number

import "math"

func Max(prob []float64) float64 {
	max := prob[0]
	for _, p := range prob {
		max = math.Max(max, p)
	}

	return max
}

func Min(prob []float64) float64 {
	min := prob[0]
	for _, p := range prob {
		min = math.Min(min, p)
	}

	return min
}

func Sum(prob []float64) float64 {
	var sum float64
	for _, p := range prob {
		sum = sum + p
	}

	return sum
}
