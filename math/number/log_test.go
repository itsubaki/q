package number_test

import (
	"testing"

	"github.com/itsubaki/q/math/number"
)

func TestLog2(t *testing.T) {
	cases := []struct {
		n    int
		want int
	}{
		{1, 0},
		{2, 1},
		{4, 2},
		{8, 3},
		{16, 4},
		{32, 5},
		{64, 6},
		{128, 7},
		{256, 8},
		{512, 9},
		{1024, 10},
		{2048, 11},
		{4096, 12},
		{8192, 13},
	}

	for _, c := range cases {
		got := number.Must(number.Log2(c.n))
		if got != c.want {
			t.Errorf("got=%v want=%d", got, c.want)
		}
	}
}
