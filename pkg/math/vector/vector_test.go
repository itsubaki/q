package vector

import (
	"fmt"
	"testing"
)

func ExampleVector_OuterProduct() {
	v := New(1, 0)
	op := v.OuterProduct(v)
	fmt.Println(op)

	// Output:
	// [[(1+0i) (0+0i)] [(0+0i) (0+0i)]]
}

func TestVector(t *testing.T) {
	cases := []struct {
		v0, v1       Vector
		innerProduct complex128
		isOrthogonal bool
		isUnit       bool
		norm         complex128
	}{
		{New(1, 0), New(1, 0), complex(1, 0), false, true, complex(1, 0)},
		{New(0, 1), New(0, 1), complex(1, 0), false, true, complex(1, 0)},
		{New(1, 0), New(0, 1), complex(0, 0), true, true, complex(1, 0)},
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
