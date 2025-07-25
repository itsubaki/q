# q

[![PkgGoDev](https://pkg.go.dev/badge/github.com/itsubaki/q)](https://pkg.go.dev/github.com/itsubaki/q)
[![Go Report Card](https://goreportcard.com/badge/github.com/itsubaki/q?style=flat-square)](https://goreportcard.com/report/github.com/itsubaki/q)
[![tests](https://github.com/itsubaki/q/workflows/tests/badge.svg)](https://github.com/itsubaki/q/actions)
[![codecov](https://codecov.io/gh/itsubaki/q/branch/main/graph/badge.svg?token=iNccCs1Tez)](https://codecov.io/gh/itsubaki/q)

- Quantum computation simulator for Go
- Pure Go implementation using only the standard library

## Examples

### Bell state

```go
qsim := q.New()

// generate qubits of |0>|0>
q0 := qsim.Zero()
q1 := qsim.Zero()

// apply quantum circuit
qsim.H(q0).CNOT(q0, q1)

for _, s := range qsim.State() {
  fmt.Println(s)
}

// [00][  0]( 0.7071 0.0000i): 0.5000
// [11][  3]( 0.7071 0.0000i): 0.5000

m0 := qsim.Measure(q0)
m1 := qsim.Measure(q1)
fmt.Println(m0.Equals(m1)) // always true

for _, s := range qsim.State() {
  fmt.Println(s)
}

// [00][  0]( 1.0000 0.0000i): 1.0000
// or
// [11][  3]( 1.0000 0.0000i): 1.0000
```

### Quantum teleportation

```go
qsim := q.New()

// generate qubits of |phi>|0>|0>
phi := qsim.New(1, 2)
q0 := qsim.Zero()
q1 := qsim.Zero()

// |phi> is normalized. |phi> = a|0> + b|1>, |a|^2 = 0.2, |b|^2 = 0.8
for _, s := range qsim.State(phi) {
  fmt.Println(s)
}

// [0][  0]( 0.4472 0.0000i): 0.2000
// [1][  1]( 0.8944 0.0000i): 0.8000

qsim.H(q0).CNOT(q0, q1)
qsim.CNOT(phi, q0).H(phi)

// Alice send mz, mx to Bob
mz := qsim.Measure(phi)
mx := qsim.Measure(q0)

// Bob Apply X and Z
qsim.CondX(mx.IsOne(), q1)
qsim.CondZ(mz.IsOne(), q1)

// Bob got |phi> state with q1
for _, s := range qsim.State(q1) {
  fmt.Println(s)
}

// [0][  0]( 0.4472 0.0000i): 0.2000
// [1][  1]( 0.8944 0.0000i): 0.8000
```

### Error correction

```go
qsim := q.New()

q0 := qsim.New(1, 2) // (0.2, 0.8)

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

qsim.CondX(m3.IsOne() && m4.IsZero(), q0)
qsim.CondX(m3.IsOne() && m4.IsOne(), q1)
qsim.CondX(m3.IsZero() && m4.IsOne(), q2)

// decoding
qsim.CNOT(q0, q2).CNOT(q0, q1)

for _, s := range qsim.State(q0) {
  fmt.Println(s)
}

// [0][  0]( 0.4472 0.0000i): 0.2000
// [1][  1]( 0.8944 0.0000i): 0.8000
```

### Grover's search algorithm

```go
qsim := q.New()

// initial state
q0 := qsim.Zero()
q1 := qsim.Zero()
q2 := qsim.Zero()
q3 := qsim.Zero()

// superposition
qsim.H(q0, q1, q2, q3)

// iteration
N := number.Pow(2, qsim.NumQubits())
r := math.Floor(math.Pi / 4 * math.Sqrt(float64(N)))
for range int(r) {
  // oracle for |110>|x>
  qsim.X(q2, q3)
  qsim.H(q3).CCCNOT(q0, q1, q2, q3).H(q3)
  qsim.X(q2, q3)

  // amplification
  qsim.H(q0, q1, q2, q3)
  qsim.X(q0, q1, q2, q3)
  qsim.H(q3).CCCNOT(q0, q1, q2, q3).H(q3)
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

### Shor's factoring algorithm

```go
N := 15
a := 7 // co-prime

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
  qsim.CNOT(q3, q5).CCNOT(q1, q5, q3).CNOT(q3, q5)
  qsim.CNOT(q6, q4).CCNOT(q1, q4, q6).CNOT(q6, q4)

  // inverse QFT
  qsim.Swap(q0, q2)
  qsim.InvQFT(q0, q1, q2)

  // measure q0, q1, q2
  m := qsim.Measure(q0, q1, q2).BinaryString()

  // find s/r. 0.010 -> 0.25 -> 1/4, 0.110 -> 0.75 -> 3/4, ...
  s, r, d, ok := number.FindOrder(a, N, fmt.Sprintf("0.%s", m))
  if !ok || number.IsOdd(r) {
    continue
  }

  // gcd(a^(r/2)-1, N), gcd(a^(r/2)+1, N)
  p0 := number.GCD(number.Pow(a, r/2)-1, N)
  p1 := number.GCD(number.Pow(a, r/2)+1, N)
  if number.IsTrivial(N, p0, p1) {
    continue
  }

  // result
  fmt.Printf("i=%d: N=%d, a=%d. p=%v, q=%v. s/r=%d/%d ([0.%v]~%.3f)\n", i, N, a, p0, p1, s, r, m, d)
}

// i=2: N=15, a=7. p=3, q=5. s/r=1/4 ([0.010]~0.250)
```

### Any 1-qubit quantum gate and its controlled gate

```go
h := gate.U(math.Pi/2, 0, math.Pi)
x := gate.U(math.Pi, 0, math.Pi)

qsim := q.New()
q0 := qsim.Zero()
q1 := qsim.Zero()

qsim.Apply(h, q0)
qsim.C(x, q0, q1)

for _, s := range qsim.State() {
  fmt.Println(s)
}

// [00][  0]( 0.7071 0.0000i): 0.5000
// [11][  3]( 0.7071 0.0000i): 0.5000
```

## References

[1] M. A. Nielsen and I. L. Chuang, *Quantum Computation and Quantum Information*, 10th ed. Cambridge, U.K.: Cambridge University Press, 2010.
