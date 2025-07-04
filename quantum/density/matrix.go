package density

import (
	"math"
	"math/cmplx"
	"strconv"
	"strings"

	"github.com/itsubaki/q/math/epsilon"
	"github.com/itsubaki/q/math/matrix"
	"github.com/itsubaki/q/math/number"
	"github.com/itsubaki/q/quantum/gate"
	"github.com/itsubaki/q/quantum/qubit"
)

// Qubit is a quantum bit.
type Qubit int

// Index returns the index of qubit.
func (q Qubit) Index() int {
	return int(q)
}

// Matrix is a density matrix.
type Matrix struct {
	m *matrix.Matrix
}

// New returns a new density matrix.
func New(ensemble []State) *Matrix {
	var m *matrix.Matrix
	for _, s := range Normalize(ensemble) {
		n := s.Qubit.Dimension()
		if m == nil {
			m = matrix.Zero(n, n)
		}

		op := s.Qubit.OuterProduct(s.Qubit).Mul(complex(s.Probability, 0))
		for i := range n {
			for j := range n {
				m.AddAt(i, j, op.At(i, j))
			}
		}
	}

	return &Matrix{
		m: m,
	}
}

// NewPure returns a new pure density matrix.
func NewPure(qb *qubit.Qubit) *Matrix {
	return New([]State{
		{
			Probability: 1.0,
			Qubit:       qb,
		},
	})
}

func (m *Matrix) At(i, j int) complex128 {
	return m.m.At(i, j)
}

// Qubits returns the qubits of the density matrix.
func (m *Matrix) Qubits() []Qubit {
	n := m.NumQubits()

	qubits := make([]Qubit, n)
	for i := range n {
		qubits[i] = Qubit(i)
	}

	return qubits
}

// Underlying returns the internal matrix.
func (m *Matrix) Underlying() *matrix.Matrix {
	return m.m
}

// Dimension returns the dimension of the density matrix.
func (m *Matrix) Dimension() (rows int, cols int) {
	return m.m.Dimension()
}

// IsPure returns true if the density matrix is pure.
func (m *Matrix) IsPure(eps ...float64) bool {
	return math.Abs(1-m.Purity()) < epsilon.E13(eps...)
}

// IsMixed returns true if the density matrix is mixed.
func (m *Matrix) IsMixed(eps ...float64) bool {
	return !m.IsPure(eps...)
}

// IsHermite returns true if the density matrix is Hermitian.
func (m *Matrix) IsHermite(eps ...float64) bool {
	return m.m.IsHermite(eps...)
}

// IsZero returns true if the density matrix is zero.
func (m *Matrix) IsZero(eps ...float64) bool {
	e := epsilon.E13(eps...)
	for i := range m.m.Data {
		if cmplx.Abs(m.m.Data[i]) > e {
			return false
		}
	}

	return true
}

// NumQubits returns the number of qubits.
func (m *Matrix) NumQubits() int {
	p, _ := m.Dimension()
	return number.Log2(p)
}

// Apply applies a unitary matrix to the density matrix.
func (m *Matrix) Apply(u *matrix.Matrix) *Matrix {
	m.m = matrix.MatMul(u, m.m, u.Dagger())
	return m
}

// Probability returns the probability of the qubit in the given state.
func (m *Matrix) Probability(q *qubit.Qubit) float64 {
	p := q.OuterProduct(q)
	return real(m.m.Apply(p).Trace())
}

// Project returns the projection of the density matrix onto the given qubit.
func (m *Matrix) Project(q *qubit.Qubit, eps ...float64) *Matrix {
	p := q.OuterProduct(q)
	tr := m.m.Apply(p).Trace()

	if cmplx.Abs(tr) < epsilon.E13(eps...) {
		return &Matrix{
			m: matrix.ZeroLike(m.m),
		}
	}

	pmp := matrix.Apply(p, m.m, p)
	return &Matrix{
		m: pmp.Mul(1.0 / tr),
	}
}

// ExpectationValue returns the expectation value of the given operator.
func (m *Matrix) ExpectedValue(u *matrix.Matrix) float64 {
	return real(m.m.Apply(u).Trace())
}

// Trace returns the trace of the density matrix.
func (m *Matrix) Trace() float64 {
	return real(m.m.Trace())
}

// Purity returns the purity of the density matrix, defined as Tr(rho^2).
func (m *Matrix) Purity() float64 {
	return real(m.m.Apply(m.m).Trace())
}

// TensorProduct returns the tensor product of two density matrices.
func (m *Matrix) TensorProduct(n *Matrix) *Matrix {
	return &Matrix{
		m: m.m.TensorProduct(n.m),
	}
}

// PartialTrace returns the partial trace of the density matrix.
// The length of index must be less than or equal to n - 1,
// where n is the number of qubits in the matrix.
func (m *Matrix) PartialTrace(index ...Qubit) *Matrix {
	n := m.NumQubits()
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

			r := int(number.Must(strconv.ParseInt(kr, 2, 0)))
			c := int(number.Must(strconv.ParseInt(lr, 2, 0)))
			out.AddAt(r, c, m.m.At(i, j))

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
// p should be between 0 and 1, representing the probability of depolarization.
func (m *Matrix) Depolarizing(p float64) *Matrix {

	n := m.NumQubits()
	d := 1 << n

	id := gate.I(n).Mul(complex(1.0/float64(d), 0))
	i := id.Mul(complex(p, 0))
	r := m.m.Mul(complex(1-p, 0))

	return &Matrix{
		m: i.Add(r),
	}
}

// Flip applies a flip channel to the density matrix.
func (m *Matrix) Flip(p float64, qb Qubit, g *matrix.Matrix) *Matrix {
	n := m.NumQubits()
	ops := make([]*matrix.Matrix, n)
	for i := range n {
		if i == qb.Index() {
			ops[i] = g
			continue
		}

		ops[i] = gate.I()
	}

	e0 := gate.I(n).Mul(complex(math.Sqrt(p), 0))
	e1 := matrix.TensorProduct(ops...).Mul(complex(math.Sqrt(1-p), 0))

	rho0 := matrix.MatMul(e0, m.m, e0.Dagger())
	rho1 := matrix.MatMul(e1, m.m, e1.Dagger())

	return &Matrix{
		m: rho0.Add(rho1),
	}
}

// BitFlip applies a bit flip channel to the density matrix.
func (m *Matrix) BitFlip(p float64, qb Qubit) *Matrix {
	return m.Flip(p, qb, gate.X())
}

// PhaseFlip applies a phase flip channel to the density matrix.
func (m *Matrix) PhaseFlip(p float64, qb Qubit) *Matrix {
	return m.Flip(p, qb, gate.Z())
}

// BitPhaseFlip applies a bit-phase flip channel to the density matrix.
func (m *Matrix) BitPhaseFlip(p float64, qb Qubit) *Matrix {
	return m.Flip(p, qb, gate.Y())
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
