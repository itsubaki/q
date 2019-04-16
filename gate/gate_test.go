package gate_test

import (
	"fmt"
	"math/cmplx"
	"testing"

	"github.com/axamon/q"
	"github.com/axamon/q/gate"
	"github.com/axamon/q/matrix"
)

func ExampleH() {
	qsim := q.New()

	// generate qubits of |0>|1>
	q0 := qsim.Zero()
	q1 := qsim.One()

	// apply quantum circuit
	qsim.H(q0)
	qsim.H(q1)
	// estimate
	result0 := qsim.Estimate(q0).Probability()
	fmt.Printf("%.1f\n", result0)
	result1 := qsim.Estimate(q1).Probability()
	fmt.Printf("%.1f\n", result1)
	// Output:
	// [0.5 0.5]
	// [0.5 0.5]
}

func ExampleS() {
	qsim := q.New()

	// generate qubits of |0>|1>
	q0 := qsim.Zero()
	q1 := qsim.One()

	qsim.S(q0)
	qsim.S(q1)

	result0 := qsim.Estimate(q0).Probability()
	fmt.Printf("%.1f\n", result0)
	result1 := qsim.Estimate(q1).Probability()
	fmt.Printf("%.1f\n", result1)
	// Output:
	// [1.0 0.0]
	// [0.0 1.0]
}

func ExampleT() {
	qsim := q.New()

	// generate qubits of |0>|1>
	q0 := qsim.Zero()
	q1 := qsim.One()

	qsim.T(q0)
	qsim.T(q1)

	result0 := qsim.Estimate(q0).Probability()
	fmt.Printf("%.1f\n", result0)
	result1 := qsim.Estimate(q1).Probability()
	fmt.Printf("%.1f\n", result1)
	// Output:
	// [1.0 0.0]
	// [0.0 1.0]
}

func ExampleI() {
	m := gate.I(2)
	for _, r := range m {
		fmt.Println(r)
	}
	// Output:
	// [(1+0i) (0+0i) (0+0i) (0+0i)]
	// [(0+0i) (1+0i) (0+0i) (0+0i)]
	// [(0+0i) (0+0i) (1+0i) (0+0i)]
	// [(0+0i) (0+0i) (0+0i) (1+0i)]
}

func TestInverseU(t *testing.T) {
	m := gate.U(1.0, 1.1, 1.2, 1.3)

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
	expected := gate.New(
		[]complex128{1, 0, 0, 0, 0, 0, 0, 0},
		[]complex128{0, 1, 0, 0, 0, 0, 0, 0},
		[]complex128{0, 0, 1, 0, 0, 0, 0, 0},
		[]complex128{0, 0, 0, 1, 0, 0, 0, 0},
		[]complex128{0, 0, 0, 0, 1, 0, 0, 0},
		[]complex128{0, 0, 0, 0, 0, -1, 0, 0},
		[]complex128{0, 0, 0, 0, 0, 0, 1, 0},
		[]complex128{0, 0, 0, 0, 0, 0, 0, -1},
	)

	actual := gate.CZ(3, 0, 2)

	if !actual.Equals(expected) {
		t.Error(actual)
	}
}

func TestCNOT(t *testing.T) {
	expected := gate.New(
		[]complex128{1, 0, 0, 0, 0, 0, 0, 0},
		[]complex128{0, 1, 0, 0, 0, 0, 0, 0},
		[]complex128{0, 0, 1, 0, 0, 0, 0, 0},
		[]complex128{0, 0, 0, 1, 0, 0, 0, 0},
		[]complex128{0, 0, 0, 0, 0, 1, 0, 0},
		[]complex128{0, 0, 0, 0, 1, 0, 0, 0},
		[]complex128{0, 0, 0, 0, 0, 0, 0, 1},
		[]complex128{0, 0, 0, 0, 0, 0, 1, 0},
	)

	actual := gate.CNOT(3, 0, 2)

	if !actual.Equals(expected) {
		t.Error(actual)
	}
}

func TestControlledNot(t *testing.T) {
	g0 := matrix.TensorProduct(gate.I().Add(gate.Z()), gate.I())
	g1 := matrix.TensorProduct(gate.I().Sub(gate.Z()), gate.X())
	CN := g0.Add(g1).Mul(0.5)

	if !gate.ControlledNot(2, []int{0}, 1).Equals(CN) {
		t.Error(CN)
	}
}

func TestToffoli(t *testing.T) {
	g := make([]matrix.Matrix, 13)

	g[0] = matrix.TensorProduct(gate.I(2), gate.H())
	g[1] = matrix.TensorProduct(gate.I(), gate.CNOT(2, 0, 1))
	g[2] = matrix.TensorProduct(gate.I(2), gate.T().Dagger())
	g[3] = gate.CNOT(3, 0, 2)
	g[4] = matrix.TensorProduct(gate.I(2), gate.T())
	g[5] = matrix.TensorProduct(gate.I(), gate.CNOT(2, 0, 1))
	g[6] = matrix.TensorProduct(gate.I(2), gate.T().Dagger())
	g[7] = gate.CNOT(3, 0, 2)
	g[8] = matrix.TensorProduct(gate.I(), gate.T().Dagger(), gate.T())
	g[9] = matrix.TensorProduct(gate.CNOT(2, 0, 1), gate.H())
	g[10] = matrix.TensorProduct(gate.I(), gate.T().Dagger(), gate.I())
	g[11] = matrix.TensorProduct(gate.CNOT(2, 0, 1), gate.I())
	g[12] = matrix.TensorProduct(gate.T(), gate.S(), gate.I())

	expected := gate.I(3)
	for _, gate := range g {
		expected = expected.Apply(gate)
	}

	actual := gate.Toffoli()
	if !actual.Equals(expected, 1e-13) {
		t.Error(actual)
	}
}

func TestIsHermite(t *testing.T) {
	if !gate.H().IsHermite() {
		t.Error(gate.H())
	}

	if !gate.X().IsHermite() {
		t.Error(gate.X())
	}

	if !gate.Y().IsHermite() {
		t.Error(gate.Y())
	}

	if !gate.Z().IsHermite() {
		t.Error(gate.Z())
	}

}

func TestIsUnitary(t *testing.T) {
	if !gate.H().IsUnitary(1e-13) {
		t.Error(gate.H())
	}

	if !gate.X().IsUnitary() {
		t.Error(gate.X())
	}

	if !gate.Y().IsUnitary() {
		t.Error(gate.Y())
	}

	if !gate.Z().IsUnitary() {
		t.Error(gate.Z())
	}

	u := gate.U(1, 2, 3, 4)
	if !u.IsUnitary(1e-13) {
		t.Error(u)
	}
}

func TestTrace(t *testing.T) {
	trA := gate.I().Trace()
	if trA != complex(2, 0) {
		t.Error(trA)
	}

	trH := gate.H().Trace()
	if trH != complex(0, 0) {
		t.Error(trH)
	}
}

func TestTensorProductProductXY(t *testing.T) {
	x := gate.X()
	y := gate.Y()

	m, n := x.Dimension()
	tmp := []matrix.Matrix{}
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
	xx := gate.X().TensorProduct(gate.X())
	y := gate.Y()

	m, n := xx.Dimension()
	tmp := []matrix.Matrix{}
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
