package density_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/itsubaki/q/math/epsilon"
	"github.com/itsubaki/q/math/matrix"
	"github.com/itsubaki/q/quantum/channel"
	"github.com/itsubaki/q/quantum/density"
	"github.com/itsubaki/q/quantum/gate"
	"github.com/itsubaki/q/quantum/qubit"
)

func Example_channel() {
	rho := density.From(qubit.One()).
		AmplitudeDamping(0.9, 0).
		BitFlip(0.5, 0)

	p, _ := rho.Measure(qubit.Zero())
	fmt.Printf("%.4f\n", p)

	// Output:
	// 0.5000
}

func Example_compose() {
	rho := density.From(qubit.One()).ApplyChannelFunc(channel.Compose(
		channel.AmplitudeDamping(0.9, 0),
		channel.BitFlip(0.5, 0),
	))

	p, _ := rho.Measure(qubit.Zero())
	fmt.Printf("%.4f\n", p)

	// Output:
	// 0.5000
}

func ExampleDensityMatrix_Matrix() {
	rho := density.FromStates([]density.WeightedState{
		{
			Probability: 0.1,
			Qubit:       qubit.Zero(),
		},
		{
			Probability: 0.9,
			Qubit:       qubit.One(),
		},
	})

	for _, r := range rho.Matrix().Seq2() {
		fmt.Println(r)
	}

	// Output:
	// [(0.1+0i) (0+0i)]
	// [(0+0i) (0.9+0i)]
}

func ExampleDensityMatrix_ExpectedValue() {
	rho := density.FromStates([]density.WeightedState{
		{
			Probability: 0.1,
			Qubit:       qubit.Zero(),
		},
		{
			Probability: 0.9,
			Qubit:       qubit.Plus(),
		},
	})

	fmt.Printf("X: %.2v\n", rho.ExpectedValue(gate.X()))
	fmt.Printf("Y: %.2v\n", rho.ExpectedValue(gate.Y()))
	fmt.Printf("Z: %.2v\n", rho.ExpectedValue(gate.Z()))

	// Output:
	// X: 0.9
	// Y: 0
	// Z: 0.1
}

func ExampleDensityMatrix_Measure() {
	rho := density.From(qubit.Zero(2).Apply(
		gate.H().TensorProduct(gate.I()),
		gate.CNOT(2, 0, 1),
	))

	for _, basis := range []*qubit.Qubit{
		qubit.From("00"),
		qubit.From("01"),
		qubit.From("10"),
		qubit.From("11"),
	} {
		p, sigma := rho.Measure(basis)
		fmt.Printf("%v: %.2f\n", basis.BinaryString(), p)

		for _, r := range sigma.Seq2() {
			fmt.Println(r)
		}
	}

	// Output:
	// 00: 0.50
	// [(1+0i) (0+0i) (0+0i) (0+0i)]
	// [(0+0i) (0+0i) (0+0i) (0+0i)]
	// [(0+0i) (0+0i) (0+0i) (0+0i)]
	// [(0+0i) (0+0i) (0+0i) (0+0i)]
	// 01: 0.00
	// [(0+0i) (0+0i) (0+0i) (0+0i)]
	// [(0+0i) (0+0i) (0+0i) (0+0i)]
	// [(0+0i) (0+0i) (0+0i) (0+0i)]
	// [(0+0i) (0+0i) (0+0i) (0+0i)]
	// 10: 0.00
	// [(0+0i) (0+0i) (0+0i) (0+0i)]
	// [(0+0i) (0+0i) (0+0i) (0+0i)]
	// [(0+0i) (0+0i) (0+0i) (0+0i)]
	// [(0+0i) (0+0i) (0+0i) (0+0i)]
	// 11: 0.50
	// [(0+0i) (0+0i) (0+0i) (0+0i)]
	// [(0+0i) (0+0i) (0+0i) (0+0i)]
	// [(0+0i) (0+0i) (0+0i) (0+0i)]
	// [(0+0i) (0+0i) (0+0i) (1+0i)]
}

