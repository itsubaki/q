package density

import (
	"math/cmplx"
	"testing"

	"github.com/itsubaki/q/pkg/quantum/gate"
	"github.com/itsubaki/q/pkg/quantum/qubit"
)

func TestDensityMatrix(t *testing.T) {
	p0, p1 := 0.1, 0.9
	q0, q1 := qubit.Zero(), qubit.One()
	m := New().Add(p0, q0).Add(p1, q1)

	if cmplx.Abs(m.Trace()-complex(1, 0)) > 1e-13 {
		t.Error(m)
	}

	if m.Measure(q0) != complex(p0, 0) {
		t.Error(m)
	}

	if m.Measure(q1) != complex(p1, 0) {
		t.Error(m)
	}

	xm := m.Apply(gate.X())
	if xm.Measure(q0) != complex(p1, 0) {
		t.Error(m)
	}

	if xm.Measure(q1) != complex(p0, 0) {
		t.Error(m)
	}
}
