package density_test

import (
	"fmt"
	"testing"

	"github.com/itsubaki/q/quantum/density"
)

func TestFlip(t *testing.T) {
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

func ExampleBitFlip() {
	m0, m1, _ := density.BitFlip(0.5)

	for _, r := range m0 {
		fmt.Println(r)
	}
	fmt.Println()

	for _, r := range m1 {
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

	for _, r := range m0 {
		fmt.Println(r)
	}
	fmt.Println()

	for _, r := range m1 {
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

	for _, r := range m0 {
		fmt.Println(r)
	}
	fmt.Println()

	for _, r := range m1 {
		fmt.Println(r)
	}

	// Output:
	// [(0.7071067811865476+0i) (0+0i)]
	// [(0+0i) (0.7071067811865476+0i)]
	//
	// [(0+0i) (0-0.7071067811865476i)]
	// [(0+0.7071067811865476i) (0+0i)]
}
