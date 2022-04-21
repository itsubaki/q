package density_test

import (
	"fmt"
	"math"
	"math/cmplx"
	"testing"

	"github.com/itsubaki/q/pkg/math/epsilon"
	"github.com/itsubaki/q/pkg/math/matrix"
	"github.com/itsubaki/q/pkg/quantum/density"
	"github.com/itsubaki/q/pkg/quantum/gate"
	"github.com/itsubaki/q/pkg/quantum/qubit"
)

func ExampleMatrix_ExpectedValue() {
	rho, _ := density.New(
		[]float64{0.1, 0.9},
		[]*qubit.Qubit{
			qubit.Zero(),
			qubit.Zero().Apply(gate.H()),
		},
	)

	fmt.Printf("X: %.2v\n", rho.ExpectedValue(gate.X()))
	fmt.Printf("Y: %.2v\n", rho.ExpectedValue(gate.Y()))
	fmt.Printf("Z: %.2v\n", rho.ExpectedValue(gate.Z()))

	// Output:
	// X: 0.9
	// Y: 0
	// Z: 0.1
}

func ExampleMatrix_Measure() {
	rho, _ := density.New(
		[]float64{0.1, 0.9},
		[]*qubit.Qubit{
			qubit.Zero(),
			qubit.One(),
		},
	)

	fmt.Printf("0: %.2v\n", rho.Measure(qubit.Zero()))
	fmt.Printf("1: %.2v\n", rho.Measure(qubit.One()))

	// Output:
	// 0: 0.1
	// 1: 0.9
}

func ExampleMatrix_Trace() {
	pure, _ := density.New([]float64{1.0}, []*qubit.Qubit{qubit.Zero()})
	mixed, _ := density.New([]float64{0.1, 0.9}, []*qubit.Qubit{qubit.Zero(), qubit.One()})

	fmt.Printf("pure:  %.2f\n", pure.Trace())
	fmt.Printf("mixed: %.2f\n", mixed.Trace())

	// Output:
	// pure:  1.00
	// mixed: 1.00
}

func ExampleMatrix_SquareTrace() {
	pure, _ := density.New([]float64{1.0}, []*qubit.Qubit{qubit.Zero()})
	mixed, _ := density.New([]float64{0.1, 0.9}, []*qubit.Qubit{qubit.Zero(), qubit.One()})

	fmt.Printf("pure:  %.2f\n", pure.SquareTrace())
	fmt.Printf("mixed: %.2f\n", mixed.SquareTrace())

	// Output:
	// pure:  1.00
	// mixed: 0.82
}

func ExampleMatrix_PartialTrace() {
	rho, _ := density.New(
		[]float64{0.5, 0.5},
		[]*qubit.Qubit{
			qubit.Zero(2).Apply(gate.QFT(2)),
			qubit.One(2).Apply(gate.QFT(2)),
		},
	)

	for _, r := range rho.Raw() {
		fmt.Printf("%.4f\n", r)
	}
	fmt.Printf("trace: %.2v, square_trace: %.2v\n\n", rho.Trace(), rho.SquareTrace())

	p0 := rho.PartialTrace(0)
	for _, r := range p0.Raw() {
		fmt.Printf("%.4f\n", r)
	}
	fmt.Printf("trace: %.2v, square_trace: %.2v\n\n", p0.Trace(), p0.SquareTrace())

	p1 := rho.PartialTrace(1)
	for _, r := range p1.Raw() {
		fmt.Printf("%.4f\n", r)
	}
	fmt.Printf("trace: %.2v, square_trace: %.2v\n\n", p1.Trace(), p1.SquareTrace())

	// Output:
	// [(0.2500+0.0000i) (0.0000+0.0000i) (0.1250+0.1250i) (0.1250-0.1250i)]
	// [(0.0000+0.0000i) (0.2500+0.0000i) (0.1250-0.1250i) (0.1250+0.1250i)]
	// [(0.1250-0.1250i) (0.1250+0.1250i) (0.2500+0.0000i) (0.0000+0.0000i)]
	// [(0.1250+0.1250i) (0.1250-0.1250i) (0.0000+0.0000i) (0.2500+0.0000i)]
	// trace: 1, square_trace: 0.5
	//
	// [(0.5000+0.0000i) (0.0000+0.0000i)]
	// [(0.0000+0.0000i) (0.5000+0.0000i)]
	// trace: 1, square_trace: 0.5
	//
	// [(0.5000+0.0000i) (0.2500+0.2500i)]
	// [(0.2500-0.2500i) (0.5000+0.0000i)]
	// trace: 1, square_trace: 0.75
}

