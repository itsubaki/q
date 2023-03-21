package qubit

import (
	"fmt"
	"math/cmplx"
	"strconv"

	"github.com/itsubaki/q/math/epsilon"
	"github.com/itsubaki/q/math/number"
)

// State is a quantum state.
type State struct {
	Amplitude    complex128
	Probability  float64
	Int          []int64
	BinaryString []string
}

func (s *State) Add(binary string) {
	s.Int = append(s.Int, number.Must(strconv.ParseInt(binary, 2, 0)))
	s.BinaryString = append(s.BinaryString, binary)
}

func (s State) Value(index ...int) (int64, string) {
	var i int
	if len(index) > 0 {
		i = index[0]
	}

	return s.Int[i], s.BinaryString[i]
}

// Equals returns true if s equals v.
// If eps is not given, epsilon.E13 is used.
func (s State) Equals(v State, eps ...float64) bool {
	if len(s.Int) != len(v.Int) {
		return false
	}

	if len(s.BinaryString) != len(v.BinaryString) {
		return false
	}

	for i := range s.Int {
		if s.Int[i] != v.Int[i] {
			return false
		}
	}

	for i := range s.BinaryString {
		if s.BinaryString[i] != v.BinaryString[i] {
			return false
		}
	}

	return cmplx.Abs(s.Amplitude-v.Amplitude) < epsilon.E13(eps...)
}

func (s State) String() string {
	return fmt.Sprintf("%v%3v(% .4f% .4fi): %.4f", s.BinaryString, s.Int, real(s.Amplitude), imag(s.Amplitude), s.Probability)
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
