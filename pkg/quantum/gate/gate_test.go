package gate_test

import (
	"testing"

	"github.com/itsubaki/q/pkg/math/matrix"
	"github.com/itsubaki/q/pkg/math/vector"
	"github.com/itsubaki/q/pkg/quantum/gate"
)

func TestCModExp2(t *testing.T) {
	v0 := vector.TensorProduct(vector.TensorProductN(vector.New(1, 0), 6), vector.New(0, 1))
	g0 := gate.CNOT(7, 2, 4).Apply(gate.CNOT(7, 2, 5))
	g1 := gate.CNOT(7, 3, 5).Apply(gate.CCNOT(7, 1, 5, 3)).Apply(gate.CNOT(7, 3, 5))
	g2 := gate.CNOT(7, 4, 6).Apply(gate.CCNOT(7, 1, 6, 4)).Apply(gate.CNOT(7, 4, 6))
	g3 := gate.I(7)

	cases := []struct {
		b, a, j, N int
		c          int
		t          []int
		m          matrix.Matrix
		v          vector.Vector
	}{
		{7, 7, 0, 15, 2, []int{4, 5, 6, 7}, nil, v0.Clone().Apply(g0)},
		{7, 7, 1, 15, 1, []int{4, 5, 6, 7}, g1.Apply(g2), nil},
		{7, 7, 2, 15, 0, []int{4, 5, 6, 7}, g3, nil},
	}

	for _, c := range cases {
		a := gate.CModExp2(c.b, c.a, c.j, c.N, c.c, c.t)
		if !a.IsUnitary() {
			t.Errorf("modexp is not unitary")
		}

		if c.m != nil && !a.Equals(c.m) {
			t.Fail()
		}

		if c.v != nil && !v0.Clone().Apply(a).Equals(c.v) {
			t.Fail()
		}
	}
}

func TestInverse(t *testing.T) {
	cases := []struct {
		m, e matrix.Matrix
	}{
		{gate.U(1, 2, 3, 4), gate.I()},
		{gate.X(2), gate.I(2)},
		{gate.CNOT(2, 0, 1), gate.I(2)},
	}

	for _, c := range cases {
		inv := c.m.Inverse()
		mmi := c.m.Apply(inv)
		if !mmi.Equals(c.e) {
			t.Fail()
		}
	}
}

func TestIsHermite(t *testing.T) {
	cases := []struct {
		m matrix.Matrix
	}{
		{gate.H()},
		{gate.X()},
		{gate.Y()},
		{gate.Z()},
	}

	for _, c := range cases {
		if !c.m.IsHermite() {
			t.Error(c.m)
		}
	}
}

func TestIsUnitary(t *testing.T) {
	cases := []struct {
		m matrix.Matrix
	}{
		{gate.I()},
		{gate.H()},
		{gate.X()},
		{gate.Y()},
		{gate.Z()},
		{gate.U(1, 2, 3, 4)},
	}

	for _, c := range cases {
		if !c.m.IsUnitary() {
			t.Error(c.m)
		}
	}
}

func TestTrace(t *testing.T) {
	cases := []struct {
		m matrix.Matrix
		t complex128
	}{
		{gate.I(), complex(2, 0)},
		{gate.H(), complex(0, 0)},
		{gate.X(), complex(0, 0)},
		{gate.Y(), complex(0, 0)},
		{gate.Z(), complex(0, 0)},
	}

	for _, c := range cases {
		tr := c.m.Trace()
		if tr != c.t {
			t.Error(tr)
		}
	}
}

func TestCZ(t *testing.T) {
	expected := gate.New(
		[]complex128{1, 0, 0, 0, 0, 0, 0, 0},
		[]complex128{0, 1, 0, 0, 0, 0, 0, 0},
		[]complex128{0, 0, 1, 0, 0, 0, 0, 0},
		[]complex128{0, 0, 0, 1, 0, 0, 0, 0},
		[]complex128{0, 0, 0, 0, 1, 0, 0, 0},
		[]complex128{0, 0, 0, 0, 0, -1, 0, 0},
		[]complex128{0, 0, 0, 0, 0, 0, 1, 0},
		[]complex128{0, 0, 0, 0, 0, 0, 0, -1},
	)

	actual := gate.CZ(3, 0, 2)
	if !actual.Equals(expected) {
		t.Error(actual)
	}
}

