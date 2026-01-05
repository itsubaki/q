package qubit

import (
	"math"
	"math/cmplx"
	"strconv"
	"strings"

	"github.com/itsubaki/q/math/epsilon"
	"github.com/itsubaki/q/math/matrix"
	"github.com/itsubaki/q/math/number"
	"github.com/itsubaki/q/math/rand"
	"github.com/itsubaki/q/math/vector"
	"github.com/itsubaki/q/quantum/gate"
)

// Qubit is a qubit.
type Qubit struct {
	n     int
	state *vector.Vector
	Rand  func() float64 // Random number generator
}

// New returns a new qubit.
func New(v *vector.Vector) *Qubit {
	q := &Qubit{
		n:     number.Log2(len(v.Data)),
		state: v,
		Rand:  rand.Float64,
	}

	q.Normalize()
	return q
}

// NewWithIndex returns a new qubit with the given index.
func NewWithIndex(n, idx int) *Qubit {
	s := make([]complex128, 1<<n)
	s[idx] = 1

	return New(vector.New(s...))
}

// Zero returns a qubit in the zero state.
func Zero(n ...int) *Qubit {
	if len(n) == 0 {
		n = []int{1}
	}

	return NewWithIndex(n[0], 0)
}

// One returns a qubit in the one state.
func One(n ...int) *Qubit {
	if len(n) == 0 {
		n = []int{1}
	}

	return NewWithIndex(n[0], (1<<n[0])-1)
}

// Plus returns a qubit in the plus state.
// The plus state is defined as (|0> + |1>) / sqrt(2).
func Plus(n ...int) *Qubit {
	qb := Zero(n...)
	for i := range qb.NumQubits() {
		qb.H(i)
	}

	return qb
}

// Minus returns a qubit in the minus state.
// The minus state is defined as (|0> - |1>) / sqrt(2).
func Minus(n ...int) *Qubit {
	qb := One(n...)
	for i := range qb.NumQubits() {
		qb.H(i)
	}

	return qb
}

// From returns a new qubit from a binary string.
func From(binary string) *Qubit {
	list := make([]*Qubit, len(binary))
	for i, c := range binary {
		switch c {
		case '0':
			list[i] = Zero()
		case '1':
			list[i] = One()
		case '+':
			list[i] = Plus()
		case '-':
			list[i] = Minus()
		}
	}

	return TensorProduct(list...)
}

// NumQubits returns the number of qubits.
func (q *Qubit) NumQubits() int {
	return q.n
}

// Dim returns the dimension of q.
func (q *Qubit) Dim() int {
	return len(q.state.Data)
}

// IsZero returns true if q is zero qubit.
func (q *Qubit) IsZero(tol ...float64) bool {
	return q.Equal(Zero(), tol...)
}

// IsOne returns true if q is one qubit.
func (q *Qubit) IsOne(tol ...float64) bool {
	return q.Equal(One(), tol...)
}

// Amplitude returns the amplitude of q.
func (q *Qubit) Amplitude() []complex128 {
	return q.state.Data
}

// Probability returns the probability of q.
func (q *Qubit) Probability() []float64 {
	p := make([]float64, len(q.Amplitude()))
	for i, a := range q.Amplitude() {
		p[i] = math.Pow(cmplx.Abs(a), 2)
	}

	return p
}

// InnerProduct returns the inner product of q and qb.
func (q *Qubit) InnerProduct(qb *Qubit) complex128 {
	return q.state.InnerProduct(qb.state)
}

// OuterProduct returns the outer product of q and qb.
func (q *Qubit) OuterProduct(qb *Qubit) *matrix.Matrix {
	return q.state.OuterProduct(qb.state)
}

// TensorProduct returns the tensor product of q and qb.
func (q *Qubit) TensorProduct(qb *Qubit) *Qubit {
	q.n = q.n + qb.n
	q.state = q.state.TensorProduct(qb.state)
	return q
}

// Apply applies a list of gates to the qubit.
func (q *Qubit) Apply(g ...*matrix.Matrix) *Qubit {
	for _, v := range g {
		q.state = q.state.Apply(v)
	}

	return q
}

