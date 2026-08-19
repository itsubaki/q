package function_test

import (
	"fmt"

	"github.com/itsubaki/q"
	F "github.com/itsubaki/q/function"
)

func ExampleSwap() {
	qsim := q.New()
	q0 := qsim.Zero()
	q1 := qsim.Zero()
	q2 := qsim.One()
	q3 := qsim.One()

	F.Swap(qsim, q0, q1, q2, q3)
	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [1100] ( 1.0000 0.0000i): 1.0000
}
