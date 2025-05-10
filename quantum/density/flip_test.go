package density_test

import (
	"fmt"
	"testing"

	"github.com/itsubaki/q/math/matrix"
	"github.com/itsubaki/q/quantum/density"
)

func TestBitPhaseFlip(t *testing.T) {
	cases := []struct {
		in     float64
		hasErr bool
	}{
		{-1, true},
	}

	for _, c := range cases {
		_, _, err := density.BitPhaseFlip(c.in)
		if (err != nil) != c.hasErr {
			t.Errorf("err: %v", err)
			continue
		}
	}
}

func ExampleFlip() {
	if _, _, err := density.Flip(1, matrix.New(
		[]complex128{1, 1},
		[]complex128{1, 1},
		[]complex128{1, 1},
	)); err != nil {
		fmt.Println(err)
	}

	if _, _, err := density.Flip(1, matrix.New(
		[]complex128{1, 1, 1},
		[]complex128{1, 1, 1},
		[]complex128{1, 1, 1},
	)); err != nil {
		fmt.Println(err)
	}

	// Output:
	// the matrix is not square
	// the matrix dimensions is not a power of 2
}

func ExampleBitFlip() {
	m0, m1, _ := density.BitFlip(0.5)

	for _, r := range m0.Seq2() {
		fmt.Println(r)
	}
	fmt.Println()

	for _, r := range m1.Seq2() {
		fmt.Println(r)
	}

	// Output:
	// [(0.7071067811865476+0i) (0+0i)]
	// [(0+0i) (0.7071067811865476+0i)]
	//
	// [(0+0i) (0.7071067811865476+0i)]
	// [(0.7071067811865476+0i) (0+0i)]
}

func ExamplePhaseFlip() {
	m0, m1, _ := density.PhaseFlip(0.5)

	for _, r := range m0.Seq2() {
		fmt.Println(r)
	}
	fmt.Println()

	for _, r := range m1.Seq2() {
		fmt.Println(r)
	}

	// Output:
	// [(0.7071067811865476+0i) (0+0i)]
	// [(0+0i) (0.7071067811865476+0i)]
	//
	// [(0.7071067811865476+0i) (0+0i)]
	// [(0+0i) (-0.7071067811865476+0i)]
}

func ExampleBitPhaseFlip() {
	m0, m1, _ := density.BitPhaseFlip(0.5)

	for _, r := range m0.Seq2() {
		fmt.Println(r)
	}
	fmt.Println()

	for _, r := range m1.Seq2() {
		fmt.Println(r)
	}

	// Output:
	// [(0.7071067811865476+0i) (0+0i)]
	// [(0+0i) (0.7071067811865476+0i)]
	//
	// [(0+0i) (0-0.7071067811865476i)]
	// [(0+0.7071067811865476i) (0+0i)]
}
