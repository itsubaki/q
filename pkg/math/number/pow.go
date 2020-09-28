package number

import (
	"math"
	"strconv"
)

func Pow(a, r int) int {
	if a == 0 {
		return 0
	}

	if r == 0 {
		return 1
	}

	p := a
	for i := 1; i < r; i++ {
		p = p * a
	}

	return p
}

// ModExp2 returns modular exponentiation
//  a^2^j mod N
func ModExp2(a, j, N int) int {
	if a == 0 {
		return 0
	}

	if j == 0 {
		return a % N
	}

	p := a
	for i := 0; i < j; i++ {
		p = (p * p) % N
	}

	return p
}

func BaseExp(N int) (int, int, bool) {
	l := len(strconv.FormatInt(int64(N), 2))
	for i := l; 1 < i; i-- {
		a := math.Pow(float64(N), 1.0/float64(i))
		if Pow(int(a), i) == N {
			return int(a), i, true
		}
	}

	return 0, 0, false
}
