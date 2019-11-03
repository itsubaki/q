package density

import (
	"fmt"

	"github.com/itsubaki/q/pkg/math/matrix"
	"github.com/itsubaki/q/pkg/quantum/qubit"
)

type Matrix struct {
	internal matrix.Matrix
}

func New(v ...[]complex128) *Matrix {
	return &Matrix{matrix.New(v...)}
}

func (m *Matrix) init(dim int) {
	m.internal = make(matrix.Matrix, dim)
	for i := 0; i < dim; i++ {
		m.internal[i] = []complex128{}
		for j := 0; j < dim; j++ {
			m.internal[i] = append(m.internal[i], complex(0, 0))
		}
	}
}

func (m *Matrix) Add(p float64, q *qubit.Qubit) *Matrix {
	dim := q.Dimension()
	if len(m.internal) < 1 {
		m.init(dim)
	}

	if len(m.internal) != dim {
		panic(fmt.Sprintf("invalid dimension=%d", dim))
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
	return &Matrix{u.Dagger().Apply(m.internal).Apply(u)}
}

func (m *Matrix) Measure(q *qubit.Qubit) complex128 {
	op := q.OuterProduct(q)
	return op.Dagger().Apply(m.internal).Apply(op).Trace()
}

func (m *Matrix) Trace() complex128 {
	return m.internal.Trace()
}
