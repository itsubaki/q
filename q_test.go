package q_test

import (
	"fmt"
	"math"

	"github.com/itsubaki/q"
	"github.com/itsubaki/q/math/matrix"
	"github.com/itsubaki/q/math/number"
	"github.com/itsubaki/q/math/rand"
	"github.com/itsubaki/q/quantum/gate"
)

func ExampleQ_Zero() {
	qsim := q.New()

	q0 := qsim.Zero()
	q1 := qsim.Zero()
	qsim.H(q0, q1)

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [00][  0]( 0.5000 0.0000i): 0.2500
	// [01][  1]( 0.5000 0.0000i): 0.2500
	// [10][  2]( 0.5000 0.0000i): 0.2500
	// [11][  3]( 0.5000 0.0000i): 0.2500
}

func ExampleQ_Zeros() {
	qsim := q.New()

	r := qsim.Zeros(2)
	qsim.H(r...)

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [00][  0]( 0.5000 0.0000i): 0.2500
	// [01][  1]( 0.5000 0.0000i): 0.2500
	// [10][  2]( 0.5000 0.0000i): 0.2500
	// [11][  3]( 0.5000 0.0000i): 0.2500
}

func ExampleQ_ZeroLog2() {
	qsim := q.New()

	r := qsim.ZeroLog2(15)
	qsim.H(r...)

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [0000][  0]( 0.2500 0.0000i): 0.0625
	// [0001][  1]( 0.2500 0.0000i): 0.0625
	// [0010][  2]( 0.2500 0.0000i): 0.0625
	// [0011][  3]( 0.2500 0.0000i): 0.0625
	// [0100][  4]( 0.2500 0.0000i): 0.0625
	// [0101][  5]( 0.2500 0.0000i): 0.0625
	// [0110][  6]( 0.2500 0.0000i): 0.0625
	// [0111][  7]( 0.2500 0.0000i): 0.0625
	// [1000][  8]( 0.2500 0.0000i): 0.0625
	// [1001][  9]( 0.2500 0.0000i): 0.0625
	// [1010][ 10]( 0.2500 0.0000i): 0.0625
	// [1011][ 11]( 0.2500 0.0000i): 0.0625
	// [1100][ 12]( 0.2500 0.0000i): 0.0625
	// [1101][ 13]( 0.2500 0.0000i): 0.0625
	// [1110][ 14]( 0.2500 0.0000i): 0.0625
	// [1111][ 15]( 0.2500 0.0000i): 0.0625
}

func ExampleQ_Ones() {
	qsim := q.New()

	r := qsim.Ones(2)
	qsim.H(r...)

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [00][  0]( 0.5000 0.0000i): 0.2500
	// [01][  1](-0.5000 0.0000i): 0.2500
	// [10][  2](-0.5000 0.0000i): 0.2500
	// [11][  3]( 0.5000 0.0000i): 0.2500
}

func ExampleQ_Reset() {
	qsim := q.New()

	r := qsim.Zeros(2)
	qsim.Reset(r...)

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	qsim.X(r...)
	qsim.Reset(r...)

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [00][  0]( 1.0000 0.0000i): 1.0000
	// [00][  0]( 1.0000 0.0000i): 1.0000
}

func ExampleQ_Amplitude() {
	qsim := q.New()

	q0 := qsim.Zero()
	q1 := qsim.Zero()

	qsim.H(q0)
	qsim.CNOT(q0, q1)

	for _, a := range qsim.Amplitude() {
		fmt.Println(a)
	}

	// Output:
	// (0.7071067811865476+0i)
	// (0+0i)
	// (0+0i)
	// (0.7071067811865476+0i)
}

func ExampleQ_Probability() {
	qsim := q.New()

	q0 := qsim.Zero()
	q1 := qsim.Zero()

	qsim.H(q0)
	qsim.CNOT(q0, q1)

	for _, p := range qsim.Probability() {
		fmt.Printf("%.4f\n", p)
	}

	// Output:
	// 0.5000
	// 0.0000
	// 0.0000
	// 0.5000
}

