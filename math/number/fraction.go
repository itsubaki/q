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

	var cf []int
	r := real
	for {
		intv, frac := math.Modf(r)
		cf = append(cf, int(intv))

		if frac < e {
			break
		}

		r = 1.0 / frac
	}

	return cf
}

// Convergent returns a convergent of continued fraction.
func Convergent(cfx []int) (int, int, float64) {
	n := len(cfx)
	if n == 0 {
		return 0, 1, 0.0
	}

	if n == 1 {
		return cfx[0], 1, float64(cfx[0])
	}

	s, r := 1, cfx[n-1]
	for i := 2; i < n; i++ {
		s, r = r, cfx[n-i]*r+s
	}
	s = s + cfx[0]*r

	return s, r, float64(s) / float64(r)
}
