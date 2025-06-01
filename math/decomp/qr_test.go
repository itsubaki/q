package decomp_test

import (
	"fmt"
	"math"
	"math/cmplx"
	"testing"

	"github.com/itsubaki/q/math/decomp"
	"github.com/itsubaki/q/math/epsilon"
	"github.com/itsubaki/q/math/matrix"
)

func ExampleQR_rankdeficient() {
	isZero := func(row []complex128) bool {
		for _, v := range row {
			if cmplx.Abs(v) < epsilon.E13() {
				continue
			}

			return false
		}

		return true
	}

	a := matrix.New(
		[]complex128{1, 2, 3},
		[]complex128{2, 4, 6},
		[]complex128{3, 6, 9},
	)

	for _, qr := range []decomp.QRFunc{
		decomp.QR,
		decomp.QRHH,
	} {
		_, r := qr(a)
		for i, row := range r.Seq2() {
			fmt.Println(i, ":", isZero(row))
		}
	}

	// Output:
	// 0 : false
	// 1 : true
	// 2 : true
	// 0 : false
	// 1 : true
	// 2 : true
}

func TestQR(t *testing.T) {
	cases := []struct {
		in *matrix.Matrix
	}{
		{
			matrix.New(
				[]complex128{0, 1},
				[]complex128{1, 0},
			),
		},
		{
			matrix.New(
				[]complex128{0, 0, 0, 1},
				[]complex128{0, 0, 1, 0},
				[]complex128{0, 1, 0, 0},
				[]complex128{1, 0, 0, 0},
			),
		},
		{
			matrix.New(
				[]complex128{0, -1i},
				[]complex128{1i, 0},
			),
		},
		{
			matrix.New(
				[]complex128{1 / math.Sqrt2, 1 / math.Sqrt2},
				[]complex128{1 / math.Sqrt2, -1 / math.Sqrt2},
			),
		},
		{
			matrix.New(
				[]complex128{1, 0},
				[]complex128{0, cmplx.Exp(1i * math.Pi / 4)},
			),
		},
		{
			matrix.New(
				[]complex128{1, 0, 0},
				[]complex128{0, 1, 0},
				[]complex128{0, 0, 1},
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
	}

	for _, qr := range []decomp.QRFunc{
		decomp.QR,
		decomp.QRHH,
	} {
		for _, c := range cases {
			Q, R := qr(c.in)

			if !Q.IsUnitary() {
				t.Errorf("Q is not unitary")
			}

			if !decomp.IsUpperTriangular(R) {
				t.Errorf("R is not upper triangular")
			}

			if !matrix.MatMul(Q, R).Equals(c.in) {
				t.Errorf("matmul(Q, R) does not equal a")
			}
		}
	}
}
