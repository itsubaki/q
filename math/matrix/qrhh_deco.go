package matrix

import (
	"math"
	"math/cmplx"

	"github.com/itsubaki/q/math/epsilon"
)

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

		// Normalize u
		norm := norm(u)
		if norm < epsilon.E13(eps...) {
			// If the norm is less than the threshold, skip this column
			continue
		}

		for i := range u {
			u[i] /= complex(norm, 0)
		}

		// hk = I - 2 * uu^dagger
		hk := Identity(rows)
		uu := outer(u)
		for i := k; i < rows; i++ {
			for j := k; j < rows; j++ {
				hk.Set(i, j, hk.At(i, j)-2*uu.At(i-k, j-k))
			}
		}

		// Apply the Householder transformation to R
		r = hk.MatMul(r)
		q = q.MatMul(hk)
	}

	return q, r
}

func norm(x []complex128) float64 {
	var sum float64
	for _, v := range x {
		sum += cmplx.Abs(v) * cmplx.Abs(v)
	}

	return math.Sqrt(sum)
}

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
