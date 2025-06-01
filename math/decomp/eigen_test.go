package decomp_test

import (
	"fmt"
	"math"
	"math/cmplx"
	"math/rand"
	"testing"

	"github.com/itsubaki/q/math/decomp"
	"github.com/itsubaki/q/math/epsilon"
	"github.com/itsubaki/q/math/matrix"
)

func diagF(a *matrix.Matrix, f func(v complex128) complex128) {
	for i := range a.Rows {
		a.Set(i, i, f(a.At(i, i)))
	}
}

func ExampleEigenJacobi_pow0p5() {
	x := matrix.New(
		[]complex128{0, 1},
		[]complex128{1, 0},
	)

	D, V := decomp.EigenJacobi(x, 10)
	diagF(D, func(v complex128) complex128 {
		return cmplx.Pow(v, 0.5)
	})

	sqrtx := matrix.MatMul(V, D, V.Dagger())
	for _, row := range sqrtx.Seq2() {
		fmt.Printf("%.3f\n", row)
	}

	sqrtx2 := matrix.MatMul(sqrtx, sqrtx)
	fmt.Println(sqrtx2.Equals(x))

	// Output:
	// [(0.500+0.500i) (0.500-0.500i)]
	// [(0.500-0.500i) (0.500+0.500i)]
	// true
}

func ExampleEigenJacobi_pow1p5() {
	x := matrix.New(
		[]complex128{0, 1},
		[]complex128{1, 0},
	)

	D, V := decomp.EigenJacobi(x, 10)
	diagF(D, func(v complex128) complex128 {
		return cmplx.Pow(v, 1.5)
	})

	x1p5 := matrix.MatMul(V, D, V.Dagger())
	for _, row := range x1p5.Seq2() {
		fmt.Printf("%.3f\n", row)
	}

	x1p52 := matrix.MatMul(x1p5, x1p5)
	fmt.Println(x1p52.Equals(x))

	// Output:
	// [(0.500-0.500i) (0.500+0.500i)]
	// [(0.500+0.500i) (0.500-0.500i)]
	// true
}

func ExampleEigenJacobi_exp() {
	exp := func(x *matrix.Matrix, theta float64, iter int) *matrix.Matrix {
		D, V := decomp.EigenJacobi(x, iter)
		diagF(D, func(v complex128) complex128 {
			return cmplx.Exp(-1 * complex(0, theta/2) * v) // exp(-i * theta/2 * v)
		})

		return matrix.MatMul(V, D, V.Dagger())
	}

	x := matrix.New(
		[]complex128{0, 1},
		[]complex128{1, 0},
	)

	theta := rand.Float64()
	expiX := exp(x, theta, 10)
	fmt.Println(expiX.Equals(rx(theta)))

	// Output:
	// true
}

func ExampleEigenJacobi_expiX() {
	exp := func(x *matrix.Matrix, theta float64, iter int) *matrix.Matrix {
		ix := x.Mul(-1 * complex(0, theta/2))

		D, V := decomp.EigenJacobi(ix, iter)
		diagF(D, func(v complex128) complex128 {
			return cmplx.Exp(v)
		})

		return matrix.MatMul(V, D, V.Dagger())
	}

	x := matrix.New(
		[]complex128{0, 1},
		[]complex128{1, 0},
	)

	theta := rand.Float64()
	expiX := exp(x, theta, 10)
	fmt.Println(expiX.Equals(rx(theta)))

	// Output:
	// true
}

func ExampleIsDiagonal() {
	x := matrix.New(
		[]complex128{1, 0},
		[]complex128{0, 2},
		[]complex128{3, 4},
	)

	fmt.Println(decomp.IsDiagonal(x))

	// Output:
	// false
}

