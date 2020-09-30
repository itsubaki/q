package q_test

import (
	"fmt"
	"math"

	"github.com/itsubaki/q"
	"github.com/itsubaki/q/pkg/math/number"
	"github.com/itsubaki/q/pkg/math/rand"
	"github.com/itsubaki/q/pkg/quantum/gate"
)

func ExampleQ_Zero() {
	qsim := q.New()

	q0 := qsim.Zero()
	q1 := qsim.Zero()
	qsim.H(q0, q1)

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [00][  0]( 0.5000 0.0000i): 0.2500
	// [01][  1]( 0.5000 0.0000i): 0.2500
	// [10][  2]( 0.5000 0.0000i): 0.2500
	// [11][  3]( 0.5000 0.0000i): 0.2500
}

func ExampleQ_ZeroWith() {
	qsim := q.New()

	r := qsim.ZeroWith(2)
	qsim.H(r...)

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [00][  0]( 0.5000 0.0000i): 0.2500
	// [01][  1]( 0.5000 0.0000i): 0.2500
	// [10][  2]( 0.5000 0.0000i): 0.2500
	// [11][  3]( 0.5000 0.0000i): 0.2500
}

func ExampleQ_ZeroLog2() {
	qsim := q.New()

	r := qsim.ZeroLog2(3)
	qsim.H(r...)

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [00][  0]( 0.5000 0.0000i): 0.2500
	// [01][  1]( 0.5000 0.0000i): 0.2500
	// [10][  2]( 0.5000 0.0000i): 0.2500
	// [11][  3]( 0.5000 0.0000i): 0.2500
}

func ExampleQ_Amplitude() {
	qsim := q.New()

	q0 := qsim.Zero()
	q1 := qsim.Zero()

	qsim.H(q0)
	qsim.CNOT(q0, q1)

	for _, a := range qsim.Amplitude() {
		fmt.Println(a)
	}

	// Output:
	// (0.7071067811865476+0i)
	// (0+0i)
	// (0+0i)
	// (0.7071067811865476+0i)
}

func ExampleQ_Probability() {
	qsim := q.New()

	q0 := qsim.Zero()
	q1 := qsim.Zero()

	qsim.H(q0)
	qsim.CNOT(q0, q1)

	for _, p := range qsim.Probability() {
		fmt.Printf("%.4f\n", p)
	}

	// Output:
	// 0.5000
	// 0.0000
	// 0.0000
	// 0.5000
}

func ExampleQ_BinaryString() {
	qsim := q.New()
	qsim.Seed = []int64{1}
	qsim.Rand = rand.Math

	r := qsim.ZeroWith(4)
	qsim.H(r...)

	b := qsim.BinaryString(r...)
	fmt.Println(b)

	// Output:
	// 1111
}

func ExampleQ_Measure() {
	qsim := q.New()
	qsim.Seed = []int64{1}
	qsim.Rand = rand.Math

	q0 := qsim.Zero()
	qsim.H(q0)

	m := qsim.Measure(q0)
	fmt.Println(m)

	// Output:
	// [(0+0i) (1+0i)]
}

func ExampleQ_MeasureAsInt() {
	qsim := q.New()
	qsim.Seed = []int64{1}
	qsim.Rand = rand.Math

	q0 := qsim.Zero()
	q1 := qsim.Zero()
	q2 := qsim.Zero()
	qsim.H(q0, q1, q2)

	m := qsim.MeasureAsInt(q0, q1, q2)
	fmt.Println(m)

	// Output:
	// 7
}

func ExampleQ_MeasureAsBinary() {
	qsim := q.New()
	qsim.Seed = []int64{1}
	qsim.Rand = rand.Math

	q0 := qsim.Zero()
	q1 := qsim.Zero()
	q2 := qsim.Zero()
	qsim.H(q0, q1, q2)

	m := qsim.MeasureAsBinary(q0, q1, q2)
	fmt.Println(m)

	// Output:
	// [1 1 1]
}

