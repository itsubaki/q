package main

import (
	"flag"
	"fmt"
	"math"
	"sort"

	"github.com/itsubaki/q"
	"github.com/itsubaki/q/math/number"
	"github.com/itsubaki/q/quantum/qubit"
)

// controlledG applies the Grover operator for 2x2 mini-sudoku solutions.
// The number of solutions `M` is 2.
func controlledG(qsim *q.Q, c q.Qubit, r, a []q.Qubit) {
	oracle(qsim, c, r, a)
	diffuser(qsim, c, r)
}

func oracle(qsim *q.Q, c q.Qubit, r, a []q.Qubit) {
	xor := func(x, y, z q.Qubit) {
		qsim.ControlledNot([]q.Qubit{c, x}, []q.Qubit{z})
		qsim.ControlledNot([]q.Qubit{c, y}, []q.Qubit{z})
	}

	xor(r[0], r[1], a[0]) // a != b
	xor(r[2], r[3], a[1]) // c != d
	xor(r[0], r[2], a[2]) // a != c
	xor(r[1], r[3], a[3]) // b != d

	// apply Z if all a are 1
	qsim.ControlledZ([]q.Qubit{c, a[0], a[1], a[2]}, []q.Qubit{a[3]})

	// uncompute
	xor(r[1], r[3], a[3])
	xor(r[0], r[2], a[2])
	xor(r[2], r[3], a[1])
	xor(r[0], r[1], a[0])
}

func diffuser(qsim *q.Q, c q.Qubit, r []q.Qubit) {
	qsim.ControlledH([]q.Qubit{c}, r)
	qsim.ControlledX([]q.Qubit{c}, r)
	qsim.ControlledZ([]q.Qubit{c, r[0], r[1], r[2]}, []q.Qubit{r[3]})
	qsim.ControlledX([]q.Qubit{c}, r)
	qsim.ControlledH([]q.Qubit{c}, r)
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
	a := qsim.Zeros(4) // ancilla qubits for comparing Sudoku constraints

	// superposition
	qsim.H(c...)
	qsim.H(r...)

	// phase estimation
	for i := range len(c) {
		// apply controlled-G**(2**i) where control is c[len(c)-1-i]
		for range 1 << i {
			controlledG(qsim, c[len(c)-1-i], r, a)
		}
	}

	// inverse quantum Fourier transform
	qsim.InvQFT(c...)

	// measure unused registers
	qsim.Measure(r...)
	qsim.Measure(a...)

	// results
	N := number.Pow(2, len(r))
	for _, s := range top(qsim.State(c), 16) {
		theta := float64(s.Int()) / float64(number.Pow(2, len(c))) // theta = k / 2**len(c)
		M := float64(N) * math.Pow(math.Sin(math.Pi*theta), 2)     // M = N * (sin(pi * theta))**2

		// [101][  5]( 0.5585 0.0597i): 0.3155; theta=0.6250, M=13.6569
		// [011][  3]( 0.5585-0.0597i): 0.3155; theta=0.3750, M=13.6569
		// [100][  4]( 0.3757 0.0000i): 0.1412; theta=0.5000, M=16.0000
		// [111][  7](-0.0414-0.3282i): 0.1094; theta=0.8750, M=2.3431
		// [001][  1](-0.0414 0.3282i): 0.1094; theta=0.1250, M=2.3431
		// [010][  2](-0.0215-0.0633i): 0.0045; theta=0.2500, M=8.0000
		// [110][  6](-0.0215 0.0633i): 0.0045; theta=0.7500, M=8.0000
		// [000][  0]( 0.0121 0.0000i): 0.0001; theta=0.0000, M=0.0000
		fmt.Printf("%v; theta=%.4f, M=%.4f\n", s, theta, M)
	}
}
