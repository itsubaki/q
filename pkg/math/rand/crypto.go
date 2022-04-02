package rand

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"

	"github.com/itsubaki/q/pkg/math/number"
)

func Crypto(_ ...int) float64 {
	v := CryptoInt64(0, math.MaxInt64)
	return float64(v) / math.MaxInt64
}

func CryptoInt64(min, max int64) int64 {
	if min < 0 || max < 0 {
		panic(fmt.Sprintf("invalid parameter. min=%v, max=%v", min, max))
	}

	v := number.Must(rand.Int(rand.Reader, big.NewInt(max-min)))
	return v.Int64() + min
}

func Coprime(N int) int {
	for {
		a := CryptoInt64(2, int64(N-1))
		if number.GCD(N, int(a)) == 1 {
			return int(a)
		}
	}
}
