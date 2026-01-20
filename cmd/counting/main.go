package main

import (
	"flag"
	"fmt"
	"math"
	"sort"

	"github.com/itsubaki/q"
	"github.com/itsubaki/q/quantum/qubit"
)

// controlledG applies the Grover operator for 2x2 mini-sudoku solutions.
// The number of solutions `M` is 2.
func controlledG(qsim *q.Q, c q.Qubit, r, s []q.Qubit, a q.Qubit) {
	oracle(qsim, c, r, s, a)
	diffuser(qsim, c, r)
}

func oracle(qsim *q.Q, c q.Qubit, r, s []q.Qubit, a q.Qubit) {
	xor := func(x, y, z q.Qubit) {
		qsim.CNOT(x, z)
		qsim.CNOT(y, z)
	}

	xor(r[0], r[1], s[0]) // a != b
	xor(r[2], r[3], s[1]) // c != d
	xor(r[0], r[2], s[2]) // a != c
	xor(r[1], r[3], s[3]) // b != d

	// apply Z if all a are 1
	qsim.ControlledZ([]q.Qubit{c, s[0], s[1], s[2], s[3]}, []q.Qubit{a})

	// uncompute
	xor(r[1], r[3], s[3])
	xor(r[0], r[2], s[2])
	xor(r[2], r[3], s[1])
	xor(r[0], r[1], s[0])
}

func diffuser(qsim *q.Q, c q.Qubit, r []q.Qubit) {
	qsim.H(r...)
	qsim.X(r...)
	qsim.ControlledZ([]q.Qubit{c, r[0], r[1], r[2]}, []q.Qubit{r[3]})
	qsim.X(r...)
	qsim.H(r...)
}

func top(s []qubit.State, n int) []qubit.State {
	sort.Slice(s, func(i, j int) bool {
		return s[i].Probability() > s[j].Probability()
	})

	return s[:min(n, len(s))]
}

func main() {
	var t int
	flag.IntVar(&t, "t", 3, "precision bits")
	flag.Parse()

	qsim := q.New()

	// initialize
	c := qsim.Zeros(t) // for phase estimation
	r := qsim.Zeros(4) // data qubits for the Grover search space
	s := qsim.Zeros(4) // ancilla qubits for comparing Sudoku constraints
	a := qsim.Zero()   // ancilla qubit for oracle

	// superposition
	qsim.H(c...)
	qsim.H(r...)
	qsim.X(a)
	qsim.H(a) // for phase kickback

	// phase estimation
	for i := range c {
		// apply controlled-G**(2**i) where control is c[len(c)-1-i]
		for range 1 << i {
			controlledG(qsim, c[len(c)-1-i], r, s, a)
		}
	}

	// inverse quantum Fourier transform
	qsim.InvQFT(c...)

	// measurement
	qsim.Measure(r...)
	qsim.Measure(s...)
	qsim.Measure(a)

	// results
	N, size := 1<<len(r), 1<<t
	for _, s := range top(qsim.State(c), 16) {
		phi := float64(s.Int()) / float64(size)        // phi = k / 2**t
		theta := math.Pi * phi                         // theta = pi * phi
		M := float64(N) * math.Pow(math.Sin(theta), 2) // M = N * (sin(theta))**2

		// [001][  1]( 0.2024-0.4886i): 0.2797; phi=0.1250, theta=0.3927, M=2.3431
		// [111][  7]( 0.2024 0.4886i): 0.2797; phi=0.8750, theta=2.7489, M=2.3431
		// [011][  3](-0.3279 0.0350i): 0.1087; phi=0.3750, theta=1.1781, M=13.6569
		// [101][  5](-0.3279-0.0350i): 0.1087; phi=0.6250, theta=1.9635, M=13.6569
		fmt.Printf("%v; phi=%.4f, theta=%.4f, M=%.4f\n", s, phi, theta, M)
	}
}
