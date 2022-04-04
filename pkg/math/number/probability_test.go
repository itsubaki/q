package number_test

import (
	"testing"

	"github.com/itsubaki/q/pkg/math/number"
)

func TestMax(t *testing.T) {
	cases := []struct {
		in   []float64
		want float64
	}{
		{[]float64{1.0, 2.0, 3.0}, 3.0},
		{[]float64{1.0, 2.0, -3.0}, 2.0},
	}

	for _, c := range cases {
		got := number.Max(c.in)
		if got != c.want {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}

func TestMin(t *testing.T) {
	cases := []struct {
		in   []float64
		want float64
	}{
		{[]float64{1.0, 2.0, 3.0}, 1.0},
		{[]float64{1.0, 2.0, -3.0}, -3.0},
	}

	for _, c := range cases {
		got := number.Min(c.in)
		if got != c.want {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}

func TestSum(t *testing.T) {
	cases := []struct {
		in   []float64
		want float64
	}{
		{[]float64{1.0, 2.0, 3.0}, 6.0},
		{[]float64{1.0, 2.0, -3.0}, 0.0},
	}

	for _, c := range cases {
		got := number.Sum(c.in)
		if got != c.want {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}
