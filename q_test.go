package q

import (
	"fmt"
	"math"
	"testing"

	"github.com/itsubaki/q/pkg/math/matrix"
	"github.com/itsubaki/q/pkg/math/number"
	"github.com/itsubaki/q/pkg/math/rand"
	"github.com/itsubaki/q/pkg/quantum/gate"
	"github.com/itsubaki/q/pkg/quantum/qubit"
)

func TestQSimFactoringN(t *testing.T) {
	N := 21
	a := rand.Coprime(N)

	qsim := New()
	r0 := qsim.ZeroWith(4)
	r1 := qsim.ZeroLog2(N)

	qsim.X(r1[len(r1)-1])
	qsim.H(r0...)
	qsim.CModExp2(a, N, r0, r1)
	qsim.InvQFT(r0...)

	for i := 0; i < 10; i++ {
		m := qsim.Clone().MeasureAsBinary(r0...)
		d := number.BinaryFraction(m)
		_, s, r := number.ContinuedFraction(d)

		if number.IsOdd(r) || number.Pow(a, r/2)%N == -1 {
			continue
		}

		p0 := number.GCD(number.Pow(a, r/2)-1, N)
		p1 := number.GCD(number.Pow(a, r/2)+1, N)

		// result
		fmt.Printf("i=%d: N=%d, a=%d. p=%v, q=%v. s/r=%d/%d (%v=%.3f)\n", i, N, a, p0, p1, s, r, m, d)

		// check
		for _, p := range []int{p0, p1} {
			if 1 < p && p < N && N%p == 0 {
				fmt.Printf("answer: p=%v, q=%v\n", p, N/p)
				return
			}
		}
	}
}

func TestQsimFactoring85(t *testing.T) {
	N := 85
	a := 3 // 3, 6, 7, 11, 12, 14, 22, 23, 24, 27, 28, 29, 31, 37, 39, 41, 44, 46, 48, 54, 56, 57, 58, 61, 62, 63, 71, 73, 74, 78, 79, 82

	// co-prime
	if number.GCD(N, a) != 1 {
		t.Fatalf("gcd(%d, %d) != 1\n", N, a)
	}

	qsim := New()

	// initial state
	q0 := qsim.Zero()
	q1 := qsim.Zero()
	q2 := qsim.Zero()
	q3 := qsim.Zero()

	q4 := qsim.Zero()
	q5 := qsim.Zero()
	q6 := qsim.Zero()
	q7 := qsim.Zero()

	// superposition
	qsim.H(q0, q1, q2, q3)

	// Controlled-U
	qsim.CNOT(q0, q4)
	qsim.CNOT(q1, q5)
	qsim.CNOT(q2, q6)
	qsim.CNOT(q3, q7)

	// inverse QFT
	qsim.Swap(q0, q1, q2, q3)

	qsim.H(q3)
	qsim.CR(q3, q2, 2).H(q2)
	qsim.CR(q3, q1, 3).CR(q2, q1, 2).H(q1)
	qsim.CR(q3, q0, 4).CR(q2, q0, 3).CR(q1, q0, 2).H(q0)

	for i := 0; i < 10; i++ {
		// measure
		m := qsim.Clone().MeasureAsBinary(q0, q1, q2, q3)

		// find s/r
		d := number.BinaryFraction(m)
		_, s, r := number.ContinuedFraction(d)

		if number.IsOdd(r) || number.Pow(a, r/2)%N == -1 {
			continue
		}

		p0 := number.GCD(number.Pow(a, r/2)-1, N)
		p1 := number.GCD(number.Pow(a, r/2)+1, N)

		// result
		fmt.Printf("i=%d: N=%d, a=%d. p=%v, q=%v. s/r=%d/%d (%v=%.3f)\n", i, N, a, p0, p1, s, r, m, d)

		// check
		for _, p := range []int{p0, p1} {
			if 1 < p && p < N && N%p == 0 {
				fmt.Printf("answer: p=%v, q=%v\n", p, N/p)
				return
			}
		}
	}
}

