package number

import (
	"math"
)

func ContinuedFraction(f float64, eps ...float64) []int {
	e := 1e-3
	if len(eps) > 0 {
		e = eps[0]
	}

	if f < e {
		return []int{0}
	}

	list := make([]int, 0)
	r := f
	for {
		t := math.Trunc(r)
		list = append(list, int(t))

		diff := r - t
		if diff < e {
			break
		}

		r = 1.0 / diff
	}

	return list
}

func Fraction(continued []int) (int, int) {
	if len(continued) == 1 && continued[0] == 0 {
		return 0, 1
	}

	s, r := 1, continued[len(continued)-1]
	for i := 2; i < len(continued); i++ {
		s, r = r, continued[len(continued)-i]*r+s
	}

	return s, r
}

func Approximate(continued []int) float64 {
	last := len(continued) - 1
	f := 1.0 / float64(continued[last])
	for i := last - 1; i > 0; i-- {
		f = 1.0 / (float64(continued[i]) + f)
	}

	return f
}
