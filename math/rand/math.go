package rand

import "math/rand/v2"

func Float64() func() float64 {
	return rand.Float64
}
