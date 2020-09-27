package qubit

import (
	"fmt"
	"math"
	"testing"

	"github.com/itsubaki/q/pkg/quantum/gate"
)

func ExampleQubit_OuterProduct_oeprator_sum() {
	q0 := Zero()
	q1 := Zero()

	q := q0.OuterProduct(q1)
	e := gate.X().Dagger().Apply(q.Apply(gate.X()))

	fmt.Println(q)
	fmt.Println(e)

	// Output:
	// [[(1+0i) (0+0i)] [(0+0i) (0+0i)]]
	// [[(0+0i) (0+0i)] [(0+0i) (1+0i)]]
}

func ExampleQubit_OuterProduct() {
	v := New(1, 0)
	op := v.OuterProduct(v)
	fmt.Println(op)

	// Output:
	// [[(1+0i) (0+0i)] [(0+0i) (0+0i)]]
}

func TestQFT(t *testing.T) {
	q := Zero(3).Apply(gate.QFT(3))

	for _, p := range q.Probability() {
		if math.Abs(p-0.125) > 1e-13 {
			t.Error(q.Probability())
		}
	}
}

func TestNum(t *testing.T) {
	for i := 1; i < 10; i++ {
		if Zero(i).NumberOfBit() != i {
			t.Fail()
		}
	}
}

func TestIs(t *testing.T) {
	if !Zero().IsZero() {
		t.Fail()
	}

	if !One().IsOne() {
		t.Fail()
	}

	if Zero().IsOne() {
		t.Fail()
	}

	if One().IsZero() {
		t.Fail()
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

func TestQubit(t *testing.T) {
	q := Zero().Apply(gate.H())
	p := q.Probability()

	if p[0]-0.5 > 1e-13 {
		t.Error(p)
	}
	if p[1]-0.5 > 1e-13 {
		t.Error(p)
	}

	q.Measure(0)

	p2 := q.Probability()
	if p2[0] == 0 && p2[1] != 1 {
		t.Error(p2)
	}

	if p2[0] == 1 && p2[1] != 0 {
		t.Error(p2)
	}
}

func TestBellState(t *testing.T) {
	g0 := gate.H().TensorProduct(gate.I())
	g1 := gate.CNOT(2, 0, 1)

	q := Zero(2).Apply(g0.Apply(g1))
	q.Measure(0)
	q.Measure(1)

	p := q.Probability()
	if p[1] != 0 {
		t.Error(q)
	}

	if p[2] != 0 {
		t.Error(q)
	}
}

func TestNormalize(t *testing.T) {
	q := New(4, 5)

	var sum float64
	for _, p := range q.Probability() {
		sum = sum + p
	}

	if sum != 1 {
		t.Error(sum)
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
