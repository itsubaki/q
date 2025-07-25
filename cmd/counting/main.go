package main

import (
	"fmt"
	"math"
	"sort"

	"github.com/itsubaki/q"
	"github.com/itsubaki/q/math/number"
	"github.com/itsubaki/q/quantum/gate"
)

// controlledG applies the Grover operator for 2x2 mini-sudoku solutions.
// The number of solutions `M` is 2.
func controlledG(qsim *q.Q, c q.Qubit, r []q.Qubit, t []q.Qubit, a q.Qubit) {
	oracle(qsim, c, r, t, a)
	amplify(qsim, c, r)
}

func oracle(qsim *q.Q, ctrl q.Qubit, r, t []q.Qubit, a q.Qubit) {
	xor := func(x, y, t q.Qubit) {
		qsim.ControlledNot([]q.Qubit{ctrl, x}, t)
		qsim.ControlledNot([]q.Qubit{ctrl, y}, t)
	}

	xor(r[0], r[1], t[0]) // a != b
	xor(r[2], r[3], t[1]) // c != d
	xor(r[0], r[2], t[2]) // a != c
	xor(r[1], r[3], t[3]) // b != d

	// apply Z if all t are 1
	qsim.ControlledZ([]q.Qubit{ctrl, t[0], t[1], t[2], t[3]}, a)

	// uncompute
	xor(r[3], r[1], t[3])
	xor(r[2], r[0], t[2])
	xor(r[3], r[2], t[1])
	xor(r[1], r[0], t[0])
}

func amplify(qsim *q.Q, ctrl q.Qubit, r []q.Qubit) {
	qsim.Controlled(gate.H(), []q.Qubit{ctrl}, r)
	qsim.Controlled(gate.X(), []q.Qubit{ctrl}, r)
	qsim.ControlledZ([]q.Qubit{ctrl, r[0], r[1], r[2]}, r[3])
	qsim.Controlled(gate.X(), []q.Qubit{ctrl}, r)
	qsim.Controlled(gate.H(), []q.Qubit{ctrl}, r)
}

func sortedKeys(result map[int64]float64) []int64 {
	keys := make([]int64, 0, len(result))
	for k := range result {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return result[keys[i]] > result[keys[j]]
	})

	return keys
}

func main() {
	qsim := q.New()

	// initialize
	c := qsim.Zeros(3) // for phase estimation
	r := qsim.Zeros(4) // data qubits for the Grover search space
	t := qsim.Zeros(4) // ancilla qubits for comparing Sudoku constraints
	a := qsim.One()    // oracle target (phase kickback)

	// superposition
	qsim.H(c...)
	qsim.H(r...)
	qsim.H(a)

	// phase estimation
	for i := range len(c) {
		// apply controlled-G**(2**i) where control is c[len(c)-1-i]
		for range 1 << i {
			controlledG(qsim, c[len(c)-1-i], r, t, a)
		}
	}

	// inverse quantum fourier transform
	qsim.IQFT(c...)

	// results
	result := make(map[int64]float64)
	for _, s := range qsim.State(c) {
		result[s.Int()] += s.Probability()
	}

	// estimate the number of solutions `M`
	N := number.Pow(2, len(r))
	for _, k := range sortedKeys(result) {
		theta := float64(k) / float64(number.Pow(2, len(c)))   // theta = k / 2**len(c)
		M := float64(N) * math.Pow(math.Sin(math.Pi*theta), 2) // M = N * (sin(pi * theta))**2

		// result=1, prob=0.2669; theta=0.1250, M=2.3431
		// result=7, prob=0.2669; theta=0.8750, M=2.3431
		// result=5, prob=0.2019; theta=0.6250, M=13.6569
		// result=3, prob=0.2019; theta=0.3750, M=13.6569
		// result=4, prob=0.0346; theta=0.5000, M=16.0000
		// result=2, prob=0.0137; theta=0.2500, M=8.0000
		// result=6, prob=0.0137; theta=0.7500, M=8.0000
		// result=0, prob=0.0005; theta=0.0000, M=0.0000
		fmt.Printf("result=%v, prob=%.4f; theta=%.4f, M=%2.4f\n", k, result[k], theta, M)
	}
}
