package eigen_test

import (
	"fmt"
	"math"
	"math/cmplx"
	"testing"

	"github.com/itsubaki/q/math/eigen"
	"github.com/itsubaki/q/math/matrix"
)

func Example_sqrt() {
	x := matrix.New(
		[]complex128{0, 1},
		[]complex128{1, 0},
	)

	v, d := eigen.Jacobi(x, 10)
	d.Fdiag(func(v complex128) complex128 { return cmplx.Pow(v, 0.5) })

	sqrtx := matrix.MatMul(v, d, v.Dagger())
	for _, row := range sqrtx.Seq2() {
		fmt.Printf("%.3f\n", row)
	}

	sqrtx2 := matrix.MatMul(sqrtx, sqrtx)
	fmt.Println(sqrtx2.Equal(x))

	// Output:
	// [(0.500+0.500i) (0.500-0.500i)]
	// [(0.500-0.500i) (0.500+0.500i)]
	// true
}

func TestJacobi(t *testing.T) {
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
				[]complex128{0, -1i},
				[]complex128{1i, 0},
			),
		},
		{
			matrix.New(
				[]complex128{2, 1 + 2i, -1i},
				[]complex128{1 - 2i, 3, 4 + 1i},
				[]complex128{1i, 4 - 1i, -1},
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
		{
			matrix.New(
				[]complex128{1e-9, 0},
				[]complex128{0, 1},
			),
		},
	}

	for _, c := range cases {
		if !c.in.IsHermitian() {
			t.Errorf("input is not Hermitian")
		}

		v, d := eigen.Jacobi(c.in, 10)
		if !v.IsUnitary() {
			t.Errorf("v * v^dagger does not equal I")
		}

		if !d.IsDiagonal() {
			t.Errorf("d is not diagonal")
		}

		if !matrix.MatMul(v, d, v.Dagger()).Equal(c.in) {
			t.Errorf("v * d * v^dagger does not equal a")
		}
	}
}
