package channel_test

import (
	"testing"

	"github.com/itsubaki/q/quantum/channel"
	"github.com/itsubaki/q/quantum/gate"
)

func TestChannel_IsTracePreserving(t *testing.T) {
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
		if c.channel.IsTracePreserving() != c.want {
			t.Errorf("channel=%v want=%v", c.channel, c.want)
		}
	}
}
