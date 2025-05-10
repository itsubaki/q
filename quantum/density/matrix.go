package density

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/itsubaki/q/math/matrix"
	"github.com/itsubaki/q/math/number"
	"github.com/itsubaki/q/quantum/gate"
	"github.com/itsubaki/q/quantum/qubit"
)

var ErrInvalidRange = errors.New("p is out of range [0,1]")

// Qubit is a quantum bit.
type Qubit int

// Index returns the index of qubit.
func (q Qubit) Index() int {
	return int(q)
}

// Matrix is a density matrix.
type Matrix struct {
	m matrix.Matrix
}

// New returns a new density matrix.
func New(ensemble []State) *Matrix {
	m := matrix.New()
	for _, s := range Normalize(ensemble) {
		n := s.Qubit.Dimension()
		if len(m.Data) < 1 {
			m = matrix.Zero(n, n)
		}

		op := s.Qubit.OuterProduct(s.Qubit).Mul(complex(s.Probability, 0))
		for i := range n {
			for j := range n {
				m.Set(i, j, m.At(i, j)+op.At(i, j))
			}
		}
	}

	return &Matrix{
		m: m,
	}
}

func (m *Matrix) At(i, j int) complex128 {
	return m.m.At(i, j)
}

// Qubits returns the qubits of the density matrix.
func (m *Matrix) Qubits() []Qubit {
	var qubits []Qubit
	for i := range m.NumQubits() {
		qubits = append(qubits, Qubit(i))
	}

	return qubits
}

// Underlying returns the internal matrix.
func (m *Matrix) Underlying() matrix.Matrix {
	return m.m
}

// Dimension returns the dimension of the density matrix.
func (m *Matrix) Dimension() (int, int) {
	return m.m.Dimension()
}

// NumQubits returns the number of qubits.
func (m *Matrix) NumQubits() int {
	p, _ := m.Dimension()
	return number.Log2(p)
}

// Apply applies a unitary matrix to the density matrix.
func (m *Matrix) Apply(u matrix.Matrix) *Matrix {
	m.m = u.Dagger().Apply(m.m).Apply(u)
	return m
}

// Measure returns the probability of measuring the qubit in the given state.
func (m *Matrix) Measure(q *qubit.Qubit) float64 {
	return real(m.m.Apply(q.OuterProduct(q)).Trace())
}

// ExpectationValue returns the expectation value of the given operator.
func (m *Matrix) ExpectedValue(u matrix.Matrix) float64 {
	return real(m.m.Apply(u).Trace())
}

// Trace returns the trace of the density matrix.
func (m *Matrix) Trace() float64 {
	return real(m.m.Trace())
}

// SquareTrace returns the square trace of the density matrix.
func (m *Matrix) SquareTrace() float64 {
	return real(m.m.Apply(m.m).Trace())
}

// PartialTrace returns the partial trace of the density matrix.
func (m *Matrix) PartialTrace(index ...Qubit) (*Matrix, error) {
	n := m.NumQubits()
	p, q := m.Dimension()
	d := number.Pow(2, n-1)
	if len(index) > n-1 {
		return nil, fmt.Errorf("length of index must be less than or equal to %d", n-1)
	}

	out := matrix.Zero(d, d)
	for i := range p {
		k, kr := take(n, i, index)

		for j := range q {
			l, lr := take(n, j, index)

			if k != l {
				continue
			}

			r := int(number.Must(strconv.ParseInt(kr, 2, 0)))
			c := int(number.Must(strconv.ParseInt(lr, 2, 0)))
			out.Set(r, c, out.At(r, c)+m.m.At(i, j))

			// fmt.Printf("[%v][%v] = [%v][%v] + [%v][%v]\n", r, c, r, c, i, j)
			//
			// 4x4 explicit
			// index -> 0
			// out[0][0] = m.m[0][0] + m.m[2][2]
			// out[0][1] = m.m[0][1] + m.m[2][3]
			// out[1][0] = m.m[1][0] + m.m[3][2]
			// out[1][1] = m.m[1][1] + m.m[3][3]
			//
			// index -> 1
			// out[0][0] = m.m[0][0] + m.m[1][1]
			// out[0][1] = m.m[0][2] + m.m[1][3]
			// out[1][0] = m.m[2][0] + m.m[3][1]
			// out[1][1] = m.m[2][2] + m.m[3][3]
		}
	}

	return &Matrix{m: out}, nil
}

// Depolarizing returns the depolarizing channel.
func (m *Matrix) Depolarizing(p float64) (*Matrix, error) {
	if p < 0 || p > 1 {
		return nil, ErrInvalidRange
	}

	n := m.NumQubits()
	i := gate.I(n).Mul(complex(p/2, 0))
	r := m.m.Mul(complex(1-p, 0))
	return &Matrix{i.Add(r)}, nil
}

func take(n, i int, index []Qubit) (string, string) {
	idx := make(map[int]struct{}, len(index))
	for _, j := range index {
		idx[j.Index()] = struct{}{}
	}

	var out, remain strings.Builder
	for bit := range n {
		b := byte('0' + ((i >> (n - 1 - bit)) & 1))
		if _, ok := idx[bit]; ok {
			out.WriteByte(b)
			continue
		}

		remain.WriteByte(b)
	}

	return out.String(), remain.String()
}

func Must[T any](v T, err error) T {
	return number.Must(v, err)
}
