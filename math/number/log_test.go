package number_test

import (
	"fmt"
	"testing"

	"github.com/itsubaki/q/math/number"
)

func ExampleLog2() {
	fmt.Println(number.Log2(8))
	fmt.Println(number.Log2(32))
	fmt.Println(number.Log2(15))
	fmt.Println(number.Log2(20))

	// Output:
	// 3
	// 5
	// 3
	// 4
}

func ExampleIsPowOf2() {
	fmt.Println(number.IsPowOf2(8))
	fmt.Println(number.IsPowOf2(32))
	fmt.Println(number.IsPowOf2(15))
	fmt.Println(number.IsPowOf2(20))

	// Output:
	// true
	// true
	// false
	// false
}

func TestLog2(t *testing.T) {
	cases := []struct {
		n    int
		want int
	}{
		{1, 0},
		{2, 1},
		{4, 2},
		{8, 3},
		{16, 4},
		{32, 5},
		{64, 6},
		{128, 7},
		{256, 8},
		{512, 9},
		{1024, 10},
		{2048, 11},
		{4096, 12},
		{8192, 13},
	}

	for _, c := range cases {
		got := number.Log2(c.n)
		if got != c.want {
			t.Errorf("got=%v want=%d", got, c.want)
		}
	}
}