// G applies a gate to the qubit at the given index.
func (q *Qubit) G(g *matrix.Matrix, idx int) {
	stride := 1 << (q.NumQubits() - 1 - idx)
	for i := 0; i < q.Dim(); i += 2 * stride {
		for j := range stride {
			a, b := q.state.Data[i+j], q.state.Data[i+j+stride]
			q.state.Data[i+j] = g.At(0, 0)*a + g.At(0, 1)*b
			q.state.Data[i+j+stride] = g.At(1, 0)*a + g.At(1, 1)*b
		}
	}
}

// U applies a unitary gate.
func (q *Qubit) U(theta, phi, lambda float64, idx int) *Qubit {
	sin := cmplx.Sin(complex(theta/2, 0))
	cos := cmplx.Cos(complex(theta/2, 0))

	e0 := cmplx.Exp(complex(0, phi))
	e1 := cmplx.Exp(complex(0, lambda))
	e2 := cmplx.Exp(complex(0, phi+lambda))

	stride := 1 << (q.NumQubits() - 1 - idx)
	for i := 0; i < q.Dim(); i += 2 * stride {
		for j := range stride {
			a, b := q.state.Data[i+j], q.state.Data[i+j+stride]
			q.state.Data[i+j] = cos*a - e1*sin*b
			q.state.Data[i+j+stride] = e0*sin*a + e2*cos*b
		}
	}

	return q
}

// I applies I gate.
func (q *Qubit) I(idx int) *Qubit {
	return q
}

// H applies H gate.
func (q *Qubit) H(idx int) *Qubit {
	stride := 1 << (q.NumQubits() - 1 - idx)
	for i := 0; i < q.Dim(); i += 2 * stride {
		for j := range stride {
			a, b := q.state.Data[i+j], q.state.Data[i+j+stride]
			q.state.Data[i+j] = (a + b) / complex(math.Sqrt2, 0)
			q.state.Data[i+j+stride] = (a - b) / complex(math.Sqrt2, 0)
		}
	}

	return q
}

// X applies X gate.
func (q *Qubit) X(idx int) *Qubit {
	stride := 1 << (q.NumQubits() - 1 - idx)
	for i := 0; i < q.Dim(); i += 2 * stride {
		for j := range stride {
			q.state.Data[i+j], q.state.Data[i+j+stride] = q.state.Data[i+j+stride], q.state.Data[i+j]
		}
	}

	return q
}

// Y applies Y gate.
func (q *Qubit) Y(idx int) *Qubit {
	stride := 1 << (q.NumQubits() - 1 - idx)
	for i := 0; i < q.Dim(); i += 2 * stride {
		for j := range stride {
			a, b := q.state.Data[i+j], q.state.Data[i+j+stride]
			q.state.Data[i+j] = b * complex(0, -1)
			q.state.Data[i+j+stride] = a * complex(0, 1)
		}
	}

	return q
}

// Z applies Z gate.
func (q *Qubit) Z(idx int) *Qubit {
	stride := 1 << (q.NumQubits() - 1 - idx)
	for i := 0; i < q.Dim(); i += 2 * stride {
		for j := range stride {
			q.state.Data[i+j+stride] *= -1
		}
	}

	return q
}

// R applies a phase rotation of theta.
func (q *Qubit) R(theta float64, idx int) *Qubit {
	mask := 1 << (q.NumQubits() - 1 - idx)

	phase := cmplx.Exp(complex(0, theta))
	for i := range q.Dim() {
		if (i & mask) == 0 {
			continue
		}

		q.state.Data[i] *= phase
	}

	return q
}

// S applies the S gate.
func (q *Qubit) S(idx int) *Qubit {
	return q.R(math.Pi/2, idx)
}

// T applies the T gate.
func (q *Qubit) T(idx int) *Qubit {
	return q.R(math.Pi/4, idx)
}

