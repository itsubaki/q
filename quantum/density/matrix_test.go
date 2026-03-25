package density_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/itsubaki/q/math/epsilon"
	"github.com/itsubaki/q/math/matrix"
	"github.com/itsubaki/q/quantum/density"
	"github.com/itsubaki/q/quantum/gate"
	"github.com/itsubaki/q/quantum/qubit"
)

func Example_bell() {
	rho, err := density.NewPureState(qubit.Zero(2).Apply(
		gate.H().TensorProduct(gate.I()),
		gate.CNOT(2, 0, 1),
	))
	if err != nil {
		panic(err)
	}

	s0 := rho.TraceOut(1)
	s1 := rho.TraceOut(0)

	fmt.Printf("trace: %.2v, purity: %.2v\n", rho.Trace(), rho.Purity())
	fmt.Printf("trace: %.2v, purity: %.2v\n", s0.Trace(), s0.Purity())
	fmt.Printf("trace: %.2v, purity: %.2v\n", s1.Trace(), s1.Purity())

	// Output:
	// trace: 1, purity: 1
	// trace: 1, purity: 0.5
	// trace: 1, purity: 0.5
}

func ExampleMatrix_Matrix() {
	rho, err := density.New([]density.State{
		{0.1, qubit.Zero()},
		{0.9, qubit.One()},
	})
	if err != nil {
		panic(err)
	}

	for _, r := range rho.Matrix().Seq2() {
		fmt.Println(r)
	}

	// Output:
	// [(0.1+0i) (0+0i)]
	// [(0+0i) (0.9+0i)]
}

