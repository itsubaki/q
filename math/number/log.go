package number

import "errors"

var ErrNotPowerOf2 = errors.New("N must be a power of 2")

// Log2 returns the base 2 logarithm of N.
// N must be a power of 2.
func Log2(N int) (int, error) {
	if N == 1 {
		return 0, nil
	}

	var n int = 1
	for {
		if N%2 != 0 {
			return -1, ErrNotPowerOf2
		}

		if N/2 == 1 {
			break
		}

		N = N / 2
		n++
	}

	return n, nil
}
