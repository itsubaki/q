package density

import (
	"fmt"
	"math/cmplx"
	"testing"

	"github.com/itsubaki/q/pkg/quantum/gate"
	"github.com/itsubaki/q/pkg/quantum/qubit"
)

func TestPartialTrace(t *testing.T) {
	bell := qubit.Zero(2).Apply(gate.H().TensorProduct(gate.I())).Apply(gate.CNOT(2, 0, 1))
	fmt.Println(bell)

	rho := New().Add(1, bell)
	fmt.Println(rho)

	pt := rho.PartialTrace(0, qubit.Zero(), qubit.One())
	fmt.Println(pt)
}

func TestDensityMatrix(t *testing.T) {
	p0, p1 := 0.1, 0.9
	q0, q1 := qubit.Zero(), qubit.One()
	rho := New().Add(p0, q0).Add(p1, q1)

	if cmplx.Abs(rho.Trace()-complex(1, 0)) > 1e-13 {
		t.Error(rho)
	}

	if rho.Measure(q0) != complex(p0, 0) {
		t.Error(rho)
	}

	if rho.Measure(q1) != complex(p1, 0) {
		t.Error(rho)
	}

	if cmplx.Abs(rho.ExpectedValue(gate.X())) > 1e-13 {
		t.Error(rho)
	}

	xrho := rho.Apply(gate.X())
	if xrho.Measure(q0) != complex(p1, 0) {
		t.Error(xrho)
	}

	if xrho.Measure(q1) != complex(p0, 0) {
		t.Error(xrho)
	}
}

func TestDensityMatrix2(t *testing.T) {
	p0, p1 := 0.1, 0.9
	q0, q1 := qubit.Zero(), qubit.Zero().Apply(gate.H())
	rho := New().Add(p0, q0).Add(p1, q1)

	if cmplx.Abs(rho.Trace()-complex(1, 0)) > 1e-13 {
		t.Error(rho)
	}

	e0 := rho.ExpectedValue(gate.X())
	if cmplx.Abs(e0-complex(p1, 0)) > 1e-13 {
		t.Error(e0)
	}

	trrho2 := rho.Squared().Trace() // < 1
	if cmplx.Abs(trrho2-complex(1, 0)) < 1e-13 {
		t.Error(trrho2)
	}
}

func TestDensityMatrixPureState(t *testing.T) {
	rho := New().Add(1.0, qubit.Zero())

	trrho2 := rho.Squared().Trace() // -> 1
	if cmplx.Abs(trrho2-complex(1, 0)) > 1e-13 {
		t.Error(trrho2)
	}
}

func TestDepolarizing(t *testing.T) {
	p0, p1 := 0.1, 0.9
	q0, q1 := qubit.Zero(), qubit.One()
	rho := New().Add(p0, q0).Add(p1, q1)

	rho.Depolarizing(0)
	if cmplx.Abs(rho.Trace()-complex(1, 0)) > 1e-13 {
		t.Error(rho)
	}

	if rho.Measure(q0) != complex(p0, 0) {
		t.Error(rho)
	}

	if rho.Measure(q1) != complex(p1, 0) {
		t.Error(rho)
	}

	rho.Depolarizing(1)
	if cmplx.Abs(rho.Trace()-complex(1, 0)) > 1e-13 {
		t.Error(rho)
	}

	if rho.Measure(q0) != complex(0.5, 0) {
		t.Error(rho)
	}

	if rho.Measure(q1) != complex(0.5, 0) {
		t.Error(rho)
	}
}
