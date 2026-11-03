package number

import (
	"cmp"
	"math"
)

// Ldexp returns a*2^b.
func Ldexp(a, b int) float64 {
	return math.Ldexp(float64(a), b)
}

// Pow returns a^r.
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

// BaseExp returns a and b such that a^b == N. If N is not a perfect power, it returns false.
func BaseExp(N int) (int, int, bool) {
	if N < 2 {
		// 0^1 == 0, 1^1 == 1
		return N, 1, true
	}

	for b := Log2(N); b >= 2; b-- {
		lo, hi := 2, N
		for lo <= hi {
			mid := lo + (hi-lo)/2
			switch comparePow(mid, b, N) {
			case 0:
				return mid, b, true
			case 1:
				hi = mid - 1
			case -1:
				lo = mid + 1
			}
		}
	}

	return 0, 0, false
}

// comparePow compares a^b with n without overflowing.
// a and n must be positive, and b must be nonnegative.
func comparePow(a, b, n int) int {
	p := 1
	for range b {
		if p > n/a {
			return 1
		}

		p *= a
	}

	return cmp.Compare(p, n)
}