func ExampleMatrix_PartialTrace_x8() {
	rho, _ := density.New(
		[]float64{0.5, 0.5},
		[]*qubit.Qubit{
			qubit.Zero(3).Apply(gate.QFT(3)),
			qubit.One(3).Apply(gate.QFT(3)),
		},
	)

	for _, r := range rho.Raw() {
		fmt.Printf("%.4f\n", r)
	}
	fmt.Printf("trace: %.2v, square_trace: %.2v\n\n", rho.Trace(), rho.SquareTrace())

	p0 := rho.PartialTrace(0)
	for _, r := range p0.Raw() {
		fmt.Printf("%.4f\n", r)
	}
	fmt.Printf("trace: %.2v, square_trace: %.2v\n\n", p0.Trace(), p0.SquareTrace())

	p1 := rho.PartialTrace(1)
	for _, r := range p1.Raw() {
		fmt.Printf("%.4f\n", r)
	}
	fmt.Printf("trace: %.2v, square_trace: %.2v\n\n", p1.Trace(), p1.SquareTrace())

	p2 := rho.PartialTrace(2)
	for _, r := range p2.Raw() {
		fmt.Printf("%.4f\n", r)
	}
	fmt.Printf("trace: %.2v, square_trace: %.2v\n\n", p2.Trace(), p2.SquareTrace())

	// Output:
	// [(0.1250+0.0000i) (0.0000+0.0000i) (0.0625+0.0625i) (0.0625-0.0625i) (0.1067+0.0442i) (0.0183-0.0442i) (0.0183+0.0442i) (0.1067-0.0442i)]
	// [(0.0000+0.0000i) (0.1250+0.0000i) (0.0625-0.0625i) (0.0625+0.0625i) (0.0183-0.0442i) (0.1067+0.0442i) (0.1067-0.0442i) (0.0183+0.0442i)]
	// [(0.0625-0.0625i) (0.0625+0.0625i) (0.1250+0.0000i) (0.0000+0.0000i) (0.1067-0.0442i) (0.0183+0.0442i) (0.1067+0.0442i) (0.0183-0.0442i)]
	// [(0.0625+0.0625i) (0.0625-0.0625i) (0.0000+0.0000i) (0.1250+0.0000i) (0.0183+0.0442i) (0.1067-0.0442i) (0.0183-0.0442i) (0.1067+0.0442i)]
	// [(0.1067-0.0442i) (0.0183+0.0442i) (0.1067+0.0442i) (0.0183-0.0442i) (0.1250+0.0000i) (0.0000+0.0000i) (0.0625+0.0625i) (0.0625-0.0625i)]
	// [(0.0183+0.0442i) (0.1067-0.0442i) (0.0183-0.0442i) (0.1067+0.0442i) (0.0000+0.0000i) (0.1250+0.0000i) (0.0625-0.0625i) (0.0625+0.0625i)]
	// [(0.0183-0.0442i) (0.1067+0.0442i) (0.1067-0.0442i) (0.0183+0.0442i) (0.0625-0.0625i) (0.0625+0.0625i) (0.1250+0.0000i) (0.0000+0.0000i)]
	// [(0.1067+0.0442i) (0.0183-0.0442i) (0.0183+0.0442i) (0.1067-0.0442i) (0.0625+0.0625i) (0.0625-0.0625i) (0.0000+0.0000i) (0.1250+0.0000i)]
	// trace: 1, square_trace: 0.5
	//
	// [(0.2500+0.0000i) (0.0000+0.0000i) (0.1250+0.1250i) (0.1250-0.1250i)]
	// [(0.0000+0.0000i) (0.2500+0.0000i) (0.1250-0.1250i) (0.1250+0.1250i)]
	// [(0.1250-0.1250i) (0.1250+0.1250i) (0.2500+0.0000i) (0.0000+0.0000i)]
	// [(0.1250+0.1250i) (0.1250-0.1250i) (0.0000+0.0000i) (0.2500+0.0000i)]
	// trace: 1, square_trace: 0.5
	//
	// [(0.2500+0.0000i) (0.0000+0.0000i) (0.2134+0.0884i) (0.0366-0.0884i)]
	// [(0.0000+0.0000i) (0.2500+0.0000i) (0.0366-0.0884i) (0.2134+0.0884i)]
	// [(0.2134-0.0884i) (0.0366+0.0884i) (0.2500+0.0000i) (0.0000+0.0000i)]
	// [(0.0366+0.0884i) (0.2134-0.0884i) (0.0000+0.0000i) (0.2500+0.0000i)]
	// trace: 1, square_trace: 0.5
	//
	// [(0.2500+0.0000i) (0.1250+0.1250i) (0.2134+0.0884i) (0.0366+0.0884i)]
	// [(0.1250-0.1250i) (0.2500+0.0000i) (0.2134-0.0884i) (0.2134+0.0884i)]
	// [(0.2134-0.0884i) (0.2134+0.0884i) (0.2500+0.0000i) (0.1250+0.1250i)]
	// [(0.0366-0.0884i) (0.2134-0.0884i) (0.1250-0.1250i) (0.2500+0.0000i)]
	// trace: 1, square_trace: 0.71
}

