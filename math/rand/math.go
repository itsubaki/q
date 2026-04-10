package rand

import (
	"math/rand/v2"

	"github.com/itsubaki/q/math/number"
)

// Float64 returns, as a float64, a pseudo-random number in [0.0, 1.0).
var Float64 = rand.Float64

// Const returns a pseudo-random number generator initialized with the given seeds.
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

// Coprime returns a random integer in [2, N) that is coprime to N.
// It requires N to be greater than 2.
func Coprime(N int) int {
	for {
		a := rand.N(N-2) + 2
		if number.GCD(N, a) == 1 {
			return a
		}
	}
}
