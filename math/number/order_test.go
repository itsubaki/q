package number_test

import (
	"fmt"
	"testing"

	"github.com/itsubaki/q/math/number"
)

func ExampleFindOrder_mod15() {
	s, r, d, ok := number.FindOrder(7, 15, number.Ldexp(6, -3))
	fmt.Printf("%v/%v=%v %v\n", s, r, d, ok)
	fmt.Printf("%d^%d mod %d = %v\n", 7, r, 15, number.ModExp(7, r, 15))

	// Output:
	// 3/4=0.75 true
	// 7^4 mod 15 = 1
}

func ExampleFindOrder_mod21a2() {
	s, r, d, ok := number.FindOrder(2, 21, number.Ldexp(85, -9))
	fmt.Printf("%v/%v=%v %v\n", s, r, d, ok)
	fmt.Printf("%d^%d mod %d = %v\n", 2, r, 21, number.ModExp(2, r, 21))

	// Output:
	// 1/6=0.16666666666666666 true
	// 2^6 mod 21 = 1
}

func ExampleFindOrder_mod21a4() {
	s, r, d, ok := number.FindOrder(4, 21, number.Ldexp(85, -8))
	fmt.Printf("%v/%v=%v %v\n", s, r, d, ok)
	fmt.Printf("%d^%d mod %d = %v\n", 4, r, 21, number.ModExp(4, r, 21))

	// Output:
	// 1/3=0.3333333333333333 true
	// 4^3 mod 21 = 1
}

func TestFindOrder(t *testing.T) {
	cases := []struct {
		a, N int
		m    float64
		s, r int
		d    float64
		ok   bool
	}{
		{7, 15, number.Ldexp(2, -3), 1, 4, 0.25, true},
		{7, 15, number.Ldexp(4, -3), 1, 2, 0.50, false},
		{7, 15, number.Ldexp(6, -3), 3, 4, 0.75, true},
		{7, 15, number.Ldexp(1, -1), 1, 2, 0.5, false},
	}

	for _, c := range cases {
		s, r, d, ok := number.FindOrder(c.a, c.N, c.m)
		if s != c.s || r != c.r || ok != c.ok || d != c.d {
			t.Errorf("%v/%v=%v %v", s, r, d, ok)
		}
	}
}
