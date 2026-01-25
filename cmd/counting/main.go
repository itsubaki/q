package main

import (
	"flag"
	"fmt"
	"math"

	"github.com/itsubaki/q"
)

func oracle(qsim *q.Q, r, s []q.Qubit, c, a q.Qubit) {
	xor := func(x, y, z q.Qubit) {
		qsim.CNOT(x, z)
		qsim.CNOT(y, z)
	}

	xor(r[0], r[1], s[0]) // a != b
	xor(r[2], r[3], s[1]) // c != d
	xor(r[0], r[2], s[2]) // a != c
	xor(r[1], r[3], s[3]) // b != d

	// apply Z if s and c are all 1
	qsim.X(a)
	qsim.ControlledZ([]q.Qubit{c, s[0], s[1], s[2], s[3]}, []q.Qubit{a})
	qsim.X(a)

	// uncompute
	xor(r[1], r[3], s[3])
	xor(r[0], r[2], s[2])
	xor(r[2], r[3], s[1])
	xor(r[0], r[1], s[0])
}

func diffuser(qsim *q.Q, c q.Qubit, r []q.Qubit) {
	qsim.H(r...)
	qsim.X(r...)
	qsim.ControlledZ(append([]q.Qubit{c}, r[:len(r)-1]...), []q.Qubit{r[len(r)-1]})
	qsim.X(r...)
	qsim.H(r...)
}

// controlledG applies the Grover operator for 2x2 mini-sudoku solutions.
// The number of solutions `M` is 2.
func controlledG(qsim *q.Q, r, s []q.Qubit, c, a q.Qubit) {
	oracle(qsim, r, s, c, a)
	diffuser(qsim, c, r)
}

func main() {
	var t, top int
	flag.IntVar(&t, "t", 7, "precision bits")
	flag.IntVar(&top, "top", 8, "top results")
	flag.Parse()

	qsim := q.New()

	// initialize
	c := qsim.Zeros(t) // for phase estimation
	r := qsim.Zeros(4) // data qubits for the Grover search space
	s := qsim.Zeros(4) // ancilla qubits for comparing Sudoku constraints
	a := qsim.Zero()   // ancilla qubit for oracle

	// initialize
	qsim.H(c...)
	qsim.H(r...)
	qsim.H(a)

	// phase estimation
	for i := range c {
		// apply controlled-G**(2**i) where control is c[len(c)-1-i]
		for range 1 << i {
			controlledG(qsim, r, s, c[len(c)-1-i], a)
		}
	}

	// inverse quantum Fourier transform
	qsim.Swap(c...)
	qsim.InvQFT(c...)

	// results
	N, size := 1<<len(r), 1<<t
	for _, s := range q.Top(qsim.State(c, r, s, a), top) {
		phi := float64(s.Int()) / float64(size)            // phi = k / 2**t
		theta := 2 * math.Pi * min(phi, math.Abs(phi-0.5)) // theta = 2*pi*phi
		M := float64(N) * math.Pow(math.Sin(theta/2), 2)   // M = N*(sin(theta/2))**2

		fmt.Printf("%v; phi=%.4f, theta=%.4f, M=%.4f, eps=%.4f\n", s, phi, theta, M, math.Abs(2-M))
	}
}
