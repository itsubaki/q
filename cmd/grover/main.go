package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"sort"

	"github.com/itsubaki/q"
	"github.com/itsubaki/q/math/number"
	"github.com/itsubaki/q/quantum/qubit"
)

// Oracle construction adapted from:
// C. Figgatt, D. Maslov, K. A. Landsman, N. M. Linke, S. Debnath, and C. Monroe.
// Complete 3-Qubit Grover Search on a Programmable Quantum Computer.
var oracleMap = map[string]func(qsim *q.Q, r []q.Qubit, a q.Qubit){
	"000": func(qsim *q.Q, r []q.Qubit, a q.Qubit) {
		qsim.X(r...)
		qsim.CCCNOT(r[0], r[1], r[2], a)
		qsim.X(r...)
	},
	"001": func(qsim *q.Q, r []q.Qubit, a q.Qubit) {
		qsim.X(r[0], r[1])
		qsim.CCCNOT(r[0], r[1], r[2], a)
		qsim.X(r[0], r[1])
	},
	"010": func(qsim *q.Q, r []q.Qubit, a q.Qubit) {
		qsim.X(r[0], r[2])
		qsim.CCCNOT(r[0], r[1], r[2], a)
		qsim.X(r[0], r[2])
	},
	"011": func(qsim *q.Q, r []q.Qubit, a q.Qubit) {
		qsim.X(r[0])
		qsim.CCCNOT(r[0], r[1], r[2], a)
		qsim.X(r[0])
	},
	"100": func(qsim *q.Q, r []q.Qubit, a q.Qubit) {
		qsim.X(r[1], r[2])
		qsim.CCCNOT(r[0], r[1], r[2], a)
		qsim.X(r[1], r[2])
	},
	"101": func(qsim *q.Q, r []q.Qubit, a q.Qubit) {
		qsim.X(r[1])
		qsim.CCCNOT(r[0], r[1], r[2], a)
		qsim.X(r[1])
	},
	"110": func(qsim *q.Q, r []q.Qubit, a q.Qubit) {
		qsim.X(r[2])
		qsim.CCCNOT(r[0], r[1], r[2], a)
		qsim.X(r[2])
	},
	"111": func(qsim *q.Q, r []q.Qubit, a q.Qubit) {
		qsim.CCCNOT(r[0], r[1], r[2], a)
	},
}

func top(s []qubit.State, n int) []qubit.State {
	sort.Slice(s, func(i, j int) bool { return s[i].Probability() > s[j].Probability() })
	if len(s) < n {
		return s
	}

	return s[:n]
}

func main() {
	var oracle string
	flag.StringVar(&oracle, "oracle", "011", "oracle function in binary string")
	flag.Parse()

	// oracle
	ora, ok := oracleMap[oracle]
	if !ok {
		fmt.Printf("oracle=%q not found\n", oracle)
		os.Exit(1)
	}

	// initial state
	qsim := q.New()
	r := qsim.Zeros(3)
	a := qsim.One()

	// superposition
	qsim.H(r...).H(a)

	// iterations
	N := number.Pow(2, qsim.NumQubits())
	ite := math.Floor(math.Pi / 4 * math.Sqrt(float64(N)))
	for range int(ite) {
		// oracle
		ora(qsim, r, a)

		// amplification
		qsim.H(r...).H(a)
		qsim.X(r...)
		qsim.CCZ(r[0], r[1], r[2])
		qsim.X(r...)
		qsim.H(r...)
	}

	state := qsim.State([]q.Qubit{r[0], r[1], r[2]}, a)
	for _, s := range top(state, 10) {
		fmt.Println(s)
	}

	fmt.Println("result:", qsim.Measure(r...).BinaryString())
}
