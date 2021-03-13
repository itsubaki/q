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

func ExampleMatrix_Depolarizing() {
	rho := density.New().Add(1, qubit.Zero())
	fmt.Println(rho.Measure(qubit.Zero()))
	fmt.Println(rho.Measure(qubit.One()))

	rho.Depolarizing(1)
	fmt.Println(rho.Measure(qubit.Zero()))
	fmt.Println(rho.Measure(qubit.One()))

	// Output:
	// (1+0i)
	// (0+0i)
	// (0.5+0i)
	// (0.5+0i)
}

func ExampleBitFlip() {
	m0, m1 := density.BitFlip(0.5)

	for _, r := range m0 {
		fmt.Println(r)
	}

	for _, r := range m1 {
		fmt.Println(r)
	}

	// Output:
	// [(0.7071067811865476+0i) (0+0i)]
	// [(0+0i) (0.7071067811865476+0i)]
	// [(0+0i) (0.7071067811865476+0i)]
	// [(0.7071067811865476+0i) (0+0i)]
}

func ExamplePhaseFlip() {
	m0, m1 := density.PhaseFlip(0.5)

	for _, r := range m0 {
		fmt.Println(r)
	}

	for _, r := range m1 {
		fmt.Println(r)
	}

	// Output:
	// [(0.7071067811865476+0i) (0+0i)]
	// [(0+0i) (0.7071067811865476+0i)]
	// [(0.7071067811865476+0i) (0+0i)]
	// [(0+0i) (-0.7071067811865476+0i)]
}

func ExampleBitPhaseFlip() {
	m0, m1 := density.BitPhaseFlip(0.5)

	for _, r := range m0 {
		fmt.Println(r)
	}

	for _, r := range m1 {
		fmt.Println(r)
	}

	// Output:
	// [(0.7071067811865476+0i) (0+0i)]
	// [(0+0i) (0.7071067811865476+0i)]
	// [(0+0i) (0-0.7071067811865476i)]
	// [(0+0.7071067811865476i) (0+0i)]
}

func TestPartialTrace(t *testing.T) {
	qc := matrix.Apply(
		matrix.TensorProduct(gate.H(), gate.I()),
		gate.CNOT(2, 0, 1),
	)
	q := qubit.Zero(2).Apply(qc)
	rho := density.New().Add(1.0, q)

	pt := rho.PartialTrace(0)
	fmt.Println(pt)
}

func TestExpectedValue(t *testing.T) {
	cases := []struct {
		p       []float64
		q       []*qubit.Qubit
		tr, str complex128
		m       matrix.Matrix
		v       complex128
	}{
		{
			[]float64{0.1, 0.9},
			[]*qubit.Qubit{qubit.Zero(), qubit.One()},
			1, 0.82,
			gate.X(), 0.0,
		},
		{
			[]float64{0.1, 0.9},
			[]*qubit.Qubit{qubit.Zero(), qubit.Zero().Apply(gate.H())},
			1, 0.91,
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

		if cmplx.Abs(rho.ExpectedValue(c.m)-c.v) > 1e-13 {
			t.Errorf("expected value=%v", rho.ExpectedValue(c.m))
		}
	}
}

func TestMeasure(t *testing.T) {
	cases := []struct {
		p []float64
		q []*qubit.Qubit
		m *qubit.Qubit
		v complex128
	}{
		{
			[]float64{1},
			[]*qubit.Qubit{qubit.Zero()},
			qubit.Zero(),
			1,
		},
		{
			[]float64{1},
			[]*qubit.Qubit{qubit.Zero()},
			qubit.One(),
			0,
		},
	}

	for _, c := range cases {
		rho := density.New()
		for i := range c.p {
			rho.Add(c.p[i], c.q[i])
		}

		if rho.Measure(c.m) != c.v {
			t.Fail()
		}
	}
}

func TestApply(t *testing.T) {
	cases := []struct {
		p []float64
		q []*qubit.Qubit
		g matrix.Matrix
		m *qubit.Qubit
		e complex128
	}{
		{
			[]float64{1},
			[]*qubit.Qubit{qubit.Zero()},
			gate.X(),
			qubit.Zero(),
			0,
		},
		{
			[]float64{1},
			[]*qubit.Qubit{qubit.Zero()},
			gate.X(),
			qubit.One(),
			1,
		},
	}

	for _, c := range cases {
		rho := density.New()
		for i := range c.p {
			rho.Add(c.p[i], c.q[i])
		}

		if rho.Apply(c.g).Measure(c.m) != c.e {
			t.Fail()
		}
	}
}

func TestAddPanicDimenstion(t *testing.T) {
	p0, q0 := 0.1, qubit.Zero()
	p1, q1 := 0.9, qubit.New(1, 0, 0, 0)

	defer func() {
		if err := recover(); err != "invalid dimension. m=2 n=4" {
			t.Fail()
		}
	}()

	density.New().Add(p0, q0).Add(p1, q1)
	t.Fail()
}

func TestAddPanicProbabilityLess(t *testing.T) {
	cases := []struct {
		p float64
	}{
		{-1},
	}

	for _, c := range cases {
		defer func() {
			msg := fmt.Sprintf("p must be 0 <= p =< 1. p=%v", c.p)
			if err := recover(); err != msg {
				t.Fail()
			}
		}()

		density.New().Add(c.p, qubit.Zero())
		t.Fail()
	}
}

func TestAddPanicProbabilityLarge(t *testing.T) {
	cases := []struct {
		p float64
	}{
		{1.1},
	}

	for _, c := range cases {
		defer func() {
			msg := fmt.Sprintf("p must be 0 <= p =< 1. p=%v", c.p)
			if err := recover(); err != msg {
				t.Fail()
			}
		}()

		density.New().Add(c.p, qubit.Zero())
		t.Fail()
	}
}

func TestDepolarizingPanicLess(t *testing.T) {
	cases := []struct {
		p float64
	}{
		{-1},
	}

	for _, c := range cases {
		defer func() {
			msg := fmt.Sprintf("p must be 0 <= p =< 1. p=%v", c.p)
			if err := recover(); err != msg {
				t.Fail()
			}
		}()

		rho := density.New().Add(1, qubit.Zero())
		rho.Depolarizing(c.p)
		t.Fail()
	}
}

func TestDepolarizingPanicLarge(t *testing.T) {
	cases := []struct {
		p float64
	}{
		{1.1},
	}

	for _, c := range cases {
		defer func() {
			msg := fmt.Sprintf("p must be 0 <= p =< 1. p=%v", c.p)
			if err := recover(); err != msg {
				t.Fail()
			}
		}()

		rho := density.New().Add(1, qubit.Zero())
		rho.Depolarizing(c.p)
		t.Fail()
	}
}

func TestBitFlipPanicLess(t *testing.T) {
	cases := []struct {
		p float64
	}{
		{-1},
	}

	for _, c := range cases {
		defer func() {
			msg := fmt.Sprintf("p must be 0 <= p =< 1. p=%v", c.p)
			if err := recover(); err != msg {
				t.Fail()
			}
		}()

		density.BitFlip(c.p)
		t.Fail()
	}
}

func TestBitFlipPanicLarge(t *testing.T) {
	cases := []struct {
		p float64
	}{
		{1.1},
	}

	for _, c := range cases {
		defer func() {
			msg := fmt.Sprintf("p must be 0 <= p =< 1. p=%v", c.p)
			if err := recover(); err != msg {
				t.Fail()
			}
		}()

		density.BitFlip(c.p)
		t.Fail()
	}
}