func ExampleQ_Measure() {
	qsim := q.New()
	qsim.Rand = rand.Const()

	q0 := qsim.Zero()
	q1 := qsim.Zero()
	q2 := qsim.Zero()
	qsim.X(q0)

	fmt.Println(qsim.Measure(q0))
	fmt.Println(qsim.Measure(q0, q1, q2))
	fmt.Println(qsim.Measure())

	// Output:
	// [(0+0i) (1+0i)]
	// [(0+0i) (0+0i) (0+0i) (0+0i) (1+0i) (0+0i) (0+0i) (0+0i)]
	// [(0+0i) (0+0i) (0+0i) (0+0i) (1+0i) (0+0i) (0+0i) (0+0i)]
}

func ExampleQ_M() {
	qsim := q.New()
	qsim.Rand = rand.Const()

	q0 := qsim.Zero()
	q1 := qsim.Zero()
	q2 := qsim.Zero()
	qsim.X(q0)

	fmt.Println(qsim.M(q0))
	fmt.Println(qsim.M(q0, q1, q2))
	fmt.Println(qsim.M())

	// Output:
	// [(0+0i) (1+0i)]
	// [(0+0i) (0+0i) (0+0i) (0+0i) (1+0i) (0+0i) (0+0i) (0+0i)]
	// [(0+0i) (0+0i) (0+0i) (0+0i) (1+0i) (0+0i) (0+0i) (0+0i)]
}

func ExampleQ_Apply() {
	qsim := q.New()

	q0 := qsim.Zero()
	q1 := qsim.Zero()
	n := qsim.NumQubits()

	qsim.Apply(gate.H(), q0)
	qsim.Apply(gate.CNOT(n, q0.Index(), q1.Index()))

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [00][  0]( 0.7071 0.0000i): 0.5000
	// [11][  3]( 0.7071 0.0000i): 0.5000
}

func ExampleQ_U() {
	qsim := q.New()

	q0 := qsim.Zero()
	qsim.U(math.Pi, 0, math.Pi, q0)

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [1][  1]( 1.0000 0.0000i): 1.0000
}

func ExampleQ_I() {
	qsim := q.New()

	q0 := qsim.Zero()
	qsim.I(q0)

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [0][  0]( 1.0000 0.0000i): 1.0000
}

func ExampleQ_X() {
	qsim := q.New()

	q0 := qsim.Zero()
	qsim.X(q0)

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [1][  1]( 1.0000 0.0000i): 1.0000
}

func ExampleQ_Y() {
	qsim := q.New()

	q0 := qsim.One()
	qsim.Y(q0)

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [0][  0]( 0.0000-1.0000i): 1.0000
}

func ExampleQ_Z() {
	qsim := q.New()

	q0 := qsim.One()
	qsim.Z(q0)

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [1][  1](-1.0000 0.0000i): 1.0000
}

func ExampleQ_S() {
	qsim := q.New()

	q0 := qsim.One()
	qsim.S(q0)

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [1][  1]( 0.0000 1.0000i): 1.0000
}

func ExampleQ_T() {
	qsim := q.New()

	q0 := qsim.One()
	qsim.T(q0)

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [1][  1]( 0.7071 0.7071i): 1.0000
}

func ExampleQ_R() {
	qsim := q.New()

	q0 := qsim.One()
	qsim.R(2*math.Pi/4, q0)

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [1][  1]( 0.0000 1.0000i): 1.0000
}

func ExampleQ_RX() {
	qsim := q.New()

	q0 := qsim.Zero()
	qsim.RX(math.Pi, q0)

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [1][  1]( 0.0000-1.0000i): 1.0000
}

func ExampleQ_RY() {
	qsim := q.New()

	q0 := qsim.Zero()
	qsim.RY(math.Pi, q0)

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [1][  1]( 1.0000 0.0000i): 1.0000
}

