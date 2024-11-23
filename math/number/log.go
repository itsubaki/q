package number

// Log2 returns the base 2 logarithm of N.
// N must be a power of 2.
func Log2(N int) int {
	if N > 0 && (N&(N-1) != 0) {
		return -1
	}

	var count int
	for N > 1 {
		N >>= 1
		count++
	}

	return count
}
