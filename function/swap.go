package function

import "github.com/itsubaki/q"

// Swap applies the swap gate to the given qubits.
func Swap(qsim *q.Q, qb ...q.Qubit) {
	l := len(qb)

	for i := range l / 2 {
		q0, q1 := qb[i], qb[(l-1)-i]

		qsim.CNOT(q0, q1)
		qsim.CNOT(q1, q0)
		qsim.CNOT(q0, q1)
	}
}
