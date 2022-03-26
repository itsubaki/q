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
	r := rand.CryptoInt64(2, 14)

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

func TestCryptoInt64Panic(t *testing.T) {
	defer func() {
		if err := recover(); err != "crypto/rand: argument to Int is <= 0" {
			t.Fail()
		}
	}()

	rand.CryptoInt64(0, 0)
	t.Fail()
}

func TestCryptoInt64PanicMin(t *testing.T) {
	defer func() {
		if err := recover(); err != "invalid parameter. min=-1" {
			t.Fail()
		}
	}()

	rand.CryptoInt64(-1, 0)
	t.Fail()
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
