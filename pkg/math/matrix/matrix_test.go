package matrix_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/itsubaki/q/pkg/math/matrix"
)

func BenchmarkApplyN8(b *testing.B) {
	n := 8
	x := matrix.TensorProductN(matrix.New(
		[]complex128{0, 1},
		[]complex128{1, 0},
	), n)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x.Apply(x)
	}
}

func BenchmarkApplyConcurrencyN8(b *testing.B) {
	apply := func(n, m matrix.Matrix) matrix.Matrix {
		p, _ := m.Dimension()
		a, b := n.Dimension()

		wg := sync.WaitGroup{}
		out := make(matrix.Matrix, a)
		for i := 0; i < a; i++ {
			wg.Add(1)
			go func(i int, out *matrix.Matrix) {
				defer wg.Done()

				v := make([]complex128, b)
				for j := 0; j < b; j++ {
					c := complex(0, 0)
					for k := 0; k < p; k++ {
						c = c + n[i][k]*m[k][j]
					}

					v = append(v, c)
				}

				(*out)[i] = v
			}(i, &out)
		}

		wg.Wait()
		return out
	}

	n := 8
	x := matrix.TensorProductN(matrix.New(
		[]complex128{0, 1},
		[]complex128{1, 0},
	), n)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		apply(x, x)
	}
}

func BenchmarkTensorProductN6(b *testing.B) {
	n := 6
	x := matrix.TensorProductN(matrix.New(
		[]complex128{0, 1},
		[]complex128{1, 0},
	), n)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		matrix.TensorProduct(x, x)
	}
}

func BenchmarkTensorProductConcurrencyN6(b *testing.B) {
	tensorproduct := func(n, m matrix.Matrix) matrix.Matrix {
		p, q := m.Dimension()
		a, b := n.Dimension()

		wg := sync.WaitGroup{}
		out := make(matrix.Matrix, a)
		for i := 0; i < p; i++ {
			wg.Add(1)
			go func(i int, out *matrix.Matrix) {
				defer wg.Done()

				for k := 0; k < a; k++ {
					r := make([]complex128, 0)
					for j := 0; j < q; j++ {
						for l := 0; l < b; l++ {
							r = append(r, m[i][j]*n[k][l])
						}
					}

					(*out)[i] = r
				}
			}(i, &out)
		}

		wg.Wait()
		return out
	}

	n := 6
	x := matrix.TensorProductN(matrix.New(
		[]complex128{0, 1},
		[]complex128{1, 0},
	), n)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tensorproduct(x, x)
	}
}

func ExampleZero() {
	fmt.Println(matrix.Zero(0))
	fmt.Println(matrix.Zero(1))
	fmt.Println(matrix.Zero(2))

	// Output:
	// []
	// [[(0+0i)]]
	// [[(0+0i) (0+0i)] [(0+0i) (0+0i)]]
}

func ExampleMatrix_Real() {
	m := matrix.New(
		[]complex128{1 + 1i, 2 + 3i},
		[]complex128{4 + 5i, 6 + 7i},
	)

	for _, r := range m.Real() {
		fmt.Println(r)
	}

	// Output:
	// [1 2]
	// [4 6]
}

func ExampleMatrix_Imag() {
	m := matrix.New(
		[]complex128{1 + 1i, 2 + 3i},
		[]complex128{4 + 5i, 6 + 7i},
	)

	for _, r := range m.Imag() {
		fmt.Println(r)
	}

	// Output:
	// [1 3]
	// [5 7]
}

func ExampleMatrix_Mul() {
	m := matrix.New(
		[]complex128{0, 1 + 1i},
		[]complex128{1 + 1i, 0},
	)

	for _, r := range m.Mul(3i) {
		fmt.Println(r)
	}

	// Output:
	// [(0+0i) (-3+3i)]
	// [(-3+3i) (0+0i)]
}

