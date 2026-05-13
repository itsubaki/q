package qubit_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/itsubaki/q/math/epsilon"
	"github.com/itsubaki/q/math/matrix"
	"github.com/itsubaki/q/math/number"
	"github.com/itsubaki/q/math/rand"
	"github.com/itsubaki/q/math/vector"
	"github.com/itsubaki/q/quantum/gate"
	"github.com/itsubaki/q/quantum/qubit"
)

func ExampleZero() {
	z := qubit.Zero()
	for _, s := range z.State() {
		fmt.Println(s)
	}

	z2 := qubit.Zeros(2)
	for _, s := range z2.State() {
		fmt.Println(s)
	}

	// Output:
	// [0] ( 1.0000 0.0000i): 1.0000
	// [00] ( 1.0000 0.0000i): 1.0000
}

func ExampleFrom() {
	z := qubit.From("0011")
	for _, s := range z.State() {
		fmt.Println(s)
	}

	// Output:
	// [0011] ( 1.0000 0.0000i): 1.0000
}

func ExampleFrom_plus() {
	z := qubit.From("+-")
	for _, s := range z.State() {
		fmt.Println(s)
	}

	// Output:
	// [00] ( 0.5000 0.0000i): 0.2500
	// [01] (-0.5000 0.0000i): 0.2500
	// [10] ( 0.5000 0.0000i): 0.2500
	// [11] (-0.5000 0.0000i): 0.2500
}

func ExampleMinuses() {
	qb := qubit.Minuses(2)
	for _, s := range qb.State() {
		fmt.Println(s)
	}

	// Output:
	// [00] ( 0.5000 0.0000i): 0.2500
	// [01] (-0.5000 0.0000i): 0.2500
	// [10] (-0.5000 0.0000i): 0.2500
	// [11] ( 0.5000 0.0000i): 0.2500
}

func ExampleQubit_OuterProduct() {
	v := qubit.Zero()
	op := v.OuterProduct(v)

	for _, r := range op.Seq2() {
		fmt.Println(r)
	}

	// Output:
	// [(1+0i) (0+0i)]
	// [(0+0i) (0+0i)]
}

func ExampleQubit_G() {
	qb := qubit.Zeros(2)

	h := gate.H()
	cx := gate.CNOT(2, 0, 1)

	qb.G(h, 0)
	qb.Apply(cx)

	for _, s := range qb.State() {
		fmt.Println(s)
	}

	// Output:
	// [00] ( 0.7071 0.0000i): 0.5000
	// [11] ( 0.7071 0.0000i): 0.5000
}

func ExampleQubit_U() {
	qb := qubit.Zeros(2)
	qb.U(math.Pi/2, 0, 0, 0)
	qb.U(math.Pi/2, 0, 0, 1)

	for _, s := range qb.State() {
		fmt.Println(s)
	}

	// Output:
	// [00] ( 0.5000 0.0000i): 0.2500
	// [01] ( 0.5000 0.0000i): 0.2500
	// [10] ( 0.5000 0.0000i): 0.2500
	// [11] ( 0.5000 0.0000i): 0.2500
}

func ExampleQubit_I() {
	qb := qubit.Zero()
	qb.H(0)
	qb.I(0)

	for _, s := range qb.State() {
		fmt.Println(s)
	}

	// Output:
	// [0] ( 0.7071 0.0000i): 0.5000
	// [1] ( 0.7071 0.0000i): 0.5000
}

func ExampleQubit_H() {
	qb := qubit.Zero()
	qb.H(0)

	for _, s := range qb.State() {
		fmt.Println(s)
	}

	// Output:
	// [0] ( 0.7071 0.0000i): 0.5000
	// [1] ( 0.7071 0.0000i): 0.5000
}

func ExampleQubit_X() {
	qb := qubit.Zero()
	qb.X(0)

	for _, s := range qb.State() {
		fmt.Println(s)
	}

	// Output:
	// [1] ( 1.0000 0.0000i): 1.0000
}

