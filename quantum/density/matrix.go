package density

import (
	"errors"
	"strconv"
	"strings"

	"github.com/itsubaki/q/math/matrix"
	"github.com/itsubaki/q/math/number"
	"github.com/itsubaki/q/quantum/gate"
	"github.com/itsubaki/q/quantum/qubit"
)

var (
	ErrInvalidRange     = errors.New("p must be 0 <= p =< 1")
	ErrInvalidDimension = errors.New("invalid dimension")
)

// Matrix is a density matrix.
type Matrix struct {
	m matrix.Matrix
}

// New returns a new density matrix.
func New(ensemble []State) *Matrix {
	m := &Matrix{matrix.New()}

	for _, s := range Normalize(ensemble) {
		n := s.Qubit.Dimension()
		if len(m.m) < 1 {
			m.m = matrix.Zero(n, n)
		}

		op := s.Qubit.OuterProduct(s.Qubit).Mul(complex(s.Probability, 0))
		for i := range n {
			for j := range n {
				m.m[i][j] = m.m[i][j] + op[i][j]
			}
		}
	}

	return m
}

// Raw returns the raw matrix.
func (m *Matrix) Raw() matrix.Matrix {
	return m.m
}

// Dimension returns the dimension of the density matrix.
func (m *Matrix) Dimension() (int, int) {
	return len(m.m), len(m.m[0])
}

// NumberOfBit returns the number of qubits.
func (m *Matrix) NumberOfBit() int {
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
func (m *Matrix) PartialTrace(index ...int) *Matrix {
	n := m.NumberOfBit()
	p, q := m.Dimension()

	d := number.Pow(2, n-1)
	out := matrix.Zero(d, d)
	for i := range p {
		k, kr := take(n, i, index)

		for j := range q {
			l, lr := take(n, j, index)

			if k != l {
				continue
			}

			r := number.Must(strconv.ParseInt(kr, 2, 0))
			c := number.Must(strconv.ParseInt(lr, 2, 0))
			out[r][c] = out[r][c] + m.m[i][j]

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

	return &Matrix{m: out}
}

// Depolarizing returns the depolarizing channel.
func (m *Matrix) Depolarizing(p float64) (*Matrix, error) {
	if p < 0 || p > 1 {
		return nil, ErrInvalidRange
	}

	n := m.NumberOfBit()
	i := gate.I(n).Mul(complex(p/2, 0))
	r := m.m.Mul(complex(1-p, 0))
	return &Matrix{i.Add(r)}, nil
}

func take(n, i int, index []int) (string, string) {
	idx := make(map[int]struct{}, len(index))
	for _, j := range index {
		idx[j] = struct{}{}
	}

	var out, remain strings.Builder
	for j := 0; j < n; j++ {
		b := byte('0' + ((i >> (n - 1 - j)) & 1))
		if _, ok := idx[j]; ok {
			out.WriteByte(b)
			continue
		}

		remain.WriteByte(b)
	}

	return out.String(), remain.String()
}
