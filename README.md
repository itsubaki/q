# q

 - quantum computation simulator

# example

## bell state

```golang
qsim := q.New()

// generate qubits of |0>|0>
q0 := qsim.Zero()
q1 := qsim.Zero()

// apply quantum circuit
qsim.H(q0).CNOT(q0, q1)

// estimate
qsim.Estimate(q0).Probability()
// -> (0.5, 0.5)
qsim.Estimate(q1).Probability()
// -> (0.5, 0.5)

qsim.Probability()
// -> (0.5, 0, 0, 0.5)

qsim.Measure()
qsim.Probability()
// -> (1, 0, 0, 0) or (0, 0, 0, 1)

m0 := qsim.Measure(q0)
m1 := qsim.Measure(q1)
// -> m0 = |0> then m1 = |0>
// -> m0 = |1> then m1 = |1>
```

## quantum teleportation

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
qsim.Estimate(q1).Probability()
// -> (0.2, 0.8)
```

### Grover's search algorithm

```golang
qsim := q.New()

q0 := qsim.Zero()
q1 := qsim.Zero()
q2 := qsim.Zero()
q3 := qsim.One()

qsim.H(q0, q1, q2, q3)
// oracle for |011>|1>
qsim.X(q0).ControlledNot([]*q.Qubit{q0, q1, q2}, q3).X(q0)
// amp
qsim.H(q0, q1, q2, q3)
qsim.X(q0, q1, q2)
qsim.ControlledZ([]*q.Qubit{q0, q1}, q2)
qsim.H(q0, q1, q2)

qsim.Probability()
// [0 0.03125 0 0.03125 0 0.03125 0 0.78125 0 0.03125 0 0.03125 0 0.03125 0 0.03125]
```

## error correction

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

// error corretion
qsim.CNOT(q0, q3).CNOT(q1, q3)
qsim.CNOT(q1, q4).CNOT(q2, q4)

m3 := qsim.Measure(q3)
m4 := qsim.Measure(q4)

qsim.ConditionX(m3.IsOne() && m4.IsZero(), q0)
qsim.ConditionX(m3.IsOne() && m4.IsOne(), q1)
qsim.ConditionX(m3.IsZero() && m4.IsOne(), q2)

// estimate
qsim.Estimate(q0).Probability() // (0.2, 0.8)
qsim.Estimate(q1).Probability() // (0.2, 0.8)
qsim.Estimate(q2).Probability() // (0.2, 0.8)
```

# Reference

 1. Michael A. Nielsen, Issac L. Chuang, Quantum Computation and Quantum Information
 2. C. Figgatt, D. Maslov, K. A. Landsman, N. M. Linke, S. Debnath, and C. Monroe, Complete 3-Qubit Grover Search on a Programmable Quantum Computer