func ExampleDensityMatrix_IsHermitian() {
	s0 := density.From(qubit.Zero())
	s1 := density.FromStates([]density.WeightedState{
		{
			Probability: 0.1,
			Qubit:       qubit.Zero(),
		},
		{
			Probability: 0.9,
			Qubit:       qubit.One(),
		},
	})

	fmt.Println(s0.IsHermitian())
	fmt.Println(s1.IsHermitian())

	// Output:
	// true
	// true
}

func ExampleDensityMatrix_Trace() {
	s0 := density.From(qubit.Zero())
	s1 := density.FromStates([]density.WeightedState{
		{
			Probability: 0.1,
			Qubit:       qubit.Zero(),
		},
		{
			Probability: 0.9,
			Qubit:       qubit.One(),
		},
	})

	fmt.Printf("pure : %.2f\n", s0.Trace())
	fmt.Printf("mixed: %.2f\n", s1.Trace())

	// Output:
	// pure : 1.00
	// mixed: 1.00
}

func ExampleDensityMatrix_Purity() {
	s0 := density.From(qubit.Zero())
	s1 := density.FromStates([]density.WeightedState{
		{
			Probability: 0.1,
			Qubit:       qubit.Zero(),
		},
		{
			Probability: 0.9,
			Qubit:       qubit.One(),
		},
	})

	fmt.Printf("pure : %.2f, %v\n", s0.Purity(), s0.IsPure())
	fmt.Printf("mixed: %.2f, %v\n", s1.Purity(), s1.IsMixed())

	// Output:
	// pure : 1.00, true
	// mixed: 0.82, true
}

