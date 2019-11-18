package number

import (
	"testing"
)

func TestFraction(t *testing.T) {
	cases := []struct {
		float       float64
		fraction    []int
		numerator   int
		denominator int
	}{
		{0.42857, []int{0, 2, 2, 1}, 3, 7},
		{1.0 / 16.0, []int{0, 16}, 1, 16},
		{4.0 / 16.0, []int{0, 4}, 1, 4},
		{7.0 / 16.0, []int{0, 2, 3, 1, 1}, 7, 16},
		{13.0 / 16.0, []int{0, 1, 4, 3}, 13, 16},
	}

	for _, c := range cases {
		f, n, d := Fraction(c.float, 1e-3)

		if n != c.numerator {
			t.Errorf("%v/%v", n, d)
		}
		if d != c.denominator {
			t.Errorf("%v/%v", n, d)
		}

		for i := range c.fraction {
			if f[i] == c.fraction[i] {
				continue
			}
			t.Error(f)
		}
	}
}
