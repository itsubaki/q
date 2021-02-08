package density_test

import (
	"fmt"
	"math/cmplx"
	"testing"

	"github.com/itsubaki/q/pkg/math/matrix"
	"github.com/itsubaki/q/pkg/quantum/density"
	"github.com/itsubaki/q/pkg/quantum/gate"
	"github.com/itsubaki/q/pkg/quantum/qubit"
)

func ExampleMatrix_PartialTrace_bell() {
	qc := matrix.Apply(
		matrix.TensorProduct(gate.H(), gate.I()),
		gate.CNOT(2, 0, 1),
	)
	q := qubit.Zero(2).Apply(qc)
	rho := density.New().Add(1.0, q)
	for _, p := range rho.Internal() {
		fmt.Printf("%.1f\n", p)
	}
	fmt.Println()

	pt := rho.PartialTrace(0)
	for _, p := range pt.Internal() {
		fmt.Printf("%.1f\n", p)
	}

	// Output:
	// [(0.5+0.0i) (0.0+0.0i) (0.0+0.0i) (0.5+0.0i)]
	// [(0.0+0.0i) (0.0+0.0i) (0.0+0.0i) (0.0+0.0i)]
	// [(0.0+0.0i) (0.0+0.0i) (0.0+0.0i) (0.0+0.0i)]
	// [(0.5+0.0i) (0.0+0.0i) (0.0+0.0i) (0.5+0.0i)]

	// [(0.5+0.0i) (0.0+0.0i)]
	// [(0.0+0.0i) (0.5+0.0i)]
}

func ExampleMatrix_PartialTrace() {
	q := qubit.Zero(2).Apply(gate.H(2))
	rho := density.New().Add(1.0, q)
	for _, p := range rho.Internal() {
		fmt.Printf("%.2f\n", p)
	}
	fmt.Println()

	pt := rho.PartialTrace(1)
	for _, p := range pt.Internal() {
		fmt.Printf("%.2f\n", p)
	}
	// Output:
}

func TestDensityMatrix(t *testing.T) {
	cases := []struct {
		p         []float64
		q         []*qubit.Qubit
		tr, str   complex128
		expectedM matrix.Matrix
		expectedV complex128
	}{
		{
			[]float64{0.1, 0.9},
			[]*qubit.Qubit{qubit.Zero(), qubit.One()},
			complex(1, 0), complex(0.82, 0),
			gate.X(), 0.0,
		},
		{
			[]float64{0.1, 0.9},
			[]*qubit.Qubit{qubit.Zero(), qubit.Zero().Apply(gate.H())},
			complex(1, 0), complex(0.91, 0),
			gate.X(), 0.9,
		},
	}

	for _, c := range cases {
		rho := density.New()
		for i := range c.p {
			rho.Add(c.p[i], c.q[i])
		}

		if cmplx.Abs(rho.Trace()-c.tr) > 1e-13 {
			t.Errorf("trace=%v", rho.Trace())
		}

		if cmplx.Abs(rho.Squared().Trace()-c.str) > 1e-13 {
			t.Errorf("strace%v", rho.Squared().Trace())
		}

		if cmplx.Abs(rho.ExpectedValue(c.expectedM)-c.expectedV) > 1e-13 {
			t.Errorf("expected value=%v", rho.ExpectedValue(c.expectedM))
		}
	}
}
