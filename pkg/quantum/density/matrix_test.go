package density

import (
	"math/cmplx"
	"testing"

	"github.com/itsubaki/q/pkg/quantum/qubit"
)

func TestDensityMatrix(t *testing.T) {
	m := New().Add(0.5, qubit.Zero()).Add(0.5, qubit.One())

	if m.At(0, 0) != complex(0.5, 0) {
		t.Error(m)
	}
	if m.At(1, 1) != complex(0.5, 0) {
		t.Error(m)
	}

	if cmplx.Abs(m.Trace()-complex(1, 0)) > 1e-13 {
		t.Error(m)
	}

	e0 := qubit.Zero().OuterProduct(qubit.Zero())
	if m.Apply(e0).At(0, 0) != complex(0.5, 0) {
		t.Error(m)
	}

	e1 := qubit.One().OuterProduct(qubit.One())
	if m.Apply(e1).At(1, 1) != complex(0.5, 0) {
		t.Error(m)
	}
}
