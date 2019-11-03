package density

import (
	"math/cmplx"
	"testing"

	"github.com/itsubaki/q/pkg/quantum/qubit"
)

func TestDensityMatrix(t *testing.T) {
	q0 := qubit.Zero()
	q1 := qubit.One()

	m0 := q0.OuterProduct(q0).Mul(complex(0.5, 0))
	m1 := q1.OuterProduct(q1).Mul(complex(0.5, 0))

	m := m0.Add(m1)
	if m[0][0] != complex(0.5, 0) {
		t.Error(m)
	}
	if m[1][1] != complex(0.5, 0) {
		t.Error(m)
	}

	if cmplx.Abs(m.Trace()-complex(1, 0)) > 1e-13 {
		t.Error(m)
	}

	e0 := q0.OuterProduct(q0)
	t0 := e0.Dagger().Apply(m).Apply(e0).Trace()
	if t0 != complex(0.5, 0) {
		t.Error(t0)
	}

	e1 := q1.OuterProduct(q1)
	t1 := e1.Dagger().Apply(m).Apply(e1).Trace()
	if t1 != complex(0.5, 0) {
		t.Error(t1)
	}

	i := e0.Dagger().Apply(e0).Add(e1.Dagger().Apply(e1))
	if i[0][0] != complex(1, 0) || i[1][1] != complex(1, 0) {
		t.Error(i)
	}
	if i[0][1] != complex(0, 0) || i[1][0] != complex(0, 0) {
		t.Error(i)
	}
}
