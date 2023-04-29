package rand

import (
	"math/rand"
	"time"
)

func Math() func() float64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Float64
}
