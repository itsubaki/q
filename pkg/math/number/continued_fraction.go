package number

import (
	"fmt"
	"math"
)

func ContinuedFraction(f float64, eps ...float64) ([]int, int, int) {
	e := 1e-3
	if len(eps) > 0 {
		e = eps[0]
	}

	if f < e {
		return []int{0}, 0, 1
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

	if len(list) == 1 {
		return list, 1, list[0]
	}

	n, d := 1, list[len(list)-1]
	for i := 2; i < len(list); i++ {
		n, d = d, list[len(list)-i]*d+n
	}

	return list, n, d
}

func InverseContinuedFraction(c []int, max ...int) float64 {
	if len(max) < 1 {
		max[0] = len(c)
	}

	if max[0] > len(c) {
		panic(fmt.Sprintf("invalid input. %v %v", c, max))
	}

	f := 1.0 / float64(c[max[0]])
	for i := max[0] - 1; i > 0; i-- {
		f = 1.0 / (float64(c[i]) + f)
	}

	return f
}
