# q

[![Go Report Card](https://goreportcard.com/badge/github.com/itsubaki/q?style=flat-square)](https://goreportcard.com/report/github.com/itsubaki/q)

- quantum computation simulator
- pure golang implementation
- no external library used

# Example

## Bell state

```golang
qsim := q.New()

// generate qubits of |0>|0>
q0 := qsim.Zero()
q1 := qsim.Zero()

// apply quantum circuit
qsim.H(q0).CNOT(q0, q1)

// estimate
qsim.Estimate(q0)
// -> (0.5, 0.5)
qsim.Estimate(q1)
// -> (0.5, 0.5)

qsim.Probability()
// -> (0.5, 0, 0, 0.5)

m0 := qsim.Measure(q0)
m1 := qsim.Measure(q1)
// -> m0 = |0> then m1 = |0>
// -> m0 = |1> then m1 = |1>

qsim.Probability()
// -> (1, 0, 0, 0) or (0, 0, 0, 1)
```

## Quantum teleportation

```golang
qsim := q.New()

// generate qubits of |phi>|0>|0>
// |phi> is normalized. |phi> = a|0> + b|1>, |a|^2 = 0.2, |b|^2 = 0.8
phi := qsim.New(1, 2)
q0 := qsim.Zero()
q1 := qsim.Zero()

qsim.H(q0).CNOT(q0, q1)
qsim.CNOT(phi, q0).H(phi)

// Alice send mz, mx to Bob
mz := qsim.Measure(phi)
mx := qsim.Measure(q0)

// Bob Apply Z and X
qsim.ConditionZ(mz.IsOne(), q1)
qsim.ConditionX(mx.IsOne(), q1)

// Bob got |phi> state
qsim.Estimate(q1)
// -> (0.2, 0.8)
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

// estimate
qsim.Estimate(q0) // (0.2, 0.8)
```

### Grover's search algorithm

```golang
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

for _, s := range qsim.State() {
  if s.Probability == 0 {
    continue
  }

  fmt.Println(s)
}

// 0001 0.03125
// 0011 0.03125
// 0101 0.03125
// 0111 0.78125 -> answer!
// 1001 0.03125
// 1011 0.03125
// 1101 0.03125
// 1111 0.03125
```

## Factoring 15

```golang
N := 15
a := 7 // co-prime

if number.GCD(N, a) != 1 {
  t.Errorf("%v %v\n", N, a)
}

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
  qsim.CNOT(q4, q6).CCNOT(q1, q6, q4).CNOT(q4, q6)

  // inverse QFT
  qsim.Swap(q0, q2)
  qsim.InverseQFT(q0, q1, q2)

  // measure q0, q1, q2
  m := qsim.MeasureAsBinary(q0, q1, q2)

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
```

## Density Matrix

```golang
p0, p1 := 0.1, 0.9
q0, q1 := qubit.Zero(), qubit.Zero().Apply(gate.H())
rho := density.New().Add(p0, q0).Add(p1, q1)

rho.Trace() // -> 1
rho.ExpectedValue(gate.X()) // -> 0.9
```

# Reference

1. Michael A. Nielsen, Issac L. Chuang, Quantum Computation and Quantum Information
2. C. Figgatt, D. Maslov, K. A. Landsman, N. M. Linke, S. Debnath, and C. Monroe, Complete 3-Qubit Grover Search on a Programmable Quantum Computer
3. Zhengjun Cao, Zhenfu Cao, Lihua Liu, Remarks on Quantum Modular Exponentiation and Some Experimental Demonstrations of Shor’s Algorithm
4. Michael R. Geller, Zhongyuan Zhou, Factoring 51 and 85 with 8 qubits
5. Programming Quantum Computers by Eric R. Johnson, Nic Harrigan, and Merecedes Gimeno-Segovia (O'Reilly)
