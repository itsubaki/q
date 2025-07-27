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
	vec  *vector.Vector
	Rand func() float64 // Random number generator
}

// New returns a new qubit.
func New(z ...complex128) *Qubit {
	q := &Qubit{
		vec:  vector.New(z...),
		Rand: rand.Float64,
	}

	q.Normalize()
	return q
}

// Zero returns a qubit in the zero state.
func Zero(n ...int) *Qubit {
	z := vector.New(1, 0)
	v := vector.TensorProductN(z, n...)
	return New(v.Data...)
}

// One returns a qubit in the one state.
func One(n ...int) *Qubit {
	o := vector.New(0, 1)
	v := vector.TensorProductN(o, n...)
	return New(v.Data...)
}

// Plus returns a qubit in the plus state.
// The plus state is defined as (|0> + |1>) / sqrt(2).
func Plus(n ...int) *Qubit {
	plus := vector.New(
		complex(1/math.Sqrt2, 0),
		complex(1/math.Sqrt2, 0),
	)

	v := vector.TensorProductN(plus, n...)
	return New(v.Data...)
}

// Minus returns a qubit in the minus state.
// The minus state is defined as (|0> - |1>) / sqrt(2).
func Minus(n ...int) *Qubit {
	minus := vector.New(
		complex(1/math.Sqrt2, 0),
		-1*complex(1/math.Sqrt2, 0),
	)

	v := vector.TensorProductN(minus, n...)
	return New(v.Data...)
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
	return len(q.vec.Data)
}

// Amplitude returns the amplitude of q.
func (q *Qubit) Amplitude() []complex128 {
	return q.vec.Data
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
	return q.vec.InnerProduct(qb.vec)
}

// OuterProduct returns the outer product of q and qb.
func (q *Qubit) OuterProduct(qb *Qubit) *matrix.Matrix {
	return q.vec.OuterProduct(qb.vec)
}

// TensorProduct returns the tensor product of q and qb.
func (q *Qubit) TensorProduct(qb *Qubit) *Qubit {
	q.vec = q.vec.TensorProduct(qb.vec)
	return q
}

// Apply returns a qubit that is applied m.
func (q *Qubit) Apply(m ...*matrix.Matrix) *Qubit {
	for _, mm := range m {
		q.vec = q.vec.Apply(mm)
	}

	return q
}

// Measure returns a measured qubit.
func (q *Qubit) Measure(index int) *Qubit {
	n := q.NumQubits()
	mask := 1 << (n - 1 - index)

	var prob0 float64
	for i := range q.Dim() {
		if (i & mask) == 0 {
			prob0 += math.Pow(cmplx.Abs(q.vec.Data[i]), 2)
		}
	}

	collapse := func(q *Qubit, result int) {
		for i := range q.Dim() {
			if ((i & mask) >> (n - 1 - index)) == result {
				continue
			}

			q.vec.Data[i] = 0
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
	q.vec = q.vec.Mul(norm)
	return q
}

// Clone returns a clone of q.
func (q *Qubit) Clone() *Qubit {
	return &Qubit{
		vec:  q.vec.Clone(),
		Rand: q.Rand,
	}
}

// Equals returns true if q and qb are equal.
func (q *Qubit) Equals(qb *Qubit, eps ...float64) bool {
	return q.vec.Equals(qb.vec, eps...)
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
	return q.vec.String()
}

// State returns the state of the qubit at the given index.
// If no index is provided, it returns the state vector of all qubits.
func (q *Qubit) State(index ...[]int) []State {
	if len(index) < 1 {
		n := q.NumQubits()
		idx := make([]int, n)
		for i := range n {
			idx[i] = i
		}

		index = append(index, idx)
	}

	n := q.NumQubits()
	state := make([]State, 0)
	for i, a := range q.Amplitude() {
		amp := round(a)
		if cmplx.Abs(amp) < epsilon.E13() {
			continue
		}

		var binary []string
		for _, idx := range index {
			binary = append(binary, take(n, i, idx))
		}

		state = append(state, NewState(amp, binary...))
	}

	return state
}

func round(a complex128, eps ...float64) complex128 {
	e := epsilon.E13(eps...)

	if math.Abs(real(a)) < e {
		a = complex(0, imag(a))
	}

	if math.Abs(imag(a)) < e {
		a = complex(real(a), 0)
	}

	return a
}

func take(n, i int, index []int) string {
	var sb strings.Builder
	for _, bit := range index {
		if (i & (1 << (n - 1 - bit))) == 0 {
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
