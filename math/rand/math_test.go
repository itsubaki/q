package rand_test

import (
	"fmt"
	"testing"

	"github.com/itsubaki/q/math/rand"
)

func ExampleMath() {
	fmt.Printf("%.13f\n", rand.Math(1))
	fmt.Printf("%.13f\n", rand.Math(2))
	fmt.Printf("%.13f\n", rand.Math(3))

	// Output:
	// 0.6046602879796
	// 0.1672966344259
	// 0.7199826688373
}

func TestMath(t *testing.T) {
	r := rand.Math()
	if r < 0 && r > 1 {
		t.Fail()
	}
}