func ExampleQ_RZ() {
	qsim := q.New()

	q0 := qsim.Zero()
	qsim.RZ(math.Pi, q0)

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [0][  0]( 0.0000-1.0000i): 1.0000
}

func Example_qft() {
	qsim := q.New()

	q0 := qsim.Zero()
	q1 := qsim.Zero()
	q2 := qsim.One()

	qsim.H(q0)
	qsim.CR(math.Pi/2, q0, q1)
	qsim.CR(math.Pi/4, q0, q2)

	qsim.H(q1)
	qsim.CR(math.Pi/2, q1, q2)

	qsim.H(q2)

	qsim.CNOT(q0, q2)
	qsim.CNOT(q2, q0)
	qsim.CNOT(q0, q2)

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [000][  0]( 0.3536 0.0000i): 0.1250
	// [001][  1]( 0.2500 0.2500i): 0.1250
	// [010][  2]( 0.0000 0.3536i): 0.1250
	// [011][  3](-0.2500 0.2500i): 0.1250
	// [100][  4](-0.3536 0.0000i): 0.1250
	// [101][  5](-0.2500-0.2500i): 0.1250
	// [110][  6]( 0.0000-0.3536i): 0.1250
	// [111][  7]( 0.2500-0.2500i): 0.1250
}

func ExampleQ_QFT() {
	qsim := q.New()

	q0 := qsim.Zero()
	q1 := qsim.Zero()
	q2 := qsim.One()

	qsim.QFT(q0, q1, q2)
	qsim.Swap(q0, q2)

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [000][  0]( 0.3536 0.0000i): 0.1250
	// [001][  1]( 0.2500 0.2500i): 0.1250
	// [010][  2]( 0.0000 0.3536i): 0.1250
	// [011][  3](-0.2500 0.2500i): 0.1250
	// [100][  4](-0.3536 0.0000i): 0.1250
	// [101][  5](-0.2500-0.2500i): 0.1250
	// [110][  6]( 0.0000-0.3536i): 0.1250
	// [111][  7]( 0.2500-0.2500i): 0.1250
}

func ExampleQ_QFT_ket4() {
	qsim := q.New()

	q0 := qsim.One()
	q1 := qsim.Zero()
	q2 := qsim.Zero()

	qsim.QFT(q0, q1, q2)
	qsim.Swap(q0, q2)

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [000][  0]( 0.3536 0.0000i): 0.1250
	// [001][  1](-0.3536 0.0000i): 0.1250
	// [010][  2]( 0.3536 0.0000i): 0.1250
	// [011][  3](-0.3536 0.0000i): 0.1250
	// [100][  4]( 0.3536 0.0000i): 0.1250
	// [101][  5](-0.3536 0.0000i): 0.1250
	// [110][  6]( 0.3536 0.0000i): 0.1250
	// [111][  7](-0.3536 0.0000i): 0.1250
}

func ExampleQ_InvQFT() {
	qsim := q.New()

	q0 := qsim.Zero()
	q1 := qsim.One()
	q2 := qsim.Zero()

	qsim.QFT(q0, q1, q2)
	qsim.InvQFT(q0, q1, q2)
	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [010][  2]( 1.0000 0.0000i): 1.0000
}

func ExampleQ_Clone() {
	qsim := q.New()

	clone := qsim.Clone()
	clone.Zero()
	clone.Zero()

	q0 := qsim.Zero()
	q1 := qsim.Zero()
	qsim.X(q0, q1)

	fmt.Println(clone)
	fmt.Println(qsim)
	fmt.Println(qsim.Clone())

	// Output:
	// [(1+0i) (0+0i) (0+0i) (0+0i)]
	// [(0+0i) (0+0i) (0+0i) (1+0i)]
	// [(0+0i) (0+0i) (0+0i) (1+0i)]
}

