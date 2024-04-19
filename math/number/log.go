package number

import "fmt"

func Log2(N int) (int, error) {
	if N == 1 {
		return 0, nil
	}

	if N%2 != 0 {
		return 0, fmt.Errorf("N must be a power of 2. N=%v", N)
	}

	var n int = 1
	for {
		if N/2 == 1 {
			break
		}

		N = N / 2
		n++
	}

	return n, nil
}
