package main

import "github.com/axamon/q"

func main() {

	qsim := q.New()

	q0 := qsim.New(1, 2) // (0.2, 0.8)

	// encoding
	q1 := qsim.Zero()
	q2 := qsim.Zero()
	qsim.CNOT(q0, q1).CNOT(q0, q2)

	// error: first qubit is flipped
	qsim.X(q0)

	// add ancilla qubit
	q3 := qsim.Zero()
	q4 := qsim.Zero()

	// error corretion
	qsim.CNOT(q0, q3).CNOT(q1, q3)
	qsim.CNOT(q1, q4).CNOT(q2, q4)

	m3 := qsim.Measure(q3)
	m4 := qsim.Measure(q4)

	qsim.ConditionX(m3.IsOne() && m4.IsZero(), q0)
	qsim.ConditionX(m3.IsOne() && m4.IsOne(), q1)
	qsim.ConditionX(m3.IsZero() && m4.IsOne(), q2)

	// estimate
	qsim.Estimate(q0).Probability() // (0.2, 0.8)
	qsim.Estimate(q1).Probability() // (0.2, 0.8)
	qsim.Estimate(q2).Probability() // (0.2, 0.8)

}
