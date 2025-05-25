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

	q, r := matrix.QR(a)
	for _, v := range q.Seq2() {
		fmt.Printf("%.3f\n", v)
	}
	fmt.Println()

	for _, v := range r.Seq2() {
		fmt.Printf("%.3f\n", v)
	}
	fmt.Println()

	fmt.Println(matrix.MatMul(q.Dagger(), q).IsUnitary())
	fmt.Println(matrix.MatMul(q, r).Equals(a))

	// Output:
	// [(0.316+0.000i) (0.949+0.000i)]
	// [(0.949+0.000i) (-0.316+0.000i)]
	//
	// [(3.162+0.000i) (4.427+0.000i)]
	// [(0.000+0.000i) (0.632+0.000i)]
	//
	// true
	// true
}

func ExampleQR_x3() {
	a := matrix.New(
		[]complex128{1, 2, 3},
		[]complex128{3, 4, 5},
		[]complex128{7, 8, 10},
	)

	q, r := matrix.QR(a)
	for _, v := range q.Seq2() {
		fmt.Printf("%.3f\n", v)
	}
	fmt.Println()

	for _, v := range r.Seq2() {
		fmt.Printf("%.3f\n", v)
	}
	fmt.Println()

	fmt.Println(matrix.MatMul(q.Dagger(), q).IsUnitary())
	fmt.Println(matrix.MatMul(q, r).Equals(a))

	// Output:
	// [(0.130+0.000i) (0.835+0.000i) (0.535+0.000i)]
	// [(0.391+0.000i) (0.452+0.000i) (-0.802+0.000i)]
	// [(0.911+0.000i) (-0.313+0.000i) (0.267+0.000i)]
	//
	// [(7.681+0.000i) (9.113+0.000i) (11.457+0.000i)]
	// [(0.000+0.000i) (0.974+0.000i) (1.635+0.000i)]
	// [(0.000+0.000i) (0.000+0.000i) (0.267+0.000i)]
	//
	// true
	// true
}

func ExampleQR_identity() {
	a := matrix.Identity(3)
	q, r := matrix.QR(a)

	fmt.Println(q.Equals(matrix.Identity(3)))
	fmt.Println(r.Equals(matrix.Identity(3)))

	// Output:
	// true
	// true
}

func ExampleQR_orthogonal() {
	v := complex(1/math.Sqrt2, 0)
	a := matrix.New(
		[]complex128{v, v},
		[]complex128{v, -v},
	)

	q, r := matrix.QR(a)
	fmt.Println(q.Equals(a))
	fmt.Println(r.Equals(matrix.Identity(2)))

	// Output:
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

	_, r := matrix.QR(a)
	for i := range r.Rows {
		fmt.Println(i, ":", isZero(r.Row(i)))
	}

	// Output:
	// 0 : false
	// 1 : true
	// 2 : true
}

func ExampleQR_zero() {
	a := matrix.Zero(3, 3)

	_, r := matrix.QR(a)
	fmt.Println(r.IsZero())

	// Output:
	// true
}

func TestQR(t *testing.T) {
	cases := []struct {
		a   *matrix.Matrix
		eps float64
	}{
		{
			a: matrix.New(
				[]complex128{0, 1},
				[]complex128{1, 0},
			),
			eps: epsilon.E13(),
		},
		{
			a: matrix.New(
				[]complex128{0, -1i},
				[]complex128{1i, 0},
			),
			eps: epsilon.E13(),
		},
		{
			a: matrix.New(
				[]complex128{1 / math.Sqrt2, 1 / math.Sqrt2},
				[]complex128{1 / math.Sqrt2, -1 / math.Sqrt2},
			),
			eps: epsilon.E13(),
		},
		{
			a: matrix.New(
				[]complex128{1, 0},
				[]complex128{0, cmplx.Exp(1i * math.Pi / 4)},
			),
			eps: epsilon.E13(),
		},
		{
			a: matrix.New(
				[]complex128{1, 2},
				[]complex128{3, 4},
			),
			eps: epsilon.E13(),
		},
		{
			a: matrix.New(
				[]complex128{1, 2, 3},
				[]complex128{3, 4, 5},
				[]complex128{7, 8, 10},
			),
			eps: epsilon.E13(),
		},
	}

	for _, c := range cases {
		Q, R := matrix.QR(c.a)

		// check Q^dagger * Q = I
		if !matrix.MatMul(Q.Dagger(), Q).IsUnitary() {
			t.Errorf("Q is not unitary")
		}

		// check Q * Q^dagger = I
		if !matrix.MatMul(Q, Q.Dagger()).IsUnitary() {
			t.Errorf("Q is not unitary")
		}

		// check Q * R = a
		if !matrix.MatMul(Q, R).Equals(c.a) {
			t.Errorf("matmul(Q, R) does not equal a")
		}

		// check R is upper triangular
		if !R.IsUpperTriangular(c.eps) {
			t.Errorf("R is not upper triangular")
		}
	}
}
