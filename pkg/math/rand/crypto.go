package rand

import (
	"crypto/rand"
	"math/big"
)

func Crypto(_ ...int64) float64 {
	n, err := rand.Int(rand.Reader, big.NewInt(1000))
	if err != nil {
		panic(err)
	}

	return float64(n.Int64()) / 1000.0
}
