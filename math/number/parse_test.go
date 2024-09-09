package number_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/itsubaki/q/math/number"
)

func ExampleParseFloat() {
	// 0.101 -> 1/2 + 1/8 = 0.5 + 0.125 = 0.625
	f, err := number.ParseFloat("0.101")
	fmt.Println(f, err)

	// Output:
	// 0.625 <nil>
}

func TestParseFloat(t *testing.T) {
	cases := []struct {
		in   string
		want float64
		err  error
	}{
		{"0.000", 0.0, nil},
		{"0.100", 0.5, nil},
		{"0.010", 0.25, nil},
		{"0.110", 0.75, nil},
		{"0.001", 0.125, nil},
		{"0.101", 0.625, nil},
		{"0.011", 0.375, nil},
		{"0.111", 0.875, nil},
		{"11.000", 3.0, nil},
		{"11.010", 3.25, nil},
		{"111", 7.0, nil},
		{"0.01101010101", 0.41650390625, nil},
		{"0.001010101010101", 0.166656494140625, nil},
		{"a.bbb.ccc", 0, number.ErrInvalidParameter},
		{"a.bbb", 0, number.ErrInvalidParameter},
		{"a.001", 0, number.ErrInvalidParameter},
		{"0.bbb", 0, number.ErrInvalidParameter},
		{"0.1.0", 0, number.ErrInvalidParameter},
	}

	for _, c := range cases {
		got, err := number.ParseFloat(c.in)
		if err != nil && !errors.Is(err, c.err) {
			t.Errorf("parse float: %v", err)
		}

		if got == c.want {
			continue
		}

		t.Errorf("got=%v, want=%v", got, c.want)
	}
}

func FuzzParseFloat(f *testing.F) {
	seed := []string{"123", "101", "1.0101", "abc", "a.bc"}
	for i := range seed {
		f.Add(seed[i])
	}

	f.Fuzz(func(t *testing.T, v string) {
		number.ParseFloat(v)
	})
}
