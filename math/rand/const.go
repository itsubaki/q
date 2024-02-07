package rand

import "math/rand/v2"

// Const returns a constant number in [0.0, 1.0).
func Const(seed ...uint64) func() float64 {
	var s uint64
	if len(seed) > 0 {
		s = seed[0]
	}

	return rand.New(rand.NewPCG(s, s)).Float64
}
