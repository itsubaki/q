package decomp_test

import (
	"testing"

	"github.com/itsubaki/q/math/decomp"
	"github.com/itsubaki/q/math/matrix"
)

func TestSchur(t *testing.T) {
	cases := []struct {
		in *matrix.Matrix
	}{
		{
			matrix.New(
				[]complex128{1 + 1i, 2 - 1i},
				[]complex128{3 + 4i, 4},
			),
		},
		{
			matrix.New(
				[]complex128{0 + 1i, 1 - 1i, 2},
				[]complex128{0, 3 + 3i, 1 + 1i},
				[]complex128{0, 0, 4 - 4i},
			),
		},
		{
			matrix.New(
				[]complex128{1, 2},
				[]complex128{3, 4},
			),
		},
		{
			matrix.New(
				[]complex128{1, 2, 3},
				[]complex128{3, 4, 5},
				[]complex128{7, 8, 10},
			),
		},
		{
			matrix.New(
				[]complex128{1, 2},
				[]complex128{2, 3},
			),
		},
		{
			matrix.New(
				[]complex128{1, 1i},
				[]complex128{-1i, 3},
			),
		},
		{
			matrix.New(
				[]complex128{1, 2, 3},
				[]complex128{0, 4, 5},
				[]complex128{0, 0, 6},
			),
		},
	}

	for _, qr := range []decomp.QRFunc{
		decomp.QR,
		decomp.QRHH,
	} {
		for _, c := range cases {
			Q, T := decomp.Schur(c.in, qr, 20)

			if !Q.IsUnitary() {
				t.Errorf("Q is not unitary")
			}

			if !decomp.IsUpperTriangular(T) {
				t.Errorf("T is not upper triangular")
			}

			if !matrix.MatMul(Q, T, Q.Dagger()).Equals(c.in) {
				t.Errorf("Q * T * Q^dagger does not equal a")
			}

		}
	}
}

func TestIsUpperTriangular(t *testing.T) {
	cases := []struct {
		in   *matrix.Matrix
		want bool
	}{
		{
			matrix.New(
				[]complex128{1, 2},
				[]complex128{0, 3},
			),
			true,
		},
		{
			matrix.New(
				[]complex128{1, 2, 3},
				[]complex128{0, 4, 5},
				[]complex128{0, 0, 6},
			),
			true,
		},
		{
			matrix.New(
				[]complex128{1, 2},
				[]complex128{3, 4},
			),
			false,
		},
		{
			matrix.New(
				[]complex128{1, 2},
				[]complex128{0, 4},
				[]complex128{0, 6},
			),
			false,
		},
		{
			matrix.New(
				[]complex128{1, 2, 3},
				[]complex128{0, 4, 5},
			),
			false,
		},
		{
			matrix.New(
				[]complex128{complex(1, 1), complex(2, 2)},
				[]complex128{0, complex(3, -1)},
			),
			true,
		},
		{
			matrix.New(
				[]complex128{complex(1, 0), complex(2, 0)},
				[]complex128{complex(1, 1), complex(3, 0)},
			),
			false,
		},
	}

	for _, c := range cases {
		if decomp.IsUpperTriangular(c.in) != c.want {
			t.Fail()
		}
	}
}
