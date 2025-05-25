package matrix_test

import (
	"testing"

	"github.com/itsubaki/q/math/epsilon"
	"github.com/itsubaki/q/math/matrix"
)

func TestEigen(t *testing.T) {
	cases := []struct {
		a    *matrix.Matrix
		qr   matrix.QRFunc
		iter int
		eps  float64
	}{
		{
			a: matrix.New(
				[]complex128{1, 2},
				[]complex128{3, 4},
			),
			qr:   matrix.QRHH,
			iter: 20,
			eps:  epsilon.E13(),
		},
		{
			a: matrix.New(
				[]complex128{1, 2, 3},
				[]complex128{3, 4, 5},
				[]complex128{7, 8, 10},
			),
			qr:   matrix.QRHH,
			iter: 20,
			eps:  epsilon.E13(),
		},
	}

	for _, c := range cases {
		D, P := matrix.Eigen(c.a, c.qr, c.iter, c.eps)

		if !D.IsDiagonal() {
			t.Errorf("D is not diagonal")
		}

		if !matrix.MatMul(P, D, P.Inverse()).Equals(c.a, c.eps) {
			t.Errorf("P * D * P^-1 does not equal a")
		}
	}
}

func TestEigenUpperT(t *testing.T) {
	cases := []struct {
		t   *matrix.Matrix
		eps float64
	}{
		{
			t: matrix.New(
				[]complex128{1, 2},
				[]complex128{0, 3},
			),
			eps: epsilon.E13(),
		},
		{
			t: matrix.New(
				[]complex128{1, 0, 0},
				[]complex128{0, 2, 0},
				[]complex128{0, 0, 3},
			),
			eps: epsilon.E13(),
		},
		{
			t: matrix.New(
				[]complex128{1, 2, 3, 4, 5},
				[]complex128{0, 2, 3, 4, 5},
				[]complex128{0, 0, 3, 4, 5},
				[]complex128{0, 0, 0, 4, 5},
				[]complex128{0, 0, 0, 0, 5},
			),
			eps: epsilon.E13(),
		},
		{
			t: matrix.New(
				[]complex128{1 + 1i, 2 - 1i, 3 + 0.5i},
				[]complex128{0, 2 + 2i, 1 - 0.5i},
				[]complex128{0, 0, 3 - 1i},
			),
			eps: epsilon.E13(),
		},
		{

			t: matrix.New(
				[]complex128{5, 0, 0, 1},
				[]complex128{0, 3, 0, 0},
				[]complex128{0, 0, 2, 0},
				[]complex128{0, 0, 0, 1},
			),
			eps: epsilon.E13(),
		},
		{
			t: matrix.New(
				[]complex128{10, 0, 0, 0, 0, 2},
				[]complex128{0, 9, 0, 0, 0, 0},
				[]complex128{0, 0, 8, 0, 0, 0},
				[]complex128{0, 0, 0, 7, 0, 0},
				[]complex128{0, 0, 0, 0, 6, 0},
				[]complex128{0, 0, 0, 0, 0, 5},
			),
			eps: epsilon.E13(),
		},
		{
			t: matrix.New(
				[]complex128{1, 0.0001},
				[]complex128{0, 1},
			),
			eps: 1e-2,
		},
	}

	for _, c := range cases {
		D, P := matrix.EigenUpperT(c.t, c.eps)

		if !D.IsDiagonal() {
			t.Errorf("D is not diagonal")
		}

		if !matrix.MatMul(P, D, P.Inverse()).Equals(c.t, c.eps) {
			t.Errorf("P * D * P^-1 does not equal t")
		}
	}
}
