package decomp

import (
	"math/cmplx"

	"github.com/itsubaki/q/math/epsilon"
	"github.com/itsubaki/q/math/matrix"
)

// QRFunc is a function type that performs QR decomposition on a given matrix.
// It returns Q (orthonormal columns) and R (upper triangular) such that A = Q * R.
type QRFunc func(a *matrix.Matrix, eps ...float64) (q *matrix.Matrix, r *matrix.Matrix)

var (
	_ QRFunc = QR
	_ QRFunc = QRHH
)

// Schur performs the Schur decomposition of matrix a using iterative QR decomposition.
// It returns Q (unitary) and T (upper triangular) such that A = Q * T * Q^dagger.
func Schur(a *matrix.Matrix, qr QRFunc, iter int, eps ...float64) (q *matrix.Matrix, t *matrix.Matrix) {
	q, t = matrix.Identity(a.Rows), a.Clone()

	for range iter {
		qk, rk := qr(t, eps...)
		t = rk.MatMul(qk)
		q = q.MatMul(qk)

		if IsUpperTriangular(t, eps...) {
			break
		}
	}

	return q, t
}

// IsUpperTriangular returns true if m is upper triangular matrix.
func IsUpperTriangular(m *matrix.Matrix, eps ...float64) bool {
	if !m.IsSquare() {
		return false
	}

	e := epsilon.E13(eps...)
	for i := 1; i < m.Rows; i++ {
		for j := range i {
			if cmplx.Abs(m.At(i, j)) > e {
				return false
			}
		}
	}

	return true
}
