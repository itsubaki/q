package number_test

import (
	"fmt"
	"testing"

	"github.com/itsubaki/q/pkg/math/number"
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
		binary string
		float  float64
	}{
		{"0.000", 0.0},
		{"0.100", 0.5},
		{"0.010", 0.25},
		{"0.110", 0.75},
		{"0.001", 0.125},
		{"0.101", 0.625},
		{"0.011", 0.375},
		{"0.111", 0.875},
		{"11.000", 3.0},
		{"11.010", 3.25},
		{"111", 7.0},
		{"0.01101010101", 0.41650390625},
		{"0.001010101010101", 0.166656494140625},
	}

	for _, c := range cases {
		result, err := number.ParseFloat(c.binary)
		if err != nil {
			t.Errorf("parse float: %v", err)
		}

		if result == c.float {
			continue
		}

		t.Errorf("expected=%v, actual=%v", c.float, result)
	}
}
