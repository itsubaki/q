package matrix_test

import (
	"fmt"
	"math"
	"math/cmplx"
	"math/rand"
	"testing"

	"github.com/itsubaki/q/math/epsilon"
	"github.com/itsubaki/q/math/matrix"
)

func Example_exp() {
	rx := func(theta float64) *matrix.Matrix {
		v := complex(theta/2, 0)
		return matrix.New(
			[]complex128{cmplx.Cos(v), -1i * cmplx.Sin(v)},
			[]complex128{-1i * cmplx.Sin(v), cmplx.Cos(v)},
		)
	}

	expA := func(x *matrix.Matrix, theta float64, iter int) *matrix.Matrix {
		V, D := matrix.EigenJacobi(x, iter)

		for i := range D.Rows {
			D.Set(i, i, cmplx.Exp(D.At(i, i)*-1i*complex(theta/2, 0)))
		}

		return matrix.MatMul(V, D, V.Dagger())
	}

	expB := func(x *matrix.Matrix, theta float64, iter int) *matrix.Matrix {
		ix := x.Mul(-1i * complex(theta/2, 0))
		V, D := matrix.EigenJacobi(ix, iter)

		for i := range D.Rows {
			D.Set(i, i, cmplx.Exp(D.At(i, i)))
		}

		return matrix.MatMul(V, D, V.Dagger())
	}

	x := matrix.New(
		[]complex128{0, 1},
		[]complex128{1, 0},
	)

	theta := rand.Float64()
	expXa := expA(x, theta, 10)
	expXb := expB(x, theta, 10)

	fmt.Println(expXa.Equals(rx(theta)))
	fmt.Println(expXb.Equals(rx(theta)))

	// Output:
	// true
	// true
}

func Example_pow0p5() {
	a := matrix.New(
		[]complex128{0, 1},
		[]complex128{1, 0},
	)

	iter := 10
	V, D := matrix.EigenJacobi(a, iter)

	for i := range D.Rows {
		D.Set(i, i, cmplx.Pow(D.At(i, i), 0.5))
	}

	sqrtx := matrix.MatMul(V, D, V.Dagger())
	for _, row := range sqrtx.Seq2() {
		fmt.Printf("%.3f\n", row)
	}

	for _, row := range matrix.MatMul(sqrtx, sqrtx).Seq2() {
		fmt.Printf("%.3f\n", row)
	}

	// Output:
	// [(0.500+0.500i) (0.500-0.500i)]
	// [(0.500-0.500i) (0.500+0.500i)]
	// [(0.000+0.000i) (1.000-0.000i)]
	// [(1.000-0.000i) (-0.000+0.000i)]
}

func Example_pow1p5() {
	a := matrix.New(
		[]complex128{0, 1},
		[]complex128{1, 0},
	)

	iter := 10
	V, D := matrix.EigenJacobi(a, iter)

	for i := range D.Rows {
		D.Set(i, i, cmplx.Pow(D.At(i, i), 1.5))
	}

	x1p5 := matrix.MatMul(V, D, V.Dagger())
	for _, row := range x1p5.Seq2() {
		fmt.Printf("%.3f\n", row)
	}

	for _, row := range matrix.MatMul(x1p5, x1p5).Seq2() {
		fmt.Printf("%.3f\n", row)
	}

	// Output:
	// [(0.500-0.500i) (0.500+0.500i)]
	// [(0.500+0.500i) (0.500-0.500i)]
	// [(0.000+0.000i) (1.000-0.000i)]
	// [(1.000-0.000i) (-0.000+0.000i)]
}

func ExampleEigenJacobi_x() {
	a := matrix.New(
		[]complex128{0, 1},
		[]complex128{1, 0},
	)

	iter := 10
	V, D := matrix.EigenJacobi(a, iter)

	for _, row := range V.Seq2() {
		fmt.Printf("%.3f\n", row)
	}

	for _, row := range D.Seq2() {
		fmt.Printf("%.3f\n", row)
	}

	// Output:
	// [(0.707+0.000i) (-0.707+0.000i)]
	// [(0.707+0.000i) (0.707+0.000i)]
	// [(1.000+0.000i) (0.000+0.000i)]
	// [(0.000+0.000i) (-1.000+0.000i)]
}

