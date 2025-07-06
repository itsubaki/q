package density_test

import (
	"fmt"
	"math"
	"math/cmplx"
	"testing"

	"github.com/itsubaki/q/math/epsilon"
	"github.com/itsubaki/q/math/matrix"
	"github.com/itsubaki/q/quantum/density"
	"github.com/itsubaki/q/quantum/gate"
	"github.com/itsubaki/q/quantum/qubit"
)

func ExampleMatrix_bell() {
	rho := density.NewPureState(qubit.Zero(2).Apply(
		gate.H().TensorProduct(gate.I()),
		gate.CNOT(2, 0, 1),
	))

	qb := rho.Qubits()
	s1 := rho.PartialTrace(qb[0]) // Partial trace over qubit 0: returns the reduced density matrix for qubit 1
	s0 := rho.PartialTrace(qb[1]) // Partial trace over qubit 1: returns the reduced density matrix for qubit 0

	fmt.Printf("trace: %.2v, purity: %.2v\n", rho.Trace(), rho.Purity())
	fmt.Printf("trace: %.2v, purity: %.2v\n", s1.Trace(), s1.Purity())
	fmt.Printf("trace: %.2v, purity: %.2v\n", s0.Trace(), s0.Purity())

	// Output:
	// trace: 1, purity: 1
	// trace: 1, purity: 0.5
	// trace: 1, purity: 0.5
}

func ExampleMatrix_prob() {
	rho := density.NewPureState(qubit.Zero(2).Apply(
		gate.H().TensorProduct(gate.I()),
		gate.CNOT(2, 0, 1),
	))

	basis := []*qubit.Qubit{
		qubit.NewFrom("00"),
		qubit.NewFrom("01"),
		qubit.NewFrom("10"),
		qubit.NewFrom("11"),
	}

	for _, b := range basis {
		fmt.Printf("%v, %.2f\n", b.State(), rho.Probability(b))
	}
	fmt.Println()

	for _, b := range basis {
		p, sigma := rho.Project(b)
		fmt.Printf("%v, p: %.2f\n", b.State(), p)

		for _, r := range sigma.Underlying().Seq2() {
			fmt.Println(r)
		}
	}

	// Output:
	// [[00][  0]( 1.0000 0.0000i): 1.0000], 0.50
	// [[01][  1]( 1.0000 0.0000i): 1.0000], 0.00
	// [[10][  2]( 1.0000 0.0000i): 1.0000], 0.00
	// [[11][  3]( 1.0000 0.0000i): 1.0000], 0.50
	//
	// [[00][  0]( 1.0000 0.0000i): 1.0000], p: 0.50
	// [(1+0i) (0+0i) (0+0i) (0+0i)]
	// [(0+0i) (0+0i) (0+0i) (0+0i)]
	// [(0+0i) (0+0i) (0+0i) (0+0i)]
	// [(0+0i) (0+0i) (0+0i) (0+0i)]
	// [[01][  1]( 1.0000 0.0000i): 1.0000], p: 0.00
	// [(0+0i) (0+0i) (0+0i) (0+0i)]
	// [(0+0i) (0+0i) (0+0i) (0+0i)]
	// [(0+0i) (0+0i) (0+0i) (0+0i)]
	// [(0+0i) (0+0i) (0+0i) (0+0i)]
	// [[10][  2]( 1.0000 0.0000i): 1.0000], p: 0.00
	// [(0+0i) (0+0i) (0+0i) (0+0i)]
	// [(0+0i) (0+0i) (0+0i) (0+0i)]
	// [(0+0i) (0+0i) (0+0i) (0+0i)]
	// [(0+0i) (0+0i) (0+0i) (0+0i)]
	// [[11][  3]( 1.0000 0.0000i): 1.0000], p: 0.50
	// [(0+0i) (0+0i) (0+0i) (0+0i)]
	// [(0+0i) (0+0i) (0+0i) (0+0i)]
	// [(0+0i) (0+0i) (0+0i) (0+0i)]
	// [(0+0i) (0+0i) (0+0i) (1+0i)]
}
func ExampleMatrix_ExpectedValue() {
	rho := density.New([]density.State{
		{0.1, qubit.Zero()},
		{0.9, qubit.Plus()},
	})

	fmt.Printf("X: %.2v\n", rho.ExpectedValue(gate.X()))
	fmt.Printf("Y: %.2v\n", rho.ExpectedValue(gate.Y()))
	fmt.Printf("Z: %.2v\n", rho.ExpectedValue(gate.Z()))

	// Output:
	// X: 0.9
	// Y: 0
	// Z: 0.1
}