func ExampleQ_ControlledModExp2_mod21() {
	qsim := q.New()

	q0 := qsim.One()
	r1 := qsim.ZeroLog2(21)

	qsim.X(r1[len(r1)-1])
	for _, s := range qsim.State([]q.Qubit{q0}, r1) {
		fmt.Println(s)
	}

	// 2^2^0 * 1 mod 21 = 2
	qsim.ControlledModExp2(2, 0, 21, q0, r1)
	for _, s := range qsim.State([]q.Qubit{q0}, r1) {
		fmt.Println(s)
	}

	// 2^2^1 * 2 mod 21 = 8
	qsim.ControlledModExp2(2, 1, 21, q0, r1)
	for _, s := range qsim.State([]q.Qubit{q0}, r1) {
		fmt.Println(s)
	}

	// 2^2^2 * 8 mod 21 = 2
	qsim.ControlledModExp2(2, 2, 21, q0, r1)
	for _, s := range qsim.State([]q.Qubit{q0}, r1) {
		fmt.Println(s)
	}

	// 2^2^3 * 2 mod 21 = 8
	qsim.ControlledModExp2(2, 3, 21, q0, r1)
	for _, s := range qsim.State([]q.Qubit{q0}, r1) {
		fmt.Println(s)
	}

	// Output:
	// [1 00001][  1   1]( 1.0000 0.0000i): 1.0000
	// [1 00010][  1   2]( 1.0000 0.0000i): 1.0000
	// [1 01000][  1   8]( 1.0000 0.0000i): 1.0000
	// [1 00010][  1   2]( 1.0000 0.0000i): 1.0000
	// [1 01000][  1   8]( 1.0000 0.0000i): 1.0000
}

func ExampleQ_ControlledModExp2_mod15() {
	qsim := q.New()

	q0 := qsim.One()
	r1 := qsim.ZeroLog2(15)

	// 1
	qsim.X(r1[len(r1)-1])
	for _, s := range qsim.State([]q.Qubit{q0}, r1) {
		fmt.Println(s)
	}

	// 7^2^0 * 1 mod 15 = 7
	qsim.ControlledModExp2(7, 0, 15, q0, r1)
	for _, s := range qsim.State([]q.Qubit{q0}, r1) {
		fmt.Println(s)
	}

	// 7^2^1 * 7 mod 15 = 13
	qsim.ControlledModExp2(7, 1, 15, q0, r1)
	for _, s := range qsim.State([]q.Qubit{q0}, r1) {
		fmt.Println(s)
	}

	// 7^2^2 * 13 mod 15 = 13
	qsim.ControlledModExp2(7, 2, 15, q0, r1)
	for _, s := range qsim.State([]q.Qubit{q0}, r1) {
		fmt.Println(s)
	}

	// Output:
	// [1 0001][  1   1]( 1.0000 0.0000i): 1.0000
	// [1 0111][  1   7]( 1.0000 0.0000i): 1.0000
	// [1 1101][  1  13]( 1.0000 0.0000i): 1.0000
	// [1 1101][  1  13]( 1.0000 0.0000i): 1.0000
}

func Example_shorFactoring21() {
	N := 21
	a := 5

	qsim := q.New()
	qsim.Seed = []int64{1}
	qsim.Rand = rand.Math

	r0 := qsim.ZeroWith(4)
	r1 := qsim.ZeroLog2(N)

	qsim.X(r1[len(r1)-1])
	qsim.H(r0...)
	qsim.CModExp2(a, N, r0, r1)
	qsim.InvQFT(r0...)

	m := qsim.MeasureAsBinary(r0...)
	d := number.BinaryFraction(m)
	_, s, r := number.ContinuedFraction(d)

	if number.IsOdd(r) || number.Pow(a, r/2)%N == -1 {
		return
	}

	p0 := number.GCD(number.Pow(a, r/2)-1, N)
	p1 := number.GCD(number.Pow(a, r/2)+1, N)

	for _, p := range []int{p0, p1} {
		if 1 < p && p < N && N%p == 0 {
			fmt.Printf("N=%d, a=%d. p=%v, q=%v. s/r=%d/%d (%v=%.3f)\n", N, a, p0, p1, s, r, m, d)
			fmt.Printf("answer: p=%v, q=%v\n", p, N/p)
			return
		}
	}

	// Output:
	// N=21, a=5. p=3, q=1. s/r=11/16 ([1 0 1 1]=0.688)
	// answer: p=3, q=7
}

