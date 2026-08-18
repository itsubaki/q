package function_test

import (
	"fmt"

	"github.com/itsubaki/q"
	F "github.com/itsubaki/q/function"
)

func ExampleXOR() {
	qsim := q.New()
	x := qsim.Zero()
	y := qsim.Zero()
	z := qsim.Zero()
	qsim.H(x, y)

	F.XOR(qsim, x, y, z)
	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [000] ( 0.5000 0.0000i): 0.2500
	// [011] ( 0.5000 0.0000i): 0.2500
	// [101] ( 0.5000 0.0000i): 0.2500
	// [110] ( 0.5000 0.0000i): 0.2500
}

func ExampleXOR_z1() {
	qsim := q.New()
	x := qsim.Zero()
	y := qsim.Zero()
	z := qsim.One()
	qsim.H(x, y)

	F.XOR(qsim, x, y, z)
	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [001] ( 0.5000 0.0000i): 0.2500
	// [010] ( 0.5000 0.0000i): 0.2500
	// [100] ( 0.5000 0.0000i): 0.2500
	// [111] ( 0.5000 0.0000i): 0.2500
}
