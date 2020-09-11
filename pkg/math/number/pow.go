package number

func Pow(a, r int) int {
	if a == 0 {
		return 0
	}

	if r == 0 {
		return 1
	}

	p := a
	for i := 1; i < r; i++ {
		p = p * a
	}

	return p
}
