package number_test

import (
	"fmt"
	"testing"

	"github.com/itsubaki/q/math/number"
)

func ExampleGCD() {
	gcd := number.GCD(15, 7)
	fmt.Println(gcd)

	// Output:
	// 1
}

func TestGCD(t *testing.T) {
	cases := []struct {
		a, b int
		want int
	}{
		{15, 2, 1},
		{15, 4, 1},
		{15, 7, 1},
		{15, 11, 1},
		{15, 13, 1},
		{15, 14, 1},
	}

	for _, c := range cases {
		got := number.GCD(c.a, c.b)
		if got != c.want {
			t.Errorf("got=%v want=%d", got, c.want)
		}
	}
}