func Example_shorFactoring85() {
	N := 85
	a := 3 // 3, 6, 7, 11, 12, 14, 22, 23, 24, 27, 28, 29, 31, 37, 39, 41, 44, 46, 48, 54, 56, 57, 58, 61, 62, 63, 71, 73, 74, 78, 79, 82

	qsim := q.New()
	qsim.Seed = []int64{1}
	qsim.Rand = rand.Math

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

	// measure
	m := qsim.MeasureAsBinary(q0, q1, q2, q3)

	// find s/r
	d := number.BinaryFraction(m)
	_, s, r := number.ContinuedFraction(d)

	if number.IsOdd(r) || number.Pow(a, r/2)%N == -1 {
		return
	}

	p0 := number.GCD(number.Pow(a, r/2)-1, N)
	p1 := number.GCD(number.Pow(a, r/2)+1, N)

	// check
	for _, p := range []int{p0, p1} {
		if 1 < p && p < N && N%p == 0 {
			fmt.Printf("N=%d, a=%d. p=%v, q=%v. s/r=%d/%d (%v=%.3f)\n", N, a, p0, p1, s, r, m, d)
			fmt.Printf("answer: p=%v, q=%v\n", p, N/p)
			return
		}
	}

	// Output:
	// N=85, a=3. p=5, q=17. s/r=15/16 ([1 1 1 1]=0.938)
	// answer: p=5, q=17
}

func Example_shorFactoring51() {
	N := 51
	a := 5 // 5, 7, 10, 11, 14, 20, 22, 23, 28, 29, 31, 37, 40, 41, 44, 46

	qsim := q.New()
	qsim.Seed = []int64{1}
	qsim.Rand = rand.Math

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

	// measure
	m := qsim.MeasureAsBinary(q0, q1, q2, q3)

	// find s/r
	d := number.BinaryFraction(m)
	_, s, r := number.ContinuedFraction(d)

	if number.IsOdd(r) || number.Pow(a, r/2)%N == -1 {
		return
	}

	p0 := number.GCD(number.Pow(a, r/2)-1, N)
	p1 := number.GCD(number.Pow(a, r/2)+1, N)

	// check
	for _, p := range []int{p0, p1} {
		if 1 < p && p < N && N%p == 0 {
			fmt.Printf("N=%d, a=%d. p=%v, q=%v. s/r=%d/%d (%v=%.3f)\n", N, a, p0, p1, s, r, m, d)
			fmt.Printf("answer: p=%v, q=%v\n", p, N/p)
			return
		}
	}

	// Output:
	// N=51, a=5. p=3, q=17. s/r=15/16 ([1 1 1 1]=0.938)
	// answer: p=3, q=17
}

func Example_shorFactoring15() {
	N := 15
	a := 7

	qsim := q.New()
	qsim.Seed = []int64{1}
	qsim.Rand = rand.Math

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
	qsim.InvQFT(q0, q1, q2)

	// measure q0, q1, q2
	m := qsim.MeasureAsBinary(q0, q1, q2)

	// find s/r. 010 -> 0.25 -> 1/4, 110 -> 0.75 -> 3/4, ...
	d := number.BinaryFraction(m)
	_, s, r := number.ContinuedFraction(d)

	// if r is odd, algorithm is failed
	if number.IsOdd(r) || number.Pow(a, r/2)%N == -1 {
		return
	}

	// gcd(a^(r/2)-1, N), gcd(a^(r/2)+1, N)
	p0 := number.GCD(number.Pow(a, r/2)-1, N)
	p1 := number.GCD(number.Pow(a, r/2)+1, N)

	// check non-trivial factor
	for _, p := range []int{p0, p1} {
		if 1 < p && p < N && N%p == 0 {
			fmt.Printf("N=%d, a=%d. p=%v, q=%v. s/r=%d/%d (%v=%.3f)\n", N, a, p0, p1, s, r, m, d)
			fmt.Printf("answer: p=%v, q=%v\n", p, N/p)
			return
		}
	}

	// Output:
	// N=15, a=7. p=3, q=5. s/r=3/4 ([1 1 0]=0.750)
	// answer: p=3, q=5
}