func ExampleMatrix_ExpectedValue() {
	rho, err := density.New([]density.State{
		{0.1, qubit.Zero()},
		{0.9, qubit.Plus()},
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("X: %.2v\n", rho.ExpectedValue(gate.X()))
	fmt.Printf("Y: %.2v\n", rho.ExpectedValue(gate.Y()))
	fmt.Printf("Z: %.2v\n", rho.ExpectedValue(gate.Z()))

	// Output:
	// X: 0.9
	// Y: 0
	// Z: 0.1
}

func ExampleMatrix_Probability() {
	rho, err := density.New([]density.State{
		{0.1, qubit.Zero()},
		{0.9, qubit.One()},
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("0: %.2v\n", rho.Probability(qubit.Zero()))
	fmt.Printf("1: %.2v\n", rho.Probability(qubit.One()))

	// Output:
	// 0: 0.1
	// 1: 0.9
}

func ExampleMatrix_IsHermite() {
	s0, err := density.NewPureState(qubit.Zero())
	if err != nil {
		panic(err)
	}

	s1, err := density.New([]density.State{
		{0.1, qubit.Zero()},
		{0.9, qubit.One()},
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(s0.IsHermite())
	fmt.Println(s1.IsHermite())

	// Output:
	// true
	// true
}

func ExampleMatrix_Trace() {
	s0, err := density.NewPureState(qubit.Zero())
	if err != nil {
		panic(err)
	}

	s1, err := density.New([]density.State{
		{0.1, qubit.Zero()},
		{0.9, qubit.One()},
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("pure : %.2f\n", s0.Trace())
	fmt.Printf("mixed: %.2f\n", s1.Trace())

	// Output:
	// pure : 1.00
	// mixed: 1.00
}

func ExampleMatrix_Purity() {
	s0, err := density.NewPureState(qubit.Zero())
	if err != nil {
		panic(err)
	}

	s1, err := density.New([]density.State{
		{0.1, qubit.Zero()},
		{0.9, qubit.One()},
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("pure : %.2f, %v\n", s0.Purity(), s0.IsPure())
	fmt.Printf("mixed: %.2f, %v\n", s1.Purity(), s1.IsMixed())

	// Output:
	// pure : 1.00, true
	// mixed: 0.82, true
}

func ExampleMatrix_TensorProduct() {
	s0, err := density.NewPureState(qubit.Zero())
	if err != nil {
		panic(err)
	}

	s1, err := density.NewPureState(qubit.One())
	if err != nil {
		panic(err)
	}

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

func ExampleMatrix_Project() {
	rho, err := density.NewPureState(qubit.Zero(2).Apply(
		gate.H().TensorProduct(gate.I()),
		gate.CNOT(2, 0, 1),
	))
	if err != nil {
		panic(err)
	}

	computationalBasis := []*qubit.Qubit{
		qubit.From("00"),
		qubit.From("01"),
		qubit.From("10"),
		qubit.From("11"),
	}

	for _, basis := range computationalBasis {
		p, sigma := rho.Project(basis)

		fmt.Printf("%v: %.2f\n", basis.State(), p)
		for _, r := range sigma.Seq2() {
			fmt.Println(r)
		}
	}

	// Output:
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

func ExampleMatrix_TraceOut() {
	rho, err := density.New([]density.State{
		{0.5, qubit.From("00")},
		{0.5, qubit.From("10")},
	})
	if err != nil {
		panic(err)
	}

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

func ExampleMatrix_PartialTrace_x8() {
	rho, err := density.NewPureState(qubit.Zero(3).Apply(
		matrix.TensorProduct(gate.H(), gate.I(), gate.I()),
		gate.CNOT(3, 0, 1),
	))
	if err != nil {
		panic(err)
	}

	s12 := rho.PartialTrace(0) // trace out qubit 0
	s02 := rho.PartialTrace(1) // trace out qubit 1
	s01 := rho.PartialTrace(2) // trace out qubit 2

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
	rho, err := density.NewPureState(qubit.Zero())
	if err != nil {
		panic(err)
	}

	fmt.Printf("0: %.2f\n", rho.Probability(qubit.Zero()))
	fmt.Printf("1: %.2f\n", rho.Probability(qubit.One()))
	fmt.Println()

	// XrhoX = |1><1|, YrhoY = |1><1|, ZrhoZ = |0><0|
	// E(rho) = 0.7|0><0| + 0.1|1><1| + 0.1|1><1| + 0.1|0><0| = 0.8|0><0| + 0.2|1><1|
	s0 := rho.Depolarizing(0.3, 0)
	fmt.Printf("0: %.2f\n", s0.Probability(qubit.Zero()))
	fmt.Printf("1: %.2f\n", s0.Probability(qubit.One()))

	// Output:
	// 0: 1.00
	// 1: 0.00
	//
	// 0: 0.80
	// 1: 0.20
}

func ExampleMatrix_FlipChannel() {
	rho, err := density.NewPureState(qubit.Zero(2))
	if err != nil {
		panic(err)
	}

	s1 := rho.FlipChannel(0.3, gate.X(), 0)
	fmt.Printf("%.2f\n", s1.Probability(qubit.From("00")))
	fmt.Printf("%.2f\n", s1.Probability(qubit.From("10")))

	// Output:
	// 0.70
	// 0.30
}

func ExampleMatrix_FlipChannel_qb1() {
	rho, err := density.NewPureState(qubit.Zero(2))
	if err != nil {
		panic(err)
	}

	s1 := rho.FlipChannel(0.3, gate.X(), 1)
	fmt.Printf("%.2f\n", s1.Probability(qubit.From("00")))
	fmt.Printf("%.2f\n", s1.Probability(qubit.From("01")))

	// Output:
	// 0.70
	// 0.30
}

func ExampleMatrix_BitFlip() {
	rho, err := density.NewPureState(qubit.Zero())
	if err != nil {
		panic(err)
	}

	x := rho.BitFlip(0.3, 0)
	fmt.Printf("%.2f\n", x.Probability(qubit.Zero()))
	fmt.Printf("%.2f\n", x.Probability(qubit.One()))

	// Output:
	// 0.70
	// 0.30
}

func ExampleMatrix_BitPhaseFlip() {
	rho, err := density.NewPureState(qubit.Plus())
	if err != nil {
		panic(err)
	}

	y := rho.BitPhaseFlip(0.3, 0)
	fmt.Printf("%.2f\n", y.Probability(qubit.Plus()))
	fmt.Printf("%.2f\n", y.Probability(qubit.Minus()))

	// Output:
	// 0.70
	// 0.30
}

func ExampleMatrix_PhaseFlip() {
	rho, err := density.NewPureState(qubit.Plus())
	if err != nil {
		panic(err)
	}

	z := rho.PhaseFlip(0.3, 0)
	fmt.Printf("%.2f\n", z.Probability(qubit.Plus()))
	fmt.Printf("%.2f\n", z.Probability(qubit.Minus()))

	// Output:
	// 0.70
	// 0.30
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
		rho, err := density.New(c.s)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

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
		rho, err := density.New(c.s)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		got := rho.Probability(c.m)
		if !epsilon.IsCloseF64(got, c.want) {
			t.Fail()
		}
	}
}

func TestPartialTrace(t *testing.T) {
	type Case struct {
		qb   []int
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
			s: []density.State{
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
			s: []density.State{
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
			s: []density.State{
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
			s: []density.State{
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
			s: []density.State{
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
			rho, err := density.New(c.s)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				continue
			}

			got := rho.PartialTrace(s.qb...)

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
		rho, err := density.New(c.s)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		got := rho.Apply(c.u).Probability(c.m)
		if !epsilon.IsCloseF64(got, c.want) {
			t.Fail()
		}
	}
}

func TestNew(t *testing.T) {
	type Case struct {
		s      []density.State
		want   *matrix.Matrix
		hasErr bool
	}

	cases := []Case{
		{
			s:      []density.State{},
			hasErr: true,
		},
		{
			s: []density.State{
				{
					Probability: 1,
					Qubit:       qubit.Zero(),
				},
			},
			want: matrix.New([][]complex128{
				{1, 0},
				{0, 0},
			}...),
			hasErr: false,
		},
	}

	for _, c := range cases {
		rho, err := density.New(c.s)
		if err != nil {
			if c.hasErr {
				continue
			}

			t.Errorf("unexpected error: %v", err)
			continue
		}

		if !rho.Matrix().Equal(c.want) {
			t.Errorf("got=%v, want=%v", rho.Matrix(), c.want)
		}
	}
}

func TestAmplitudeDamping(t *testing.T) {
	type Case struct {
		s    *qubit.Qubit
		p    float64
		m    *qubit.Qubit
		want float64
	}

	cases := []Case{
		{
			s:    qubit.One(),
			p:    0.3,
			m:    qubit.Zero(),
			want: 0.3,
		},
		{
			s:    qubit.One(),
			p:    0.3,
			m:    qubit.One(),
			want: 0.7,
		},
		{
			s:    qubit.Zero(),
			p:    0.3,
			m:    qubit.Zero(),
			want: 1.0,
		},
		{
			s:    qubit.Zero(),
			p:    0.3,
			m:    qubit.One(),
			want: 0.0,
		},
		{
			s:    qubit.Plus(),
			p:    0.3,
			m:    qubit.Zero(),
			want: 0.5 + 0.3*0.5,
		},
		{
			s:    qubit.Plus(),
			p:    0.3,
			m:    qubit.One(),
			want: 0.5 * (1 - 0.3),
		},
		{
			s:    qubit.Plus(),
			p:    0.3,
			m:    qubit.Plus(),
			want: 0.5 + 0.5*math.Sqrt(1-0.3),
		},
		{
			s:    qubit.One(),
			p:    0.0,
			m:    qubit.One(),
			want: 1.0,
		},
		{
			s:    qubit.One(),
			p:    1.0,
			m:    qubit.Zero(),
			want: 1.0,
		},
		{
			s:    qubit.Plus(),
			p:    1.0,
			m:    qubit.Zero(),
			want: 1.0,
		},
	}

	for _, c := range cases {
		rho, err := density.New([]density.State{
			{
				Probability: 1,
				Qubit:       c.s,
			},
		})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		got := rho.AmplitudeDamping(c.p, 0)
		if !epsilon.IsCloseF64(got.Probability(c.m), c.want) {
			t.Fail()
		}
	}
}

func TestPhaseDamping(t *testing.T) {
	type Case struct {
		s    *qubit.Qubit
		g    float64
		m    *qubit.Qubit
		want float64
	}

	cases := []Case{
		{
			s:    qubit.Zero(),
			g:    0.3,
			m:    qubit.Zero(),
			want: 1.0,
		},
		{
			s:    qubit.One(),
			g:    0.3,
			m:    qubit.One(),
			want: 1.0,
		},
		{
			s:    qubit.Plus(),
			g:    0.3,
			m:    qubit.Zero(),
			want: 0.5,
		},
		{
			s:    qubit.Plus(),
			g:    0.3,
			m:    qubit.One(),
			want: 0.5,
		},
		{
			s:    qubit.Plus(),
			g:    0.3,
			m:    qubit.Plus(),
			want: 0.5 + 0.5*math.Sqrt(1-0.3),
		},
		{
			s:    qubit.Plus(),
			g:    0.3,
			m:    qubit.Minus(),
			want: 0.5 - 0.5*math.Sqrt(1-0.3),
		},
		{
			s:    qubit.Plus(),
			g:    0.0,
			m:    qubit.Plus(),
			want: 1.0,
		},
		{
			s:    qubit.Plus(),
			g:    1.0,
			m:    qubit.Plus(),
			want: 0.5,
		},
	}

	for _, c := range cases {
		rho, err := density.New([]density.State{
			{
				Probability: 1,
				Qubit:       c.s,
			},
		})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		got := rho.PhaseDamping(c.g, 0)
		if !epsilon.IsCloseF64(got.Probability(c.m), c.want) {
			t.Errorf("got %v, want %v", got.Probability(c.m), c.want)
		}
	}
}
