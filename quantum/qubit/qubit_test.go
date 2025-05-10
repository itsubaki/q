package qubit_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/itsubaki/q/math/epsilon"
	"github.com/itsubaki/q/math/matrix"
	"github.com/itsubaki/q/math/number"
	"github.com/itsubaki/q/math/rand"
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

func Example_pOVM() {
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
	fmt.Println(add.Equals(gate.I()))

	{
		q0 := qubit.Zero().Apply(E1) // E1|0>
		q1 := qubit.Zero().Apply(E2) // E2|0>
		q2 := qubit.Zero().Apply(E3) // E3|0>

		fmt.Println("zero:")
		fmt.Println(q0.InnerProduct(qubit.Zero())) // <0|E1|0>
		fmt.Println(q1.InnerProduct(qubit.Zero())) // <0|E2|0>
		fmt.Println(q2.InnerProduct(qubit.Zero())) // <0|E3|0>
	}

	{
		q0 := qubit.Zero().Apply(gate.H()).Apply(E1) // E1|+>
		q1 := qubit.Zero().Apply(gate.H()).Apply(E2) // E2|+>
		q2 := qubit.Zero().Apply(gate.H()).Apply(E3) // E3|+>

		fmt.Println("H(zero):")
		fmt.Println(q0.InnerProduct(qubit.Zero().Apply(gate.H()))) // <+|E1|+>
		fmt.Println(q1.InnerProduct(qubit.Zero().Apply(gate.H()))) // <+|E2|+>
		fmt.Println(q2.InnerProduct(qubit.Zero().Apply(gate.H()))) // <+|E3|+>
	}

	// Output:
	// true
	// zero:
	// (0+0i)
	// (0.29289321881345254+0i)
	// (0.7071067811865475+0i)
	// H(zero):
	// (0.29289321881345254+0i)
	// (0+0i)
	// (0.7071067811865477+0i)
}

func Example_bellState() {
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

func Example_grover2qubit() {
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

func Example_grover3qubit() {
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

func Example_errorCorrectionBitFlip() {
	phi := qubit.New(1, 2)

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

func Example_errorCorrectionPhaseFlip() {
	phi := qubit.New(1, 2)

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

func Example_quantumTeleportation() {
	phi := qubit.New(1, 2)
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

func Example_quantumTeleportation2() {
	phi := qubit.New(1, 2)
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

func ExampleQubit_OuterProduct() {
	v := qubit.New(1, 0)
	op := v.OuterProduct(v)
	fmt.Println(op.Data)

	// Output:
	// [[(1+0i) (0+0i)] [(0+0i) (0+0i)]]
}

func ExampleQubit_OuterProduct_operatorSum() {
	v := qubit.New(1, 0)
	q := v.OuterProduct(v)
	e := gate.X().Dagger().Apply(q.Apply(gate.X()))

	fmt.Println(e.Data)

	// Output:
	// [[(0+0i) (0+0i)] [(0+0i) (1+0i)]]
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

func TestFidelity(t *testing.T) {
	cases := []struct {
		q0, q1 *qubit.Qubit
		want   float64
	}{
		{qubit.Zero(), qubit.Zero(), 1.0},
		{qubit.Zero(), qubit.One(), 0.0},
	}

	for _, c := range cases {
		got := c.q0.Fidelity(c.q1)
		if got != c.want {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}

func TestTraceDistance(t *testing.T) {
	cases := []struct {
		q0, q1 *qubit.Qubit
		want   float64
	}{
		{qubit.Zero(), qubit.Zero(), 0.0},
		{qubit.Zero(), qubit.One(), 1.0},
	}

	for _, c := range cases {
		got := c.q0.TraceDistance(c.q1)
		if got != c.want {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}

func TestNormalize(t *testing.T) {
	cases := []struct {
		in   *qubit.Qubit
		want float64
		eps  float64
	}{
		{qubit.New(1, 0), 1.0, epsilon.E13()},
		{qubit.New(0, 1), 1.0, epsilon.E13()},
		{qubit.New(4, 5), 1.0, epsilon.E13()},
		{qubit.New(10, 5), 1.0, epsilon.E13()},
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
