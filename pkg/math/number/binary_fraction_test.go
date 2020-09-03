package number

import "testing"

func TestBinaryFraction(t *testing.T) {
	cases := []struct {
		input  []int
		output float64
	}{
		{[]int{1, 0, 0}, 0.5},
		{[]int{0, 1, 0}, 0.25},
		{[]int{0, 0, 1}, 0.125},
		{[]int{1, 1, 0}, 0.75},
		{[]int{1, 1, 1}, 0.875},
	}

	for _, c := range cases {
		result := BinaryFraction(c.input...)
		if result == c.output {
			continue
		}

		t.Errorf("expected=%v, actual=%v", c.output, result)
	}
}
