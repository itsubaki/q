package qubit

import (
	"fmt"
	"math"
	"testing"

	"github.com/itsubaki/q/pkg/math/matrix"
	"github.com/itsubaki/q/pkg/math/number"
	"github.com/itsubaki/q/pkg/quantum/gate"
)

func Example_pOVM() {
	E1 := gate.New(
		[]complex128{0, 0},
		[]complex128{0, 1},
	).Mul(complex(math.Sqrt(2)/(1.0+math.Sqrt(2)), 0))

	E2 := gate.New(
		[]complex128{1, -1},
		[]complex128{-1, 1},
	).Mul(complex(0.5, 0)).
		Mul(complex(math.Sqrt(2)/(1.0+math.Sqrt(2)), 0))

	E3 := gate.I().Sub(E1).Sub(E2)

	add := E1.Add(E2).Add(E3)
	fmt.Println(add.Equals(gate.I()))

	q0 := Zero()
	fmt.Println(q0.Apply(E1).InnerProduct(q0))

	q1 := Zero().Apply(gate.H())
	fmt.Println(q1.Apply(E2).InnerProduct(q1))

	// Output:
	// true
	// (0+0i)
	// (0+0i)
}

func Example_bellState() {
	g0 := gate.H().TensorProduct(gate.I())
	g1 := gate.CNOT(2, 0, 1)
	qc := g0.Apply(g1)

	q := Zero(2).Apply(qc)
	fmt.Println(q.Amplitude())

	// Output:
	// [(0.7071067811865476+0i) (0+0i) (0+0i) (0.7071067811865476+0i)]
}

func Example_grover2qubit() {
	oracle := gate.CZ(2, 0, 1)

	h2 := gate.H(2)
	x2 := gate.X(2)
	amp := h2.Apply(x2).Apply(gate.CZ(2, 0, 1)).Apply(x2).Apply(h2)

	qc := h2.Apply(oracle).Apply(amp)
	q := Zero(2).Apply(qc)

	q.Measure(0)
	q.Measure(1)

	for i, p := range q.Probability() {
		if p == 0 {
			continue
		}

		fmt.Printf("%2d: %.5f\n", i, p)
	}

	// Output:
	// 3: 1.00000
}

func Example_quantumTeleportation() {
	phi := New(1, 2)
	for _, p := range phi.Probability() {
		if p == 0 {
			continue
		}

		fmt.Printf("%.5f\n", p)
	}

	g0 := matrix.TensorProduct(gate.H(), gate.I())
	g1 := gate.CNOT(2, 0, 1)
	bell := Zero(2).Apply(g0).Apply(g1)
	phi.TensorProduct(bell)

	g2 := gate.CNOT(3, 0, 1)
	g3 := matrix.TensorProduct(gate.H(), gate.I(2))
	phi.Apply(g2).Apply(g3)

	mz := phi.Measure(0)
	mx := phi.Measure(1)

	if mz.IsOne() {
		z := matrix.TensorProduct(gate.I(2), gate.Z())
		phi.Apply(z)
	}

	if mx.IsOne() {
		x := matrix.TensorProduct(gate.I(2), gate.X())
		phi.Apply(x)
	}

	for _, p := range phi.Probability() {
		if p == 0 {
			continue
		}

		fmt.Printf("%.5f\n", p)
	}

	// Output:
	// 0.20000
	// 0.80000
	// 0.20000
	// 0.80000
}

func Example_quantumTeleportation2() {
	phi := New(1, 2)
	for _, p := range phi.Probability() {
		if p == 0 {
			continue
		}

		fmt.Printf("%.5f\n", p)
	}

	g0 := matrix.TensorProduct(gate.H(), gate.I())
	g1 := gate.CNOT(2, 0, 1)
	bell := Zero(2).Apply(g0).Apply(g1)
	phi.TensorProduct(bell)

	g2 := gate.CNOT(3, 0, 1)
	g3 := matrix.TensorProduct(gate.H(), gate.I(2))
	g4 := gate.CNOT(3, 1, 2)
	g5 := gate.CZ(3, 0, 2)

	phi.Apply(g2).Apply(g3).Apply(g4).Apply(g5)
	phi.Measure(0)
	phi.Measure(1)

	for _, p := range phi.Probability() {
		if p == 0 {
			continue
		}

		fmt.Printf("%.5f\n", p)
	}

	// Output:
	// 0.20000
	// 0.80000
	// 0.20000
	// 0.80000
}

func Example_grover3qubit() {
	x := matrix.TensorProduct(gate.X(), gate.I(3))
	oracle := x.Apply(gate.ControlledNot(4, []int{0, 1, 2}, 3)).Apply(x)

	h4 := matrix.TensorProduct(gate.H(3), gate.H())
	x3 := matrix.TensorProduct(gate.X(3), gate.I())
	cz := matrix.TensorProduct(gate.ControlledZ(3, []int{0, 1}, 2), gate.I())
	h3 := matrix.TensorProduct(gate.H(3), gate.I())
	amp := h4.Apply(x3).Apply(cz).Apply(x3).Apply(h3)

	q := TensorProduct(Zero(3), One())
	q.Apply(gate.H(4)).Apply(oracle).Apply(amp)

	for i, p := range q.Probability() {
		if p == 0 {
			continue
		}

		fmt.Printf("%2d: %.5f\n", i, p)
	}

	// Output:
	//  1: 0.03125
	//  3: 0.03125
	//  5: 0.03125
	//  7: 0.78125
	//  9: 0.03125
	// 11: 0.03125
	// 13: 0.03125
	// 15: 0.03125
}

