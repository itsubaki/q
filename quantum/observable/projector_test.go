package observable_test

import (
	"testing"

	"github.com/itsubaki/q/math/matrix"
	"github.com/itsubaki/q/quantum/gate"
	"github.com/itsubaki/q/quantum/observable"
	"github.com/itsubaki/q/quantum/qubit"
)

func TestProjector(t *testing.T) {
	cases := []struct {
		in   []*qubit.Qubit
		want *matrix.Matrix
	}{
		{
			in: []*qubit.Qubit{qubit.Zero()},
			want: gate.New(
				[]complex128{1, 0},
				[]complex128{0, 0},
			),
		},
		{
			in: []*qubit.Qubit{qubit.One()},
			want: gate.New(
				[]complex128{0, 0},
				[]complex128{0, 1},
			),
		},
		{
			in: []*qubit.Qubit{
				qubit.Zero(),
				qubit.Zero(),
			},
			want: gate.New(
				[]complex128{1, 0, 0, 0},
				[]complex128{0, 0, 0, 0},
				[]complex128{0, 0, 0, 0},
				[]complex128{0, 0, 0, 0},
			),
		},
	}

	for _, c := range cases {
		got := observable.Projector(c.in...)
		if !got.Equal(c.want) {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}
