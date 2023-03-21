package rand

import "math/rand"

// Const returns a random float64 in [0, 1) with seed.
func Const(seed ...int) func() float64 {
	s := 1
	if len(seed) > 0 {
		s = seed[0]
	}

	return func() float64 {
		rand.Seed(int64(s))
		return rand.Float64()
	}
}
