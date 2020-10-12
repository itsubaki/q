package number_test

import (
	"fmt"
	"testing"

	"github.com/itsubaki/q/pkg/math/number"
)

func ExampleFindOrder_mod15() {
	m := []int{1, 1, 0}
	s, r, d, ok := number.FindOrder(7, 15, m)
	fmt.Printf("%v/%v=%v %v\n", s, r, d, ok)
	fmt.Printf("%d^%d mod %d = %v\n", 7, r, 15, number.ModExp(7, r, 15))

	// Output:
	// 3/4=0.75 true
	// 7^4 mod 15 = 1
}

func ExampleFindOrder_mod21a2() {
	m := []int{0, 0, 1, 0, 1, 0, 1, 0, 1}
	s, r, d, ok := number.FindOrder(2, 21, m)
	fmt.Printf("%v/%v=%v %v\n", s, r, d, ok)
	fmt.Printf("%d^%d mod %d = %v\n", 2, r, 21, number.ModExp(2, r, 21))

	// Output:
	// 1/6=0.16666666666666666 true
	// 2^6 mod 21 = 1
}

func ExampleFindOrder_mod21a4() {
	m := []int{0, 1, 0, 1, 0, 1, 0, 1}
	s, r, d, ok := number.FindOrder(4, 21, m)
	fmt.Printf("%v/%v=%v %v\n", s, r, d, ok)
	fmt.Printf("%d^%d mod %d = %v\n", 4, r, 21, number.ModExp(4, r, 21))

	// Output:
	// 1/3=0.3333333333333333 true
	// 4^3 mod 21 = 1
}

func TestFindOrder(t *testing.T) {
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
		s, r, d, ok := number.FindOrder(c.a, c.N, c.m)
		if s != c.s || r != c.r || ok != c.ok || d != c.d {
			t.Errorf("%v/%v=%v %v", s, r, d, ok)
		}
	}
}
