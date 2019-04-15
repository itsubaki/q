package main

import "github.com/axamon/q"

func main() {

	// Creates a new simulation.
	qsim := q.New()

	q0 := qsim.Zero()
	q1 := qsim.Zero()
	q2 := qsim.Zero()
	q3 := qsim.One()

	qsim.H(q0, q1, q2, q3)
	// oracle for |011>|1>
	qsim.X(q0).ControlledNot([]*q.Qubit{q0, q1, q2}, q3).X(q0)
	// amp
	qsim.H(q0, q1, q2, q3)
	qsim.X(q0, q1, q2)
	qsim.ControlledZ([]*q.Qubit{q0, q1}, q2)
	qsim.H(q0, q1, q2)

	qsim.Probability()
	// [0 0.03125 0 0.03125 0 0.03125 0 0.78125 0 0.03125 0 0.03125 0 0.03125 0 0.03125]

}
