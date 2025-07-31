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
	fmt.Print("z : ")
	for _, s := range z.State() {
		fmt.Println(s)
	}

	z2 := qubit.Zero(2)
	fmt.Print("z2:")
	for _, s := range z2.State() {
		fmt.Println(s)
	}

	// Output:
	// z : [0][  0]( 1.0000 0.0000i): 1.0000
	// z2:[00][  0]( 1.0000 0.0000i): 1.0000
}

func ExampleNewFrom() {
	z := qubit.NewFrom("0011")
	for _, s := range z.State() {
		fmt.Println(s)
	}

	// Output:
	// [0011][  3]( 1.0000 0.0000i): 1.0000
}

func ExampleNewFrom_plus() {
	z := qubit.NewFrom("+-")
	for _, s := range z.State() {
		fmt.Println(s)
	}

	// Output:
	// [00][  0]( 0.5000 0.0000i): 0.2500
	// [01][  1](-0.5000 0.0000i): 0.2500
	// [10][  2]( 0.5000 0.0000i): 0.2500
	// [11][  3](-0.5000 0.0000i): 0.2500
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

func ExampleQubit_OuterProduct_operatorSum() {
	v := qubit.Zero()
	q := v.OuterProduct(v)
	e := gate.X().Dagger().Apply(q.Apply(gate.X()))

	for _, r := range e.Seq2() {
		fmt.Println(r)
	}

	// Output:
	// [(0+0i) (0+0i)]
	// [(0+0i) (1+0i)]
}

func ExampleQubit_ApplyAt() {
	qb := qubit.Zero(2)

	h := gate.H()
	cnot := gate.CNOT(2, 0, 1)

	qb.ApplyAt(h, 0)
	qb.Apply(cnot)

	for _, s := range qb.State() {
		fmt.Println(s)
	}

	// Output:
	// [00][  0]( 0.7071 0.0000i): 0.5000
	// [11][  3]( 0.7071 0.0000i): 0.5000
}

func ExampleQubit_U() {
	qb := qubit.Zero(2)
	qb.U(math.Pi/2, 0, 0, 0)
	qb.U(math.Pi/2, 0, 0, 1)

	for _, s := range qb.State() {
		fmt.Println(s)
	}

	// Output:
	// [00][  0]( 0.5000 0.0000i): 0.2500
	// [01][  1]( 0.5000 0.0000i): 0.2500
	// [10][  2]( 0.5000 0.0000i): 0.2500
	// [11][  3]( 0.5000 0.0000i): 0.2500
}

func ExampleQubit_I() {
	qb := qubit.Zero()
	qb.H(0)
	qb.I(0)

	for _, s := range qb.State() {
		fmt.Println(s)
	}

	// Output:
	// [0][  0]( 0.7071 0.0000i): 0.5000
	// [1][  1]( 0.7071 0.0000i): 0.5000
}

func ExampleQubit_H() {
	qb := qubit.Zero()
	qb.H(0)

	for _, s := range qb.State() {
		fmt.Println(s)
	}

	// Output:
	// [0][  0]( 0.7071 0.0000i): 0.5000
	// [1][  1]( 0.7071 0.0000i): 0.5000
}

func ExampleQubit_X() {
	qb := qubit.Zero()
	qb.X(0)

	for _, s := range qb.State() {
		fmt.Println(s)
	}

	// Output:
	// [1][  1]( 1.0000 0.0000i): 1.0000
}

func ExampleQubit_Y() {
	qb := qubit.Zero()
	qb.H(0)
	qb.Y(0)

	for _, s := range qb.State() {
		fmt.Println(s)
	}

	// Output:
	// [0][  0]( 0.0000-0.7071i): 0.5000
	// [1][  1]( 0.0000 0.7071i): 0.5000
}

func ExampleQubit_Z() {
	qb := qubit.Zero()
	qb.H(0)
	qb.Z(0)

	for _, s := range qb.State() {
		fmt.Println(s)
	}

	// Output:
	// [0][  0]( 0.7071 0.0000i): 0.5000
	// [1][  1](-0.7071 0.0000i): 0.5000
}

func ExampleQubit_R() {
	qb := qubit.Zero()
	qb.H(0)
	qb.R(math.Pi, 0)

	for _, s := range qb.State() {
		fmt.Println(s)
	}

	// Output:
	// [0][  0]( 0.7071 0.0000i): 0.5000
	// [1][  1](-0.7071 0.0000i): 0.5000
}

func ExampleQubit_S() {
	qb := qubit.Zero()
	qb.H(0)
	qb.S(0)

	for _, s := range qb.State() {
		fmt.Println(s)
	}

	// Output:
	// [0][  0]( 0.7071 0.0000i): 0.5000
	// [1][  1]( 0.0000 0.7071i): 0.5000
}

func ExampleQubit_T() {
	qb := qubit.Zero()
	qb.H(0)
	qb.T(0)

	for _, s := range qb.State() {
		fmt.Println(s)
	}

	// Output:
	// [0][  0]( 0.7071 0.0000i): 0.5000
	// [1][  1]( 0.5000 0.5000i): 0.5000
}

func ExampleQubit_RX() {
	qb := qubit.Zero()
	qb.H(0)
	qb.RX(math.Pi, 0)

	for _, s := range qb.State() {
		fmt.Println(s)
	}

	// Output:
	// [0][  0]( 0.0000-0.7071i): 0.5000
	// [1][  1]( 0.0000-0.7071i): 0.5000
}

func ExampleQubit_RY() {
	qb := qubit.Zero()
	qb.H(0)
	qb.RY(math.Pi, 0)

	for _, s := range qb.State() {
		fmt.Println(s)
	}

	// Output:
	// [0][  0](-0.7071 0.0000i): 0.5000
	// [1][  1]( 0.7071 0.0000i): 0.5000
}

func ExampleQubit_RZ() {
	qb := qubit.Zero()
	qb.H(0)
	qb.RZ(math.Pi, 0)

	for _, s := range qb.State() {
		fmt.Println(s)
	}

	// Output:
	// [0][  0]( 0.0000-0.7071i): 0.5000
	// [1][  1]( 0.0000 0.7071i): 0.5000
}

func ExampleQubit_C() {
	qb := qubit.Zero(2)
	qb.H(0)
	qb.C(gate.X(), 0, 1)

	for _, s := range qb.State() {
		fmt.Println(s)
	}

	// Output:
	// [00][  0]( 0.7071 0.0000i): 0.5000
	// [11][  3]( 0.7071 0.0000i): 0.5000
}

func ExampleQubit_CU() {
	qb := qubit.Zero(2)
	qb.H(0)
	qb.CU(math.Pi, 0, math.Pi, 0, 1)

	for _, s := range qb.State() {
		fmt.Println(s)
	}

	// Output:
	// [00][  0]( 0.7071 0.0000i): 0.5000
	// [11][  3]( 0.7071 0.0000i): 0.5000
}

func ExampleQubit_CH() {
	qb := qubit.Zero(2)
	qb.H(0)
	qb.CH(0, 1)

	for _, s := range qb.State() {
		fmt.Println(s)
	}

	// Output:
	// [00][  0]( 0.7071 0.0000i): 0.5000
	// [10][  2]( 0.5000 0.0000i): 0.2500
	// [11][  3]( 0.5000 0.0000i): 0.2500
}

func ExampleQubit_CX() {
	qb := qubit.Zero(2)
	qb.H(0)
	qb.CX(0, 1)

	for _, s := range qb.State() {
		fmt.Println(s)
	}

	// Output:
	// [00][  0]( 0.7071 0.0000i): 0.5000
	// [11][  3]( 0.7071 0.0000i): 0.5000
}

func ExampleQubit_CZ() {
	qb := qubit.Zero(2)
	qb.H(0)
	qb.H(1)
	qb.CZ(0, 1)

	for _, s := range qb.State() {
		fmt.Println(s)
	}

	// Output:
	// [00][  0]( 0.5000 0.0000i): 0.2500
	// [01][  1]( 0.5000 0.0000i): 0.2500
	// [10][  2]( 0.5000 0.0000i): 0.2500
	// [11][  3](-0.5000 0.0000i): 0.2500
}

func ExampleQubit_ControlledH() {
	qb := qubit.Zero(2)
	qb.H(0)
	qb.ControlledH([]int{0}, 1)

	for _, s := range qb.State() {
		fmt.Println(s)
	}

	// Output:
	// [00][  0]( 0.7071 0.0000i): 0.5000
	// [10][  2]( 0.5000 0.0000i): 0.2500
	// [11][  3]( 0.5000 0.0000i): 0.2500
}

func ExampleQubit_ControlledX() {
	qb := qubit.Zero(3)
	qb.X(0)
	qb.X(1)
	qb.ControlledX([]int{0, 1}, 2)

	for _, s := range qb.State() {
		fmt.Println(s)
	}

	// Output:
	// [111][  7]( 1.0000 0.0000i): 1.0000
}

func ExampleQubit_ControlledZ() {
	qb := qubit.Zero(2)
	qb.H(0)
	qb.H(1)
	qb.ControlledZ([]int{0}, 1)

	for _, s := range qb.State() {
		fmt.Println(s)
	}

	// Output:
	// [00][  0]( 0.5000 0.0000i): 0.2500
	// [01][  1]( 0.5000 0.0000i): 0.2500
	// [10][  2]( 0.5000 0.0000i): 0.2500
	// [11][  3](-0.5000 0.0000i): 0.2500
}

func ExampleQubit_QFT() {
	qb := qubit.Zero(3)
	qb.X(2)
	qb.QFT()
	qb.Swap(0, 2)

	for _, s := range qb.State() {
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

func ExampleQubit_InvQFT() {
	qb := qubit.Zero(3)
	qb.X(2)
	qb.QFT()
	qb.InvQFT()

	for _, s := range qb.State() {
		fmt.Println(s)
	}

	// Output:
	// [001][  1]( 1.0000 0.0000i): 1.0000
}

func ExampleQubit_Update() {
	qb := qubit.Zero(2)
	qb.Update(vector.New(1, 0, 0, 1))

	for _, s := range qb.State() {
		fmt.Println(s)
	}

	// Output:
	// [00][  0]( 0.7071 0.0000i): 0.5000
	// [11][  3]( 0.7071 0.0000i): 0.5000
}

func ExampleQubit_State() {
	v := qubit.Zero(4).Apply(gate.H(4))

	for _, s := range v.State([]int{0, 1, 2, 3}) {
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

func ExampleQubit_State_order() {
	v := qubit.Zero(4).Apply(gate.H(4))

	for _, s := range v.State([]int{0}, []int{1, 2, 3}) {
		fmt.Println(s)
	}

	for _, s := range v.State([]int{1, 2, 3}, []int{0}) {
		fmt.Println(s)
	}

	// Output:
	// [0 000][  0   0]( 0.2500 0.0000i): 0.0625
	// [0 001][  0   1]( 0.2500 0.0000i): 0.0625
	// [0 010][  0   2]( 0.2500 0.0000i): 0.0625
	// [0 011][  0   3]( 0.2500 0.0000i): 0.0625
	// [0 100][  0   4]( 0.2500 0.0000i): 0.0625
	// [0 101][  0   5]( 0.2500 0.0000i): 0.0625
	// [0 110][  0   6]( 0.2500 0.0000i): 0.0625
	// [0 111][  0   7]( 0.2500 0.0000i): 0.0625
	// [1 000][  1   0]( 0.2500 0.0000i): 0.0625
	// [1 001][  1   1]( 0.2500 0.0000i): 0.0625
	// [1 010][  1   2]( 0.2500 0.0000i): 0.0625
	// [1 011][  1   3]( 0.2500 0.0000i): 0.0625
	// [1 100][  1   4]( 0.2500 0.0000i): 0.0625
	// [1 101][  1   5]( 0.2500 0.0000i): 0.0625
	// [1 110][  1   6]( 0.2500 0.0000i): 0.0625
	// [1 111][  1   7]( 0.2500 0.0000i): 0.0625
	// [000 0][  0   0]( 0.2500 0.0000i): 0.0625
	// [001 0][  1   0]( 0.2500 0.0000i): 0.0625
	// [010 0][  2   0]( 0.2500 0.0000i): 0.0625
	// [011 0][  3   0]( 0.2500 0.0000i): 0.0625
	// [100 0][  4   0]( 0.2500 0.0000i): 0.0625
	// [101 0][  5   0]( 0.2500 0.0000i): 0.0625
	// [110 0][  6   0]( 0.2500 0.0000i): 0.0625
	// [111 0][  7   0]( 0.2500 0.0000i): 0.0625
	// [000 1][  0   1]( 0.2500 0.0000i): 0.0625
	// [001 1][  1   1]( 0.2500 0.0000i): 0.0625
	// [010 1][  2   1]( 0.2500 0.0000i): 0.0625
	// [011 1][  3   1]( 0.2500 0.0000i): 0.0625
	// [100 1][  4   1]( 0.2500 0.0000i): 0.0625
	// [101 1][  5   1]( 0.2500 0.0000i): 0.0625
	// [110 1][  6   1]( 0.2500 0.0000i): 0.0625
	// [111 1][  7   1]( 0.2500 0.0000i): 0.0625
}

func Example_bell() {
	q := qubit.Zero(2).Apply(
		gate.H().TensorProduct(gate.I()),
		gate.CNOT(2, 0, 1),
	)

	for _, s := range q.State() {
		fmt.Println(s)
	}

	// Output:
	// [00][  0]( 0.7071 0.0000i): 0.5000
	// [11][  3]( 0.7071 0.0000i): 0.5000
}

func Example_grover2() {
	oracle := gate.CZ(2, 0, 1)
	amp := matrix.Apply(
		gate.H(2),
		gate.X(2),
		gate.CZ(2, 0, 1),
		gate.X(2),
		gate.H(2),
	)

	q := qubit.Zero(2).Apply(
		gate.H(2),
		oracle,
		amp,
	)

	q.Measure(0)
	q.Measure(1)

	for _, s := range q.State() {
		fmt.Println(s)
	}

	// Output:
	// [11][  3](-1.0000 0.0000i): 1.0000
}

func Example_grover3() {
	oracle := matrix.Apply(
		matrix.TensorProduct(gate.X(), gate.I(3)),
		gate.ControlledNot(4, []int{0, 1, 2}, 3),
		matrix.TensorProduct(gate.X(), gate.I(3)),
	)

	amp := matrix.Apply(
		matrix.TensorProduct(gate.H(3), gate.H()),
		matrix.TensorProduct(gate.X(3), gate.I()),
		matrix.TensorProduct(gate.ControlledZ(3, []int{0, 1}, 2), gate.I()),
		matrix.TensorProduct(gate.X(3), gate.I()),
		matrix.TensorProduct(gate.H(3), gate.I()),
	)

	q := qubit.TensorProduct(
		qubit.Zero(3),
		qubit.One(),
	).Apply(
		gate.H(4),
		oracle,
		amp,
	)

	for _, s := range q.State() {
		fmt.Println(s)
	}

	// Output:
	// [0001][  1](-0.1768 0.0000i): 0.0313
	// [0011][  3](-0.1768 0.0000i): 0.0313
	// [0101][  5](-0.1768 0.0000i): 0.0313
	// [0111][  7](-0.8839 0.0000i): 0.7813
	// [1001][  9](-0.1768 0.0000i): 0.0313
	// [1011][ 11](-0.1768 0.0000i): 0.0313
	// [1101][ 13](-0.1768 0.0000i): 0.0313
	// [1111][ 15](-0.1768 0.0000i): 0.0313
}

func Example_eccBitFlip() {
	phi := qubit.New(vector.New(1, 2))

	// encoding
	phi.TensorProduct(qubit.Zero(2))
	phi.Apply(
		gate.CNOT(3, 0, 1),
		gate.CNOT(3, 0, 2),
	)

	// error: first qubit is flipped
	phi.Apply(matrix.TensorProduct(gate.X(), gate.I(2)))

	// add ancilla qubit
	phi.TensorProduct(qubit.Zero(2))

	// z1z2
	phi.Apply(
		gate.CNOT(5, 0, 3),
		gate.CNOT(5, 1, 3),
	)

	// z2z3
	phi.Apply(
		gate.CNOT(5, 1, 4),
		gate.CNOT(5, 2, 4),
	)

	// measure
	m3 := phi.Measure(3)
	m4 := phi.Measure(4)

	// recover
	if m3.IsOne() && m4.IsZero() {
		phi.Apply(matrix.TensorProduct(gate.X(), gate.I(4)))
	}

	if m3.IsOne() && m4.IsOne() {
		phi.Apply(matrix.TensorProduct(gate.I(), gate.X(), gate.I(3)))
	}

	if m3.IsZero() && m4.IsOne() {
		phi.Apply(matrix.TensorProduct(gate.I(2), gate.X(), gate.I(2)))
	}

	// decoding
	phi.Apply(
		gate.CNOT(5, 0, 2),
		gate.CNOT(5, 0, 1),
	)

	for _, s := range phi.State() {
		fmt.Println(s)
	}

	// Output:
	// [00010][  2]( 0.4472 0.0000i): 0.2000
	// [10010][ 18]( 0.8944 0.0000i): 0.8000
}

func Example_eccPhaseFlip() {
	phi := qubit.New(vector.New(1, 2))

	// encoding
	phi.TensorProduct(qubit.Zero(2))
	phi.Apply(
		gate.CNOT(3, 0, 1),
		gate.CNOT(3, 0, 2),
		gate.H(3),
	)

	// error: first qubit is flipped
	phi.Apply(matrix.TensorProduct(gate.Z(), gate.I(2)))

	// H
	phi.Apply(gate.H(3))

	// add ancilla qubit
	phi.TensorProduct(qubit.Zero(2))

	// x1x2
	phi.Apply(
		gate.CNOT(5, 0, 3),
		gate.CNOT(5, 1, 3),
	)

	// x2x3
	phi.Apply(
		gate.CNOT(5, 1, 4),
		gate.CNOT(5, 2, 4),
	)

	// H
	phi.Apply(matrix.TensorProduct(gate.H(3), gate.I(2)))

	// measure
	m3 := phi.Measure(3)
	m4 := phi.Measure(4)

	// recover
	if m3.IsOne() && m4.IsZero() {
		phi.Apply(matrix.TensorProduct(gate.Z(), gate.I(4)))
	}

	if m3.IsOne() && m4.IsOne() {
		phi.Apply(matrix.TensorProduct(gate.I(), gate.Z(), gate.I(3)))
	}

	if m3.IsZero() && m4.IsOne() {
		phi.Apply(matrix.TensorProduct(gate.I(2), gate.Z(), gate.I(2)))
	}

	// decoding
	phi.Apply(
		matrix.TensorProduct(gate.H(3), gate.I(2)),
		gate.CNOT(5, 0, 2),
		gate.CNOT(5, 0, 1),
	)

	for _, s := range phi.State() {
		fmt.Println(s)
	}

	// Output:
	// [00010][  2]( 0.4472 0.0000i): 0.2000
	// [10010][ 18]( 0.8944 0.0000i): 0.8000
}
func Example_teleportation() {
	phi := qubit.New(vector.New(1, 2))
	phi.Rand = rand.Const()

	fmt.Println("before:")
	for _, s := range phi.State() {
		fmt.Println(s)
	}

	bell := qubit.Zero(2).Apply(
		matrix.TensorProduct(gate.H(), gate.I()),
		gate.CNOT(2, 0, 1),
	)
	phi.TensorProduct(bell)

	phi.Apply(
		gate.CNOT(3, 0, 1),
		matrix.TensorProduct(gate.H(), gate.I(2)),
		gate.CNOT(3, 1, 2),
		gate.CZ(3, 0, 2),
	)

	phi.Measure(0)
	phi.Measure(1)

	fmt.Println("after:")
	for _, s := range phi.State() {
		fmt.Println(s)
	}

	// Output:
	// before:
	// [0][  0]( 0.4472 0.0000i): 0.2000
	// [1][  1]( 0.8944 0.0000i): 0.8000
	// after:
	// [110][  6]( 0.4472 0.0000i): 0.2000
	// [111][  7]( 0.8944 0.0000i): 0.8000
}

func Example_teleportationCond() {
	phi := qubit.New(vector.New(1, 2))
	phi.Rand = rand.Const()

	fmt.Println("before:")
	for _, s := range phi.State() {
		fmt.Println(s)
	}

	bell := qubit.Zero(2).Apply(
		matrix.TensorProduct(gate.H(), gate.I()),
		gate.CNOT(2, 0, 1),
	)
	phi.TensorProduct(bell)

	phi.Apply(
		gate.CNOT(3, 0, 1),
		matrix.TensorProduct(gate.H(), gate.I(2)),
	)

	mz := phi.Measure(0)
	mx := phi.Measure(1)

	if mx.IsOne() {
		phi.Apply(matrix.TensorProduct(gate.I(2), gate.X()))
	}

	if mz.IsOne() {
		phi.Apply(matrix.TensorProduct(gate.I(2), gate.Z()))
	}

	fmt.Println("after:")
	for _, s := range phi.State() {
		fmt.Println(s)
	}

	// Output:
	// before:
	// [0][  0]( 0.4472 0.0000i): 0.2000
	// [1][  1]( 0.8944 0.0000i): 0.8000
	// after:
	// [110][  6]( 0.4472 0.0000i): 0.2000
	// [111][  7]( 0.8944 0.0000i): 0.8000
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

	add := E1.Add(E2).Add(E3)
	fmt.Println("euqlas:", add.Equals(gate.I()))

	{
		q0 := qubit.Zero().Apply(E1) // E1|0>
		q1 := qubit.Zero().Apply(E2) // E2|0>
		q2 := qubit.Zero().Apply(E3) // E3|0>

		fmt.Println("zero:")
		fmt.Printf("%.4v\n", q0.InnerProduct(qubit.Zero())) // <0|E1|0>
		fmt.Printf("%.4v\n", q1.InnerProduct(qubit.Zero())) // <0|E2|0>
		fmt.Printf("%.4v\n", q2.InnerProduct(qubit.Zero())) // <0|E3|0>
	}

	{
		q0 := qubit.Plus().Apply(E1) // E1|+>
		q1 := qubit.Plus().Apply(E2) // E2|+>
		q2 := qubit.Plus().Apply(E3) // E3|+>

		fmt.Println("H(zero):")
		fmt.Printf("%.4v\n", q0.InnerProduct(qubit.Plus())) // <+|E1|+>
		fmt.Printf("%.4v\n", q1.InnerProduct(qubit.Plus())) // <+|E2|+>
		fmt.Printf("%.4v\n", q2.InnerProduct(qubit.Plus())) // <+|E3|+>
	}

	// Output:
	// euqlas: true
	// zero:
	// (0+0i)
	// (0.2929+0i)
	// (0.7071+0i)
	// H(zero):
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
	// [0][  0]( 0.0000 0.7071i): 0.5000
	// [1][  1]( 0.7071 0.0000i): 0.5000
}

func TestNumQubits(t *testing.T) {
	for i := 1; i < 10; i++ {
		if qubit.Zero(i).NumQubits() != i {
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
		eps  float64
	}{
		{qubit.Zero(), 1.0, epsilon.E13()},
		{qubit.One(), 1.0, epsilon.E13()},
		{qubit.New(vector.New(4, 5)), 1.0, epsilon.E13()},
		{qubit.New(vector.New(10, 5)), 1.0, epsilon.E13()},
	}

	for _, c := range cases {
		got := number.Sum(c.in.Probability())
		if math.Abs(got-c.want) > c.eps {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}

func TestMeasure(t *testing.T) {
	eps := epsilon.E13()

	q := qubit.Zero(3).Apply(gate.H(3))
	for _, p := range q.Probability() {
		if p != 0 && math.Abs(p-0.125) > eps {
			t.Errorf("probability=%v", q.Probability())
		}
	}

	q.Measure(0)
	for _, p := range q.Probability() {
		if p != 0 && math.Abs(p-0.25) > eps {
			t.Errorf("probability=%v", q.Probability())
		}
	}

	q.Measure(1)
	for _, p := range q.Probability() {
		if p != 0 && math.Abs(p-0.5) > eps {
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
	in := qubit.Zero(2).Apply(gate.H(2))
	got := in.Clone()

	if !in.Equals(got) {
		t.Fail()
	}
}

func TestInt(t *testing.T) {
	cases := []struct {
		in   *qubit.Qubit
		want int64
	}{
		{qubit.Zero(), 0},
		{qubit.One(), 1},
	}

	for _, c := range cases {
		if c.in.Int() != c.want {
			t.Fail()
		}
	}
}
func TestBinaryString(t *testing.T) {
	cases := []struct {
		in   *qubit.Qubit
		want string
	}{
		{qubit.Zero(3), "000"},
		{qubit.One(3), "111"},
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
		{qubit.Zero(2), "[(1+0i) (0+0i) (0+0i) (0+0i)]"},
		{qubit.One(2), "[(0+0i) (0+0i) (0+0i) (1+0i)]"},
	}

	for _, c := range cases {
		if c.in.String() != c.want {
			t.Fail()
		}
	}
}
