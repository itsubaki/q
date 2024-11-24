package q

import (
	"github.com/itsubaki/q/math/matrix"
	"github.com/itsubaki/q/math/number"
	"github.com/itsubaki/q/math/rand"
	"github.com/itsubaki/q/quantum/gate"
	"github.com/itsubaki/q/quantum/qubit"
)

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

// Theta returns 2 * pi / 2**k
func Theta(k int) float64 {
	return gate.Theta(k)
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
		q.qb = qubit.New(v...)
		q.qb.Rand = q.Rand
		return Qubit(0)
	}

	q.qb.TensorProduct(qubit.New(v...))
	return Qubit(q.NumberOfBit() - 1)
}

// NewOf returns a new qubit from binary string.
func (q *Q) NewOf(binary string) []Qubit {
	qb := make([]Qubit, len(binary))
	for i, b := range binary {
		if b == '0' {
			qb[i] = q.Zero()
			continue
		}

		qb[i] = q.One()
	}

	return qb
}

// Zero returns a qubit in the zero state.
func (q *Q) Zero() Qubit {
	return q.New(1, 0)
}

// One returns a qubit in the one state.
func (q *Q) One() Qubit {
	return q.New(0, 1)
}

// ZeroWith returns n qubits in the zero state.
func (q *Q) ZeroWith(n int) []Qubit {
	qb := make([]Qubit, n)
	for i := 0; i < n; i++ {
		qb[i] = q.Zero()
	}

	return qb
}

// One returns n qubits in the one state.
func (q *Q) OneWith(n int) []Qubit {
	qb := make([]Qubit, n)
	for i := 0; i < n; i++ {
		qb[i] = q.One()
	}

	return qb
}

// ZeroLog2 returns n qubits in the zero state.
// n is greater than or equal to log2(N).
func (q *Q) ZeroLog2(N int) []Qubit {
	return q.ZeroWith(number.Log2(N) + 1)
}

