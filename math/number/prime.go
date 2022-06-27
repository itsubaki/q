package number

import "math"

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

func IsTrivial(N int, factor ...int) bool {
	for _, p := range factor {
		if 1 < p && p < N && N%p == 0 {
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
