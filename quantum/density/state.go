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

	for i := range ensemble {
		ensemble[i].Probability = ensemble[i].Probability / sum
	}

	return ensemble
}
