package gate_test

import (
	"fmt"
	"math"
	"strconv"
	"testing"

	"github.com/itsubaki/q/math/matrix"
	"github.com/itsubaki/q/math/number"
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

func ExampleControlledModExp2() {
	n, a, j, N := 5, 7, 0, 15
	g := gate.ControlledModExp2(n, a, j, N, 0, []int{1, 2, 3, 4})

	for i, r := range g.Transpose().Seq2() {
		bin := fmt.Sprintf("%0*b", n, i)
		if bin[:1] == "0" { // control qubit is |0>
			continue
		}

		// decimal number representation of target qubits
		k := number.Must(strconv.ParseInt(bin[1:], 2, 64))
		if k >= int64(N) {
			continue
		}

		for l, e := range r {
			if e == complex(0, 0) {
				continue
			}

			// decimal number representation of a^2^j * k mod N
			a2jkmodNs := fmt.Sprintf("%0*s", n, strconv.FormatInt(int64(l), 2)[1:])
			a2jkmodN := number.Must(strconv.ParseInt(a2jkmodNs, 2, 64))
			got := (int64(number.ModExp2(a, j, N)) * k) % int64(N)

			fmt.Printf("%s:%s=%2d %s:%s=%2d %2d\n", bin[:1], bin[1:], k, bin[:1], a2jkmodNs[1:], a2jkmodN, got)
		}
	}

	// Output:
	// 1:0000= 0 1:0000= 0  0
	// 1:0001= 1 1:0111= 7  7
	// 1:0010= 2 1:1110=14 14
	// 1:0011= 3 1:0110= 6  6
	// 1:0100= 4 1:1101=13 13
	// 1:0101= 5 1:0101= 5  5
	// 1:0110= 6 1:1100=12 12
	// 1:0111= 7 1:0100= 4  4
	// 1:1000= 8 1:1011=11 11
	// 1:1001= 9 1:0011= 3  3
	// 1:1010=10 1:1010=10 10
	// 1:1011=11 1:0010= 2  2
	// 1:1100=12 1:1001= 9  9
	// 1:1101=13 1:0001= 1  1
	// 1:1110=14 1:1000= 8  8
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

func TestControlledModExp2(t *testing.T) {
	g1 := matrix.Apply(
		gate.CNOT(7, 3, 5),
		gate.CCNOT(7, 1, 5, 3),
		gate.CNOT(7, 3, 5),
	)
	g2 := matrix.Apply(
		gate.CNOT(7, 4, 6),
		gate.CCNOT(7, 1, 6, 4),
		gate.CNOT(7, 4, 6),
	)
	g3 := gate.I(7)

	cases := []struct {
		n, a, j, N int
		c          int
		t          []int
		want       *matrix.Matrix
	}{
		{7, 7, 1, 15, 1, []int{4, 5, 6, 7}, matrix.Apply(g1, g2)},
		{7, 7, 2, 15, 0, []int{4, 5, 6, 7}, g3},
	}

	for _, c := range cases {
		got := gate.ControlledModExp2(c.n, c.a, c.j, c.N, c.c, c.t)
		if !got.IsUnitary() {
			t.Errorf("modexp is not unitary")
		}

		if !got.Equals(c.want) {
			t.Fail()
		}
	}
}
