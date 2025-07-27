package gate_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/itsubaki/q/math/matrix"
	"github.com/itsubaki/q/quantum/gate"
)

func ExampleX() {
	g := gate.X()
	for _, r := range g.Seq2() {
		fmt.Println(r)
	}

	// Output:
	// [(0+0i) (1+0i)]
	// [(1+0i) (0+0i)]
}

func ExampleX_xX() {
	g := gate.X(2)
	for _, r := range g.Seq2() {
		fmt.Println(r)
	}

	// Output:
	// [(0+0i) (0+0i) (0+0i) (1+0i)]
	// [(0+0i) (0+0i) (1+0i) (0+0i)]
	// [(0+0i) (1+0i) (0+0i) (0+0i)]
	// [(1+0i) (0+0i) (0+0i) (0+0i)]
}

func ExampleH() {
	g := gate.H()
	for _, r := range g.Seq2() {
		fmt.Printf("%.4v\n", r)
	}

	// Output:
	// [(0.7071+0i) (0.7071+0i)]
	// [(0.7071+0i) (-0.7071+0i)]
}

func ExampleH_hH() {
	g := gate.H(2)
	for _, r := range g.Seq2() {
		fmt.Printf("%.4v\n", r)
	}

	// Output:
	// [(0.5+0i) (0.5+0i) (0.5+0i) (0.5+0i)]
	// [(0.5+0i) (-0.5+0i) (0.5+0i) (-0.5+0i)]
	// [(0.5+0i) (0.5+0i) (-0.5+0i) (-0.5+0i)]
	// [(0.5+0i) (-0.5+0i) (-0.5+0i) (0.5-0i)]
}

func ExampleCNOT() {
	g := gate.CNOT(2, 0, 1)
	for _, r := range g.Seq2() {
		fmt.Println(r)
	}

	// Output:
	// [(1+0i) (0+0i) (0+0i) (0+0i)]
	// [(0+0i) (1+0i) (0+0i) (0+0i)]
	// [(0+0i) (0+0i) (0+0i) (1+0i)]
	// [(0+0i) (0+0i) (1+0i) (0+0i)]
}

func ExampleSwap() {
	g := gate.Swap(2, 0, 1)
	for _, r := range g.Seq2() {
		fmt.Println(r)
	}

	// Output:
	// [(1+0i) (0+0i) (0+0i) (0+0i)]
	// [(0+0i) (0+0i) (1+0i) (0+0i)]
	// [(0+0i) (1+0i) (0+0i) (0+0i)]
	// [(0+0i) (0+0i) (0+0i) (1+0i)]
}

func ExampleTensorProduct() {
	for _, r := range gate.TensorProduct(gate.X(), 2, []int{1}).Seq2() {
		fmt.Printf("%.4v\n", r)
	}
	fmt.Println()

	for _, r := range gate.TensorProduct(gate.X(), 2, []int{0}).Seq2() {
		fmt.Printf("%.4v\n", r)
	}

	// Output:
	// [(0+0i) (1+0i) (0+0i) (0+0i)]
	// [(1+0i) (0+0i) (0+0i) (0+0i)]
	// [(0+0i) (0+0i) (0+0i) (1+0i)]
	// [(0+0i) (0+0i) (1+0i) (0+0i)]
	//
	// [(0+0i) (0+0i) (1+0i) (0+0i)]
	// [(0+0i) (0+0i) (0+0i) (1+0i)]
	// [(1+0i) (0+0i) (0+0i) (0+0i)]
	// [(0+0i) (1+0i) (0+0i) (0+0i)]
}

func TestU(t *testing.T) {
	cases := []struct {
		in, want *matrix.Matrix
	}{
		{gate.U(0, 0, 0), gate.I()},
		{gate.U(math.Pi, 0, math.Pi), gate.X()},
		{gate.U(math.Pi, math.Pi/2, math.Pi/2), gate.Y()},
		{gate.U(0, math.Pi, 0), gate.Z()},
		{gate.U(math.Pi/2, 0, math.Pi), gate.H()},
		{gate.U(math.Pi/4, -1*math.Pi/2, math.Pi/2), gate.RX(math.Pi / 4)},
		{gate.U(math.Pi/4, 0, 0), gate.RY(math.Pi / 4)},
	}

	for _, c := range cases {
		if !c.in.Equals(c.want) {
			t.Fail()
		}
	}
}

