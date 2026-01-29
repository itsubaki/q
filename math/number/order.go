package number

// FindOrder returns convergent s/r and its real number such that a**r mod N = 1.
func FindOrder(a, N int, m float64, tol ...float64) (int, int, float64, bool) {
	c := ContinuedFraction(m, tol...)

	s, r, d := Convergent(c[:1])
	for i := 1; i < len(c); i++ {
		s, r, d = Convergent(c[:i+1])
		if ModExp(a, r, N) == 1 {
			return s, r, d, true
		}
	}

	return s, r, d, false
}