// NumberOfBit returns the number of qubits.
func (q *Q) NumberOfBit() int {
	return q.qb.NumberOfBit()
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

// U applies U gate.
func (q *Q) U(theta, phi, lambda float64, qb ...Qubit) *Q {
	return q.Apply(gate.U(theta, phi, lambda), qb...)
}

// I applies I gate.
func (q *Q) I(qb ...Qubit) *Q {
	return q.Apply(gate.I(), qb...)
}

// X applies X gate.
func (q *Q) X(qb ...Qubit) *Q {
	return q.Apply(gate.X(), qb...)
}

// Y applies Y gate.
func (q *Q) Y(qb ...Qubit) *Q {
	return q.Apply(gate.Y(), qb...)
}

// Z applies Z gate.
func (q *Q) Z(qb ...Qubit) *Q {
	return q.Apply(gate.Z(), qb...)
}

// H applies H gate.
func (q *Q) H(qb ...Qubit) *Q {
	return q.Apply(gate.H(), qb...)
}

// S applies S gate.
func (q *Q) S(qb ...Qubit) *Q {
	return q.Apply(gate.S(), qb...)
}

// T applies T gate.
func (q *Q) T(qb ...Qubit) *Q {
	return q.Apply(gate.T(), qb...)
}

// R applies R gate with theta.
func (q *Q) R(theta float64, qb ...Qubit) *Q {
	return q.Apply(gate.R(theta), qb...)
}

// RX applies RX gate with theta.
func (q *Q) RX(theta float64, qb ...Qubit) *Q {
	return q.Apply(gate.RX(theta), qb...)
}

// RY applies RY gate with theta.
func (q *Q) RY(theta float64, qb ...Qubit) *Q {
	return q.Apply(gate.RY(theta), qb...)
}

// RZ applies RZ gate with theta.
func (q *Q) RZ(theta float64, qb ...Qubit) *Q {
	return q.Apply(gate.RZ(theta), qb...)
}

// Apply applies matrix to qubits.
func (q *Q) Apply(m matrix.Matrix, qb ...Qubit) *Q {
	if len(qb) < 1 {
		q.qb.Apply(m)
		return q
	}

	idx := make(map[int]bool)
	for _, i := range Index(qb...) {
		idx[i] = true
	}

	g := gate.I()
	if _, ok := idx[0]; ok {
		g = m
	}

	for i := 1; i < q.NumberOfBit(); i++ {
		if _, ok := idx[i]; ok {
			g = g.TensorProduct(m)
			continue
		}

		g = g.TensorProduct(gate.I())
	}

	q.qb.Apply(g)
	return q
}

func (q *Q) Controlled(m matrix.Matrix, control []Qubit, target Qubit) *Q {
	n := q.NumberOfBit()
	g := gate.Controlled(m, n, Index(control...), target.Index())
	q.qb.Apply(g)
	return q
}

func (q *Q) C(m matrix.Matrix, control, target Qubit) *Q {
	return q.Controlled(m, []Qubit{control}, target)
}

// ControlledNot applies CNOT gate.
func (q *Q) ControlledNot(control []Qubit, target Qubit) *Q {
	n := q.NumberOfBit()
	g := gate.ControlledNot(n, Index(control...), target.Index())
	q.qb.Apply(g)
	return q
}

// CNOT applies CNOT gate.
func (q *Q) CNOT(control, target Qubit) *Q {
	return q.ControlledNot([]Qubit{control}, target)
}

// CCNOT applies CCNOT gate.
func (q *Q) CCNOT(control0, control1, target Qubit) *Q {
	return q.ControlledNot([]Qubit{control0, control1}, target)
}

// CCCNOT applies CCCNOT gate.
func (q *Q) CCCNOT(control0, control1, control2, target Qubit) *Q {
	return q.ControlledNot([]Qubit{control0, control1, control2}, target)
}

// Toffoli applies Toffoli gate.
func (q *Q) Toffoli(control0, control1, target Qubit) *Q {
	return q.CCNOT(control0, control1, target)
}

// ControlledZ applies Controlled-Z gate.
func (q *Q) ControlledZ(control []Qubit, target Qubit) *Q {
	n := q.NumberOfBit()
	g := gate.ControlledZ(n, Index(control...), target.Index())
	q.qb.Apply(g)
	return q
}

func (q *Q) CZ(control, target Qubit) *Q {
	return q.ControlledZ([]Qubit{control}, target)
}

func (q *Q) CCZ(control0, control1, target Qubit) *Q {
	return q.ControlledZ([]Qubit{control0, control1}, target)
}

func (q *Q) ControlledR(theta float64, control []Qubit, target Qubit) *Q {
	n := q.NumberOfBit()
	g := gate.ControlledR(theta, n, Index(control...), target.Index())
	q.qb.Apply(g)
	return q
}

// CR applies Controlled-R gate.
func (q *Q) CR(theta float64, control, target Qubit) *Q {
	return q.ControlledR(theta, []Qubit{control}, target)
}

// ControlledModExp2 applies Controlled-ModExp2 gate.
func (q *Q) ControlledModExp2(a, j, N int, control Qubit, target []Qubit) *Q {
	n := q.NumberOfBit()
	g := gate.ControlledModExp2(n, a, j, N, control.Index(), Index(target...))
	q.qb.Apply(g)
	return q
}

// CModExp2 applies Controlled-ModExp2 gate.
func (q *Q) CModExp2(a, N int, control []Qubit, target []Qubit) *Q {
	for i := 0; i < len(control); i++ {
		q.ControlledModExp2(a, i, N, control[i], target)
	}

	return q
}

// Cond applies m if condition is true.
func (q *Q) Cond(condition bool, m matrix.Matrix, qb ...Qubit) *Q {
	if condition {
		return q.Apply(m, qb...)
	}

	return q
}

// CondX applies X gate if condition is true.
func (q *Q) CondX(condition bool, qb ...Qubit) *Q {
	return q.Cond(condition, gate.X(), qb...)
}

// CondZ applies Z gate if condition is true.
func (q *Q) CondZ(condition bool, qb ...Qubit) *Q {
	return q.Cond(condition, gate.Z(), qb...)
}

// Swap applies Swap gate.
func (q *Q) Swap(qb ...Qubit) *Q {
	n := q.NumberOfBit()
	l := len(qb)

	for i := 0; i < l/2; i++ {
		q0, q1 := qb[i], qb[(l-1)-i]
		g := gate.Swap(n, q0.Index(), q1.Index())
		q.qb.Apply(g)
	}

	return q
}

// QFT applies Quantum Fourier Transform.
func (q *Q) QFT(qb ...Qubit) *Q {
	l := len(qb)
	for i := 0; i < l; i++ {
		q.H(qb[i])

		k := 2
		for j := i + 1; j < l; j++ {
			q.CR(Theta(k), qb[j], qb[i])
			k++
		}
	}

	return q
}

// InverseQFT applies Inverse Quantum Fourier Transform.
func (q *Q) InverseQFT(qb ...Qubit) *Q {
	l := len(qb)
	for i := l - 1; i > -1; i-- {
		k := l - i
		for j := l - 1; j > i; j-- {
			q.CR(-1*Theta(k), qb[j], qb[i])
			k--
		}

		q.H(qb[i])
	}

	return q
}

// InvQFT applies Inverse Quantum Fourier Transform.
func (q *Q) InvQFT(qb ...Qubit) *Q {
	return q.InverseQFT(qb...)
}

// IQFT applies Inverse Quantum Fourier Transform.
func (q *Q) IQFT(qb ...Qubit) *Q {
	return q.InverseQFT(qb...)
}

// M returns the measured state of qubits.
func (q *Q) M(qb ...Qubit) *qubit.Qubit {
	return q.Measure(qb...)
}

// Measure returns the measured state of qubits.
func (q *Q) Measure(qb ...Qubit) *qubit.Qubit {
	if len(qb) < 1 {
		n := q.NumberOfBit()
		m := make([]*qubit.Qubit, n)
		for i := 0; i < n; i++ {
			m[i] = q.qb.Measure(i)
		}

		return qubit.TensorProduct(m...)
	}

	m := make([]*qubit.Qubit, len(qb))
	for i := 0; i < len(qb); i++ {
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
		Rand: q.qb.Rand,
	}
}

// Raw returns the internal qubit.
func (q *Q) Raw() *qubit.Qubit {
	return q.qb
}

// String returns the string representation of a quantum computation simulator.
func (q *Q) String() string {
	return q.qb.String()
}

// State returns the state of qubits.
func (q *Q) State(reg ...any) []qubit.State {
	idx := make([][]int, 0)
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
