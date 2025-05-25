package matrix_test

import (
	"testing"

	"github.com/itsubaki/q/math/matrix"
)

func TestSchur(t *testing.T) {
	cases := []struct {
		a *matrix.Matrix
	}{
		{
			a: matrix.New(
				[]complex128{1 + 1i, 2 - 1i},
				[]complex128{3 + 4i, 4},
			),
		},
		{
			a: matrix.New(
				[]complex128{0 + 1i, 1 - 1i, 2},
				[]complex128{0, 3 + 3i, 1 + 1i},
				[]complex128{0, 0, 4 - 4i},
			),
		},
		{
			a: matrix.New(
				[]complex128{1, 2},
				[]complex128{3, 4},
			),
		},
		{
			a: matrix.New(
				[]complex128{1, 2, 3},
				[]complex128{3, 4, 5},
				[]complex128{7, 8, 10},
			),
		},
		{
			a: matrix.New(
				[]complex128{1, 2},
				[]complex128{2, 3},
			),
		},
		{
			a: matrix.New(
				[]complex128{1, 1i},
				[]complex128{-1i, 3},
			),
		},
		{
			a: matrix.New(
				[]complex128{1, 2, 3},
				[]complex128{0, 4, 5},
				[]complex128{0, 0, 6},
			),
		},
	}

	for _, qr := range []matrix.QRFunc{
		matrix.QR,
		matrix.QRHH,
	} {
		for _, c := range cases {
			Q, T := matrix.Schur(c.a, qr, 20)

			if !Q.IsUnitary() {
				t.Errorf("Q is not unitary")
			}

			if !matrix.MatMul(Q, T, Q.Dagger()).Equals(c.a) {
				t.Errorf("Q * T * Q^dagger does not equal a")
			}

			if !T.IsUpperTriangular() {
				t.Errorf("T is not upper triangular")
			}
		}
	}
}