func ExampleMatrix_Apply() {
	x := matrix.New(
		[]complex128{0, 1},
		[]complex128{1, 0},
	)

	fmt.Println("x:")
	for _, r := range x {
		fmt.Println(r)
	}

	xx := x.Apply(x)
	fmt.Println("xx:")
	for _, r := range xx {
		fmt.Println(r)
	}

	// Output:
	// x:
	// [(0+0i) (1+0i)]
	// [(1+0i) (0+0i)]
	// xx:
	// [(1+0i) (0+0i)]
	// [(0+0i) (1+0i)]
}

func ExampleTensorProduct() {
	x := matrix.New(
		[]complex128{0, 1},
		[]complex128{1, 0},
	)

	xx := matrix.TensorProduct(x, x)
	for _, r := range xx {
		fmt.Println(r)
	}

	// Output:
	// [(0+0i) (0+0i) (0+0i) (1+0i)]
	// [(0+0i) (0+0i) (1+0i) (0+0i)]
	// [(0+0i) (1+0i) (0+0i) (0+0i)]
	// [(1+0i) (0+0i) (0+0i) (0+0i)]
}

func ExampleTensorProductN() {
	x := matrix.New(
		[]complex128{0, 1},
		[]complex128{1, 0},
	)

	xx := matrix.TensorProductN(x, 2)
	for _, r := range xx {
		fmt.Println(r)
	}

	// Output:
	// [(0+0i) (0+0i) (0+0i) (1+0i)]
	// [(0+0i) (0+0i) (1+0i) (0+0i)]
	// [(0+0i) (1+0i) (0+0i) (0+0i)]
	// [(1+0i) (0+0i) (0+0i) (0+0i)]
}

func ExampleApply() {
	x := matrix.New(
		[]complex128{0, 1},
		[]complex128{1, 0},
	)

	xx := matrix.Apply(x, x)
	for _, r := range xx {
		fmt.Println(r)
	}

	// Output:
	// [(1+0i) (0+0i)]
	// [(0+0i) (1+0i)]
}

func ExampleMatrix_TensorProduct() {
	x := matrix.New(
		[]complex128{0, 1},
		[]complex128{1, 0},
	)
	fmt.Println("x:")
	for _, r := range x {
		fmt.Println(r)
	}

	xx := x.TensorProduct(x)
	fmt.Println("xx:")
	for _, r := range xx {
		fmt.Println(r)
	}

	// Output:
	// x:
	// [(0+0i) (1+0i)]
	// [(1+0i) (0+0i)]
	// xx:
	// [(0+0i) (0+0i) (0+0i) (1+0i)]
	// [(0+0i) (0+0i) (1+0i) (0+0i)]
	// [(0+0i) (1+0i) (0+0i) (0+0i)]
	// [(1+0i) (0+0i) (0+0i) (0+0i)]
}

func TestInverse(t *testing.T) {
	cases := []struct {
		m, e matrix.Matrix
	}{
		{
			matrix.New(
				[]complex128{1, 2, 0, -1},
				[]complex128{-1, 1, 2, 0},
				[]complex128{2, 0, 1, 1},
				[]complex128{1, -2, -1, 1},
			),
			matrix.New(
				[]complex128{1, 0, 0, 0},
				[]complex128{0, 1, 0, 0},
				[]complex128{0, 0, 1, 0},
				[]complex128{0, 0, 0, 1},
			),
		},
		{
			matrix.New(
				[]complex128{0, 1},
				[]complex128{1, 0},
			),
			matrix.New(
				[]complex128{1, 0},
				[]complex128{0, 1},
			),
		},
	}

	for _, c := range cases {
		inv := c.m.Inverse()
		mmi := c.m.Apply(inv)
		if !mmi.Equals(c.e) {
			t.Fail()
		}
	}
}

func TestCommutator(t *testing.T) {
	cases := []struct {
		x, y, e matrix.Matrix
	}{
		{
			matrix.New(
				[]complex128{0, 1},
				[]complex128{1, 0},
			),
			matrix.New(
				[]complex128{0, -1i},
				[]complex128{1i, 0},
			),
			matrix.New(
				[]complex128{2i, 0},
				[]complex128{0, -2i},
			),
		},
	}

	for _, c := range cases {
		if !matrix.Commutator(c.x, c.y).Equals(c.e) {
			t.Fail()
		}
	}
}

