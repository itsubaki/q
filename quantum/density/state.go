package density

import (
	"github.com/itsubaki/q/math/epsilon"
	"github.com/itsubaki/q/quantum/qubit"
)

// State is a quantum state.
type State struct {
	Probability float64
	Qubit       *qubit.Qubit
}

// Normalize normalizes the probabilities of a set of states.
func Normalize(states []State) []State {
	var sum float64
	for _, s := range states {
		sum += s.Probability
	}

	out := make([]State, len(states))
	for i := range states {
		out[i] = State{
			Probability: states[i].Probability / sum,
			Qubit:       states[i].Qubit,
		}
	}

	return out
}

// IsValid checks if the given set of states is valid for constructing a density matrix.
// A valid set must satisfy the following conditions:
// 1. The set must not be empty.
// 2. All qubits in the set must have the same dimension.
// 3. All probabilities in the set must be non-negative.
// 4. The sum of probabilities in the set must be equal to 1 (within a specified tolerance).
func IsValid(states []State, tol ...float64) bool {
	if len(states) == 0 {
		return false
	}

	n := states[0].Qubit.Dim()
	for _, s := range states {
		if s.Qubit.Dim() != n {
			return false
		}
	}

	for _, s := range states {
		if s.Probability < 0 {
			return false
		}
	}

	var sum float64
	for _, s := range states {
		sum += s.Probability
	}

	return epsilon.IsOneF64(sum, tol...)
}
