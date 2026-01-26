package main

import (
	"fmt"
	"testing"

	"github.com/itsubaki/q"
	"github.com/itsubaki/q/math/epsilon"
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

func TestEigenVector(t *testing.T) {
	cases := []struct {
		N, a, t int
		bin     []string
		amp     []complex128
	}{
		{
			15, 7, 3,
			[]string{"0001", "0100", "0111", "1101"},
			[]complex128{1, 0, 0, 0},
		},
	}

	for _, c := range cases {
		qsim := q.New()
		r0 := qsim.Zeros(c.t)
		r1 := qsim.ZeroLog2(c.N)

		qsim.X(r1[len(r1)-1])
		qsim.H(r0...)
		for j := range r0 {
			CModExp2(qsim, c.a, j, c.N, r0[j], r1)
		}
		qsim.InvQFT(r0...)

		us := make(map[string]complex128)
		for _, s := range qsim.State(r1) {
			m := s.BinaryString()[0]
			if v, ok := us[m]; ok {
				us[m] = v + s.Amplitude()
				continue
			}

			us[m] = s.Amplitude()
		}

		if len(us) != len(c.bin) {
			t.Fail()
		}

		for i := range c.bin {
			if !epsilon.IsClose(us[c.bin[i]], c.amp[i]) {
				t.Fail()
			}
		}
	}
}

func ExampleControlledModExp2g() {
	n, a, j, N := 5, 7, 0, 15
	g := ControlledModExp2g(n, a, j, N, 0, []int{1, 2, 3, 4})

	mask := (1 << (n - 1)) - 1
	a2j := number.ModExp2(a, j, N)

	for i, r := range g.Transpose().Seq2() {
		if (i>>(n-1))&1 == 0 {
			// control qubit is zero
			continue
		}

		// decimal number representation of target qubits
		k := i & mask
		if k >= N {
			continue
		}

		for l, e := range r {
			if epsilon.IsZero(e) {
				continue
			}

			expected := l & mask
			got := (a2j * k) % N
			fmt.Printf("1:%04b=%2d 1:%04b=%2d %2d\n", k, k, expected, expected, got)
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

func TestControlledModExp2g(t *testing.T) {
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

		if !got.Equal(c.want) {
			t.Fail()
		}
	}
}

// ControlledModExp2g returns gate of controlled modular exponentiation operation.
// |j>|k> -> |j>|a**(2**j) * k mod N>.
func ControlledModExp2g(n, a, j, N, c int, t []int) *matrix.Matrix {
	m := gate.I(n)
	r1len := len(t)
	a2jmodN := number.ModExp2(a, j, N)

	d, _ := m.Dim()
	idx := make([]int, d)
	for i := range d {
		if (i>>(n-1-c))&1 == 0 {
			// control bit is 0, then do nothing
			idx[i] = i
			continue
		}

		// r1len bits of i
		mask := (1 << r1len) - 1
		k := i & mask
		if k > N-1 {
			idx[i] = i
			continue
		}

		// r0len bits of i + a2jkmodN bits
		a2jkmodN := a2jmodN * k % N
		idx[i] = (i >> r1len << r1len) | a2jkmodN
	}

	data := make([][]complex128, d)
	for i, j := range idx {
		data[j] = m.Row(i)
	}

	return matrix.New(data...)
}