func TestQsimFactoring51(t *testing.T) {
	N := 51
	a := 5 // 5, 7, 10, 11, 14, 20, 22, 23, 28, 29, 31, 37, 40, 41, 44, 46

	// co-prime
	if number.GCD(N, a) != 1 {
		t.Fatalf("gcd(%d, %d) != 1\n", N, a)
	}

	qsim := New()

	// initial state
	q0 := qsim.Zero()
	q1 := qsim.Zero()
	q2 := qsim.Zero()
	q3 := qsim.Zero()

	q4 := qsim.Zero()
	q5 := qsim.Zero()
	q6 := qsim.Zero()
	q7 := qsim.Zero()

	// superposition
	qsim.H(q0, q1, q2, q3)

	// Controlled-U
	qsim.CNOT(q0, q4)
	qsim.CNOT(q1, q5)
	qsim.CNOT(q2, q6)
	qsim.CNOT(q3, q7)

	// inverse QFT
	qsim.Swap(q0, q1, q2, q3)
	qsim.H(q3)
	qsim.CR(q3, q2, 2)
	qsim.H(q2)
	qsim.CR(q3, q1, 3)
	qsim.CR(q2, q1, 2)
	qsim.H(q1)
	qsim.CR(q3, q0, 4)
	qsim.CR(q2, q0, 3)
	qsim.CR(q1, q0, 2)
	qsim.H(q0)

	for i := 0; i < 10; i++ {
		// measure
		m := qsim.Clone().MeasureAsBinary(q0, q1, q2, q3)

		// find s/r
		d := number.BinaryFraction(m)
		_, s, r := number.ContinuedFraction(d)

		if number.IsOdd(r) || number.Pow(a, r/2)%N == -1 {
			continue
		}

		p0 := number.GCD(number.Pow(a, r/2)-1, N)
		p1 := number.GCD(number.Pow(a, r/2)+1, N)

		// result
		fmt.Printf("i=%d: N=%d, a=%d. p=%v, q=%v. s/r=%d/%d (%v=%.3f)\n", i, N, a, p0, p1, s, r, m, d)

		// check
		for _, p := range []int{p0, p1} {
			if 1 < p && p < N && N%p == 0 {
				fmt.Printf("answer: p=%v, q=%v\n", p, N/p)
				return
			}
		}
	}
}

func TestQSimFactoring15(t *testing.T) {
	N := 15
	a := 7

	// co-prime
	if number.GCD(N, a) != 1 {
		t.Fatalf("gcd(%d, %d) != 1\n", N, a)
	}

	qsim := New()

	// initial state
	q0 := qsim.Zero()
	q1 := qsim.Zero()
	q2 := qsim.Zero()

	q3 := qsim.Zero()
	q4 := qsim.Zero()
	q5 := qsim.Zero()
	q6 := qsim.One()

	// superposition
	qsim.H(q0, q1, q2)

	// Controlled-U^(2^0)
	qsim.CNOT(q2, q4)
	qsim.CNOT(q2, q5)

	// Controlled-U^(2^1)
	qsim.CNOT(q3, q5).CCNOT(q1, q5, q3).CNOT(q3, q5)
	qsim.CNOT(q4, q6).CCNOT(q1, q6, q4).CNOT(q4, q6)

	// inverse QFT
	qsim.Swap(q0, q2)
	qsim.InverseQFT(q0, q1, q2)

	for i := 0; i < 10; i++ {
		// measure q0, q1, q2
		m := qsim.Clone().MeasureAsBinary(q0, q1, q2)

		// find s/r. 010 -> 0.25 -> 1/4, 110 -> 0.75 -> 3/4, ...
		d := number.BinaryFraction(m)
		_, s, r := number.ContinuedFraction(d)

		// if r is odd, algorithm is failed
		if number.IsOdd(r) || number.Pow(a, r/2)%N == -1 {
			continue
		}

		// gcd(a^(r/2)-1, N), gcd(a^(r/2)+1, N)
		p0 := number.GCD(number.Pow(a, r/2)-1, N)
		p1 := number.GCD(number.Pow(a, r/2)+1, N)

		// result
		fmt.Printf("i=%d: N=%d, a=%d. p=%v, q=%v. s/r=%d/%d (%v=%.3f)\n", i, N, a, p0, p1, s, r, m, d)

		// check non-trivial factor
		for _, p := range []int{p0, p1} {
			if 1 < p && p < N && N%p == 0 {
				fmt.Printf("answer: p=%v, q=%v\n", p, N/p)
				return
			}
		}
	}
}

