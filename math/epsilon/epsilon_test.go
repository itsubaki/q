package epsilon_test

import (
	"fmt"
	"testing"

	"github.com/itsubaki/q/math/epsilon"
)

func ExampleIsZero() {
	fmt.Println(epsilon.IsZero(complex(1e-13, 1e-13)))

	// Output:
	// true
}

func ExampleIsZeroF64() {
	fmt.Println(epsilon.IsZeroF64(1e-13))

	// Output:
	// true
}

func TestTol(t *testing.T) {
	cases := []struct {
		in   []float64
		want []float64
	}{
		{[]float64{}, []float64{1e-8, 1e-5}},
		{[]float64{1e-3}, []float64{1e-3, 1e-3}},
		{[]float64{1e-3, 1e-5}, []float64{1e-3, 1e-5}},
	}

	for _, c := range cases {
		got0, got1 := epsilon.Tol(c.in...)
		if got0 != c.want[0] || got1 != c.want[1] {
			t.Errorf("got=%v,%v want=%v,%v", got0, got1, c.want[0], c.want[1])
		}
	}
}

func TestIsClose(t *testing.T) {
	cases := []struct {
		a, b complex128
		tol  []float64
		want bool
	}{
		{1 + 1i, 1 + 1i, []float64{}, true},
		{1 + 1i, 1 + 1.00000001i, []float64{}, true},
		{1 + 1i, 1 + 1.1i, []float64{}, false},
		{1 + 1i, 1 + 1.0001i, []float64{1e-3}, true},
	}

	for _, c := range cases {
		got := epsilon.IsClose(c.a, c.b, c.tol...)
		if got != c.want {
			t.Errorf("got=%v want=%v", got, c.want)
		}
	}
}

func TestIsCloseF64(t *testing.T) {
	cases := []struct {
		a, b float64
		tol  []float64
		want bool
	}{
		{1.0, 1.0, []float64{}, true},
		{1.0, 1.00000001, []float64{}, true},
		{1.0, 1.1, []float64{}, false},
		{1.0, 1.0001, []float64{1e-3}, true},
	}

	for _, c := range cases {
		got := epsilon.IsCloseF64(c.a, c.b, c.tol...)
		if got != c.want {
			t.Errorf("got=%v want=%v", got, c.want)
		}
	}
}
