# q

[![Build Status](https://travis-ci.org/axamon/q.svg?branch=master)](https://travis-ci.org/axamon/q)
[![GoDoc](https://godoc.org/github.com/axamon/q?status.svg)](https://godoc.org/github.com/axamon/q)
[![codecov](https://codecov.io/gh/axamon/q/branch/master/graph/badge.svg)](https://codecov.io/gh/axamon/q)

 - quantum computation simulator
 - pure golang implementation
 - no external library used

# example

in the example folder

 - bell state
 - quantum teleportation
 - fatoring 15


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

## Factoring 15



# Reference

 1. Michael A. Nielsen, Issac L. Chuang, Quantum Computation and Quantum Information
 2. C. Figgatt, D. Maslov, K. A. Landsman, N. M. Linke, S. Debnath, and C. Monroe, Complete 3-Qubit Grover Search on a Programmable Quantum Computer
