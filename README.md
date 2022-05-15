# q

[![PkgGoDev](https://pkg.go.dev/badge/github.com/itsubaki/q)](https://pkg.go.dev/github.com/itsubaki/q)
[![Go Report Card](https://goreportcard.com/badge/github.com/itsubaki/q?style=flat-square)](https://goreportcard.com/report/github.com/itsubaki/q)
[![tests](https://github.com/itsubaki/q/workflows/tests/badge.svg?branch=main)](https://github.com/itsubaki/q/actions)
[![codecov](https://codecov.io/gh/itsubaki/q/branch/main/graph/badge.svg?token=iNccCs1Tez)](https://codecov.io/gh/itsubaki/q)

- quantum computation simulator
- pure golang implementation
- full scratch

# Example

## Bell state

```golang
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
fmt.Printf("%v %v\n", m0.IsZero(), m1.IsZero())
// if m0.IsZero() is true, m1.IsZero() is also true
// if m0.IsOne()  is true, m1.IsOne()  is also true

for _, s := range qsim.State() {
  fmt.Println(s)
}
// [00][  0]( 1.0000 0.0000i): 1.0000
// or
// [11][  3]( 1.0000 0.0000i): 1.0000
```

## Quantum teleportation

```golang
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
qsim.ConditionX(mx.IsOne(), q1)
qsim.ConditionZ(mz.IsOne(), q1)

// Bob got |phi> state with q1
for _, s := range qsim.State(q1) {
  fmt.Println(s)
}
// [0][  0]( 0.4472 0.0000i): 0.2000
// [1][  1]( 0.8944 0.0000i): 0.8000
```

## Error correction

```golang
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

qsim.ConditionX(m3.IsOne() && m4.IsZero(), q0)
qsim.ConditionX(m3.IsOne() && m4.IsOne(), q1)
qsim.ConditionX(m3.IsZero() && m4.IsOne(), q2)

// decoding
qsim.CNOT(q0, q2).CNOT(q0, q1)

for _, s := range qsim.State(q0) {
  fmt.Println(s)
}
// [0][  0]( 0.4472 0.0000i): 0.2000
// [1][  1]( 0.8944 0.0000i): 0.8000
```

## Grover's search algorithm

```golang
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
// [1100][ 12](-0.9805 0.0000i): 0.9613 -> answer!
// [1101][ 13]( 0.0508 0.0000i): 0.0026
// [1110][ 14]( 0.0508 0.0000i): 0.0026
// [1111][ 15]( 0.0508 0.0000i): 0.0026
```

## Shor's factoring algorithm

```golang
N := 15
a := 7 // co-prime

for i := 0; i < 10; i++{
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

- In general, See [`cmd/shor`](./cmd/shor/main.go)

# Reference

- Michael A. Nielsen, Issac L. Chuang. Quantum Computation and Quantum Information.
