package number

import "math"

func IsPrime(N int) bool {
	if N < 2 || IsEven(N) {
		return false
	}

	if N == 2 {
		return true
	}

	for i := 3; i < int(math.Sqrt(float64(N)))+1; i = i + 2 {
		if N%i == 0 {
			return false
		}
	}

	return true
}

func IsOdd(v int) bool {
	return !IsEven(v)
}

func IsEven(v int) bool {
	if v%2 == 0 {
		return true
	}

	return false
}
