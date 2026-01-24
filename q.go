package q

import (
	"sort"

	"github.com/itsubaki/q/math/matrix"
	"github.com/itsubaki/q/math/number"
	"github.com/itsubaki/q/math/rand"
	"github.com/itsubaki/q/math/vector"
	"github.com/itsubaki/q/quantum/gate"
	"github.com/itsubaki/q/quantum/qubit"
)

// Theta returns 2 * pi / 2**k
func Theta(k int) float64 {
	return gate.Theta(k)
}

// Qubit is a quantum bit.
type Qubit int

// Index returns the index of qubit.
func (q Qubit) Index() int {
	return int(q)
}

// Index returns the index list of qubits.
func Index(qb ...Qubit) []int {
	idx := make([]int, len(qb))
	for i := range qb {
		idx[i] = qb[i].Index()
	}

	return idx
}

// Q is a quantum computation simulator.
type Q struct {
	qb   *qubit.Qubit
	Rand func() float64
}

// New returns a new quantum computation simulator.
func New() *Q {
	return &Q{
		qb:   nil,
		Rand: rand.Float64,
	}
}

// New returns a new qubit.
func (q *Q) New(v ...complex128) Qubit {
	if q.qb == nil {
		q.qb = qubit.New(vector.New(v...))
		q.qb.Rand = q.Rand
		return Qubit(0)
	}

	q.qb.TensorProduct(qubit.New(vector.New(v...)))
	return Qubit(q.NumQubits() - 1)
}

// Zero returns a qubit in the zero state.
func (q *Q) Zero() Qubit {
	return q.New(1, 0)
}

// One returns a qubit in the one state.
func (q *Q) One() Qubit {
	return q.New(0, 1)
}

// Zeros returns n qubits in the zero state.
func (q *Q) Zeros(n int) []Qubit {
	qb := make([]Qubit, n)
	for i := range n {
		qb[i] = q.Zero()
	}

	return qb
}

// Ones returns n qubits in the one state.
func (q *Q) Ones(n int) []Qubit {
	qb := make([]Qubit, n)
	for i := range n {
		qb[i] = q.One()
	}

	return qb
}

// ZeroLog2 returns log2(N)+1 qubits in the zero state.
func (q *Q) ZeroLog2(N int) []Qubit {
	return q.Zeros(number.Log2(N) + 1)
}

// NumQubits returns the number of qubits.
func (q *Q) NumQubits() int {
	if q.qb == nil {
		return 0
	}

	return q.qb.NumQubits()
}

// Amplitude returns the amplitude of qubits.
func (q *Q) Amplitude() []complex128 {
	return q.qb.Amplitude()
}

// Probability returns the probability of qubits.
func (q *Q) Probability() []float64 {
	return q.qb.Probability()
}

// Reset sets qubits to the zero state.
func (q *Q) Reset(qb ...Qubit) {
	for i := range qb {
		if q.Measure(qb[i]).IsOne() {
			q.X(qb[i])
		}
	}
}

// Apply applies a list of gates to the qubit.
func (q *Q) Apply(g ...*matrix.Matrix) *Q {
	q.qb.Apply(g...)
	return q
}

// G applies a gate.
func (q *Q) G(g *matrix.Matrix, qb ...Qubit) *Q {
	for i := range qb {
		q.qb.G(g, qb[i].Index())
	}

	return q
}

// U applies U gate.
func (q *Q) U(theta, phi, lambda float64, qb ...Qubit) *Q {
	for i := range qb {
		q.qb.U(theta, phi, lambda, qb[i].Index())
	}

	return q
}

// I applies I gate.
func (q *Q) I(qb ...Qubit) *Q {
	for i := range qb {
		q.qb.I(qb[i].Index())
	}

	return q
}

// X applies X gate.
func (q *Q) X(qb ...Qubit) *Q {
	for i := range qb {
		q.qb.X(qb[i].Index())
	}

	return q
}

// Y applies Y gate.
func (q *Q) Y(qb ...Qubit) *Q {
	for i := range qb {
		q.qb.Y(qb[i].Index())
	}

	return q
}

