package number

import "math"

func Fraction(f float64, eps float64) ([]int, int, int) {
	continued := make([]int, 0)

	reciprocal := f
	for {
		continued = append(continued, int(reciprocal))
		diff := reciprocal - math.Trunc(reciprocal)
		if diff < eps {
			break
		}

		reciprocal = 1.0 / diff
	}

	// TODO return numerator denominator

	return continued, 3, 7
}
