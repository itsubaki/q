package rand

import (
	"math/rand"
	"time"
)

func Math(seed ...int) float64 {
	if len(seed) > 0 {
		rand.Seed(int64(seed[0]))
		return rand.Float64()
	}

	rand.Seed(time.Now().UnixNano())
	return rand.Float64()
}
