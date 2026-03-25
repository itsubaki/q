package density_test

import (
	"fmt"

	"github.com/itsubaki/q/quantum/density"
	"github.com/itsubaki/q/quantum/qubit"
)

func ExampleChannel() {
	rho := density.NewPureState(qubit.Zero())
	s0 := rho.ApplyChannel([]*density.Channel{
		density.Depolarizing(0.1, rho.NumQubits(), 0),
		density.AmplitudeDamping(0.7, rho.NumQubits(), 0),
		density.PhaseDamping(0.7, rho.NumQubits(), 0),
		density.BitFlip(0.1, rho.NumQubits(), 0),
	}...)

	fmt.Printf("%.4f\n", s0.Probability(qubit.Zero()))

	// Output:
	// 0.8840
}
