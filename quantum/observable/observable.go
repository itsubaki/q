package observable

import (
	"github.com/itsubaki/q/math/matrix"
	"github.com/itsubaki/q/quantum/gate"
)

// Pauli returns a Pauli observable from the string representation.
func Pauli(s string) *matrix.Matrix {
	list := make([]*matrix.Matrix, 0)
	for _, c := range s {
		switch c {
		case 'I':
			list = append(list, I())
		case 'X':
			list = append(list, X())
		case 'Y':
			list = append(list, Y())
		case 'Z':
			list = append(list, Z())
		}
	}

	return matrix.TensorProduct(list...)
}

// I returns an identity observable.
func I(n ...int) *matrix.Matrix {
	return gate.I(n...)
}

// X returns a Pauli-X observable.
func X(n ...int) *matrix.Matrix {
	return gate.X(n...)
}

// Y returns a Pauli-Y observable.
func Y(n ...int) *matrix.Matrix {
	return gate.Y(n...)
}

// Z returns a Pauli-Z observable.
func Z(n ...int) *matrix.Matrix {
	return gate.Z(n...)
}