func ExampleQ_CCCNOT_grover_4qubit() {
	qsim := New()

	// initial state
	q0 := qsim.Zero()
	q1 := qsim.Zero()
	q2 := qsim.Zero()
	q3 := qsim.Zero()

	// superposition
	qsim.H(q0, q1, q2, q3)

	// iteration
	N := number.Pow(2, qsim.NumberOfBit())
	r := math.Floor(math.Pi / 4 * math.Sqrt(float64(N)))
	for i := 0; i < int(r); i++ {
		qsim.X(q2, q3)
		qsim.H(q3).CCCNOT(q0, q1, q2, q3).H(q3)
		qsim.X(q2, q3)

		qsim.H(q0, q1, q2, q3)
		qsim.X(q0, q1, q2, q3)
		qsim.H(q3).CCCNOT(q0, q1, q2, q3).H(q3)
		qsim.X(q0, q1, q2, q3)
		qsim.H(q0, q1, q2, q3)
	}

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [0000][  0]( 0.0508 0.0000i): 0.0026
	// [0001][  1]( 0.0508 0.0000i): 0.0026
	// [0010][  2]( 0.0508 0.0000i): 0.0026
	// [0011][  3]( 0.0508 0.0000i): 0.0026
	// [0100][  4]( 0.0508 0.0000i): 0.0026
	// [0101][  5]( 0.0508 0.0000i): 0.0026
	// [0110][  6]( 0.0508 0.0000i): 0.0026
	// [0111][  7]( 0.0508 0.0000i): 0.0026
	// [1000][  8]( 0.0508 0.0000i): 0.0026
	// [1001][  9]( 0.0508 0.0000i): 0.0026
	// [1010][ 10]( 0.0508 0.0000i): 0.0026
	// [1011][ 11]( 0.0508 0.0000i): 0.0026
	// [1100][ 12](-0.9805 0.0000i): 0.9613
	// [1101][ 13]( 0.0508 0.0000i): 0.0026
	// [1110][ 14]( 0.0508 0.0000i): 0.0026
	// [1111][ 15]( 0.0508 0.0000i): 0.0026
}

func ExampleQ_CCCNOT_grover_3qubit() {
	qsim := New()

	// initial state
	q0 := qsim.Zero()
	q1 := qsim.Zero()
	q2 := qsim.Zero()
	q3 := qsim.One()

	// superposition
	qsim.H(q0, q1, q2, q3)

	// oracle for |011>|1>
	qsim.X(q0).CCCNOT(q0, q1, q2, q3).X(q0)

	// amplification
	qsim.H(q0, q1, q2, q3)
	qsim.X(q0, q1, q2).CCZ(q0, q1, q2).X(q0, q1, q2)
	qsim.H(q0, q1, q2)

	for _, s := range qsim.State([]Qubit{q0, q1, q2}, []Qubit{q3}) {
		if s.Probability == 0 {
			continue
		}

		fmt.Println(s)
	}

	// Output:
	// [000 1][  0   1](-0.1768 0.0000i): 0.0313
	// [001 1][  1   1](-0.1768 0.0000i): 0.0313
	// [010 1][  2   1](-0.1768 0.0000i): 0.0313
	// [011 1][  3   1](-0.8839 0.0000i): 0.7813
	// [100 1][  4   1](-0.1768 0.0000i): 0.0313
	// [101 1][  5   1](-0.1768 0.0000i): 0.0313
	// [110 1][  6   1](-0.1768 0.0000i): 0.0313
	// [111 1][  7   1](-0.1768 0.0000i): 0.0313
}

func ExampleQ_Apply_bellstate() {
	qsim := New()

	i0 := qsim.Zero().Index()
	i1 := qsim.Zero().Index()
	n := qsim.NumberOfBit()

	g0 := gate.H().TensorProduct(gate.I())
	g1 := gate.CNOT(n, i0, i1)
	qc := g0.Apply(g1)

	qsim.Apply(qc)

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [00][  0]( 0.7071 0.0000i): 0.5000
	// [11][  3]( 0.7071 0.0000i): 0.5000
}

