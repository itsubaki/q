package qubit

import (
	"fmt"
	"math/cmplx"
	"strconv"

	"github.com/itsubaki/q/pkg/math/epsilon"
	"github.com/itsubaki/q/pkg/math/number"
)

type State struct {
	Amplitude    complex128
	Probability  float64
	Int          []int64
	BinaryString []string
}

func (s *State) Add(binary string) {
	i := number.Must(strconv.ParseInt(binary, 2, 0))

	s.Int = append(s.Int, i)
	s.BinaryString = append(s.BinaryString, binary)
}

func (s State) Value(index ...int) (int64, string) {
	var i int
	if len(index) > 0 {
		i = index[0]
	}

	return s.Int[i], s.BinaryString[i]
}

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

	e := epsilon.E13(eps...)
	if cmplx.Abs(s.Amplitude-v.Amplitude) > e {
		return false
	}

	return true
}

func (s State) String() string {
	return fmt.Sprintf("%v%3v(% .4f% .4fi): %.4f", s.BinaryString, s.Int, real(s.Amplitude), imag(s.Amplitude), s.Probability)
}

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
