package density

import (
	"iter"

	"github.com/itsubaki/q/math/epsilon"
	"github.com/itsubaki/q/math/matrix"
	"github.com/itsubaki/q/math/number"
	"github.com/itsubaki/q/quantum/qubit"
)

// DensityMatrix represents a quantum state in the density matrix formalism.
type DensityMatrix struct {
	rho *matrix.Matrix
}

// FromStates returns a density matrix constructed from a set of states.
func FromStates(states []WeightedState) *DensityMatrix {
	if len(states) == 0 {
		return &DensityMatrix{
			rho: matrix.New(),
		}
	}

	n := states[0].Qubit.Dim()
	rho := matrix.Zero(n, n)
	for _, s := range states {
		op := s.Qubit.OuterProduct(s.Qubit)
		rho = rho.Add(op.Mul(complex(s.Probability, 0)))
	}

	return &DensityMatrix{
		rho: rho,
	}
}

// From returns a density matrix constructed from a pure state represented by a qubit.
func From(qb *qubit.Qubit) *DensityMatrix {
	return FromStates([]WeightedState{
		{
			Probability: 1.0,
			Qubit:       qb,
		},
	})
}

// At returns the value at (i, j).
func (m *DensityMatrix) At(i, j int) complex128 {
	return m.rho.At(i, j)
}

// Matrix returns the internal matrix.
func (m *DensityMatrix) Matrix() *matrix.Matrix {
	return m.rho
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

// Probability returns the probability of the qubit in the given state.
func (m *DensityMatrix) Probability(q *qubit.Qubit) float64 {
	return real(matrix.MatMul(m.rho, q.OuterProduct(q)).Trace())
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

// Clone returns a copy of m.
func (m *DensityMatrix) Clone() *DensityMatrix {
	return &DensityMatrix{
		rho: m.rho.Clone(),
	}
}

// Project returns the probability and post-measurement density matrix.
func (m *DensityMatrix) Project(q *qubit.Qubit, tol ...float64) (float64, *DensityMatrix) {
	p := m.Probability(q)
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

// Pauli returns the density matrix after applying a Pauli channel to the specified qubit.
func (m *DensityMatrix) Pauli(px, py, pz float64, qb int) *DensityMatrix {
	return m.ApplyChannelFunc(Pauli(px, py, pz, qb))
}

// Depolarizing returns the density matrix after applying a depolarizing channel to the specified qubit.
func (m *DensityMatrix) Depolarizing(p float64, qb int) *DensityMatrix {
	return m.ApplyChannelFunc(Depolarizing(p, qb))
}

// AmplitudeDamping returns the density matrix after applying an amplitude damping channel to the specified qubit.
func (m *DensityMatrix) AmplitudeDamping(gamma float64, qb int) *DensityMatrix {
	return m.ApplyChannelFunc(AmplitudeDamping(gamma, qb))
}

// PhaseDamping returns the density matrix after applying a phase damping channel to the specified qubit.
func (m *DensityMatrix) PhaseDamping(gamma float64, qb int) *DensityMatrix {
	return m.ApplyChannelFunc(PhaseDamping(gamma, qb))
}

// Flip returns the density matrix after applying a flip channel to the specified qubit.
func (m *DensityMatrix) Flip(p float64, u *matrix.Matrix, qb int) *DensityMatrix {
	return m.ApplyChannelFunc(Flip(p, u, qb))
}

// BitFlip returns the density matrix after applying a bit flip channel to the specified qubit.
func (m *DensityMatrix) BitFlip(p float64, qb int) *DensityMatrix {
	return m.ApplyChannelFunc(BitFlip(p, qb))
}

// PhaseFlip returns the density matrix after applying a phase flip channel to the specified qubit.
func (m *DensityMatrix) PhaseFlip(p float64, qb int) *DensityMatrix {
	return m.ApplyChannelFunc(PhaseFlip(p, qb))
}

// BitPhaseFlip returns the density matrix after applying a bit-phase flip channel to the specified qubit.
func (m *DensityMatrix) BitPhaseFlip(p float64, qb int) *DensityMatrix {
	return m.ApplyChannelFunc(BitPhaseFlip(p, qb))
}

// ApplyChannelFunc returns the density matrix after applying a quantum channel.
func (m *DensityMatrix) ApplyChannelFunc(channel ...ChannelFunc) *DensityMatrix {
	ch := make([]*Channel, len(channel))
	for i, f := range channel {
		ch[i] = f(m.NumQubits())
	}

	return m.ApplyChannel(ch...)
}

// ApplyChannel returns the density matrix after applying a quantum channel.
func (m *DensityMatrix) ApplyChannel(channel ...*Channel) *DensityMatrix {
	if len(channel) == 0 {
		return m.Clone()
	}

	out := m.ApplyKraus(channel[0].Kraus...)
	for _, ch := range channel[1:] {
		out = out.ApplyKraus(ch.Kraus...)
	}

	return out
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

// Apply returns the density matrix after applying a unitary operator.
func (m *DensityMatrix) Apply(u *matrix.Matrix) *DensityMatrix {
	return &DensityMatrix{
		rho: matrix.MatMul(u, m.rho, u.Dagger()),
	}
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
