package number

import (
	"math"

	"github.com/itsubaki/q/math/epsilon"
)

// ContinuedFraction returns a continued fraction of real.
func ContinuedFraction(real float64, tol ...float64) []int {
	if epsilon.IsZeroF64(real, tol...) {
		return []int{0}
	}

	var cf []int
	r := real
	for {
		a, frac := math.Modf(r)
		cf = append(cf, int(a))
		if epsilon.IsZeroF64(frac, tol...) {
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

	p0, p1 := 1, cfx[0]
	q0, q1 := 0, 1

	for i := 1; i < n; i++ {
		p0, p1 = p1, cfx[i]*p1+p0
		q0, q1 = q1, cfx[i]*q1+q0
	}

	return p1, q1, float64(p1) / float64(q1)
}
