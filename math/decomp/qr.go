package decomp

import (
	"math"
	"math/cmplx"

	"github.com/itsubaki/q/math/epsilon"
	"github.com/itsubaki/q/math/matrix"
)

// QR performs the QR decomposition of matrix a using the Modified Gram-Schmidt process.
// It returns Q (orthonormal columns) and R (upper triangular) such that A = Q * R.
func QR(a *matrix.Matrix, eps ...float64) (q *matrix.Matrix, r *matrix.Matrix) {
	rows, cols := a.Dimension()
	q, r = a.Clone(), matrix.Zero(cols, cols)
	e := epsilon.E13(eps...)

	for k := range cols {
		nqk := norm(column(q, k))
		if nqk < e {
			// If the norm is smaller than a small threshold (effectively zero),
			// treat the k-th vector as numerically zero to avoid division by zero
			// and preserve numerical stability.
			zero(q, k)
			continue
		}

		r.Set(k, k, complex(nqk, 0))

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
func column(a *matrix.Matrix, j int) []complex128 {
	col := make([]complex128, a.Rows)
	for i := range a.Rows {
		col[i] = a.At(i, j)
	}

	return col
}

// div divides the j-th column of matrix a by v.
func div(a *matrix.Matrix, j int, v complex128) {
	for i := range a.Rows {
		a.DivAt(i, j, v)
	}
}

// zero sets the j-th column of matrix a to zero.
func zero(a *matrix.Matrix, j int) {
	for i := range a.Rows {
		a.Set(i, j, 0)
	}
}

// QRHH performs the QR decomposition using the Householder transformation.
// It returns Q (orthonormal columns) and R (upper triangular) such that A = Q * R.
func QRHH(a *matrix.Matrix, eps ...float64) (q *matrix.Matrix, r *matrix.Matrix) {
	rows, cols := a.Dimension()
	q, r = matrix.Identity(rows), a.Clone()

	for k := 0; k < cols && k < rows-1; k++ {
		x := make([]complex128, rows-k)
		for i := k; i < rows; i++ {
			x[i-k] = r.At(i, k)
		}

		u, ok := householder(x, eps...)
		if !ok {
			continue
		}

		// Apply the Householder transformation
		for j := k; j < cols; j++ {
			var v complex128
			for i := range u {
				v += cmplx.Conj(u[i]) * r.At(i+k, j)
			}

			for i := range u {
				r.SubAt(i+k, j, 2*v*u[i])
			}
		}

		for i := range rows {
			var v complex128
			for j := range u {
				v += q.At(i, j+k) * u[j]
			}

			for j := range u {
				q.SubAt(i, j+k, 2*v*cmplx.Conj(u[j]))
			}
		}
	}

	return q, r
}

func householder(x []complex128, eps ...float64) ([]complex128, bool) {
	e := epsilon.E13(eps...)

	alpha := complex(norm(x), 0)
	if cmplx.Abs(x[0]) > e {
		alpha *= -x[0] / complex(cmplx.Abs(x[0]), 0)
	}

	u := make([]complex128, len(x))
	u[0] = x[0] - alpha
	copy(u[1:], x[1:])

	nu := norm(u)
	if nu < e {
		return nil, false
	}

	for i := range u {
		u[i] /= complex(nu, 0)
	}

	return u, true
}
