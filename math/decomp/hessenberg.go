package decomp

import (
	"math/cmplx"

	"github.com/itsubaki/q/math/epsilon"
	"github.com/itsubaki/q/math/matrix"
)

// Hessenberg performs the Hessenberg decomposition of matrix a using Householder transformations.
// It returns Q (unitary) and H (upper Hessenberg) such that A = Q * H * Q^dagger.
func Hessenberg(a *matrix.Matrix) (q *matrix.Matrix, h *matrix.Matrix) {
	n := a.Rows
	q, h = matrix.Identity(n), a.Clone()

	for k := range n - 2 {
		x := make([]complex128, n-k-1)
		for i := range x {
			x[i] = h.At(k+1+i, k)
		}

		u, ok := householder(x)
		if !ok {
			continue
		}

		for j := k; j < n; j++ {
			var proj complex128
			for i := range u {
				proj += cmplx.Conj(u[i]) * h.At(k+1+i, j)
			}

			for i := range u {
				h.SubAt(k+1+i, j, 2*u[i]*proj)
			}
		}

		for i := range n {
			var proj complex128
			for j := range u {
				proj += h.At(i, k+1+j) * u[j]
			}

			for j := range u {
				h.SubAt(i, k+1+j, 2*proj*cmplx.Conj(u[j]))
			}
		}

		for i := range n {
			var proj complex128
			for j := range u {
				proj += q.At(i, k+1+j) * u[j]
			}

			for j := range u {
				q.SubAt(i, k+1+j, 2*proj*cmplx.Conj(u[j]))
			}
		}
	}

	return q, h
}

// IsHessenberg returns true if m is Hessenberg matrix.
func IsHessenberg(m *matrix.Matrix, eps ...float64) bool {
	if !m.IsSquare() {
		return false
	}

	e := epsilon.E13(eps...)
	for i := 2; i < m.Rows; i++ {
		for j := range i - 1 {
			if cmplx.Abs(m.At(i, j)) > e {
				return false
			}
		}
	}

	return true
}