func ExampleEigenJacobi_cx() {
	a := matrix.New(
		[]complex128{1, 0, 0, 0},
		[]complex128{0, 1, 0, 0},
		[]complex128{0, 0, 0, 1},
		[]complex128{0, 0, 1, 0},
	)

	iter := 10
	V, D := matrix.EigenJacobi(a, iter)

	for _, row := range V.Seq2() {
		fmt.Printf("%.3f\n", row)
	}

	for _, row := range D.Seq2() {
		fmt.Printf("%.3f\n", row)
	}

	// Output:
	// [(1.000+0.000i) (0.000+0.000i) (0.000+0.000i) (0.000+0.000i)]
	// [(0.000+0.000i) (1.000+0.000i) (0.000+0.000i) (0.000+0.000i)]
	// [(0.000+0.000i) (0.000+0.000i) (0.707+0.000i) (-0.707+0.000i)]
	// [(0.000+0.000i) (0.000+0.000i) (0.707+0.000i) (0.707+0.000i)]
	// [(1.000+0.000i) (0.000+0.000i) (0.000+0.000i) (0.000+0.000i)]
	// [(0.000+0.000i) (1.000+0.000i) (0.000+0.000i) (0.000+0.000i)]
	// [(0.000+0.000i) (0.000+0.000i) (1.000+0.000i) (0.000+0.000i)]
	// [(0.000+0.000i) (0.000+0.000i) (0.000+0.000i) (-1.000+0.000i)]
}

func ExampleEigenJacobi_h() {
	a := matrix.New(
		[]complex128{1 / math.Sqrt2, 1 / math.Sqrt2},
		[]complex128{1 / math.Sqrt2, -1 / math.Sqrt2},
	)

	iter := 10
	V, D := matrix.EigenJacobi(a, iter)

	for _, row := range V.Seq2() {
		fmt.Printf("%.3f\n", row)
	}

	for _, row := range D.Seq2() {
		fmt.Printf("%.3f\n", row)
	}

	// Output:
	// [(0.924+0.000i) (-0.383+0.000i)]
	// [(0.383+0.000i) (0.924+0.000i)]
	// [(1.000+0.000i) (0.000+0.000i)]
	// [(0.000+0.000i) (-1.000+0.000i)]
}

func TestEigenJacobi(t *testing.T) {
	rx := func(theta float64) *matrix.Matrix {
		v := complex(theta/2, 0)
		return matrix.New(
			[]complex128{cmplx.Cos(v), -1i * cmplx.Sin(v)},
			[]complex128{-1i * cmplx.Sin(v), cmplx.Cos(v)},
		)
	}

	cases := []struct {
		a *matrix.Matrix
	}{
		{
			a: matrix.New(
				[]complex128{0, 1},
				[]complex128{1, 0},
			),
		},
		{
			a: matrix.New(
				[]complex128{0, -1i},
				[]complex128{1i, 0},
			),
		},
		{
			a: matrix.New(
				[]complex128{1, 0},
				[]complex128{0, -1},
			),
		},
		{
			a: matrix.New(
				[]complex128{1 / math.Sqrt2, 1 / math.Sqrt2},
				[]complex128{1 / math.Sqrt2, -1 / math.Sqrt2},
			),
		},
		{
			a: matrix.New(
				[]complex128{1, 0},
				[]complex128{0, 1i},
			),
		},
		{
			a: matrix.New(
				[]complex128{1, 0},
				[]complex128{0, cmplx.Exp(1i * math.Pi / 4)},
			),
		},
		{
			a: matrix.New(
				[]complex128{1, 0},
				[]complex128{0, cmplx.Exp(complex(0, rand.Float64()))},
			),
		},
		{
			a: matrix.New(
				[]complex128{0, 0, 0, 1},
				[]complex128{0, 0, 1, 0},
				[]complex128{0, 1, 0, 0},
				[]complex128{1, 0, 0, 0},
			),
		},
		{
			a: matrix.New(
				[]complex128{1, 0, 0, 0},
				[]complex128{0, 1, 0, 0},
				[]complex128{0, 0, 0, 1},
				[]complex128{0, 0, 1, 0},
			),
		},
		{
			a: rx(rand.Float64()),
		},
	}

	for _, c := range cases {
		V, D := matrix.EigenJacobi(c.a, 10)

		if !D.IsDiagonal() {
			t.Errorf("D is not diagonal")
		}

		if !V.IsUnitary() {
			t.Errorf("V * V^dagger does not equal I")
		}

		if !matrix.MatMul(V, D, V.Dagger()).Equals(c.a) {
			t.Errorf("V * D * V^dagger does not equal a")
			for _, row := range V.Seq2() {
				t.Log(row)
			}
		}
	}
}

func TestEigenQR(t *testing.T) {
	cases := []struct {
		a *matrix.Matrix
	}{
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
	}

	for _, qr := range []matrix.QRFunc{
		matrix.QR,
		matrix.QRHH,
	} {
		for _, c := range cases {
			D, P := matrix.EigenQR(c.a, qr, 20)

			if !D.IsDiagonal() {
				t.Errorf("D is not diagonal")
			}

			if !matrix.MatMul(P, D, P.Inverse()).Equals(c.a) {
				t.Errorf("P * D * P^-1 does not equal a")
			}
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

		if !D.IsDiagonal(c.eps) {
			t.Errorf("D is not diagonal")
		}

		if !matrix.MatMul(P, D, P.Inverse()).Equals(c.t, c.eps) {
			t.Errorf("P * D * P^-1 does not equal t")
		}
	}
}
