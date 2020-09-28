package gate

import (
	"testing"

	"github.com/itsubaki/q/pkg/math/matrix"
	"github.com/itsubaki/q/pkg/math/vector"
)

func TestCModExp2(t *testing.T) {
	v0 := vector.TensorProduct(vector.TensorProductN(vector.New(1, 0), 6), vector.New(0, 1))
	g0 := CNOT(7, 2, 4).Apply(CNOT(7, 2, 5))
	g1 := CNOT(7, 3, 5).Apply(CCNOT(7, 1, 5, 3)).Apply(CNOT(7, 3, 5))
	g2 := CNOT(7, 4, 6).Apply(CCNOT(7, 1, 6, 4)).Apply(CNOT(7, 4, 6))
	g3 := I(7)

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
		a := CModExp2(c.b, c.a, c.j, c.N, c.c, c.t)
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
		{U(1, 2, 3, 4), I()},
		{X(2), I(2)},
		{CNOT(2, 0, 1), I(2)},
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
		{H()},
		{X()},
		{Y()},
		{Z()},
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
		{I()},
		{H()},
		{X()},
		{Y()},
		{Z()},
		{U(1, 2, 3, 4)},
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
		{I(), complex(2, 0)},
		{H(), complex(0, 0)},
		{X(), complex(0, 0)},
		{Y(), complex(0, 0)},
		{Z(), complex(0, 0)},
	}

	for _, c := range cases {
		tr := c.m.Trace()
		if tr != c.t {
			t.Error(tr)
		}
	}
}

func TestCZ(t *testing.T) {
	expected := New(
		[]complex128{1, 0, 0, 0, 0, 0, 0, 0},
		[]complex128{0, 1, 0, 0, 0, 0, 0, 0},
		[]complex128{0, 0, 1, 0, 0, 0, 0, 0},
		[]complex128{0, 0, 0, 1, 0, 0, 0, 0},
		[]complex128{0, 0, 0, 0, 1, 0, 0, 0},
		[]complex128{0, 0, 0, 0, 0, -1, 0, 0},
		[]complex128{0, 0, 0, 0, 0, 0, 1, 0},
		[]complex128{0, 0, 0, 0, 0, 0, 0, -1},
	)

	actual := CZ(3, 0, 2)
	if !actual.Equals(expected) {
		t.Error(actual)
	}
}

func TestCNOT(t *testing.T) {
	expected := New(
		[]complex128{1, 0, 0, 0, 0, 0, 0, 0},
		[]complex128{0, 1, 0, 0, 0, 0, 0, 0},
		[]complex128{0, 0, 1, 0, 0, 0, 0, 0},
		[]complex128{0, 0, 0, 1, 0, 0, 0, 0},
		[]complex128{0, 0, 0, 0, 0, 1, 0, 0},
		[]complex128{0, 0, 0, 0, 1, 0, 0, 0},
		[]complex128{0, 0, 0, 0, 0, 0, 0, 1},
		[]complex128{0, 0, 0, 0, 0, 0, 1, 0},
	)

	actual := CNOT(3, 0, 2)
	if !actual.Equals(expected) {
		t.Error(actual)
	}
}

func TestControlledNot(t *testing.T) {
	g0 := I().Add(Z()).TensorProduct(I())
	g1 := I().Sub(Z()).TensorProduct(X())
	expected := g0.Add(g1).Mul(0.5)

	actual := ControlledNot(2, []int{0}, 1)
	if !actual.Equals(expected) {
		t.Error(actual)
	}
}

func TestToffoli(t *testing.T) {
	g := Zero(13)
	g[0] = I(2).TensorProduct(H())
	g[1] = I(1).TensorProduct(CNOT(2, 0, 1))
	g[2] = I(2).TensorProduct(T().Dagger())
	g[3] = CNOT(3, 0, 2)
	g[4] = I(2).TensorProduct(T())
	g[5] = I(1).TensorProduct(CNOT(2, 0, 1))
	g[6] = I(2).TensorProduct(T().Dagger())
	g[7] = CNOT(3, 0, 2)
	g[8] = I(1).TensorProduct(T().Dagger()).TensorProduct(T())
	g[9] = CNOT(2, 0, 1).TensorProduct(H())
	g[10] = I(1).TensorProduct(T().Dagger()).TensorProduct(I())
	g[11] = CNOT(2, 0, 1).TensorProduct(I())
	g[12] = T().TensorProduct(S()).TensorProduct(I())

	expected := I(3)
	for _, gate := range g {
		expected = expected.Apply(gate)
	}

	actual := Toffoli(3, 0, 1, 2)
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

	actual := Fredkin(3, 0, 1, 2)
	if !actual.Equals(m, 1e-13) {
		t.Error(actual)
	}
}
