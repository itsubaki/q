package density

import "github.com/itsubaki/q/quantum/qubit"

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
