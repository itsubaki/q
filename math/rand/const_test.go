package rand_test

import (
	"fmt"
	"testing"

	"github.com/itsubaki/q/math/rand"
)

func ExampleConst() {
	fmt.Printf("%.13f\n", rand.Const()())
	fmt.Printf("%.13f\n", rand.Const(1)())
	fmt.Printf("%.13f\n", rand.Const(2)())
	fmt.Printf("%.13f\n", rand.Const(3)())

	// Output:
	// 0.6046602879796
	// 0.6046602879796
	// 0.1672966344259
	// 0.7199826688373
}

func TestConst(t *testing.T) {
	r := rand.Const()()
	if r >= 0 && r < 1 {
		return
	}

	t.Fail()
}