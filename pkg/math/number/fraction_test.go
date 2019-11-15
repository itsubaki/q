package number

import (
	"testing"
)

func TestFraction(t *testing.T) {
	list, _, _ := Fraction(0.42857, 1e-3)

	e := []int{0, 2, 2, 1}
	for i := range e {
		if list[i] == e[i] {
			continue
		}

		t.Error(list)
	}
}
