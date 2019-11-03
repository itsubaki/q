package gate

import (
	"fmt"
	"math/cmplx"
	"testing"

	"github.com/itsubaki/q/pkg/math/matrix"
)

func TestInverseU(t *testing.T) {
	m := U(1.0, 1.1, 1.2, 1.3)

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

func TestCZ(t *testing.T) {
	expected := New(
		[]complex128{1, 0, 0, 0, 0, 0, 0, 0},
		[]complex128{0, 1, 0, 0, 0, 0, 0, 0},
		[]complex128{0, 0, 1, 0, 0, 0, 0, 0},
		[]complex128{0, 0, 0, 1, 0, 0, 0, 0},
		[]complex128{0, 0, 0, 0, 1, 0, 0, 0},
		[]complex128{0, 0, 0, 0, 0, -1, 0, 0},
		[]complex128{0, 0, 0, 0, 0, 0, 1, 0},
		[]complex128{0, 0, 0, 0, 0, 0, 0, -1},
	)

	actual := CZ(3, 0, 2)

	if !actual.Equals(expected) {
		t.Error(actual)
	}
}

func TestCNOT(t *testing.T) {
	expected := New(
		[]complex128{1, 0, 0, 0, 0, 0, 0, 0},
		[]complex128{0, 1, 0, 0, 0, 0, 0, 0},
		[]complex128{0, 0, 1, 0, 0, 0, 0, 0},
		[]complex128{0, 0, 0, 1, 0, 0, 0, 0},
		[]complex128{0, 0, 0, 0, 0, 1, 0, 0},
		[]complex128{0, 0, 0, 0, 1, 0, 0, 0},
		[]complex128{0, 0, 0, 0, 0, 0, 0, 1},
		[]complex128{0, 0, 0, 0, 0, 0, 1, 0},
	)

	actual := CNOT(3, 0, 2)

	if !actual.Equals(expected) {
		t.Error(actual)
	}
}

func TestControlledNot(t *testing.T) {
	g0 := matrix.TensorProduct(I().Add(Z()), I())
	g1 := matrix.TensorProduct(I().Sub(Z()), X())
	CN := g0.Add(g1).Mul(0.5)

	if !ControlledNot(2, []int{0}, 1).Equals(CN) {
		t.Error(CN)
	}
}

func TestToffoli(t *testing.T) {
	g := make([]matrix.Matrix, 13)

	g[0] = matrix.TensorProduct(I(2), H())
	g[1] = matrix.TensorProduct(I(), CNOT(2, 0, 1))
	g[2] = matrix.TensorProduct(I(2), T().Dagger())
	g[3] = CNOT(3, 0, 2)
	g[4] = matrix.TensorProduct(I(2), T())
	g[5] = matrix.TensorProduct(I(), CNOT(2, 0, 1))
	g[6] = matrix.TensorProduct(I(2), T().Dagger())
	g[7] = CNOT(3, 0, 2)
	g[8] = matrix.TensorProduct(I(), T().Dagger(), T())
	g[9] = matrix.TensorProduct(CNOT(2, 0, 1), H())
	g[10] = matrix.TensorProduct(I(), T().Dagger(), I())
	g[11] = matrix.TensorProduct(CNOT(2, 0, 1), I())
	g[12] = matrix.TensorProduct(T(), S(), I())

	expected := I(3)
	for _, gate := range g {
		expected = expected.Apply(gate)
	}

	actual := Toffoli()
	if !actual.Equals(expected, 1e-13) {
		t.Error(actual)
	}
}

func TestIsHermite(t *testing.T) {
	if !H().IsHermite() {
		t.Error(H())
	}

	if !X().IsHermite() {
		t.Error(X())
	}

	if !Y().IsHermite() {
		t.Error(Y())
	}

	if !Z().IsHermite() {
		t.Error(Z())
	}

}

func TestIsUnitary(t *testing.T) {
	if !H().IsUnitary(1e-13) {
		t.Error(H())
	}

	if !X().IsUnitary() {
		t.Error(X())
	}

	if !Y().IsUnitary() {
		t.Error(Y())
	}

	if !Z().IsUnitary() {
		t.Error(Z())
	}

	u := U(1, 2, 3, 4)
	if !u.IsUnitary(1e-13) {
		t.Error(u)
	}
}

func TestTrace(t *testing.T) {
	trA := I().Trace()
	if trA != complex(2, 0) {
		t.Error(trA)
	}

	trH := H().Trace()
	if trH != complex(0, 0) {
		t.Error(trH)
	}
}

func TestTensorProductProductXY(t *testing.T) {
	x := X()
	y := Y()

	m, n := x.Dimension()
	tmp := make([]matrix.Matrix, 0)
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			tmp = append(tmp, y.Mul(x[i][j]))
		}
	}

	fmt.Printf("%v %v %v %v\n", tmp[0][0][0], tmp[0][0][1], tmp[1][0][0], tmp[1][0][1])
	fmt.Printf("%v %v %v %v\n", tmp[0][1][0], tmp[0][1][1], tmp[1][1][0], tmp[1][1][1])
	fmt.Printf("%v %v %v %v\n", tmp[2][0][0], tmp[2][0][1], tmp[3][0][0], tmp[3][0][1])
	fmt.Printf("%v %v %v %v\n", tmp[2][1][0], tmp[2][1][1], tmp[3][1][0], tmp[3][1][1])
	fmt.Println()
}