func ExampleQ_CR_qft_3qubit() {
	qsim := New()

	q0 := qsim.Zero()
	q1 := qsim.One()
	q2 := qsim.Zero()

	qsim.H(q0)
	qsim.CR(q1, q0, 2)
	qsim.CR(q2, q0, 3)
	qsim.H(q1)
	qsim.CR(q2, q1, 2)
	qsim.H(q2)

	qsim.Swap(q0, q2)

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [000][  0]( 0.3536 0.0000i): 0.1250
	// [001][  1]( 0.0000 0.3536i): 0.1250
	// [010][  2](-0.3536 0.0000i): 0.1250
	// [011][  3]( 0.0000-0.3536i): 0.1250
	// [100][  4]( 0.3536 0.0000i): 0.1250
	// [101][  5]( 0.0000 0.3536i): 0.1250
	// [110][  6](-0.3536 0.0000i): 0.1250
	// [111][  7]( 0.0000-0.3536i): 0.1250
}

func ExampleQ_CNOT_bellstate() {
	qsim := New()

	q0 := qsim.Zero()
	q1 := qsim.Zero()

	qsim.H(q0).CNOT(q0, q1)

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [00][  0]( 0.7071 0.0000i): 0.5000
	// [11][  3]( 0.7071 0.0000i): 0.5000
}

func ExampleQ_CZ_teleportation() {
	qsim := New()

	phi := qsim.New(1, 2)
	q0 := qsim.Zero()
	q1 := qsim.Zero()

	for _, s := range qsim.State([]Qubit{phi}) {
		fmt.Println(s)
	}

	qsim.H(q0).CNOT(q0, q1)
	qsim.CNOT(phi, q0).H(phi)

	qsim.CNOT(q0, q1)
	qsim.CZ(phi, q1)

	qsim.Measure(phi, q0)

	for _, s := range qsim.State([]Qubit{q1}) {
		fmt.Println(s)
	}

	// Output:
	// [0][  0]( 0.4472 0.0000i): 0.2000
	// [1][  1]( 0.8944 0.0000i): 0.8000
	// [0][  0]( 0.4472 0.0000i): 0.2000
	// [1][  1]( 0.8944 0.0000i): 0.8000
}

func ExampleQ_ConditionZ_teleportation() {
	qsim := New()

	phi := qsim.New(1, 2)
	q0 := qsim.Zero()
	q1 := qsim.Zero()

	for _, s := range qsim.State([]Qubit{phi}) {
		fmt.Println(s)
	}

	qsim.H(q0).CNOT(q0, q1)
	qsim.CNOT(phi, q0).H(phi)

	mz := qsim.Measure(phi)
	mx := qsim.Measure(q0)

	qsim.ConditionX(mx.IsOne(), q1)
	qsim.ConditionZ(mz.IsOne(), q1)

	for _, s := range qsim.State([]Qubit{q1}) {
		fmt.Println(s)
	}

	// Output:
	// [0][  0]( 0.4472 0.0000i): 0.2000
	// [1][  1]( 0.8944 0.0000i): 0.8000
	// [0][  0]( 0.4472 0.0000i): 0.2000
	// [1][  1]( 0.8944 0.0000i): 0.8000
}

func ExampleQ_ConditionX_error_correction() {
	qsim := New()

	q0 := qsim.New(1, 2)

	for _, s := range qsim.State([]Qubit{q0}) {
		fmt.Println(s)
	}

	// encoding
	q1 := qsim.Zero()
	q2 := qsim.Zero()
	qsim.CNOT(q0, q1).CNOT(q0, q2)

	// error: first qubit is flipped
	qsim.X(q0)

	// add ancilla qubit
	q3 := qsim.Zero()
	q4 := qsim.Zero()

	// error correction
	qsim.CNOT(q0, q3).CNOT(q1, q3)
	qsim.CNOT(q1, q4).CNOT(q2, q4)

	m3 := qsim.Measure(q3)
	m4 := qsim.Measure(q4)

	qsim.ConditionX(m3.IsOne() && m4.IsZero(), q0)
	qsim.ConditionX(m3.IsOne() && m4.IsOne(), q1)
	qsim.ConditionX(m3.IsZero() && m4.IsOne(), q2)

	// decoding
	qsim.CNOT(q0, q2).CNOT(q0, q1)

	for _, s := range qsim.State([]Qubit{q0}) {
		fmt.Println(s)
	}

	// Output:
	// [0][  0]( 0.4472 0.0000i): 0.2000
	// [1][  1]( 0.8944 0.0000i): 0.8000
	// [0][  0]( 0.4472 0.0000i): 0.2000
	// [1][  1]( 0.8944 0.0000i): 0.8000
}

