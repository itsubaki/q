package qubit

import (
	"math"
	"testing"

	"github.com/itsubaki/q/pkg/quantum/gate"
)

func TestOperatorSum(t *testing.T) {
	q0 := Zero()
	q1 := Zero()

	q := q0.OuterProduct(q1)
	e := gate.X().Dagger().Apply(q.Apply(gate.X()))

	if q[0][0] != complex(1, 0) {
		t.Error(q)
	}
	if e[1][1] != complex(1, 0) {
		t.Error(e)
	}
}

func TestOuterProduct(t *testing.T) {
	v0 := New(1, 0)
	out := v0.OuterProduct(v0)

	if out[0][0] != complex(1, 0) {
		t.Fail()
	}

	if out[1][0] != complex(0, 0) {
		t.Fail()
	}

	if out[0][1] != complex(0, 0) {
		t.Fail()
	}

	if out[1][1] != complex(0, 0) {
		t.Fail()
	}
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
	f0 := Zero().Fidelity(Zero())
	if f0 != 1 {
		t.Error(f0)
	}

	f1 := Zero().Fidelity(One())
	if f1 != 0 {
		t.Error(f1)
	}
}

func TestTraceDistance(t *testing.T) {
	d0 := Zero().TraceDistance(Zero())
	if d0 != 0 {
		t.Error(d0)
	}

	d1 := Zero().TraceDistance(One())
	if d1 != 1 {
		t.Error(d1)
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

	q.MeasureAll()

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

	bell := Zero(2).Apply(g0.Apply(g1))
	bell.MeasureAll()

	if bell.Probability()[1] != 0 {
		t.Error(bell)
	}

	if bell.Probability()[2] != 0 {
		t.Error(bell)
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
