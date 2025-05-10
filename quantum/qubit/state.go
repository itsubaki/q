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
	for i, bin := range binary {
		intv[i] = number.Must(strconv.ParseInt(bin, 2, 0))
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

func (s State) BinaryString(index ...int) string {
	var i int
	if len(index) > 0 {
		i = index[0]
	}

	return s.binaryString[i]
}

func (s State) Int(index ...int) int64 {
	var i int
	if len(index) > 0 {
		i = index[0]
	}

	return s.intValue[i]
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

// Equals returns true if s equals v.
// If eps is not given, epsilon.E13 is used.
func (s State) Equals(v State, eps ...float64) bool {
	if len(s.binaryString) != len(v.binaryString) {
		return false
	}

	for i := range s.binaryString {
		if s.binaryString[i] != v.binaryString[i] {
			return false
		}
	}

	return cmplx.Abs(s.amp-v.amp) < epsilon.E13(eps...)
}

// Equals returns true if s equals v.
// If eps is not given, epsilon.E13 is used.
func Equals(s, v []State, eps ...float64) bool {
	if len(s) != len(v) {
		return false
	}

	for i := range s {
		if !s[i].Equals(v[i], eps...) {
			return false
		}
	}

	return true
}