func ExampleQ_String() {
	qsim := q.New()

	q0 := qsim.Zero()
	q1 := qsim.One()
	qsim.X(q0, q1)

	fmt.Println(qsim)

	// Output:
	// [(0+0i) (0+0i) (1+0i) (0+0i)]
}

func ExampleQ_C() {
	qsim := q.New()

	q0 := qsim.Zero()
	q1 := qsim.Zero()

	qsim.H(q0)
	qsim.C(gate.X(), q0, q1) // qsim.CNOT(q0, q1)

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [00][  0]( 0.7071 0.0000i): 0.5000
	// [11][  3]( 0.7071 0.0000i): 0.5000
}

func ExampleQ_CondX() {
	qsim := q.New()

	q0 := qsim.Zero()

	qsim.CondX(false, q0)
	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	qsim.CondX(true, q0)
	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [0][  0]( 1.0000 0.0000i): 1.0000
	// [1][  1]( 1.0000 0.0000i): 1.0000
}

func ExampleQ_CondZ() {
	qsim := q.New()
	q0 := qsim.One()

	qsim.CondZ(false, q0)
	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	qsim.CondZ(true, q0)
	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [1][  1]( 1.0000 0.0000i): 1.0000
	// [1][  1](-1.0000 0.0000i): 1.0000
}

func Example_bellState() {
	qsim := q.New()

	q0 := qsim.Zero()
	q1 := qsim.Zero()

	qsim.H(q0).CNOT(q0, q1)

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	m0 := qsim.Measure(q0)
	m1 := qsim.Measure(q1)
	fmt.Println(m0.Equals(m1))

	// Output:
	// [00][  0]( 0.7071 0.0000i): 0.5000
	// [11][  3]( 0.7071 0.0000i): 0.5000
	// true
}

func Example_bellState2() {
	qsim := q.New()
	qsim.Rand = rand.Const()

	r := qsim.Zeros(2)

	qsim.H(r[0])
	qsim.CNOT(r[0], r[1])
	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [00][  0]( 0.7071 0.0000i): 0.5000
	// [11][  3]( 0.7071 0.0000i): 0.5000
}

func Example_quantumTeleportation() {
	qsim := q.New()

	phi := qsim.New(1, 2)
	q0 := qsim.Zero()
	q1 := qsim.Zero()

	fmt.Println("phi:")
	for _, s := range qsim.State(phi) {
		fmt.Println(s)
	}

	qsim.H(q0).CNOT(q0, q1)
	qsim.CNOT(phi, q0).H(phi)

	qsim.CNOT(q0, q1)
	qsim.CZ(phi, q1)

	qsim.Measure(phi, q0)

	fmt.Println("q1:")
	for _, s := range qsim.State(q1) {
		fmt.Println(s)
	}

	// Output:
	// phi:
	// [0][  0]( 0.4472 0.0000i): 0.2000
	// [1][  1]( 0.8944 0.0000i): 0.8000
	// q1:
	// [0][  0]( 0.4472 0.0000i): 0.2000
	// [1][  1]( 0.8944 0.0000i): 0.8000
}

func Example_quantumTeleportation2() {
	qsim := q.New()

	phi := qsim.New(1, 2)
	q0 := qsim.Zero()
	q1 := qsim.Zero()

	fmt.Println("phi:")
	for _, s := range qsim.State(phi) {
		fmt.Println(s)
	}

	qsim.H(q0).CNOT(q0, q1)
	qsim.CNOT(phi, q0).H(phi)

	mz := qsim.Measure(phi)
	mx := qsim.Measure(q0)

	qsim.CondX(mx.IsOne(), q1)
	qsim.CondZ(mz.IsOne(), q1)

	fmt.Println("q1:")
	for _, s := range qsim.State(q1) {
		fmt.Println(s)
	}

	// Output:
	// phi:
	// [0][  0]( 0.4472 0.0000i): 0.2000
	// [1][  1]( 0.8944 0.0000i): 0.8000
	// q1:
	// [0][  0]( 0.4472 0.0000i): 0.2000
	// [1][  1]( 0.8944 0.0000i): 0.8000
}

