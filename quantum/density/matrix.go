package density

import (
	"fmt"
	"iter"
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
	rho *matrix.Matrix
}

// New returns a new density matrix.
func New(ensemble []State) *Matrix {
	normalized := Normalize(ensemble)
	n := normalized[0].Qubit.Dim()

	rho := matrix.Zero(n, n)
	for _, s := range normalized {
		op := s.Qubit.OuterProduct(s.Qubit)
		rho = rho.Add(op.Mul(complex(s.Probability, 0)))
	}

	return &Matrix{
		rho: rho,
	}
}

// NewPureState returns a new pure state density matrix for the given qubit.
func NewPureState(qb *qubit.Qubit) *Matrix {
	return New([]State{
		{
			Probability: 1.0,
			Qubit:       qb,
		},
	})
}

// NewZeroState returns a new zero state density matrix for the given number of qubits.
func NewZeroState(n ...int) *Matrix {
	return NewPureState(qubit.Zero(n...))
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

func (m *Matrix) ComputationalBasis() []*qubit.Qubit {
	n := m.NumQubits()

	basis := make([]*qubit.Qubit, 1<<n)
	for i := range 1 << n {
		basis[i] = qubit.NewFrom(fmt.Sprintf("%0*b", n, i))
	}

	return basis
}

// At returns a value of matrix at (i,j).
func (m *Matrix) At(i, j int) complex128 {
	return m.rho.At(i, j)
}

// Underlying returns the internal matrix.
func (m *Matrix) Underlying() *matrix.Matrix {
	return m.rho
}

// Seq2 returns a sequence of rows.
func (m *Matrix) Seq2() iter.Seq2[int, []complex128] {
	return m.rho.Seq2()
}

// Dim returns the dimension of the density matrix.
func (m *Matrix) Dim() (rows int, cols int) {
	return m.rho.Dim()
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
	return m.rho.IsHermite(eps...)
}

// IsZero returns true if the density matrix is zero.
func (m *Matrix) IsZero(eps ...float64) bool {
	e := epsilon.E13(eps...)
	for i := range m.rho.Data {
		if cmplx.Abs(m.rho.Data[i]) > e {
			return false
		}
	}

	return true
}

// NumQubits returns the number of qubits.
func (m *Matrix) NumQubits() int {
	p, _ := m.Dim()
	return number.Log2(p)
}

// Apply applies a unitary matrix to the density matrix.
func (m *Matrix) Apply(u *matrix.Matrix) *Matrix {
	m.rho = matrix.MatMul(u, m.rho, u.Dagger())
	return m
}

// Probability returns the probability of the qubit in the given state.
func (m *Matrix) Probability(q *qubit.Qubit) float64 {
	p := q.OuterProduct(q)
	tr := matrix.MatMul(m.rho, p).Trace()
	return real(tr)
}

// Project returns the probability and post-measurement density matrix.
func (m *Matrix) Project(q *qubit.Qubit, eps ...float64) (float64, *Matrix) {
	p := m.Probability(q)
	if math.Abs(p) < epsilon.E13(eps...) {
		return 0, &Matrix{
			rho: matrix.ZeroLike(m.rho),
		}
	}

	op := q.OuterProduct(q)
	rho := matrix.MatMul(op, m.rho, op)
	return p, &Matrix{
		rho: rho.Mul(1.0 / complex(p, 0)),
	}
}

// ExpectedValue returns the expectation value of the given operator.
func (m *Matrix) ExpectedValue(u *matrix.Matrix) float64 {
	return real(m.rho.Apply(u).Trace())
}

// Trace returns the trace of the density matrix.
func (m *Matrix) Trace() float64 {
	return real(m.rho.Trace())
}

// Purity returns the purity of the density matrix, defined as Tr(rho^2).
func (m *Matrix) Purity() float64 {
	return real(matrix.MatMul(m.rho, m.rho).Trace())
}

// TensorProduct returns the tensor product of two density matrices.
func (m *Matrix) TensorProduct(n *Matrix) *Matrix {
	return &Matrix{
		rho: m.rho.TensorProduct(n.rho),
	}
}

// PartialTrace returns the partial trace of the density matrix.
// The length of index must be less than or equal to n - 1,
// where n is the number of qubits in the matrix.
func (m *Matrix) PartialTrace(index ...Qubit) *Matrix {
	n := m.NumQubits()
	p, q := m.Dim()
	d := number.Pow(2, n-1)

	rho := matrix.Zero(d, d)
	for i := range p {
		k, kr := take(n, i, index)

		for j := range q {
			l, lr := take(n, j, index)

			if k != l {
				continue
			}

			r := int(number.Must(strconv.ParseInt(kr, 2, 0)))
			c := int(number.Must(strconv.ParseInt(lr, 2, 0)))
			rho.AddAt(r, c, m.At(i, j))

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

	return &Matrix{rho: rho}
}

// Depolarizing returns the depolarizing channel.
// It applies the identity with probability (1 - p),
// and applies each of the Pauli gates X, Y, and Z with probability p/3.
func (m *Matrix) Depolarizing(p float64) *Matrix {
	n := m.NumQubits()

	id := m.rho.Mul(complex(1-p, 0))
	x := matrix.MatMul(gate.X(n), m.rho, gate.X(n)).Mul(complex(p/3, 0))
	y := matrix.MatMul(gate.Y(n), m.rho, gate.Y(n)).Mul(complex(p/3, 0))
	z := matrix.MatMul(gate.Z(n), m.rho, gate.Z(n)).Mul(complex(p/3, 0))

	return &Matrix{
		rho: id.Add(x).Add(y).Add(z),
	}
}

// ApplyChannel applies a channel to the density matrix.
// It applies the identity with probability 1-p, and applies the gate g with probability p.
func (m *Matrix) ApplyChannel(p float64, g *matrix.Matrix, qb ...Qubit) *Matrix {
	n, k := m.NumQubits(), len(qb)
	id := gate.I()
	e0 := gate.I().Mul(complex(math.Sqrt(1-p), 0))
	e1 := g.Mul(complex(math.Sqrt(p), 0))

	index := make([]int, k)
	for i, v := range qb {
		index[i] = v.Index()
	}

	// create kraus operators
	kraus := make([]*matrix.Matrix, 1<<k)
	for i := range 1 << k {
		ops := make([]*matrix.Matrix, n)
		for j := range n {
			ops[j] = id
		}

		for j, idx := range index {
			if (i>>j)&1 == 1 {
				ops[idx] = e1
				continue
			}

			ops[idx] = e0
		}

		kraus[i] = matrix.TensorProduct(ops...)
	}

	// E(rho) = sum(E * rho * E^dagger)
	rho := matrix.ZeroLike(m.rho)
	for _, e := range kraus {
		rho = rho.Add(matrix.MatMul(e, m.rho, e.Dagger()))
	}

	return &Matrix{
		rho: rho,
	}
}

// BitFlip applies a bit flip channel to the density matrix.
func (m *Matrix) BitFlip(p float64, qb Qubit) *Matrix {
	return m.ApplyChannel(p, gate.X(), qb)
}

// BitPhaseFlip applies a bit-phase flip channel to the density matrix.
func (m *Matrix) BitPhaseFlip(p float64, qb Qubit) *Matrix {
	return m.ApplyChannel(p, gate.Y(), qb)
}

// PhaseFlip applies a phase flip channel to the density matrix.
func (m *Matrix) PhaseFlip(p float64, qb Qubit) *Matrix {
	return m.ApplyChannel(p, gate.Z(), qb)
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