func ExampleQubit_Y() {
	qb := qubit.Zero()
	qb.H(0)
	qb.Y(0)

	for _, s := range qb.State() {
		fmt.Println(s)
	}

	// Output:
	// [0] ( 0.0000-0.7071i): 0.5000
	// [1] ( 0.0000 0.7071i): 0.5000
}

func ExampleQubit_Z() {
	qb := qubit.Zero()
	qb.H(0)
	qb.Z(0)

	for _, s := range qb.State() {
		fmt.Println(s)
	}

	// Output:
	// [0] ( 0.7071 0.0000i): 0.5000
	// [1] (-0.7071 0.0000i): 0.5000
}

func ExampleQubit_R() {
	qb := qubit.Zero()
	qb.H(0)
	qb.R(math.Pi, 0)

	for _, s := range qb.State() {
		fmt.Println(s)
	}

	// Output:
	// [0] ( 0.7071 0.0000i): 0.5000
	// [1] (-0.7071 0.0000i): 0.5000
}

func ExampleQubit_S() {
	qb := qubit.Zero()
	qb.H(0)
	qb.S(0)

	for _, s := range qb.State() {
		fmt.Println(s)
	}

	// Output:
	// [0] ( 0.7071 0.0000i): 0.5000
	// [1] ( 0.0000 0.7071i): 0.5000
}

func ExampleQubit_T() {
	qb := qubit.Zero()
	qb.H(0)
	qb.T(0)

	for _, s := range qb.State() {
		fmt.Println(s)
	}

	// Output:
	// [0] ( 0.7071 0.0000i): 0.5000
	// [1] ( 0.5000 0.5000i): 0.5000
}

func ExampleQubit_RX() {
	qb := qubit.Zero()
	qb.H(0)
	qb.RX(math.Pi, 0)

	for _, s := range qb.State() {
		fmt.Println(s)
	}

	// Output:
	// [0] ( 0.0000-0.7071i): 0.5000
	// [1] ( 0.0000-0.7071i): 0.5000
}

func ExampleQubit_RY() {
	qb := qubit.Zero()
	qb.H(0)
	qb.RY(math.Pi, 0)

	for _, s := range qb.State() {
		fmt.Println(s)
	}

	// Output:
	// [0] (-0.7071 0.0000i): 0.5000
	// [1] ( 0.7071 0.0000i): 0.5000
}

func ExampleQubit_RZ() {
	qb := qubit.Zero()
	qb.H(0)
	qb.RZ(math.Pi, 0)

	for _, s := range qb.State() {
		fmt.Println(s)
	}

	// Output:
	// [0] ( 0.0000-0.7071i): 0.5000
	// [1] ( 0.0000 0.7071i): 0.5000
}

func ExampleQubit_C() {
	qb := qubit.Zeros(2)
	qb.H(0)
	qb.C(gate.X(), 0, 1)

	for _, s := range qb.State() {
		fmt.Println(s)
	}

	// Output:
	// [00] ( 0.7071 0.0000i): 0.5000
	// [11] ( 0.7071 0.0000i): 0.5000
}

func ExampleQubit_CU() {
	qb := qubit.Zeros(2)
	qb.H(0)
	qb.CU(math.Pi, 0, math.Pi, 0, 1)

	for _, s := range qb.State() {
		fmt.Println(s)
	}

	// Output:
	// [00] ( 0.7071 0.0000i): 0.5000
	// [11] ( 0.7071 0.0000i): 0.5000
}

func ExampleQubit_CH() {
	qb := qubit.Zeros(2)
	qb.H(0)
	qb.CH(0, 1)

	for _, s := range qb.State() {
		fmt.Println(s)
	}

	// Output:
	// [00] ( 0.7071 0.0000i): 0.5000
	// [10] ( 0.5000 0.0000i): 0.2500
	// [11] ( 0.5000 0.0000i): 0.2500
}

func ExampleQubit_CX() {
	qb := qubit.Zeros(2)
	qb.H(0)
	qb.CX(0, 1)

	for _, s := range qb.State() {
		fmt.Println(s)
	}

	// Output:
	// [00] ( 0.7071 0.0000i): 0.5000
	// [11] ( 0.7071 0.0000i): 0.5000
}

