package number_test

import (
	"fmt"
	"testing"

	"github.com/itsubaki/q/pkg/math/number"
)

func ExampleContinuedFraction() {
	f, n, d := number.ContinuedFraction(0.8125)
	fmt.Printf("%v %v/%v\n", f, n, d)

	// Output:
	// [0 1 4 3] 13/16
}

func ExampleInverseContinuedFraction() {
	m := []int{0, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1}
	v := number.BinaryFraction(m)
	f, n, d := number.ContinuedFraction(v)
	fmt.Printf("%v %v %v/%v\n", f, v, n, d)

	for i := 1; i < len(f); i++ {
		inv := number.InverseContinuedFraction(f, i)
		ff, nn, dd := number.ContinuedFraction(inv)
		fmt.Printf("%v: %v %v %v/%v\n", f[:i+1], inv, ff, nn, dd)
	}

	// Output:
	// [0 6 170 1 1] 0.16650390625 341/2048
	// [0 6]: 0.16666666666666666 [0 6] 1/6
	// [0 6 170]: 0.1665034280117532 [0 6 170] 170/1021
	// [0 6 170 1]: 0.1665043816942551 [0 6 170 1] 171/1027
	// [0 6 170 1 1]: 0.16650390625 [0 6 170 1 1] 341/2048
}

func TestContinuedFraction(t *testing.T) {
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
		{0.0, []int{0}, 0, 1},
		{0.166656494140625, []int{0, 6}, 1, 6},
	}

	for _, c := range cases {
		f, n, d := number.ContinuedFraction(c.float)

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
