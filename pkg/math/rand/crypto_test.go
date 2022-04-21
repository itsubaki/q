package rand_test

import (
	"fmt"
	"testing"

	"github.com/itsubaki/q/pkg/math/rand"
)

func ExampleCryptoInt64() {
	r, _ := rand.CryptoInt64(2, 14)

	for _, e := range []int64{2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13} {
		if r == e {
			fmt.Println("found")
			break
		}
	}

	// Output:
	// found
}

func ExampleCoprime() {
	p := rand.Coprime(15)

	for _, e := range []int{2, 4, 7, 8, 11, 13, 14} {
		if p == e {
			fmt.Println("found")
			break
		}
	}

	// Output:
	// found
}

func TestCrypto(t *testing.T) {
	r := rand.Crypto()
	if r < 0 && r > 1 {
		t.Fail()
	}
}

func TestCryptoInt64Error(t *testing.T) {
	cases := []struct {
		min, max int64
		want     error
	}{
		{-1, 0, fmt.Errorf("invalid parameter. min=-1, max=0")},
		{0, -1, fmt.Errorf("invalid parameter. min=0, max=-1")},
	}

	for _, c := range cases {
		_, got := rand.CryptoInt64(c.min, c.max)
		if got.Error() != c.want.Error() {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
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
