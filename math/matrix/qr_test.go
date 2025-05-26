package matrix_test

import (
	"fmt"
	"math"
	"math/cmplx"
	"testing"

	"github.com/itsubaki/q/math/epsilon"
	"github.com/itsubaki/q/math/matrix"
)

func ExampleQR() {
	a := matrix.New(
		[]complex128{1, 2},
		[]complex128{3, 4},
	)

	for _, qr := range []matrix.QRFunc{
		matrix.QR,
		matrix.QRHH,
	} {
		q, r := qr(a)
		for _, v := range q.Seq2() {
			fmt.Printf("%.3f\n", v)
		}

		for _, v := range r.Seq2() {
			fmt.Printf("%.3f\n", v)
		}

		fmt.Println(q.IsUnitary())
		fmt.Println(matrix.MatMul(q, r).Equals(a))
	}

	// Output:
	// [(0.316+0.000i) (0.949+0.000i)]
	// [(0.949+0.000i) (-0.316+0.000i)]
	// [(3.162+0.000i) (4.427+0.000i)]
	// [(0.000+0.000i) (0.632+0.000i)]
	// true
	// true
	// [(-0.316+0.000i) (-0.949+0.000i)]
	// [(-0.949+0.000i) (0.316+0.000i)]
	// [(-3.162+0.000i) (-4.427+0.000i)]
	// [(0.000+0.000i) (-0.632+0.000i)]
	// true
	// true
}

func ExampleQR_x3() {
	a := matrix.New(
		[]complex128{1, 2, 3},
		[]complex128{3, 4, 5},
		[]complex128{7, 8, 10},
	)

	for _, qr := range []matrix.QRFunc{
		matrix.QR,
		matrix.QRHH,
	} {
		q, r := qr(a)
		for _, v := range q.Seq2() {
			fmt.Printf("%.3f\n", v)
		}

		for _, v := range r.Seq2() {
			fmt.Printf("%.3f\n", v)
		}

		fmt.Println(q.IsUnitary())
		fmt.Println(matrix.MatMul(q, r).Equals(a))
	}

	// Output:
	// [(0.130+0.000i) (0.835+0.000i) (0.535+0.000i)]
	// [(0.391+0.000i) (0.452+0.000i) (-0.802+0.000i)]
	// [(0.911+0.000i) (-0.313+0.000i) (0.267+0.000i)]
	// [(7.681+0.000i) (9.113+0.000i) (11.457+0.000i)]
	// [(0.000+0.000i) (0.974+0.000i) (1.635+0.000i)]
	// [(0.000+0.000i) (0.000+0.000i) (0.267+0.000i)]
	// true
	// true
	// [(-0.130+0.000i) (-0.835+0.000i) (-0.535+0.000i)]
	// [(-0.391+0.000i) (-0.452+0.000i) (0.802+0.000i)]
	// [(-0.911+0.000i) (0.313+0.000i) (-0.267+0.000i)]
	// [(-7.681+0.000i) (-9.113+0.000i) (-11.457+0.000i)]
	// [(0.000+0.000i) (-0.974+0.000i) (-1.635+0.000i)]
	// [(0.000+0.000i) (0.000+0.000i) (-0.267+0.000i)]
	// true
	// true
}

func ExampleQR_identity() {
	a := matrix.Identity(2)

	for _, qr := range []matrix.QRFunc{
		matrix.QR,
		matrix.QRHH,
	} {
		q, r := qr(a)
		for _, v := range q.Seq2() {
			fmt.Printf("%.3f\n", v)
		}

		for _, v := range r.Seq2() {
			fmt.Printf("%.3f\n", v)
		}

		fmt.Println(q.IsUnitary())
		fmt.Println(matrix.MatMul(q, r).Equals(a))
	}

	// Output:
	// [(1.000+0.000i) (0.000+0.000i)]
	// [(0.000+0.000i) (1.000+0.000i)]
	// [(1.000+0.000i) (0.000+0.000i)]
	// [(0.000+0.000i) (1.000+0.000i)]
	// true
	// true
	// [(-1.000+0.000i) (0.000+0.000i)]
	// [(0.000+0.000i) (1.000+0.000i)]
	// [(-1.000+0.000i) (0.000+0.000i)]
	// [(0.000+0.000i) (1.000+0.000i)]
	// true
	// true
}

func ExampleQR_orthogonal() {
	v := complex(1/math.Sqrt2, 0)
	a := matrix.New(
		[]complex128{v, v},
		[]complex128{v, -v},
	)

	for _, qr := range []matrix.QRFunc{
		matrix.QR,
		matrix.QRHH,
	} {
		q, r := qr(a)
		for _, v := range q.Seq2() {
			fmt.Printf("%.3f\n", v)
		}

		for _, v := range r.Seq2() {
			fmt.Printf("%.3f\n", v)
		}

		fmt.Println(q.IsUnitary())
		fmt.Println(matrix.MatMul(q, r).Equals(a))
	}

	// Output:
	// [(0.707+0.000i) (0.707+0.000i)]
	// [(0.707+0.000i) (-0.707+0.000i)]
	// [(1.000+0.000i) (0.000+0.000i)]
	// [(0.000+0.000i) (1.000+0.000i)]
	// true
	// true
	// [(-0.707+0.000i) (-0.707+0.000i)]
	// [(-0.707+0.000i) (0.707+0.000i)]
	// [(-1.000+0.000i) (0.000+0.000i)]
	// [(0.000+0.000i) (-1.000-0.000i)]
	// true
	// true
}

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

	for _, qr := range []matrix.QRFunc{
		matrix.QR,
		matrix.QRHH,
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

func ExampleQR_zero() {
	a := matrix.Zero(3, 3)

	for _, qr := range []matrix.QRFunc{
		matrix.QR,
		matrix.QRHH,
	} {
		_, r := qr(a)
		fmt.Println(r.IsZero())
	}

	// Output:
	// true
	// true
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

	for _, qr := range []matrix.QRFunc{
		matrix.QR,
		matrix.QRHH,
	} {
		for _, c := range cases {
			Q, R := qr(c.in)

			if !Q.IsUnitary() {
				t.Errorf("Q is not unitary")
			}

			if !matrix.MatMul(Q, R).Equals(c.in) {
				t.Errorf("matmul(Q, R) does not equal a")
			}

			if !R.IsUpperTriangular() {
				t.Errorf("R is not upper triangular")
			}
		}
	}
}
