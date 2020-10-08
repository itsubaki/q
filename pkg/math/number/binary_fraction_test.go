package number_test

import (
	"fmt"
	"testing"

	"github.com/itsubaki/q/pkg/math/number"
)

func ExampleBinaryFraction() {
	// 0.101 -> 1/2 + 1/8 = 0.5 + 0.125 = 0.625
	f := number.BinaryFraction("101")
	fmt.Println(f)

	// Output:
	// 0.625
}

func TestBinaryFraction(t *testing.T) {
	cases := []struct {
		binary string
		float  float64
	}{
		{"000", 0.0},
		{"100", 0.5},
		{"010", 0.25},
		{"110", 0.75},
		{"001", 0.125},
		{"101", 0.625},
		{"011", 0.375},
		{"111", 0.875},
		{"01101010101", 0.41650390625},
	}

	for _, c := range cases {
		result := number.BinaryFraction(c.binary)
		if result == c.float {
			continue
		}

		t.Errorf("expected=%v, actual=%v", c.float, result)
	}
}
