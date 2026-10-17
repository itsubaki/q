package channel

import (
	"github.com/itsubaki/q/math/matrix"
)

// Channel represents a quantum channel defined by Kraus operators.
type Channel struct {
	Kraus []*matrix.Matrix
}

// New returns a new quantum channel with the given Kraus operators.
func New(kraus ...*matrix.Matrix) *Channel {
	return &Channel{
		Kraus: kraus,
	}
}

// IsTracePreserving returns true if the quantum channel is trace-preserving.
func (c *Channel) IsTracePreserving(tol ...float64) bool {
	if len(c.Kraus) == 0 {
		return false
	}

	sum := matrix.ZeroLike(c.Kraus[0])
	for _, k := range c.Kraus {
		sum = sum.Add(matrix.MatMul(k.Dagger(), k))
	}

	return sum.IsIdentity(tol...)
}

// Compose returns a new quantum channel that is the composition of the current channel and another channel.
// another is applied after the current channel.
func (c *Channel) Compose(another *Channel) *Channel {
	var kraus []*matrix.Matrix
	for _, k1 := range c.Kraus {
		for _, k2 := range another.Kraus {
			kraus = append(kraus, matrix.MatMul(k2, k1))
		}
	}

	return New(kraus...)
}
