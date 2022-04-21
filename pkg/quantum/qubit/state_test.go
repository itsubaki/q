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
	// 4 0100 <nil>
	// 4 0100 <nil>
	// 10 1010 <nil>
	// 8 1000 <nil>
	// [0100 1010 1000][  4  10   8]( 1.0000 0.0000i): 1.0000
}

func TestValue(t *testing.T) {
	s := qubit.State{
		Amplitude:    1,
		Probability:  1,
		Int:          []int64{4, 10, 8},
		BinaryString: []string{"0100", "1010", "1000"},
	}

	cases := []struct {
		in   []int
		want error
	}{
		{[]int{1, 1}, fmt.Errorf("invalid parameter. len(index)=2")},
		{[]int{-1}, fmt.Errorf("invalid parameter. index=[-1]")},
	}

	for _, c := range cases {
		_, _, got := s.Value(c.in...)
		if got.Error() != c.want.Error() {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}