func TestPOVM(t *testing.T) {
	E1 := gate.New(
		[]complex128{0, 0},
		[]complex128{0, 1},
	).Mul(complex(math.Sqrt(2)/(1.0+math.Sqrt(2)), 0))

	E2 := gate.New(
		[]complex128{1, -1},
		[]complex128{-1, 1},
	).Mul(complex(0.5, 0)).
		Mul(complex(math.Sqrt(2)/(1.0+math.Sqrt(2)), 0))

	E3 := gate.I().Sub(E1).Sub(E2)

	if !E1.Add(E2).Add(E3).Equals(gate.I()) {
		t.Fail()
	}

	q0 := qubit.Zero()
	if q0.Apply(E1).InnerProduct(q0) != complex(0, 0) {
		t.Fail()
	}

	q1 := qubit.Zero().Apply(gate.H())
	if q1.Apply(E2).InnerProduct(q1) != complex(0, 0) {
		t.Fail()
	}
}

func TestGrover3qubit(t *testing.T) {
	x := matrix.TensorProduct(gate.X(), gate.I(3))
	oracle := x.Apply(gate.ControlledNot(4, []int{0, 1, 2}, 3)).Apply(x)

	h4 := matrix.TensorProduct(gate.H(3), gate.H())
	x3 := matrix.TensorProduct(gate.X(3), gate.I())
	cz := matrix.TensorProduct(gate.ControlledZ(3, []int{0, 1}, 2), gate.I())
	h3 := matrix.TensorProduct(gate.H(3), gate.I())
	amp := h4.Apply(x3).Apply(cz).Apply(x3).Apply(h3)

	q := qubit.TensorProduct(qubit.Zero(3), qubit.One())
	q.Apply(gate.H(4)).Apply(oracle).Apply(amp)

	// q3 is always |1>
	q3 := q.Measure(3)
	if !q3.IsOne() {
		t.Error(q3)
	}

	p := q.Probability()
	for i, pp := range p {
		// |011>|1>
		if i == 7 {
			if math.Abs(pp-0.78125) > 1e-13 {
				t.Error(q.Probability())
			}
			continue
		}

		if i%2 == 0 {
			if math.Abs(pp) > 1e-13 {
				t.Error(q.Probability())
			}
			continue
		}

		if math.Abs(pp-0.03125) > 1e-13 {
			t.Error(q.Probability())
		}
	}
}

func TestGrover2qubit(t *testing.T) {
	oracle := gate.CZ(2, 0, 1)

	h2 := gate.H(2)
	x2 := gate.X(2)
	amp := h2.Apply(x2).Apply(gate.CZ(2, 0, 1)).Apply(x2).Apply(h2)

	qc := h2.Apply(oracle).Apply(amp)
	q := qubit.Zero(2).Apply(qc)

	q.Measure(0)
	q.Measure(1)

	p := q.Probability()
	if math.Abs(p[3]-1) > 1e-13 {
		t.Error(p)
	}
}

