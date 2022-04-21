package rand_test

import (
	"testing"

	"github.com/itsubaki/q/pkg/math/rand"
)

func TestCrypto(t *testing.T) {
	r := rand.Crypto()
	if r < 0 && r > 1 {
		t.Fail()
	}
}

func TestCryptoInt64(t *testing.T) {
	r, err := rand.CryptoInt64(2, 14)
	if err != nil {
		t.Errorf("crypto: %v", err)
	}

	found := false
	for _, e := range []int64{2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13} {
		if r == e {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("coprime=%v", r)
	}
}

func TestCoprime(t *testing.T) {
	p := rand.Coprime(15)

	found := false
	for _, e := range []int{2, 4, 7, 8, 11, 13, 14} {
		if p == e {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("coprime=%v", p)
	}
}

func FuzzCryptoInt64(f *testing.F) {
	f.Add(int64(0), int64(3))
	f.Fuzz(func(t *testing.T, min, max int64) {
		if min >= max || max < 0 || min < 0 {
			return
		}

		rand.CryptoInt64(min, max)
	})
}
