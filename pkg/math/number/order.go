package number

func FindOrder(a, N int, binary []int, eps ...float64) (int, int, float64, bool) {
	if len(binary) < 1 {
		return 0, 1, 0, false
	}

	f := BinaryToFloat64(binary)
	c := ContinuedFraction(f, eps...)

	s, r, d := Fraction(c[:1])
	for i := 2; i < len(c)+1; i++ {
		s, r, d = Fraction(c[:i])
		if r < N && ModExp(a, r, N) == 1 {
			return s, r, d, true
		}
	}

	return s, r, d, false
}