func TestQuantumTeleportation(t *testing.T) {
	g0 := matrix.TensorProduct(gate.H(), gate.I())
	g1 := gate.CNOT(2, 0, 1)
	bell := qubit.Zero(2).Apply(g0).Apply(g1)

	phi := qubit.New(1, 2)
	phi.TensorProduct(bell)

	g2 := gate.CNOT(3, 0, 1)
	g3 := matrix.TensorProduct(gate.H(), gate.I(2))
	phi.Apply(g2).Apply(g3)

	mz := phi.Measure(0)
	mx := phi.Measure(1)

	if mz.IsOne() {
		z := matrix.TensorProduct(gate.I(2), gate.Z())
		phi.Apply(z)
	}

	if mx.IsOne() {
		x := matrix.TensorProduct(gate.I(2), gate.X())
		phi.Apply(x)
	}

	var test = []struct {
		zero int
		one  int
		zval float64
		oval float64
		eps  float64
		mz   *qubit.Qubit
		mx   *qubit.Qubit
	}{
		{0, 1, 0.2, 0.8, 1e-13, qubit.Zero(), qubit.Zero()},
		{2, 3, 0.2, 0.8, 1e-13, qubit.Zero(), qubit.One()},
		{4, 5, 0.2, 0.8, 1e-13, qubit.One(), qubit.Zero()},
		{6, 7, 0.2, 0.8, 1e-13, qubit.One(), qubit.One()},
	}

	p := phi.Probability()
	if math.Abs(number.Sum(p)-1) > 1e-13 {
		t.Error(p)
	}

	for _, tt := range test {
		if p[tt.zero] == 0 {
			continue
		}

		if math.Abs(p[tt.zero]-tt.zval) > tt.eps {
			t.Error(p)
		}
		if math.Abs(p[tt.one]-tt.oval) > tt.eps {
			t.Error(p)
		}

		if !mz.Equals(tt.mz) {
			t.Error(p)
		}

		if !mx.Equals(tt.mx) {
			t.Error(p)
		}
	}
}

func TestQuantumTeleportation2(t *testing.T) {
	g0 := matrix.TensorProduct(gate.H(), gate.I())
	g1 := gate.CNOT(2, 0, 1)
	bell := qubit.Zero(2).Apply(g0).Apply(g1)

	phi := qubit.New(1, 2)
	phi.TensorProduct(bell)

	g2 := gate.CNOT(3, 0, 1)
	g3 := matrix.TensorProduct(gate.H(), gate.I(2))
	g4 := gate.CNOT(3, 1, 2)
	g5 := gate.CZ(3, 0, 2)

	phi.Apply(g2).Apply(g3).Apply(g4).Apply(g5)

	mz := phi.Measure(0)
	mx := phi.Measure(1)

	var test = []struct {
		zero int
		one  int
		zval float64
		oval float64
		eps  float64
		mz   *qubit.Qubit
		mx   *qubit.Qubit
	}{
		{0, 1, 0.2, 0.8, 1e-13, qubit.Zero(), qubit.Zero()},
		{2, 3, 0.2, 0.8, 1e-13, qubit.Zero(), qubit.One()},
		{4, 5, 0.2, 0.8, 1e-13, qubit.One(), qubit.Zero()},
		{6, 7, 0.2, 0.8, 1e-13, qubit.One(), qubit.One()},
	}

	p := phi.Probability()
	if math.Abs(number.Sum(p)-1) > 1e-13 {
		t.Error(p)
	}

	for _, tt := range test {
		if p[tt.zero] == 0 {
			continue
		}

		if math.Abs(p[tt.zero]-tt.zval) > tt.eps {
			t.Error(p)
		}

		if math.Abs(p[tt.one]-tt.oval) > tt.eps {
			t.Error(p)
		}

		if !mz.Equals(tt.mz) {
			t.Error(p)
		}

		if !mx.Equals(tt.mx) {
			t.Error(p)
		}
	}
}

func TestErrorCorrectionZero(t *testing.T) {
	phi := qubit.Zero()

	// encoding
	phi.TensorProduct(qubit.Zero(2))
	phi.Apply(gate.CNOT(3, 0, 1))
	phi.Apply(gate.CNOT(3, 0, 2))

	// error: first qubit is flipped
	phi.Apply(matrix.TensorProduct(gate.X(), gate.I(2)))

	// add ancilla qubit
	phi.TensorProduct(qubit.Zero(2))

	// z1z2
	phi.Apply(gate.CNOT(5, 0, 3)).Apply(gate.CNOT(5, 1, 3))

	// z2z3
	phi.Apply(gate.CNOT(5, 1, 4)).Apply(gate.CNOT(5, 2, 4))

	// measure
	m3 := phi.Measure(3)
	m4 := phi.Measure(4)

	// recover
	if m3.IsOne() && m4.IsZero() {
		phi.Apply(matrix.TensorProduct(gate.X(), gate.I(4)))
	}

	if m3.IsOne() && m4.IsOne() {
		phi.Apply(matrix.TensorProduct(gate.I(), gate.X(), gate.I(3)))
	}

	if m3.IsZero() && m4.IsOne() {
		phi.Apply(matrix.TensorProduct(gate.I(2), gate.X(), gate.I(2)))
	}

	// answer is |000>|10>
	if phi.Probability()[2] != 1 {
		t.Error(phi.Probability())
	}
}

