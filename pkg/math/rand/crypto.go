package rand

import (
	"crypto/rand"
	"math/big"

	"github.com/itsubaki/q/pkg/math/number"
)

func Crypto(_ ...int64) float64 {
	v := CryptoInt(1000)
	return float64(v) / 1000
}

func CryptoInt(max int) int {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		panic(err)
	}

	return int(n.Int64())
}

func Coprime(N int) int {
	min, max := 2, (N - 1)
	for {
		a := CryptoInt(max-min) + min
		if number.GCD(N, a) == 1 {
			return a
		}
	}
}
