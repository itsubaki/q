package rand

import (
	"math/rand"
	"time"
)

func Math() func() float64 {
	s := time.Now().UnixNano()
	return rand.New(rand.NewSource(s)).Float64
}
