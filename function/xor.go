package function

import "github.com/itsubaki/q"

func XOR(qsim *q.Q, x, y, z q.Qubit) {
	qsim.CNOT(x, z)
	qsim.CNOT(y, z)
}
