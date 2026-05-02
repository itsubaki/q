package eigen

import (
	"math"
	"math/cmplx"

	"github.com/itsubaki/q/math/epsilon"
	"github.com/itsubaki/q/math/matrix"
)

// Jacobi returns the eigenvectors and eigenvalues of a matrix using the Jacobi method.
// The input matrix a must be hermitian.
func Jacobi(a *matrix.Matrix, iter int, tol ...float64) (vectors *matrix.Matrix, lambdas *matrix.Matrix) {
	n := a.Rows
	v, ak := matrix.Identity(n), a.Clone()

	for range iter {
		var max float64
		for i := range n - 1 {
			for j := i + 1; j < n; j++ {
				c := ak.At(i, j)
				abs := cmplx.Abs(c)
				if abs > max {
					max = abs
				}

				if epsilon.IsZeroF64(abs, tol...) {
					continue
				}

				a, b := ak.At(i, i), ak.At(j, j)
				diff, phi := b-a, complex(math.Pi/4, 0)
				if !epsilon.IsZero(diff, tol...) {
					phi = 0.5 * cmplx.Atan(2*c/diff)
				}

				phase := cmplx.Rect(1, cmplx.Phase(c))
				sin := cmplx.Sin(phi) * phase
				cos := cmplx.Cos(phi)

				g := matrix.Identity(n)
				g.Set(i, i, cos)
				g.Set(i, j, -cmplx.Conj(sin))
				g.Set(j, i, sin)
				g.Set(j, j, cos)

				v = matrix.MatMul(v, g)
				ak = matrix.MatMul(g.Dagger(), ak, g)
			}
		}

		if epsilon.IsZeroF64(max, tol...) {
			break
		}
	}

	d := matrix.ZeroLike(ak)
	for i := range n {
		d.Set(i, i, ak.At(i, i))
	}

	return v, d
}
