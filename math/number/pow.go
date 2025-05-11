package number

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