func ExampleQubit_CZ() {
	qb := qubit.Zeros(2)
	qb.H(0)
	qb.H(1)
	qb.CZ(0, 1)

	for _, s := range qb.State() {
		fmt.Println(s)
	}

	// Output:
	// [00] ( 0.5000 0.0000i): 0.2500
	// [01] ( 0.5000 0.0000i): 0.2500
	// [10] ( 0.5000 0.0000i): 0.2500
	// [11] (-0.5000 0.0000i): 0.2500
}

func ExampleQubit_ControlledH() {
	qb := qubit.Zeros(2)
	qb.H(0)
	qb.ControlledH([]int{0}, 1)

	for _, s := range qb.State() {
		fmt.Println(s)
	}

	// Output:
	// [00] ( 0.7071 0.0000i): 0.5000
	// [10] ( 0.5000 0.0000i): 0.2500
	// [11] ( 0.5000 0.0000i): 0.2500
}

func ExampleQubit_ControlledX() {
	qb := qubit.Zeros(3)
	qb.X(0)
	qb.X(1)
	qb.ControlledX([]int{0, 1}, 2)

	for _, s := range qb.State() {
		fmt.Println(s)
	}

	// Output:
	// [111] ( 1.0000 0.0000i): 1.0000
}

func ExampleQubit_ControlledZ() {
	qb := qubit.Zeros(2)
	qb.H(0)
	qb.H(1)
	qb.ControlledZ([]int{0}, 1)

	for _, s := range qb.State() {
		fmt.Println(s)
	}

	// Output:
	// [00] ( 0.5000 0.0000i): 0.2500
	// [01] ( 0.5000 0.0000i): 0.2500
	// [10] ( 0.5000 0.0000i): 0.2500
	// [11] (-0.5000 0.0000i): 0.2500
}

func ExampleQubit_Swap() {
	qb := qubit.Zeros(2)
	qb.X(0)

	qb.Swap(0, 1)
	for _, s := range qb.State() {
		fmt.Println(s)
	}

	// Output:
	// [01] ( 1.0000 0.0000i): 1.0000
}

func ExampleQubit_Swap_eq() {
	qb := qubit.Zeros(2)
	qb.X(0)

	qb.Swap(0, 0)
	for _, s := range qb.State() {
		fmt.Println(s)
	}

	// Output:
	// [10] ( 1.0000 0.0000i): 1.0000
}

