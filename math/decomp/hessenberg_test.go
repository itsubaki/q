package decomp_test

import (
	"fmt"
	"math/cmplx"
	"math/rand"
	"testing"

	"github.com/itsubaki/q/math/decomp"
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
	q, h := decomp.Hessenberg(aa)

	fmt.Println(q.IsUnitary())
	fmt.Println(decomp.IsHessenberg(h))
	fmt.Println(matrix.MatMul(q, h, q.Dagger()).Equals(aa))

	// Output:
	// true
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

	qq, h := decomp.Hessenberg(aa)
	q, r := decomp.QRHH(h)

	fmt.Println(q.IsUnitary())
	fmt.Println(decomp.IsUpperTriangular(r))
	fmt.Println(matrix.MatMul(q, r).Equals(h))
	fmt.Println(matrix.MatMul(qq, q, r, qq.Dagger()).Equals(aa))

	// Output:
	// true
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
		Q, T := decomp.Hessenberg(c.in)

		if !Q.IsUnitary() {
			t.Errorf("Q is not unitary")
		}

		if !decomp.IsHessenberg(T) {
			t.Errorf("T is not in Hessenberg form")
		}

		if !matrix.MatMul(Q, T, Q.Dagger()).Equals(c.in) {
			t.Errorf("Q * T * Q^dagger does not equal a")
		}
	}
}

func TestIsHessenberg(t *testing.T) {
	cases := []struct {
		in   *matrix.Matrix
		want bool
	}{
		{
			matrix.New(
				[]complex128{1, 2},
				[]complex128{3, 4},
			),
			true,
		},
		{
			matrix.New(
				[]complex128{1, 2, 3},
				[]complex128{4, 5, 6},
				[]complex128{0, 7, 8},
			),
			true,
		},
		{
			matrix.New(
				[]complex128{1, 2, 3},
				[]complex128{4, 5, 6},
				[]complex128{9, 7, 8},
			),
			false,
		},
		{
			matrix.Identity(4),
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
				[]complex128{1, 0, 0},
				[]complex128{2, 3, 0},
				[]complex128{4, 5, 6},
			),
			false,
		},
		{
			matrix.New(
				[]complex128{1, 2, 3},
				[]complex128{4, 5, 6},
				[]complex128{0, 7, 8},
				[]complex128{0, 0, 9},
			),
			false,
		},
	}

	for _, c := range cases {
		if decomp.IsHessenberg(c.in) != c.want {
			t.Fail()
		}
	}
}
