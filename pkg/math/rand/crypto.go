package rand

import (
	"crypto/rand"
	"math/big"

	"github.com/itsubaki/q/pkg/math/number"
)

func Crypto(_ ...int64) float64 {
	n := CryptoInt(1000)
	v := n.Int64()
	return float64(v) / 1000
}

func CryptoInt(max int64) *big.Int {
	n, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		panic(err)
	}

	return n
}

func Coprime(N int) int {
	min, max := 2, (N - 1)
	for {
		v := CryptoInt(int64(max - min))
		a := int(v.Int64()) + min
		if number.GCD(N, a) == 1 {
			return a
		}
	}
}
