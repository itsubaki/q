package rand

import (
	"math/rand"
	"time"
)

func Math() func() float64 {
	src := rand.NewSource(time.Now().UnixNano())
	return rand.New(src).Float64
}
