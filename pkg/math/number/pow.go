package number

import "math"

func Pow(a, r int) int {
	return int(math.Pow(float64(a), float64(r)))
}
