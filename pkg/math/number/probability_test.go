package number_test

import (
	"testing"

	"github.com/itsubaki/q/pkg/math/number"
)

func TestMax(t *testing.T) {
	cases := []struct {
		input []float64
		max   float64
	}{
		{[]float64{1.0, 2.0, 3.0}, 3.0},
		{[]float64{1.0, 2.0, -3.0}, 2.0},
	}

	for _, c := range cases {
		actual := number.Max(c.input)
		if actual != c.max {
			t.Errorf("actual=%v, expected=%v", actual, c.max)
		}
	}
}

func TestMin(t *testing.T) {
	cases := []struct {
		input []float64
		min   float64
	}{
		{[]float64{1.0, 2.0, 3.0}, 1.0},
		{[]float64{1.0, 2.0, -3.0}, -3.0},
	}

	for _, c := range cases {
		actual := number.Min(c.input)
		if actual != c.min {
			t.Errorf("actual=%v, expected=%v", actual, c.min)
		}
	}
}

func TestSum(t *testing.T) {
	cases := []struct {
		input []float64
		sum   float64
	}{
		{[]float64{1.0, 2.0, 3.0}, 6.0},
		{[]float64{1.0, 2.0, -3.0}, 0.0},
	}

	for _, c := range cases {
		actual := number.Sum(c.input)
		if actual != c.sum {
			t.Errorf("actual=%v, expected=%v", actual, c.sum)
		}
	}
}