func Example_superDenseCoding() {
	sdc := func(g *matrix.Matrix) string {
		qsim := q.New()

		// initial state
		q0 := qsim.Zero()
		q1 := qsim.Zero()

		qsim.H(q0)
		qsim.CNOT(q0, q1)

		// encode
		qsim.Apply(g, q0)

		// decode
		qsim.CNOT(q0, q1)
		qsim.H(q0)

		// measure
		return qsim.M(q0, q1).BinaryString()
	}

	fmt.Printf("I : %v\n", sdc(gate.I()))
	fmt.Printf("X : %v\n", sdc(gate.X()))
	fmt.Printf("Z : %v\n", sdc(gate.Z()))
	fmt.Printf("ZX: %v\n", sdc(gate.Z().Apply(gate.X())))

	// Output:
	// I : 00
	// X : 01
	// Z : 10
	// ZX: 11
}

func Example_errorCorrection() {
	qsim := q.New()

	q0 := qsim.New(1, 2)

	fmt.Println("q0:")
	for _, s := range qsim.State(q0) {
		fmt.Println(s)
	}

	// encoding
	q1 := qsim.Zero()
	q2 := qsim.Zero()
	qsim.CNOT(q0, q1).CNOT(q0, q2)

	// error: first qubit is flipped
	qsim.X(q0)

	fmt.Println("q0(flipped):")
	for _, s := range qsim.State(q0) {
		fmt.Println(s)
	}

	// add ancilla qubit
	q3 := qsim.Zero()
	q4 := qsim.Zero()

	// error correction
	qsim.CNOT(q0, q3).CNOT(q1, q3)
	qsim.CNOT(q1, q4).CNOT(q2, q4)

	m3 := qsim.Measure(q3)
	m4 := qsim.Measure(q4)

	qsim.CondX(m3.IsOne() && m4.IsZero(), q0)
	qsim.CondX(m3.IsOne() && m4.IsOne(), q1)
	qsim.CondX(m3.IsZero() && m4.IsOne(), q2)

	// decoding
	qsim.CNOT(q0, q2).CNOT(q0, q1)

	fmt.Println("q0(corrected):")
	for _, s := range qsim.State(q0) {
		fmt.Println(s)
	}

	// Output:
	// q0:
	// [0][  0]( 0.4472 0.0000i): 0.2000
	// [1][  1]( 0.8944 0.0000i): 0.8000
	// q0(flipped):
	// [0][  0]( 0.8944 0.0000i): 0.8000
	// [1][  1]( 0.4472 0.0000i): 0.2000
	// q0(corrected):
	// [0][  0]( 0.4472 0.0000i): 0.2000
	// [1][  1]( 0.8944 0.0000i): 0.8000
}

func Example_deutschJozsa() {
	type FuncType int
	const (
		Constant FuncType = iota
		Balanced
	)

	oracle := func(qsim *q.Q, q0, q1 q.Qubit) FuncType {
		if rand.Float64() > 0.5 {
			return Constant
		}

		qsim.CNOT(q0, q1)
		return Balanced
	}

	qsim := q.New()
	q0 := qsim.Zero()
	q1 := qsim.One()

	qsim.H(q0, q1)
	ans := oracle(qsim, q0, q1)
	qsim.H(q0)
	m0 := qsim.M(q0)

	if m0.IsZero() && ans == Constant {
		fmt.Println("Correct!")
	}

	if m0.IsOne() && ans == Balanced {
		fmt.Println("Correct!")
	}

	// Output:
	// Correct!
}

