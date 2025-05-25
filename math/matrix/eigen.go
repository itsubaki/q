package matrix

import (
	"math/cmplx"

	"github.com/itsubaki/q/math/epsilon"
)

// EigenJacobi performs eigen decomposition of a matrix using the Jacobi method.
func EigenJacobi(a *Matrix, iter int, eps ...float64) (vectors *Matrix, lambdas *Matrix) {
	n := a.Rows
	v, ak := Identity(n), a.Clone()

	for range iter {
		for i := range n - 1 {
			for j := i + 1; j < n; j++ {
				a, b, c := ak.At(i, i), ak.At(j, j), ak.At(i, j)
				if cmplx.Abs(c) < epsilon.E13(eps...) {
					continue
				}

				phi := 0.5 * cmplx.Atan(2*c/(b-a))
				cos := cmplx.Cos(phi)
				sin := cmplx.Sin(phi) * cmplx.Rect(1, cmplx.Phase(c))

				g := Identity(n)
				g.Set(i, i, cos)
				g.Set(j, j, cos)
				g.Set(i, j, -cmplx.Conj(sin))
				g.Set(j, i, sin)

				v = MatMul(v, g)
				ak = MatMul(g.Dagger(), ak, g)
			}
		}
	}

	d := ZeroLike(ak)
	for i := range n {
		d.Set(i, i, ak.At(i, i))
	}

	return v, d
}

// EigenQR performs eigen decomposition of a matrix using the Schur decomposition.
func EigenQR(m *Matrix, qr QRFunc, iter int, eps ...float64) (lambdas *Matrix, vectors *Matrix) {
	q, t := Schur(m, qr, iter, eps...)
	lambdas, vectors = EigenUpperT(t, eps...)
	return lambdas, MatMul(q, vectors)
}

// EigenUpperT performs eigen decomposition of an upper triangular matrix.
func EigenUpperT(t *Matrix, eps ...float64) (lambdas *Matrix, vectors *Matrix) {
	lambdas = ZeroLike(t)
	for i := range t.Rows {
		lambdas.Set(i, i, t.At(i, i))
	}

	vectors = Zero(t.Rows, t.Rows)
	for k := range t.Rows {
		x := make([]complex128, t.Rows)
		x[k] = 1.0

		lambdak := t.At(k, k)
		for i := k - 1; i >= 0; i-- {
			diff := t.At(i, i) - lambdak
			if cmplx.Abs(diff) < epsilon.E13(eps...) {
				x[i] = 0.0
				continue
			}

			var sum complex128
			for j := i + 1; j <= k; j++ {
				sum += t.At(i, j) * x[j]
			}

			x[i] = -sum / diff
		}

		nx := norm(x)
		if nx > epsilon.E13(eps...) {
			for i := range x {
				x[i] /= complex(nx, 0)
			}
		}

		for i := range t.Rows {
			vectors.Set(i, k, x[i])
		}
	}

	return lambdas, vectors
}
