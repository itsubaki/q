package function

import "github.com/itsubaki/q"

// Fredkin applies the Fredkin gate.
// It swaps the target qubits target0 and target1 if the control qubit is 1.
func Fredkin(qsim *q.Q, control q.Qubit, target0, target1 q.Qubit) {
	qsim.CNOT(target0, target1)
	qsim.CCNOT(control, target1, target0)
	qsim.CNOT(target0, target1)
}