func Example_errorCorrectionBitFlip() {
	phi := New(1, 2)

	// encoding
	phi.TensorProduct(Zero(2))
	phi.Apply(gate.CNOT(3, 0, 1))
	phi.Apply(gate.CNOT(3, 0, 2))

	// error: first qubit is flipped
	phi.Apply(matrix.TensorProduct(gate.X(), gate.I(2)))

	// add ancilla qubit
	phi.TensorProduct(Zero(2))

	// z1z2
	phi.Apply(gate.CNOT(5, 0, 3)).Apply(gate.CNOT(5, 1, 3))

	// z2z3
	phi.Apply(gate.CNOT(5, 1, 4)).Apply(gate.CNOT(5, 2, 4))

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

	// answer is 0.2|000>|10> + 0.8|111>|10>
	for i, p := range phi.Probability() {
		if p == 0 {
			continue
		}

		fmt.Printf("%2d: %.5f\n", i, p)
	}

	// Output:
	//  2: 0.20000
	// 30: 0.80000
}

func Example_errorCorrectionPhaseFlip() {
	phi := New(1, 2)

	// encoding
	phi.TensorProduct(Zero(2))
	phi.Apply(gate.CNOT(3, 0, 1))
	phi.Apply(gate.CNOT(3, 0, 2))
	phi.Apply(gate.H(3))

	// error: first qubit is flipped
	phi.Apply(matrix.TensorProduct(gate.Z(), gate.I(2)))

	// H
	phi.Apply(gate.H(3))

	// add ancilla qubit
	phi.TensorProduct(Zero(2))

	// x1x2
	phi.Apply(gate.CNOT(5, 0, 3)).Apply(gate.CNOT(5, 1, 3))

	// x2x3
	phi.Apply(gate.CNOT(5, 1, 4)).Apply(gate.CNOT(5, 2, 4))

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

	phi.Apply(matrix.TensorProduct(gate.H(3), gate.I(2)))

	for i, p := range phi.Probability() {
		if math.Abs(p) < 1e-13 {
			continue
		}

		fmt.Printf("%2d: %.5f\n", i, p)
	}

	// Output:
	//  2: 0.20000
	// 30: 0.80000
}

func ExampleQubit_OuterProduct() {
	v := New(1, 0)
	op := v.OuterProduct(v)
	fmt.Println(op)

	// Output:
	// [[(1+0i) (0+0i)] [(0+0i) (0+0i)]]
}

func ExampleQubit_OuterProduct_operatorSum() {
	v := New(1, 0)
	q := v.OuterProduct(v)
	e := gate.X().Dagger().Apply(q.Apply(gate.X()))

	fmt.Println(e)

	// Output:
	// [[(0+0i) (0+0i)] [(0+0i) (1+0i)]]
}

func TestNumberOfBit(t *testing.T) {
	for i := 1; i < 10; i++ {
		if Zero(i).NumberOfBit() != i {
			t.Fail()
		}
	}
}

func TestIsZero(t *testing.T) {
	cases := []struct {
		qb     *Qubit
		isZero bool
		isOne  bool
	}{
		{Zero(), true, false},
		{One(), false, true},
	}

	for _, c := range cases {
		if c.qb.IsZero() != c.isZero {
			t.Fail()
		}
		if c.qb.IsOne() != c.isOne {
			t.Fail()
		}
	}
}

func TestFidelity(t *testing.T) {
	cases := []struct {
		q0, q1 *Qubit
		f      float64
	}{
		{Zero(), Zero(), 1.0},
		{Zero(), One(), 0.0},
	}

	for _, c := range cases {
		f := c.q0.Fidelity(c.q1)
		if f != c.f {
			t.Error(f)
		}
	}
}

func TestTraceDistance(t *testing.T) {
	cases := []struct {
		q0, q1 *Qubit
		d      float64
	}{
		{Zero(), Zero(), 0.0},
		{Zero(), One(), 1.0},
	}

	for _, c := range cases {
		d := c.q0.TraceDistance(c.q1)
		if d != c.d {
			t.Error(d)
		}
	}
}

func TestNormalize(t *testing.T) {
	cases := []struct {
		q *Qubit
	}{
		{New(1, 0)},
		{New(0, 1)},
		{New(4, 5)},
		{New(10, 5)},
	}

	for _, c := range cases {
		sum := number.Sum(c.q.Probability())
		if math.Abs(sum-1.0) > 1e-13 {
			t.Error(sum)
		}
	}
}

func TestMeasure(t *testing.T) {
	q := Zero(3).Apply(gate.H(3))

	q.Measure(0)
	for _, p := range q.Probability() {
		if p != 0 && p != 0.25 {
			t.Error(q.Probability())
		}
	}

	q.Measure(1)
	for _, p := range q.Probability() {
		if p != 0 && p != 0.4999999999999999 {
			t.Error(q.Probability())
		}
	}

	q.Measure(2)
	for _, p := range q.Probability() {
		if p != 0 && p != 1 {
			t.Error(q.Probability())
		}
	}
}
