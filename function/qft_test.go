package function_test

import (
	"fmt"

	"github.com/itsubaki/q"
	F "github.com/itsubaki/q/function"
	"github.com/itsubaki/q/math/number"
	"github.com/itsubaki/q/math/rand"
)

func ExampleQFT() {
	qsim := q.New()
	q0 := qsim.Zero()
	q1 := qsim.Zero()
	q2 := qsim.One()

	F.QFT(qsim, q0, q1, q2)
	F.Swap(qsim, q0, q2)

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [000] ( 0.3536 0.0000i): 0.1250
	// [001] ( 0.2500 0.2500i): 0.1250
	// [010] ( 0.0000 0.3536i): 0.1250
	// [011] (-0.2500 0.2500i): 0.1250
	// [100] (-0.3536 0.0000i): 0.1250
	// [101] (-0.2500-0.2500i): 0.1250
	// [110] ( 0.0000-0.3536i): 0.1250
	// [111] ( 0.2500-0.2500i): 0.1250
}

func ExampleInvQFT() {
	qsim := q.New()
	q0 := qsim.Zero()
	q1 := qsim.One()
	q2 := qsim.Zero()

	F.QFT(qsim, q0, q1, q2)
	F.InvQFT(qsim, q0, q1, q2)

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [010] ( 1.0000 0.0000i): 1.0000
}

func Example_shor15() {
	// Reference: Zhengjun Cao, Zhenfu Cao, Lihua Liu. Remarks on Quantum Modular Exponentiation and Some Experimental Demonstrations of Shor's Algorithm.
	N := 15
	a := 7

	qsim := q.New()
	qsim.SetRand(rand.Const())

	q0 := qsim.Zero()
	q1 := qsim.Zero()
	q2 := qsim.Zero()

	q3 := qsim.Zero()
	q4 := qsim.Zero()
	q5 := qsim.Zero()
	q6 := qsim.One()

	qsim.H(q0, q1, q2)

	// Controlled-U^(2^0)
	qsim.CNOT(q2, q4)
	qsim.CNOT(q2, q5)

	// Controlled-U^(2^1)
	qsim.CNOT(q3, q5)
	qsim.CCNOT(q1, q5, q3)
	qsim.CNOT(q3, q5)

	qsim.CNOT(q6, q4)
	qsim.CCNOT(q1, q4, q6)
	qsim.CNOT(q6, q4)

	// inverse QFT
	F.Swap(qsim, q0, q2)
	F.InvQFT(qsim, q0, q1, q2)

	qsim.Measure(q3, q4, q5, q6)
	for _, s := range qsim.State([]q.Qubit{q0, q1, q2}) {
		fmt.Println(s)
	}

	// measure q0, q1, q2
	m := qsim.Measure(q0, q1, q2)
	k := number.MustParseInt(m.BinaryString())
	phi := number.Ldexp(k, -m.NumQubits())

	// find s/r. 0.010 -> 0.25 -> 1/4, 0.110 -> 0.75 -> 3/4, ...
	s, r, d, ok := number.FindOrder(a, N, phi)
	if !ok || number.IsOdd(r) {
		return
	}

	// gcd(a^(r/2)-1, N), gcd(a^(r/2)+1, N)
	p0 := number.GCD(number.Pow(a, r/2)-1, N)
	p1 := number.GCD(number.Pow(a, r/2)+1, N)
	if number.IsTrivial(N, p0, p1) {
		return
	}

	fmt.Printf("N=%d, a=%d. p=%v, q=%v. s/r=%d/%d ([0.%v]~%.3f)\n", N, a, p0, p1, s, r, m.BinaryString(), d)

	// Output:
	// [000] ( 0.5000 0.0000i): 0.2500
	// [010] ( 0.0000 0.5000i): 0.2500
	// [100] (-0.5000 0.0000i): 0.2500
	// [110] ( 0.0000-0.5000i): 0.2500
	// N=15, a=7. p=3, q=5. s/r=1/4 ([0.010]~0.250)
}
