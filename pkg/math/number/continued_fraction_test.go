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
	v := 0.41650390625
	f, n, d := number.ContinuedFraction(v)
	fmt.Printf("%v %v %v/%v\n", f, v, n, d)

	for i := 1; i < len(f); i++ {
		inv := number.InverseContinuedFraction(f, i)
		ff, nn, dd := number.ContinuedFraction(inv)
		fmt.Printf("%v: %v %v %v/%v\n", f[:i+1], inv, ff, nn, dd)
	}

	// Output:
	// [0 2 2 2 42 3 1] 0.41650390625 853/2048
	// [0 2]: 0.5 [0 2] 1/2
	// [0 2 2]: 0.4 [0 2 2] 2/5
	// [0 2 2 2]: 0.4166666666666667 [0 2 2 1 1] 5/12
	// [0 2 2 2 42]: 0.4165029469548134 [0 2 2 2 42] 212/509
	// [0 2 2 2 42 3]: 0.41650422352176736 [0 2 2 2 42 2 1] 641/1539
	// [0 2 2 2 42 3 1]: 0.41650390625000006 [0 2 2 2 42 3 1] 853/2048
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
