package vector_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/itsubaki/q/math/matrix"
	"github.com/itsubaki/q/math/vector"
)

func BenchmarkApplyN12(b *testing.B) {
	n := 12
	v := vector.TensorProductN(vector.New(1, 2), n)
	x := matrix.TensorProductN(matrix.New(
		[]complex128{0, 1},
		[]complex128{1, 0},
	), n)

	b.ResetTimer()
	for range b.N {
		v.Apply(x)
	}
}

func BenchmarkApplyConcurrencyN12(b *testing.B) {
	apply := func(v *vector.Vector, m *matrix.Matrix) *vector.Vector {
		p, q := m.Dim()

		wg := sync.WaitGroup{}
		data := make([]complex128, p)
		for i := range p {
			wg.Add(1)

			go func(i int) {
				defer wg.Done()

				row := m.Row(i)
				for j := range q {
					data[i] += row[j] * v.Data[j]
				}
			}(i)
		}

		wg.Wait()
		return &vector.Vector{
			Data: data,
		}
	}

	n := 12
	v := vector.TensorProductN(vector.New(1, 2), n)
	x := matrix.TensorProductN(matrix.New(
		[]complex128{0, 1},
		[]complex128{1, 0},
	), n)

	b.ResetTimer()
	for range b.N {
		apply(v, x)
	}
}

func ExampleZero() {
	v := vector.Zero(3)
	fmt.Println(v)

	// Output:
	// [(0+0i) (0+0i) (0+0i)]
}

func ExampleNew() {
	v := vector.New(1, 0)
	fmt.Println(v)

	// Output:
	// [(1+0i) (0+0i)]
}

func ExampleTensorProductN() {
	v := vector.New(1, 0)
	vv := vector.TensorProductN(v, 2)
	fmt.Println(vv)

	// Output:
	// [(1+0i) (0+0i) (0+0i) (0+0i)]
}

func ExampleVector_Apply() {
	v := vector.New(1, 2)
	fmt.Println(v)

	m := matrix.New(
		[]complex128{1, 2},
		[]complex128{1, 3},
	)

	vv := v.Apply(m)
	fmt.Println(vv)

	// Output:
	// [(1+0i) (2+0i)]
	// [(5+0i) (7+0i)]
}

func ExampleVector_InnerProduct() {
	v := vector.New(1, 0)
	vv := v.InnerProduct(v)
	fmt.Println(vv)

	// Output:
	// (1+0i)
}

func ExampleVector_OuterProduct() {
	v := vector.New(1, 0)
	vv := v.OuterProduct(v)

	for _, r := range vv.Seq2() {
		fmt.Println(r)
	}

	// Output:
	// [(1+0i) (0+0i)]
	// [(0+0i) (0+0i)]
}

func ExampleVector_TensorProduct() {
	v := vector.New(1, 0)
	vv := v.TensorProduct(v)
	fmt.Println(vv)

	// Output:
	// [(1+0i) (0+0i) (0+0i) (0+0i)]
}

func ExampleVector_Real() {
	v := vector.New(1+2i, 3+4i)
	for _, r := range v.Real() {
		fmt.Println(r)
	}

	// Output:
	// 1
	// 3
}
func ExampleVector_Imag() {
	v := vector.New(1+2i, 3+4i)
	for _, r := range v.Imag() {
		fmt.Println(r)
	}

	// Output:
	// 2
	// 4
}

func TestMul(t *testing.T) {
	cases := []struct {
		v    *vector.Vector
		c    complex128
		want *vector.Vector
	}{
		{vector.New(1, 2, 3, 4, 5), 3, vector.New(3, 6, 9, 12, 15)},
	}

	for _, c := range cases {
		if !c.v.Mul(c.c).Equal(c.want) {
			t.Fail()
		}
	}
}

func TestClone(t *testing.T) {
	in := vector.New(1, 2, 3, 4, 5, 6)
	got := in.Clone()

	if !in.Equal(got) {
		t.Fail()
	}
}

func TestTensorProductN(t *testing.T) {
	cases := []struct {
		in *vector.Vector
	}{
		{vector.New(1, 2, 3, 4, 5, 6)},
	}

	for _, c := range cases {
		if !vector.TensorProductN(c.in).Equal(c.in) {
			t.Fail()
		}
	}
}

func TestEqual(t *testing.T) {
	cases := []struct {
		v0, v1 *vector.Vector
		want   bool
	}{
		{vector.New(1, 2), vector.New(1, 2), true},
		{vector.New(1, 2), vector.New(3, 4), false},
		{vector.New(1, 2, 3, 4, 5, 6), vector.New(1, 2), false},
	}

	for _, c := range cases {
		if c.v0.Equal(c.v1) != c.want {
			t.Fail()
		}
	}
}
