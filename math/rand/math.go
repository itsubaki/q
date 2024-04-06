package rand

import (
	"math/rand/v2"

	"github.com/itsubaki/q/math/number"
)

var Float64 = rand.Float64

// Const returns a constant number in [0.0, 1.0).
func Const(seed ...uint64) func() float64 {
	var s0, s1 uint64
	if len(seed) > 0 {
		s0 = seed[0]
	}

	if len(seed) > 1 {
		s1 = seed[1]
	}

	return rand.New(rand.NewPCG(s0, s1)).Float64
}

// Coprime returns a random coprime number in [2, N).
func Coprime(N int) int {
	min, max := 2, N-2
	for {
		a := rand.N(max-1) + min
		if number.GCD(N, a) == 1 {
			return a
		}
	}
}
