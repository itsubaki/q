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
)

// Qubit is a qubit.
type Qubit struct {
	state *vector.Vector
	Rand  func() float64 // Random number generator
}

// New returns a new qubit.
func New(v *vector.Vector) *Qubit {
	q := &Qubit{
		state: v,
		Rand:  rand.Float64,
	}

	q.Normalize()
	return q
}

// Zero returns a qubit in the zero state.
func Zero(n ...int) *Qubit {
	return New(vector.TensorProductN(
		vector.New(1, 0),
		n...,
	))
}

// One returns a qubit in the one state.
func One(n ...int) *Qubit {
	return New(vector.TensorProductN(
		vector.New(0, 1),
		n...,
	))
}

// Plus returns a qubit in the plus state.
// The plus state is defined as (|0> + |1>) / sqrt(2).
func Plus(n ...int) *Qubit {
	return New(vector.TensorProductN(vector.New(
		complex(1/math.Sqrt2, 0),
		complex(1/math.Sqrt2, 0),
	), n...))
}

// Minus returns a qubit in the minus state.
// The minus state is defined as (|0> - |1>) / sqrt(2).
func Minus(n ...int) *Qubit {
	return New(vector.TensorProductN(vector.New(
		complex(1/math.Sqrt2, 0),
		complex(1/math.Sqrt2, 0)*-1,
	), n...))
}

// NewFrom returns a new qubit from a binary string.
func NewFrom(binary string) *Qubit {
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
	return number.Log2(q.Dim())
}

// IsZero returns true if q is zero qubit.
func (q *Qubit) IsZero(eps ...float64) bool {
	return q.Equals(Zero(), eps...)
}

// IsOne returns true if q is one qubit.
func (q *Qubit) IsOne(eps ...float64) bool {
	return q.Equals(One(), eps...)
}

// Dim returns the dimension of q.
func (q *Qubit) Dim() int {
	return len(q.state.Data)
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
	q.state = q.state.TensorProduct(qb.state)
	return q
}

// Apply returns a qubit that is applied m.
func (q *Qubit) Apply(m ...*matrix.Matrix) *Qubit {
	for _, v := range m {
		q.state = q.state.Apply(v)
	}

	return q
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
		state: q.state.Clone(),
		Rand:  q.Rand,
	}
}

// Equals returns true if q and qb are equal.
func (q *Qubit) Equals(qb *Qubit, eps ...float64) bool {
	return q.state.Equals(qb.state, eps...)
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
	if len(idx) < 1 {
		n := q.NumQubits()
		all := make([]int, n)
		for i := range n {
			all[i] = i
		}

		idx = append(idx, all)
	}

	n := q.NumQubits()
	state := make([]State, 0)
	for i, a := range q.Amplitude() {
		amp, ok := round(a)
		if !ok {
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

func round(a complex128, eps ...float64) (complex128, bool) {
	e := epsilon.E13(eps...)
	r, i := math.Abs(real(a)), math.Abs(imag(a))

	if r < e && i < e {
		return complex(0, 0), false
	}

	if r < e {
		a = complex(0, imag(a))
	}

	if i < e {
		a = complex(real(a), 0)
	}

	return a, true
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
