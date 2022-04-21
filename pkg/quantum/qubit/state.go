package qubit

import "fmt"

type State struct {
	Amplitude    complex128
	Probability  float64
	Int          []int64
	BinaryString []string
}

func (s State) Value(index ...int) (int64, string, error) {
	if len(index) > 1 {
		return 0, "", fmt.Errorf("invalid parameter. len(index)=%v", len(index))
	}

	i := 0
	if len(index) > 0 {
		i = index[0]
	}

	if i < 0 || i > len(s.Int)-1 {
		return 0, "", fmt.Errorf("invalid parameter. index=%v", index)
	}

	return s.Int[i], s.BinaryString[i], nil
}

func (s State) String() string {
	return fmt.Sprintf("%v%3v(% .4f% .4fi): %.4f", s.BinaryString, s.Int, real(s.Amplitude), imag(s.Amplitude), s.Probability)
}
