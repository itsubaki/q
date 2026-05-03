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
		var p, q int

		// find the largest off-diagonal element in ak.
		for i := range n - 1 {
			for j := i + 1; j < n; j++ {
				val := cmplx.Abs(ak.At(i, j))
				if val > max {
					max, p, q = val, i, j
				}
			}
		}

		if epsilon.IsZeroF64(max, tol...) {
			break
		}

		// compute the rotation angle.
		a, b, c := ak.At(p, p), ak.At(q, q), ak.At(p, q)
		theta := math.Pi / 4
		if !epsilon.IsZero(b-a, tol...) {
			theta = 0.5 * math.Atan2(2*cmplx.Abs(c), real(b-a))
		}

		phase := cmplx.Rect(1, cmplx.Phase(c))
		cos := complex(math.Cos(theta), 0)
		sin := complex(math.Sin(theta), 0) * phase

		// construct the Givens rotation matrix.
		for i := range n {
			if i == p || i == q {
				continue
			}

			aip, aiq := ak.At(i, p), ak.At(i, q)
			ak.Set(i, p, cos*aip-cmplx.Conj(sin)*aiq)
			ak.Set(p, i, cmplx.Conj(ak.At(i, p)))
			ak.Set(i, q, sin*aip+cos*aiq)
			ak.Set(q, i, cmplx.Conj(ak.At(i, q)))
		}

		// update the diagonal elements.
		absSin2 := real(sin * cmplx.Conj(sin))
		ak.Set(p, p, cos*cos*a+complex(absSin2, 0)*b-2*cos*cmplx.Conj(sin)*c)
		ak.Set(q, q, complex(absSin2, 0)*a+cos*cos*b+2*cos*cmplx.Conj(sin)*c)
		ak.Set(p, q, 0)
		ak.Set(q, p, 0)

		// update the eigenvector matrix.
		for i := range n {
			vip, viq := v.At(i, p), v.At(i, q)
			v.Set(i, p, cos*vip-cmplx.Conj(sin)*viq)
			v.Set(i, q, sin*vip+cos*viq)
		}
	}

	// construct the diagonal matrix of eigenvalues.
	d := matrix.ZeroLike(ak)
	for i := range n {
		val := ak.At(i, i)
		if epsilon.IsZero(val, tol...) {
			val = 0
		}

		d.Set(i, i, val)
	}

	return v, d
}
