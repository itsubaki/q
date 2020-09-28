package matrix

import "testing"

func TestInverse(t *testing.T) {
	cases := []struct {
		m, e Matrix
	}{
		{
			New(
				[]complex128{1, 2, 0, -1},
				[]complex128{-1, 1, 2, 0},
				[]complex128{2, 0, 1, 1},
				[]complex128{1, -2, -1, 1},
			),
			New(
				[]complex128{1, 0, 0, 0},
				[]complex128{0, 1, 0, 0},
				[]complex128{0, 0, 1, 0},
				[]complex128{0, 0, 0, 1},
			),
		},
		{
			New(
				[]complex128{0, 1},
				[]complex128{1, 0},
			),
			New(
				[]complex128{1, 0},
				[]complex128{0, 1},
			),
		},
	}

	for _, c := range cases {
		inv := c.m.Inverse()
		mmi := c.m.Apply(inv)
		if !mmi.Equals(c.e) {
			t.Fail()
		}
	}
}

func TestCommutator(t *testing.T) {
	cases := []struct {
		x, y, e Matrix
	}{
		{
			New(
				[]complex128{0, 1},
				[]complex128{1, 0},
			),
			New(
				[]complex128{0, complex(0, -1)},
				[]complex128{complex(0, 1), 0},
			),
			New(
				[]complex128{complex(0, 2), 0},
				[]complex128{0, complex(0, -2)},
			),
		},
	}

	for _, c := range cases {
		if !Commutator(c.x, c.y).Equals(c.e) {
			t.Fail()
		}
	}
}

func TestAntiCommutator(t *testing.T) {
	cases := []struct {
		x, y, e Matrix
	}{
		{
			New(
				[]complex128{0, 1},
				[]complex128{1, 0},
			),
			New(
				[]complex128{0, complex(0, -1)},
				[]complex128{complex(0, 1), 0},
			),
			New(
				[]complex128{0, 0},
				[]complex128{0, 0},
			),
		},
	}

	for _, c := range cases {
		if !AntiCommutator(c.x, c.y).Equals(c.e) {
			t.Fail()
		}
	}
}
