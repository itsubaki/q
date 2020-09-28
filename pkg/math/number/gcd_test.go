package number_test

import (
	"testing"

	"github.com/itsubaki/q/pkg/math/number"
)

func TestGCD(t *testing.T) {
	cases := []struct {
		a, b, c int
	}{
		{15, 2, 1},
		{15, 4, 1},
		{15, 7, 1},
		{15, 11, 1},
		{15, 13, 1},
		{15, 14, 1},
	}

	for _, c := range cases {
		gcd := number.GCD(c.a, c.b)
		if gcd != c.c {
			t.Errorf("gcd=%d", gcd)
		}
	}
}