// Z applies Z gate.
func (q *Q) Z(qb ...Qubit) *Q {
	for i := range qb {
		q.qb.Z(qb[i].Index())
	}

	return q
}

// H applies H gate.
func (q *Q) H(qb ...Qubit) *Q {
	for i := range qb {
		q.qb.H(qb[i].Index())
	}

	return q
}

// S applies S gate.
func (q *Q) S(qb ...Qubit) *Q {
	for i := range qb {
		q.qb.S(qb[i].Index())
	}

	return q
}

// T applies T gate.
func (q *Q) T(qb ...Qubit) *Q {
	for i := range qb {
		q.qb.T(qb[i].Index())
	}

	return q
}

// R applies R gate with theta.
func (q *Q) R(theta float64, qb ...Qubit) *Q {
	for i := range qb {
		q.qb.R(theta, qb[i].Index())
	}

	return q
}

// RX applies RX gate with theta.
func (q *Q) RX(theta float64, qb ...Qubit) *Q {
	for i := range qb {
		q.qb.RX(theta, qb[i].Index())
	}

	return q
}

// RY applies RY gate with theta.
func (q *Q) RY(theta float64, qb ...Qubit) *Q {
	for i := range qb {
		q.qb.RY(theta, qb[i].Index())
	}

	return q
}

// RZ applies RZ gate with theta.
func (q *Q) RZ(theta float64, qb ...Qubit) *Q {
	for i := range qb {
		q.qb.RZ(theta, qb[i].Index())
	}

	return q
}

// C applies controlled operation with a gate.
func (q *Q) C(g *matrix.Matrix, control, target Qubit) *Q {
	return q.Controlled(g, []Qubit{control}, []Qubit{target})
}

// CU applies controlled unitary operation.
func (q *Q) CU(theta, phi, lambda float64, control, target Qubit) *Q {
	return q.ControlledU(theta, phi, lambda, []Qubit{control}, []Qubit{target})
}

// CX applies CNOT gate.
func (q *Q) CX(control, target Qubit) *Q {
	return q.ControlledNot([]Qubit{control}, []Qubit{target})
}

// CNOT applies CNOT gate.
func (q *Q) CNOT(control, target Qubit) *Q {
	return q.ControlledNot([]Qubit{control}, []Qubit{target})
}

// CCNOT applies CCNOT gate.
func (q *Q) CCNOT(control0, control1, target Qubit) *Q {
	return q.ControlledNot([]Qubit{control0, control1}, []Qubit{target})
}

// CCCNOT applies CCCNOT gate.
func (q *Q) CCCNOT(control0, control1, control2, target Qubit) *Q {
	return q.ControlledNot([]Qubit{control0, control1, control2}, []Qubit{target})
}

// CZ applies Controlled-Z gate.
func (q *Q) CZ(control, target Qubit) *Q {
	return q.ControlledZ([]Qubit{control}, []Qubit{target})
}

// CCZ applies CCZ gate with two controls.
func (q *Q) CCZ(control0, control1, target Qubit) *Q {
	return q.ControlledZ([]Qubit{control0, control1}, []Qubit{target})
}

// CR applies Controlled-R gate.
func (q *Q) CR(theta float64, control, target Qubit) *Q {
	return q.ControlledR(theta, []Qubit{control}, []Qubit{target})
}

// Controlled applies controlled operation with a gate.
func (q *Q) Controlled(g *matrix.Matrix, control, target []Qubit) *Q {
	for i := range target {
		q.qb.Controlled(g, Index(control...), target[i].Index())
	}

	return q
}

// ControlledU applies controlled unitary operation.
func (q *Q) ControlledU(theta, phi, lambda float64, control, target []Qubit) *Q {
	for i := range target {
		q.qb.ControlledU(theta, phi, lambda, Index(control...), target[i].Index())
	}

	return q
}

// ControlledH applies controlled-Hadamard gate.
func (q *Q) ControlledH(control, target []Qubit) *Q {
	for i := range target {
		q.qb.ControlledH(Index(control...), target[i].Index())
	}

	return q
}

