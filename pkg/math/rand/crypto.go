package rand

import (
	"crypto/rand"
	"math/big"

	"github.com/itsubaki/q/pkg/math/number"
)

func Crypto(_ ...int64) float64 {
	v := CryptoInt(0, 1000)
	return float64(v) / 1000
}

func CryptoInt(min, max int) int {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max-min)))
	if err != nil {
		panic(err)
	}

	return int(n.Int64()) + min
}

func Coprime(N int) int {
	for {
		a := CryptoInt(2, N-1)
		if number.GCD(N, a) == 1 {
			return a
		}
	}
}