// RX applies the rotation around X-axis.
func (q *Qubit) RX(theta float64, idx int) *Qubit {
	sin := cmplx.Sin(complex(theta/2, 0))
	cos := cmplx.Cos(complex(theta/2, 0))

	stride := 1 << (q.NumQubits() - 1 - idx)
	for i := 0; i < q.Dim(); i += 2 * stride {
		for j := range stride {
			a, b := q.state.Data[i+j], q.state.Data[i+j+stride]
			q.state.Data[i+j] = cos*a - 1i*sin*b
			q.state.Data[i+j+stride] = -1i*sin*a + cos*b
		}
	}

	return q
}

// RY applies the rotation around Y-axis.
func (q *Qubit) RY(theta float64, idx int) *Qubit {
	sin := cmplx.Sin(complex(theta/2, 0))
	cos := cmplx.Cos(complex(theta/2, 0))

	stride := 1 << (q.NumQubits() - 1 - idx)
	for i := 0; i < q.Dim(); i += 2 * stride {
		for j := range stride {
			a, b := q.state.Data[i+j], q.state.Data[i+j+stride]
			q.state.Data[i+j] = cos*a - 1*sin*b
			q.state.Data[i+j+stride] = sin*a + cos*b
		}
	}

	return q
}

// RZ applies the rotation around Z-axis.
func (q *Qubit) RZ(theta float64, idx int) *Qubit {
	e0 := cmplx.Exp(complex(0, -theta/2))
	e1 := cmplx.Exp(complex(0, theta/2))

	stride := 1 << (q.NumQubits() - 1 - idx)
	for i := 0; i < q.Dim(); i += 2 * stride {
		for j := range stride {
			q.state.Data[i+j] *= e0
			q.state.Data[i+j+stride] *= e1
		}
	}

	return q
}

// C applies a controlled gate.
func (q *Qubit) C(g *matrix.Matrix, control, target int) *Qubit {
	return q.Controlled(g, []int{control}, target)
}

// CU applies a controlled unitary operation.
func (q *Qubit) CU(theta, phi, lambda float64, control, target int) *Qubit {
	return q.ControlledU(theta, phi, lambda, []int{control}, target)
}

// CH applies the controlled-H gate.
func (q *Qubit) CH(control, target int) *Qubit {
	return q.ControlledH([]int{control}, target)
}

// CX applies the controlled-X gate.
func (q *Qubit) CX(control, target int) *Qubit {
	return q.ControlledX([]int{control}, target)
}

// CZ applies the controlled-Z gate.
func (q *Qubit) CZ(control, target int) *Qubit {
	return q.ControlledZ([]int{control}, target)
}

// CR applies a controlled rotation around the Z-axis.
func (q *Qubit) CR(thehta float64, control, target int) *Qubit {
	return q.ControlledR(thehta, []int{control}, target)
}

// Controlled applies a controlled 2x2 unitary gate U to the target qubit.
func (q *Qubit) Controlled(g *matrix.Matrix, control []int, target int) *Qubit {
	n := q.NumQubits()

	var cmask int
	for _, c := range control {
		cmask |= 1 << (n - 1 - c)
	}
	tmask := 1 << (n - 1 - target)

	for i := range q.Dim() {
		if (i & cmask) != cmask {
			continue
		}

		j := i ^ tmask
		if i > j {
			continue
		}

		a, b := q.state.Data[i], q.state.Data[j]
		q.state.Data[i] = g.At(0, 0)*a + g.At(0, 1)*b
		q.state.Data[j] = g.At(1, 0)*a + g.At(1, 1)*b
	}

	return q
}

// ControlledU applies a controlled unitary operation.
func (q *Qubit) ControlledU(theta, phi, lambda float64, control []int, target int) *Qubit {
	n := q.NumQubits()

	var cmask int
	for _, c := range control {
		cmask |= 1 << (n - 1 - c)
	}
	tmask := 1 << (n - 1 - target)

	sin := cmplx.Sin(complex(theta/2, 0))
	cos := cmplx.Cos(complex(theta/2, 0))

	e0 := cmplx.Exp(complex(0, phi))
	e1 := cmplx.Exp(complex(0, lambda))
	e2 := cmplx.Exp(complex(0, phi+lambda))

	for i := range q.Dim() {
		if (i & cmask) != cmask {
			continue
		}

		j := i ^ tmask
		if i > j {
			continue
		}

		a, b := q.state.Data[i], q.state.Data[j]
		q.state.Data[i] = cos*a - e1*sin*b
		q.state.Data[j] = e0*sin*a + e2*cos*b
	}

	return q
}

