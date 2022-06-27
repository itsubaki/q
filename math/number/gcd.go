package number

func GCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}

	return a
}
