package rand

import (
	"crypto/rand"
	"math"
	"math/big"

	"github.com/itsubaki/q/math/number"
)

// Crypto returns a random float64 in [0, 1).
func Crypto() float64 {
	v := CryptoInt64(0, math.MaxInt64)
	return float64(v) / math.MaxInt64
}

// Coprime returns a random coprime number in [2, N).
func Coprime(N int) int {
	for {
		a := int(CryptoInt64(2, int64(N-1)))
		if number.GCD(N, a) == 1 {
			return a
		}
	}
}

// CryptoInt64 returns a random int64 in [min, max).
func CryptoInt64(min, max int64) int64 {
	v := number.Must(rand.Int(rand.Reader, big.NewInt(max-min)))
	return v.Int64() + min
}