func ExampleQubit_QFT() {
	qb := qubit.Zeros(3)
	qb.X(2)
	qb.QFT()
	qb.Swap(0, 2)

	for _, s := range qb.State() {
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

func ExampleQubit_InvQFT() {
	qb := qubit.Zeros(3)
	qb.X(2)
	qb.QFT()
	qb.InvQFT()

	for _, s := range qb.State() {
		fmt.Println(s)
	}

	// Output:
	// [001] ( 1.0000 0.0000i): 1.0000
}

func ExampleQubit_Set() {
	qb := qubit.Zeros(2)
	qb.Set(vector.New(1, 0, 0, 1))

	for _, s := range qb.State() {
		fmt.Println(s)
	}

	// Output:
	// [00] ( 0.7071 0.0000i): 0.5000
	// [11] ( 0.7071 0.0000i): 0.5000
}

func ExampleQubit_State() {
	v := qubit.Pluses(4)
	for _, s := range v.State([]int{0, 1, 2, 3}) {
		fmt.Println(s)
	}

	// Output:
	// [0000] ( 0.2500 0.0000i): 0.0625
	// [0001] ( 0.2500 0.0000i): 0.0625
	// [0010] ( 0.2500 0.0000i): 0.0625
	// [0011] ( 0.2500 0.0000i): 0.0625
	// [0100] ( 0.2500 0.0000i): 0.0625
	// [0101] ( 0.2500 0.0000i): 0.0625
	// [0110] ( 0.2500 0.0000i): 0.0625
	// [0111] ( 0.2500 0.0000i): 0.0625
	// [1000] ( 0.2500 0.0000i): 0.0625
	// [1001] ( 0.2500 0.0000i): 0.0625
	// [1010] ( 0.2500 0.0000i): 0.0625
	// [1011] ( 0.2500 0.0000i): 0.0625
	// [1100] ( 0.2500 0.0000i): 0.0625
	// [1101] ( 0.2500 0.0000i): 0.0625
	// [1110] ( 0.2500 0.0000i): 0.0625
	// [1111] ( 0.2500 0.0000i): 0.0625
}

func ExampleQubit_State_grouping() {
	v := qubit.Pluses(4)
	for _, s := range v.State([]int{0}, []int{1, 2, 3}) {
		fmt.Println(s)
	}

	for _, s := range v.State([]int{1, 2, 3}, []int{0}) {
		fmt.Println(s)
	}

	// Output:
	// [0 000] ( 0.2500 0.0000i): 0.0625
	// [0 001] ( 0.2500 0.0000i): 0.0625
	// [0 010] ( 0.2500 0.0000i): 0.0625
	// [0 011] ( 0.2500 0.0000i): 0.0625
	// [0 100] ( 0.2500 0.0000i): 0.0625
	// [0 101] ( 0.2500 0.0000i): 0.0625
	// [0 110] ( 0.2500 0.0000i): 0.0625
	// [0 111] ( 0.2500 0.0000i): 0.0625
	// [1 000] ( 0.2500 0.0000i): 0.0625
	// [1 001] ( 0.2500 0.0000i): 0.0625
	// [1 010] ( 0.2500 0.0000i): 0.0625
	// [1 011] ( 0.2500 0.0000i): 0.0625
	// [1 100] ( 0.2500 0.0000i): 0.0625
	// [1 101] ( 0.2500 0.0000i): 0.0625
	// [1 110] ( 0.2500 0.0000i): 0.0625
	// [1 111] ( 0.2500 0.0000i): 0.0625
	// [000 0] ( 0.2500 0.0000i): 0.0625
	// [001 0] ( 0.2500 0.0000i): 0.0625
	// [010 0] ( 0.2500 0.0000i): 0.0625
	// [011 0] ( 0.2500 0.0000i): 0.0625
	// [100 0] ( 0.2500 0.0000i): 0.0625
	// [101 0] ( 0.2500 0.0000i): 0.0625
	// [110 0] ( 0.2500 0.0000i): 0.0625
	// [111 0] ( 0.2500 0.0000i): 0.0625
	// [000 1] ( 0.2500 0.0000i): 0.0625
	// [001 1] ( 0.2500 0.0000i): 0.0625
	// [010 1] ( 0.2500 0.0000i): 0.0625
	// [011 1] ( 0.2500 0.0000i): 0.0625
	// [100 1] ( 0.2500 0.0000i): 0.0625
	// [101 1] ( 0.2500 0.0000i): 0.0625
	// [110 1] ( 0.2500 0.0000i): 0.0625
	// [111 1] ( 0.2500 0.0000i): 0.0625
}

func Example_bell() {
	q := qubit.Zeros(2).Apply(
		gate.From("HI"),
		gate.CNOT(2, 0, 1),
	)

	for _, s := range q.State() {
		fmt.Println(s)
	}

	// Output:
	// [00] ( 0.7071 0.0000i): 0.5000
	// [11] ( 0.7071 0.0000i): 0.5000
}

func Example_grover2() {
	oracle := gate.CZ(2, 0, 1)
	diffuser := matrix.Apply(
		gate.H(2),
		gate.X(2),
		gate.CZ(2, 0, 1),
		gate.X(2),
		gate.H(2),
	)

	q := qubit.Zeros(2).Apply(
		gate.H(2),
		oracle,
		diffuser,
	)

	q.Measure(0)
	q.Measure(1)

	for _, s := range q.State() {
		fmt.Println(s)
	}

	// Output:
	// [11] (-1.0000 0.0000i): 1.0000
}

func Example_grover3() {
	oracle := matrix.Apply(
		gate.From("XIII"),
		gate.ControlledNot(4, []int{0, 1, 2}, 3),
		gate.From("XIII"),
	)

	diffuser := matrix.Apply(
		gate.From("HHHH"),
		gate.From("XXXI"),
		matrix.TensorProduct(gate.ControlledZ(3, []int{0, 1}, 2), gate.I()),
		gate.From("XXXI"),
		gate.From("HHHI"),
	)

	q := qubit.TensorProduct(
		qubit.Zeros(3),
		qubit.One(),
	).Apply(
		gate.H(4),
		oracle,
		diffuser,
	)

	for _, s := range q.State() {
		fmt.Println(s)
	}

	// Output:
	// [0001] (-0.1768 0.0000i): 0.0313
	// [0011] (-0.1768 0.0000i): 0.0313
	// [0101] (-0.1768 0.0000i): 0.0313
	// [0111] (-0.8839 0.0000i): 0.7813
	// [1001] (-0.1768 0.0000i): 0.0313
	// [1011] (-0.1768 0.0000i): 0.0313
	// [1101] (-0.1768 0.0000i): 0.0313
	// [1111] (-0.1768 0.0000i): 0.0313
}

func Example_bitFlip() {
	psi := qubit.New(vector.New(1, 2))

	// encoding
	psi.TensorProduct(qubit.Zeros(2))
	psi.Apply(
		gate.CNOT(3, 0, 1),
		gate.CNOT(3, 0, 2),
	)

	// error: first qubit is flipped
	psi.Apply(gate.From("XII"))

	// add ancilla qubit
	psi.TensorProduct(qubit.Zeros(2))

	// z1z2
	psi.Apply(
		gate.CNOT(5, 0, 3),
		gate.CNOT(5, 1, 3),
	)

	// z2z3
	psi.Apply(
		gate.CNOT(5, 1, 4),
		gate.CNOT(5, 2, 4),
	)

	// measure
	m3 := psi.Measure(3)
	m4 := psi.Measure(4)

	// recover
	if m3.IsOne() && m4.IsZero() {
		psi.Apply(gate.From("XIIII"))
	}

	if m3.IsOne() && m4.IsOne() {
		psi.Apply(gate.From("IXIII"))
	}

	if m3.IsZero() && m4.IsOne() {
		psi.Apply(gate.From("IIXII"))
	}

	// decoding
	psi.Apply(
		gate.CNOT(5, 0, 2),
		gate.CNOT(5, 0, 1),
	)

	for _, s := range psi.State() {
		fmt.Println(s)
	}

	// Output:
	// [00010] ( 0.4472 0.0000i): 0.2000
	// [10010] ( 0.8944 0.0000i): 0.8000
}

func Example_phaseFlip() {
	psi := qubit.New(vector.New(1, 2))

	// encoding
	psi.TensorProduct(qubit.Zeros(2))
	psi.Apply(
		gate.CNOT(3, 0, 1),
		gate.CNOT(3, 0, 2),
		gate.H(3),
	)

	// error: first qubit is flipped
	psi.Apply(gate.From("ZII"))

	// H
	psi.Apply(gate.H(3))

	// add ancilla qubit
	psi.TensorProduct(qubit.Zeros(2))

	// x1x2
	psi.Apply(
		gate.CNOT(5, 0, 3),
		gate.CNOT(5, 1, 3),
	)

	// x2x3
	psi.Apply(
		gate.CNOT(5, 1, 4),
		gate.CNOT(5, 2, 4),
	)

	// H
	psi.Apply(gate.From("HHHII"))

	// measure
	m3 := psi.Measure(3)
	m4 := psi.Measure(4)

	// recover
	if m3.IsOne() && m4.IsZero() {
		psi.Apply(gate.From("ZIIII"))
	}

	if m3.IsOne() && m4.IsOne() {
		psi.Apply(gate.From("IZIII"))
	}

	if m3.IsZero() && m4.IsOne() {
		psi.Apply(gate.From("IIZII"))
	}

	// decoding
	psi.Apply(
		gate.From("HHHII"),
		gate.CNOT(5, 0, 2),
		gate.CNOT(5, 0, 1),
	)

	for _, s := range psi.State() {
		fmt.Println(s)
	}

	// Output:
	// [00010] ( 0.4472 0.0000i): 0.2000
	// [10010] ( 0.8944 0.0000i): 0.8000
}

func Example_quantumTeleportation() {
	psi := qubit.New(vector.New(1, 2))
	psi.SetRand(rand.Const())

	for _, s := range psi.State() {
		fmt.Println(s)
	}

	bell := qubit.Zeros(2).Apply(
		gate.From("HI"),
		gate.CNOT(2, 0, 1),
	)
	psi.TensorProduct(bell)

	psi.Apply(
		gate.CNOT(3, 0, 1),
		gate.From("HII"),
		gate.CNOT(3, 1, 2),
		gate.CZ(3, 0, 2),
	)

	psi.Measure(0)
	psi.Measure(1)

	for _, s := range psi.State([]int{2}) {
		fmt.Println(s)
	}

	// Output:
	// [0] ( 0.4472 0.0000i): 0.2000
	// [1] ( 0.8944 0.0000i): 0.8000
	// [0] ( 0.4472 0.0000i): 0.2000
	// [1] ( 0.8944 0.0000i): 0.8000
}

func Example_povm() {
	E1 := gate.New(
		[]complex128{0, 0},
		[]complex128{0, 1},
	).Mul(complex(math.Sqrt(2)/(1.0+math.Sqrt(2)), 0))

	E2 := gate.New(
		[]complex128{1, -1},
		[]complex128{-1, 1},
	).Mul(complex(math.Sqrt(2)/(1.0+math.Sqrt(2)), 0)).Mul(complex(0.5, 0))

	E3 := gate.I().Sub(E1).Sub(E2)

	{
		add := E1.Add(E2).Add(E3)
		fmt.Println(add.Equal(gate.I()))
	}

	{
		q0 := qubit.Zero().Apply(E1) // E1|0>
		q1 := qubit.Zero().Apply(E2) // E2|0>
		q2 := qubit.Zero().Apply(E3) // E3|0>

		fmt.Printf("%.4v\n", q0.InnerProduct(qubit.Zero())) // <0|E1|0>
		fmt.Printf("%.4v\n", q1.InnerProduct(qubit.Zero())) // <0|E2|0>
		fmt.Printf("%.4v\n", q2.InnerProduct(qubit.Zero())) // <0|E3|0>
	}

	{
		q0 := qubit.Plus().Apply(E1) // E1|+>
		q1 := qubit.Plus().Apply(E2) // E2|+>
		q2 := qubit.Plus().Apply(E3) // E3|+>

		fmt.Printf("%.4v\n", q0.InnerProduct(qubit.Plus())) // <+|E1|+>
		fmt.Printf("%.4v\n", q1.InnerProduct(qubit.Plus())) // <+|E2|+>
		fmt.Printf("%.4v\n", q2.InnerProduct(qubit.Plus())) // <+|E3|+>
	}

	// Output:
	// true
	// (0+0i)
	// (0.2929+0i)
	// (0.7071+0i)
	// (0.2929+0i)
	// (0+0i)
	// (0.7071+0i)
}

func Example_round() {
	qb := qubit.New(vector.New(
		complex(1e-15, 0.5),
		complex(0.5, 1e-15),
	))
	for _, s := range qb.State() {
		fmt.Println(s)
	}

	// Output:
	// [0] ( 0.0000 0.7071i): 0.5000
	// [1] ( 0.7071 0.0000i): 0.5000
}

func ExampleQubit_Rand() {
	qb := qubit.New(vector.New())
	qb.SetRand(rand.Const())

	for range 10 {
		fmt.Println(qb.Rand()())
	}

	// Output:
	// 0.9999275824802834
	// 0.8856419373528862
	// 0.38147752771154886
	// 0.4812673234167829
	// 0.44417259544314847
	// 0.5210016660132573
	// 0.8861088591612437
	// 0.6769530468231688
	// 0.9850412088603281
	// 0.98505615011337
}

func TestBloch(t *testing.T) {
	cases := []struct {
		theta, phi float64
		want       *qubit.Qubit
	}{
		{
			theta: 0,
			phi:   0,
			want:  qubit.Zero(),
		},
		{
			theta: math.Pi,
			phi:   0,
			want:  qubit.One(),
		},
		{
			theta: math.Pi / 2,
			phi:   0,
			want:  qubit.Plus(),
		},
		{
			theta: math.Pi / 2,
			phi:   math.Pi,
			want:  qubit.Minus(),
		},
		{
			theta: math.Pi / 2,
			phi:   math.Pi / 2,
			want:  qubit.Zero().Apply(gate.H(), gate.S()),
		},
	}

	for _, c := range cases {
		got := qubit.Bloch(c.theta, c.phi)
		if !got.Equal(c.want) {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}

func TestNumQubits(t *testing.T) {
	for i := 1; i < 10; i++ {
		if qubit.Zeros(i).NumQubits() != i {
			t.Fail()
		}
	}
}

func TestIsZero(t *testing.T) {
	cases := []struct {
		in   *qubit.Qubit
		want bool
	}{
		{qubit.Zero(), true},
		{qubit.One(), false},
	}

	for _, c := range cases {
		if c.in.IsZero() != c.want {
			t.Fail()
		}
	}
}

func TestIsOne(t *testing.T) {
	cases := []struct {
		in   *qubit.Qubit
		want bool
	}{
		{qubit.Zero(), false},
		{qubit.One(), true},
	}

	for _, c := range cases {
		if c.in.IsOne() != c.want {
			t.Fail()
		}
	}
}

func TestNormalize(t *testing.T) {
	cases := []struct {
		in   *qubit.Qubit
		want float64
	}{
		{qubit.Zero(), 1.0},
		{qubit.One(), 1.0},
		{qubit.New(vector.New(4, 5)), 1.0},
		{qubit.New(vector.New(10, 5)), 1.0},
	}

	for _, c := range cases {
		got := number.Sum(c.in.Probability())
		if !epsilon.IsCloseF64(got, c.want) {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}

func TestMeasure(t *testing.T) {
	q := qubit.Pluses(3)
	for _, p := range q.Probability() {
		if p != 0 && !epsilon.IsCloseF64(p, 0.125) {
			t.Errorf("probability=%v", q.Probability())
		}
	}

	q.Measure(0)
	for _, p := range q.Probability() {
		if p != 0 && !epsilon.IsCloseF64(p, 0.25) {
			t.Errorf("probability=%v", q.Probability())
		}
	}

	q.Measure(1)
	for _, p := range q.Probability() {
		if p != 0 && !epsilon.IsCloseF64(p, 0.5) {
			t.Errorf("probability=%v", q.Probability())
		}
	}

	q.Measure(2)
	for _, p := range q.Probability() {
		if p != 0 && p != 1 {
			t.Error(q.Probability())
		}
	}
}

func TestClone(t *testing.T) {
	in := qubit.Pluses(2)
	got := in.Clone()

	if !in.Equal(got) {
		t.Fail()
	}
}

func TestBinaryString(t *testing.T) {
	cases := []struct {
		in   *qubit.Qubit
		want string
	}{
		{qubit.Zeros(3), "000"},
		{qubit.Ones(3), "111"},
	}

	for _, c := range cases {
		if c.in.BinaryString() != c.want {
			t.Fail()
		}
	}
}

func TestString(t *testing.T) {
	cases := []struct {
		in   *qubit.Qubit
		want string
	}{
		{qubit.Zeros(2), "[(1+0i) (0+0i) (0+0i) (0+0i)]"},
		{qubit.Ones(2), "[(0+0i) (0+0i) (0+0i) (1+0i)]"},
	}

	for _, c := range cases {
		if c.in.String() != c.want {
			t.Fail()
		}
	}
}
