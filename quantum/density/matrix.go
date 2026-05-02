package density

import (
	"iter"
	"math/cmplx"

	"github.com/itsubaki/q/math/eigen"
	"github.com/itsubaki/q/math/epsilon"
	"github.com/itsubaki/q/math/matrix"
	"github.com/itsubaki/q/math/number"
	"github.com/itsubaki/q/quantum/channel"
	"github.com/itsubaki/q/quantum/qubit"
)

// DensityMatrix represents a quantum state in the density matrix formalism.
type DensityMatrix struct {
	rho *matrix.Matrix
}

// New returns a density matrix constructed from a pure state represented by a qubit.
func New(qb *qubit.Qubit) *DensityMatrix {
	return NewMixed([]WeightedState{
		{
			Probability: 1.0,
			Qubit:       qb,
		},
	})
}

// NewMixed returns a density matrix constructed from a set of states.
func NewMixed(states []WeightedState) *DensityMatrix {
	if len(states) == 0 {
		return nil
	}

	n := states[0].Qubit.Dim()
	rho := matrix.Zero(n, n)
	for _, s := range Normalize(states) {
		op := s.Qubit.OuterProduct(s.Qubit)
		rho = rho.Add(op.Mul(complex(s.Probability, 0)))
	}

	return &DensityMatrix{
		rho: rho,
	}
}

// At returns the value at (i, j).
func (m *DensityMatrix) At(i, j int) complex128 {
	return m.rho.At(i, j)
}

// Seq2 returns a sequence of rows.
func (m *DensityMatrix) Seq2() iter.Seq2[int, []complex128] {
	return m.rho.Seq2()
}

// Dim returns the dimension of the density matrix.
func (m *DensityMatrix) Dim() (rows int, cols int) {
	return m.rho.Dim()
}

// IsPure returns true if the density matrix is pure.
func (m *DensityMatrix) IsPure(tol ...float64) bool {
	return epsilon.IsOneF64(m.Purity(), tol...)
}

// IsMixed returns true if the density matrix is mixed.
func (m *DensityMatrix) IsMixed(tol ...float64) bool {
	return !m.IsPure(tol...)
}

// IsHermitian returns true if the density matrix is Hermitian.
func (m *DensityMatrix) IsHermitian(tol ...float64) bool {
	return m.rho.IsHermitian(tol...)
}

// NumQubits returns the number of qubits.
func (m *DensityMatrix) NumQubits() int {
	p, _ := m.Dim()
	return number.Log2(p)
}

// ExpectedValue returns the expectation value of the given operator.
func (m *DensityMatrix) ExpectedValue(u *matrix.Matrix) float64 {
	return real(matrix.MatMul(m.rho, u).Trace())
}

// Trace returns the trace of the density matrix.
func (m *DensityMatrix) Trace() float64 {
	return real(m.rho.Trace())
}

// Purity returns the purity of the density matrix, defined as Tr(rho^2).
func (m *DensityMatrix) Purity() float64 {
	return real(matrix.MatMul(m.rho, m.rho).Trace())
}

// TensorProduct returns the tensor product of two density matrices.
func (m *DensityMatrix) TensorProduct(n *DensityMatrix) *DensityMatrix {
	return &DensityMatrix{
		rho: m.rho.TensorProduct(n.rho),
	}
}

// Sqrt returns the square root of the density matrix.
func (m *DensityMatrix) Sqrt(tol ...float64) *DensityMatrix {
	v, d := eigen.Jacobi(m.rho, 100, tol...)
	d.Fdiag(func(v complex128) complex128 { return cmplx.Pow(v, 0.5) })
	return &DensityMatrix{
		rho: matrix.MatMul(v, d, v.Dagger()),
	}
}

// Clone returns a copy of m.
func (m *DensityMatrix) Clone() *DensityMatrix {
	return &DensityMatrix{
		rho: m.rho.Clone(),
	}
}

// Measure returns the probability and post-measurement density matrix.
func (m *DensityMatrix) Measure(q *qubit.Qubit, tol ...float64) (float64, *DensityMatrix) {
	p := real(matrix.MatMul(m.rho, q.OuterProduct(q)).Trace())
	if epsilon.IsZeroF64(p, tol...) {
		return 0, &DensityMatrix{
			rho: matrix.ZeroLike(m.rho),
		}
	}

	op := q.OuterProduct(q)
	rho := matrix.MatMul(op, m.rho, op)
	return p, &DensityMatrix{
		rho: rho.Mul(1.0 / complex(p, 0)),
	}
}

// PartialTrace returns the density matrix obtained by tracing out the specified qubits.
// The number of qubits to trace out must be less than or equal to n - 1, where n is the number of qubits in the matrix.
func (m *DensityMatrix) PartialTrace(qb ...int) *DensityMatrix {
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

	return &DensityMatrix{
		rho: rho,
	}
}

// TraceOut is an alias for PartialTrace.
func (m *DensityMatrix) TraceOut(qb ...int) *DensityMatrix {
	return m.PartialTrace(qb...)
}

