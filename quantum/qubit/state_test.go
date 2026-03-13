package qubit_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/itsubaki/q/quantum/qubit"
)

func ExampleState() {
	s := qubit.NewState(complex(1, 0), []string{"0100", "1010", "1000"}...)

	fmt.Println(s.BinaryString(), s.Int(), s.Amplitude(), s.Probability())
	fmt.Println(s.String())

	// Output:
	// [0100 1010 1000] [4 10 8] (1+0i) 1
	// [0100 1010 1000][  4  10   8]( 1.0000 0.0000i): 1.0000
}

func TestState_Equal(t *testing.T) {
	cases := []struct {
		s    qubit.State
		v    qubit.State
		want bool
	}{
		{
			qubit.NewState(complex(1, 0)),
			qubit.State{},
			false,
		},
		{
			qubit.NewState(complex(1, 0), "001"),
			qubit.NewState(complex(1, 0), "100"),
			false,
		},
		{
			qubit.NewState(complex(0, 1), "001"),
			qubit.NewState(complex(1, 0), "001"),
			false,
		},
		{
			qubit.NewState(complex(1, 0), "001", "010"),
			qubit.NewState(complex(1, 0), "001"),
			false,
		},
		{
			qubit.NewState(complex(1, 0), "001", "010"),
			qubit.NewState(complex(1, 0), "001", "010"),
			true,
		},
	}

	for _, c := range cases {
		got := c.s.Equal(c.v)
		if got != c.want {
			t.Errorf("got=%v want=%v", got, c.want)
		}
	}
}

func TestEqual(t *testing.T) {
	cases := []struct {
		s    []qubit.State
		v    []qubit.State
		want bool
	}{
		{
			[]qubit.State{},
			[]qubit.State{},
			true,
		},
		{
			[]qubit.State{{}},
			[]qubit.State{{}, {}},
			false,
		},
		{
			[]qubit.State{
				qubit.NewState(1, "000"),
			},
			[]qubit.State{
				qubit.NewState(1/math.Sqrt2, "000"),
				qubit.NewState(1/math.Sqrt2, "111"),
			},
			false,
		},
		{
			[]qubit.State{
				qubit.NewState(complex(1, 0), "000"),
			},
			[]qubit.State{
				qubit.NewState(complex(0, 1), "000"),
			},
			false,
		},
	}

	for _, c := range cases {
		got := qubit.Equal(c.s, c.v)
		if got != c.want {
			t.Errorf("got=%v want=%v", got, c.want)
		}
	}
}

func TestEqualUpToGlobalPhase(t *testing.T) {
	cases := []struct {
		s    []qubit.State
		v    []qubit.State
		want bool
	}{
		{
			s:    []qubit.State{qubit.NewState(0.5)},
			v:    []qubit.State{},
			want: false,
		},
		{
			s:    []qubit.State{qubit.NewState(complex(1, 0))},
			v:    []qubit.State{qubit.NewState(complex(1, 0))},
			want: true,
		},
		{
			s:    []qubit.State{qubit.NewState(complex(1, 0))},
			v:    []qubit.State{qubit.NewState(complex(0, 1))},
			want: true,
		},
	}

	for _, c := range cases {
		got := qubit.EqualUpToGlobalPhase(c.s, c.v)
		if got != c.want {
			t.Errorf("got=%v want=%v", got, c.want)
		}
	}
}
