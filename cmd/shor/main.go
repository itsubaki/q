package main

import (
	"flag"
	"fmt"
	"slices"
	"strings"

	"github.com/itsubaki/q"
	"github.com/itsubaki/q/math/matrix"
	"github.com/itsubaki/q/math/number"
	"github.com/itsubaki/q/math/rand"
	"github.com/itsubaki/q/math/vector"
	"github.com/itsubaki/q/quantum/gate"
	"github.com/itsubaki/q/quantum/qubit"
)

// go run main.go --N 15
func main() {
	var N, t, a int
	var seed uint64
	flag.IntVar(&N, "N", 15, "positive integer")
	flag.IntVar(&t, "t", 3, "precision bits")
	flag.IntVar(&a, "a", -1, "coprime number of N")
	flag.Uint64Var(&seed, "seed", 0, "PRNG seed for measurements")
	flag.Parse()

	if N < 2 {
		fmt.Printf("N=%d. N must be greater than 1.\n", N)
		return
	}

	if number.IsEven(N) {
		fmt.Printf("N=%d is even. p=%d, q=%d.\n", N, 2, N/2)
		return
	}

	if a, b, ok := number.BaseExp(N); ok {
		fmt.Printf("N=%d. N is exponentiation. %d^%d.\n", N, a, b)
		return
	}

	if number.IsPrime(N) {
		fmt.Printf("N=%d is prime.\n", N)
		return
	}

	if a < 0 {
		a = rand.Coprime(N)
	}

	if a < 2 || a > N-1 {
		fmt.Printf("N=%d, a=%d. a must be 1 < a < N.\n", N, a)
		return
	}

	if number.GCD(N, a) != 1 {
		fmt.Printf("N=%d, a=%d. a is not coprime. a is non-trivial factor.\n", N, a)
		return
	}

	fmt.Printf("N=%d, a=%d, t=%d, seed=%d.\n\n", N, a, t, seed)

	qsim := q.New()
	if seed > 0 {
		qsim.Rand = rand.Const(seed)
	}

	r0 := qsim.Zeros(t)
	r1 := qsim.ZeroLog2(N)

	qsim.X(r1[len(r1)-1])
	print("initial state", qsim, r0, r1)

	qsim.H(r0...)
	print("create superposition", qsim, r0, r1)

	for j := range r0 {
		CModExp2(qsim, a, j, N, r0[j], r1)
		print(fmt.Sprintf("apply controlled-U[%d]", j), qsim, r0, r1)
	}

	qsim.InvQFT(r0...)
	print("apply inverse QFT", qsim, r0, r1)

	qsim.Measure(r1...)
	print("measure reg1", qsim, r0, r1)

	var prop float64
	for _, s := range qsim.State(r0) {
		i, m := s.Int(), s.BinaryString()
		ss, r, d, ok := number.FindOrder(a, N, fmt.Sprintf("0.%s", m))
		if !ok || number.IsOdd(r) {
			fmt.Printf("  i=%4d: N=%d, a=%d, t=%d; s/r=%4d/%4d ([0.%v]~%.4f);\n", i, N, a, t, ss, r, m, d)
			continue
		}

		p0 := number.GCD(number.Pow(a, r/2)-1, N)
		p1 := number.GCD(number.Pow(a, r/2)+1, N)
		if number.IsTrivial(N, p0, p1) {
			fmt.Printf("  i=%4d: N=%d, a=%d, t=%d; s/r=%4d/%4d ([0.%v]~%.4f); p=%v, q=%v.\n", i, N, a, t, ss, r, m, d, p0, p1)
			continue
		}

		fmt.Printf("* i=%4d: N=%d, a=%d, t=%d; s/r=%4d/%4d ([0.%v]~%.4f); p=%v, q=%v.\n", i, N, a, t, ss, r, m, d, p0, p1)
		prop += s.Probability()
	}

	fmt.Printf("total probability: %.8f\n", prop)
}

func print(desc string, qsim *q.Q, reg ...any) {
	fmt.Println(desc)

	max := slices.Max(qsim.Probability())
	for _, s := range qsim.State(reg...) {
		p := strings.Repeat("*", int(s.Probability()/max*32))
		fmt.Printf("%s: %s\n", s, p)
	}

	fmt.Println()
}

// CModExp2 applies controlled modular exponentiation operation.
func CModExp2(qsim *q.Q, a, j, N int, control q.Qubit, target []q.Qubit) {
	ControlledModExp2(qsim.Underlying(), a, j, N, control.Index(), q.Index(target...))
}

// ControlledModExp2 applies the controlled modular exponentiation operation.
// |j>|k> -> |j>|a**(2**j) * k mod N>.
func ControlledModExp2(qb *qubit.Qubit, a, j, N, control int, target []int) {
	n := qb.NumQubits()
	state := qb.Amplitude()
	a2jModN := number.ModExp2(a, j, N)
	cmask := 1 << (n - 1 - control)

	newState := make([]complex128, qb.Dim())
	for i := range qb.Dim() {
		if (i & cmask) == 0 {
			newState[i] = state[i]
			continue
		}

		// binary to integer
		var k int
		for j, t := range target {
			k |= ((i >> (n - 1 - t)) & 1) << (len(target) - 1 - j)
		}

		// (a**(2**j) * k) mod N
		a2jkModN := (a2jModN * k) % N

		// integer to binary
		newIdx := i
		for j, t := range target {
			bit := (a2jkModN >> (len(target) - 1 - j)) & 1
			pos := n - 1 - t
			newIdx = (newIdx &^ (1 << pos)) | (bit << pos)
		}

		// update the state
		newState[newIdx] = state[i]
	}

	// update the qubit state
	qb.Update(vector.New(newState...))
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