// ControlledX applies CNOT gate.
func (q *Q) ControlledX(control, target []Qubit) *Q {
	return q.ControlledNot(control, target)
}

// ControlledNot applies CNOT gate.
func (q *Q) ControlledNot(control, target []Qubit) *Q {
	for i := range target {
		q.qb.ControlledX(Index(control...), target[i].Index())
	}

	return q
}

// ControlledZ applies Controlled-Z gate.
func (q *Q) ControlledZ(control, target []Qubit) *Q {
	for i := range target {
		q.qb.ControlledZ(Index(control...), target[i].Index())
	}

	return q
}

// ControlledR applies Controlled-R gate.
func (q *Q) ControlledR(theta float64, control, target []Qubit) *Q {
	for i := range target {
		q.qb.ControlledR(theta, Index(control...), target[i].Index())
	}

	return q
}

// CondX applies X gate if condition is true.
func (q *Q) CondX(condition bool, qb ...Qubit) *Q {
	if condition {
		return q.X(qb...)
	}

	return q
}

// CondZ applies Z gate if condition is true.
func (q *Q) CondZ(condition bool, qb ...Qubit) *Q {
	if condition {
		return q.Z(qb...)
	}

	return q
}

// Cond applies gate if condition is true.
func (q *Q) Cond(condition bool, g *matrix.Matrix, qb ...Qubit) *Q {
	if condition {
		return q.G(g, qb...)
	}

	return q
}

// Swap applies Swap gate.
func (q *Q) Swap(qb ...Qubit) *Q {
	l := len(qb)
	for i := range l / 2 {
		q0, q1 := qb[i], qb[(l-1)-i]
		q.qb.Swap(q0.Index(), q1.Index())
	}

	return q
}

// QFT applies Quantum Fourier Transform.
func (q *Q) QFT(qb ...Qubit) *Q {
	q.qb.QFT(Index(qb...)...)
	return q
}

// InvQFT applies Inverse Quantum Fourier Transform.
func (q *Q) InvQFT(qb ...Qubit) *Q {
	q.qb.InvQFT(Index(qb...)...)
	return q
}

// M returns the measured state of qubits.
func (q *Q) M(qb ...Qubit) *qubit.Qubit {
	return q.Measure(qb...)
}

// Measure returns the measured state of qubits.
func (q *Q) Measure(qb ...Qubit) *qubit.Qubit {
	if len(qb) < 1 {
		qb = make([]Qubit, q.NumQubits())
		for i := range qb {
			qb[i] = Qubit(i)
		}
	}

	m := make([]*qubit.Qubit, len(qb))
	for i := range qb {
		m[i] = q.qb.Measure(qb[i].Index())
	}

	return qubit.TensorProduct(m...)
}

// Clone returns a clone of a quantum computation simulator.
func (q *Q) Clone() *Q {
	if q.qb == nil {
		return &Q{
			qb:   nil,
			Rand: q.Rand,
		}
	}

	return &Q{
		qb:   q.qb.Clone(),
		Rand: q.Rand,
	}
}

// Underlying returns the internal qubit.
func (q *Q) Underlying() *qubit.Qubit {
	return q.qb
}

// String returns the string representation of a quantum computation simulator.
func (q *Q) String() string {
	return q.qb.String()
}

// State returns the state of qubits.
func (q *Q) State(reg ...any) []qubit.State {
	var idx [][]int
	for _, r := range reg {
		switch r := r.(type) {
		case Qubit:
			idx = append(idx, []int{r.Index()})
		case []Qubit:
			idx = append(idx, Index(r...))
		}
	}

	return q.qb.State(idx...)
}

// Top returns the top n states by probability.
// if n < 0, returns input states.
func Top(s []qubit.State, n int) []qubit.State {
	if n < 0 {
		return s
	}

	sort.Slice(s, func(i, j int) bool {
		return s[i].Probability() > s[j].Probability()
	})

	return s[:min(n, len(s))]
}
