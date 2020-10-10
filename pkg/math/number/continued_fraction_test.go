package number_test

import (
	"fmt"
	"testing"

	"github.com/itsubaki/q/pkg/math/number"
)

func ExampleContinuedFraction() {
	c := number.ContinuedFraction(0.8125)
	s, r := number.Fraction(c)
	fmt.Printf("%v %v/%v=%v\n", c, s, r, float64(s)/float64(r))

	// Output:
	// [0 1 4 3] 13/16=0.8125
}

func ExampleFraction() {
	m := []int{0, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1}
	v := number.BinaryFraction(m)
	fmt.Printf("%v\n", v)

	c := number.ContinuedFraction(v)
	s, r := number.Fraction(c)
	fmt.Printf("%v %v/%v=%v\n", c, s, r, v)

	for i := 2; i < len(c)+1; i++ {
		s, r := number.Fraction(c[:i])
		fmt.Printf("%v: %v/%v=%v\n", c[:i], s, r, float64(s)/float64(r))
	}

	// Output:
	// 0.16650390625
	// [0 6 170 1 1] 341/2048=0.16650390625
	// [0 6]: 1/6=0.16666666666666666
	// [0 6 170]: 170/1021=0.1665034280117532
	// [0 6 170 1]: 171/1027=0.1665043816942551
	// [0 6 170 1 1]: 341/2048=0.16650390625
}

func TestContinuedFraction(t *testing.T) {
	cases := []struct {
		value    float64
		fraction []int
		s        int
		r        int
	}{
		{0.42857, []int{0, 2, 2, 1}, 3, 7},
		{1.0 / 16.0, []int{0, 16}, 1, 16},
		{4.0 / 16.0, []int{0, 4}, 1, 4},
		{7.0 / 16.0, []int{0, 2, 3, 1, 1}, 7, 16},
		{13.0 / 16.0, []int{0, 1, 4, 3}, 13, 16},
		{0.0, []int{0}, 0, 1},
		{1.0, []int{1}, 1, 1},
		{0.166656494140625, []int{0, 6}, 1, 6},
	}

	for _, c := range cases {
		f := number.ContinuedFraction(c.value)
		s, r := number.Fraction(f)

		if s != c.s || r != c.r {
			t.Errorf("%v/%v", s, r)
		}

		for i := range c.fraction {
			if f[i] == c.fraction[i] {
				continue
			}
			t.Error(f)
		}
	}
}
