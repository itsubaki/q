package main

import (
	"fmt"

	"github.com/axamon/q/q"
)

func main() {

	qsim := q.New()

	// generate qubits of |0>|0>
	q0 := qsim.Zero()
	q1 := qsim.Zero()
	q2 := qsim.One()

	// apply quantum circuit
	qsim.H(q0).CNOT(q0, q1)

	// estimate
	qsim.Estimate(q0).Probability()
	// -> (0.5, 0.5)
	qsim.Estimate(q1).Probability()
	// -> (0.5, 0.5)

	qsim.Probability()
	// -> (0.5, 0, 0, 0.5)

	qsim.Measure()
	qsim.Probability()
	// -> (1, 0, 0, 0) or (0, 0, 0, 1)

	m0 := qsim.Measure(q0)
	m1 := qsim.Measure(q1)
	// -> m0 = |0> then m1 = |0>
	// -> m0 = |1> then m1 = |1>

	fmt.Println(m0, m1)

}