func ExampleMatrix_Underlying() {
	rho := density.New([]density.State{
		{0.1, qubit.Zero()},
		{0.9, qubit.One()},
	})

	for _, r := range rho.Underlying().Seq2() {
		fmt.Println(r)
	}

	// Output:
	// [(0.1+0i) (0+0i)]
	// [(0+0i) (0.9+0i)]
}

func ExampleMatrix_Probability() {
	rho := density.New([]density.State{
		{0.1, qubit.Zero()},
		{0.9, qubit.One()},
	})

	fmt.Printf("0: %.2v\n", rho.Probability(qubit.Zero()))
	fmt.Printf("1: %.2v\n", rho.Probability(qubit.One()))

	// Output:
	// 0: 0.1
	// 1: 0.9
}

func ExampleMatrix_IsHermite() {
	s0 := density.NewPureState(qubit.Zero())
	s1 := density.New([]density.State{{0.1, qubit.Zero()}, {0.9, qubit.One()}})

	fmt.Println(s0.IsHermite())
	fmt.Println(s1.IsHermite())

	// Output:
	// true
	// true
}

func ExampleMatrix_Trace() {
	s0 := density.NewPureState(qubit.Zero())
	s1 := density.New([]density.State{{0.1, qubit.Zero()}, {0.9, qubit.One()}})

	fmt.Printf("pure : %.2f\n", s0.Trace())
	fmt.Printf("mixed: %.2f\n", s1.Trace())

	// Output:
	// pure : 1.00
	// mixed: 1.00
}

func ExampleMatrix_Purity() {
	s0 := density.NewPureState(qubit.Zero())
	s1 := density.New([]density.State{{0.1, qubit.Zero()}, {0.9, qubit.One()}})

	fmt.Printf("pure : %.2f, %v\n", s0.Purity(), s0.IsPure())
	fmt.Printf("mixed: %.2f, %v\n", s1.Purity(), s1.IsMixed())

	// Output:
	// pure : 1.00, true
	// mixed: 0.82, true
}

func ExampleMatrix_TensorProduct() {
	s0 := density.NewPureState(qubit.Zero())
	s1 := density.NewPureState(qubit.One())

	s01 := s0.TensorProduct(s1)
	for _, r := range s01.Underlying().Seq2() {
		fmt.Println(r)
	}

	// Output:
	// [(0+0i) (0+0i) (0+0i) (0+0i)]
	// [(0+0i) (1+0i) (0+0i) (0+0i)]
	// [(0+0i) (0+0i) (0+0i) (0+0i)]
	// [(0+0i) (0+0i) (0+0i) (0+0i)]
}

func ExampleMatrix_PartialTrace() {
	rho := density.New([]density.State{
		{0.5, qubit.NewFrom("00")},
		{0.5, qubit.NewFrom("10")},
	})

	for _, r := range rho.Underlying().Seq2() {
		fmt.Println(r)
	}
	fmt.Println()

	qb := rho.Qubits()
	s1 := rho.PartialTrace(qb[0])
	s0 := rho.PartialTrace(qb[1])

	fmt.Printf("trace: %.2v, purity: %.2v\n", rho.Trace(), rho.Purity())
	fmt.Printf("trace: %.2v, purity: %.2v\n", s1.Trace(), s1.Purity())
	fmt.Printf("trace: %.2v, purity: %.2v\n", s0.Trace(), s0.Purity())

	// Output:
	// [(0.5+0i) (0+0i) (0+0i) (0+0i)]
	// [(0+0i) (0+0i) (0+0i) (0+0i)]
	// [(0+0i) (0+0i) (0.5+0i) (0+0i)]
	// [(0+0i) (0+0i) (0+0i) (0+0i)]
	//
	// trace: 1, purity: 0.5
	// trace: 1, purity: 1
	// trace: 1, purity: 0.5
}

func ExampleMatrix_PartialTrace_x8() {
	rho := density.NewPureState(qubit.Zero(3).Apply(
		matrix.TensorProduct(gate.H(), gate.I(), gate.I()),
		gate.CNOT(3, 0, 1),
	))

	qb := rho.Qubits()
	s12 := rho.PartialTrace(qb[0])
	s02 := rho.PartialTrace(qb[1])
	s01 := rho.PartialTrace(qb[2])

	fmt.Printf("trace: %.2v, purity: %.2v\n", rho.Trace(), rho.Purity())
	fmt.Printf("trace: %.2v, purity: %.2v\n", s12.Trace(), s12.Purity())
	fmt.Printf("trace: %.2v, purity: %.2v\n", s02.Trace(), s02.Purity())
	fmt.Printf("trace: %.2v, purity: %.2v\n", s01.Trace(), s01.Purity())

	// Output:
	// trace: 1, purity: 1
	// trace: 1, purity: 0.5
	// trace: 1, purity: 0.5
	// trace: 1, purity: 1
}

