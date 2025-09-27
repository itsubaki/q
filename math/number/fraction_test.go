package number_test

import (
	"fmt"
	"testing"

	"github.com/itsubaki/q/math/number"
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
	for i := range c {
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

func ExampleConvergent_n0() {
	var cfx []int
	s, r, d := number.Convergent(cfx)
	fmt.Println(s, r, d)

	// Output:
	// 0 1 0
}

func TestContinuedFraction(t *testing.T) {
	cases := []struct {
		in   float64
		cf   []int
		s, r int
		d    float64
		tol  []float64
	}{
		{
			in: 1.0 / 16.0,
			cf: []int{0, 16},
			s:  1,
			r:  16,
			d:  0.0625},
		{
			in: 4.0 / 16.0,
			cf: []int{0, 4},
			s:  1,
			r:  4,
			d:  0.25,
		},
		{
			in: 7.0 / 16.0,
			cf: []int{0, 2, 3, 1, 1},
			s:  7,
			r:  16,
			d:  0.4375,
		},
		{
			in: 13.0 / 16.0,
			cf: []int{0, 1, 4, 3},
			s:  13,
			r:  16,
			d:  0.8125,
		},
		{
			in:  0.42857,
			cf:  []int{0, 2, 2, 1},
			s:   3,
			r:   7,
			d:   0.42857142857142855,
			tol: []float64{1e-3, 1e-3},
		},
		{
			in:  0.166656494140625,
			cf:  []int{0, 6},
			s:   1,
			r:   6,
			d:   0.16666666666666666,
			tol: []float64{1e-3, 1e-3},
		},
		{
			in: 2.38461538462,
			cf: []int{2, 2, 1, 1, 2},
			s:  31,
			r:  13,
			d:  2.3846153846153846,
		},
		{
			in: 0.0,
			cf: []int{0},
			s:  0,
			r:  1,
			d:  0,
		},
		{
			in: 1.0,
			cf: []int{1},
			s:  1,
			r:  1,
			d:  1,
		},
		{
			in: 1.5,
			cf: []int{1, 2},
			s:  3,
			r:  2,
			d:  1.5,
		},
		{
			in: 2.0,
			cf: []int{2},
			s:  2,
			r:  1,
			d:  2.0,
		},
	}

	for _, c := range cases {
		f := number.ContinuedFraction(c.in, c.tol...)
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
