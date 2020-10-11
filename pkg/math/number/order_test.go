package number_test

import (
	"fmt"
	"testing"

	"github.com/itsubaki/q/pkg/math/number"
)

func ExampleOrderFinding() {
	s, r, d, ok := number.OrderFinding(2, 21, []int{0, 0, 1, 0, 1, 0, 1, 0, 1})
	fmt.Printf("%v/%v=%v %v %v\n", s, r, d, ok, number.ModExp(2, r, 21))

	// Output:
	// 1/6=0.16666666666666666 true 1
}

func TestOrderFinding(t *testing.T) {
	cases := []struct {
		a, N int
		m    []int
		s, r int
		d    float64
		ok   bool
	}{
		{7, 15, []int{0, 1, 0}, 1, 4, 0.25, true},
		{7, 15, []int{1, 0, 0}, 1, 2, 0.50, false},
		{7, 15, []int{1, 1, 0}, 3, 4, 0.75, true},
		{7, 15, []int{}, 0, 1, 0, false},
		{7, 15, []int{1}, 1, 2, 0.5, false},
	}

	for _, c := range cases {
		s, r, d, ok := number.OrderFinding(c.a, c.N, c.m)
		if s != c.s || r != c.r || ok != c.ok || d != c.d {
			t.Errorf("%v/%v=%v %v", s, r, d, ok)
		}
	}
}
