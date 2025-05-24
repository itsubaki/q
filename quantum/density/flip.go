package density

import (
	"errors"
	"math"

	"github.com/itsubaki/q/math/matrix"
	"github.com/itsubaki/q/math/number"
	"github.com/itsubaki/q/quantum/gate"
)

var (
	ErrInvalidRange     = errors.New("the probabikity is out of range [0, 1]")
	ErrNotSquare        = errors.New("the matrix is not square")
	ErrInvalidDimension = errors.New("the matrix dimensions is not a power of 2")
)

// Flip returns the flip channel.
func Flip(p float64, m *matrix.Matrix) (*matrix.Matrix, *matrix.Matrix, error) {
	if p < 0 || p > 1 {
		return nil, nil, ErrInvalidRange
	}

	d, d2 := m.Dimension()
	if d != d2 {
		return nil, nil, ErrNotSquare
	}

	if !number.IsPowOf2(d) {
		return nil, nil, ErrInvalidDimension
	}

	n := number.Log2(d)
	e0 := gate.I(n).Mul(complex(math.Sqrt(p), 0))
	e1 := m.Mul(complex(math.Sqrt(1-p), 0))
	return e0, e1, nil
}

// BitFlip returns the bit flip channel.
func BitFlip(p float64) (*matrix.Matrix, *matrix.Matrix, error) {
	return Flip(p, gate.X())
}

// PhaseFlip returns the phase flip channel.
func PhaseFlip(p float64) (*matrix.Matrix, *matrix.Matrix, error) {
	return Flip(p, gate.Z())
}

// BitPhaseFlip returns the bit-phase flip channel.
func BitPhaseFlip(p float64) (*matrix.Matrix, *matrix.Matrix, error) {
	return Flip(p, gate.Y())
}
