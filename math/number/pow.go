package number

import (
	"math"
	"strconv"

	"github.com/itsubaki/q/math/epsilon"
)

// Pow returns a**r, the base-a exponential of r.
func Pow(a, r int) int {
	if a == 0 {
		return 0
	}

	if r == 0 {
		return 1
	}

	p := a
	for i := 1; i < r; i++ {
		p = p * a
	}

	return p
}

// ModExp returns a**r mod N.
func ModExp(a, r, N int) int {
	if a == 0 {
		return 0
	}

	if r == 0 {
		return 1
	}

	p := a
	for i := 1; i < r; i++ {
		p = (p * a) % N
	}

	return p
}

// ModExp2 returns a**(2**j) mod N.
func ModExp2(a, j, N int) int {
	if a == 0 {
		return 0
	}

	if j == 0 {
		return a % N
	}

	p := a
	for i := 0; i < j; i++ {
		p = (p * p) % N
	}

	return p
}

// BaseExp returns a and b such that a**b == N. If N is not a power of a, return false.
func BaseExp(N int, eps ...float64) (int, int, bool) {
	e := epsilon.E13(eps...)

	l := len(strconv.FormatInt(int64(N), 2))
	for i := l; 1 < i; i-- {
		a := math.Pow(float64(N), 1.0/float64(i))
		if a-math.Trunc(a) < e {
			if Pow(int(a), i) == N {
				return int(a), i, true
			}
		}

		if 1-(a-math.Trunc(a)) < e {
			if Pow(int(a)+1, i) == N {
				return int(a) + 1, i, true
			}
		}
	}

	return 0, 0, false
}
