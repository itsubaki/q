package number

func FindOrder(a, N int, binary []int, eps ...float64) (int, int, float64, bool) {
	if len(binary) < 1 {
		return 0, 1, 0, false
	}

	f := BinaryToFloat64(binary)
	c := ContinuedFraction(f, eps...)

	s, r, d := Convergent(c[:1])
	for i := 1; i < len(c); i++ {
		s, r, d = Convergent(c[:i+1])
		if ModExp(a, r, N) == 1 {
			return s, r, d, true
		}
	}

	return s, r, d, false
}
