package number_test

import (
	"fmt"
	"testing"

	"github.com/itsubaki/q/math/number"
)

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
		want    int
	}{
		{0, 15, 1, 0},
		{15, 0, 1, 0},
		{15, 3, 1, 0},
	}

	for _, c := range cases {
		got := number.ModExp(c.a, c.r, c.N)
		if got != c.want {
			t.Errorf("got=%d, want=%v; %d^%d mod %d\n", got, c.want, c.a, c.r, c.N)
		}
	}
}

func TestModExp2(t *testing.T) {
	cases := []struct {
		a, j, N int
		want    int
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
		got := number.ModExp2(c.a, c.j, c.N)
		if got != c.want {
			t.Errorf("got=%d, %d^(2^%d) mod %d = %d\n", got, c.a, c.j, c.N, c.want)
		}
	}
}