// ControlledH applies the controlled Hadamard gate.
func (q *Qubit) ControlledH(control []int, target int) *Qubit {
	n := q.NumQubits()

	var cmask int
	for _, c := range control {
		cmask |= 1 << (n - 1 - c)
	}
	tmask := 1 << (n - 1 - target)

	// iterate over all states
	sqrt2 := complex(1/math.Sqrt2, 0)
	for i := range q.Dim() {
		if (i & cmask) != cmask {
			continue
		}

		j := i ^ tmask
		if i > j {
			continue
		}

		a, b := q.state.Data[i], q.state.Data[j]
		q.state.Data[i] = (a + b) * sqrt2
		q.state.Data[j] = (a - b) * sqrt2
	}

	return q
}

// ControlledX applies the controlled-X gate.
func (q *Qubit) ControlledX(control []int, target int) *Qubit {
	n := q.NumQubits()

	var cmask int
	for _, c := range control {
		cmask |= 1 << (n - 1 - c)
	}
	tmask := 1 << (n - 1 - target)

	// iterate over all states
	for i := 0; i < q.Dim(); i++ {
		if (i & cmask) != cmask {
			continue
		}

		j := i ^ tmask
		if i > j {
			continue
		}

		// swap
		q.state.Data[i], q.state.Data[j] = q.state.Data[j], q.state.Data[i]
	}

	return q
}

// ControlledZ applies the controlled-Z gate.
func (q *Qubit) ControlledZ(control []int, target int) *Qubit {
	n := q.NumQubits()

	var cmask int
	for _, c := range control {
		cmask |= 1 << (n - 1 - c)
	}

	tmask := 1 << (n - 1 - target)

	// iterate over all states
	for i := 0; i < q.Dim(); i++ {
		if (i & cmask) != cmask {
			continue
		}

		if (i & tmask) == tmask {
			q.state.Data[i] *= -1
		}
	}

	return q
}

// ControlledR applies a controlled rotation around the Z-axis.
func (q *Qubit) ControlledR(theta float64, control []int, target int) *Qubit {
	n := q.NumQubits()

	var cmask int
	for _, c := range control {
		cmask |= 1 << (n - 1 - c)
	}
	tmask := 1 << (n - 1 - target)

	// iterate over all states
	phase := cmplx.Exp(complex(0, theta))
	for i := 0; i < q.Dim(); i++ {
		if (i & cmask) != cmask {
			continue
		}

		if (i & tmask) == tmask {
			q.state.Data[i] *= phase
		}
	}

	return q
}

// Swap swaps the states of two qubits.
func (q *Qubit) Swap(i, j int) *Qubit {
	if i == j {
		return q
	}

	n := q.NumQubits()
	imask := 1 << (n - 1 - i)
	jmask := 1 << (n - 1 - j)

	for k := range q.Dim() {
		ibit := (k & imask) >> (n - 1 - i)
		jbit := (k & jmask) >> (n - 1 - j)

		if ibit == jbit {
			continue
		}

		l := k ^ (imask | jmask)
		if k > l {
			continue
		}

		q.state.Data[k], q.state.Data[l] = q.state.Data[l], q.state.Data[k]
	}

	return q
}

// QFT applies the quantum Fourier transform.
func (q *Qubit) QFT(idx ...int) *Qubit {
	if len(idx) == 0 {
		n := q.NumQubits()
		idx = make([]int, n)
		for i := range n {
			idx[i] = i
		}
	}

	for i := range idx {
		q.H(idx[i])

		k := 2
		for j := i + 1; j < len(idx); j++ {
			q.CR(gate.Theta(k), idx[i], idx[j])
			k++
		}
	}

	return q
}

