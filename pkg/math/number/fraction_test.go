package number_test

import (
	"fmt"
	"testing"

	"github.com/itsubaki/q/pkg/math/epsilon"
	"github.com/itsubaki/q/pkg/math/number"
)

func ExampleContinuedFraction() {
	c := number.ContinuedFraction(0.8125)
	s, r, d := number.Convergent(c)
	fmt.Printf("%v %v/%v=%v\n", c, s, r, d)

	// Output:
	// [0 1 4 3] 13/16=0.8125
}

func ExampleConvergent() {
	m := "0.00101010101"
	v, err := number.ParseFloat(m)
	fmt.Printf("%v=%v %v\n", m, v, err)

	c := number.ContinuedFraction(v)
	for i := 0; i < len(c); i++ {
		s, r, d := number.Convergent(c[:i+1])
		fmt.Printf("%v: %v/%v=%v\n", c[:i+1], s, r, d)
	}

	// Output:
	// 0.00101010101=0.16650390625 <nil>
	// [0]: 0/1=0
	// [0 6]: 1/6=0.16666666666666666
	// [0 6 170]: 170/1021=0.1665034280117532
	// [0 6 170 1]: 171/1027=0.1665043816942551
	// [0 6 170 1 1]: 341/2048=0.16650390625
}

func TestContinuedFraction(t *testing.T) {
	cases := []struct {
		in   float64
		cf   []int
		s, r int
		d    float64
		eps  float64
	}{
		{1.0 / 16.0, []int{0, 16}, 1, 16, 0.0625, epsilon.E3()},
		{4.0 / 16.0, []int{0, 4}, 1, 4, 0.25, epsilon.E3()},
		{7.0 / 16.0, []int{0, 2, 3, 1, 1}, 7, 16, 0.4375, epsilon.E3()},
		{13.0 / 16.0, []int{0, 1, 4, 3}, 13, 16, 0.8125, epsilon.E3()},
		{0.42857, []int{0, 2, 2, 1}, 3, 7, 0.42857142857142855, epsilon.E3()},
		{0.166656494140625, []int{0, 6}, 1, 6, 0.16666666666666666, epsilon.E3()},
		{2.38461538462, []int{2, 2, 1, 1, 2}, 31, 13, 2.3846153846153846, epsilon.E3()},
		{0.0, []int{0}, 0, 1, 0, epsilon.E3()},
		{1.0, []int{1}, 1, 1, 1, epsilon.E3()},
		{1.5, []int{1, 2}, 3, 2, 1.5, epsilon.E3()},
		{2.0, []int{2}, 2, 1, 2.0, epsilon.E3()},
	}

	for _, c := range cases {
		f := number.ContinuedFraction(c.in, c.eps)
		for i := range c.cf {
			if f[i] == c.cf[i] {
				continue
			}
			t.Errorf("%v %v", f, c.cf)
		}

		s, r, d := number.Convergent(f)
		if s != c.s || r != c.r || d != c.d {
			t.Errorf("%v/%v=%v %v/%v=%v", c.s, c.r, c.d, s, r, d)
		}
	}
}
