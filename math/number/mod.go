package number

// ModExp returns a**r mod N.
func ModExp(a, r, N int) int {
	if a == 0 || N == 1 {
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

	p := a % N
	for range j {
		p = (p * p) % N
	}

	return p
}
