package epsilon_test

import (
	"testing"

	"github.com/itsubaki/q/math/epsilon"
)

func TestE13(t *testing.T) {
	cases := []struct {
		in   []float64
		want float64
	}{
		{[]float64{}, 1e-13},
		{[]float64{1e-5}, 1e-5},
	}

	for _, c := range cases {
		got := epsilon.E13(c.in...)
		if got != c.want {
			t.Errorf("got=%v want=%v", got, c.want)
		}
	}
}

func TestE3(t *testing.T) {
	cases := []struct {
		in   []float64
		want float64
	}{
		{[]float64{}, 1e-3},
		{[]float64{1e-5}, 1e-5},
	}

	for _, c := range cases {
		got := epsilon.E3(c.in...)
		if got != c.want {
			t.Errorf("got=%v want=%v", got, c.want)
		}
	}
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
