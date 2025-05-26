package matrix_test

import (
	"fmt"
	"math/cmplx"
	"math/rand"
	"testing"

	"github.com/itsubaki/q/math/matrix"
)

func ExampleHessenberg() {
	ry := func(theta float64) *matrix.Matrix {
		v := complex(theta/2, 0)
		return matrix.New(
			[]complex128{cmplx.Cos(v), -1 * cmplx.Sin(v)},
			[]complex128{cmplx.Sin(v), cmplx.Cos(v)},
		)
	}

	a := ry(rand.Float64())
	aa := matrix.TensorProduct(a, a)
	q, h := matrix.Hessenberg(aa)

	fmt.Println(q.IsUnitary())
	fmt.Println(matrix.MatMul(q, h, q.Dagger()).Equals(aa))

	// Output:
	// true
	// true
}

func ExampleHessenberg_qr() {
	ry := func(theta float64) *matrix.Matrix {
		v := complex(theta/2, 0)
		return matrix.New(
			[]complex128{cmplx.Cos(v), -1 * cmplx.Sin(v)},
			[]complex128{cmplx.Sin(v), cmplx.Cos(v)},
		)
	}

	a := ry(rand.Float64())
	aa := matrix.TensorProduct(a, a)

	qq, h := matrix.Hessenberg(aa)
	q, r := matrix.QRHH(h)

	fmt.Println(q.IsUnitary())
	fmt.Println(matrix.MatMul(q, r).Equals(h))
	fmt.Println(matrix.MatMul(qq, q, r, qq.Dagger()).Equals(aa))

	// Output:
	// true
	// true
	// true
}

func TestHessenberg(t *testing.T) {
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
				[]complex128{1, 2, 3},
				[]complex128{2, 5, 6},
				[]complex128{3, 6, 9},
			),
		},
		{
			matrix.New(
				[]complex128{2 + 0i, 1 - 1i, 0},
				[]complex128{1 + 1i, 3 + 0i, 4 - 2i},
				[]complex128{0, 4 + 2i, 1 + 0i},
			),
		},
		{
			matrix.Identity(3),
		},
		{
			matrix.New(
				[]complex128{1 + 1i, 2, 3 - 1i, 4},
				[]complex128{0, 5 + 2i, 6, 7 - 1i},
				[]complex128{0, 0, 8 + 1i, 9},
				[]complex128{0, 0, 0, 10 + 3i},
			),
		},
	}

	for _, c := range cases {
		Q, T := matrix.Hessenberg(c.in)

		if !Q.IsUnitary() {
			t.Errorf("Q is not unitary")
		}

		if !T.IsHessenberg() {
			t.Errorf("T is not in Hessenberg form")
		}

		if !matrix.MatMul(Q, T, Q.Dagger()).Equals(c.in) {
			t.Errorf("Q * T * Q^dagger does not equal a")
		}
	}
}