func TestAntiCommutator(t *testing.T) {
	cases := []struct {
		x, y, e matrix.Matrix
	}{
		{
			matrix.New(
				[]complex128{0, 1},
				[]complex128{1, 0},
			),
			matrix.New(
				[]complex128{0, -1i},
				[]complex128{1i, 0},
			),
			matrix.New(
				[]complex128{0, 0},
				[]complex128{0, 0},
			),
		},
	}

	for _, c := range cases {
		if !matrix.AntiCommutator(c.x, c.y).Equals(c.e) {
			t.Fail()
		}
	}
}

func TestTrace(t *testing.T) {
	cases := []struct {
		m matrix.Matrix
		t complex128
	}{
		{
			matrix.New(
				[]complex128{0, 1},
				[]complex128{1, 0},
			),
			0,
		},
		{
			matrix.New(
				[]complex128{2, 1},
				[]complex128{1, 4},
			),
			6,
		},
		{
			matrix.New(
				[]complex128{1 + 1i, 2 + 3i},
				[]complex128{4 + 5i, 6 + 7i},
			),
			7 + 8i,
		},
	}

	for _, c := range cases {
		if c.m.Trace() != c.t {
			t.Fail()
		}
	}
}

func TestDagger(t *testing.T) {
	cases := []struct {
		m matrix.Matrix
	}{
		{
			matrix.New(
				[]complex128{1 + 1i, 2 + 3i},
				[]complex128{4 + 5i, 6 + 7i},
			),
		},
	}

	for _, c := range cases {
		if !c.m.Transpose().Conjugate().Equals(c.m.Dagger()) {
			t.Fail()
		}
	}
}

func TestEquals(t *testing.T) {
	cases := []struct {
		m0, m1 matrix.Matrix
		yes    bool
	}{
		{
			matrix.New(
				[]complex128{1 + 1i, 2 + 3i},
				[]complex128{4 + 5i, 6 + 7i},
			),
			matrix.New(
				[]complex128{1 + 1i, 2 + 3i},
				[]complex128{4 + 5i, 6 + 7i},
			),
			true,
		},
		{
			matrix.New(
				[]complex128{1 + 1i, 2 + 3i},
				[]complex128{4 + 5i, 6 + 7i},
			),
			matrix.New(
				[]complex128{10 + 10i, 20 + 30i},
				[]complex128{40 + 50i, 60 + 70i},
			),
			false,
		},
		{
			matrix.New(
				[]complex128{1 + 1i, 2 + 3i},
			),
			matrix.New(
				[]complex128{1 + 1i, 2 + 3i},
				[]complex128{4 + 5i, 6 + 7i},
			),
			false,
		},
		{
			matrix.New(
				[]complex128{1 + 1i},
				[]complex128{4 + 5i},
			),
			matrix.New(
				[]complex128{1 + 1i, 2 + 3i},
				[]complex128{4 + 5i, 6 + 7i},
			),
			false,
		},
	}

	for _, c := range cases {
		if c.m0.Equals(c.m1) != c.yes {
			t.Fail()
		}
	}
}

func TestIsHermite(t *testing.T) {
	cases := []struct {
		m matrix.Matrix
	}{
		{
			matrix.New(
				[]complex128{0, 1},
				[]complex128{1, 0},
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
				[]complex128{1, 0},
				[]complex128{0, -1},
			),
		},
	}

	for _, c := range cases {
		if !c.m.IsHermite() {
			t.Error(c.m)
		}
	}
}

func TestIsUnitary(t *testing.T) {
	cases := []struct {
		m matrix.Matrix
	}{
		{
			matrix.New(
				[]complex128{0, 1},
				[]complex128{1, 0},
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
				[]complex128{1, 0},
				[]complex128{0, -1},
			),
		},
	}

	for _, c := range cases {
		if !c.m.IsUnitary() {
			t.Error(c.m)
		}
	}
}
