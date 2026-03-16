package number

import "math"

// Ldexp returns a * 2**b.
func Ldexp(a, b int) float64 {
	return math.Ldexp(float64(a), b)
}

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

// BaseExp returns a and b such that a**b == N. If N is not a power of a, return false.
func BaseExp(N int) (int, int, bool) {
	if N < 2 {
		// 0**1 == 0, 1**1 == 1
		return N, 1, true
	}

	for b := Log2(N); b >= 2; b-- {
		lo, hi := 2, N
		for lo < hi+1 {
			mid := (lo + hi) / 2
			p := Pow(mid, b)

			if p == N {
				return mid, b, true
			}

			if p < 0 || p > N {
				hi = mid - 1
				continue
			}

			lo = mid + 1
		}
	}

	return 0, 0, false
}
