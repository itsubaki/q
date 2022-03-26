package rand

import (
	"crypto/rand"
	"math/big"

	"github.com/itsubaki/q/pkg/math/number"
)

func Crypto(_ ...int) float64 {
	v := CryptoInt64(0, 1000)
	return float64(v) / 1000
}

func CryptoInt64(min, max int64) int64 {
	n, err := rand.Int(rand.Reader, big.NewInt(max-min))
	if err != nil {
		panic(err)
	}

	return n.Int64() + min
}

func Coprime(N int) int {
	for {
		a := CryptoInt64(2, int64(N-1))
		if number.GCD(N, int(a)) == 1 {
			return int(a)
		}
	}
}
