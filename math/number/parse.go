package number

import (
	"math"
	"strconv"
)

// Ldexp returns a * 2**b.
func Ldexp(a, b int) float64 {
	return math.Ldexp(float64(a), b)
}

// MustParseInt returns int from binary string.
// It panics if error occurs.
func MustParseInt(binary string) int {
	return int(Must(strconv.ParseInt(binary, 2, 0)))
}
