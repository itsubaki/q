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
	"github.com/itsubaki/q/quantum/observable"
	"github.com/itsubaki/q/quantum/qubit"
)

func Example_channel() {
	p, post := density.New(qubit.One()).
		AmplitudeDamping(0.9).
		BitFlip(0.5).
		Measure(observable.Projector(qubit.Zero()))

	fmt.Printf("%.4f\n", p)
	for _, r := range post.Seq2() {
		fmt.Println(r)
	}

	// Output:
	// 0.5000
	// [(1+0i) (0+0i)]
	// [(0+0i) (0+0i)]
}

func Example_composeFunc() {
	composed := channel.ComposeFunc(
		channel.AmplitudeDamping(0.9, 0),
		channel.BitFlip(0.5, 1),
	)

	rho := density.New(qubit.Ones(2)).
		ApplyChannelFunc(composed)

	p, _ := rho.Measure(observable.Projector(
		qubit.Zeros(2),
	))
	fmt.Printf("%.4f\n", p)

	// Output:
	// 0.4500
}

func Example_nonCommutative() {
	ch1 := channel.AmplitudeDamping(0.9, 0)
	ch2 := channel.BitFlip(0.5, 0)

	rhoA := density.New(qubit.One()).
		ApplyChannelFunc(ch1).
		ApplyChannelFunc(ch2)

	rhoB := density.New(qubit.One()).
		ApplyChannelFunc(ch2).
		ApplyChannelFunc(ch1)

	fmt.Println(rhoA.Equal(rhoB))

	// Output:
	// false
}

func Example_classical() {
	rho := density.NewMixed([]density.WeightedState{
		{Probability: 0.5, Qubit: qubit.Zeros(3)},
		{Probability: 0.5, Qubit: qubit.Ones(3)},
	})

	for _, ob := range []*matrix.Matrix{
		observable.Pauli("XXX"), // 0
		observable.Pauli("ZZZ"), // 0
		observable.Pauli("ZZI"), // 1
		observable.Pauli("ZIZ"), // 1
		observable.Pauli("IZZ"), // 1
	} {
		fmt.Println(rho.Expect(ob))
	}

	// Output:
	// 0
	// 0
	// 1
	// 1
	// 1
}

func ExampleDensityMatrix_Expect() {
	rho := density.NewMixed([]density.WeightedState{
		{Probability: 0.1, Qubit: qubit.Zero()},
		{Probability: 0.9, Qubit: qubit.Plus()},
	})

	fmt.Printf("X: %.2v\n", rho.Expect(observable.X()))
	fmt.Printf("Y: %.2v\n", rho.Expect(observable.Y()))
	fmt.Printf("Z: %.2v\n", rho.Expect(observable.Z()))

	// Output:
	// X: 0.9
	// Y: 0
	// Z: 0.1
}

