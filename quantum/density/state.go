package density

import (
	"github.com/itsubaki/q/math/matrix"
	"github.com/itsubaki/q/quantum/qubit"
)

// WeightedState is a quantum state with an associated probability.
type WeightedState struct {
	Probability float64
	Qubit       *qubit.Qubit
}

func (s WeightedState) DensityOperator() *matrix.Matrix {
	op := s.Qubit.OuterProduct(s.Qubit)
	return op.Mul(complex(s.Probability, 0))
}

// Normalize normalizes the probabilities of a set of states.
func Normalize(states []WeightedState) []WeightedState {
	var sum float64
	for _, s := range states {
		sum += s.Probability
	}

	out := make([]WeightedState, len(states))
	for i := range states {
		out[i] = WeightedState{
			Probability: states[i].Probability / sum,
			Qubit:       states[i].Qubit,
		}
	}

	return out
}