func TestCNOT(t *testing.T) {
	expected := gate.New(
		[]complex128{1, 0, 0, 0, 0, 0, 0, 0},
		[]complex128{0, 1, 0, 0, 0, 0, 0, 0},
		[]complex128{0, 0, 1, 0, 0, 0, 0, 0},
		[]complex128{0, 0, 0, 1, 0, 0, 0, 0},
		[]complex128{0, 0, 0, 0, 0, 1, 0, 0},
		[]complex128{0, 0, 0, 0, 1, 0, 0, 0},
		[]complex128{0, 0, 0, 0, 0, 0, 0, 1},
		[]complex128{0, 0, 0, 0, 0, 0, 1, 0},
	)

	actual := gate.CNOT(3, 0, 2)
	if !actual.Equals(expected) {
		t.Error(actual)
	}
}

func TestControlledNot(t *testing.T) {
	g0 := gate.I().Add(gate.Z()).TensorProduct(gate.I())
	g1 := gate.I().Sub(gate.Z()).TensorProduct(gate.X())
	expected := g0.Add(g1).Mul(0.5)

	actual := gate.ControlledNot(2, []int{0}, 1)
	if !actual.Equals(expected) {
		t.Error(actual)
	}
}

func TestToffoli(t *testing.T) {
	g := gate.Zero(13)
	g[0] = gate.I(2).TensorProduct(gate.H())
	g[1] = gate.I(1).TensorProduct(gate.CNOT(2, 0, 1))
	g[2] = gate.I(2).TensorProduct(gate.T().Dagger())
	g[3] = gate.CNOT(3, 0, 2)
	g[4] = gate.I(2).TensorProduct(gate.T())
	g[5] = gate.I(1).TensorProduct(gate.CNOT(2, 0, 1))
	g[6] = gate.I(2).TensorProduct(gate.T().Dagger())
	g[7] = gate.CNOT(3, 0, 2)
	g[8] = gate.I(1).TensorProduct(gate.T().Dagger()).TensorProduct(gate.T())
	g[9] = gate.CNOT(2, 0, 1).TensorProduct(gate.H())
	g[10] = gate.I(1).TensorProduct(gate.T().Dagger()).TensorProduct(gate.I())
	g[11] = gate.CNOT(2, 0, 1).TensorProduct(gate.I())
	g[12] = gate.T().TensorProduct(gate.S()).TensorProduct(gate.I())

	expected := gate.I(3)
	for _, gate := range g {
		expected = expected.Apply(gate)
	}

	actual := gate.Toffoli(3, 0, 1, 2)
	if !actual.Equals(expected, 1e-13) {
		t.Error(actual)
	}
}

func TestFredkin(t *testing.T) {
	m := make(matrix.Matrix, 8)
	m[0] = []complex128{1, 0, 0, 0, 0, 0, 0, 0}
	m[1] = []complex128{0, 1, 0, 0, 0, 0, 0, 0}
	m[2] = []complex128{0, 0, 1, 0, 0, 0, 0, 0}
	m[3] = []complex128{0, 0, 0, 1, 0, 0, 0, 0}
	m[4] = []complex128{0, 0, 0, 0, 1, 0, 0, 0}
	m[5] = []complex128{0, 0, 0, 0, 0, 0, 1, 0}
	m[6] = []complex128{0, 0, 0, 0, 0, 1, 0, 0}
	m[7] = []complex128{0, 0, 0, 0, 0, 0, 0, 1}

	actual := gate.Fredkin(3, 0, 1, 2)
	if !actual.Equals(m, 1e-13) {
		t.Error(actual)
	}
}
