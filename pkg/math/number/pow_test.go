package number

import (
	"testing"
)

func TestModExp2(t *testing.T) {
	cases := []struct {
		a, j, N int
		out     int64
	}{
		{7, 0, 15, 7},
		{7, 1, 15, 4},
		{7, 2, 15, 13},
		{7, 3, 15, 1},
		{7, 4, 15, 7},
		{7, 5, 15, 4},
		{7, 6, 15, 13},
		{7, 7, 15, 1},
		{7, 8, 15, 7},
		{7, 9, 15, 4},
		{7, 10, 15, 13},
		{7, 11, 15, 1},
		{7, 12, 15, 7},
		{7, 13, 15, 4},
		{7, 14, 15, 13},
	}

	for _, c := range cases {
		v := ModExp2(c.a, c.j, c.N)
		if v != c.out {
			t.Errorf("%d^(2^%d) mod %d = %d. actual=%d\n", c.a, c.j, c.N, c.out, v)
		}
	}
}