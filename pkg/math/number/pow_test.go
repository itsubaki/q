package number_test

import (
	"testing"

	"github.com/itsubaki/q/pkg/math/number"
)

func TestModExp2(t *testing.T) {
	cases := []struct {
		a, j, N int
		out     int
	}{
		{7, 0, 15, 7},
		{7, 1, 15, 4},
		{7, 2, 15, 1},
		{7, 3, 15, 1},
		{7, 4, 15, 1},
		{7, 5, 15, 1},
		{7, 6, 15, 1},
		{7, 7, 15, 1},
		{7, 8, 15, 1},
		{7, 9, 15, 1},
		{7, 10, 15, 1},
		{7, 11, 15, 1},
		{7, 12, 15, 1},
		{7, 13, 15, 1},
		{7, 14, 15, 1},
	}

	for _, c := range cases {
		v := number.ModExp2(c.a, c.j, c.N)
		if v != c.out {
			t.Errorf("%d^(2^%d) mod %d = %d. actual=%d\n", c.a, c.j, c.N, c.out, v)
		}
	}
}

func TestBaseExp(t *testing.T) {
	cases := []struct {
		N, a, b int
		ok      bool
	}{
		{9, 3, 2, true},
		{25, 5, 2, true},
		{49, 7, 2, true},
		{125, 5, 3, true},
		{10, 0, 0, false},
	}

	for _, c := range cases {
		a, b, ok := number.BaseExp(c.N)
		if a != c.a || b != c.b || ok != c.ok {
			t.Errorf("N=%v, a=%v, b=%v, ok=%v", c.N, a, b, ok)
		}
	}
}
