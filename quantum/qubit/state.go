package qubit

import (
	"fmt"
	"math"
	"math/cmplx"
	"strconv"

	"github.com/itsubaki/q/math/epsilon"
	"github.com/itsubaki/q/math/number"
)

// State is a quantum state.
type State struct {
	amp          complex128
	prob         float64
	binaryString []string
	intValue     []int64
}

// NewState returns a new State.
func NewState(amp complex128, binary ...string) State {
	intv := make([]int64, len(binary))
	for i, b := range binary {
		intv[i] = number.Must(strconv.ParseInt(b, 2, 0))
	}

	return State{
		amp:          amp,
		prob:         math.Pow(cmplx.Abs(amp), 2),
		binaryString: binary,
		intValue:     intv,
	}
}

func (s State) Probability() float64 {
	return s.prob
}

func (s State) Amplitude() complex128 {
	return s.amp
}

func (s State) BinaryString() []string {
	return s.binaryString
}

func (s State) Int() []int64 {
	return s.intValue
}

func (s State) String() string {
	return fmt.Sprintf("%v%3v(% .4f% .4fi): %.4f",
		s.binaryString,
		s.intValue,
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

// Equal returns true if s equals v.
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