func ExampleMatrix_PartialTrace_x16() {
	rho, _ := density.New(
		[]float64{1.0},
		[]*qubit.Qubit{
			qubit.Zero(4).
				Apply(matrix.TensorProduct(gate.H(2), gate.X(), gate.Z()).
					Apply(gate.CNOT(4, 1, 3)).
					Apply(gate.CNOT(4, 0, 2))),
		},
	)

	p00 := rho.PartialTrace(0).PartialTrace(0)
	fmt.Printf("trace: %.2f\n", p00.Trace())
	fmt.Printf("square_trace: %.2f\n", p00.SquareTrace())

	// Output:
	// trace: 1.00
	// square_trace: 0.25
}

func ExampleMatrix_Depolarizing() {
	rho, _ := density.New([]float64{1.0}, []*qubit.Qubit{qubit.Zero()})
	fmt.Printf("0: %.2f\n", rho.Measure(qubit.Zero()))
	fmt.Printf("1: %.2f\n", rho.Measure(qubit.One()))
	fmt.Println()

	dep, _ := rho.Depolarizing(1)
	fmt.Printf("0: %.2f\n", dep.Measure(qubit.Zero()))
	fmt.Printf("1: %.2f\n", dep.Measure(qubit.One()))

	// Output:
	// 0: 1.00
	// 1: 0.00
	//
	// 0: 0.50
	// 1: 0.50
}

func ExampleBitFlip() {
	m0, m1, _ := density.BitFlip(0.5)

	for _, r := range m0 {
		fmt.Println(r)
	}
	fmt.Println()

	for _, r := range m1 {
		fmt.Println(r)
	}

	// Output:
	// [(0.7071067811865476+0i) (0+0i)]
	// [(0+0i) (0.7071067811865476+0i)]
	//
	// [(0+0i) (0.7071067811865476+0i)]
	// [(0.7071067811865476+0i) (0+0i)]
}

