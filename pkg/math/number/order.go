package number

func FindOrder(a, N int, f float64) (int, int, bool) {
	c := ContinuedFraction(f)
	for i := 1; i < len(c)+1; i++ {
		s, r := Fraction(c[:i])
		if ModExp(a, r, N) == 1 {
			return s, r, true
		}
	}

	return 0, 0, false
}
