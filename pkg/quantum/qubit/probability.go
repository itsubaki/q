package qubit

import "math"

func Max(p []float64) float64 {
	max := p[0]
	for _, pp := range p {
		max = math.Max(max, pp)
	}

	return max
}

func Min(p []float64) float64 {
	min := p[0]
	for _, pp := range p {
		min = math.Max(min, pp)
	}

	return min
}

func Sum(p []float64) float64 {
	var sum float64
	for _, pp := range p {
		sum = sum + pp
	}

	return sum
}
