package main

import (
	"fmt"
	"math"

	"github.com/itsubaki/q"
)

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
	S := EAB + EABp + EApB - EApBp

	fmt.Printf("shots     = %d\n", shots)
	fmt.Printf("E(A,B)    = %+.4f\n", EAB)
	fmt.Printf("E(A,B')   = %+.4f\n", EABp)
	fmt.Printf("E(A',B)   = %+.4f\n", EApB)
	fmt.Printf("E(A',B')  = %+.4f\n", EApBp)
	fmt.Printf("S         = %+.4f\n", S)
	fmt.Printf("2*sqrt(2) = %+.4f\n", 2*math.Sqrt(2))
}
