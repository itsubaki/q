package observable

import (
	"github.com/itsubaki/q/math/matrix"
	"github.com/itsubaki/q/quantum/qubit"
)

// Projector returns a projector observable for the given qubit.
func Projector(qb ...*qubit.Qubit) *matrix.Matrix {
	q := qubit.TensorProduct(qb...)
	return q.OuterProduct(q)
}
