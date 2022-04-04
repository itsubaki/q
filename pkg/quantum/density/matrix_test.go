package density_test

import (
	"fmt"
	"math/cmplx"
	"testing"

	"github.com/itsubaki/q/pkg/math/epsilon"
	"github.com/itsubaki/q/pkg/math/matrix"
	"github.com/itsubaki/q/pkg/quantum/density"
	"github.com/itsubaki/q/pkg/quantum/gate"
	"github.com/itsubaki/q/pkg/quantum/qubit"
)

func ExampleMatrix_ExpectedValue() {
	p0, q0 := 0.1, qubit.Zero()
	p1, q1 := 0.9, qubit.Zero().Apply(gate.H())
	rho := density.New().Add(p0, q0).Add(p1, q1)

	fmt.Printf("%.2v\n", rho.ExpectedValue(gate.X()))
	fmt.Printf("%.2v\n", rho.ExpectedValue(gate.Y()))
	fmt.Printf("%.2v\n", rho.ExpectedValue(gate.Z()))

	// Output:
	// (0.9+0i)
	// (0+0i)
	// (0.1+0i)
}

func ExampleMatrix_Measure() {
	p0, q0 := 0.1, qubit.Zero()
	p1, q1 := 0.9, qubit.One()
	rho := density.New().Add(p0, q0).Add(p1, q1)

	fmt.Println(rho.Measure(qubit.Zero()))
	fmt.Println(rho.Measure(qubit.One()))

	// Output:
	// (0.1+0i)
	// (0.9+0i)
}

func ExampleMatrix_Trace() {
	pure := density.New().Add(1.0, qubit.Zero())
	mix := density.New().Add(0.1, qubit.Zero()).Add(0.9, qubit.One())

	fmt.Printf("%.2f\n", pure.Squared().Trace())
	fmt.Printf("%.2f\n", mix.Squared().Trace())

	// Output:
	// (1.00+0.00i)
	// (0.82+0.00i)
}

