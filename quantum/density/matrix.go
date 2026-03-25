package density

import (
	"errors"
	"iter"
	"math"

	"github.com/itsubaki/q/math/epsilon"
	"github.com/itsubaki/q/math/matrix"
	"github.com/itsubaki/q/math/number"
	"github.com/itsubaki/q/quantum/gate"
	"github.com/itsubaki/q/quantum/qubit"
)

var ErrInvalidStates = errors.New("invalid states")

// Matrix is a density matrix.
type Matrix struct {
	rho *matrix.Matrix
}

// New returns a new density matrix.
func New(states []State) (*Matrix, error) {
	if !IsValid(states) {
		return nil, ErrInvalidStates
	}

	n := states[0].Qubit.Dim()
	rho := matrix.Zero(n, n)
	for _, s := range states {
		op := s.Qubit.OuterProduct(s.Qubit)
		rho = rho.Add(op.Mul(complex(s.Probability, 0)))
	}

	return &Matrix{
		rho: rho,
	}, nil
}

// NewPureState returns a new pure state density matrix for the given qubit.
func NewPureState(qb *qubit.Qubit) (*Matrix, error) {
	return New([]State{
		{
			Probability: 1.0,
			Qubit:       qb,
		},
	})
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
	return epsilon.IsOneF64(m.Purity(), tol...)
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

// Probability returns the probability of the qubit in the given state.
func (m *Matrix) Probability(q *qubit.Qubit) float64 {
	return real(matrix.MatMul(m.rho, q.OuterProduct(q)).Trace())
}

// ExpectedValue returns the expectation value of the given operator.
func (m *Matrix) ExpectedValue(u *matrix.Matrix) float64 {
	return real(matrix.MatMul(m.rho, u).Trace())
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

// ApplyKraus returns the density matrix after applying a set of Kraus operators.
func (m *Matrix) ApplyKraus(ops ...*matrix.Matrix) *Matrix {
	rho := matrix.ZeroLike(m.rho)
	for _, k := range ops {
		rho = rho.Add(matrix.MatMul(k, m.rho, k.Dagger()))
	}

	return &Matrix{
		rho: rho,
	}
}

// Apply returns the density matrix after applying a unitary operator.
func (m *Matrix) Apply(u *matrix.Matrix) *Matrix {
	return m.ApplyKraus(u)
}

// AmplitudeDamping returns the density matrix after applying an amplitude damping channel to the specified qubit.
func (m *Matrix) AmplitudeDamping(gamma float64, qb int) *Matrix {
	e0 := matrix.New(
		[]complex128{1, 0},
		[]complex128{0, complex(math.Sqrt(1-gamma), 0)},
	)

	e1 := matrix.New(
		[]complex128{0, complex(math.Sqrt(gamma), 0)},
		[]complex128{0, 0},
	)

	n := m.NumQubits()
	k0 := gate.TensorProduct(e0, n, []int{qb})
	k1 := gate.TensorProduct(e1, n, []int{qb})
	return m.ApplyKraus(k0, k1)
}

// PhaseDamping returns the density matrix after applying a phase damping channel to the specified qubit.
func (m *Matrix) PhaseDamping(gamma float64, qb int) *Matrix {
	e0 := matrix.New(
		[]complex128{1, 0},
		[]complex128{0, complex(math.Sqrt(1-gamma), 0)},
	)

	e1 := matrix.New(
		[]complex128{0, 0},
		[]complex128{0, complex(math.Sqrt(gamma), 0)},
	)

	n := m.NumQubits()
	k0 := gate.TensorProduct(e0, n, []int{qb})
	k1 := gate.TensorProduct(e1, n, []int{qb})
	return m.ApplyKraus(k0, k1)
}

// PauliChannel returns the density matrix after applying a Pauli channel to the specified qubit.
func (m *Matrix) PauliChannel(px, py, pz float64, qb int) *Matrix {
	e0 := gate.I().Mul(complex(math.Sqrt(1-px-py-pz), 0))
	e1 := gate.X().Mul(complex(math.Sqrt(px), 0))
	e2 := gate.Y().Mul(complex(math.Sqrt(py), 0))
	e3 := gate.Z().Mul(complex(math.Sqrt(pz), 0))

	n := m.NumQubits()
	k0 := gate.TensorProduct(e0, n, []int{qb})
	k1 := gate.TensorProduct(e1, n, []int{qb})
	k2 := gate.TensorProduct(e2, n, []int{qb})
	k3 := gate.TensorProduct(e3, n, []int{qb})
	return m.ApplyKraus(k0, k1, k2, k3)
}

// Depolarizing returns the density matrix after applying a depolarizing channel to the specified qubit.
// It applies the identity with probability 1-p, and applies X, Y, Z with probability p/3 each.
func (m *Matrix) Depolarizing(p float64, qb int) *Matrix {
	return m.PauliChannel(p/3, p/3, p/3, qb)
}

// FlipChannel returns the density matrix after applying a flip channel to the specified qubit.
// It applies the identity with probability 1-p, and applies the gate g with probability p.
func (m *Matrix) FlipChannel(p float64, u *matrix.Matrix, qb ...int) *Matrix {
	e0 := gate.I().Mul(complex(math.Sqrt(1-p), 0))
	e1 := u.Mul(complex(math.Sqrt(p), 0))

	n := m.NumQubits()
	k0 := gate.TensorProduct(e0, n, qb)
	k1 := gate.TensorProduct(e1, n, qb)
	return m.ApplyKraus(k0, k1)
}

// BitFlip returns the density matrix after applying a bit flip channel to the specified qubit.
func (m *Matrix) BitFlip(p float64, qb int) *Matrix {
	return m.FlipChannel(p, gate.X(), qb)
}

// BitPhaseFlip returns the density matrix after applying a bit-phase flip channel to the specified qubit.
func (m *Matrix) BitPhaseFlip(p float64, qb int) *Matrix {
	return m.FlipChannel(p, gate.Y(), qb)
}

// PhaseFlip returns the density matrix after applying a phase flip channel to the specified qubit.
func (m *Matrix) PhaseFlip(p float64, qb int) *Matrix {
	return m.FlipChannel(p, gate.Z(), qb)
}

// PartialTrace returns the density matrix obtained by tracing out the specified qubits from the original density matrix.
// The length of index must be less than or equal to n - 1, where n is the number of qubits in the matrix.
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

// TraceOut is an alias for PartialTrace.
func (m *Matrix) TraceOut(qb ...int) *Matrix {
	return m.PartialTrace(qb...)
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
