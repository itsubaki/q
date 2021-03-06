package number_test

import (
	"fmt"
	"testing"

	"github.com/itsubaki/q/pkg/math/number"
)

func ExampleBaseExp() {
	a, b, ok := number.BaseExp(125)
	fmt.Printf("%v^%v %v", a, b, ok)

	// Output:
	// 5^3 true
}

func ExampleModExp() {
	// 4^3 mod 21
	v := number.ModExp(4, 3, 21)
	fmt.Println(v)

	// Output:
	// 1
}

func ExampleModExp2() {
	// 7^2^4 mod 15
	v := number.ModExp2(7, 4, 15)
	fmt.Println(v)

	// Output:
	// 1
}

func TestModExp(t *testing.T) {
	cases := []struct {
		a, r, N int
		out     int
	}{
		{0, 15, 1, 0},
		{15, 0, 1, 1},
	}

	for _, c := range cases {
		v := number.ModExp(c.a, c.r, c.N)
		if v != c.out {
			t.Errorf("%d^%d) mod %d = %d. actual=%d\n", c.a, c.r, c.N, c.out, v)
		}
	}
}

func TestModExp2(t *testing.T) {
	cases := []struct {
		a, j, N int
		out     int
	}{
		{7, 0, 15, 7},
		{7, 1, 15, 4},
		{7, 2, 15, 1},
		{7, 3, 15, 1},
		{7, 4, 15, 1},
		{7, 5, 15, 1},
		{7, 6, 15, 1},
		{7, 7, 15, 1},
		{7, 8, 15, 1},
		{7, 9, 15, 1},
		{7, 10, 15, 1},
		{7, 11, 15, 1},
		{7, 12, 15, 1},
		{7, 13, 15, 1},
		{7, 14, 15, 1},
		{0, 15, 15, 0},
	}

	for _, c := range cases {
		v := number.ModExp2(c.a, c.j, c.N)
		if v != c.out {
			t.Errorf("%d^(2^%d) mod %d = %d. actual=%d\n", c.a, c.j, c.N, c.out, v)
		}
	}
}

func TestBaseExp(t *testing.T) {
	cases := []struct {
		N, a, b int
		ok      bool
	}{
		{9, 3, 2, true},
		{25, 5, 2, true},
		{27, 3, 3, true},
		{49, 7, 2, true},
		{125, 5, 3, true},
		{10, 0, 0, false},
	}

	for _, c := range cases {
		a, b, ok := number.BaseExp(c.N)
		if a != c.a || b != c.b || ok != c.ok {
			t.Errorf("N=%v, a=%v, b=%v, ok=%v", c.N, a, b, ok)
		}
	}
}

func TestPow(t *testing.T) {
	cases := []struct {
		a, r int
		p    int
	}{
		{0, 4, 0},
		{2, 0, 1},
	}

	for _, c := range cases {
		p := number.Pow(c.a, c.r)
		if p != c.p {
			t.Errorf("actual=%v, expected=%v", p, c.p)
		}
	}
}
