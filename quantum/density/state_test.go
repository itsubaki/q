package density_test

import (
	"fmt"

	"github.com/itsubaki/q/quantum/density"
	"github.com/itsubaki/q/quantum/qubit"
)

func ExampleNormalize() {
	states := []density.State{
		{Probability: 1, Qubit: qubit.Zero()},
		{Probability: 2, Qubit: qubit.One()},
		{Probability: 3, Qubit: qubit.Plus()},
	}

	normalized := density.Normalize(states)
	for _, s := range normalized {
		fmt.Printf("%.4f\n", s.Probability)
	}

	// Output:
	// 0.1667
	// 0.3333
	// 0.5000
}
