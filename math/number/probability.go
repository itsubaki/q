package number

// Sum returns the sum of p.
func Sum(p []float64) float64 {
	var sum float64
	for _, v := range p {
		sum = sum + v
	}

	return sum
}
