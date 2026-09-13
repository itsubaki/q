package q_test

import (
	"fmt"
	"math"

	"github.com/itsubaki/q"
	"github.com/itsubaki/q/math/matrix"
	"github.com/itsubaki/q/math/number"
	"github.com/itsubaki/q/math/rand"
	"github.com/itsubaki/q/quantum/gate"
	"github.com/itsubaki/q/quantum/qubit"
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
	// [00] ( 0.5000 0.0000i): 0.2500
	// [01] ( 0.5000 0.0000i): 0.2500
	// [10] ( 0.5000 0.0000i): 0.2500
	// [11] ( 0.5000 0.0000i): 0.2500
}

func ExampleQ_Zeros() {
	qsim := q.New()
	qb := qsim.Zeros(2)

	qsim.H(qb...)
	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [00] ( 0.5000 0.0000i): 0.2500
	// [01] ( 0.5000 0.0000i): 0.2500
	// [10] ( 0.5000 0.0000i): 0.2500
	// [11] ( 0.5000 0.0000i): 0.2500
}

func ExampleQ_Ones() {
	qsim := q.New()
	qb := qsim.Ones(2)

	qsim.H(qb...)
	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [00] ( 0.5000 0.0000i): 0.2500
	// [01] (-0.5000 0.0000i): 0.2500
	// [10] (-0.5000 0.0000i): 0.2500
	// [11] ( 0.5000 0.0000i): 0.2500
}

func ExampleQ_Pluses() {
	qsim := q.New()
	qb := qsim.Pluses(2)

	for _, s := range qsim.State(qb) {
		fmt.Println(s)
	}

	// Output:
	// [00] ( 0.5000 0.0000i): 0.2500
	// [01] ( 0.5000 0.0000i): 0.2500
	// [10] ( 0.5000 0.0000i): 0.2500
	// [11] ( 0.5000 0.0000i): 0.2500
}

func ExampleQ_Minuses() {
	qsim := q.New()
	qb := qsim.Minuses(2)

	for _, s := range qsim.State(qb) {
		fmt.Println(s)
	}

	// Output:
	// [00] ( 0.5000 0.0000i): 0.2500
	// [01] (-0.5000 0.0000i): 0.2500
	// [10] (-0.5000 0.0000i): 0.2500
	// [11] ( 0.5000 0.0000i): 0.2500
}

func ExampleQ_From() {
	qsim := q.New()
	qb := qsim.From("01+-")

	for _, s := range qsim.State(qb) {
		fmt.Println(s)
	}

	// Output:
	// [0100] ( 0.5000 0.0000i): 0.2500
	// [0101] (-0.5000 0.0000i): 0.2500
	// [0110] ( 0.5000 0.0000i): 0.2500
	// [0111] (-0.5000 0.0000i): 0.2500
}

func ExampleQ_Reset() {
	qsim := q.New()
	qb := qsim.Zeros(2)

	qsim.X(qb[0])
	qsim.Reset(qb...)

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [00] ( 1.0000 0.0000i): 1.0000
}

func ExampleQ_NumQubits() {
	qsim := q.New()
	fmt.Println(qsim.NumQubits())

	qsim.Zeros(3)
	fmt.Println(qsim.NumQubits())

	qsim.Zeros(10)
	fmt.Println(qsim.NumQubits())

	// Output:
	// 0
	// 3
	// 13
}

func ExampleQ_Amplitude() {
	qsim := q.New()
	q0 := qsim.Zero()
	q1 := qsim.Zero()

	qsim.H(q0)
	qsim.CNOT(q0, q1)

	for _, a := range qsim.Amplitude() {
		fmt.Printf("%.4f\n", a)
	}

	// Output:
	// (0.7071+0.0000i)
	// (0.0000+0.0000i)
	// (0.0000+0.0000i)
	// (0.7071+0.0000i)
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
	qsim.SetRand(rand.Const())

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
	qsim.SetRand(rand.Const())

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
	h := gate.H()
	cnot := gate.CNOT(n, q0.Index(), q1.Index())

	qsim.G(h, q0)
	qsim.Apply(cnot)

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [00] ( 0.7071 0.0000i): 0.5000
	// [11] ( 0.7071 0.0000i): 0.5000
}

