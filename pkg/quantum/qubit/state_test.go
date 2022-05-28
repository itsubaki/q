package qubit_test

import (
	"fmt"
	"testing"

	"github.com/itsubaki/q/pkg/quantum/qubit"
)

func ExampleState() {
	s := qubit.State{
		Amplitude:    1,
		Probability:  1,
		Int:          []int64{4, 10, 8},
		BinaryString: []string{"0100", "1010", "1000"},
	}

	fmt.Println(s.Value())
	fmt.Println(s.Value(0))
	fmt.Println(s.Value(1))
	fmt.Println(s.Value(2))
	fmt.Println(s.String())

	// Output:
	// 4 0100
	// 4 0100
	// 10 1010
	// 8 1000
	// [0100 1010 1000][  4  10   8]( 1.0000 0.0000i): 1.0000
}

func TestState_Equals(t *testing.T) {
	cases := []struct {
		s    qubit.State
		v    qubit.State
		want bool
	}{
		{
			qubit.State{Int: []int64{1}},
			qubit.State{},
			false,
		},
		{
			qubit.State{Int: []int64{1}},
			qubit.State{Int: []int64{1}, BinaryString: []string{"001"}},
			false,
		},
		{
			qubit.State{Int: []int64{1}, BinaryString: []string{"001"}, Amplitude: complex(1, 1)},
			qubit.State{Int: []int64{1}, BinaryString: []string{"001"}, Amplitude: complex(1, 2)},
			false,
		},
		{
			qubit.State{Int: []int64{1}, BinaryString: []string{"001"}, Amplitude: complex(1, 1)},
			qubit.State{Int: []int64{1}, BinaryString: []string{"001"}, Amplitude: complex(1, 1)},
			true,
		},
	}

	for _, c := range cases {
		got := c.s.Equals(c.v)
		if got != c.want {
			t.Errorf("got=%v want=%v", got, c.want)
		}
	}
}
