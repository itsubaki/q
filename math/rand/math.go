package rand

import (
	"math/rand/v2"

	"github.com/itsubaki/q/math/number"
)

var Float64 = rand.Float64

// Const returns a sequence of numbers fixed by seeds.
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
// It requires N to be greater than 2.
func Coprime(N int) int {
	for {
		a := rand.N(N-2) + 2
		if number.GCD(N, a) == 1 {
			return a
		}
	}
}
