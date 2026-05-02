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

	V, D := eigen.Jacobi(x, 10)
	D.Fdiag(func(v complex128) complex128 { return cmplx.Pow(v, 0.5) })

	sqrtx := matrix.MatMul(V, D, V.Dagger())
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
	}

	for _, c := range cases {
		if !c.in.Equal(c.in.Dagger()) {
			t.Errorf("input is not Hermitian")
		}

		V, D := eigen.Jacobi(c.in, 10)

		if !V.IsUnitary() {
			t.Errorf("V * V^dagger does not equal I")
		}

		if !D.IsDiagonal() {
			t.Errorf("D is not diagonal")
		}

		if !matrix.MatMul(V, D, V.Dagger()).Equal(c.in) {
			t.Errorf("V * D * V^dagger does not equal a")
			for _, row := range V.Seq2() {
				t.Log(row)
			}
		}
	}
}