func ExamplePhaseFlip() {
	m0, m1, _ := density.PhaseFlip(0.5)

	for _, r := range m0 {
		fmt.Println(r)
	}
	fmt.Println()

	for _, r := range m1 {
		fmt.Println(r)
	}

	// Output:
	// [(0.7071067811865476+0i) (0+0i)]
	// [(0+0i) (0.7071067811865476+0i)]
	//
	// [(0.7071067811865476+0i) (0+0i)]
	// [(0+0i) (-0.7071067811865476+0i)]
}

func ExampleBitPhaseFlip() {
	m0, m1, _ := density.BitPhaseFlip(0.5)

	for _, r := range m0 {
		fmt.Println(r)
	}
	fmt.Println()

	for _, r := range m1 {
		fmt.Println(r)
	}

	// Output:
	// [(0.7071067811865476+0i) (0+0i)]
	// [(0+0i) (0.7071067811865476+0i)]
	//
	// [(0+0i) (0-0.7071067811865476i)]
	// [(0+0.7071067811865476i) (0+0i)]
}

func TestExpectedValue(t *testing.T) {
	cases := []struct {
		p        []float64
		q        []*qubit.Qubit
		tr, sqtr float64
		m        matrix.Matrix
		v        float64
		eps      float64
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
		rho, _ := density.New(c.p, c.q)
		if math.Abs(rho.Trace()-c.tr) > c.eps {
			t.Errorf("trace=%v", rho.Trace())
		}

		if math.Abs(rho.SquareTrace()-c.sqtr) > c.eps {
			t.Errorf("squared_trace=%v", rho.SquareTrace())
		}

		if math.Abs(rho.ExpectedValue(c.m)-c.v) > c.eps {
			t.Errorf("expected_value=%v", rho.ExpectedValue(c.m))
		}
	}
}

func TestMeasure(t *testing.T) {
	cases := []struct {
		p    []float64
		q    []*qubit.Qubit
		m    *qubit.Qubit
		want float64
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
		got, _ := density.New(c.p, c.q)
		if got.Measure(c.m) != c.want {
			t.Fail()
		}
	}
}

