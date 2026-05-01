# q

[![PkgGoDev](https://pkg.go.dev/badge/github.com/itsubaki/q)](https://pkg.go.dev/github.com/itsubaki/q)
[![Go Report Card](https://goreportcard.com/badge/github.com/itsubaki/q?style=flat-square)](https://goreportcard.com/report/github.com/itsubaki/q)
[![tests](https://github.com/itsubaki/q/workflows/tests/badge.svg)](https://github.com/itsubaki/q/actions)
[![codecov](https://codecov.io/gh/itsubaki/q/branch/main/graph/badge.svg?token=iNccCs1Tez)](https://codecov.io/gh/itsubaki/q)

A quantum computing simulator in Go using only the standard library.

## Installation

```shell
go get github.com/itsubaki/q@latest
```

## Examples

### Bell State

```go
qsim := q.New()

// generate qubits in the state |0>|0>
q0 := qsim.Zero()
q1 := qsim.Zero()

// apply the quantum circuit
qsim.H(q0)
qsim.CNOT(q0, q1)

for _, s := range qsim.State() {
	fmt.Println(s)
}

// [00][  0]( 0.7071 0.0000i): 0.5000
// [11][  3]( 0.7071 0.0000i): 0.5000

m0 := qsim.Measure(q0)
m1 := qsim.Measure(q1)
fmt.Println(m0.Equal(m1)) // always true

for _, s := range qsim.State() {
	fmt.Println(s)
}

// [00][  0]( 1.0000 0.0000i): 1.0000
// or
// [11][  3]( 1.0000 0.0000i): 1.0000
```

### Quantum Teleportation

```go
qsim := q.New()

// generate qubits in the state |psi>|0>|0>
psi := qsim.New(1, 2)
q0 := qsim.Zero()
q1 := qsim.Zero()

// |psi> is normalized. |psi> = a|0> + b|1>, where |a|^2 = 0.2 and |b|^2 = 0.8
for _, s := range qsim.State(psi) {
	fmt.Println(s)
}

// [0][  0]( 0.4472 0.0000i): 0.2000
// [1][  1]( 0.8944 0.0000i): 0.8000

qsim.H(q0)
qsim.CNOT(q0, q1)
qsim.CNOT(psi, q0)
qsim.H(psi)

// Alice sends mz and mx to Bob
mz := qsim.Measure(psi)
mx := qsim.Measure(q0)

// Bob applies X and Z
qsim.CondX(mx.IsOne(), q1)
qsim.CondZ(mz.IsOne(), q1)

// Bob obtains the |psi> state in q1
for _, s := range qsim.State(q1) {
	fmt.Println(s)
}

// [0][  0]( 0.4472 0.0000i): 0.2000
// [1][  1]( 0.8944 0.0000i): 0.8000
```

### Grover's Search Algorithm

```go
qsim := q.New()

// initial state
q0 := qsim.Zero()
q1 := qsim.Zero()
q2 := qsim.Zero()
q3 := qsim.Zero()

// superposition
qsim.H(q0, q1, q2, q3)

// iterations
N := number.Pow(2, qsim.NumQubits())
R := math.Floor(math.Pi / 4 * math.Sqrt(float64(N)))
for range int(R) {
	// oracle for |110>|x>
	qsim.X(q2, q3)
	qsim.H(q3)
	qsim.CCCNOT(q0, q1, q2, q3)
	qsim.H(q3)
	qsim.X(q2, q3)

	// diffuser
	qsim.H(q0, q1, q2, q3)
	qsim.X(q0, q1, q2, q3)
	qsim.H(q3)
	qsim.CCCNOT(q0, q1, q2, q3)
	qsim.H(q3)
	qsim.X(q0, q1, q2, q3)
	qsim.H(q0, q1, q2, q3)
}

for _, s := range qsim.State([]q.Qubit{q0, q1, q2}, q3) {
	fmt.Println(s)
}

// [000 0][  0   0]( 0.0508 0.0000i): 0.0026
// [000 1][  0   1]( 0.0508 0.0000i): 0.0026
// [001 0][  1   0]( 0.0508 0.0000i): 0.0026
// [001 1][  1   1]( 0.0508 0.0000i): 0.0026
// [010 0][  2   0]( 0.0508 0.0000i): 0.0026
// [010 1][  2   1]( 0.0508 0.0000i): 0.0026
// [011 0][  3   0]( 0.0508 0.0000i): 0.0026
// [011 1][  3   1]( 0.0508 0.0000i): 0.0026
// [100 0][  4   0]( 0.0508 0.0000i): 0.0026
// [100 1][  4   1]( 0.0508 0.0000i): 0.0026
// [101 0][  5   0]( 0.0508 0.0000i): 0.0026
// [101 1][  5   1]( 0.0508 0.0000i): 0.0026
// [110 0][  6   0](-0.9805 0.0000i): 0.9613 --> answer!
// [110 1][  6   1]( 0.0508 0.0000i): 0.0026
// [111 0][  7   0]( 0.0508 0.0000i): 0.0026
// [111 1][  7   1]( 0.0508 0.0000i): 0.0026
```

### Shor's Factoring Algorithm

```go
N := 15
a := 7 // co-prime with N

for i := range 10 {
	qsim := q.New()

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

	// Controlled-U
	qsim.CNOT(q2, q4)
	qsim.CNOT(q2, q5)

	// Controlled-U^2
	qsim.CNOT(q3, q5)
	qsim.CCNOT(q1, q5, q3)
	qsim.CNOT(q3, q5)

	qsim.CNOT(q6, q4)
	qsim.CCNOT(q1, q4, q6)
	qsim.CNOT(q6, q4)

	// inverse QFT
	qsim.Swap(q0, q2)
	qsim.InvQFT(q0, q1, q2)

	// measure q0, q1, q2
	m := qsim.Measure(q0, q1, q2)
	phi := number.Ldexp(m.Int(), -m.NumQubits())

	// find s/r. 0.010 -> 0.25 -> 1/4, 0.110 -> 0.75 -> 3/4, ...
	s, r, d, ok := number.FindOrder(a, N, phi)
	if !ok || number.IsOdd(r) {
		continue
	}

	// gcd(a^(r/2)-1, N), gcd(a^(r/2)+1, N)
	p0 := number.GCD(number.Pow(a, r/2)-1, N)
	p1 := number.GCD(number.Pow(a, r/2)+1, N)
	if number.IsTrivial(N, p0, p1) {
		continue
  	}

	// output result
	fmt.Printf("i=%d: N=%d, a=%d. p=%v, q=%v. s/r=%d/%d ([0.%v]~%.3f)\n", i, N, a, p0, p1, s, r, m, d)
}

// i=2: N=15, a=7. p=3, q=5. s/r=1/4 ([0.010]~0.250)
```

### Building Arbitrary 1-Qubit and Controlled Gates

```go
h := gate.U(math.Pi/2, 0, math.Pi)
x := gate.U(math.Pi, 0, math.Pi)

qsim := q.New()
q0 := qsim.Zero()
q1 := qsim.Zero()

qsim.G(h, q0)
qsim.C(x, q0, q1)

for _, s := range qsim.State() {
	fmt.Println(s)
}

// [00][  0]( 0.7071 0.0000i): 0.5000
// [11][  3]( 0.7071 0.0000i): 0.5000
```


### Density Matrix and Channels

```go
p, post := density.New(qubit.One()).
	AmplitudeDamping(0.9).
	BitFlip(0.5).
	Measure(qubit.Zero())

fmt.Printf("%.4f\n", p)
for _, r := range post.Seq2() {
	fmt.Println(r)
}

// 0.5000
// [(1+0i) (0+0i)]
// [(0+0i) (0+0i)]
```
