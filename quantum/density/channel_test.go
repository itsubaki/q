package density_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/itsubaki/q/quantum/density"
	"github.com/itsubaki/q/quantum/gate"
)

func ExampleChannel_IsValid() {
	pauli := density.Pauli(0.3, 0.3, 0.3, 0)(1)
	fmt.Println(pauli.IsValid())

	// Output:
	// true
}

func FuzzPauli(f *testing.F) {
	f.Add(0.0, 0.0, 0.0)
	f.Add(0.1, 0.2, 0.3)
	f.Add(0.3, 0.3, 0.3)

	f.Fuzz(func(t *testing.T, pX, pY, pZ float64) {
		if math.IsNaN(pX) || math.IsNaN(pY) || math.IsNaN(pZ) {
			return
		}

		if math.IsInf(pX, 0) || math.IsInf(pY, 0) || math.IsInf(pZ, 0) {
			return
		}

		if pX < 0 || pY < 0 || pZ < 0 || pX+pY+pZ > 1 {
			return
		}

		pauli := density.Pauli(pX, pY, pZ, 0)(1)
		if !pauli.IsValid() {
			t.Errorf("pX=%v pY=%v pZ=%v", pX, pY, pZ)
		}
	})
}

func TestChannel_IsValid(t *testing.T) {
	cases := []struct {
		channel *density.Channel
		want    bool
	}{
		{
			channel: density.NewChannel(),
			want:    false,
		},
		{
			channel: density.Pauli(0.3, 0.3, 0.3, 0)(1),
			want:    true,
		},
		{
			channel: density.Depolarizing(0.5, 0)(1),
			want:    true,
		},
		{
			channel: density.PhaseDamping(0.5, 0)(1),
			want:    true,
		},
		{
			channel: density.AmplitudeDamping(0.5, 0)(1),
			want:    true,
		},
		{
			channel: density.Flip(0.5, gate.Z(), 0)(1),
			want:    true,
		},
	}

	for _, c := range cases {
		if c.channel.IsValid() != c.want {
			t.Errorf("channel=%v want=%v", c.channel, c.want)
		}
	}
}

func TestCompose(t *testing.T) {
	cases := []struct {
		channel []density.ChannelFunc
		isValid bool
	}{
		{
			channel: []density.ChannelFunc{
				density.Pauli(0.3, 0.3, 0.3, 0),
			},
			isValid: true,
		},
		{
			channel: []density.ChannelFunc{
				density.Depolarizing(0.1, 0),
				density.AmplitudeDamping(0.7, 0),
				density.PhaseDamping(0.7, 0),
				density.BitFlip(0.1, 0),
			},
			isValid: true,
		},
		{
			channel: []density.ChannelFunc{},
			isValid: false,
		},
	}

	for _, c := range cases {
		noise := density.Compose(c.channel...)(1)
		if noise.IsValid() != c.isValid {
			t.Errorf("channel=%v", c.channel)
		}
	}
}
