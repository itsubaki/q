package density_test

import (
	"fmt"
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

func ExampleMatrix_ComputationalBasis() {
	rho := density.NewPureState(qubit.Zero(2).Apply(
		gate.H().TensorProduct(gate.I()),
		gate.CNOT(2, 0, 1),
	))

	basis := rho.ComputationalBasis()
	for _, b := range basis {
		fmt.Printf("%v: %.2f\n", b.State(), rho.Probability(b))
	}
	fmt.Println()

	for _, b := range basis {
		p, sigma := rho.Project(b)
		fmt.Printf("%v: %.2f\n", b.State(), p)

		for _, r := range sigma.Seq2() {
			fmt.Println(r)
		}
	}

	// Output:
	// [[00][  0]( 1.0000 0.0000i): 1.0000]: 0.50
	// [[01][  1]( 1.0000 0.0000i): 1.0000]: 0.00
	// [[10][  2]( 1.0000 0.0000i): 1.0000]: 0.00
	// [[11][  3]( 1.0000 0.0000i): 1.0000]: 0.50
	//
	// [[00][  0]( 1.0000 0.0000i): 1.0000]: 0.50
	// [(1+0i) (0+0i) (0+0i) (0+0i)]
	// [(0+0i) (0+0i) (0+0i) (0+0i)]
	// [(0+0i) (0+0i) (0+0i) (0+0i)]
	// [(0+0i) (0+0i) (0+0i) (0+0i)]
	// [[01][  1]( 1.0000 0.0000i): 1.0000]: 0.00
	// [(0+0i) (0+0i) (0+0i) (0+0i)]
	// [(0+0i) (0+0i) (0+0i) (0+0i)]
	// [(0+0i) (0+0i) (0+0i) (0+0i)]
	// [(0+0i) (0+0i) (0+0i) (0+0i)]
	// [[10][  2]( 1.0000 0.0000i): 1.0000]: 0.00
	// [(0+0i) (0+0i) (0+0i) (0+0i)]
	// [(0+0i) (0+0i) (0+0i) (0+0i)]
	// [(0+0i) (0+0i) (0+0i) (0+0i)]
	// [(0+0i) (0+0i) (0+0i) (0+0i)]
	// [[11][  3]( 1.0000 0.0000i): 1.0000]: 0.50
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
	s0 := density.NewZeroState()
	s1 := density.New([]density.State{
		{0.1, qubit.Zero()},
		{0.9, qubit.One()},
	})

	fmt.Println(s0.IsHermite())
	fmt.Println(s1.IsHermite())

	// Output:
	// true
	// true
}

func ExampleMatrix_Trace() {
	s0 := density.NewZeroState()
	s1 := density.New([]density.State{
		{0.1, qubit.Zero()},
		{0.9, qubit.One()},
	})

	fmt.Printf("pure : %.2f\n", s0.Trace())
	fmt.Printf("mixed: %.2f\n", s1.Trace())

	// Output:
	// pure : 1.00
	// mixed: 1.00
}

func ExampleMatrix_Purity() {
	s0 := density.NewZeroState()
	s1 := density.New([]density.State{
		{0.1, qubit.Zero()},
		{0.9, qubit.One()},
	})

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
	for _, r := range s01.Seq2() {
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
		{0.5, qubit.From("00")},
		{0.5, qubit.From("10")},
	})

	for _, r := range rho.Seq2() {
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
	rho := density.NewZeroState()
	fmt.Printf("0: %.2f\n", rho.Probability(qubit.Zero()))
	fmt.Printf("1: %.2f\n", rho.Probability(qubit.One()))
	fmt.Println()

	// XrhoX = |1><1|, YrhoY = |1><1|, ZrhoZ = |0><0|
	// E(rho) = 0.7|0><0| + 0.1|1><1| + 0.1|1><1| + 0.1|0><0| = 0.8|0><0| + 0.2|1><1|
	dep := rho.Depolarizing(0.3)
	fmt.Printf("0: %.2f\n", dep.Probability(qubit.Zero()))
	fmt.Printf("1: %.2f\n", dep.Probability(qubit.One()))

	// Output:
	// 0: 1.00
	// 1: 0.00
	//
	// 0: 0.80
	// 1: 0.20
}

func ExampleMatrix_ApplyChannel() {
	rho := density.NewZeroState(2)

	qb := rho.Qubits()
	x := rho.ApplyChannel(0.3, gate.X(), qb[0])

	fmt.Printf("%.2f\n", x.Probability(qubit.From("00")))
	fmt.Printf("%.2f\n", x.Probability(qubit.From("10")))

	// Output:
	// 0.70
	// 0.30
}

func ExampleMatrix_ApplyChannel_qb1() {
	rho := density.NewZeroState(2)

	qb := rho.Qubits()
	x := rho.ApplyChannel(0.3, gate.X(), qb[1])

	fmt.Printf("%.2f\n", x.Probability(qubit.From("00")))
	fmt.Printf("%.2f\n", x.Probability(qubit.From("01")))

	// Output:
	// 0.70
	// 0.30
}

func ExampleMatrix_BitFlip() {
	rho := density.NewZeroState()

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
	rho := density.NewZeroState()

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
	}{
		{
			s: []density.State{
				{
					Probability: 0.1,
					Qubit:       qubit.Zero(),
				},
				{
					Probability: 0.9,
					Qubit:       qubit.One(),
				},
			},
			tr:   1,
			sqtr: 0.82,
			m:    gate.X(),
			v:    0.0,
		},
		{
			s: []density.State{
				{
					Probability: 0.1,
					Qubit:       qubit.Zero(),
				},
				{
					Probability: 0.9,
					Qubit:       qubit.Plus(),
				},
			},
			tr:   1,
			sqtr: 0.91,
			m:    gate.X(),
			v:    0.9,
		},
	}

	for _, c := range cases {
		rho := density.New(c.s)

		if !epsilon.IsCloseF64(rho.Trace(), c.tr) {
			t.Errorf("trace=%v", rho.Trace())
		}

		if !epsilon.IsCloseF64(rho.Purity(), c.sqtr) {
			t.Errorf("purity=%v", rho.Purity())
		}

		if !epsilon.IsCloseF64(rho.ExpectedValue(c.m), c.v) {
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
			s: []density.State{
				{
					Probability: 1,
					Qubit:       qubit.Zero(),
				},
			},
			m:    qubit.Zero(),
			want: 1,
		},
		{
			s: []density.State{
				{
					Probability: 1,
					Qubit:       qubit.Zero(),
				},
			},
			m:    qubit.One(),
			want: 0,
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
		idx  []density.Qubit
		want [][]complex128
	}

	cases := []struct {
		s  []density.State
		cs []Case
	}{
		{
			s: []density.State{
				{
					Probability: 1.0,
					Qubit:       qubit.Zero(2),
				},
			},
			cs: []Case{
				{
					idx: []density.Qubit{0},
					want: [][]complex128{
						{1, 0},
						{0, 0},
					},
				},
				{
					idx: []density.Qubit{1},
					want: [][]complex128{
						{1, 0},
						{0, 0},
					},
				},
			},
		},
		{
			s: []density.State{
				{
					Probability: 1.0,
					Qubit:       qubit.One(2),
				},
			},
			cs: []Case{
				{
					idx: []density.Qubit{0},
					want: [][]complex128{
						{0, 0},
						{0, 1},
					},
				},
				{
					idx: []density.Qubit{1},
					want: [][]complex128{
						{0, 0},
						{0, 1},
					},
				},
			},
		},
		{
			s: []density.State{
				{
					Probability: 1.0,
					Qubit:       qubit.Zero(2).Apply(gate.H(2)),
				},
			},
			cs: []Case{
				{
					idx: []density.Qubit{0},
					want: [][]complex128{
						{0.5, 0.5},
						{0.5, 0.5},
					},
				},
				{
					idx: []density.Qubit{1},
					want: [][]complex128{
						{0.5, 0.5},
						{0.5, 0.5},
					},
				},
			},
		},
		{
			s: []density.State{
				{
					Probability: 0.5,
					Qubit:       qubit.Zero(2),
				},
				{
					Probability: 0.5,
					Qubit:       qubit.One(2),
				},
			},
			cs: []Case{
				{
					idx: []density.Qubit{0},
					want: [][]complex128{
						{0.5, 0},
						{0, 0.5},
					},
				},
				{
					idx: []density.Qubit{1},
					want: [][]complex128{
						{0.5, 0},
						{0, 0.5},
					},
				},
			},
		},
		{
			s: []density.State{
				{
					Probability: 0.5,
					Qubit:       qubit.Zero(2).Apply(gate.H(2)),
				},
				{
					Probability: 0.5,
					Qubit:       qubit.One(2),
				},
			},
			cs: []Case{
				{
					idx: []density.Qubit{0},
					want: [][]complex128{
						{0.25, 0.25},
						{0.25, 0.75},
					},
				},
				{
					idx: []density.Qubit{1},
					want: [][]complex128{
						{0.25, 0.25},
						{0.25, 0.75},
					},
				},
			},
		},
		{
			s: []density.State{
				{
					Probability: 0.75,
					Qubit:       qubit.Zero(2).Apply(gate.H(2)),
				},
				{
					Probability: 0.25,
					Qubit:       qubit.One(2).Apply(gate.H(2)),
				},
			},
			cs: []Case{
				{
					idx: []density.Qubit{0},
					want: [][]complex128{
						{0.5, 0.25},
						{0.25, 0.5},
					},
				},
				{
					idx: []density.Qubit{1},
					want: [][]complex128{
						{0.5, 0.25},
						{0.25, 0.5},
					},
				},
			},
		},
		{
			s: []density.State{
				{
					Probability: 0.25,
					Qubit:       qubit.Zero(2).Apply(gate.H(2)),
				},
				{
					Probability: 0.75,
					Qubit:       qubit.One(2).Apply(gate.H(2)),
				},
			},
			cs: []Case{
				{
					idx: []density.Qubit{0},
					want: [][]complex128{
						{0.5, -0.25},
						{-0.25, 0.5},
					},
				},
				{
					idx: []density.Qubit{1},
					want: [][]complex128{
						{0.5, -0.25},
						{-0.25, 0.5},
					},
				},
			},
		},
		{
			s: []density.State{
				{
					Probability: 1.0,
					Qubit:       qubit.Zero(3),
				},
			},
			cs: []Case{
				{
					idx: []density.Qubit{0, 1},
					want: [][]complex128{
						{1, 0},
						{0, 0},
					},
				},
				{
					idx: []density.Qubit{1, 2},
					want: [][]complex128{
						{1, 0},
						{0, 0},
					},
				},
			},
		},
		{
			s: []density.State{
				{
					Probability: 1.0,
					Qubit:       qubit.One(3),
				},
			},
			cs: []Case{
				{
					idx: []density.Qubit{0, 1},
					want: [][]complex128{
						{0, 0},
						{0, 1},
					},
				},
				{
					idx: []density.Qubit{1, 2},
					want: [][]complex128{
						{0, 0},
						{0, 1},
					},
				},
			},
		},
		{
			s: []density.State{
				{
					Probability: 1.0,
					Qubit:       qubit.Zero(3).Apply(gate.H(3)),
				},
			},
			cs: []Case{
				{
					idx: []density.Qubit{0, 1},
					want: [][]complex128{
						{0.5, 0.5},
						{0.5, 0.5},
					},
				},
				{
					idx: []density.Qubit{1, 2},
					want: [][]complex128{
						{0.5, 0.5},
						{0.5, 0.5},
					},
				},
			},
		},
	}

	for _, c := range cases {
		for _, s := range c.cs {
			got := density.New(c.s).PartialTrace(s.idx...)

			p, q := got.Dim()
			if p != len(s.want) || q != len(s.want) {
				t.Errorf("got=%v, %v want=%v", p, q, s.want)
			}

			for i := range s.want {
				for j := range s.want[0] {
					if !epsilon.IsClose(got.At(i, j), s.want[i][j]) {
						t.Errorf("%v:%v, got=%v, want=%v", i, j, got.At(i, j), s.want[i][j])
					}
				}
			}

			if !epsilon.IsCloseF64(got.Trace(), 1) {
				t.Errorf("trace: got=%v, want=%v", got.Trace(), 1)
			}

			if got.Purity() > 1 && !epsilon.IsCloseF64(got.Purity(), 1) {
				t.Errorf("purity: got=%v, want<=%v", got.Purity(), 1)
			}
		}
	}
}

func TestApply(t *testing.T) {
	cases := []struct {
		s    []density.State
		u    *matrix.Matrix
		m    *qubit.Qubit
		want float64
	}{
		{
			s: []density.State{
				{
					Probability: 1,
					Qubit:       qubit.Zero(),
				},
			},
			u:    gate.X(),
			m:    qubit.Zero(),
			want: 0,
		},
		{
			s: []density.State{
				{
					Probability: 1,
					Qubit:       qubit.Zero(),
				},
			},
			u:    gate.X(),
			m:    qubit.One(),
			want: 1,
		},
	}

	for _, c := range cases {
		if density.New(c.s).Apply(c.u).Probability(c.m) != c.want {
			t.Fail()
		}
	}
}