func ExampleMatrix_Depolarizing() {
	rho := density.NewPureState(qubit.Zero())
	fmt.Printf("0: %.2f\n", rho.Probability(qubit.Zero()))
	fmt.Printf("1: %.2f\n", rho.Probability(qubit.One()))
	fmt.Println()

	dep := rho.Depolarizing(1)
	fmt.Printf("0: %.2f\n", dep.Probability(qubit.Zero()))
	fmt.Printf("1: %.2f\n", dep.Probability(qubit.One()))

	// Output:
	// 0: 1.00
	// 1: 0.00
	//
	// 0: 0.50
	// 1: 0.50
}

func ExampleMatrix_ApplyChannel() {
	rho := density.NewPureState(qubit.Zero(2))

	qb := rho.Qubits()
	x := rho.ApplyChannel(0.3, gate.X(), qb[0])

	fmt.Printf("%.2f\n", x.Probability(qubit.NewFrom("00")))
	fmt.Printf("%.2f\n", x.Probability(qubit.NewFrom("10")))

	// Output:
	// 0.70
	// 0.30
}

func ExampleMatrix_ApplyChannel_qb1() {
	rho := density.NewPureState(qubit.Zero(2))

	qb := rho.Qubits()
	x := rho.ApplyChannel(0.3, gate.X(), qb[1])

	fmt.Printf("%.2f\n", x.Probability(qubit.NewFrom("00")))
	fmt.Printf("%.2f\n", x.Probability(qubit.NewFrom("01")))

	// Output:
	// 0.70
	// 0.30
}

func ExampleMatrix_BitFlip() {
	rho := density.NewPureState(qubit.Zero())

	qb := rho.Qubits()
	x := rho.BitFlip(0.3, qb[0])

	fmt.Printf("%.2f\n", x.Probability(qubit.Zero()))
	fmt.Printf("%.2f\n", x.Probability(qubit.One()))

	// Output:
	// 0.70
	// 0.30
}

func ExampleMatrix_BitPhaseFlip() {
	rho := density.NewPureState(qubit.Plus())

	qb := rho.Qubits()
	y := rho.BitPhaseFlip(0.3, qb[0])

	fmt.Printf("%.2f\n", y.Probability(qubit.Plus()))
	fmt.Printf("%.2f\n", y.Probability(qubit.Minus()))

	// Output:
	// 0.70
	// 0.30
}

func ExampleMatrix_PhaseFlip() {
	rho := density.NewPureState(qubit.Plus())

	qb := rho.Qubits()
	z := rho.PhaseFlip(0.3, qb[0])

	fmt.Printf("%.2f\n", z.Probability(qubit.Plus()))
	fmt.Printf("%.2f\n", z.Probability(qubit.Minus()))

	// Output:
	// 0.70
	// 0.30
}

func ExampleMatrix_phaseAndBitPhaseFlip() {
	rho := density.NewPureState(qubit.Zero())

	qb := rho.Qubits()
	y := rho.BitPhaseFlip(0.3, qb[0])
	z := rho.PhaseFlip(0.3, qb[0])

	fmt.Printf("%.2f\n", y.Probability(qubit.Zero()))
	fmt.Printf("%.2f\n", y.Probability(qubit.One()))
	fmt.Printf("%.2f\n", z.Probability(qubit.Zero()))
	fmt.Printf("%.2f\n", z.Probability(qubit.One()))

	// Output:
	// 0.70
	// 0.30
	// 1.00
	// 0.00
}

func TestExpectedValue(t *testing.T) {
	cases := []struct {
		s        []density.State
		tr, sqtr float64
		m        *matrix.Matrix
		v        float64
		eps      float64
	}{
		{
			[]density.State{{0.1, qubit.Zero()}, {0.9, qubit.One()}},
			1, 0.82,
			gate.X(), 0.0,
			epsilon.E13(),
		},
		{
			[]density.State{{0.1, qubit.Zero()}, {0.9, qubit.Plus()}},
			1, 0.91,
			gate.X(), 0.9,
			epsilon.E13(),
		},
	}

	for _, c := range cases {
		rho := density.New(c.s)

		if math.Abs(rho.Trace()-c.tr) > c.eps {
			t.Errorf("trace=%v", rho.Trace())
		}

		if math.Abs(rho.Purity()-c.sqtr) > c.eps {
			t.Errorf("purity=%v", rho.Purity())
		}

		if math.Abs(rho.ExpectedValue(c.m)-c.v) > c.eps {
			t.Errorf("expected_value=%v", rho.ExpectedValue(c.m))
		}
	}
}

