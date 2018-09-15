package vector

import (
	"testing"

	"github.com/itsubaki/q/circuit/gate"
	"github.com/itsubaki/q/math/matrix"
)

func TestVector(t *testing.T) {
	v0 := New(1, 1)
	v1 := New(1, -1)

	if v0.InnerProduct(v1) != complex(0, 0) {
		t.Error(v0.InnerProduct(v1))
	}

	if !v0.IsOrthogonal(v1) {
		t.Error(v0.InnerProduct(v1))
	}

	v3 := New(1, 0)
	if v3.InnerProduct(v3) != complex(1, 0) {
		t.Error(v3.InnerProduct(v3))
	}

	if v3.IsOrthogonal(v3) {
		t.Error(v3.InnerProduct(v3))
	}

	if !v3.IsUnit() {
		t.Error(v3.IsUnit())
	}

	if v3.Norm() != complex(1, 0) {
		t.Error(v3.Norm())
	}

	v4 := New(0, 1)
	if v3.InnerProduct(v4) != complex(0, 0) {
		t.Error(v3.InnerProduct(v4))
	}
	if !v3.IsOrthogonal(v4) {
		t.Error(v3.InnerProduct(v4))
	}

}

func TestTensorProduct(t *testing.T) {
	v := New(1, 0)

	v4 := TensorProduct(v, v)
	x4 := matrix.TensorProduct(gate.X(), gate.X())

	xv4 := v4.Apply(x4)
	expected := TensorProduct(New(0, 1), New(0, 1))
	if !xv4.Equals(expected) {
		t.Error(xv4)
	}

	v16 := TensorProduct(v4, v4)
	x16 := matrix.TensorProduct(x4, x4)
	xv16 := v16.Apply(x16)

	expected16 := TensorProduct(New(0, 0, 0, 1), New(0, 0, 0, 1))
	if !xv16.Equals(expected16) {
		t.Error(xv16)
	}

}
