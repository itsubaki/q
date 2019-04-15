package main

import (
	"github.com/axamon/q"
)

func main() {

	// Creates a new simulation.
	qsim := q.New()

	// generate qubits of |phi>|0>|0>
	// |phi> is normalized. |phi> = a|0> + b|1>, |a|^2 = 0.2, |b|^2 = 0.8
	phi := qsim.New(1, 2)
	q0 := qsim.Zero()
	q1 := qsim.Zero()

	qsim.H(q0).CNOT(q0, q1)
	qsim.CNOT(phi, q0).H(phi)

	// Alice send mz, mx to Bob
	mz := qsim.Measure(phi)
	mx := qsim.Measure(q0)

	// Bob Apply Z and X
	qsim.ConditionZ(mz.IsOne(), q1)
	qsim.ConditionX(mx.IsOne(), q1)

	// Bob got |phi> state
	qsim.Estimate(q1).Probability()
	// -> (0.2, 0.8)

}