func TestProbability(t *testing.T) {
	cases := []struct {
		s    []density.State
		m    *qubit.Qubit
		want float64
	}{
		{
			[]density.State{{1, qubit.Zero()}},
			qubit.Zero(),
			1,
		},
		{
			[]density.State{{1, qubit.Zero()}},
			qubit.One(),
			0,
		},
	}

	for _, c := range cases {
		if density.New(c.s).Probability(c.m) != c.want {
			t.Fail()
		}
	}
}

func TestPartialTrace(t *testing.T) {
	type Case struct {
		index density.Qubit
		want  [][]complex128
	}

	cases := []struct {
		s   []density.State
		cs  []Case
		eps float64
	}{
		{
			[]density.State{{1.0, qubit.Zero(2)}},
			[]Case{
				{0, [][]complex128{{1, 0}, {0, 0}}},
				{1, [][]complex128{{1, 0}, {0, 0}}},
			},
			epsilon.E13(),
		},
		{
			[]density.State{{1.0, qubit.One(2)}},
			[]Case{
				{0, [][]complex128{{0, 0}, {0, 1}}},
				{1, [][]complex128{{0, 0}, {0, 1}}},
			},
			epsilon.E13(),
		},
		{
			[]density.State{{1.0, qubit.Zero(2).Apply(gate.H(2))}},
			[]Case{
				{0, [][]complex128{{0.5, 0.5}, {0.5, 0.5}}},
				{1, [][]complex128{{0.5, 0.5}, {0.5, 0.5}}},
			},
			epsilon.E13(),
		},
		{
			[]density.State{{0.5, qubit.Zero(2)}, {0.5, qubit.One(2)}},
			[]Case{
				{0, [][]complex128{{0.5, 0}, {0, 0.5}}},
				{1, [][]complex128{{0.5, 0}, {0, 0.5}}},
			},
			epsilon.E13(),
		},
		{
			[]density.State{{0.5, qubit.Zero(2).Apply(gate.H(2))}, {0.5, qubit.One(2)}},
			[]Case{
				{0, [][]complex128{{0.25, 0.25}, {0.25, 0.75}}},
				{1, [][]complex128{{0.25, 0.25}, {0.25, 0.75}}},
			},
			epsilon.E13(),
		},
		{
			[]density.State{{0.75, qubit.Zero(2).Apply(gate.H(2))}, {0.25, qubit.One(2).Apply(gate.H(2))}},
			[]Case{
				{0, [][]complex128{{0.5, 0.25}, {0.25, 0.5}}},
				{1, [][]complex128{{0.5, 0.25}, {0.25, 0.5}}},
			},
			epsilon.E13(),
		},
		{
			[]density.State{{0.25, qubit.Zero(2).Apply(gate.H(2))}, {0.75, qubit.One(2).Apply(gate.H(2))}},
			[]Case{
				{0, [][]complex128{{0.5, -0.25}, {-0.25, 0.5}}},
				{1, [][]complex128{{0.5, -0.25}, {-0.25, 0.5}}},
			},
			epsilon.E13(),
		},
	}

	for _, c := range cases {
		for _, cs := range c.cs {
			got := density.New(c.s).PartialTrace(cs.index)

			p, q := got.Dimension()
			if p != len(cs.want) || q != len(cs.want) {
				t.Errorf("got=%v, %v want=%v", p, q, cs.want)
			}

			for i := range cs.want {
				for j := range cs.want[0] {
					if cmplx.Abs(got.At(i, j)-cs.want[i][j]) > c.eps {
						t.Errorf("%v:%v, got=%v, want=%v", i, j, got.At(i, j), cs.want[i][j])
					}
				}
			}

			if math.Abs(got.Trace()-1) > c.eps {
				t.Errorf("trace: got=%v, want=%v", got.Trace(), 1)
			}

			if got.Purity() > 1+c.eps {
				t.Errorf("purity: got=%v > 1", got.Purity())
			}
		}
	}
}

func TestApply(t *testing.T) {
	cases := []struct {
		s    []density.State
		g    *matrix.Matrix
		m    *qubit.Qubit
		want float64
	}{
		{
			[]density.State{{1, qubit.Zero()}},
			gate.X(),
			qubit.Zero(),
			0,
		},
		{
			[]density.State{{1, qubit.Zero()}},
			gate.X(),
			qubit.One(),
			1,
		},
	}

	for _, c := range cases {
		if density.New(c.s).Apply(c.g).Probability(c.m) != c.want {
			t.Fail()
		}
	}
}

func TestMatrix_IsZero(t *testing.T) {
	cases := []struct {
		in   *density.Matrix
		want bool
	}{
		{
			density.New([]density.State{
				{1.0, qubit.Zero()},
			}),
			false,
		},
	}

	for _, c := range cases {
		if c.in.IsZero() != c.want {
			t.Fail()
		}
	}
}
