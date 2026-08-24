package function_test

import (
	"fmt"

	"github.com/itsubaki/q"
	F "github.com/itsubaki/q/function"
)

func ExampleFredkin() {
	qsim := q.New()
	c := qsim.Zero()
	t0 := qsim.Zero()
	t1 := qsim.Zero()
	qsim.H(c)
	qsim.X(t0)

	F.Fredkin(qsim, c, t0, t1)
	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [010] ( 0.7071 0.0000i): 0.5000
	// [101] ( 0.7071 0.0000i): 0.5000
}
