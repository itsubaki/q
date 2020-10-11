package number_test

import (
	"fmt"
	"testing"

	"github.com/itsubaki/q/pkg/math/number"
)

func ExampleBinaryToFloat64() {
	// 0.101 -> 1/2 + 1/8 = 0.5 + 0.125 = 0.625
	f := number.BinaryToFloat64([]int{1, 0, 1})
	fmt.Println(f)

	// Output:
	// 0.625
}

func ExampleContinuedFraction() {
	c := number.ContinuedFraction(0.8125)
	s, r, d := number.Convergent(c)
	fmt.Printf("%v %v/%v=%v\n", c, s, r, d)

	// Output:
	// [0 1 4 3] 13/16=0.8125
}

func ExampleConvergent() {
	m := []int{0, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1}
	v := number.BinaryToFloat64(m)
	fmt.Printf("%v=%v\n", m, v)

	c := number.ContinuedFraction(v)
	for i := 0; i < len(c); i++ {
		s, r, d := number.Convergent(c[:i+1])
		fmt.Printf("%v: %v/%v=%v\n", c[:i+1], s, r, d)
	}

	// Output:
	// [0 0 1 0 1 0 1 0 1 0 1]=0.16650390625
	// [0]: 0/1=0
	// [0 6]: 1/6=0.16666666666666666
	// [0 6 170]: 170/1021=0.1665034280117532
	// [0 6 170 1]: 171/1027=0.1665043816942551
	// [0 6 170 1 1]: 341/2048=0.16650390625
}

func TestBinaryToFloat64(t *testing.T) {
	cases := []struct {
		binary []int
		float  float64
	}{
		{[]int{0, 0, 0}, 0.0},
		{[]int{1, 0, 0}, 0.5},
		{[]int{0, 1, 0}, 0.25},
		{[]int{1, 1, 0}, 0.75},
		{[]int{0, 0, 1}, 0.125},
		{[]int{1, 0, 1}, 0.625},
		{[]int{0, 1, 1}, 0.375},
		{[]int{1, 1, 1}, 0.875},
		{[]int{0, 1, 1, 0, 1, 0, 1, 0, 1, 0, 1}, 0.41650390625},
		{[]int{0, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1}, 0.166656494140625},
	}

	for _, c := range cases {
		result := number.BinaryToFloat64(c.binary)
		if result == c.float {
			continue
		}

		t.Errorf("expected=%v, actual=%v", c.float, result)
	}
}

func TestContinuedFraction(t *testing.T) {
	cases := []struct {
		input float64
		cf    []int
		s, r  int
		d     float64
	}{
		{1.0 / 16.0, []int{0, 16}, 1, 16, 0.0625},
		{4.0 / 16.0, []int{0, 4}, 1, 4, 0.25},
		{7.0 / 16.0, []int{0, 2, 3, 1, 1}, 7, 16, 0.4375},
		{13.0 / 16.0, []int{0, 1, 4, 3}, 13, 16, 0.8125},
		{0.42857, []int{0, 2, 2, 1}, 3, 7, 0.42857142857142855},
		{0.166656494140625, []int{0, 6}, 1, 6, 0.16666666666666666},
		{2.38461538462, []int{2, 2, 1, 1, 2}, 31, 13, 2.3846153846153846},
		{0.0, []int{0}, 0, 1, 0},
		{1.0, []int{1}, 1, 1, 1},
		{1.5, []int{1, 2}, 3, 2, 1.5},
		{2.0, []int{2}, 2, 1, 2.0},
	}

	for _, c := range cases {
		f := number.ContinuedFraction(c.input)
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