func Example_grover4qubit() {
	qsim := q.New()

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

func Example_grover3qubit() {
	qsim := q.New()

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

	for _, s := range qsim.State([]q.Qubit{q0, q1, q2}, []q.Qubit{q3}) {
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

func Example_qFT3qubit() {
	qsim := q.New()

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

func Example_bellState() {
	qsim := q.New()

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

func Example_quantumTeleportation() {
	qsim := q.New()

	phi := qsim.New(1, 2)
	q0 := qsim.Zero()
	q1 := qsim.Zero()

	fmt.Println("phi:")
	for _, s := range qsim.State([]q.Qubit{phi}) {
		fmt.Println(s)
	}

	qsim.H(q0).CNOT(q0, q1)
	qsim.CNOT(phi, q0).H(phi)

	qsim.CNOT(q0, q1)
	qsim.CZ(phi, q1)

	qsim.Measure(phi, q0)

	fmt.Println("q1:")
	for _, s := range qsim.State([]q.Qubit{q1}) {
		fmt.Println(s)
	}

	// Output:
	// phi:
	// [0][  0]( 0.4472 0.0000i): 0.2000
	// [1][  1]( 0.8944 0.0000i): 0.8000
	// q1:
	// [0][  0]( 0.4472 0.0000i): 0.2000
	// [1][  1]( 0.8944 0.0000i): 0.8000
}

func ExampleQ_Apply() {
	qsim := q.New()

	q0 := qsim.Zero()
	q1 := qsim.Zero()
	n := qsim.NumberOfBit()

	qsim.Apply(gate.H(), q0)
	qsim.Apply(gate.CNOT(n, q0.Index(), q1.Index()))

	for _, s := range qsim.State() {
		fmt.Println(s)
	}

	// Output:
	// [00][  0]( 0.7071 0.0000i): 0.5000
	// [11][  3]( 0.7071 0.0000i): 0.5000
}

func ExampleQ_ConditionZ_quantumTeleportation() {
	qsim := q.New()

	phi := qsim.New(1, 2)
	q0 := qsim.Zero()
	q1 := qsim.Zero()

	fmt.Println("phi:")
	for _, s := range qsim.State([]q.Qubit{phi}) {
		fmt.Println(s)
	}

	qsim.H(q0).CNOT(q0, q1)
	qsim.CNOT(phi, q0).H(phi)

	mz := qsim.Measure(phi)
	mx := qsim.Measure(q0)

	qsim.ConditionX(mx.IsOne(), q1)
	qsim.ConditionZ(mz.IsOne(), q1)

	fmt.Println("q1:")
	for _, s := range qsim.State([]q.Qubit{q1}) {
		fmt.Println(s)
	}

	// Output:
	// phi:
	// [0][  0]( 0.4472 0.0000i): 0.2000
	// [1][  1]( 0.8944 0.0000i): 0.8000
	// q1:
	// [0][  0]( 0.4472 0.0000i): 0.2000
	// [1][  1]( 0.8944 0.0000i): 0.8000
}

func ExampleQ_ConditionX_errorCorrection() {
	qsim := q.New()

	q0 := qsim.New(1, 2)

	fmt.Println("q0:")
	for _, s := range qsim.State([]q.Qubit{q0}) {
		fmt.Println(s)
	}

	// encoding
	q1 := qsim.Zero()
	q2 := qsim.Zero()
	qsim.CNOT(q0, q1).CNOT(q0, q2)

	// error: first qubit is flipped
	qsim.X(q0)

	fmt.Println("q0(flipped):")
	for _, s := range qsim.State([]q.Qubit{q0}) {
		fmt.Println(s)
	}

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

	fmt.Println("q0(corrected):")
	for _, s := range qsim.State([]q.Qubit{q0}) {
		fmt.Println(s)
	}

	// Output:
	// q0:
	// [0][  0]( 0.4472 0.0000i): 0.2000
	// [1][  1]( 0.8944 0.0000i): 0.8000
	// q0(flipped):
	// [0][  0]( 0.8944 0.0000i): 0.8000
	// [1][  1]( 0.4472 0.0000i): 0.2000
	// q0(corrected):
	// [0][  0]( 0.4472 0.0000i): 0.2000
	// [1][  1]( 0.8944 0.0000i): 0.8000
}