func ExampleQ_U() {
	qsim := q.New()
	qb := qsim.Zero()

	qsim.U(math.Pi, 0, math.Pi, qb)
	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [1] ( 1.0000 0.0000i): 1.0000
}

func ExampleQ_I() {
	qsim := q.New()
	qb := qsim.Zero()

	qsim.I(qb)
	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [0] ( 1.0000 0.0000i): 1.0000
}

func ExampleQ_X() {
	qsim := q.New()
	qb := qsim.Zero()

	qsim.X(qb)
	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [1] ( 1.0000 0.0000i): 1.0000
}

func ExampleQ_Y() {
	qsim := q.New()
	qb := qsim.Zero()

	qsim.H(qb)
	qsim.Y(qb)

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [0] ( 0.0000-0.7071i): 0.5000
	// [1] ( 0.0000 0.7071i): 0.5000
}

func ExampleQ_Z() {
	qsim := q.New()
	qb := qsim.Zero()

	qsim.H(qb)
	qsim.Z(qb)

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [0] ( 0.7071 0.0000i): 0.5000
	// [1] (-0.7071 0.0000i): 0.5000
}

func ExampleQ_S() {
	qsim := q.New()
	qb := qsim.Zero()

	qsim.H(qb)
	qsim.S(qb)

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [0] ( 0.7071 0.0000i): 0.5000
	// [1] ( 0.0000 0.7071i): 0.5000
}

func ExampleQ_T() {
	qsim := q.New()
	qb := qsim.Zero()

	qsim.H(qb)
	qsim.T(qb)

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [0] ( 0.7071 0.0000i): 0.5000
	// [1] ( 0.5000 0.5000i): 0.5000
}

func ExampleQ_R() {
	qsim := q.New()
	qb := qsim.Zero()

	qsim.H(qb)
	qsim.R(2*math.Pi/4, qb)

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [0] ( 0.7071 0.0000i): 0.5000
	// [1] ( 0.0000 0.7071i): 0.5000
}

func ExampleQ_RX() {
	qsim := q.New()
	qb := qsim.Zero()

	qsim.H(qb)
	qsim.RX(math.Pi, qb)

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [0] ( 0.0000-0.7071i): 0.5000
	// [1] ( 0.0000-0.7071i): 0.5000
}

func ExampleQ_RY() {
	qsim := q.New()
	qb := qsim.Zero()

	qsim.H(qb)
	qsim.RY(math.Pi, qb)

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [0] (-0.7071 0.0000i): 0.5000
	// [1] ( 0.7071 0.0000i): 0.5000
}

func ExampleQ_RZ() {
	qsim := q.New()
	qb := qsim.Zero()

	qsim.H(qb)
	qsim.RZ(math.Pi, qb)

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [0] ( 0.0000-0.7071i): 0.5000
	// [1] ( 0.0000 0.7071i): 0.5000
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
	// [00] ( 0.7071 0.0000i): 0.5000
	// [11] ( 0.7071 0.0000i): 0.5000
}

func ExampleQ_CU() {
	qsim := q.New()
	q0 := qsim.Zero()
	q1 := qsim.Zero()

	qsim.H(q0)
	qsim.CU(math.Pi, 0, math.Pi, q0, q1)

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [00] ( 0.7071 0.0000i): 0.5000
	// [11] ( 0.7071 0.0000i): 0.5000
}

func ExampleQ_CX() {
	qsim := q.New()
	q0 := qsim.Zero()
	q1 := qsim.Zero()

	qsim.H(q0)
	qsim.CX(q0, q1)

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [00] ( 0.7071 0.0000i): 0.5000
	// [11] ( 0.7071 0.0000i): 0.5000
}

func ExampleQ_CZ() {
	qsim := q.New()
	q0 := qsim.Zero()
	q1 := qsim.Zero()

	qsim.H(q0)
	qsim.H(q1)
	qsim.CZ(q0, q1)

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [00] ( 0.5000 0.0000i): 0.2500
	// [01] ( 0.5000 0.0000i): 0.2500
	// [10] ( 0.5000 0.0000i): 0.2500
	// [11] (-0.5000 0.0000i): 0.2500
}

func ExampleQ_ControlledH() {
	qsim := q.New()
	q0 := qsim.Zero()
	q1 := qsim.Zero()

	qsim.H(q0)
	qsim.ControlledH([]q.Qubit{q0}, []q.Qubit{q1})

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [00] ( 0.7071 0.0000i): 0.5000
	// [10] ( 0.5000 0.0000i): 0.2500
	// [11] ( 0.5000 0.0000i): 0.2500
}

