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
		p, q := m.Dimension()

		wg := sync.WaitGroup{}
		data := make([]complex128, p)
		for i := range p {
			wg.Add(1)

			go func(i int) {
				defer wg.Done()

				for j := range q {
					data[i] += m.At(i, j) * v.Data[j]
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

func ExampleVector_IsOrthogonal() {
	v0 := vector.New(1, 0)
	v1 := vector.New(0, 1)
	o := v0.IsOrthogonal(v1)
	fmt.Println(o)

	// Output:
	// true
}

func ExampleVector_Norm() {
	v := vector.New(1, 2)
	n := v.Norm()
	fmt.Printf("%.4f", n)

	// Output:
	// (2.2361+0.0000i)
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

func TestVector(t *testing.T) {
	cases := []struct {
		v0, v1       *vector.Vector
		innerProduct complex128
		isOrthogonal bool
		isUnit       bool
		norm         complex128
	}{
		{vector.New(1, 0), vector.New(1, 0), 1, false, true, 1},
		{vector.New(0, 1), vector.New(0, 1), 1, false, true, 1},
		{vector.New(1, 0), vector.New(0, 1), 0, true, true, 1},
	}

	for _, c := range cases {
		if c.v0.InnerProduct(c.v1) != c.innerProduct {
			t.Errorf("%v", c.v0.InnerProduct(c.v1))
		}

		if c.v0.IsOrthogonal(c.v1) != c.isOrthogonal {
			t.Errorf("%v", c.v0.IsOrthogonal(c.v1))
		}

		if c.v0.IsUnit() != c.isUnit {
			t.Errorf("%v", c.v0.IsUnit())
		}

		if c.v0.Norm() != c.norm {
			t.Errorf("%v", c.v0.Norm())
		}
	}
}

func TestAdd(t *testing.T) {
	cases := []struct {
		v0, v1 *vector.Vector
		want   *vector.Vector
	}{
		{vector.New(1, 2, 3, 4, 5), vector.New(6, 7, 8, 9, 10), vector.New(7, 9, 11, 13, 15)},
	}

	for _, c := range cases {
		if !c.v0.Add(c.v1).Equals(c.want) {
			t.Fail()
		}
	}
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
		if !c.v.Mul(c.c).Equals(c.want) {
			t.Fail()
		}
	}
}

func TestClone(t *testing.T) {
	in := vector.New(1, 2, 3, 4, 5, 6)
	got := in.Clone()

	if !in.Equals(got) {
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
		if !vector.TensorProductN(c.in).Equals(c.in) {
			t.Fail()
		}
	}
}

func TestEquals(t *testing.T) {
	cases := []struct {
		v0, v1 *vector.Vector
		want   bool
	}{
		{vector.New(1, 2), vector.New(1, 2), true},
		{vector.New(1, 2), vector.New(3, 4), false},
		{vector.New(1, 2, 3, 4, 5, 6), vector.New(1, 2), false},
	}

	for _, c := range cases {
		if c.v0.Equals(c.v1) != c.want {
			t.Fail()
		}
	}
}
