package channel_test

import (
	"math"
	"testing"

	"github.com/itsubaki/q/quantum/channel"
	"github.com/itsubaki/q/quantum/gate"
)

func TestChannel_IsValid(t *testing.T) {
	cases := []struct {
		channel *channel.Channel
		want    bool
	}{
		{
			channel: channel.New(),
			want:    false,
		},
		{
			channel: channel.Pauli(0.3, 0.3, 0.3, 0)(1),
			want:    true,
		},
		{
			channel: channel.Depolarizing(0.5, 0)(1),
			want:    true,
		},
		{
			channel: channel.PhaseDamping(0.5, 0)(1),
			want:    true,
		},
		{
			channel: channel.AmplitudeDamping(0.5, 0)(1),
			want:    true,
		},
		{
			channel: channel.Flip(0.5, gate.Z(), 0)(1),
			want:    true,
		},
		{
			channel: channel.BitFlip(0.5, 0)(1),
			want:    true,
		},
		{
			channel: channel.PhaseFlip(0.5, 0)(1),
			want:    true,
		},
		{
			channel: channel.BitPhaseFlip(0.5, 0)(1),
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
		channel []channel.ChannelFunc
		isValid bool
	}{
		{
			channel: []channel.ChannelFunc{
				channel.Pauli(0.3, 0.3, 0.3, 0),
			},
			isValid: true,
		},
		{
			channel: []channel.ChannelFunc{
				channel.Depolarizing(0.1, 0),
				channel.AmplitudeDamping(0.7, 0),
				channel.PhaseDamping(0.7, 0),
				channel.BitFlip(0.1, 0),
			},
			isValid: true,
		},
		{
			channel: []channel.ChannelFunc{},
			isValid: true,
		},
	}

	for _, c := range cases {
		ch := channel.Compose(c.channel...)(1)
		if ch.IsValid() != c.isValid {
			t.Errorf("channel=%v", c.channel)
		}
	}
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

		pauli := channel.Pauli(pX, pY, pZ, 0)(1)
		if !pauli.IsValid() {
			t.Errorf("pX=%v pY=%v pZ=%v", pX, pY, pZ)
		}
	})
}
