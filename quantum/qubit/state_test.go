package qubit_test

import (
	"fmt"
	"testing"

	"github.com/itsubaki/q/quantum/qubit"
)

func ExampleState() {
	s := qubit.NewState(complex(1, 0), []string{"0100", "1010", "1000"}...)

	fmt.Println(s.Int(), s.BinaryString())
	fmt.Println(s.Int(0), s.BinaryString(0))
	fmt.Println(s.Int(1), s.BinaryString(1))
	fmt.Println(s.Int(2), s.BinaryString(2))
	fmt.Println(s.String())
	fmt.Println(s.Amplitude(), s.Probability())

	// Output:
	// 4 0100
	// 4 0100
	// 10 1010
	// 8 1000
	// [0100 1010 1000][  4  10   8]( 1.0000 0.0000i): 1.0000
	// (1+0i) 1
}

func TestState_Equals(t *testing.T) {
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
		got := c.s.Equals(c.v)
		if got != c.want {
			t.Errorf("got=%v want=%v", got, c.want)
		}
	}
}

func TestEquals(t *testing.T) {
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
			[]qubit.State{qubit.NewState(complex(1, 0))},
			[]qubit.State{qubit.NewState(complex(0, 1))},
			false,
		},
	}

	for _, c := range cases {
		got := qubit.Equals(c.s, c.v)
		if got != c.want {
			t.Errorf("got=%v want=%v", got, c.want)
		}
	}
}
