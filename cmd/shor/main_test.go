package main

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/itsubaki/q"
	"github.com/itsubaki/q/math/matrix"
	"github.com/itsubaki/q/math/number"
	"github.com/itsubaki/q/quantum/gate"
)

func ExampleCModExp2_mod15() {
	qsim := q.New()
	c := qsim.Zero()
	t := qsim.ZeroLog2(15)

	qsim.X(c)
	qsim.X(t[len(t)-1])
	for _, s := range qsim.State(c, t) {
		fmt.Println(s)
	}

	a, N := 7, 15
	for j := range 3 {
		CModExp2(qsim, a, j, N, c, t)
		for _, s := range qsim.State(c, t) {
			fmt.Println(s)
		}
	}

	// Output:
	// [1 0001][  1   1]( 1.0000 0.0000i): 1.0000
	// [1 0111][  1   7]( 1.0000 0.0000i): 1.0000
	// [1 1101][  1  13]( 1.0000 0.0000i): 1.0000
	// [1 1101][  1  13]( 1.0000 0.0000i): 1.0000
}

func ExampleCModExp2_mod21() {
	qsim := q.New()
	c := qsim.Zero()
	t := qsim.ZeroLog2(21)

	qsim.X(c)
	qsim.X(t[len(t)-1])

	for _, s := range qsim.State(c, t) {
		fmt.Println(s)
	}

	a, N := 2, 21
	for j := range 4 {
		CModExp2(qsim, a, j, N, c, t)
		for _, s := range qsim.State(c, t) {
			fmt.Println(s)
		}
	}

	// Output:
	// [1 00001][  1   1]( 1.0000 0.0000i): 1.0000
	// [1 00010][  1   2]( 1.0000 0.0000i): 1.0000
	// [1 01000][  1   8]( 1.0000 0.0000i): 1.0000
	// [1 00010][  1   2]( 1.0000 0.0000i): 1.0000
	// [1 01000][  1   8]( 1.0000 0.0000i): 1.0000
}

func ExampleControlledModExp2g() {
	n, a, j, N := 5, 7, 0, 15
	g := ControlledModExp2g(n, a, j, N, 0, []int{1, 2, 3, 4})

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
		got := ControlledModExp2g(c.n, c.a, c.j, c.N, c.c, c.t)
		if !got.IsUnitary() {
			t.Errorf("modexp is not unitary")
		}

		if !got.Equals(c.want) {
			t.Fail()
		}
	}
}
