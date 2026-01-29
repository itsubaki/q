package number_test

import (
	"fmt"
	"testing"

	"github.com/itsubaki/q/math/number"
)

func ExampleLdexp() {
	fmt.Println(number.Ldexp(1, -2))
	fmt.Println(number.Ldexp(3, -2))
	fmt.Println(number.Ldexp(3, 4))

	// Output:
	// 0.25
	// 0.75
	// 48
}

func ExampleBaseExp() {
	a, b, ok := number.BaseExp(125)
	fmt.Printf("%v^%v %v", a, b, ok)

	// Output:
	// 5^3 true
}

func TestBaseExp(t *testing.T) {
	cases := []struct {
		in   int
		a, b int
		ok   bool
	}{
		{0, 0, 1, true},
		{1, 1, 1, true},
		{9, 3, 2, true},
		{25, 5, 2, true},
		{27, 3, 3, true},
		{49, 7, 2, true},
		{125, 5, 3, true},
		{10000, 10, 4, true},
		{10, 0, 0, false},
	}

	for _, c := range cases {
		a, b, ok := number.BaseExp(c.in)
		if a != c.a || b != c.b || ok != c.ok {
			t.Errorf("N=%v, a=%v, b=%v, ok=%v", c.in, a, b, ok)
		}
	}
}

func TestPow(t *testing.T) {
	cases := []struct {
		a, r int
		want int
	}{
		{0, 4, 0},
		{2, 0, 1},
	}

	for _, c := range cases {
		got := number.Pow(c.a, c.r)
		if got != c.want {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}
