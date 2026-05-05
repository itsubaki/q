package qubit

import (
	"fmt"
	"math"
	"math/cmplx"

	"github.com/itsubaki/q/math/epsilon"
)

// State is a quantum state.
type State struct {
	prob         float64
	amp          complex128
	binaryString []string
}

// NewState returns a new State.
func NewState(amp complex128, binary ...string) State {
	return State{
		prob:         math.Pow(cmplx.Abs(amp), 2),
		amp:          amp,
		binaryString: binary,
	}
}

// Probability returns the probability of the state.
func (s State) Probability() float64 {
	return s.prob
}

// Amplitude returns the amplitude of the state.
func (s State) Amplitude() complex128 {
	return s.amp
}

// BinaryString returns the binary string representation of the state.
func (s State) BinaryString() []string {
	return s.binaryString
}

// String returns the string representation of the state.
func (s State) String() string {
	return fmt.Sprintf("%v (% .4f% .4fi): %.4f",
		s.binaryString,
		real(s.amp),
		imag(s.amp),
		s.prob,
	)
}

// Equal returns true if s equals v.
func (s State) Equal(v State, tol ...float64) bool {
	if len(s.binaryString) != len(v.binaryString) {
		return false
	}

	for i := range s.binaryString {
		if s.binaryString[i] != v.binaryString[i] {
			return false
		}
	}

	return epsilon.IsClose(s.amp, v.amp, tol...)
}

// Equal returns true if s and v are equal.
func Equal(s, v []State, tol ...float64) bool {
	if len(s) != len(v) {
		return false
	}

	for i := range s {
		if !s[i].Equal(v[i], tol...) {
			return false
		}
	}

	return true
}

// EqualUpToGlobalPhase returns true if s equals v up to global phase.
func EqualUpToGlobalPhase(s, v []State, tol ...float64) bool {
	if len(s) != len(v) {
		return false
	}

	var dot complex128
	for i := range s {
		dot += s[i].amp * cmplx.Conj(v[i].amp)
	}

	return epsilon.IsOneF64(cmplx.Abs(dot), tol...)
}
