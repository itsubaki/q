package main

import (
	"flag"
	"fmt"
	"math"

	"github.com/itsubaki/q"
	"github.com/itsubaki/q/math/number"
)

// oracle constructs a Grover oracle for validating 2x2 mini-sudoku solutions.
// It applies a phase flip (−1) only to states that represent valid solutions.
//
// The 2x2 sudoku grid is represented as follows:
//
//	| a | b |
//	| c | d |
//
// The input `r` is a slice of 4 qubits representing the cells a, b, c, and d
// in row-major order: [a, b, c, d].
// The `s` slice must contain 4 ancilla qubits used for intermediate checks.
// The `a` qubit is the oracle’s phase flag (target).
// THe `a` is initialized to |−> = (|0> − |1>)/√2 to facilitate phase kickback.
//
// The oracle checks the following uniqueness constraints:
//
//   - a != b
//   - c != d
//   - a != c
//   - b != d
//
// If **all** constraints are satisfied (i.e., the input represents a valid mini-sudoku solution),
// the oracle applies a X gate to qubit `a`, flipping the sign of the amplitude (−1 phase).
// This marks the valid state for Grover’s amplitude amplification.
//
// Finally, the ancilla qubits `s` are uncomputed (returned to |0>) to clean up
// any entanglement and avoid side effects in the rest of the algorithm.
//
// Note: The important aspect of this oracle is that it can verify whether
// a state is a valid solution **without knowing in advance what the solution is**.
// This aligns with Grover's algorithm, which assumes only a condition-checking black box (oracle),
// not prior knowledge of the answer itself.
func oracle(qsim *q.Q, r, s []q.Qubit, a q.Qubit) {
	xor := func(x, y, z q.Qubit) {
		qsim.CNOT(x, z)
		qsim.CNOT(y, z)
	}

	xor(r[0], r[1], s[0]) // a != b
	xor(r[2], r[3], s[1]) // c != d
	xor(r[0], r[2], s[2]) // a != c
	xor(r[1], r[3], s[3]) // b != d

	// apply X if all s are 1
	qsim.ControlledX(s, []q.Qubit{a})

	// uncompute
	xor(r[1], r[3], s[3])
	xor(r[0], r[2], s[2])
	xor(r[2], r[3], s[1])
	xor(r[0], r[1], s[0])
}

func diffuser(qsim *q.Q, r []q.Qubit) {
	qsim.H(r...)
	qsim.X(r...)
	qsim.ControlledZ(r[:len(r)-1], []q.Qubit{r[len(r)-1]})
	qsim.X(r...)
	qsim.H(r...)
}

func G(qsim *q.Q, r, s []q.Qubit, a q.Qubit) {
	oracle(qsim, r, s, a)
	diffuser(qsim, r)
}

func main() {
	var top int
	flag.IntVar(&top, "top", 8, "top results")
	flag.Parse()

	qsim := q.New()

	// initialize
	r := qsim.Zeros(4)
	s := qsim.Zeros(4)
	a := qsim.Zero()

	// iteration count
	N := float64(number.Pow(2, len(r)))
	M := float64(2)                        // there are 2 solutions: [0,1,1,0] and [1,0,0,1].
	R := int(math.Pi / 4 * math.Sqrt(N/M)) // floor(pi/4 * sqrt(N/M))

	// initialize
	qsim.H(r...)
	qsim.X(a)
	qsim.H(a)

	// iterations
	for range R {
		G(qsim, r, s, a)
	}

	// quantum states
	for _, state := range q.Top(qsim.State(r, s, a), top) {
		fmt.Println(state)
	}

	// measure
	fmt.Printf("result=%v, R=%v\n", qsim.Measure(r...).BinaryString(), R)
}
