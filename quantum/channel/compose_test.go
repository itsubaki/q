package channel_test

import (
	"testing"

	"github.com/itsubaki/q/quantum/channel"
)

func TestComposeFunc(t *testing.T) {
	cases := []struct {
		channel           []channel.ChannelFunc
		isTracePreserving bool
	}{
		{
			channel: []channel.ChannelFunc{
				channel.Pauli(0.3, 0.3, 0.3, 0),
			},
			isTracePreserving: true,
		},
		{
			channel: []channel.ChannelFunc{
				channel.Depolarizing(0.1, 0),
				channel.AmplitudeDamping(0.7, 0),
				channel.PhaseDamping(0.7, 0),
				channel.BitFlip(0.1, 0),
			},
			isTracePreserving: true,
		},
		{
			channel:           []channel.ChannelFunc{},
			isTracePreserving: true,
		},
	}

	for _, c := range cases {
		ch := channel.ComposeFunc(c.channel...)(1)
		if ch.IsTracePreserving() != c.isTracePreserving {
			t.Errorf("channel=%v", c.channel)
		}
	}
}
