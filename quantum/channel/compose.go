package channel

import "github.com/itsubaki/q/quantum/gate"

// ChannelFunc is a function type that generates a quantum channel for a given number of qubits.
type ChannelFunc func(n int) *Channel

// ComposeFunc returns a new quantum channel function that is the composition of multiple quantum channel functions.
func ComposeFunc(fn ...ChannelFunc) ChannelFunc {
	return func(n int) *Channel {
		ch := make([]*Channel, len(fn))
		for i, f := range fn {
			ch[i] = f(n)
		}

		return Compose(ch...)
	}
}

// Compose returns a new quantum channel that is the composition of multiple quantum channels.
func Compose(channels ...*Channel) *Channel {
	if len(channels) == 0 {
		return New(gate.I(1))
	}

	result := channels[0]
	for _, c := range channels[1:] {
		result = result.Compose(c)
	}

	return result
}
