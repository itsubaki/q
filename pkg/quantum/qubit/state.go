package qubit

import "fmt"

type State struct {
	Amplitude    complex128
	Probability  float64
	Int          []int64
	BinaryString []string
}

func (s State) Value(index ...int) (int64, string) {
	var i int
	if len(index) > 0 {
		i = index[0]
	}

	return s.Int[i], s.BinaryString[i]
}

func (s State) String() string {
	return fmt.Sprintf("%v%3v(% .4f% .4fi): %.4f", s.BinaryString, s.Int, real(s.Amplitude), imag(s.Amplitude), s.Probability)
}
