package qubit_test

import (
	"fmt"

	"github.com/itsubaki/q/pkg/quantum/qubit"
)

func ExampleState() {
	s := qubit.State{
		Amplitude:    complex(1, 0),
		Probability:  1,
		Int:          []int{4, 10},
		BinaryString: []string{"0100", "1010"},
	}

	fmt.Println(s.Value())
	fmt.Println(s.Value(0))
	fmt.Println(s.Value(1))
	fmt.Println(s.String())

	// Output:
	// 4 0100
	// 4 0100
	// 10 1010
	// [0100 1010][  4  10]( 1.0000 0.0000i): 1.0000
}
