package number_test

import (
	"testing"

	"github.com/itsubaki/q/math/number"
)

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
