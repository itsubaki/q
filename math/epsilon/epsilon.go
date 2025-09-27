package epsilon

import (
	"math"
	"math/cmplx"
)

// Tol returns absolute and relative tolerances from the provided arguments.
// If only one value is provided, both absolute and relative tolerances are set to that value.
func Tol(tol ...float64) (float64, float64) {
	if len(tol) == 0 {
		return 1e-8, 1e-5
	}

	if len(tol) == 1 {
		return tol[0], tol[0]
	}

	return tol[0], tol[1]
}

// IsClose checks if two complex numbers are close within the given tolerances.
func IsClose(a, b complex128, tol ...float64) bool {
	atol, rtol := Tol(tol...)
	return cmplx.Abs(a-b) <= atol+rtol*math.Max(cmplx.Abs(a), cmplx.Abs(b))
}

// IsCloseF64 checks if two float64 numbers are close within the given tolerances.
func IsCloseF64(a, b float64, tol ...float64) bool {
	atol, rtol := Tol(tol...)
	return math.Abs(a-b) <= atol+rtol*math.Max(math.Abs(a), math.Abs(b))
}
