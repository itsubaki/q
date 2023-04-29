package rand_test

import (
	"testing"

	"github.com/itsubaki/q/math/rand"
)

func TestMath(t *testing.T) {
	r := rand.Math()()
	if r >= 0 && r < 1 {
		return
	}

	t.Fail()
}