func TestC(t *testing.T) {
	cases := []struct {
		in, want *matrix.Matrix
	}{
		{gate.C(gate.I(), 2, 0, 1), gate.I(2)},
		{gate.C(gate.X(), 2, 0, 1), gate.CNOT(2, 0, 1)},
		{gate.C(gate.Z(), 2, 0, 1), gate.CZ(2, 0, 1)},
		{gate.C(gate.S(), 2, 0, 1), gate.CS(2, 0, 1)},
		{gate.C(gate.R(gate.Theta(4)), 2, 0, 1), gate.CR(gate.Theta(4), 2, 0, 1)},
		{gate.C(gate.X(), 3, 0, 2), gate.CNOT(3, 0, 2)},
		{gate.C(gate.X(), 3, 0, 1), gate.CNOT(3, 0, 1)},
		{gate.C(gate.X(), 3, 1, 0), gate.CNOT(3, 1, 0)},
		{gate.C(gate.U(math.Pi, 0, math.Pi), 3, 1, 0), gate.CNOT(3, 1, 0)},
	}

	for _, c := range cases {
		if !c.in.Equals(c.want) {
			t.Fail()
		}
	}
}

func TestControlled(t *testing.T) {
	cases := []struct {
		in, want *matrix.Matrix
	}{
		{gate.Controlled(gate.I(), 2, []int{0}, 1), gate.I(2)},
		{gate.Controlled(gate.X(), 3, []int{1, 2}, 0), gate.CCNOT(3, 1, 2, 0)},
		{gate.Controlled(gate.X(), 3, []int{2, 1}, 0), gate.CCNOT(3, 2, 1, 0)},
		{gate.Controlled(gate.X(), 3, []int{0, 2}, 1), gate.CCNOT(3, 0, 2, 1)},
		{gate.Controlled(gate.X(), 3, []int{2, 0}, 1), gate.CCNOT(3, 2, 0, 1)},
		{gate.Controlled(gate.X(), 3, []int{0, 1}, 2), gate.CCNOT(3, 0, 1, 2)},
		{gate.Controlled(gate.X(), 3, []int{1, 0}, 2), gate.CCNOT(3, 1, 0, 2)},
	}

	for _, c := range cases {
		if !c.in.Equals(c.want) {
			t.Fail()
		}
	}
}

func TestInverse(t *testing.T) {
	cases := []struct {
		in, want *matrix.Matrix
	}{
		{gate.U(1, 2, 3), gate.I()},
		{gate.X(2), gate.I(2)},
		{gate.CNOT(2, 0, 1), gate.I(2)},
	}

	for _, c := range cases {
		if !c.in.Apply(c.in.Inverse()).Equals(c.want) {
			t.Fail()
		}
	}
}

func TestIsHermite(t *testing.T) {
	cases := []struct {
		in   *matrix.Matrix
		want bool
	}{
		{gate.H(), true},
		{gate.X(), true},
		{gate.Y(), true},
		{gate.Z(), true},
		{gate.New(
			[]complex128{1, 2},
			[]complex128{3, 4},
		), false},
	}

	for _, c := range cases {
		if c.in.IsHermite() != c.want {
			t.Fail()
		}
	}
}

func TestIsUnitary(t *testing.T) {
	cases := []struct {
		in   *matrix.Matrix
		want bool
	}{
		{gate.I(), true},
		{gate.H(), true},
		{gate.X(), true},
		{gate.Y(), true},
		{gate.Z(), true},
		{gate.U(1, 2, 3), true},
		{gate.R(gate.Theta(4)), true},
		{gate.RX(1.23), true},
		{gate.RY(1.23), true},
		{gate.RZ(1.23), true},
		{gate.ControlledS(2, []int{0}, 1), true},
		{gate.ControlledR(gate.Theta(4), 2, []int{0}, 1), true},
		{gate.CS(2, 0, 1), true},
		{gate.CR(gate.Theta(4), 2, 0, 1), true},
		{gate.QFT(2), true},
		{gate.New(
			[]complex128{1, 2},
			[]complex128{3, 4},
		), false},
	}

	for _, c := range cases {
		if c.in.IsUnitary() != c.want {
			t.Fail()
		}
	}
}

