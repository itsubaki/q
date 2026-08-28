package function

import "github.com/itsubaki/q"

// XOR applies the XOR operation.
// It flips the target qubit z if either x or y is 1, but not both.
func XOR(qsim *q.Q, x, y, z q.Qubit) {
	qsim.CNOT(x, z)
	qsim.CNOT(y, z)
}
