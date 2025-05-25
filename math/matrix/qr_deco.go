package matrix

import (
	"math"
	"math/cmplx"

	"github.com/itsubaki/q/math/epsilon"
)

// QR performs the QR decomposition of matrix a using the Modified Gram-Schmidt process.
// It returns Q (orthonormal columns) and R (upper triangular) such that A = Q * R.
func QR(a *Matrix, eps ...float64) (q *Matrix, r *Matrix) {
	rows, cols := a.Dimension()
	q, r = a.Clone(), Zero(cols, cols)

	for k := range cols {
		nqk := norm(column(q, k))
		r.Set(k, k, complex(nqk, 0))

		if nqk < epsilon.E13(eps...) {
			// If the norm is smaller than a small threshold (effectively zero),
			// treat the k-th vector as numerically zero to avoid division by zero
			// and preserve numerical stability.
			zero(q, k)
			continue
		}

		// Normalize the k-th column of Q
		div(q, k, complex(nqk, 0))

		// Compute the k-th row of R
		for j := k + 1; j < cols; j++ {
			rkj := dot(column(q, k), column(q, j))
			r.Set(k, j, rkj)

			// Subtract the projection of the j-th column onto the k-th column of Q
			for i := range rows {
				// qij = qij - rkj * qik
				q.SubAt(i, j, rkj*q.At(i, k))
			}
		}
	}

	return q, r
}

// norm computes the Euclidean norm of a vector x.
func norm(x []complex128) float64 {
	var sum float64
	for _, v := range x {
		sum += cmplx.Abs(v) * cmplx.Abs(v)
	}

	return math.Sqrt(sum)
}

// dot computes the complex inner product of vectors a and b.
func dot(a, b []complex128) complex128 {
	var sum complex128
	for i := range a {
		sum += cmplx.Conj(a[i]) * b[i]
	}

	return sum
}

// column returns the j-th column of matrix a as a slice.
func column(a *Matrix, j int) []complex128 {
	col := make([]complex128, a.Rows)
	for i := range a.Rows {
		col[i] = a.At(i, j)
	}

	return col
}

// div divides the j-th column of matrix a by v.
func div(a *Matrix, j int, v complex128) {
	for i := range a.Rows {
		a.DivAt(i, j, v)
	}
}

// zero sets the j-th column of matrix a to zero.
func zero(a *Matrix, j int) {
	for i := range a.Rows {
		a.Set(i, j, 0)
	}
}
