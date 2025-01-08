package number

// Log2 returns the integer part of the base 2 logarithm of N.
func Log2(N int) int {
	var count int
	for N > 1 {
		N >>= 1
		count++
	}

	return count
}

func IsPowOf2(N int) bool {
	return N > 0 && N&(N-1) == 0
}