func ExampleQ_ControlledX() {
	qsim := q.New()
	q0 := qsim.Zero()
	q1 := qsim.Zero()

	qsim.H(q0)
	qsim.ControlledX([]q.Qubit{q0}, []q.Qubit{q1})

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [00] ( 0.7071 0.0000i): 0.5000
	// [11] ( 0.7071 0.0000i): 0.5000
}

func ExampleQ_CCZ() {
	qsim := q.New()
	q0 := qsim.Zero()
	q1 := qsim.Zero()
	q2 := qsim.One()

	qsim.H(q0)
	qsim.H(q1)
	qsim.CCZ(q0, q1, q2)

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [001] ( 0.5000 0.0000i): 0.2500
	// [011] ( 0.5000 0.0000i): 0.2500
	// [101] ( 0.5000 0.0000i): 0.2500
	// [111] (-0.5000 0.0000i): 0.2500
}

func ExampleQ_CCNOT() {
	qsim := q.New()
	q0 := qsim.Zero()
	q1 := qsim.Zero()
	q2 := qsim.Zero()

	qsim.H(q0)
	qsim.H(q1)
	qsim.CCNOT(q0, q1, q2)

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [000] ( 0.5000 0.0000i): 0.2500
	// [010] ( 0.5000 0.0000i): 0.2500
	// [100] ( 0.5000 0.0000i): 0.2500
	// [111] ( 0.5000 0.0000i): 0.2500
}

func ExampleQ_CondX() {
	qsim := q.New()
	qb := qsim.Zero()

	for _, b := range []bool{false, true} {
		qsim.CondX(b, qb)

		for _, s := range qsim.State() {
			fmt.Println(s)
		}
	}

	// Output:
	// [0] ( 1.0000 0.0000i): 1.0000
	// [1] ( 1.0000 0.0000i): 1.0000
}

func ExampleQ_CondZ() {
	qsim := q.New()
	qb := qsim.One()

	for _, b := range []bool{false, true} {
		qsim.CondZ(b, qb)

		for _, s := range qsim.State() {
			fmt.Println(s)
		}
	}

	// Output:
	// [1] ( 1.0000 0.0000i): 1.0000
	// [1] (-1.0000 0.0000i): 1.0000
}

func ExampleQ_Cond() {
	qsim := q.New()
	qb := qsim.Zero()

	for _, b := range []bool{false, true} {
		qsim.Cond(b, gate.X(), qb)

		for _, s := range qsim.State() {
			fmt.Println(s)
		}
	}

	// Output:
	// [0] ( 1.0000 0.0000i): 1.0000
	// [1] ( 1.0000 0.0000i): 1.0000
}

func ExampleQ_Swap() {
	qsim := q.New()
	q0 := qsim.Zero()
	q1 := qsim.One()

	qsim.Swap(q0, q1)
	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [10] ( 1.0000 0.0000i): 1.0000
}