func TestErrorCorrectionOne(t *testing.T) {
	phi := qubit.One()

	// encoding
	phi.TensorProduct(qubit.Zero(2))
	phi.Apply(gate.CNOT(3, 0, 1))
	phi.Apply(gate.CNOT(3, 0, 2))

	// error: first qubit is flipped
	phi.Apply(matrix.TensorProduct(gate.X(), gate.I(2)))

	// add ancilla qubit
	phi.TensorProduct(qubit.Zero(2))

	// z1z2
	phi.Apply(gate.CNOT(5, 0, 3)).Apply(gate.CNOT(5, 1, 3))

	// z2z3
	phi.Apply(gate.CNOT(5, 1, 4)).Apply(gate.CNOT(5, 2, 4))

	// measure
	m3 := phi.Measure(3)
	m4 := phi.Measure(4)

	// recover
	if m3.IsOne() && m4.IsZero() {
		phi.Apply(matrix.TensorProduct(gate.X(), gate.I(4)))
	}

	if m3.IsOne() && m4.IsOne() {
		phi.Apply(matrix.TensorProduct(gate.I(), gate.X(), gate.I(3)))
	}

	if m3.IsZero() && m4.IsOne() {
		phi.Apply(matrix.TensorProduct(gate.I(2), gate.X(), gate.I(2)))
	}

	// answer is |111>|10>
	if phi.Probability()[30] != 1 {
		t.Error(phi.Probability())
	}
}

func TestErrorCorrectionBitFlip1(t *testing.T) {
	phi := qubit.New(1, 2)

	// encoding
	phi.TensorProduct(qubit.Zero(2))
	phi.Apply(gate.CNOT(3, 0, 1))
	phi.Apply(gate.CNOT(3, 0, 2))

	// error: first qubit is flipped
	phi.Apply(matrix.TensorProduct(gate.X(), gate.I(2)))

	// add ancilla qubit
	phi.TensorProduct(qubit.Zero(2))

	// z1z2
	phi.Apply(gate.CNOT(5, 0, 3)).Apply(gate.CNOT(5, 1, 3))

	// z2z3
	phi.Apply(gate.CNOT(5, 1, 4)).Apply(gate.CNOT(5, 2, 4))

	// measure
	m3 := phi.Measure(3)
	m4 := phi.Measure(4)

	// recover
	if m3.IsOne() && m4.IsZero() {
		phi.Apply(matrix.TensorProduct(gate.X(), gate.I(4)))
	}

	if m3.IsOne() && m4.IsOne() {
		phi.Apply(matrix.TensorProduct(gate.I(), gate.X(), gate.I(3)))
	}

	if m3.IsZero() && m4.IsOne() {
		phi.Apply(matrix.TensorProduct(gate.I(2), gate.X(), gate.I(2)))
	}

	// answer is 0.2|000>|10> + 0.8|111>|10>
	p := phi.Probability()
	if math.Abs(p[2]-0.2) > 1e-13 {
		t.Error(p)
	}

	if math.Abs(p[30]-0.8) > 1e-13 {
		t.Error(p)
	}
}

func TestErrorCorrectionBitFlip2(t *testing.T) {
	phi := qubit.New(1, 2)

	// encoding
	phi.TensorProduct(qubit.Zero(2))
	phi.Apply(gate.CNOT(3, 0, 1))
	phi.Apply(gate.CNOT(3, 0, 2))

	// error: second qubit is flipped
	phi.Apply(matrix.TensorProduct(gate.I(), gate.X(), gate.I()))

	// add ancilla qubit
	phi.TensorProduct(qubit.Zero(2))

	// z1z2
	phi.Apply(gate.CNOT(5, 0, 3)).Apply(gate.CNOT(5, 1, 3))

	// z2z3
	phi.Apply(gate.CNOT(5, 1, 4)).Apply(gate.CNOT(5, 2, 4))

	// measure
	m3 := phi.Measure(3)
	m4 := phi.Measure(4)

	// recover
	if m3.IsOne() && m4.IsZero() {
		phi.Apply(matrix.TensorProduct(gate.X(), gate.I(4)))
	}

	if m3.IsOne() && m4.IsOne() {
		phi.Apply(matrix.TensorProduct(gate.I(), gate.X(), gate.I(3)))
	}

	if m3.IsZero() && m4.IsOne() {
		phi.Apply(matrix.TensorProduct(gate.I(2), gate.X(), gate.I(2)))
	}

	// answer is 0.2|000>|11> + 0.8|111>|11>
	p := phi.Probability()
	if math.Abs(p[3]-0.2) > 1e-13 {
		t.Error(p)
	}

	if math.Abs(p[31]-0.8) > 1e-13 {
		t.Error(p)
	}
}

