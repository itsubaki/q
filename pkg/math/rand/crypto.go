package rand

import (
	"crypto/rand"
	"math"
	"math/big"

	"github.com/itsubaki/q/pkg/math/number"
)

func Crypto(_ ...int) float64 {
	v := MustCryptoInt64(0, math.MaxInt64)
	return float64(v) / math.MaxInt64
}

func Coprime(N int) int {
	for {
		a := int(MustCryptoInt64(2, int64(N-1)))
		if number.GCD(N, a) == 1 {
			return a
		}
	}
}

func MustCryptoInt64(min, max int64) int64 {
	v := number.Must(rand.Int(rand.Reader, big.NewInt(max-min)))
	return v.Int64() + min
}