func ExampleQ_Qubit() {
	qsim := q.New()
	qsim.Zero()

	for _, s := range qsim.Qubit().State() {
		fmt.Println(s)
	}

	// Output:
	// [0] ( 1.0000 0.0000i): 1.0000
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

func ExampleTop() {
	qsim := q.New()
	qb := qsim.Zeros(3)

	qsim.H(qb...)
	for _, s := range q.Top(qsim.State(), 2) {
		fmt.Println(s)
	}

	// Output:
	// [000] ( 0.3536 0.0000i): 0.1250
	// [001] ( 0.3536 0.0000i): 0.1250
}

func ExampleTop_all() {
	qsim := q.New()
	qb := qsim.Zeros(3)

	qsim.H(qb...)
	for _, s := range q.Top(qsim.State(), -1) {
		fmt.Println(s)
	}

	// Output:
	// [000] ( 0.3536 0.0000i): 0.1250
	// [001] ( 0.3536 0.0000i): 0.1250
	// [010] ( 0.3536 0.0000i): 0.1250
	// [011] ( 0.3536 0.0000i): 0.1250
	// [100] ( 0.3536 0.0000i): 0.1250
	// [101] ( 0.3536 0.0000i): 0.1250
	// [110] ( 0.3536 0.0000i): 0.1250
	// [111] ( 0.3536 0.0000i): 0.1250
}

func Example_bell() {
	qsim := q.New()
	q0 := qsim.Zero()
	q1 := qsim.Zero()

	qsim.H(q0)
	qsim.CNOT(q0, q1)

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	m0 := qsim.Measure(q0)
	m1 := qsim.Measure(q1)
	fmt.Println(m0.Equal(m1))

	// Output:
	// [00] ( 0.7071 0.0000i): 0.5000
	// [11] ( 0.7071 0.0000i): 0.5000
	// true
}

func Example_quantumTeleportation() {
	qsim := q.New()
	psi := qsim.New(1, 2)
	q0 := qsim.Zero()
	q1 := qsim.Zero()

	for _, s := range qsim.State(psi) {
		fmt.Println(s)
	}

	qsim.H(q0)
	qsim.CNOT(q0, q1)
	qsim.CNOT(psi, q0)
	qsim.H(psi)

	mz := qsim.Measure(psi)
	mx := qsim.Measure(q0)

	qsim.CondX(mx.IsOne(), q1)
	qsim.CondZ(mz.IsOne(), q1)

	for _, s := range qsim.State(q1) {
		fmt.Println(s)
	}

	// Output:
	// [0] ( 0.4472 0.0000i): 0.2000
	// [1] ( 0.8944 0.0000i): 0.8000
	// [0] ( 0.4472 0.0000i): 0.2000
	// [1] ( 0.8944 0.0000i): 0.8000
}

func Example_deutschJozsa() {
	constant := func(qsim *q.Q, q0, q1 q.Qubit) string {
		return "Constant"
	}

	balanced := func(qsim *q.Q, q0, q1 q.Qubit) string {
		qsim.CNOT(q0, q1)
		return "Balanced"
	}

	deutschJozsa := func(oracle func(qsim *q.Q, q0, q1 q.Qubit) string) (string, int) {
		qsim := q.New()
		q0 := qsim.Zero()
		q1 := qsim.One()

		qsim.H(q0, q1)
		ans := oracle(qsim, q0, q1)
		qsim.H(q0)

		if qsim.M(q0).IsZero() {
			return ans, 0
		}

		return ans, 1
	}

	fmt.Println(deutschJozsa(constant))
	fmt.Println(deutschJozsa(balanced))

	// Output:
	// Constant 0
	// Balanced 1
}

func Example_grover() {
	qsim := q.New()
	q0 := qsim.Zero()
	q1 := qsim.Zero()
	q2 := qsim.Zero()
	q3 := qsim.One()

	qsim.H(q0, q1, q2, q3)

	N := number.Pow(2, qsim.NumQubits())
	R := int(math.Pi / 4 * math.Sqrt(float64(N)))
	for range R {
		// oracle for |110>|x>
		qsim.X(q2, q3)
		qsim.H(q3)
		qsim.CCCNOT(q0, q1, q2, q3)
		qsim.H(q3)
		qsim.X(q2, q3)

		// diffuser
		qsim.H(q0, q1, q2)
		qsim.X(q0, q1, q2)
		qsim.CCCNOT(q0, q1, q2, q3)
		qsim.X(q0, q1, q2)
		qsim.H(q0, q1, q2)
	}

	for _, s := range qsim.State([]q.Qubit{q0, q1, q2}, q3) {
		fmt.Println(s)
	}

	// Output:
	// [000 0] ( 0.0508 0.0000i): 0.0026
	// [000 1] (-0.0508 0.0000i): 0.0026
	// [001 0] ( 0.0508 0.0000i): 0.0026
	// [001 1] (-0.0508 0.0000i): 0.0026
	// [010 0] ( 0.0508 0.0000i): 0.0026
	// [010 1] (-0.0508 0.0000i): 0.0026
	// [011 0] ( 0.0508 0.0000i): 0.0026
	// [011 1] (-0.0508 0.0000i): 0.0026
	// [100 0] ( 0.0508 0.0000i): 0.0026
	// [100 1] (-0.0508 0.0000i): 0.0026
	// [101 0] ( 0.0508 0.0000i): 0.0026
	// [101 1] (-0.0508 0.0000i): 0.0026
	// [110 0] (-0.9805 0.0000i): 0.9613
	// [110 1] (-0.0508 0.0000i): 0.0026
	// [111 0] ( 0.0508 0.0000i): 0.0026
	// [111 1] (-0.0508 0.0000i): 0.0026
}

func Example_qft() {
	qsim := q.New()
	q0 := qsim.Zero()
	q1 := qsim.Zero()
	q2 := qsim.One()

	qsim.H(q0)
	qsim.CR(q.Theta(2), q0, q1)
	qsim.CR(q.Theta(3), q0, q2)

	qsim.H(q1)
	qsim.CR(q.Theta(2), q1, q2)

	qsim.H(q2)

	// swap
	qsim.CNOT(q0, q2)
	qsim.CNOT(q2, q0)
	qsim.CNOT(q0, q2)

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

func Example_superDenseCoding() {
	sdc := func(g *matrix.Matrix) string {
		qsim := q.New()
		q0 := qsim.Zero()
		q1 := qsim.Zero()

		qsim.H(q0)
		qsim.CNOT(q0, q1)

		// encode
		qsim.G(g, q0)

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

func Example_ecc() {
	qsim := q.New()
	q0 := qsim.New(1, 2)

	fmt.Println("q0:")
	for _, s := range qsim.State(q0) {
		fmt.Println(s)
	}

	q1 := qsim.Zero()
	q2 := qsim.Zero()
	qsim.CNOT(q0, q1)
	qsim.CNOT(q0, q2)

	fmt.Println("q0(encoded):")
	for _, s := range qsim.State(q0, []q.Qubit{q1, q2}) {
		fmt.Println(s)
	}

	// error: the first qubit is flipped
	qsim.X(q0)

	fmt.Println("q0(flipped):")
	for _, s := range qsim.State(q0, []q.Qubit{q1, q2}) {
		fmt.Println(s)
	}

	q3 := qsim.Zero()
	q4 := qsim.Zero()

	// error correction
	qsim.CNOT(q0, q3)
	qsim.CNOT(q1, q3)
	qsim.CNOT(q1, q4)
	qsim.CNOT(q2, q4)

	m3 := qsim.Measure(q3)
	m4 := qsim.Measure(q4)

	qsim.CondX(m3.IsOne() && m4.IsZero(), q0)
	qsim.CondX(m3.IsOne() && m4.IsOne(), q1)
	qsim.CondX(m3.IsZero() && m4.IsOne(), q2)

	// decode
	qsim.CNOT(q0, q2)
	qsim.CNOT(q0, q1)

	fmt.Println("q0(corrected):")
	for _, s := range qsim.State(q0, []q.Qubit{q1, q2}, []q.Qubit{q3, q4}) {
		fmt.Println(s)
	}

	// Output:
	// q0:
	// [0] ( 0.4472 0.0000i): 0.2000
	// [1] ( 0.8944 0.0000i): 0.8000
	// q0(encoded):
	// [0 00] ( 0.4472 0.0000i): 0.2000
	// [1 11] ( 0.8944 0.0000i): 0.8000
	// q0(flipped):
	// [0 11] ( 0.8944 0.0000i): 0.8000
	// [1 00] ( 0.4472 0.0000i): 0.2000
	// q0(corrected):
	// [0 00 10] ( 0.4472 0.0000i): 0.2000
	// [1 00 10] ( 0.8944 0.0000i): 0.8000
}

func Example_gateTeleportation() {
	qsim := q.New()
	psi := qsim.New(1, 2)
	a := qsim.Zero()

	qsim.H(a)
	qsim.T(a) // magic state

	qsim.CNOT(a, psi)
	m0 := qsim.Measure(psi)
	qsim.Cond(m0.IsOne(), gate.X(), a)
	qsim.Cond(m0.IsOne(), gate.S(), a)

	{
		qs := q.New()
		qb := qs.New(1, 2)
		qs.T(qb)

		fmt.Println(qubit.EqualUpToGlobalPhase(qsim.State(a), qs.State(qb)))
	}

	// Output:
	// true
}

func Example_any() {
	h := gate.U(math.Pi/2, 0, math.Pi)
	x := gate.U(math.Pi, 0, math.Pi)

	qsim := q.New()
	q0 := qsim.Zero()
	q1 := qsim.Zero()

	qsim.G(h, q0)
	qsim.C(x, q0, q1)

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [00] ( 0.7071 0.0000i): 0.5000
	// [11] ( 0.7071 0.0000i): 0.5000
}