func Example_grover3qubit() {
	// NOTE: C. Figgatt, D. Maslov, K. A. Landsman, N. M. Linke, S. Debnath, and C. Monroe. Complete 3-Qubit Grover Search on a Programmable Quantum Computer.
	qsim := q.New()

	// initial state
	r := qsim.Zeros(3)
	a := qsim.One()

	// superposition
	qsim.H(r...).H(a)

	// oracle for |011>|1>
	qsim.X(r[0]).CCCNOT(r[0], r[1], r[2], a).X(r[0])

	// amplification
	qsim.H(r...).H(a)
	qsim.X(r...).CCZ(r[0], r[1], r[2]).X(r...)
	qsim.H(r...)

	for _, s := range qsim.State(r, a) {
		fmt.Println(s)
	}

	// Output:
	// [000 1][  0   1](-0.1768 0.0000i): 0.0313
	// [001 1][  1   1](-0.1768 0.0000i): 0.0313
	// [010 1][  2   1](-0.1768 0.0000i): 0.0313
	// [011 1][  3   1](-0.8839 0.0000i): 0.7813
	// [100 1][  4   1](-0.1768 0.0000i): 0.0313
	// [101 1][  5   1](-0.1768 0.0000i): 0.0313
	// [110 1][  6   1](-0.1768 0.0000i): 0.0313
	// [111 1][  7   1](-0.1768 0.0000i): 0.0313
}

func Example_grover4qubit() {
	// NOTE: Eric R. Johnson, Nic Harrigan, and Merecedes Gimeno-Segovia. Programming Quantum Computers. O'Reilly.
	qsim := q.New()

	// initial state
	q0 := qsim.Zero()
	q1 := qsim.Zero()
	q2 := qsim.Zero()
	q3 := qsim.Zero()

	// superposition
	qsim.H(q0, q1, q2, q3)

	// iteration
	N := number.Pow(2, qsim.NumQubits())
	r := math.Floor(math.Pi / 4 * math.Sqrt(float64(N)))
	for range int(r) {
		// oracle for |110>|x>
		qsim.X(q2, q3)
		qsim.H(q3).CCCNOT(q0, q1, q2, q3).H(q3)
		qsim.X(q2, q3)

		// amplification
		qsim.H(q0, q1, q2, q3)
		qsim.X(q0, q1, q2, q3)
		qsim.H(q3).CCCNOT(q0, q1, q2, q3).H(q3)
		qsim.X(q0, q1, q2, q3)
		qsim.H(q0, q1, q2, q3)
	}

	for _, s := range qsim.State([]q.Qubit{q0, q1, q2}, q3) {
		fmt.Println(s)
	}

	// Output:
	// [000 0][  0   0]( 0.0508 0.0000i): 0.0026
	// [000 1][  0   1]( 0.0508 0.0000i): 0.0026
	// [001 0][  1   0]( 0.0508 0.0000i): 0.0026
	// [001 1][  1   1]( 0.0508 0.0000i): 0.0026
	// [010 0][  2   0]( 0.0508 0.0000i): 0.0026
	// [010 1][  2   1]( 0.0508 0.0000i): 0.0026
	// [011 0][  3   0]( 0.0508 0.0000i): 0.0026
	// [011 1][  3   1]( 0.0508 0.0000i): 0.0026
	// [100 0][  4   0]( 0.0508 0.0000i): 0.0026
	// [100 1][  4   1]( 0.0508 0.0000i): 0.0026
	// [101 0][  5   0]( 0.0508 0.0000i): 0.0026
	// [101 1][  5   1]( 0.0508 0.0000i): 0.0026
	// [110 0][  6   0](-0.9805 0.0000i): 0.9613
	// [110 1][  6   1]( 0.0508 0.0000i): 0.0026
	// [111 0][  7   0]( 0.0508 0.0000i): 0.0026
	// [111 1][  7   1]( 0.0508 0.0000i): 0.0026
}

