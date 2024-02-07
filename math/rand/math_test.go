package rand_test

import (
	"testing"

	"github.com/itsubaki/q/math/rand"
)

func TestFloat64(t *testing.T) {
	r := rand.Float64()
	if r >= 0 && r < 1 {
		return
	}

	t.Fail()
}
