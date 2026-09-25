package main

import (
	"fmt"
	"math"

	"github.com/itsubaki/q"
)


// E estimates the correlation E(A, B) between two measurement settings.
//
// The correlation is defined as:
//   +1 when Alice's and Bob's measurement results are equal
//   -1 when they are different
//
// By repeating the experiment over many shots, we estimate the
// expectation value of the correlation.
func E(thetaA, thetaB float64, shots int) float64 {
	var sum int
	for range shots {
		qsim := q.New()

		a := qsim.Zero()
		b := qsim.Zero()

		qsim.H(a)
		qsim.CNOT(a, b)

		qsim.RY(-thetaA, a)
		qsim.RY(-thetaB, b)

		ma := qsim.Measure(a)
		mb := qsim.Measure(b)

		if ma.Equal(mb) {
			sum++
		} else {
			sum--
		}
	}

	return float64(sum) / float64(shots)
}

func main() {
	const shots = 1000

	A, Ap := 0.0, math.Pi/2
	B, Bp := math.Pi/4, -math.Pi/4

	EAB := E(A, B, shots)
	EABp := E(A, Bp, shots)
	EApB := E(Ap, B, shots)
	EApBp := E(Ap, Bp, shots)

	// CHSH parameter:
	//  S = E(A,B) + E(A,B') + E(A',B) - E(A',B')
	//
	// Local hidden-variable theories satisfy |S| <= 2.
	// Quantum mechanics allows |S| <= 2*sqrt(2) (Tsirelson bound).
	S := EAB + EABp + EApB - EApBp

	fmt.Printf("shots     = %d\n", shots)
	fmt.Printf("E(A,B)    = %+.4f\n", EAB)
	fmt.Printf("E(A,B')   = %+.4f\n", EABp)
	fmt.Printf("E(A',B)   = %+.4f\n", EApB)
	fmt.Printf("E(A',B')  = %+.4f\n", EApBp)
	fmt.Printf("S         = %+.4f\n", S)
	fmt.Printf("2*sqrt(2) = %+.4f\n", 2*math.Sqrt(2))
}
