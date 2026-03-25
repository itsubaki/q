package density_test

import (
	"fmt"
	"testing"

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

func TestIsValid(t *testing.T) {
	type Case struct {
		s    []density.State
		want bool
	}

	cases := []Case{
		{
			s: []density.State{},
		},
		{
			s: []density.State{
				{
					Probability: 1,
					Qubit:       qubit.Zero(),
				},
			},
			want: true,
		},
		{
			s: []density.State{
				{
					Probability: 0.5,
					Qubit:       qubit.Zero(),
				},
				{
					Probability: 0.5,
					Qubit:       qubit.One(),
				},
			},
			want: true,
		},
		{
			s: []density.State{
				{
					Probability: -0.1,
					Qubit:       qubit.Zero(),
				},
				{
					Probability: 1.1,
					Qubit:       qubit.One(),
				},
			},
			want: false,
		},
		{
			s: []density.State{
				{
					Probability: 0.5,
					Qubit:       qubit.Zero(),
				},
				{
					Probability: 0.5,
					Qubit:       qubit.Zero(2),
				},
			},
			want: false,
		},
		{
			s: []density.State{
				{
					Probability: 0.5,
					Qubit:       qubit.Zero(),
				},
				{
					Probability: 0.4,
					Qubit:       qubit.Zero(),
				},
			},
			want: false,
		},
	}

	for _, c := range cases {
		got := density.IsValid(c.s)
		if got != c.want {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}
