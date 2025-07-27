package main

import (
	"fmt"

	"github.com/itsubaki/q"
)

func Example_controlledModExp2mod15() {
	qsim := q.New()
	c := qsim.Zero()
	t := qsim.ZeroLog2(15)

	qsim.X(c)
	qsim.X(t[len(t)-1])
	for _, s := range qsim.State(c, t) {
		fmt.Println(s)
	}

	a, N := 7, 15
	for i := range 3 {
		ControlledModExp2(qsim, a, i, N, c, t)
		for _, s := range qsim.State(c, t) {
			fmt.Println(s)
		}
	}

	// Output:
	// [1 0001][  1   1]( 1.0000 0.0000i): 1.0000
	// [1 0111][  1   7]( 1.0000 0.0000i): 1.0000
	// [1 1101][  1  13]( 1.0000 0.0000i): 1.0000
	// [1 1101][  1  13]( 1.0000 0.0000i): 1.0000
}

func Example_controlledModExp2mod21() {
	qsim := q.New()
	c := qsim.Zero()
	t := qsim.ZeroLog2(21)

	qsim.X(c)
	qsim.X(t[len(t)-1])

	for _, s := range qsim.State(c, t) {
		fmt.Println(s)
	}

	a, N := 2, 21
	for i := range 4 {
		ControlledModExp2(qsim, a, i, N, c, t)
		for _, s := range qsim.State(c, t) {
			fmt.Println(s)
		}
	}

	// Output:
	// [1 00001][  1   1]( 1.0000 0.0000i): 1.0000
	// [1 00010][  1   2]( 1.0000 0.0000i): 1.0000
	// [1 01000][  1   8]( 1.0000 0.0000i): 1.0000
	// [1 00010][  1   2]( 1.0000 0.0000i): 1.0000
	// [1 01000][  1   8]( 1.0000 0.0000i): 1.0000
}
