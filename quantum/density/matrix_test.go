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
	phi := qubit.Zero(2).Apply(
		gate.H().TensorProduct(gate.I()),
		gate.CNOT(2, 0, 1),
	)

	rho := density.New([]density.State{
		{1.0, phi}, // pure state
	})

	qb := rho.Qubits()
	p0 := rho.PartialTrace(qb[0]) // Partial trace over qubit 0: returns the reduced density matrix for qubit 1
	p1 := rho.PartialTrace(qb[1]) // Partial trace over qubit 1: returns the reduced density matrix for qubit 0

	fmt.Printf("trace: %.2v, purity: %.2v\n", rho.Trace(), rho.Purity())
	fmt.Printf("trace: %.2v, purity: %.2v\n", p0.Trace(), p0.Purity())
	fmt.Printf("trace: %.2v, purity: %.2v\n", p1.Trace(), p1.Purity())

	q00 := qubit.TensorProduct(qubit.Zero(), qubit.Zero())
	q01 := qubit.TensorProduct(qubit.Zero(), qubit.One())
	q10 := qubit.TensorProduct(qubit.One(), qubit.Zero())
	q11 := qubit.TensorProduct(qubit.One(), qubit.One())

	m00 := rho.Probability(q00) // 0.5
	m01 := rho.Probability(q01) // zero
	m10 := rho.Probability(q10) // zero
	m11 := rho.Probability(q11) // 0.5
	fmt.Printf("%.2f, %.2f, %.2f, %.2f\n", m00, m01, m10, m11)

	fmt.Println(rho.Project(q00).Underlying().Data)
	fmt.Println(rho.Project(q01).IsZero())
	fmt.Println(rho.Project(q10).IsZero())
	fmt.Println(rho.Project(q11).Underlying().Data)

	// Output:
	// trace: 1, purity: 1
	// trace: 1, purity: 0.5
	// trace: 1, purity: 0.5
	// 0.50, 0.00, 0.00, 0.50
	// [(1+0i) (0+0i) (0+0i) (0+0i) (0+0i) (0+0i) (0+0i) (0+0i) (0+0i) (0+0i) (0+0i) (0+0i) (0+0i) (0+0i) (0+0i) (0+0i)]
	// true
	// true
	// [(0+0i) (0+0i) (0+0i) (0+0i) (0+0i) (0+0i) (0+0i) (0+0i) (0+0i) (0+0i) (0+0i) (0+0i) (0+0i) (0+0i) (0+0i) (1+0i)]
}

func ExampleMatrix_ExpectedValue() {
	rho := density.New([]density.State{
		{0.1, qubit.Zero()},
		{0.9, qubit.Zero().Apply(gate.H())},
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

func ExampleMatrix_Trace() {
	pure := density.New([]density.State{{1.0, qubit.Zero()}})
	mixed := density.New([]density.State{{0.1, qubit.Zero()}, {0.9, qubit.One()}})

	fmt.Printf("pure:  %.2f\n", pure.Trace())
	fmt.Printf("mixed: %.2f\n", mixed.Trace())

	// Output:
	// pure:  1.00
	// mixed: 1.00
}

func ExampleMatrix_Purity() {
	pure := density.New([]density.State{{1.0, qubit.Zero()}})
	mixed := density.New([]density.State{{0.1, qubit.Zero()}, {0.9, qubit.One()}})

	fmt.Printf("pure:  %.2f\n", pure.Purity())
	fmt.Printf("mixed: %.2f\n", mixed.Purity())

	// Output:
	// pure:  1.00
	// mixed: 0.82
}

func ExampleMatrix_PartialTrace() {
	rho := density.New([]density.State{
		{0.5, qubit.Zero(2)},
		{0.5, qubit.One().TensorProduct(qubit.Zero())},
	})

	for _, r := range rho.Underlying().Seq2() {
		fmt.Println(r)
	}
	fmt.Println()

	qb := rho.Qubits()
	p0 := rho.PartialTrace(qb[0]) // Partial trace over qubit 0: returns the reduced density matrix for qubit 1
	p1 := rho.PartialTrace(qb[1]) // Partial trace over qubit 1: returns the reduced density matrix for qubit 0

	fmt.Printf("trace: %.2v, purity: %.2v\n", rho.Trace(), rho.Purity())
	fmt.Printf("trace: %.2v, purity: %.2v\n", p0.Trace(), p0.Purity()) // qubit 1: pure |0⟩
	fmt.Printf("trace: %.2v, purity: %.2v\n", p1.Trace(), p1.Purity()) // qubit 0: maximally mixed

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
	phi := qubit.Zero(3).Apply(
		matrix.TensorProduct(gate.H(), gate.I(), gate.I()),
		gate.CNOT(3, 0, 1),
	)

	rho := density.New([]density.State{
		{1.0, phi},
	})

	qb := rho.Qubits()
	p0 := rho.PartialTrace(qb[0])
	p1 := rho.PartialTrace(qb[1])
	p2 := rho.PartialTrace(qb[2])

	fmt.Printf("trace: %.2v, purity: %.2v\n", rho.Trace(), rho.Purity())
	fmt.Printf("trace: %.2v, purity: %.2v\n", p0.Trace(), p0.Purity())
	fmt.Printf("trace: %.2v, purity: %.2v\n", p1.Trace(), p1.Purity())
	fmt.Printf("trace: %.2v, purity: %.2v\n", p2.Trace(), p2.Purity())

	// Output:
	// trace: 1, purity: 1
	// trace: 1, purity: 0.5
	// trace: 1, purity: 0.5
	// trace: 1, purity: 1
}

func ExampleMatrix_Depolarizing() {
	rho := density.New([]density.State{{1.0, qubit.Zero()}})
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
			[]density.State{{0.1, qubit.Zero()}, {0.9, qubit.Zero().Apply(gate.H())}},
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

func TestProbabilityOf(t *testing.T) {
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