// InvQFT applies the inverse quantum Fourier transform.
func (q *Qubit) InvQFT(idx ...int) *Qubit {
	if len(idx) == 0 {
		n := q.NumQubits()
		idx = make([]int, n)
		for i := range n {
			idx[i] = i
		}
	}

	len := len(idx)
	for i := len - 1; i > -1; i-- {
		k := len - i
		for j := len - 1; j > i; j-- {
			q.CR(-1*gate.Theta(k), idx[j], idx[i])
			k--
		}

		q.H(idx[i])
	}

	return q
}

// Set sets the state vector of the qubit and normalized it.
func (q *Qubit) Set(state *vector.Vector) {
	q.n = number.Log2(len(state.Data))
	q.state = state
	q.Normalize()
}

// Measure returns a measured qubit.
func (q *Qubit) Measure(idx int) *Qubit {
	n := q.NumQubits()
	mask := 1 << (n - 1 - idx)

	var prob0 float64
	for i := range q.Dim() {
		if (i & mask) == 0 {
			prob0 += math.Pow(cmplx.Abs(q.state.Data[i]), 2)
		}
	}

	collapse := func(q *Qubit, result int) {
		for i := range q.Dim() {
			if ((i & mask) >> (n - 1 - idx)) == result {
				continue
			}

			q.state.Data[i] = 0
		}

		q.Normalize()
	}

	if q.Rand() < prob0 {
		collapse(q, 0)
		return Zero()
	}

	collapse(q, 1)
	return One()
}

// Normalize returns a normalized qubit.
func (q *Qubit) Normalize() *Qubit {
	sum := number.Sum(q.Probability())
	norm := complex(1/math.Sqrt(sum), 0)
	q.state = q.state.Mul(norm)
	return q
}

// Clone returns a clone of q.
func (q *Qubit) Clone() *Qubit {
	return &Qubit{
		n:     q.n,
		state: q.state.Clone(),
		Rand:  q.Rand,
	}
}

// Equal returns true if q and qb are equal.
func (q *Qubit) Equal(qb *Qubit, tol ...float64) bool {
	return q.state.Equal(qb.state, tol...)
}

// BinaryString measures the quantum state and returns its binary string representation.
func (q *Qubit) BinaryString() string {
	c := q.Clone()
	var sb strings.Builder
	for i := range q.NumQubits() {
		if c.Measure(i).IsZero() {
			sb.WriteByte('0')
			continue
		}

		sb.WriteByte('1')
	}

	return sb.String()
}

// Int measures the quantum state and returns its int representation.
func (q *Qubit) Int() int64 {
	return number.Must(strconv.ParseInt(q.BinaryString(), 2, 0))
}

// String returns the string representation of q.
func (q *Qubit) String() string {
	return q.state.String()
}

// State returns the state of the qubit at the given index.
// If no index is provided, it returns the state vector of all qubits.
func (q *Qubit) State(idx ...[]int) []State {
	n := q.NumQubits()
	if len(idx) < 1 {
		idx = make([][]int, 1)
		idx[0] = make([]int, n)
		for i := range n {
			idx[0][i] = i
		}
	}

	var state []State
	for i, a := range q.Amplitude() {
		amp, isZero := round(a)
		if isZero {
			continue
		}

		var binary []string
		for _, j := range idx {
			binary = append(binary, take(n, i, j))
		}

		state = append(state, NewState(amp, binary...))
	}

	return state
}

func round(a complex128, tol ...float64) (complex128, bool) {
	rz := epsilon.IsZeroF64(math.Abs(real(a)), tol...)
	iz := epsilon.IsZeroF64(math.Abs(imag(a)), tol...)
	if rz && iz {
		return complex(0, 0), true
	}

	if rz {
		a = complex(0, imag(a))
	}

	if iz {
		a = complex(real(a), 0)
	}

	return a, false
}

func take(n, i int, idx []int) string {
	var sb strings.Builder
	for _, j := range idx {
		if (i & (1 << (n - 1 - j))) == 0 {
			sb.WriteByte('0')
			continue
		}

		sb.WriteByte('1')
	}

	return sb.String()
}

func TensorProduct(qb ...*Qubit) *Qubit {
	q := qb[0]
	for i := 1; i < len(qb); i++ {
		q = q.TensorProduct(qb[i])
	}

	return q
}
