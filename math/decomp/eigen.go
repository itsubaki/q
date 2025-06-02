package decomp

import (
	"math"
	"math/cmplx"

	"github.com/itsubaki/q/math/epsilon"
	"github.com/itsubaki/q/math/matrix"
)

// EigenJacobi performs eigen decomposition of a matrix using the Jacobi method.
// The matrix `a` must be Hermitian.
func EigenJacobi(a *matrix.Matrix, iter int, eps ...float64) (lambdas *matrix.Matrix, vectors *matrix.Matrix) {
	n := a.Rows
	v, ak := matrix.Identity(n), a.Clone()
	e := epsilon.E13(eps...)

	for range iter {
		for i := range n - 1 {
			for j := i + 1; j < n; j++ {
				a, b, c := ak.At(i, i), ak.At(j, j), ak.At(i, j)
				if cmplx.Abs(c) < e {
					continue
				}

				diff, phi := b-a, complex(math.Pi/4, 0)
				if cmplx.Abs(diff) > e {
					phi = 0.5 * cmplx.Atan(2*c/diff)
				}

				cos := cmplx.Cos(phi)
				sin := cmplx.Sin(phi) * cmplx.Rect(1, cmplx.Phase(c)) // cmplx.Exp(complex(0, cmplx.Phase(c)))

				g := matrix.Identity(n)
				g.Set(i, i, cos)
				g.Set(i, j, -cmplx.Conj(sin))
				g.Set(j, i, sin)
				g.Set(j, j, cos)

				v = matrix.MatMul(v, g)
				ak = matrix.MatMul(g.Dagger(), ak, g)
			}
		}
	}

	d := matrix.ZeroLike(ak)
	for i := range n {
		d.Set(i, i, ak.At(i, i))
	}

	return d, v
}

// EigenQR performs eigen decomposition of a matrix using the Schur decomposition.
func EigenQR(m *matrix.Matrix, qr QRFunc, iter int, eps ...float64) (lambdas *matrix.Matrix, vectors *matrix.Matrix) {
	q, t := Schur(m, qr, iter, eps...)
	lambdas, vectors = EigenUpperT(t, eps...)
	return lambdas, matrix.MatMul(q, vectors)
}

// EigenUpperT performs eigen decomposition of an upper triangular matrix.
func EigenUpperT(t *matrix.Matrix, eps ...float64) (lambdas *matrix.Matrix, vectors *matrix.Matrix) {
	lambdas = matrix.ZeroLike(t)
	for i := range t.Rows {
		lambdas.Set(i, i, t.At(i, i))
	}

	vectors = matrix.Zero(t.Rows, t.Rows)
	e := epsilon.E13(eps...)

	for k := range t.Rows {
		x := make([]complex128, t.Rows)
		x[k] = 1.0

		lambdak := t.At(k, k)
		for i := k - 1; i >= 0; i-- {
			diff := t.At(i, i) - lambdak
			if cmplx.Abs(diff) < e {
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
		if nx > e {
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

// IsDiagonal returns true if m is diagonal matrix.
func IsDiagonal(m *matrix.Matrix, eps ...float64) bool {
	if !m.IsSquare() {
		return false
	}

	e := epsilon.E13(eps...)
	for i := range m.Rows {
		for j := range m.Cols {
			if i == j {
				continue
			}

			if cmplx.Abs(m.At(i, j)) > e {
				return false
			}
		}
	}

	return true
}
