package matrix

import (
	"math"
	"math/cmplx"

	"github.com/itsubaki/q/math/epsilon"
)

// QR performs the QR decomposition of matrix a using the Gram-Schmidt process.
// It returns Q (orthonormal columns) and R (upper triangular) such that A = Q * R.
func QR(a *Matrix, eps ...float64) (q *Matrix, r *Matrix) {
	rows, cols := a.Dimension()
	q, r = Zero(rows, cols), Zero(cols, cols)

	for j := range cols {
		// Initialize the j-th column of Q with the j-th column of A
		for i := range rows {
			q.Set(i, j, a.At(i, j))
		}

		for k := range j {
			// Compute the inner product of the k-th column of Q and the j-th column of A
			var inprod complex128
			for i := range rows {
				inprod += cmplx.Conj(q.At(i, k)) * a.At(i, j)
			}

			r.Set(k, j, inprod)
			for i := range rows {
				q.Set(i, j, q.At(i, j)-inprod*q.At(i, k))
			}
		}

		// Compute the norm of the j-th column of Q
		var norm2 float64
		for i := range rows {
			norm2 += cmplx.Abs(q.At(i, j)) * cmplx.Abs(q.At(i, j))
		}
		norm := math.Sqrt(norm2)

		// Set the j-th diagonal element of R to the norm of the j-th column of Q
		r.Set(j, j, complex(norm, 0))

		// If the norm is less than the threshold, set the j-th column of Q to zero
		if norm < epsilon.E13(eps...) {
			for i := range rows {
				q.Set(i, j, 0)
			}

			continue
		}

		// Normalize the j-th column of Q
		for i := range rows {
			q.Set(i, j, q.At(i, j)/complex(norm, 0))
		}
	}

	return q, r
}