func ExampleMatrix_Depolarizing() {
	rho := density.New().Add(1, qubit.Zero())
	fmt.Println(rho.Measure(qubit.Zero()))
	fmt.Println(rho.Measure(qubit.One()))

	dep := rho.Depolarizing(1)
	fmt.Println(dep.Measure(qubit.Zero()))
	fmt.Println(dep.Measure(qubit.One()))

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

func ExampleMatrix_PartialTrace() {
	bell := matrix.Apply(
		matrix.TensorProduct(gate.H(), gate.I()),
		gate.CNOT(2, 0, 1),
	)

	q := qubit.Zero(2).Apply(bell)
	rho := density.New().Add(1.0, q)
	for _, r := range rho.Raw() {
		fmt.Printf("%.1f\n", r)
	}
	fmt.Println()

	p0 := rho.PartialTrace(0)
	for _, r := range p0.Raw() {
		fmt.Printf("%.1f\n", r)
	}
	fmt.Println()

	p1 := rho.PartialTrace(1)
	for _, r := range p1.Raw() {
		fmt.Printf("%.1f\n", r)
	}

	// Output:
	// [(0.5+0.0i) (0.0+0.0i) (0.0+0.0i) (0.5+0.0i)]
	// [(0.0+0.0i) (0.0+0.0i) (0.0+0.0i) (0.0+0.0i)]
	// [(0.0+0.0i) (0.0+0.0i) (0.0+0.0i) (0.0+0.0i)]
	// [(0.5+0.0i) (0.0+0.0i) (0.0+0.0i) (0.5+0.0i)]
	//
	// [(0.5+0.0i) (0.0+0.0i)]
	// [(0.0+0.0i) (0.5+0.0i)]
	//
	// [(0.5+0.0i) (0.0+0.0i)]
	// [(0.0+0.0i) (0.5+0.0i)]
}

func TestPartialTrace(t *testing.T) {
	type Case struct {
		index int
		want  [][]complex128
	}

	cases := []struct {
		rho *density.Matrix
		cs  []Case
		eps float64
	}{
		{
			density.New().Add(1.0, qubit.Zero(2)),
			[]Case{
				{0, [][]complex128{{1, 0}, {0, 0}}},
				{1, [][]complex128{{1, 0}, {0, 0}}},
			},
			epsilon.E13(),
		},
		{
			density.New().Add(1.0, qubit.One(2)),
			[]Case{
				{0, [][]complex128{{0, 0}, {0, 1}}},
				{1, [][]complex128{{0, 0}, {0, 1}}},
			},
			epsilon.E13(),
		},
		{
			density.New().Add(1, qubit.Zero(2).Apply(gate.H(2))),
			[]Case{
				{0, [][]complex128{{0.5, 0.5}, {0.5, 0.5}}},
				{1, [][]complex128{{0.5, 0.5}, {0.5, 0.5}}},
			},
			epsilon.E13(),
		},
		{
			density.New().
				Add(0.5, qubit.Zero(2)).
				Add(0.5, qubit.One(2)),
			[]Case{
				{0, [][]complex128{{0.5, 0}, {0, 0.5}}},
				{1, [][]complex128{{0.5, 0}, {0, 0.5}}},
			}, epsilon.E13(),
		},
		{
			density.New().
				Add(0.5, qubit.Zero(2).Apply(gate.H(2))).
				Add(0.5, qubit.One(2)),
			[]Case{
				{0, [][]complex128{{0.25, 0.25}, {0.25, 0.75}}},
				{1, [][]complex128{{0.25, 0.25}, {0.25, 0.75}}},
			},
			epsilon.E13(),
		},
		{
			density.New().
				Add(0.75, qubit.Zero(2).Apply(gate.H(2))).
				Add(0.25, qubit.One(2).Apply(gate.H(2))),
			[]Case{
				{0, [][]complex128{{0.5, 0.25}, {0.25, 0.5}}},
				{1, [][]complex128{{0.5, 0.25}, {0.25, 0.5}}},
			},
			epsilon.E13(),
		},
		{
			density.New().
				Add(0.25, qubit.Zero(2).Apply(gate.H(2))).
				Add(0.75, qubit.One(2).Apply(gate.H(2))),
			[]Case{
				{0, [][]complex128{{0.5, -0.25}, {-0.25, 0.5}}},
				{1, [][]complex128{{0.5, -0.25}, {-0.25, 0.5}}},
			},
			epsilon.E13(),
		},
	}

	for _, c := range cases {
		for _, cs := range c.cs {
			got := c.rho.PartialTrace(cs.index).Raw()
			if len(got) != len(cs.want) {
				t.Errorf("got=%v want=%v", got, cs.want)
			}

			for i := 0; i < len(cs.want); i++ {
				for j := 0; j < len(cs.want[0]); j++ {
					if cmplx.Abs(got[i][j]-cs.want[i][j]) > c.eps {
						t.Errorf("%v:%v, got=%v want=%v", i, j, got[i][j], cs.want[i][j])
					}
				}
			}
		}
	}
}

func TestExpectedValue(t *testing.T) {
	cases := []struct {
		p       []float64
		q       []*qubit.Qubit
		tr, str complex128
		m       matrix.Matrix
		v       complex128
		eps     float64
	}{
		{
			[]float64{0.1, 0.9},
			[]*qubit.Qubit{qubit.Zero(), qubit.One()},
			1, 0.82,
			gate.X(), 0.0,
			epsilon.E13(),
		},
		{
			[]float64{0.1, 0.9},
			[]*qubit.Qubit{qubit.Zero(), qubit.Zero().Apply(gate.H())},
			1, 0.91,
			gate.X(), 0.9,
			epsilon.E13(),
		},
	}

	for _, c := range cases {
		rho := density.New()
		for i := range c.p {
			rho.Add(c.p[i], c.q[i])
		}

		if cmplx.Abs(rho.Trace()-c.tr) > c.eps {
			t.Errorf("trace=%v", rho.Trace())
		}

		if cmplx.Abs(rho.Squared().Trace()-c.str) > c.eps {
			t.Errorf("strace%v", rho.Squared().Trace())
		}

		if cmplx.Abs(rho.ExpectedValue(c.m)-c.v) > c.eps {
			t.Errorf("expected value=%v", rho.ExpectedValue(c.m))
		}
	}
}

func TestMeasure(t *testing.T) {
	cases := []struct {
		p    []float64
		q    []*qubit.Qubit
		m    *qubit.Qubit
		want complex128
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
		got := density.New()
		for i := range c.p {
			got.Add(c.p[i], c.q[i])
		}

		if got.Measure(c.m) != c.want {
			t.Fail()
		}
	}
}

func TestApply(t *testing.T) {
	cases := []struct {
		p    []float64
		q    []*qubit.Qubit
		g    matrix.Matrix
		m    *qubit.Qubit
		want complex128
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
		got := density.New()
		for i := range c.p {
			got.Add(c.p[i], c.q[i])
		}

		if got.Apply(c.g).Measure(c.m) != c.want {
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
