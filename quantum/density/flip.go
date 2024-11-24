package density

import (
	"fmt"
	"math"

	"github.com/itsubaki/q/math/matrix"
	"github.com/itsubaki/q/math/number"
	"github.com/itsubaki/q/quantum/gate"
)

// Flip returns the flip channel.
func Flip(p float64, m matrix.Matrix) (matrix.Matrix, matrix.Matrix, error) {
	if p < 0 || p > 1 {
		return nil, nil, ErrInvalidRange
	}

	d, d2 := m.Dimension()
	if d != d2 {
		return nil, nil, fmt.Errorf("matrix must be square: %w", ErrInvalidDimension)
	}

	if !number.IsPowOf2(d) {
		return nil, nil, fmt.Errorf("matrix dimension must be a power of 2: %w", ErrInvalidDimension)
	}

	n := number.Log2(d)
	e0 := gate.I(n).Mul(complex(math.Sqrt(p), 0))
	e1 := m.Mul(complex(math.Sqrt(1-p), 0))
	return e0, e1, nil
}

// BitFlip returns the bit flip channel.
func BitFlip(p float64) (matrix.Matrix, matrix.Matrix, error) {
	return Flip(p, gate.X())
}

// PhaseFlip returns the phase flip channel.
func PhaseFlip(p float64) (matrix.Matrix, matrix.Matrix, error) {
	return Flip(p, gate.Z())
}

// BitPhaseFlip returns the bit-phase flip channel.
func BitPhaseFlip(p float64) (matrix.Matrix, matrix.Matrix, error) {
	return Flip(p, gate.Y())
}
