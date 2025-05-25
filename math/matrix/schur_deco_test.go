package matrix_test

import (
	"math"
	"math/cmplx"
	"testing"

	"github.com/itsubaki/q/math/epsilon"
	"github.com/itsubaki/q/math/matrix"
)

func TestSchur(t *testing.T) {
	cases := []struct {
		a    *matrix.Matrix
		qr   matrix.QRFunc
		iter int
		eps  float64
	}{
		{
			a: matrix.New(
				[]complex128{1 + 1i, 2 - 1i},
				[]complex128{3 + 4i, 4},
			),
			qr:   matrix.QR,
			iter: 20,
			eps:  epsilon.E13(),
		},
		{
			a: matrix.New(
				[]complex128{0 + 1i, 1 - 1i, 2},
				[]complex128{0, 3 + 3i, 1 + 1i},
				[]complex128{0, 0, 4 - 4i},
			),
			qr:   matrix.QR,
			iter: 20,
			eps:  epsilon.E13(),
		},
		{
			a: matrix.New(
				[]complex128{1, 2},
				[]complex128{3, 4},
			),
			qr:   matrix.QR,
			iter: 20,
			eps:  epsilon.E13(),
		},
		{
			a: matrix.New(
				[]complex128{1, 2, 3},
				[]complex128{3, 4, 5},
				[]complex128{7, 8, 10},
			),
			qr:   matrix.QR,
			iter: 20,
			eps:  epsilon.E13(),
		},
		{
			a: matrix.New(
				[]complex128{1, 2},
				[]complex128{2, 3},
			),
			qr:   matrix.QR,
			iter: 20,
			eps:  epsilon.E13(),
		},
		{
			a: matrix.New(
				[]complex128{1, 1i},
				[]complex128{-1i, 3},
			),
			qr:   matrix.QR,
			iter: 20,
			eps:  epsilon.E13(),
		},
		{
			a: matrix.New(
				[]complex128{1, 2, 3},
				[]complex128{0, 4, 5},
				[]complex128{0, 0, 6},
			),
			qr:   matrix.QR,
			iter: 10,
			eps:  epsilon.E13(),
		},
	}

	for _, c := range cases {
		Q, T := matrix.Schur(c.a, c.qr, c.iter, c.eps)

		// check Q^dagger * Q = I
		if !matrix.MatMul(Q.Dagger(), Q).IsUnitary() {
			t.Errorf("Q is not unitary")
		}

		// check Q * Q^dagger = I
		if !matrix.MatMul(Q, Q.Dagger()).IsUnitary() {
			t.Errorf("Q is not unitary")
		}

		// check Q * T * Q^dagger = a
		if !matrix.MatMul(Q, T, Q.Dagger()).Equals(c.a) {
			t.Errorf("Q * T * Q^dagger does not equal a")
		}

		// check T is upper triangular
		if !T.IsUpperTriangular(c.eps) {
			t.Errorf("T is not upper triangular")
		}
	}
}

func TestShurQRHH(t *testing.T) {
	cases := []struct {
		a    *matrix.Matrix
		qr   matrix.QRFunc
		iter int
		eps  float64
	}{
		{
			a: matrix.New(
				[]complex128{1 + 1i, 2 - 1i},
				[]complex128{3 + 4i, 4},
			),
			qr:   matrix.QRHH,
			iter: 20,
			eps:  epsilon.E13(),
		},
		{
			a: matrix.New(
				[]complex128{0 + 1i, 1 - 1i, 2},
				[]complex128{0, 3 + 3i, 1 + 1i},
				[]complex128{0, 0, 4 - 4i},
			),
			qr:   matrix.QRHH,
			iter: 20,
			eps:  epsilon.E13(),
		},
		{
			a: matrix.New(
				[]complex128{1, 0},
				[]complex128{0, cmplx.Exp(1i * math.Pi / 4)},
			),
			qr:   matrix.QRHH,
			iter: 20,
			eps:  epsilon.E13(),
		},
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
		{
			a: matrix.New(
				[]complex128{1, 2},
				[]complex128{2, 3},
			),
			qr:   matrix.QRHH,
			iter: 20,
			eps:  epsilon.E13(),
		},
		{
			a: matrix.New(
				[]complex128{complex(1, 0), complex(0, 1)},
				[]complex128{complex(0, -1), complex(3, 0)},
			),
			qr:   matrix.QRHH,
			iter: 20,
			eps:  epsilon.E13(),
		},
		{
			a: matrix.New(
				[]complex128{1, 2, 3},
				[]complex128{0, 4, 5},
				[]complex128{0, 0, 6},
			),
			qr:   matrix.QRHH,
			iter: 10,
			eps:  epsilon.E13(),
		},
	}

	for _, c := range cases {
		Q, T := matrix.Schur(c.a, c.qr, c.iter, c.eps)

		// check Q^dagger * Q = I
		if !matrix.MatMul(Q.Dagger(), Q).IsUnitary() {
			t.Errorf("Q is not unitary")
		}

		// check Q * Q^dagger = I
		if !matrix.MatMul(Q, Q.Dagger()).IsUnitary() {
			t.Errorf("Q is not unitary")
		}

		// check Q * T * Q^dagger = a
		if !matrix.MatMul(Q, T, Q.Dagger()).Equals(c.a) {
			t.Errorf("Q * T * Q^dagger does not equal a")
		}

		// check T is upper triangular
		if !T.IsUpperTriangular(c.eps) {
			t.Errorf("T is not upper triangular")
		}
	}
}
