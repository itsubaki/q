package matrix_test

import (
	"fmt"
	"math/cmplx"
	"sync"
	"testing"

	"github.com/itsubaki/q/math/matrix"
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
					var c complex128
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

func BenchmarkDaggerN8(b *testing.B) {
	n := 8
	m := matrix.TensorProductN(matrix.New(
		[]complex128{0, 1},
		[]complex128{1, 0},
	), n)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Dagger()
	}
}

func BenchmarkDaggerConcurrencyN8(b *testing.B) {
	n := 8
	m := matrix.TensorProductN(matrix.New(
		[]complex128{0, 1},
		[]complex128{1, 0},
	), n)

	dagger := func(m matrix.Matrix) {
		p, q := m.Dimension()

		wg := sync.WaitGroup{}
		out := make(matrix.Matrix, p)
		for i := 0; i < p; i++ {
			wg.Add(1)

			go func(i int, out *matrix.Matrix) {
				defer wg.Done()

				v := make([]complex128, 0, q)
				for j := 0; j < q; j++ {
					v = append(v, cmplx.Conj(m[j][i]))
				}

				(*out)[i] = v
			}(i, &out)
		}

		wg.Wait()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dagger(m)
	}
}

func BenchmarkTensorProductN8(b *testing.B) {
	n := 8
	m := matrix.New(
		[]complex128{0, 1},
		[]complex128{1, 0},
	)

	for i := 0; i < b.N; i++ {
		matrix.TensorProductN(m, n)
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

	fmt.Println("xx:")
	for _, r := range x.Apply(x) {
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

	for _, r := range matrix.TensorProduct(x, x) {
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

	for _, r := range matrix.TensorProductN(x, 2) {
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

	for _, r := range matrix.Apply(x, x) {
		fmt.Println(r)
	}

	// Output:
	// [(1+0i) (0+0i)]
	// [(0+0i) (1+0i)]
}

func ExampleApplyN() {
	x := matrix.New(
		[]complex128{0, 1},
		[]complex128{1, 0},
	)

	for _, r := range matrix.ApplyN(x) {
		fmt.Println(r)
	}
	fmt.Println()

	for _, r := range matrix.ApplyN(x, 2) {
		fmt.Println(r)
	}

	// Output:
	// [(0+0i) (1+0i)]
	// [(1+0i) (0+0i)]
	//
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

	fmt.Println("xx:")
	for _, r := range x.TensorProduct(x) {
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
		in   matrix.Matrix
		want matrix.Matrix
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
		got := c.in.Apply(c.in.Inverse())
		if !got.Equals(c.want) {
			t.Fail()
		}
	}
}

func TestCommutator(t *testing.T) {
	cases := []struct {
		x, y, want matrix.Matrix
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
		if !matrix.Commutator(c.x, c.y).Equals(c.want) {
			t.Fail()
		}
	}
}

func TestAntiCommutator(t *testing.T) {
	cases := []struct {
		x, y, want matrix.Matrix
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
		if !matrix.AntiCommutator(c.x, c.y).Equals(c.want) {
			t.Fail()
		}
	}
}

func TestTrace(t *testing.T) {
	cases := []struct {
		in   matrix.Matrix
		want complex128
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
		if c.in.Trace() != c.want {
			t.Fail()
		}
	}
}

func TestDagger(t *testing.T) {
	cases := []struct {
		in matrix.Matrix
	}{
		{
			matrix.New(
				[]complex128{1 + 1i, 2 + 3i},
				[]complex128{4 + 5i, 6 + 7i},
			),
		},
	}

	for _, c := range cases {
		if !c.in.Dagger().Equals(c.in.Transpose().Conjugate()) {
			t.Fail()
		}
	}
}

func TestEquals(t *testing.T) {
	cases := []struct {
		m0, m1 matrix.Matrix
		want   bool
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
		if c.m0.Equals(c.m1) != c.want {
			t.Fail()
		}
	}
}

func TestIsSquare(t *testing.T) {
	cases := []struct {
		in   matrix.Matrix
		want bool
	}{
		{
			matrix.New(
				[]complex128{0, 1},
				[]complex128{1, 0},
			),
			true,
		},
		{
			matrix.New(
				[]complex128{1, 2, 3},
				[]complex128{4, 5, 6},
			),
			false,
		},
	}

	for _, c := range cases {
		if c.in.IsSquare() != c.want {
			t.Fail()
		}
	}
}

func TestIsHermite(t *testing.T) {
	cases := []struct {
		in   matrix.Matrix
		want bool
	}{
		{
			matrix.New(
				[]complex128{0, 1},
				[]complex128{1, 0},
			),
			true,
		},
		{
			matrix.New(
				[]complex128{0, -1i},
				[]complex128{1i, 0},
			),
			true,
		},
		{
			matrix.New(
				[]complex128{1, 0},
				[]complex128{0, -1},
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
		{
			matrix.New(
				[]complex128{1, 2, 3},
				[]complex128{4, 5, 6},
			),
			false,
		},
	}

	for _, c := range cases {
		if c.in.IsHermite() != c.want {
			t.Fail()
		}
	}
}

func TestIsUnitary(t *testing.T) {
	cases := []struct {
		in   matrix.Matrix
		want bool
	}{
		{
			matrix.New(
				[]complex128{0, 1},
				[]complex128{1, 0},
			),
			true,
		},
		{
			matrix.New(
				[]complex128{0, -1i},
				[]complex128{1i, 0},
			),
			true,
		},
		{
			matrix.New(
				[]complex128{1, 0},
				[]complex128{0, -1},
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
		{
			matrix.New(
				[]complex128{1, 2, 3},
				[]complex128{4, 5, 6},
			),
			false,
		},
	}

	for i, c := range cases {
		if c.in.IsUnitary() != c.want {
			t.Errorf("case[%v] is failed", i)
		}
	}
}

func TestTensorProductN(t *testing.T) {
	cases := []struct {
		in matrix.Matrix
	}{
		{
			matrix.New(
				[]complex128{0, 1},
				[]complex128{1, 0},
			),
		},
	}

	for _, c := range cases {
		if !matrix.TensorProductN(c.in).Equals(c.in) {
			t.Fail()
		}
	}
}
