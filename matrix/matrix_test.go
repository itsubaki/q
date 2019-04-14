package matrix_test

import (
	"math/cmplx"
	"testing"

	"github.com/axamon/q/matrix"
)

func TestInverse(t *testing.T) {
	m := matrix.New(
		[]complex128{1, 2, 0, -1},
		[]complex128{-1, 1, 2, 0},
		[]complex128{2, 0, 1, 1},
		[]complex128{1, -2, -1, 1},
	)

	inv := m.Inverse()
	im := m.Apply(inv)

	mm, nn := im.Dimension()
	for i := 0; i < mm; i++ {
		for j := 0; j < nn; j++ {
			if i == j {
				if cmplx.Abs(im[i][j]-complex(1, 0)) > 1e-13 {
					t.Errorf("[%v][%v] %v\n", i, j, im[i][j])
				}
				continue
			}
			if cmplx.Abs(im[i][j]-complex(0, 0)) > 1e-13 {
				t.Errorf("[%v][%v] %v\n", i, j, im[i][j])
			}
		}
	}
}

func TestCommutator(t *testing.T) {
	x := matrix.New(
		[]complex128{0, 1},
		[]complex128{1, 0},
	)

	y := matrix.New(
		[]complex128{0, complex(0, -1)},
		[]complex128{complex(0, 1), 0},
	)

	z := matrix.Commutator(x, y)

	expected := matrix.New(
		[]complex128{complex(0, 2), 0},
		[]complex128{0, complex(0, -2)},
	)

	if !z.Equals(expected) {
		t.Fail()
	}
}

func TestAntiCommutator(t *testing.T) {
	x := matrix.New(
		[]complex128{0, 1},
		[]complex128{1, 0},
	)

	y := matrix.New(
		[]complex128{0, complex(0, -1)},
		[]complex128{complex(0, 1), 0},
	)

	z := matrix.Commutator(x, y).Add(matrix.AntiCommutator(x, y))

	expected := y.Apply(x).Mul(2)
	if !z.Equals(expected) {
		t.Fail()
	}
}
