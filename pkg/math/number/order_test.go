package number_test

import (
	"fmt"
	"testing"

	"github.com/itsubaki/q/pkg/math/number"
)

func ExampleFindOrder() {
	s, r, ok := number.FindOrder(7, 15, 0.75)
	fmt.Printf("%v/%v %v\n", s, r, ok)

	// Output:
	// 3/4 true
}

func TestFindOrder(t *testing.T) {
	cases := []struct {
		a, N int
		m    float64
		s, r int
		ok   bool
	}{
		{7, 15, 0.25, 1, 4, true},
		{7, 15, 0.50, 0, 0, false},
		{7, 15, 0.75, 3, 4, true},
	}

	for _, c := range cases {
		s, r, ok := number.FindOrder(c.a, c.N, c.m)
		if s != c.s || r != c.r || ok != c.ok {
			t.Errorf("%v/%v %v", s, r, ok)
		}
	}
}
