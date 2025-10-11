package epsilon

import (
	"math"
	"math/cmplx"
)

const (
	AbsTol = 1e-8
	RelTol = 1e-5
)

// Tol returns absolute and relative tolerances from the provided arguments.
// If only one value is provided, both absolute and relative tolerances are set to that value.
func Tol(tol ...float64) (float64, float64) {
	if len(tol) == 0 {
		return AbsTol, RelTol
	}

	if len(tol) == 1 {
		return tol[0], tol[0]
	}

	return tol[0], tol[1]
}

// IsZero returns true if a complex128 number is close to zero within the given tolerances.
func IsZero(a complex128, tol ...float64) bool {
	return IsClose(a, 0.0, tol...)
}

// IsClose returns true if two complex128 numbers are close within the given tolerances.
func IsClose(a, b complex128, tol ...float64) bool {
	atol, rtol := Tol(tol...)
	return cmplx.Abs(a-b) <= atol+rtol*math.Max(cmplx.Abs(a), cmplx.Abs(b))
}

// IsZeroF64 returns true if a float64 number is close to zero within the given tolerances.
func IsZeroF64(a float64, tol ...float64) bool {
	return IsCloseF64(a, 0.0, tol...)
}

// IsCloseF64 returns true if two float64 numbers are close within the given tolerances.
func IsCloseF64(a, b float64, tol ...float64) bool {
	atol, rtol := Tol(tol...)
	return math.Abs(a-b) <= atol+rtol*math.Max(math.Abs(a), math.Abs(b))
}
