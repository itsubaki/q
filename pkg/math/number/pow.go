package number

func Pow(a, r int) int {
	o := a
	for i := 1; i < r; i++ {
		o = o * a
	}

	return o
}
