package number

import (
	"fmt"
	"math"
)

func BinaryToFloat64(binary []int) float64 {
	var d float64
	for i, b := range binary {
		if b != 0 && b != 1 {
			panic(fmt.Sprintf("invalid input: %v", binary))
		}

		if b == 0 {
			continue
		}

		d = d + math.Pow(0.5, float64(i+1))
	}

	return d
}

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

func Fraction(cf []int) (int, int, float64) {
	if len(cf) == 1 && cf[0] == 0 {
		return 0, 1, 0
	}

	s, r := 1, cf[len(cf)-1]
	for i := 2; i < len(cf); i++ {
		s, r = r, cf[len(cf)-i]*r+s
	}

	return s, r, float64(s) / float64(r)
}

func Approximate(cf []int) float64 {
	last := len(cf) - 1
	f := 1.0 / float64(cf[last])
	for i := last - 1; i > 0; i-- {
		f = 1.0 / (float64(cf[i]) + f)
	}

	return f
}