func TestTrace(t *testing.T) {
	cases := []struct {
		in   *matrix.Matrix
		want complex128
	}{
		{gate.I(), complex(2, 0)},
		{gate.H(), complex(0, 0)},
		{gate.X(), complex(0, 0)},
		{gate.Y(), complex(0, 0)},
		{gate.Z(), complex(0, 0)},
	}

	for _, c := range cases {
		got := c.in.Trace()
		if got != c.want {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}

func TestMultiQubitGate(t *testing.T) {
	cases := []struct {
		in   *matrix.Matrix
		want *matrix.Matrix
	}{
		{
			in: gate.CZ(3, 0, 2),
			want: gate.New(
				[]complex128{1, 0, 0, 0, 0, 0, 0, 0},
				[]complex128{0, 1, 0, 0, 0, 0, 0, 0},
				[]complex128{0, 0, 1, 0, 0, 0, 0, 0},
				[]complex128{0, 0, 0, 1, 0, 0, 0, 0},
				[]complex128{0, 0, 0, 0, 1, 0, 0, 0},
				[]complex128{0, 0, 0, 0, 0, -1, 0, 0},
				[]complex128{0, 0, 0, 0, 0, 0, 1, 0},
				[]complex128{0, 0, 0, 0, 0, 0, 0, -1},
			),
		},
		{
			in: gate.CNOT(3, 0, 2),
			want: gate.New(
				[]complex128{1, 0, 0, 0, 0, 0, 0, 0},
				[]complex128{0, 1, 0, 0, 0, 0, 0, 0},
				[]complex128{0, 0, 1, 0, 0, 0, 0, 0},
				[]complex128{0, 0, 0, 1, 0, 0, 0, 0},
				[]complex128{0, 0, 0, 0, 0, 1, 0, 0},
				[]complex128{0, 0, 0, 0, 1, 0, 0, 0},
				[]complex128{0, 0, 0, 0, 0, 0, 0, 1},
				[]complex128{0, 0, 0, 0, 0, 0, 1, 0},
			),
		},
		{
			in: gate.CCNOT(3, 0, 1, 2),
			want: func() *matrix.Matrix {
				g := make([]*matrix.Matrix, 13)
				g[0] = gate.I(2).TensorProduct(gate.H())
				g[1] = gate.I(1).TensorProduct(gate.CNOT(2, 0, 1))
				g[2] = gate.I(2).TensorProduct(gate.T().Dagger())
				g[3] = gate.CNOT(3, 0, 2)
				g[4] = gate.I(2).TensorProduct(gate.T())
				g[5] = gate.I(1).TensorProduct(gate.CNOT(2, 0, 1))
				g[6] = gate.I(2).TensorProduct(gate.T().Dagger())
				g[7] = gate.CNOT(3, 0, 2)
				g[8] = gate.I(1).TensorProduct(gate.T().Dagger()).TensorProduct(gate.T())
				g[9] = gate.CNOT(2, 0, 1).TensorProduct(gate.H())
				g[10] = gate.I(1).TensorProduct(gate.T().Dagger()).TensorProduct(gate.I())
				g[11] = gate.CNOT(2, 0, 1).TensorProduct(gate.I())
				g[12] = gate.T().TensorProduct(gate.S()).TensorProduct(gate.I())

				w := gate.I(3)
				for _, v := range g {
					w = w.Apply(v)
				}

				return w
			}(),
		},
		{
			in: gate.ControlledNot(2, []int{0}, 1),
			want: func() *matrix.Matrix {
				g0 := gate.I().Add(gate.Z()).TensorProduct(gate.I())
				g1 := gate.I().Sub(gate.Z()).TensorProduct(gate.X())
				return g0.Add(g1).Mul(0.5)
			}(),
		},
	}

	for _, c := range cases {
		if !c.in.Equals(c.want) {
			t.Errorf("got=%v, want=%v", c.in, c.want)
		}
	}
}