// Pauli returns the density matrix after applying a Pauli channel to the specified qubits.
func (m *DensityMatrix) Pauli(px, py, pz float64, qb ...int) *DensityMatrix {
	return m.applyChannel(func(q int) channel.ChannelFunc {
		return channel.Pauli(px, py, pz, q)
	}, qb...)
}

// Depolarizing returns the density matrix after applying a depolarizing channel to the specified qubits.
func (m *DensityMatrix) Depolarizing(p float64, qb ...int) *DensityMatrix {
	return m.applyChannel(func(q int) channel.ChannelFunc {
		return channel.Depolarizing(p, q)
	}, qb...)
}

// AmplitudeDamping returns the density matrix after applying an amplitude damping channel to the specified qubits.
func (m *DensityMatrix) AmplitudeDamping(gamma float64, qb ...int) *DensityMatrix {
	return m.applyChannel(func(q int) channel.ChannelFunc {
		return channel.AmplitudeDamping(gamma, q)
	}, qb...)
}

// PhaseDamping returns the density matrix after applying a phase damping channel to the specified qubits.
func (m *DensityMatrix) PhaseDamping(gamma float64, qb ...int) *DensityMatrix {
	return m.applyChannel(func(q int) channel.ChannelFunc {
		return channel.PhaseDamping(gamma, q)
	}, qb...)
}

// Flip returns the density matrix after applying a flip channel to the specified qubits.
func (m *DensityMatrix) Flip(p float64, u *matrix.Matrix, qb ...int) *DensityMatrix {
	return m.applyChannel(func(q int) channel.ChannelFunc {
		return channel.Flip(p, u, q)
	}, qb...)
}

// BitFlip returns the density matrix after applying a bit flip channel to the specified qubits.
func (m *DensityMatrix) BitFlip(p float64, qb ...int) *DensityMatrix {
	return m.applyChannel(func(q int) channel.ChannelFunc {
		return channel.BitFlip(p, q)
	}, qb...)
}

// PhaseFlip returns the density matrix after applying a phase flip channel to the specified qubits.
func (m *DensityMatrix) PhaseFlip(p float64, qb ...int) *DensityMatrix {
	return m.applyChannel(func(q int) channel.ChannelFunc {
		return channel.PhaseFlip(p, q)
	}, qb...)
}

// BitPhaseFlip returns the density matrix after applying a bit-phase flip channel to the specified qubits.
func (m *DensityMatrix) BitPhaseFlip(p float64, qb ...int) *DensityMatrix {
	return m.applyChannel(func(q int) channel.ChannelFunc {
		return channel.BitPhaseFlip(p, q)
	}, qb...)
}

// applyChannel is a helper function that applies a quantum channel to the specified qubits.
// If no qubits are specified, the channel is applied to all qubits in the density matrix.
func (m *DensityMatrix) applyChannel(f func(int) channel.ChannelFunc, qb ...int) *DensityMatrix {
	if len(qb) == 0 {
		qb = make([]int, m.NumQubits())
		for i := range qb {
			qb[i] = i
		}
	}

	fn := make([]channel.ChannelFunc, len(qb))
	for i, q := range qb {
		fn[i] = f(q)
	}

	return m.ApplyChannelFunc(fn...)
}

// ApplyChannelFunc returns the density matrix after applying a quantum channel.
func (m *DensityMatrix) ApplyChannelFunc(fn ...channel.ChannelFunc) *DensityMatrix {
	ch := make([]*channel.Channel, len(fn))
	for i, f := range fn {
		ch[i] = f(m.NumQubits())
	}

	return m.ApplyChannel(ch...)
}

// ApplyChannel returns the density matrix after applying a quantum channel.
func (m *DensityMatrix) ApplyChannel(ch ...*channel.Channel) *DensityMatrix {
	if len(ch) == 0 {
		return m.Clone()
	}

	out := m.ApplyKraus(ch[0].Kraus...)
	for _, ch := range ch[1:] {
		out = out.ApplyKraus(ch.Kraus...)
	}

	return out
}

// Apply returns the density matrix after applying a unitary operator.
func (m *DensityMatrix) Apply(u *matrix.Matrix) *DensityMatrix {
	return m.ApplyKraus(u)
}

// ApplyKraus returns the density matrix after applying a set of Kraus operators.
func (m *DensityMatrix) ApplyKraus(ops ...*matrix.Matrix) *DensityMatrix {
	if len(ops) == 0 {
		return m.Clone()
	}

	rho := matrix.ZeroLike(m.rho)
	for _, k := range ops {
		rho = rho.Add(matrix.MatMul(k, m.rho, k.Dagger()))
	}

	return &DensityMatrix{
		rho: rho,
	}
}

// MatMul returns the matrix product of two density matrices.
func MatMul(m, n *DensityMatrix) *DensityMatrix {
	return &DensityMatrix{
		rho: matrix.MatMul(m.rho, n.rho),
	}
}

// Equal returns true if two density matrices are equal within a specified tolerance.
func Equal(m, n *DensityMatrix, tol ...float64) bool {
	return m.rho.Equal(m.rho)
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
