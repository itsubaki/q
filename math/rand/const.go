package rand

import "math/rand"

// Const returns a constant number in [0.0, 1.0).
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
