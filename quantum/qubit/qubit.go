package qubit

import (
	"fmt"
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
	vector vector.Vector
	Rand   func() float64 // Random number generator
}

// New returns a new qubit.
// z is a vector of complex128.
func New(z ...complex128) *Qubit {
	q := &Qubit{
		vector: vector.New(z...),
		Rand:   rand.Float64,
	}

	q.Normalize()
	return q
}

// Zero returns a qubit in the zero state.
// n is the number of qubits.
func Zero(n ...int) *Qubit {
	v := vector.TensorProductN(vector.Vector{1, 0}, n...)
	return New(v.Complex()...)
}

// One returns a qubit in the one state.
// n is the number of qubits.
func One(n ...int) *Qubit {
	v := vector.TensorProductN(vector.Vector{0, 1}, n...)
	return New(v.Complex()...)
}

// NumberOfBit returns the number of qubits.
func (q *Qubit) NumberOfBit() int {
	return number.Log2(q.Dimension())
}

// IsZero returns true if q is zero qubit.
func (q *Qubit) IsZero(eps ...float64) bool {
	return q.Equals(Zero(), eps...)
}

// IsOne returns true if q is one qubit.
func (q *Qubit) IsOne(eps ...float64) bool {
	return q.Equals(One(), eps...)
}

// InnerProduct returns the inner product of q and qb.
func (q *Qubit) InnerProduct(qb *Qubit) complex128 {
	return q.vector.InnerProduct(qb.vector)
}

// OuterProduct returns the outer product of q and qb.
func (q *Qubit) OuterProduct(qb *Qubit) matrix.Matrix {
	return q.vector.OuterProduct(qb.vector)
}

// Dimension returns the dimension of q.
func (q *Qubit) Dimension() int {
	return q.vector.Dimension()
}

// Clone returns a clone of q.
func (q *Qubit) Clone() *Qubit {
	return &Qubit{
		vector: q.vector.Clone(),
		Rand:   q.Rand,
	}
}

// Fidelity returns the fidelity of q and qb.
func (q *Qubit) Fidelity(qb *Qubit) float64 {
	p0 := qb.Probability()
	p1 := q.Probability()

	var sum float64
	for i := range p0 {
		sum = sum + math.Sqrt(p0[i]*p1[i])
	}

	return sum
}

// TraceDistance returns the trace distance of q and qb.
func (q *Qubit) TraceDistance(qb *Qubit) float64 {
	p0 := qb.Probability()
	p1 := q.Probability()

	var sum float64
	for i := range p0 {
		sum = sum + math.Abs(p0[i]-p1[i])
	}

	return sum / 2
}

// Equals returns true if q and qb are equal.
func (q *Qubit) Equals(qb *Qubit, eps ...float64) bool {
	return q.vector.Equals(qb.vector, eps...)
}

// TensorProduct returns the tensor product of q and qb.
func (q *Qubit) TensorProduct(qb *Qubit) *Qubit {
	q.vector = q.vector.TensorProduct(qb.vector)
	return q
}

// Apply returns a qubit that is applied m.
func (q *Qubit) Apply(m ...matrix.Matrix) *Qubit {
	for _, mm := range m {
		q.vector = q.vector.Apply(mm)
	}

	return q
}

// Normalize returns a normalized qubit.
func (q *Qubit) Normalize() *Qubit {
	sum := number.Sum(q.Probability())
	z := 1 / math.Sqrt(sum)
	q.vector = q.vector.Mul(complex(z, 0))
	return q
}

// Amplitude returns the amplitude of q.
func (q *Qubit) Amplitude() []complex128 {
	return q.vector.Complex()
}

// Probability returns the probability of q.
func (q *Qubit) Probability() []float64 {
	p := make([]float64, len(q.Amplitude()))
	for i, a := range q.Amplitude() {
		p[i] = math.Pow(cmplx.Abs(a), 2)
	}

	return p
}

// Measure returns a measured qubit.
func (q *Qubit) Measure(index int) *Qubit {
	n := q.NumberOfBit()
	mask := 1 << (n - 1 - index)

	zidx, zprop := make([]int, 0), make([]float64, 0)
	oidx := make([]int, 0)
	for i, p := range q.Probability() {
		if i&mask == 0 {
			zidx, zprop = append(zidx, i), append(zprop, p)
			continue
		}

		oidx = append(oidx, i)
	}

	// One()
	if q.Rand() > number.Sum(zprop) {
		for _, i := range zidx {
			q.vector[i] = complex(0, 0)
		}

		q.Normalize()
		return One()
	}

	// Zero()
	for _, i := range oidx {
		q.vector[i] = complex(0, 0)
	}

	q.Normalize()
	return Zero()
}

// Int measures the quantum state and returns its int representation.
func (q *Qubit) Int() int64 {
	return number.Must(strconv.ParseInt(q.BinaryString(), 2, 0))
}

// BinaryString measures the quantum state and returns its binary string representation.
func (q *Qubit) BinaryString() string {
	c := q.Clone()

	var sb strings.Builder
	for i := range q.NumberOfBit() {
		if c.Measure(i).IsZero() {
			sb.WriteString("0")
			continue
		}

		sb.WriteString("1")
	}

	return sb.String()
}

// String returns the string representation of q.
func (q *Qubit) String() string {
	return fmt.Sprintf("%v", q.vector)
}

// State returns the state of q with index.
func (q *Qubit) State(index ...[]int) []State {
	if len(index) < 1 {
		n := q.NumberOfBit()
		idx := make([]int, n)
		for i := range n {
			idx[i] = i
		}

		index = append(index, idx)
	}

	n := q.NumberOfBit()
	state := make([]State, 0)
	for i, a := range q.Amplitude() {
		amp := round(a)
		if amp == 0 {
			continue
		}

		var bin []string
		for _, idx := range index {
			bin = append(bin, take(n, i, idx))
		}

		state = append(state, NewState(amp, bin...))
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
	s := fmt.Sprintf("%0*b", n, i)

	var sb strings.Builder
	for _, i := range index {
		sb.WriteString(s[i : i+1])
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