func ExampleDensityMatrix_TensorProduct() {
	s0 := density.From(qubit.Zero())
	s1 := density.From(qubit.One())

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

func ExampleDensityMatrix_TraceOut() {
	rho := density.FromStates([]density.WeightedState{
		{
			Probability: 0.5,
			Qubit:       qubit.From("00"),
		},
		{
			Probability: 0.5,
			Qubit:       qubit.From("10"),
		},
	})

	s0 := rho.TraceOut(1) // trace out qubit 1
	s1 := rho.TraceOut(0) // trace out qubit 0

	fmt.Printf("trace: %.2v, purity: %.2v\n", rho.Trace(), rho.Purity())
	fmt.Printf("trace: %.2v, purity: %.2v\n", s0.Trace(), s0.Purity())
	fmt.Printf("trace: %.2v, purity: %.2v\n", s1.Trace(), s1.Purity())

	// Output:
	// trace: 1, purity: 0.5
	// trace: 1, purity: 0.5
	// trace: 1, purity: 1
}

func ExampleDensityMatrix_TraceOut_x8() {
	rho := density.From(qubit.Zero(3).Apply(
		matrix.TensorProduct(gate.H(), gate.I(), gate.I()),
		gate.CNOT(3, 0, 1),
	))

	s12 := rho.TraceOut(0)
	s02 := rho.TraceOut(1)
	s01 := rho.TraceOut(2)

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

func ExampleDensityMatrix_Depolarizing() {
	rho := density.From(qubit.Zero())
	p0, _ := rho.Measure(qubit.Zero())
	p1, _ := rho.Measure(qubit.One())
	fmt.Printf("0: %.2f\n", p0)
	fmt.Printf("1: %.2f\n", p1)
	fmt.Println()

	// XrhoX = |1><1|, YrhoY = |1><1|, ZrhoZ = |0><0|
	// E(rho) = 0.7|0><0| + 0.1|1><1| + 0.1|1><1| + 0.1|0><0| = 0.8|0><0| + 0.2|1><1|
	s0 := rho.Depolarizing(0.3, 0)
	p0, _ = s0.Measure(qubit.Zero())
	p1, _ = s0.Measure(qubit.One())
	fmt.Printf("0: %.2f\n", p0)
	fmt.Printf("1: %.2f\n", p1)

	// Output:
	// 0: 1.00
	// 1: 0.00
	//
	// 0: 0.80
	// 1: 0.20
}

func ExampleDensityMatrix_AmplitudeDamping() {
	rho := density.From(qubit.Plus())
	for _, p := range []float64{0.3, 0.6, 0.9, 1.0} {
		p0, _ := rho.AmplitudeDamping(p, 0).Measure(qubit.Zero())
		fmt.Printf("%.4f\n", p0)
	}

	// Output:
	// 0.6500
	// 0.8000
	// 0.9500
	// 1.0000
}

func ExampleDensityMatrix_Flip() {
	rho := density.From(qubit.Zero(2))
	s1 := rho.Flip(0.3, gate.X(), 0)

	p0, _ := s1.Measure(qubit.From("00"))
	p1, _ := s1.Measure(qubit.From("10"))
	fmt.Printf("%.2f\n", p0)
	fmt.Printf("%.2f\n", p1)

	// Output:
	// 0.70
	// 0.30
}

func ExampleDensityMatrix_Flip_qb1() {
	rho := density.From(qubit.Zero(2))
	s1 := rho.Flip(0.3, gate.X(), 1)

	p0, _ := s1.Measure(qubit.From("00"))
	p1, _ := s1.Measure(qubit.From("01"))
	fmt.Printf("%.2f\n", p0)
	fmt.Printf("%.2f\n", p1)

	// Output:
	// 0.70
	// 0.30
}

func ExampleDensityMatrix_BitFlip() {
	rho := density.From(qubit.Zero())
	x := rho.BitFlip(0.3, 0)

	p0, _ := x.Measure(qubit.Zero())
	p1, _ := x.Measure(qubit.One())
	fmt.Printf("%.2f\n", p0)
	fmt.Printf("%.2f\n", p1)

	// Output:
	// 0.70
	// 0.30
}

func ExampleDensityMatrix_BitPhaseFlip() {
	rho := density.From(qubit.Plus())
	y := rho.BitPhaseFlip(0.3, 0)

	p0, _ := y.Measure(qubit.Plus())
	p1, _ := y.Measure(qubit.Minus())
	fmt.Printf("%.2f\n", p0)
	fmt.Printf("%.2f\n", p1)

	// Output:
	// 0.70
	// 0.30
}

func ExampleDensityMatrix_PhaseFlip() {
	rho := density.From(qubit.Plus())
	z := rho.PhaseFlip(0.3, 0)

	p0, _ := z.Measure(qubit.Plus())
	p1, _ := z.Measure(qubit.Minus())
	fmt.Printf("%.2f\n", p0)
	fmt.Printf("%.2f\n", p1)

	// Output:
	// 0.70
	// 0.30
}

func TestFromStates(t *testing.T) {
	cases := []struct {
		s    []density.WeightedState
		want float64
	}{
		{
			s:    []density.WeightedState{},
			want: 0,
		},
	}

	for _, c := range cases {
		rho := density.FromStates(c.s)
		prop, _ := rho.Measure(qubit.Zero())

		if !epsilon.IsCloseF64(prop, c.want) {
			t.Errorf("got=%v, want=%v", prop, c.want)
		}
	}

}

func TestExpectedValue(t *testing.T) {
	cases := []struct {
		s        []density.WeightedState
		tr, sqtr float64
		m        *matrix.Matrix
		v        float64
	}{
		{
			s: []density.WeightedState{
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
			s: []density.WeightedState{
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
		rho := density.FromStates(c.s)
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
		s    []density.WeightedState
		m    *qubit.Qubit
		want float64
	}{
		{
			s: []density.WeightedState{
				{
					Probability: 1,
					Qubit:       qubit.Zero(),
				},
			},
			m:    qubit.Zero(),
			want: 1,
		},
		{
			s: []density.WeightedState{
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
		got, _ := density.FromStates(c.s).Measure(c.m)
		if !epsilon.IsCloseF64(got, c.want) {
			t.Fail()
		}
	}
}

func TestApply(t *testing.T) {
	cases := []struct {
		s    []density.WeightedState
		u    *matrix.Matrix
		m    *qubit.Qubit
		want float64
	}{
		{
			s: []density.WeightedState{
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
			s: []density.WeightedState{
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
		got, _ := density.FromStates(c.s).Apply(c.u).Measure(c.m)
		if !epsilon.IsCloseF64(got, c.want) {
			t.Fail()
		}
	}
}

func TestNewMixedState(t *testing.T) {
	cases := []struct {
		s    []density.WeightedState
		want *matrix.Matrix
	}{
		{
			s: []density.WeightedState{
				{
					Probability: 1,
					Qubit:       qubit.Zero(),
				},
			},
			want: matrix.New([][]complex128{
				{1, 0},
				{0, 0},
			}...),
		},
	}

	for _, c := range cases {
		rho := density.FromStates(c.s)
		if !rho.Matrix().Equal(c.want) {
			t.Errorf("got=%v, want=%v", rho.Matrix(), c.want)
		}
	}
}

func TestTraceOut(t *testing.T) {
	type Case struct {
		qb   []int
		want [][]complex128
	}

	cases := []struct {
		s  []density.WeightedState
		cs []Case
	}{
		{
			s: []density.WeightedState{
				{
					Probability: 1.0,
					Qubit:       qubit.Zero(2),
				},
			},
			cs: []Case{
				{
					qb: []int{0},
					want: [][]complex128{
						{1, 0},
						{0, 0},
					},
				},
				{
					qb: []int{1},
					want: [][]complex128{
						{1, 0},
						{0, 0},
					},
				},
			},
		},
		{
			s: []density.WeightedState{
				{
					Probability: 1.0,
					Qubit:       qubit.One(2),
				},
			},
			cs: []Case{
				{
					qb: []int{0},
					want: [][]complex128{
						{0, 0},
						{0, 1},
					},
				},
				{
					qb: []int{1},
					want: [][]complex128{
						{0, 0},
						{0, 1},
					},
				},
			},
		},
		{
			s: []density.WeightedState{
				{
					Probability: 1.0,
					Qubit:       qubit.Zero(2).Apply(gate.H(2)),
				},
			},
			cs: []Case{
				{
					qb: []int{0},
					want: [][]complex128{
						{0.5, 0.5},
						{0.5, 0.5},
					},
				},
				{
					qb: []int{1},
					want: [][]complex128{
						{0.5, 0.5},
						{0.5, 0.5},
					},
				},
			},
		},
		{
			s: []density.WeightedState{
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
					qb: []int{0},
					want: [][]complex128{
						{0.5, 0},
						{0, 0.5},
					},
				},
				{
					qb: []int{1},
					want: [][]complex128{
						{0.5, 0},
						{0, 0.5},
					},
				},
			},
		},
		{
			s: []density.WeightedState{
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
					qb: []int{0},
					want: [][]complex128{
						{0.25, 0.25},
						{0.25, 0.75},
					},
				},
				{
					qb: []int{1},
					want: [][]complex128{
						{0.25, 0.25},
						{0.25, 0.75},
					},
				},
			},
		},
		{
			s: []density.WeightedState{
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
					qb: []int{0},
					want: [][]complex128{
						{0.5, 0.25},
						{0.25, 0.5},
					},
				},
				{
					qb: []int{1},
					want: [][]complex128{
						{0.5, 0.25},
						{0.25, 0.5},
					},
				},
			},
		},
		{
			s: []density.WeightedState{
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
					qb: []int{0},
					want: [][]complex128{
						{0.5, -0.25},
						{-0.25, 0.5},
					},
				},
				{
					qb: []int{1},
					want: [][]complex128{
						{0.5, -0.25},
						{-0.25, 0.5},
					},
				},
			},
		},
		{
			s: []density.WeightedState{
				{
					Probability: 1.0,
					Qubit:       qubit.Zero(3),
				},
			},
			cs: []Case{
				{
					qb: []int{0, 1},
					want: [][]complex128{
						{1, 0},
						{0, 0},
					},
				},
				{
					qb: []int{1, 2},
					want: [][]complex128{
						{1, 0},
						{0, 0},
					},
				},
			},
		},
		{
			s: []density.WeightedState{
				{
					Probability: 1.0,
					Qubit:       qubit.One(3),
				},
			},
			cs: []Case{
				{
					qb: []int{0, 1},
					want: [][]complex128{
						{0, 0},
						{0, 1},
					},
				},
				{
					qb: []int{1, 2},
					want: [][]complex128{
						{0, 0},
						{0, 1},
					},
				},
			},
		},
		{
			s: []density.WeightedState{
				{
					Probability: 1.0,
					Qubit:       qubit.Zero(3).Apply(gate.H(3)),
				},
			},
			cs: []Case{
				{
					qb: []int{0, 1},
					want: [][]complex128{
						{0.5, 0.5},
						{0.5, 0.5},
					},
				},
				{
					qb: []int{1, 2},
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
			rho := density.FromStates(c.s)

			got := rho.TraceOut(s.qb...)
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

func TestAmplitudeDamping(t *testing.T) {
	cases := []struct {
		qubit *qubit.Qubit
		p     float64
		m     *qubit.Qubit
		want  float64
	}{
		{
			qubit: qubit.One(),
			p:     0.3,
			m:     qubit.Zero(),
			want:  0.3,
		},
		{
			qubit: qubit.One(),
			p:     0.3,
			m:     qubit.One(),
			want:  0.7,
		},
		{
			qubit: qubit.Zero(),
			p:     0.3,
			m:     qubit.Zero(),
			want:  1.0,
		},
		{
			qubit: qubit.Zero(),
			p:     0.3,
			m:     qubit.One(),
			want:  0.0,
		},
		{
			qubit: qubit.Plus(),
			p:     0.3,
			m:     qubit.Zero(),
			want:  0.5 + 0.3*0.5,
		},
		{
			qubit: qubit.Plus(),
			p:     0.3,
			m:     qubit.One(),
			want:  0.5 * (1 - 0.3),
		},
		{
			qubit: qubit.Plus(),
			p:     0.3,
			m:     qubit.Plus(),
			want:  0.5 + 0.5*math.Sqrt(1-0.3),
		},
		{
			qubit: qubit.One(),
			p:     0.0,
			m:     qubit.One(),
			want:  1.0,
		},
		{
			qubit: qubit.One(),
			p:     1.0,
			m:     qubit.Zero(),
			want:  1.0,
		},
		{
			qubit: qubit.Plus(),
			p:     1.0,
			m:     qubit.Zero(),
			want:  1.0,
		},
	}

	for _, c := range cases {
		rho := density.FromStates([]density.WeightedState{
			{
				Probability: 1,
				Qubit:       c.qubit,
			},
		})

		got, _ := rho.AmplitudeDamping(c.p, 0).Measure(c.m)
		if !epsilon.IsCloseF64(got, c.want) {
			t.Fail()
		}
	}
}

func TestPhaseDamping(t *testing.T) {
	cases := []struct {
		qubit *qubit.Qubit
		g     float64
		m     *qubit.Qubit
		want  float64
	}{
		{
			qubit: qubit.Zero(),
			g:     0.3,
			m:     qubit.Zero(),
			want:  1.0,
		},
		{
			qubit: qubit.One(),
			g:     0.3,
			m:     qubit.One(),
			want:  1.0,
		},
		{
			qubit: qubit.Plus(),
			g:     0.3,
			m:     qubit.Zero(),
			want:  0.5,
		},
		{
			qubit: qubit.Plus(),
			g:     0.3,
			m:     qubit.One(),
			want:  0.5,
		},
		{
			qubit: qubit.Plus(),
			g:     0.3,
			m:     qubit.Plus(),
			want:  0.5 + 0.5*math.Sqrt(1-0.3),
		},
		{
			qubit: qubit.Plus(),
			g:     0.3,
			m:     qubit.Minus(),
			want:  0.5 - 0.5*math.Sqrt(1-0.3),
		},
		{
			qubit: qubit.Plus(),
			g:     0.0,
			m:     qubit.Plus(),
			want:  1.0,
		},
		{
			qubit: qubit.Plus(),
			g:     1.0,
			m:     qubit.Plus(),
			want:  0.5,
		},
	}

	for _, c := range cases {
		rho := density.FromStates([]density.WeightedState{
			{
				Probability: 1,
				Qubit:       c.qubit,
			},
		})

		got, _ := rho.PhaseDamping(c.g, 0).Measure(c.m)
		if !epsilon.IsCloseF64(got, c.want) {
			t.Errorf("got %v, want %v", got, c.want)
		}
	}
}

func TestPauli(t *testing.T) {
	cases := []struct {
		qubit *qubit.Qubit
		px    float64
		py    float64
		pz    float64
		m     *qubit.Qubit
		want  float64
	}{
		{
			qubit: qubit.Zero(),
			px:    0,
			py:    0,
			pz:    0,
			m:     qubit.Zero(),
			want:  1.0,
		},
		{
			qubit: qubit.Zero(),
			px:    1.0,
			py:    0,
			pz:    0,
			m:     qubit.One(),
			want:  1.0,
		},
		{
			qubit: qubit.Zero(),
			px:    0.5,
			py:    0,
			pz:    0,
			m:     qubit.Zero(),
			want:  0.5,
		},
		{
			qubit: qubit.Zero(),
			px:    0.5,
			py:    0,
			pz:    0,
			m:     qubit.One(),
			want:  0.5,
		},
		{
			qubit: qubit.Zero(),
			px:    0,
			py:    0,
			pz:    1.0,
			m:     qubit.Zero(),
			want:  1.0,
		},
		{
			qubit: qubit.Plus(),
			px:    0,
			py:    0,
			pz:    1.0,
			m:     qubit.Minus(),
			want:  1.0,
		},
		{
			qubit: qubit.Zero(),
			px:    0,
			py:    1.0,
			pz:    0,
			m:     qubit.One(),
			want:  1.0,
		},
		{
			qubit: qubit.Zero(),
			px:    1.0 / 3,
			py:    1.0 / 3,
			pz:    1.0 / 3,
			m:     qubit.Zero(),
			want:  1.0 / 3,
		},
		{
			qubit: qubit.Zero(),
			px:    1.0 / 3,
			py:    1.0 / 3,
			pz:    1.0 / 3,
			m:     qubit.One(),
			want:  1.0/3 + 1.0/3,
		},
	}

	for _, c := range cases {
		rho := density.FromStates([]density.WeightedState{
			{
				Probability: 1,
				Qubit:       c.qubit,
			},
		})

		got, _ := rho.Pauli(c.px, c.py, c.pz, 0).Measure(c.m)
		if !epsilon.IsCloseF64(got, c.want) {
			t.Errorf("got %v, want %v", got, c.want)
		}
	}
}

func TestDensityMatrix_ApplyChannelFunc(t *testing.T) {
	cases := []struct {
		channelFunc []channel.ChannelFunc
		want        float64
	}{
		{
			channelFunc: []channel.ChannelFunc{},
			want:        1.0,
		},
		{
			channelFunc: []channel.ChannelFunc{
				channel.Depolarizing(0.1, 0),
				channel.AmplitudeDamping(0.7, 0),
				channel.PhaseDamping(0.7, 0),
				channel.BitFlip(0.1, 0),
			},
			want: 0.884,
		},
	}

	for _, c := range cases {
		rho := density.From(qubit.Zero())
		s0 := rho.ApplyChannelFunc(c.channelFunc...)

		got, _ := s0.Measure(qubit.Zero())
		if !epsilon.IsCloseF64(got, c.want) {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}

func TestDensityMatrix_ApplyKraus(t *testing.T) {
	cases := []struct {
		kraus []*matrix.Matrix
		want  float64
	}{
		{
			kraus: []*matrix.Matrix{},
			want:  1.0,
		},
		{
			kraus: channel.Depolarizing(0.1, 0)(1).Kraus,
			want:  0.9333333333333332,
		},
	}

	for _, c := range cases {
		rho := density.From(qubit.Zero())
		s0 := rho.ApplyKraus(c.kraus...)

		got, _ := s0.Measure(qubit.Zero())
		if !epsilon.IsCloseF64(got, c.want) {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}
