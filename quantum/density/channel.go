package density

import (
	"math"

	"github.com/itsubaki/q/math/matrix"
	"github.com/itsubaki/q/quantum/gate"
)

// ChannelFunc is a function type that generates a quantum channel for a given number of qubits.
type ChannelFunc func(n int) *Channel

// Channel represents a quantum channel defined by Kraus operators.
type Channel struct {
	Kraus []*matrix.Matrix
}

// NewChannel returns a new quantum channel with the given Kraus operators.
func NewChannel(kraus ...*matrix.Matrix) *Channel {
	return &Channel{
		Kraus: kraus,
	}
}

// Pauli returns a new quantum channel that applies a Pauli channel to the specified qubit.
func Pauli(pX, pY, pZ float64, qb int) ChannelFunc {
	return func(n int) *Channel {
		e0 := gate.I().Mul(complex(math.Sqrt(1-pX-pY-pZ), 0))
		e1 := gate.X().Mul(complex(math.Sqrt(pX), 0))
		e2 := gate.Y().Mul(complex(math.Sqrt(pY), 0))
		e3 := gate.Z().Mul(complex(math.Sqrt(pZ), 0))

		k0 := gate.TensorProduct(e0, n, []int{qb})
		k1 := gate.TensorProduct(e1, n, []int{qb})
		k2 := gate.TensorProduct(e2, n, []int{qb})
		k3 := gate.TensorProduct(e3, n, []int{qb})
		return NewChannel(k0, k1, k2, k3)
	}
}

// Depolarizing returns a new quantum channel that applies a depolarizing channel to the specified qubit.
func Depolarizing(p float64, qb int) ChannelFunc {
	return Pauli(p/3, p/3, p/3, qb)
}

// Flip returns a new quantum channel that applies a flip channel to the specified qubit.
func Flip(p float64, u *matrix.Matrix, qb int) ChannelFunc {
	return func(n int) *Channel {
		e0 := gate.I().Mul(complex(math.Sqrt(1-p), 0))
		e1 := u.Mul(complex(math.Sqrt(p), 0))

		k0 := gate.TensorProduct(e0, n, []int{qb})
		k1 := gate.TensorProduct(e1, n, []int{qb})
		return NewChannel(k0, k1)
	}
}

// BitFlip returns a new quantum channel that applies a bit flip channel to the specified qubit.
func BitFlip(p float64, qb int) ChannelFunc {
	return Flip(p, gate.X(), qb)
}

// PhaseFlip returns a new quantum channel that applies a phase flip channel to the specified qubit.
func PhaseFlip(p float64, qb int) ChannelFunc {
	return Flip(p, gate.Z(), qb)
}

// BitPhaseFlip returns a new quantum channel that applies a bit-phase flip channel to the specified qubit.
func BitPhaseFlip(p float64, qb int) ChannelFunc {
	return Flip(p, gate.Y(), qb)
}

// AmplitudeDamping returns a new quantum channel that applies an amplitude damping channel to the specified qubit.
func AmplitudeDamping(gamma float64, qb int) ChannelFunc {
	return func(n int) *Channel {
		e0 := matrix.New(
			[]complex128{1, 0},
			[]complex128{0, complex(math.Sqrt(1-gamma), 0)},
		)

		e1 := matrix.New(
			[]complex128{0, complex(math.Sqrt(gamma), 0)},
			[]complex128{0, 0},
		)

		k0 := gate.TensorProduct(e0, n, []int{qb})
		k1 := gate.TensorProduct(e1, n, []int{qb})
		return NewChannel(k0, k1)
	}
}

// PhaseDamping returns a new quantum channel that applies a phase damping channel to the specified qubit.
func PhaseDamping(gamma float64, qb int) ChannelFunc {
	return func(n int) *Channel {
		e0 := matrix.New(
			[]complex128{1, 0},
			[]complex128{0, complex(math.Sqrt(1-gamma), 0)},
		)

		e1 := matrix.New(
			[]complex128{0, 0},
			[]complex128{0, complex(math.Sqrt(gamma), 0)},
		)

		k0 := gate.TensorProduct(e0, n, []int{qb})
		k1 := gate.TensorProduct(e1, n, []int{qb})
		return NewChannel(k0, k1)
	}
}