func Example_qFT() {
	qsim := q.New()

	q0 := qsim.Zero()
	q1 := qsim.One()
	q2 := qsim.Zero()

	qsim.H(q0)
	qsim.CR(q.Theta(2), q1, q0)
	qsim.CR(q.Theta(3), q2, q0)
	qsim.H(q1)
	qsim.CR(q.Theta(2), q2, q1)
	qsim.H(q2)
	qsim.Swap(q0, q2)

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [000][  0]( 0.3536 0.0000i): 0.1250
	// [001][  1]( 0.0000 0.3536i): 0.1250
	// [010][  2](-0.3536 0.0000i): 0.1250
	// [011][  3]( 0.0000-0.3536i): 0.1250
	// [100][  4]( 0.3536 0.0000i): 0.1250
	// [101][  5]( 0.0000 0.3536i): 0.1250
	// [110][  6](-0.3536 0.0000i): 0.1250
	// [111][  7]( 0.0000-0.3536i): 0.1250
}

func Example_shorFactoring15() {
	// NOTE: Zhengjun Cao, Zhenfu Cao, Lihua Liu. Remarks on Quantum Modular Exponentiation and Some Experimental Demonstrations of Shor’s Algorithm.
	N := 15
	a := 7

	qsim := q.New()
	qsim.Rand = rand.Const()

	// initial state
	q0 := qsim.Zero()
	q1 := qsim.Zero()
	q2 := qsim.Zero()

	q3 := qsim.Zero()
	q4 := qsim.Zero()
	q5 := qsim.Zero()
	q6 := qsim.One()

	// superposition
	qsim.H(q0, q1, q2)

	// Controlled-U^(2^0)
	qsim.CNOT(q2, q4)
	qsim.CNOT(q2, q5)

	// Controlled-U^(2^1)
	qsim.CNOT(q3, q5).CCNOT(q1, q5, q3).CNOT(q3, q5)
	qsim.CNOT(q6, q4).CCNOT(q1, q4, q6).CNOT(q6, q4)

	// inverse QFT
	qsim.Swap(q0, q2)
	qsim.InvQFT(q0, q1, q2)

	// measure q0, q1, q2
	m := qsim.Measure(q0, q1, q2).BinaryString()

	// find s/r. 0.010 -> 0.25 -> 1/4, 0.110 -> 0.75 -> 3/4, ...
	s, r, d, ok := number.FindOrder(a, N, fmt.Sprintf("0.%s", m))
	if !ok || number.IsOdd(r) {
		return
	}

	// gcd(a^(r/2)-1, N), gcd(a^(r/2)+1, N)
	p0 := number.GCD(number.Pow(a, r/2)-1, N)
	p1 := number.GCD(number.Pow(a, r/2)+1, N)

	// check non-trivial factor
	if number.IsTrivial(N, p0, p1) {
		return
	}

	fmt.Printf("N=%d, a=%d. p=%v, q=%v. s/r=%d/%d ([0.%v]~%.3f)\n", N, a, p0, p1, s, r, m, d)

	// Output:
	// N=15, a=7. p=3, q=5. s/r=3/4 ([0.110]~0.750)
}

