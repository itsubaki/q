package density

import (
	"iter"
	"math"

	"github.com/itsubaki/q/math/epsilon"
	"github.com/itsubaki/q/math/matrix"
	"github.com/itsubaki/q/math/number"
	"github.com/itsubaki/q/quantum/gate"
	"github.com/itsubaki/q/quantum/qubit"
)

// Matrix is a density matrix.
type Matrix struct {
	rho *matrix.Matrix
}

// New returns a new density matrix.
func New(ensemble []State) *Matrix {
	n := ensemble[0].Qubit.Dim()

	rho := matrix.Zero(n, n)
	for _, s := range Normalize(ensemble) {
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

func IsValid(ensemble []State) bool {
	if len(ensemble) == 0 {
		return false
	}

	n := ensemble[0].Qubit.Dim()
	for _, s := range ensemble {
		if s.Qubit.Dim() != n {
			return false
		}
	}

	for _, s := range ensemble {
		if s.Probability < 0 {
			return false
		}
	}

	return true
}

// At returns a value of matrix at (i,j).
func (m *Matrix) At(i, j int) complex128 {
	return m.rho.At(i, j)
}

// Matrix returns the internal matrix.
func (m *Matrix) Matrix() *matrix.Matrix {
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
func (m *Matrix) IsPure(tol ...float64) bool {
	return epsilon.IsZeroF64(1-m.Purity(), tol...)
}

// IsMixed returns true if the density matrix is mixed.
func (m *Matrix) IsMixed(tol ...float64) bool {
	return !m.IsPure(tol...)
}

// IsHermite returns true if the density matrix is Hermitian.
func (m *Matrix) IsHermite(tol ...float64) bool {
	return m.rho.IsHermite(tol...)
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
func (m *Matrix) Project(q *qubit.Qubit, tol ...float64) (float64, *Matrix) {
	p := m.Probability(q)
	if epsilon.IsZeroF64(p, tol...) {
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
func (m *Matrix) PartialTrace(qb ...int) *Matrix {
	n := m.NumQubits()

	// mask for the qubits to be traced out
	var mask int
	for _, q := range qb {
		mask |= 1 << (n - 1 - q)
	}

	p, q := m.Dim()
	d := 1 << (n - len(qb))
	rho := matrix.Zero(d, d)
	for i := range p {
		ti, ki := split(i, n, mask)

		for j := range q {
			tj, kj := split(j, n, mask)

			if ti != tj {
				continue
			}

			rho.AddAt(ki, kj, m.At(i, j))
		}
	}

	return &Matrix{
		rho: rho,
	}
}

// Depolarizing returns the depolarizing channel.
// It applies the identity with probability (1 - p),
// and applies each of the Pauli gates X, Y, and Z with probability p/3.
func (m *Matrix) Depolarizing(p float64, qb int) *Matrix {
	n := m.NumQubits()

	idx := []int{qb}
	id := m.rho.Mul(complex(1-p, 0))
	xg := gate.TensorProduct(gate.X(), n, idx)
	yg := gate.TensorProduct(gate.Y(), n, idx)
	zg := gate.TensorProduct(gate.Z(), n, idx)

	x := matrix.MatMul(xg, m.rho, xg.Dagger()).Mul(complex(p/3, 0))
	y := matrix.MatMul(yg, m.rho, yg.Dagger()).Mul(complex(p/3, 0))
	z := matrix.MatMul(zg, m.rho, zg.Dagger()).Mul(complex(p/3, 0))

	return &Matrix{
		rho: id.Add(x).Add(y).Add(z),
	}
}

// ApplyChannel applies a channel to the density matrix.
// It applies the identity with probability 1-p, and applies the gate g with probability p.
func (m *Matrix) ApplyChannel(p float64, u *matrix.Matrix, qb ...int) *Matrix {
	n := m.NumQubits()

	e0 := gate.I().Mul(complex(math.Sqrt(1-p), 0))
	e1 := u.Mul(complex(math.Sqrt(p), 0))

	k0 := gate.TensorProduct(e0, n, qb)
	k1 := gate.TensorProduct(e1, n, qb)

	rho := matrix.ZeroLike(m.rho)
	rho = rho.Add(matrix.MatMul(k0, m.rho, k0.Dagger()))
	rho = rho.Add(matrix.MatMul(k1, m.rho, k1.Dagger()))

	return &Matrix{
		rho: rho,
	}
}

// BitFlip applies a bit flip channel to the density matrix.
func (m *Matrix) BitFlip(p float64, qb int) *Matrix {
	return m.ApplyChannel(p, gate.X(), qb)
}

// BitPhaseFlip applies a bit-phase flip channel to the density matrix.
func (m *Matrix) BitPhaseFlip(p float64, qb int) *Matrix {
	return m.ApplyChannel(p, gate.Y(), qb)
}

// PhaseFlip applies a phase flip channel to the density matrix.
func (m *Matrix) PhaseFlip(p float64, qb int) *Matrix {
	return m.ApplyChannel(p, gate.Z(), qb)
}

// split separates the bits of x into two integers according to mask.
//
// Bits where mask has value 1 are extracted into the returned trace value,
// preserving their relative order. Bits where mask has value 0 are extracted
// into the returned kept value.
//
// The n parameter specifies the number of bits of x to consider.
//
// For example:
//
//	n = 3
//	x = 0b101
//	mask = 0b010
//
// The bit at position 1 is traced out, so:
//
//	trace = 0b0
//	kept  = 0b11
//
// This helper is used when computing partial traces of density matrices.
func split(x, n, mask int) (int, int) {
	var trace, kept, trPos, kpPos int
	for i := range n {
		bit := (x >> i) & 1

		if (mask>>i)&1 == 1 {
			trace |= bit << trPos
			trPos++
			continue
		}

		kept |= bit << kpPos
		kpPos++
	}

	return trace, kept
}
