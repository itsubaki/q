package density

import (
	"fmt"
	"math"

	"github.com/itsubaki/q/pkg/math/matrix"
	"github.com/itsubaki/q/pkg/quantum/gate"
	"github.com/itsubaki/q/pkg/quantum/qubit"
)

type Matrix struct {
	internal matrix.Matrix
}

func New(v ...[]complex128) *Matrix {
	return &Matrix{matrix.New(v...)}
}

func (m *Matrix) Zero(dim int) {
	m.internal = make(matrix.Matrix, dim)
	for i := 0; i < dim; i++ {
		m.internal[i] = make([]complex128, 0)
		for j := 0; j < dim; j++ {
			m.internal[i] = append(m.internal[i], complex(0, 0))
		}
	}
}

func (m *Matrix) Add(p float64, q *qubit.Qubit) *Matrix {
	dim := q.Dimension()
	if len(m.internal) < 1 {
		m.Zero(dim)
	}

	if len(m.internal) != dim {
		panic(fmt.Sprintf("dimension invalid. m=%d n=%d", len(m.internal), dim))
	}

	op := q.OuterProduct(q).Mul(complex(p, 0))
	for i := 0; i < dim; i++ {
		for j := 0; j < dim; j++ {
			m.internal[i][j] = m.internal[i][j] + op[i][j]
		}
	}

	return m
}

func (m *Matrix) Apply(u matrix.Matrix) *Matrix {
	m.internal = u.Dagger().Apply(m.internal).Apply(u)
	return m
}

func (m *Matrix) Measure(q *qubit.Qubit) complex128 {
	return m.internal.Apply(q.OuterProduct(q)).Trace()
}

func (m *Matrix) ExpectedValue(u matrix.Matrix) complex128 {
	return m.internal.Apply(u).Trace()
}

func (m *Matrix) Trace() complex128 {
	return m.internal.Trace()
}

func (m *Matrix) PartialTrace() complex128 {
	return m.internal.PartialTrace()
}

func (m *Matrix) NumberOfBit() int {
	mm, _ := m.internal.Dimension()
	log := math.Log2(float64(mm))
	return int(log)
}

func (m *Matrix) Depolarizing(p float64) {
	n := m.NumberOfBit()
	i := gate.I(n).Mul(complex(p/2, 0))
	r := m.internal.Mul(complex(1-p, 0))

	m.internal = i.Add(r)
}

func Flip(p float64, m matrix.Matrix) (matrix.Matrix, matrix.Matrix) {
	n, _ := m.Dimension()
	e0 := gate.I(n).Mul(complex(math.Sqrt(p), 0))
	e1 := m.Mul(complex(math.Sqrt(1-p), 0))
	return e0, e1
}

func BitFlip(p float64) (matrix.Matrix, matrix.Matrix) {
	return Flip(p, gate.X())
}

func PhaseFlip(p float64) (matrix.Matrix, matrix.Matrix) {
	return Flip(p, gate.Z())
}

func BitPhaseFlip(p float64) (matrix.Matrix, matrix.Matrix) {
	return Flip(p, gate.Y())
}
