package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/axamon/q"
	"github.com/axamon/q/number"
)

func main() {

	N := 15
	a := 7

	// a is co-prime
	if number.GCD(N, a) != 1 {
		log.Fatalf("%v %v\n", N, a)
	}

	// Creates a new simulation.
	qsim := q.New()

	q0 := qsim.Zero()
	q1 := qsim.Zero()
	q2 := qsim.Zero()

	q3 := qsim.Zero()
	q4 := qsim.Zero()
	q5 := qsim.Zero()
	q6 := qsim.One()

	// superposition
	qsim.H(q0, q1, q2)

	// Controlled-U
	qsim.CNOT(q2, q4)
	qsim.CNOT(q2, q5)

	// Controlled-U^2
	qsim.ControlledNot([]*q.Qubit{q1, q4}, q6)
	qsim.ControlledNot([]*q.Qubit{q1, q6}, q4)
	qsim.ControlledNot([]*q.Qubit{q1, q4}, q6)

	qsim.ControlledNot([]*q.Qubit{q1, q3}, q5)
	qsim.ControlledNot([]*q.Qubit{q1, q5}, q3)
	qsim.ControlledNot([]*q.Qubit{q1, q3}, q5)

	// QFT
	qsim.H(q0)
	qsim.CR(q1, q0, 2)
	qsim.CR(q2, q0, 3)

	qsim.H(q1)
	qsim.CR(q2, q1, 2)

	qsim.H(q2)

	qsim.Swap(q0, q2)

	// measure q0, q1, q2
	qsim.Measure(q0)
	qsim.Measure(q1)
	qsim.Measure(q2)

	p := qsim.Probability()
	for i := range p {
		if p[i] == 0 {
			continue
		}
		fmt.Printf("%07s %v\n", strconv.FormatInt(int64(i), 2), p[i])
	}
	// 010,0001(1)  0.25 -> 1/16
	// 010,0100(4)  0.25 -> 4/16 -> 1/4
	// 010,0111(7)  0.25 -> 7/16
	// 010,1101(13) 0.25 -> 13/16
	// r = 16 is trivial. r < N.
	// r -> 4

	// gcd(a^(r/2)-1, N) -> gcd(7^(4/2)-1, 15)
	// gcd(a^(r/2)+1, N) -> gcd(7^(4/2)+1, 15)
	p0 := number.GCD(a*a-1, N)
	p1 := number.GCD(a*a+1, N)
	if p0 != 3 || p1 != 5 {
		log.Fatalf("%v %v\n", p0, p1)
	}

	if p0 == 3 && p1 == 5 {
		log.Println("ok")
	}

}
