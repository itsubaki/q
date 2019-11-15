package number

import (
	"testing"
)

func TestFraction(t *testing.T) {
	list, n, d := Fraction(0.42857, 1e-3)

	if n != 3 || d != 7 {
		t.Errorf("%v/%v", n, d)
	}

	e := []int{0, 2, 2, 1}
	for i := range e {
		if list[i] == e[i] {
			continue
		}

		t.Error(list)
	}

	{
		list, n, d := Fraction(1.0/16.0, 1e-3)
		if n != 1 || d != 16 {
			t.Errorf("%v %v/%v", list, n, d)
		}
	}
	{
		list, n, d := Fraction(4.0/16.0, 1e-3)
		if n != 1 || d != 4 {
			t.Errorf("%v %v/%v", list, n, d)
		}
	}
	{
		list, n, d := Fraction(7.0/16.0, 1e-3)
		if n != 7 || d != 16 {
			t.Errorf("%v %v/%v", list, n, d)
		}
	}
	{
		list, n, d := Fraction(13.0/16.0, 1e-3)
		if n != 13 || d != 16 {
			t.Errorf("%v %v/%v", list, n, d)
		}
	}
}
