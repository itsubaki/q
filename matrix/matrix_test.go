package matrix_test

import (
	"fmt"
	"math/cmplx"
	"testing"

	"github.com/axamon/q/matrix"
)

var m = matrix.New(
	[]complex128{1, 2, 3, 4},
	[]complex128{0, 1, 1 + 1i, 2 + 2i},
	[]complex128{0, 0, 1, 1},
	[]complex128{0, 0, 0, 1},
)

func ExampleMatrix_Transpose() {
	for _, r := range m {
		fmt.Println(r)
	}
	fmt.Println()
	mt := m.Transpose()
	for _, r := range mt {
		fmt.Println(r)
	}
	// Output:
	// [(1+0i) (2+0i) (3+0i) (4+0i)]
	// [(0+0i) (1+0i) (1+1i) (2+2i)]
	// [(0+0i) (0+0i) (1+0i) (1+0i)]
	// [(0+0i) (0+0i) (0+0i) (1+0i)]
	//
	// [(1+0i) (0+0i) (0+0i) (0+0i)]
	// [(2+0i) (1+0i) (0+0i) (0+0i)]
	// [(3+0i) (1+1i) (1+0i) (0+0i)]
	// [(4+0i) (2+2i) (1+0i) (1+0i)]
}

func ExampleMatrix_Conjugate() {
	for _, r := range m {
		fmt.Println(r)
	}
	fmt.Println()
	mC := m.Conjugate()
	for _, r := range mC {
		fmt.Println(r)
	}
	// Output:
	// [(1+0i) (2+0i) (3+0i) (4+0i)]
	// [(0+0i) (1+0i) (1+1i) (2+2i)]
	// [(0+0i) (0+0i) (1+0i) (1+0i)]
	// [(0+0i) (0+0i) (0+0i) (1+0i)]
	//
	// [(1-0i) (2-0i) (3-0i) (4-0i)]
	// [(0-0i) (1-0i) (1-1i) (2-2i)]
	// [(0-0i) (0-0i) (1-0i) (1-0i)]
	// [(0-0i) (0-0i) (0-0i) (1-0i)]
}

func ExampleMatrix_Dagger() {
	for _, r := range m {
		fmt.Println(r)
	}
	fmt.Println()
	mD := m.Dagger()
	for _, r := range mD {
		fmt.Println(r)
	}
	// Output:
	// [(1+0i) (2+0i) (3+0i) (4+0i)]
	// [(0+0i) (1+0i) (1+1i) (2+2i)]
	// [(0+0i) (0+0i) (1+0i) (1+0i)]
	// [(0+0i) (0+0i) (0+0i) (1+0i)]
	//
	// [(1-0i) (0-0i) (0-0i) (0-0i)]
	// [(2-0i) (1-0i) (0-0i) (0-0i)]
	// [(3-0i) (1-1i) (1-0i) (0-0i)]
	// [(4-0i) (2-2i) (1-0i) (1-0i)]
}

func ExampleMatrix_IsHermite() {
	mH := matrix.New(
		[]complex128{2 + 0i, 2 + 1i, 4 + 2i},
		[]complex128{2 - 1i, 3 + 0i, 3 + 3i},
		[]complex128{4 - 2i, 3 - 3i, 3 + 0i},
	)
	fmt.Println(m.IsHermite())
	fmt.Println(mH.IsHermite())
	// Output:
	// false
	// true
}

func ExampleMatrix_IsUnitary() {
	fmt.Println(m.IsUnitary())
	mU := matrix.New(
		[]complex128{1, 0, 0},
		[]complex128{0, 1, 0},
		[]complex128{0, 0, 1},
	)
	mU1 := matrix.New(
		[]complex128{0.5 + 0.5i, 0.5 - 0.5i},
		[]complex128{0.4 - 0.5i, 0.5 + 0.5i},
	)
	fmt.Println(mU.IsUnitary())
	fmt.Println(mU1.IsUnitary(0.1))
	// Output:
	// false
	// true
	// true
}

func ExampleMatrix_Trace() {
	fmt.Println(m.Trace())
	// Output:
	// (4+0i)
}

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

func ExampleMatrix_Mul() {
	mMul := m.Mul(2 + 1i)
	for _, r := range mMul {
		fmt.Println(r)
	}
	// Output:
	// [(2+1i) (4+2i) (6+3i) (8+4i)]
	// [(0+0i) (2+1i) (1+3i) (2+6i)]
	// [(0+0i) (0+0i) (2+1i) (2+1i)]
	// [(0+0i) (0+0i) (0+0i) (2+1i)]
}

func ExampleMatrix_TensorProduct() {
	m := matrix.New(
		[]complex128{1 + 1i, 2 + 2i},
		[]complex128{3 + 3i, 0 + 0i},
	)
	mTP := m.TensorProduct(m)
	for _, r := range mTP {
		fmt.Println(r)
	}
	// Output:
	// [(0+2i) (0+4i) (0+4i) (0+8i)]
	// [(0+6i) (0+0i) (0+12i) (0+0i)]
	// [(0+6i) (0+12i) (0+0i) (0+0i)]
	// [(0+18i) (0+0i) (0+0i) (0+0i)]
}

func ExampleTensorProduct() {
	m := matrix.New(
		[]complex128{1 + 1i, 2 + 2i},
		[]complex128{3 + 3i, 0 + 0i},
	)
	mM := matrix.TensorProduct(m, m)
	for _, r := range mM {
		fmt.Println(r)
	}
	// Output:
	// [(0+2i) (0+4i) (0+4i) (0+8i)]
	// [(0+6i) (0+0i) (0+12i) (0+0i)]
	// [(0+6i) (0+12i) (0+0i) (0+0i)]
	// [(0+18i) (0+0i) (0+0i) (0+0i)]
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
