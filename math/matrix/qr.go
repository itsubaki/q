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

// QRHH performs the QR decomposition using the Householder transformation.
// It returns Q (orthonormal columns) and R (upper triangular) such that A = Q * R.
func QRHH(a *Matrix, eps ...float64) (q *Matrix, r *Matrix) {
	rows, cols := a.Dimension()
	q, r = Identity(rows), a.Clone()

	for k := 0; k < cols && k < rows-1; k++ {
		x := make([]complex128, rows-k)
		for i := k; i < rows; i++ {
			x[i-k] = r.At(i, k)
		}

		// alpha defines the magnitude and direction to reflect x
		// onto the first basis vector e_1 in the Householder transform.
		alpha := complex(norm(x), 0)
		if cmplx.Abs(x[0]) != 0 {
			alpha *= x[0] / complex(cmplx.Abs(x[0]), 0)
		}

		// u = x + alpha * e_1
		u := make([]complex128, len(x))
		for i := range x {
			if i == 0 {
				u[i] = x[i] + alpha
				continue
			}

			u[i] = x[i]
		}

		nu := norm(u)
		if nu < epsilon.E13(eps...) {
			// If the norm is less than the threshold, skip this column
			continue
		}

		// Normalize u
		for i := range u {
			u[i] /= complex(nu, 0)
		}

		// h = I - 2 * uu^dagger
		h, uu := Identity(rows), outer(u)
		for i := k; i < rows; i++ {
			for j := k; j < rows; j++ {
				// h[i][j] -= 2 * uu.At(i-k, j-k)
				h.SubAt(i, j, 2*uu.At(i-k, j-k))
			}
		}

		// Apply the Householder transformation
		q = q.MatMul(h)
		r = h.MatMul(r)
	}

	return q, r
}

// outer computes the outer product of vector u with itself.
func outer(u []complex128) *Matrix {
	n := len(u)
	out := Zero(n, n)
	for i := range n {
		for j := range n {
			out.Set(i, j, u[i]*cmplx.Conj(u[j]))
		}
	}

	return out
}