func TestEigenJacobi(t *testing.T) {
	cases := []struct {
		in *matrix.Matrix
	}{
		{
			matrix.New(
				[]complex128{1, 0},
				[]complex128{0, 1},
			),
		},
		{
			matrix.New(
				[]complex128{0, 1},
				[]complex128{1, 0},
			),
		},
		{
			matrix.New(
				[]complex128{1, 0},
				[]complex128{0, -1},
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
				[]complex128{0, 0, 0, 1},
				[]complex128{0, 0, 1, 0},
				[]complex128{0, 1, 0, 0},
				[]complex128{1, 0, 0, 0},
			),
		},
		{
			matrix.New(
				[]complex128{1, 0, 0, 0},
				[]complex128{0, 1, 0, 0},
				[]complex128{0, 0, 0, 1},
				[]complex128{0, 0, 1, 0},
			),
		},
	}

	for _, c := range cases {
		if !c.in.Equals(c.in.Dagger()) {
			t.Errorf("input is not Hermitian")
		}

		D, V := decomp.EigenJacobi(c.in, 10)

		if !decomp.IsDiagonal(D) {
			t.Errorf("D is not diagonal")
		}

		if !V.IsUnitary() {
			t.Errorf("V * V^dagger does not equal I")
		}

		if !matrix.MatMul(V, D, V.Dagger()).Equals(c.in) {
			t.Errorf("V * D * V^dagger does not equal a")
			for _, row := range V.Seq2() {
				t.Log(row)
			}
		}
	}
}

func TestEigenQR(t *testing.T) {
	cases := []struct {
		in *matrix.Matrix
	}{
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
			D, P := decomp.EigenQR(c.in, qr, 20)

			if !decomp.IsDiagonal(D) {
				t.Errorf("D is not diagonal")
			}

			if !matrix.MatMul(P, D, P.Inverse()).Equals(c.in) {
				t.Errorf("P * D * P^-1 does not equal a")
			}
		}
	}
}

func TestEigenUpperT(t *testing.T) {
	cases := []struct {
		in  *matrix.Matrix
		eps float64
	}{
		{
			matrix.New(
				[]complex128{1, 2},
				[]complex128{0, 3},
			),
			epsilon.E13(),
		},
		{
			matrix.New(
				[]complex128{1, 0, 0},
				[]complex128{0, 2, 0},
				[]complex128{0, 0, 3},
			),
			epsilon.E13(),
		},
		{
			matrix.New(
				[]complex128{1, 2, 3, 4, 5},
				[]complex128{0, 2, 3, 4, 5},
				[]complex128{0, 0, 3, 4, 5},
				[]complex128{0, 0, 0, 4, 5},
				[]complex128{0, 0, 0, 0, 5},
			),
			epsilon.E13(),
		},
		{
			matrix.New(
				[]complex128{1 + 1i, 2 - 1i, 3 + 0.5i},
				[]complex128{0, 2 + 2i, 1 - 0.5i},
				[]complex128{0, 0, 3 - 1i},
			),
			epsilon.E13(),
		},
		{

			matrix.New(
				[]complex128{5, 0, 0, 1},
				[]complex128{0, 3, 0, 0},
				[]complex128{0, 0, 2, 0},
				[]complex128{0, 0, 0, 1},
			),
			epsilon.E13(),
		},
		{
			matrix.New(
				[]complex128{10, 0, 0, 0, 0, 2},
				[]complex128{0, 9, 0, 0, 0, 0},
				[]complex128{0, 0, 8, 0, 0, 0},
				[]complex128{0, 0, 0, 7, 0, 0},
				[]complex128{0, 0, 0, 0, 6, 0},
				[]complex128{0, 0, 0, 0, 0, 5},
			),
			epsilon.E13(),
		},
		{
			matrix.New(
				[]complex128{1, 0.0001},
				[]complex128{0, 1},
			),
			1e-2,
		},
	}

	for _, c := range cases {
		D, P := decomp.EigenUpperT(c.in, c.eps)

		if !decomp.IsDiagonal(D, c.eps) {
			t.Errorf("D is not diagonal")
		}

		if !matrix.MatMul(P, D, P.Inverse()).Equals(c.in, c.eps) {
			t.Errorf("P * D * P^-1 does not equal t")
		}
	}
}

func TestIsDiagonal(t *testing.T) {
	cases := []struct {
		in   *matrix.Matrix
		want bool
	}{
		{
			matrix.New(
				[]complex128{1, 0},
				[]complex128{0, 2},
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
	}

	for _, c := range cases {
		if decomp.IsDiagonal(c.in) != c.want {
			t.Fail()
		}
	}
}
