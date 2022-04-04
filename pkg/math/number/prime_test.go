package number_test

import (
	"fmt"
	"testing"

	"github.com/itsubaki/q/pkg/math/number"
)

func ExampleIsTrivial() {
	fmt.Println(number.IsTrivial(21, 1, 21))
	fmt.Println(number.IsTrivial(21, 1, 7))

	// Output:
	// true
	// false
}

func TestIsPrime(t *testing.T) {
	cases := []struct {
		in   int
		want bool
	}{
		{1, false},
		{2, true},
		{3, true},
		{4, false},
		{5, true},
		{6, false},
		{7, true},
		{8, false},
		{9, false},
		{10, false},
		{11, true},
		{12, false},
		{13, true},
		{14, false},
		{15, false},
		{16, false},
		{17, true},
		{18, false},
		{19, true},
		{20, false},
		{21, false},
	}

	for _, c := range cases {
		got := number.IsPrime(c.in)
		if got != c.want {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}

func TestIsOdd(t *testing.T) {
	cases := []struct {
		in   int
		want bool
	}{
		{1, true},
		{2, false},
		{3, true},
		{4, false},
	}

	for _, c := range cases {
		if number.IsOdd(c.in) != c.want {
			t.Fail()
		}
	}
}