func TestPartialTrace(t *testing.T) {
	type Case struct {
		index int
		want  [][]complex128
	}

	cases := []struct {
		p   []float64
		q   []*qubit.Qubit
		cs  []Case
		eps float64
	}{
		{
			[]float64{1.0},
			[]*qubit.Qubit{qubit.Zero(2)},
			[]Case{
				{0, [][]complex128{{1, 0}, {0, 0}}},
				{1, [][]complex128{{1, 0}, {0, 0}}},
			},
			epsilon.E13(),
		},
		{
			[]float64{1.0},
			[]*qubit.Qubit{qubit.One(2)},
			[]Case{
				{0, [][]complex128{{0, 0}, {0, 1}}},
				{1, [][]complex128{{0, 0}, {0, 1}}},
			},
			epsilon.E13(),
		},
		{
			[]float64{1.0},
			[]*qubit.Qubit{qubit.Zero(2).Apply(gate.H(2))},
			[]Case{
				{0, [][]complex128{{0.5, 0.5}, {0.5, 0.5}}},
				{1, [][]complex128{{0.5, 0.5}, {0.5, 0.5}}},
			},
			epsilon.E13(),
		},
		{
			[]float64{0.5, 0.5},
			[]*qubit.Qubit{qubit.Zero(2), qubit.One(2)},
			[]Case{
				{0, [][]complex128{{0.5, 0}, {0, 0.5}}},
				{1, [][]complex128{{0.5, 0}, {0, 0.5}}},
			},
			epsilon.E13(),
		},
		{
			[]float64{0.5, 0.5},
			[]*qubit.Qubit{
				qubit.Zero(2).Apply(gate.H(2)),
				qubit.One(2),
			},
			[]Case{
				{0, [][]complex128{{0.25, 0.25}, {0.25, 0.75}}},
				{1, [][]complex128{{0.25, 0.25}, {0.25, 0.75}}},
			},
			epsilon.E13(),
		},
		{
			[]float64{0.75, 0.25},
			[]*qubit.Qubit{
				qubit.Zero(2).Apply(gate.H(2)),
				qubit.One(2).Apply(gate.H(2)),
			},
			[]Case{
				{0, [][]complex128{{0.5, 0.25}, {0.25, 0.5}}},
				{1, [][]complex128{{0.5, 0.25}, {0.25, 0.5}}},
			},
			epsilon.E13(),
		},
		{
			[]float64{0.25, 0.75},
			[]*qubit.Qubit{
				qubit.Zero(2).Apply(gate.H(2)),
				qubit.One(2).Apply(gate.H(2)),
			},
			[]Case{
				{0, [][]complex128{{0.5, -0.25}, {-0.25, 0.5}}},
				{1, [][]complex128{{0.5, -0.25}, {-0.25, 0.5}}},
			},
			epsilon.E13(),
		},
	}

	for _, c := range cases {
		for _, cs := range c.cs {
			rho, _ := density.New(c.p, c.q)

			got := rho.PartialTrace(cs.index)
			p, q := got.Dimension()
			if p != len(cs.want) || q != len(cs.want) {
				t.Errorf("got=%v, %v want=%v", p, q, cs.want)
			}

			for i := 0; i < len(cs.want); i++ {
				for j := 0; j < len(cs.want[0]); j++ {
					if cmplx.Abs(got.Raw()[i][j]-cs.want[i][j]) > c.eps {
						t.Errorf("%v:%v, got=%v, want=%v", i, j, got.Raw()[i][j], cs.want[i][j])
					}
				}
			}

			if math.Abs(got.Trace()-1) > c.eps {
				t.Errorf("trace: got=%v, want=%v", got.Trace(), 1)
			}

			if got.SquareTrace() > 1+c.eps {
				t.Errorf("square_trace: got=%v > 1", got.SquareTrace())
			}
		}
	}
}

func TestApply(t *testing.T) {
	cases := []struct {
		p    []float64
		q    []*qubit.Qubit
		g    matrix.Matrix
		m    *qubit.Qubit
		want float64
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
		got, _ := density.New(c.p, c.q)
		if got.Apply(c.g).Measure(c.m) != c.want {
			t.Fail()
		}
	}
}

func TestMatrixNew(t *testing.T) {
	cases := []struct {
		in   []float64
		want error
	}{
		{[]float64{1.0, 2.0}, fmt.Errorf("invalid length. len(p)=2, len(q)=1")},
		{[]float64{1.5}, fmt.Errorf("add: p must be 0 <= p =< 1. p=1.5")},
	}

	for _, c := range cases {
		_, got := density.New(c.in, []*qubit.Qubit{qubit.Zero()})
		if got.Error() != c.want.Error() {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}

func TestMatrixAddError(t *testing.T) {
	rho, _ := density.New([]float64{.50}, []*qubit.Qubit{qubit.Zero()})

	got := rho.Add(0.5, qubit.One(2)).Error()
	want := "invalid dimension. m=2 n=4"
	if got != want {
		t.Errorf("got=%v, want=%v", got, want)
	}
}

func TestDepolarizingError(t *testing.T) {
	rho, _ := density.New([]float64{1.0}, []*qubit.Qubit{qubit.Zero()})

	cases := []struct {
		in   float64
		want error
	}{
		{-1, fmt.Errorf("p must be 0 <= p =< 1. p=-1")},
	}

	for _, c := range cases {
		_, got := rho.Depolarizing(c.in)
		if got.Error() != c.want.Error() {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}

func TestFlipError(t *testing.T) {
	cases := []struct {
		in   float64
		want error
	}{
		{-1, fmt.Errorf("p must be 0 <= p =< 1. p=-1")},
	}

	for _, c := range cases {
		_, _, got := density.BitPhaseFlip(c.in)
		if got.Error() != c.want.Error() {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}
