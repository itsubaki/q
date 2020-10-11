package number

func FindOrder(a, N int, binary []int) (int, int, bool) {
	if len(binary) < 1 {
		return 0, 1, false
	}

	f := BinaryToFloat64(binary)
	c := ContinuedFraction(f)

	s, r := Fraction(c[:1])
	for i := 2; i < len(c)+1; i++ {
		s, r = Fraction(c[:i])
		if r < N && ModExp(a, r, N) == 1 {
			return s, r, true
		}
	}

	return s, r, false
}
