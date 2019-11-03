package density

import (
	"github.com/itsubaki/q/pkg/math/matrix"
	"github.com/itsubaki/q/pkg/quantum/qubit"
)

type Matrix struct {
	internal matrix.Matrix
}

func New() *Matrix {
	return &Matrix{make(matrix.Matrix, 0)}
}

func (m *Matrix) Add(p float64, q *qubit.Qubit) *Matrix {
	dim := q.Dimension()
	if len(m.internal) != dim {
		m.internal = make(matrix.Matrix, dim)
		for i := 0; i < dim; i++ {
			m.internal[i] = []complex128{}
			for j := 0; j < dim; j++ {
				m.internal[i] = append(m.internal[i], complex(0, 0))
			}
		}
	}

	tmp := q.OuterProduct(q).Mul(complex(p, 0))
	for i := 0; i < dim; i++ {
		for j := 0; j < dim; j++ {
			m.internal[i][j] = m.internal[i][j] + tmp[i][j]
		}
	}

	return m
}

func (m *Matrix) Apply(u matrix.Matrix) *Matrix {
	return &Matrix{u.Dagger().Apply(m.internal).Apply(u)}
}

func (m *Matrix) At(i, j int) complex128 {
	return m.internal[i][j]
}

func (m *Matrix) Trace() complex128 {
	return m.internal.Trace()
}
