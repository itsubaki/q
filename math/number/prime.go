package number

import "math"

// IsPrime returns true if N is prime number.
func IsPrime(N int) bool {
	if N < 2 {
		return false
	}

	if N == 2 {
		return true
	}

	if IsEven(N) {
		return false
	}

	for i := 3; i < int(math.Sqrt(float64(N)))+1; i = i + 2 {
		if N%i == 0 {
			return false
		}
	}

	return true
}

// IsTrivial returns true if factor is trivial factor of N.
func IsTrivial(N int, factor ...int) bool {
	for _, p := range factor {
		if 1 < p && p < N && N%p == 0 {
			return false
		}
	}

	return true
}

// IsOdd returns true if v is odd number.
func IsOdd(v int) bool {
	return !IsEven(v)
}

// IsEven returns true if v is even number.
func IsEven(v int) bool {
	if v%2 == 0 {
		return true
	}

	return false
}
