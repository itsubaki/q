package decomp

import (
	"math/cmplx"

	"github.com/itsubaki/q/math/epsilon"
	"github.com/itsubaki/q/math/matrix"
)

type ParlettF func(z complex128) complex128

// Parlett computes the matrix function f(t) using Parlett recursion.
// The input matrix t must be an upper triangular matrix.
func Parlett(t *matrix.Matrix, f, df ParlettF, eps ...float64) *matrix.Matrix {
	n := t.Rows
	a := matrix.Zero(n, n)

	for i := range n {
		a.Set(i, i, f(t.At(i, i)))
	}

	for d := 1; d < n; d++ {
		for i := range n - d {
			j := i + d
			tij := t.At(i, j)

			// numnerator
			num := tij * (a.At(j, j) - a.At(i, i))
			for k := i + 1; k < j; k++ {
				num += t.At(i, k)*a.At(k, j) - a.At(i, k)*t.At(k, j)
			}

			// denominator
			denom := t.At(j, j) - t.At(i, i)
			if cmplx.Abs(denom) < epsilon.E13(eps...) {
				a.Set(i, j, tij*df(t.At(i, i))+num)
				continue
			}

			// set
			a.Set(i, j, num/denom)
		}
	}

	return a
}
