package qubit

import "math"

type Probability float64

func Max(p []Probability) Probability {
	max := float64(p[0])
	for _, pp := range p {
		max = math.Max(max, float64(pp))
	}

	return Probability(max)
}

func Min(p []Probability) Probability {
	min := float64(p[0])
	for _, pp := range p {
		min = math.Max(min, float64(pp))
	}

	return Probability(min)
}

func Sum(p []Probability) Probability {
	var sum Probability
	for _, pp := range p {
		sum = sum + pp
	}
	return sum
}
