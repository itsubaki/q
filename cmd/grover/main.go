package main

import (
	"fmt"
	"math"
	"sort"

	"github.com/itsubaki/q"
	"github.com/itsubaki/q/math/number"
	"github.com/itsubaki/q/quantum/qubit"
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
// The `t` slice must contain 4 ancilla qubits used for intermediate checks.
// The `a` qubit is the oracle’s phase flag (target) and should be initialized
// to (|0> − |1>)/sqrt(2) before calling this function.
//
// The oracle checks the following uniqueness constraints:
//
//   - a != b
//   - c != d
//   - a != c
//   - b != d
//
// If **all** constraints are satisfied (i.e., the input represents a valid mini-sudoku solution),
// the oracle applies a Z gate to qubit `a`, flipping the sign of the amplitude (−1 phase).
// This marks the valid state for Grover’s amplitude amplification.
//
// Finally, the ancilla qubits `t` are uncomputed (returned to |0>) to clean up
// any entanglement and avoid side effects in the rest of the algorithm.
//
// Note: The important aspect of this oracle is that it can verify whether
// a state is a valid solution **without knowing in advance what the solution is**.
// This aligns with Grover's algorithm, which assumes only a condition-checking black box (oracle),
// not prior knowledge of the answer itself.
func oracle(qsim *q.Q, r, t []q.Qubit, a q.Qubit) {
	// check a != b, c != d, a != c, b != d
	qsim.CNOT(r[0], t[0])
	qsim.CNOT(r[1], t[0])
	qsim.CNOT(r[2], t[1])
	qsim.CNOT(r[3], t[1])
	qsim.CNOT(r[0], t[2])
	qsim.CNOT(r[2], t[2])
	qsim.CNOT(r[1], t[3])
	qsim.CNOT(r[3], t[3])

	// apply Z if all t are 1
	qsim.ControlledZ(t, a)

	// uncompute
	qsim.CNOT(r[3], t[3])
	qsim.CNOT(r[1], t[3])
	qsim.CNOT(r[2], t[2])
	qsim.CNOT(r[0], t[2])
	qsim.CNOT(r[3], t[1])
	qsim.CNOT(r[2], t[1])
	qsim.CNOT(r[1], t[0])
	qsim.CNOT(r[0], t[0])
}

func amplify(qsim *q.Q, r []q.Qubit) {
	qsim.H(r...)
	qsim.X(r...)
	qsim.ControlledZ([]q.Qubit{r[0], r[1], r[2]}, r[3])
	qsim.X(r...)
	qsim.H(r...)
}

func top(s []qubit.State, n int) []qubit.State {
	sort.Slice(s, func(i, j int) bool {
		return s[i].Probability() > s[j].Probability()
	})

	if len(s) < n {
		return s
	}

	return s[:n]
}

func main() {
	qsim := q.New()

	// initialize
	r := qsim.Zeros(4)
	t := qsim.Zeros(4)
	a := qsim.One()

	// superposition
	qsim.H(r...).H(a)

	// iteration count
	N := float64(number.Pow(2, len(r)))
	M := float64(2)                        // there are 2 solutions: [0,1,1,0] and [1,0,0,1].
	R := int(math.Pi / 4 * math.Sqrt(N/M)) // floor(pi/4 * sqrt(N/M))

	// iterations
	for range R {
		oracle(qsim, r, t, a)
		amplify(qsim, r)
	}

	for _, s := range top(qsim.State(r), 5) {
		// [0110][  6](-0.4861 0.0000i): 0.2363
		// [1001][  9](-0.4861 0.0000i): 0.2363
		// [1000][  8]( 0.1768 0.0000i): 0.0313
		// [1011][ 11]( 0.1768 0.0000i): 0.0313
		// [0010][  2]( 0.1768 0.0000i): 0.0313
		fmt.Println(s)
	}

	fmt.Printf("result=%v, R=%v\n", qsim.Measure(r...).BinaryString(), R)
}
