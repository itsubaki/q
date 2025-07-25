package number

import (
	"math"

	"github.com/itsubaki/q/math/epsilon"
)

// ContinuedFraction returns a continued fraction of real.
// if eps is empty, epsilon.E3 is used.
func ContinuedFraction(real float64, eps ...float64) []int {
	e := epsilon.E3(eps...)
	if real < e {
		return []int{0}
	}

	list := make([]int, 0)
	r := real
	for {
		list = append(list, int(r))

		diff := r - math.Trunc(r)
		if diff < e {
			break
		}

		r = 1.0 / diff
	}

	return list
}

// Convergent returns a convergent of continued fraction.
func Convergent(cfx []int) (int, int, float64) {
	l := len(cfx)
	if l == 1 {
		return cfx[0], 1, float64(cfx[0])
	}

	s, r := 1, cfx[l-1]
	for i := 2; i < l; i++ {
		s, r = r, cfx[l-i]*r+s
	}
	s = s + cfx[0]*r

	return s, r, float64(s) / float64(r)
}
