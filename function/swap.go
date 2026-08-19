package function

import "github.com/itsubaki/q"

// Swap applies the swap gate to the given qubits.
func Swap(qsim *q.Q, qb ...q.Qubit) {
	l := len(qb)
	for i := range l / 2 {
		qsim.Swap(qb[i], qb[(l-1)-i])
	}
}
