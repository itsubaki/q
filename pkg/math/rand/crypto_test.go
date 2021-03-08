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

func TestCryptoInt(t *testing.T) {
	r := rand.CryptoInt(2, 14)

	found := false
	for _, e := range []int{2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13} {
		if r == e {
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("coprime=%v", r)
	}
}

func TestCoprime(t *testing.T) {
	p := rand.Coprime(15)

	found := false
	for _, e := range []int{4, 7, 8, 11, 13, 14} {
		if p == e {
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("coprime=%v", p)
	}
}
