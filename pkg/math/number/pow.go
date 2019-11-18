package number

func Pow(a, r int) int {
	p := a
	for i := 1; i < r; i++ {
		p = p * a
	}

	return p
}