func Example_shorFactoring51() {
	// NOTE: Michael R. Geller, Zhongyuan Zhou. Factoring 51 and 85 with 8 qubits.
	N := 51
	a := 5 // 5, 7, 10, 11, 14, 20, 22, 23, 28, 29, 31, 37, 40, 41, 44, 46
	rng := rand.Const()

	for {
		qsim := q.New()
		qsim.Rand = rng

		q0 := qsim.Zero()
		q1 := qsim.Zero()
		q2 := qsim.Zero()
		q3 := qsim.Zero()

		q4 := qsim.Zero()
		q5 := qsim.Zero()
		q6 := qsim.Zero()
		q7 := qsim.Zero()

		qsim.H(q0, q1, q2, q3)

		qsim.CNOT(q0, q4)
		qsim.CNOT(q1, q5)
		qsim.CNOT(q2, q6)
		qsim.CNOT(q3, q7)

		// inverse QFT
		qsim.Swap(q0, q1, q2, q3)
		qsim.H(q3)
		qsim.CR(-1*q.Theta(2), q3, q2)
		qsim.H(q2)
		qsim.CR(-1*q.Theta(3), q3, q1)
		qsim.CR(-1*q.Theta(2), q2, q1)
		qsim.H(q1)
		qsim.CR(-1*q.Theta(4), q3, q0)
		qsim.CR(-1*q.Theta(3), q2, q0)
		qsim.CR(-1*q.Theta(2), q1, q0)
		qsim.H(q0)

		m := qsim.Measure(q0, q1, q2, q3).BinaryString()
		s, r, d, ok := number.FindOrder(a, N, fmt.Sprintf("0.%s", m))
		if !ok || number.IsOdd(r) {
			continue
		}

		p0 := number.GCD(number.Pow(a, r/2)-1, N)
		p1 := number.GCD(number.Pow(a, r/2)+1, N)
		if number.IsTrivial(N, p0, p1) {
			continue
		}

		fmt.Printf("N=%d, a=%d. p=%v, q=%v. s/r=%d/%d ([0.%v]~%.3f)\n", N, a, p0, p1, s, r, m, d)
		break
	}

	// Output:
	// N=51, a=5. p=3, q=17. s/r=7/16 ([0.0111]~0.438)
}

func Example_shorFactoring85() {
	// NOTE: Michael R. Geller, Zhongyuan Zhou. Factoring 51 and 85 with 8 qubits.
	N := 85
	a := 14
	rng := rand.Const()

	for {
		qsim := q.New()
		qsim.Rand = rng

		q0 := qsim.Zero()
		q1 := qsim.Zero()
		q2 := qsim.Zero()
		q3 := qsim.Zero()

		q4 := qsim.Zero()
		q5 := qsim.Zero()
		q6 := qsim.Zero()
		q7 := qsim.Zero()

		qsim.H(q0, q1, q2, q3)

		qsim.CNOT(q0, q4)
		qsim.CNOT(q1, q5)
		qsim.CNOT(q2, q6)
		qsim.CNOT(q3, q7)

		// inverse QFT
		qsim.Swap(q0, q1, q2, q3)
		qsim.H(q3)
		qsim.CR(-1*q.Theta(2), q3, q2).H(q2)
		qsim.CR(-1*q.Theta(3), q3, q1).CR(-1*q.Theta(2), q2, q1).H(q1)
		qsim.CR(-1*q.Theta(4), q3, q0).CR(-1*q.Theta(3), q2, q0).CR(-1*q.Theta(2), q1, q0).H(q0)

		m := qsim.Measure(q0, q1, q2, q3).BinaryString()
		s, r, d, ok := number.FindOrder(a, N, fmt.Sprintf("0.%s", m))
		if !ok || number.IsOdd(r) {
			continue
		}

		p0 := number.GCD(number.Pow(a, r/2)-1, N)
		p1 := number.GCD(number.Pow(a, r/2)+1, N)
		if number.IsTrivial(N, p0, p1) {
			continue
		}

		fmt.Printf("N=%d, a=%d. p=%v, q=%v. s/r=%d/%d ([0.%v]~%.3f)\n", N, a, p0, p1, s, r, m, d)
		break
	}

	// Output:
	// N=85, a=14. p=5, q=17. s/r=7/16 ([0.0111]~0.438)
}

func Example_any() {
	h := gate.U(math.Pi/2, 0, math.Pi)
	x := gate.U(math.Pi, 0, math.Pi)

	qsim := q.New()
	q0 := qsim.Zero()
	q1 := qsim.Zero()

	qsim.Apply(h, q0)
	qsim.C(x, q0, q1)

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [00][  0]( 0.7071 0.0000i): 0.5000
	// [11][  3]( 0.7071 0.0000i): 0.5000
}

func ExampleQ_Underlying() {
	qsim := q.New()
	qsim.Zero()

	for _, s := range qsim.Underlying().State() {
		fmt.Println(s)
	}

	// Output:
	// [0][  0]( 1.0000 0.0000i): 1.0000
}