func TestErrorCorrectionBitFlip3(t *testing.T) {
	phi := qubit.New(1, 2)

	// encoding
	phi.TensorProduct(qubit.Zero(2))
	phi.Apply(gate.CNOT(3, 0, 1))
	phi.Apply(gate.CNOT(3, 0, 2))

	// error: third qubit is flipped
	phi.Apply(matrix.TensorProduct(gate.I(), gate.I(), gate.X()))

	// add ancilla qubit
	phi.TensorProduct(qubit.Zero(2))

	// z1z2
	phi.Apply(gate.CNOT(5, 0, 3)).Apply(gate.CNOT(5, 1, 3))

	// z2z3
	phi.Apply(gate.CNOT(5, 1, 4)).Apply(gate.CNOT(5, 2, 4))

	// measure
	m3 := phi.Measure(3)
	m4 := phi.Measure(4)

	// recover
	if m3.IsOne() && m4.IsZero() {
		phi.Apply(matrix.TensorProduct(gate.X(), gate.I(4)))
	}

	if m3.IsOne() && m4.IsOne() {
		phi.Apply(matrix.TensorProduct(gate.I(), gate.X(), gate.I(3)))
	}

	if m3.IsZero() && m4.IsOne() {
		phi.Apply(matrix.TensorProduct(gate.I(2), gate.X(), gate.I(2)))
	}

	// answer is 0.2|000>|01> + 0.8|111>|01>
	p := phi.Probability()
	if math.Abs(p[1]-0.2) > 1e-13 {
		t.Error(p)
	}

	if math.Abs(p[29]-0.8) > 1e-13 {
		t.Error(p)
	}
}

func TestErrorCorrectionPhaseFlip1(t *testing.T) {
	phi := qubit.New(1, 2)

	// encoding
	phi.TensorProduct(qubit.Zero(2))
	phi.Apply(gate.CNOT(3, 0, 1))
	phi.Apply(gate.CNOT(3, 0, 2))
	phi.Apply(gate.H(3))

	// error: first qubit is flipped
	phi.Apply(matrix.TensorProduct(gate.Z(), gate.I(2)))

	// H
	phi.Apply(gate.H(3))

	// add ancilla qubit
	phi.TensorProduct(qubit.Zero(2))

	// x1x2
	phi.Apply(gate.CNOT(5, 0, 3)).Apply(gate.CNOT(5, 1, 3))

	// x2x3
	phi.Apply(gate.CNOT(5, 1, 4)).Apply(gate.CNOT(5, 2, 4))

	// H
	phi.Apply(matrix.TensorProduct(gate.H(3), gate.I(2)))

	// measure
	m3 := phi.Measure(3)
	m4 := phi.Measure(4)

	// recover
	if m3.IsOne() && m4.IsZero() {
		phi.Apply(matrix.TensorProduct(gate.Z(), gate.I(4)))
	}

	if m3.IsOne() && m4.IsOne() {
		phi.Apply(matrix.TensorProduct(gate.I(), gate.Z(), gate.I(3)))
	}

	if m3.IsZero() && m4.IsOne() {
		phi.Apply(matrix.TensorProduct(gate.I(2), gate.Z(), gate.I(2)))
	}

	phi.Apply(matrix.TensorProduct(gate.H(3), gate.I(2)))

	p := phi.Probability()
	if math.Abs(p[2]-0.2) > 1e-13 {
		t.Error(p)
	}

	if math.Abs(p[30]-0.8) > 1e-13 {
		t.Error(p)
	}
}
