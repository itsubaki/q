package rand_test

import (
	"fmt"
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

func TestCryptoInt64Panic2(t *testing.T) {
	defer func() {
		if err := recover(); err != "crypto/rand: argument to Int is <= 0" {
			t.Fail()
		}
	}()

	rand.CryptoInt64(5, 2)
	t.Fail()
}

func TestCryptoInt64PanicMin(t *testing.T) {
	defer func() {
		if err := recover(); err != "invalid parameter. min=-1, max=0" {
			t.Fail()
		}
	}()

	rand.CryptoInt64(-1, 0)
	t.Fail()
}

func TestCryptoInt64PanicMax(t *testing.T) {
	defer func() {
		if err := recover(); err != "invalid parameter. min=0, max=-1" {
			t.Fail()
		}
	}()

	rand.CryptoInt64(0, -1)
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

func TestMustPanic(t *testing.T) {
	defer func() {
		if rec := recover(); rec != nil {
			err, ok := rec.(error)
			if !ok {
				t.Fail()
			}

			if err.Error() != "something went wrong" {
				t.Fail()
			}
		}
	}()

	rand.Must(nil, fmt.Errorf("something went wrong"))
	t.Fail()
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
