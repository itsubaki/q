package number

import "math"

// Max returns the maximum value in p.
func Max(p []float64) float64 {
	max := p[0]
	for _, pp := range p {
		max = math.Max(max, pp)
	}

	return max
}

// Min returns the minimum value in p.
func Min(p []float64) float64 {
	min := p[0]
	for _, pp := range p {
		min = math.Min(min, pp)
	}

	return min
}

// Sum returns the sum of p.
func Sum(p []float64) float64 {
	var sum float64
	for _, pp := range p {
		sum = sum + pp
	}

	return sum
}
