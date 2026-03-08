package density

import "github.com/itsubaki/q/quantum/qubit"

// State is a quantum state.
type State struct {
	Probability float64
	Qubit       *qubit.Qubit
}

// Normalize normalizes the probabilities of an ensemble.
func Normalize(ensemble []State) []State {
	var sum float64
	for _, s := range ensemble {
		sum += s.Probability
	}

	out := make([]State, len(ensemble))
	for i := range ensemble {
		out[i] = State{
			Probability: ensemble[i].Probability / sum,
			Qubit:       ensemble[i].Qubit,
		}
	}

	return out
}
