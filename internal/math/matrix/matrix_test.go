package matrix

import (
	"math/cmplx"
	"testing"
)

func TestInverse(t *testing.T) {
	m := New(
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