func TestTensorProductProductXXY(t *testing.T) {
	xx := X().TensorProduct(X())
	y := Y()

	m, n := xx.Dimension()
	tmp := make([]matrix.Matrix, 0)
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			tmp = append(tmp, y.Mul(xx[i][j]))
		}
	}

	fmt.Printf("%v %v %v %v %v %v %v %v\n", tmp[0][0][0], tmp[0][0][1], tmp[1][0][0], tmp[1][0][1], tmp[2][0][0], tmp[2][0][1], tmp[3][0][0], tmp[3][0][1])
	fmt.Printf("%v %v %v %v %v %v %v %v\n", tmp[0][1][0], tmp[0][1][1], tmp[1][1][0], tmp[1][1][1], tmp[2][1][0], tmp[2][1][1], tmp[3][1][0], tmp[3][1][1])
	fmt.Printf("%v %v %v %v %v %v %v %v\n", tmp[4][0][0], tmp[4][0][1], tmp[5][0][0], tmp[5][0][1], tmp[6][0][0], tmp[6][0][1], tmp[7][0][0], tmp[7][0][1])
	fmt.Printf("%v %v %v %v %v %v %v %v\n", tmp[4][1][0], tmp[4][1][1], tmp[5][1][0], tmp[5][1][1], tmp[6][1][0], tmp[6][1][1], tmp[7][1][0], tmp[7][1][1])
	fmt.Printf("%v %v %v %v %v %v %v %v\n", tmp[8][0][0], tmp[8][0][1], tmp[9][0][0], tmp[9][0][1], tmp[10][0][0], tmp[10][0][1], tmp[11][0][0], tmp[11][0][1])
	fmt.Printf("%v %v %v %v %v %v %v %v\n", tmp[8][1][0], tmp[8][1][1], tmp[9][1][0], tmp[9][1][1], tmp[10][1][0], tmp[10][1][1], tmp[11][1][0], tmp[11][1][1])
	fmt.Printf("%v %v %v %v %v %v %v %v\n", tmp[12][0][0], tmp[12][0][1], tmp[13][0][0], tmp[13][0][1], tmp[14][0][0], tmp[14][0][1], tmp[15][0][0], tmp[15][0][1])
	fmt.Printf("%v %v %v %v %v %v %v %v\n", tmp[12][1][0], tmp[12][1][1], tmp[13][1][0], tmp[13][1][1], tmp[14][1][0], tmp[14][1][1], tmp[15][1][0], tmp[15][1][1])
	fmt.Println()
}
