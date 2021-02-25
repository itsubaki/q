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

func ExampleMatrix_Real() {
	m := matrix.New(
		[]complex128{complex(1, 1), complex(2, 3)},
		[]complex128{complex(4, 5), complex(6, 7)},
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
		[]complex128{complex(1, 1), complex(2, 3)},
		[]complex128{complex(4, 5), complex(6, 7)},
	)

	for _, r := range m.Imag() {
		fmt.Println(r)
	}

	// Output:
	// [1 3]
	// [5 7]
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
				[]complex128{0, complex(0, -1)},
				[]complex128{complex(0, 1), 0},
			),
			matrix.New(
				[]complex128{complex(0, 2), 0},
				[]complex128{0, complex(0, -2)},
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
				[]complex128{0, complex(0, -1)},
				[]complex128{complex(0, 1), 0},
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
