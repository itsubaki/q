package number_test

import (
	"fmt"
	"testing"

	"github.com/itsubaki/q/pkg/math/number"
)

func ExampleFindOrder() {
	s, r, ok := number.FindOrder(2, 21, []int{0, 0, 1, 0, 1, 0, 1, 0, 1})
	fmt.Printf("%v/%v %v %v\n", s, r, ok, number.ModExp(2, r, 21))

	// Output:
	// 1/6 true 1
}

func TestFindOrder(t *testing.T) {
	cases := []struct {
		a, N int
		m    []int
		s, r int
		ok   bool
	}{
		{7, 15, []int{0, 1, 0}, 1, 4, true},
		{7, 15, []int{1, 0, 0}, 1, 2, false},
		{7, 15, []int{1, 1, 0}, 3, 4, true},
		{7, 15, []int{}, 0, 1, false},
		{7, 15, []int{1}, 1, 2, false},
	}

	for _, c := range cases {
		s, r, ok := number.FindOrder(c.a, c.N, c.m)
		if s != c.s || r != c.r || ok != c.ok {
			t.Errorf("%v/%v %v", s, r, ok)
		}
	}
}