func ExampleDensityMatrix_Measure() {
	rho := density.New(qubit.Zeros(2).Apply(
		matrix.TensorProduct(gate.H(), gate.I()),
		gate.CNOT(2, 0, 1),
	))

	for _, qb := range []*qubit.Qubit{
		qubit.From("00"),
		qubit.From("01"),
		qubit.From("10"),
		qubit.From("11"),
	} {
		p, post := rho.Measure(observable.Projector(qb))

		fmt.Printf("%v: %.2f\n", qb.BinaryString(), p)
		for _, r := range post.Seq2() {
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

func ExampleDensityMatrix_TensorProduct() {
	s0 := density.New(qubit.Zero())
	s1 := density.New(qubit.One())

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
	rho := density.NewMixed([]density.WeightedState{
		{Probability: 0.5, Qubit: qubit.From("00")},
		{Probability: 0.5, Qubit: qubit.From("10")},
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
	rho := density.New(qubit.Zeros(3).Apply(
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
	rho := density.New(qubit.Zero())
	s := rho.Depolarizing(0.3)

	for _, qb := range []*qubit.Qubit{
		qubit.Zero(),
		qubit.One(),
	} {
		p, _ := s.Measure(observable.Projector(qb))
		fmt.Printf("%.2f\n", p)
	}

	// Output:
	// 0.80
	// 0.20
}

func ExampleDensityMatrix_BitFlip() {
	rho := density.New(qubit.Zero())
	s := rho.BitFlip(0.3)

	for _, qb := range []*qubit.Qubit{
		qubit.Zero(),
		qubit.One(),
	} {
		p, _ := s.Measure(observable.Projector(qb))
		fmt.Printf("%.2f\n", p)
	}

	// Output:
	// 0.70
	// 0.30
}

func ExampleDensityMatrix_BitPhaseFlip() {
	rho := density.New(qubit.Plus())
	s := rho.BitPhaseFlip(0.3)

	for _, qb := range []*qubit.Qubit{
		qubit.Plus(),
		qubit.Minus(),
	} {
		p, _ := s.Measure(observable.Projector(qb))
		fmt.Printf("%.2f\n", p)
	}

	// Output:
	// 0.70
	// 0.30
}

func ExampleDensityMatrix_PhaseFlip() {
	rho := density.New(qubit.Plus())
	s := rho.PhaseFlip(0.3)

	for _, qb := range []*qubit.Qubit{
		qubit.Plus(),
		qubit.Minus(),
	} {
		p, _ := s.Measure(observable.Projector(qb))
		fmt.Printf("%.2f\n", p)
	}

	// Output:
	// 0.70
	// 0.30
}

func ExampleDensityMatrix_VonNeumannEntropy() {
	rho := density.New(qubit.Zeros(2).Apply(
		matrix.TensorProduct(gate.H(), gate.I()),
		gate.CNOT(2, 0, 1),
	))

	s0 := rho.TraceOut(1)
	s1 := rho.TraceOut(0)

	fmt.Println(epsilon.IsZeroF64(rho.VonNeumannEntropy()))
	fmt.Println(epsilon.IsOneF64(s0.VonNeumannEntropy()))
	fmt.Println(epsilon.IsOneF64(s1.VonNeumannEntropy()))

	// Output:
	// true
	// true
	// true
}

func ExampleDensityMatrix_FidelitySquared() {
	rho := density.New(qubit.Zero())
	sigma := density.New(qubit.Plus())

	fmt.Printf("%.4f\n", rho.Fidelity(sigma))
	fmt.Printf("%.4f\n", rho.FidelitySquared(sigma))

	// Output:
	// 0.7071
	// 0.5000
}

func TestNewMixed(t *testing.T) {
	cases := []struct {
		s    []density.WeightedState
		want float64
	}{
		{
			s:    []density.WeightedState{},
			want: 0,
		},
		{
			s: []density.WeightedState{
				{Probability: 1, Qubit: qubit.Zero()},
			},
			want: 1,
		},
	}

	for _, c := range cases {
		p, _ := density.NewMixed(c.s).Measure(observable.Projector(
			qubit.Zero(),
		))

		if !epsilon.IsCloseF64(p, c.want) {
			t.Errorf("got=%v, want=%v", p, c.want)
		}
	}

}

func TestDensityMatrix_IsHermitian(t *testing.T) {
	cases := []struct {
		s    []density.WeightedState
		want bool
	}{
		{
			s: []density.WeightedState{
				{Probability: 1, Qubit: qubit.Zero()},
			},
			want: true,
		},
		{
			s: []density.WeightedState{
				{Probability: 0.1, Qubit: qubit.Zero()},
				{Probability: 0.9, Qubit: qubit.One()},
				{Probability: 0.9, Qubit: qubit.One()},
			},
			want: true,
		},
	}

	for _, c := range cases {
		got := density.NewMixed(c.s).IsHermitian()
		if got != c.want {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}

func TestDensityMatrix_Trace(t *testing.T) {
	cases := []struct {
		s    []density.WeightedState
		want float64
	}{
		{
			s: []density.WeightedState{
				{Probability: 1, Qubit: qubit.Zero()},
			},
			want: 1.0,
		},
		{
			s: []density.WeightedState{
				{Probability: 0.1, Qubit: qubit.Zero()},
				{Probability: 0.1, Qubit: qubit.Zero()},
				{Probability: 0.9, Qubit: qubit.One()},
			},
			want: 1.0,
		},
	}

	for _, c := range cases {
		got := density.NewMixed(c.s).Trace()
		if !epsilon.IsCloseF64(got, c.want) {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}

func TestDensityMatrix_Purity(t *testing.T) {
	cases := []struct {
		s       []density.WeightedState
		purity  float64
		isPure  bool
		isMixed bool
	}{
		{
			s: []density.WeightedState{
				{Probability: 1, Qubit: qubit.Zero()},
			},
			purity:  1.0,
			isPure:  true,
			isMixed: false,
		},
		{
			s: []density.WeightedState{
				{Probability: 0.1, Qubit: qubit.Zero()},
				{Probability: 0.9, Qubit: qubit.One()},
			},
			purity:  0.82,
			isPure:  false,
			isMixed: true,
		},
	}

	for _, c := range cases {
		rho := density.NewMixed(c.s)
		if !epsilon.IsCloseF64(rho.Purity(), c.purity) {
			t.Errorf("got=%v, want=%v", rho.Purity(), c.purity)
		}

		if rho.IsPure() != c.isPure {
			t.Errorf("got=%v, want=%v", rho.IsPure(), c.isPure)
		}

		if rho.IsMixed() != c.isMixed {
			t.Errorf("got=%v, want=%v", rho.IsMixed(), c.isMixed)
		}
	}
}

func TestDensityMatrix_Expect(t *testing.T) {
	cases := []struct {
		s        []density.WeightedState
		tr, sqtr float64
		m        *matrix.Matrix
		v        float64
	}{
		{
			s: []density.WeightedState{
				{Probability: 0.1, Qubit: qubit.Zero()},
				{Probability: 0.9, Qubit: qubit.One()},
			},
			tr:   1,
			sqtr: 0.82,
			m:    gate.X(),
			v:    0.0,
		},
		{
			s: []density.WeightedState{
				{Probability: 0.1, Qubit: qubit.Zero()},
				{Probability: 0.9, Qubit: qubit.Plus()},
			},
			tr:   1,
			sqtr: 0.91,
			m:    gate.X(),
			v:    0.9,
		},
	}

	for _, c := range cases {
		rho := density.NewMixed(c.s)
		if !epsilon.IsCloseF64(rho.Trace(), c.tr) {
			t.Errorf("trace=%v", rho.Trace())
		}

		if !epsilon.IsCloseF64(rho.Purity(), c.sqtr) {
			t.Errorf("purity=%v", rho.Purity())
		}

		if !epsilon.IsCloseF64(rho.Expect(c.m), c.v) {
			t.Errorf("expected_value=%v", rho.Expect(c.m))
		}
	}
}

func TestDensityMatrix_Apply(t *testing.T) {
	cases := []struct {
		s    []density.WeightedState
		u    *matrix.Matrix
		m    *qubit.Qubit
		want float64
	}{
		{
			s: []density.WeightedState{
				{Probability: 1, Qubit: qubit.Zero()},
			},
			u:    gate.X(),
			m:    qubit.Zero(),
			want: 0,
		},
		{
			s: []density.WeightedState{
				{Probability: 1, Qubit: qubit.Zero()},
			},
			u:    gate.X(),
			m:    qubit.One(),
			want: 1,
		},
	}

	for _, c := range cases {
		got, _ := density.NewMixed(c.s).
			Apply(c.u).
			Measure(observable.Projector(c.m))

		if !epsilon.IsCloseF64(got, c.want) {
			t.Fail()
		}
	}
}
func TestDensityMatrix_TraceOut(t *testing.T) {
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
				{Probability: 1.0, Qubit: qubit.Zeros(2)},
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
				{Probability: 1.0, Qubit: qubit.Ones(2)},
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
				{Probability: 1.0, Qubit: qubit.Zeros(2).Apply(gate.H(2))},
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
				{Probability: 0.5, Qubit: qubit.Zeros(2)},
				{Probability: 0.5, Qubit: qubit.Ones(2)},
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
				{Probability: 0.5, Qubit: qubit.Zeros(2).Apply(gate.H(2))},
				{Probability: 0.5, Qubit: qubit.Ones(2)},
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
				{Probability: 0.75, Qubit: qubit.Zeros(2).Apply(gate.H(2))},
				{Probability: 0.25, Qubit: qubit.Ones(2).Apply(gate.H(2))},
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
				{Probability: 0.25, Qubit: qubit.Zeros(2).Apply(gate.H(2))},
				{Probability: 0.75, Qubit: qubit.Ones(2).Apply(gate.H(2))},
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
				{Probability: 1.0, Qubit: qubit.Zeros(3)},
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
				{Probability: 1.0, Qubit: qubit.Ones(3)},
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
				{Probability: 1.0, Qubit: qubit.Zeros(3).Apply(gate.H(3))},
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
			got := density.NewMixed(c.s).TraceOut(s.qb...)

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
		got, _ := density.New(qubit.Zero()).
			ApplyChannelFunc(c.channelFunc...).
			Measure(observable.Projector(
				qubit.Zero(),
			))

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
		got, _ := density.New(qubit.Zero()).
			ApplyKraus(c.kraus...).
			Measure(observable.Projector(
				qubit.Zero(),
			))

		if !epsilon.IsCloseF64(got, c.want) {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}

func TestDensityMatrix_Flip(t *testing.T) {
	cases := []struct {
		s  *qubit.Qubit
		g  *matrix.Matrix
		p  float64
		qb int
		m  []struct {
			q *qubit.Qubit
			p float64
		}
	}{
		{
			s:  qubit.From("00"),
			g:  gate.X(),
			p:  0.3,
			qb: 0,
			m: []struct {
				q *qubit.Qubit
				p float64
			}{
				{
					q: qubit.From("00"),
					p: 0.7,
				},
				{
					q: qubit.From("10"),
					p: 0.3,
				},
			},
		},
		{
			s:  qubit.From("00"),
			g:  gate.X(),
			p:  0.3,
			qb: 1,
			m: []struct {
				q *qubit.Qubit
				p float64
			}{
				{
					q: qubit.From("00"),
					p: 0.7,
				},
				{
					q: qubit.From("01"),
					p: 0.3,
				},
			},
		},
	}

	for _, c := range cases {
		s := density.New(c.s).Flip(c.p, c.g, c.qb)

		for _, m := range c.m {
			p, _ := s.Measure(observable.Projector(
				m.q,
			))

			if !epsilon.IsCloseF64(p, m.p) {
				t.Errorf("got=%v, want=%v", p, m.p)
			}
		}
	}
}

func TestDensityMatrix_AmplitudeDamping(t *testing.T) {
	cases := []struct {
		q     *qubit.Qubit
		gamma float64
		qb    int
		m     *qubit.Qubit
		p     float64
	}{
		{
			q:     qubit.Plus(),
			gamma: 0.3,
			qb:    0,
			m:     qubit.Zero(),
			p:     0.65,
		},
		{
			q:     qubit.Plus(),
			gamma: 0.6,
			qb:    0,
			m:     qubit.Zero(),
			p:     0.8,
		},
		{
			q:     qubit.Plus(),
			gamma: 0.9,
			qb:    0,
			m:     qubit.Zero(),
			p:     0.95,
		},
		{
			q:     qubit.Plus(),
			gamma: 1.0,
			qb:    0,
			m:     qubit.Zero(),
			p:     1.0,
		},
	}

	for _, c := range cases {
		p, _ := density.New(c.q).
			AmplitudeDamping(c.gamma, c.qb).
			Measure(observable.Projector(c.m))

		if !epsilon.IsCloseF64(p, c.p) {
			t.Errorf("got=%v, want=%v", p, c.p)
		}
	}
}

func TestDensityMatrix_PhaseDamping(t *testing.T) {
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
		got, _ := density.New(c.qubit).
			PhaseDamping(c.g, 0).
			Measure(observable.Projector(
				c.m,
			))

		if !epsilon.IsCloseF64(got, c.want) {
			t.Errorf("got %v, want %v", got, c.want)
		}
	}
}

func TestDensityMatrix_Pauli(t *testing.T) {
	cases := []struct {
		qubit *qubit.Qubit
		pX    float64
		pY    float64
		pZ    float64
		m     *qubit.Qubit
		want  float64
	}{
		{
			qubit: qubit.Zero(),
			pX:    0,
			pY:    0,
			pZ:    0,
			m:     qubit.Zero(),
			want:  1.0,
		},
		{
			qubit: qubit.Zero(),
			pX:    1.0,
			pY:    0,
			pZ:    0,
			m:     qubit.One(),
			want:  1.0,
		},
		{
			qubit: qubit.Zero(),
			pX:    0.5,
			pY:    0,
			pZ:    0,
			m:     qubit.Zero(),
			want:  0.5,
		},
		{
			qubit: qubit.Zero(),
			pX:    0.5,
			pY:    0,
			pZ:    0,
			m:     qubit.One(),
			want:  0.5,
		},
		{
			qubit: qubit.Zero(),
			pX:    0,
			pY:    0,
			pZ:    1.0,
			m:     qubit.Zero(),
			want:  1.0,
		},
		{
			qubit: qubit.Plus(),
			pX:    0,
			pY:    0,
			pZ:    1.0,
			m:     qubit.Minus(),
			want:  1.0,
		},
		{
			qubit: qubit.Zero(),
			pX:    0,
			pY:    1.0,
			pZ:    0,
			m:     qubit.One(),
			want:  1.0,
		},
		{
			qubit: qubit.Zero(),
			pX:    1.0 / 3,
			pY:    1.0 / 3,
			pZ:    1.0 / 3,
			m:     qubit.Zero(),
			want:  1.0 / 3,
		},
		{
			qubit: qubit.Zero(),
			pX:    1.0 / 3,
			pY:    1.0 / 3,
			pZ:    1.0 / 3,
			m:     qubit.One(),
			want:  1.0/3 + 1.0/3,
		},
	}

	for _, c := range cases {
		got, _ := density.New(c.qubit).
			Pauli(c.pX, c.pY, c.pZ, 0).
			Measure(observable.Projector(
				c.m,
			))

		if !epsilon.IsCloseF64(got, c.want) {
			t.Errorf("got %v, want %v", got, c.want)
		}
	}
}

func TestDensityMatrix_TraceDistance(t *testing.T) {
	cases := []struct {
		s1   []density.WeightedState
		s2   []density.WeightedState
		want float64
	}{
		{
			s1: []density.WeightedState{
				{Probability: 1, Qubit: qubit.Zero()},
			},
			s2: []density.WeightedState{
				{Probability: 1, Qubit: qubit.Zero()},
			},
			want: 0,
		},
		{
			s1: []density.WeightedState{
				{Probability: 1, Qubit: qubit.Zero()},
			},
			s2: []density.WeightedState{
				{Probability: 1, Qubit: qubit.One()},
			},
			want: 1,
		},
		{
			s1: []density.WeightedState{
				{Probability: 1, Qubit: qubit.Zero()},
			},
			s2: []density.WeightedState{
				{Probability: 1, Qubit: qubit.Plus()},
			},
			want: math.Sqrt(0.5),
		},
		{
			s1: []density.WeightedState{
				{Probability: 0.5, Qubit: qubit.Zero()},
				{Probability: 0.5, Qubit: qubit.One()},
			},
			s2: []density.WeightedState{
				{Probability: 1, Qubit: qubit.Zero()},
			},
			want: 0.5,
		},
		{
			s1: []density.WeightedState{
				{Probability: 1, Qubit: qubit.Plus()},
			},
			s2: []density.WeightedState{
				{Probability: 0.5, Qubit: qubit.Zero()},
				{Probability: 0.5, Qubit: qubit.One()},
			},
			want: 0.5,
		},
	}

	for _, c := range cases {
		rhoA := density.NewMixed(c.s1)
		selfA := rhoA.TraceDistance(rhoA)
		if !epsilon.IsZeroF64(selfA) {
			t.Errorf("got=%v, want=%v", selfA, 0)
		}

		rhoB := density.NewMixed(c.s2)
		selfB := rhoB.TraceDistance(rhoB)
		if !epsilon.IsZeroF64(selfB) {
			t.Errorf("got=%v, want=%v", selfB, 0)
		}

		gotAB := rhoA.TraceDistance(rhoB)
		if !epsilon.IsCloseF64(gotAB, c.want) {
			t.Errorf("got=%v, want=%v", gotAB, c.want)
		}

		gotBA := rhoB.TraceDistance(rhoA)
		if !epsilon.IsCloseF64(gotBA, c.want) {
			t.Errorf("got=%v, want=%v", gotBA, c.want)
		}

		if !epsilon.IsCloseF64(gotAB, gotBA) {
			t.Errorf("gotAB=%v, gotBA=%v", gotAB, gotBA)
		}
	}
}

func TestDensityMatrix_Fidelity(t *testing.T) {
	cases := []struct {
		s1   []density.WeightedState
		s2   []density.WeightedState
		want float64
	}{
		{
			s1: []density.WeightedState{
				{Probability: 1, Qubit: qubit.Zero()},
			},
			s2: []density.WeightedState{
				{Probability: 1, Qubit: qubit.Zero()},
			},
			want: 1,
		},
		{
			s1: []density.WeightedState{
				{Probability: 1, Qubit: qubit.Zero()},
			},
			s2: []density.WeightedState{
				{Probability: 1, Qubit: qubit.One()},
			},
			want: 0,
		},
		{
			s1: []density.WeightedState{
				{Probability: 1, Qubit: qubit.Zero()},
			},
			s2: []density.WeightedState{
				{Probability: 1, Qubit: qubit.Plus()},
			},
			want: math.Sqrt(0.5),
		},
		{
			s1: []density.WeightedState{
				{Probability: 0.5, Qubit: qubit.Zero()},
				{Probability: 0.5, Qubit: qubit.One()},
			},
			s2: []density.WeightedState{
				{Probability: 1, Qubit: qubit.Zero()},
			},
			want: math.Sqrt(0.5),
		},
		{
			s1: []density.WeightedState{
				{Probability: 1, Qubit: qubit.Plus()},
			},
			s2: []density.WeightedState{
				{Probability: 0.5, Qubit: qubit.Zero()},
				{Probability: 0.5, Qubit: qubit.One()},
			},
			want: math.Sqrt(0.5),
		},
	}

	for _, c := range cases {
		rhoA := density.NewMixed(c.s1)
		selfA := rhoA.Fidelity(rhoA)
		if !epsilon.IsOneF64(selfA) {
			t.Errorf("got=%v, want=%v", selfA, 1)
		}

		rhoB := density.NewMixed(c.s2)
		selfB := rhoB.Fidelity(rhoB)
		if !epsilon.IsOneF64(selfB) {
			t.Errorf("got=%v, want=%v", selfB, 1)
		}

		gotAB := rhoA.Fidelity(rhoB)
		if !epsilon.IsCloseF64(gotAB, c.want) {
			t.Errorf("got=%v, want=%v", gotAB, c.want)
		}

		gotBA := rhoB.Fidelity(rhoA)
		if !epsilon.IsCloseF64(gotBA, c.want) {
			t.Errorf("got=%v, want=%v", gotBA, c.want)
		}
	}
}

func TestDensityMatrix_VonNeumannEntropy(t *testing.T) {
	cases := []struct {
		s    []density.WeightedState
		want float64
	}{
		{
			s: []density.WeightedState{
				{Probability: 1, Qubit: qubit.Zero()},
			},
			want: 0,
		},
		{
			s: []density.WeightedState{
				{Probability: 1, Qubit: qubit.Zero()},
				{Probability: 0, Qubit: qubit.One()},
			},
			want: 0,
		},
		{
			s: []density.WeightedState{
				{Probability: 0.5, Qubit: qubit.Zero()},
				{Probability: 0.5, Qubit: qubit.One()},
			},
			want: 1,
		},
		{
			s: []density.WeightedState{
				{Probability: 0.9, Qubit: qubit.Zero()},
				{Probability: 0.1, Qubit: qubit.One()},
			},
			want: -(0.9*math.Log2(0.9) + 0.1*math.Log2(0.1)),
		},
		{
			s: []density.WeightedState{
				{Probability: 0.5, Qubit: qubit.Zero()},
				{Probability: 0.5, Qubit: qubit.Zero()},
			},
			want: 0,
		},
		{
			s: []density.WeightedState{
				{Probability: 0.5, Qubit: qubit.Zero()},
				{Probability: 0.5, Qubit: qubit.Plus()},
			},
			want: -((0.5+1/(2*math.Sqrt2))*math.Log2(0.5+1/(2*math.Sqrt2)) + (0.5-1/(2*math.Sqrt2))*math.Log2(0.5-1/(2*math.Sqrt2))),
		},
		{
			s: []density.WeightedState{
				{Probability: 0.25, Qubit: qubit.From("00")},
				{Probability: 0.25, Qubit: qubit.From("01")},
				{Probability: 0.25, Qubit: qubit.From("10")},
				{Probability: 0.25, Qubit: qubit.From("11")},
			},
			want: 2,
		},
	}

	for _, c := range cases {
		rho := density.NewMixed(c.s)
		got := rho.VonNeumannEntropy()

		if !epsilon.IsCloseF64(got, c.want) {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}

func TestDensityMatrix_RelativeEntropy(t *testing.T) {
	cases := []struct {
		s1   []density.WeightedState
		s2   []density.WeightedState
		want float64
	}{
		{
			s1: []density.WeightedState{
				{Probability: 1, Qubit: qubit.Zero()},
			},
			s2: []density.WeightedState{
				{Probability: 1, Qubit: qubit.Zero()},
			},
			want: 0,
		},
		{
			s1: []density.WeightedState{
				{Probability: 0.3, Qubit: qubit.Zero()},
				{Probability: 0.7, Qubit: qubit.One()},
			},
			s2: []density.WeightedState{
				{Probability: 0.3, Qubit: qubit.Zero()},
				{Probability: 0.7, Qubit: qubit.One()},
			},
			want: 0,
		},
		{
			s1: []density.WeightedState{
				{Probability: 1, Qubit: qubit.Zero()},
			},
			s2: []density.WeightedState{
				{Probability: 0.5, Qubit: qubit.Zero()},
				{Probability: 0.5, Qubit: qubit.One()},
			},
			want: 1,
		},
		{
			s1: []density.WeightedState{
				{Probability: 1, Qubit: qubit.Plus()},
			},
			s2: []density.WeightedState{
				{Probability: 0.5, Qubit: qubit.Zero()},
				{Probability: 0.5, Qubit: qubit.One()},
			},
			want: 1,
		},
		{
			s1: []density.WeightedState{
				{Probability: 0.5, Qubit: qubit.Zero()},
				{Probability: 0.5, Qubit: qubit.One()},
			},
			s2: []density.WeightedState{
				{Probability: 0.75, Qubit: qubit.Zero()},
				{Probability: 0.25, Qubit: qubit.One()},
			},
			want: 0.5*math.Log2(0.5/0.75) + 0.5*math.Log2(0.5/0.25),
		},
		{

			s1: []density.WeightedState{
				{Probability: 0.75, Qubit: qubit.Zero()},
				{Probability: 0.25, Qubit: qubit.One()},
			},
			s2: []density.WeightedState{
				{Probability: 0.5, Qubit: qubit.Zero()},
				{Probability: 0.5, Qubit: qubit.One()},
			},
			want: 0.75*math.Log2(0.75/0.5) + 0.25*math.Log2(0.25/0.5),
		},
		{
			s1: []density.WeightedState{
				{Probability: 1, Qubit: qubit.Zero()},
			},
			s2: []density.WeightedState{
				{Probability: 1, Qubit: qubit.One()},
			},
			want: math.Inf(1),
		},
		{
			s1: []density.WeightedState{
				{Probability: 0.5, Qubit: qubit.Zero()},
				{Probability: 0.5, Qubit: qubit.One()},
			},
			s2: []density.WeightedState{
				{Probability: 1, Qubit: qubit.Zero()},
			},
			want: math.Inf(1),
		},
		{
			s1: []density.WeightedState{
				{Probability: 1, Qubit: qubit.Plus()},
			},
			s2: []density.WeightedState{
				{Probability: 1, Qubit: qubit.Zero()},
			},
			want: math.Inf(1),
		},
		{
			s1: []density.WeightedState{
				{Probability: 1, Qubit: qubit.Zero()},
			},
			s2: []density.WeightedState{
				{Probability: 1 - 1e-12, Qubit: qubit.Zero()},
				{Probability: 1e-12, Qubit: qubit.One()},
			},
			want: -math.Log2(1 - 1e-12),
		},
	}

	for _, c := range cases {
		rho := density.NewMixed(c.s1)
		sigma := density.NewMixed(c.s2)
		got := rho.RelativeEntropy(sigma)

		if math.IsInf(got, 1) && math.IsInf(c.want, 1) {
			continue
		}

		if !epsilon.IsCloseF64(got, c.want) {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}
