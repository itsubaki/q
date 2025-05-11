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

	// exponentiation by squaring
	p := 1
	for r > 0 {
		if r&1 == 1 {
			p = p * a
		}

		a = a * a
		r >>= 1
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

	if N == 1 {
		return 0
	}

	// exponentiation by squaring
	p := 1
	for r > 0 {
		if r&1 == 1 {
			p = (p * a) % N
		}

		a = (a * a) % N
		r >>= 1
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
	for range j {
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
