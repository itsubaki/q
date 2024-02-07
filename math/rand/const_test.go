package rand_test

import (
	"fmt"
	"testing"

	"github.com/itsubaki/q/math/rand"
)

func ExampleConst() {
	c := rand.Const()
	fmt.Printf("%.13f\n", c())
	fmt.Printf("%.13f\n", c())
	fmt.Printf("%.13f\n", c())
	fmt.Printf("%.13f\n", rand.Const(1)())
	fmt.Printf("%.13f\n", rand.Const(1)())
	fmt.Printf("%.13f\n", rand.Const(1)())
	fmt.Printf("%.13f\n", rand.Const(2)())
	fmt.Printf("%.13f\n", rand.Const(3)())

	// Output:
	// 0.9999275824803
	// 0.8856419373529
	// 0.3814775277115
	// 0.3402859786606
	// 0.3402859786606
	// 0.3402859786606
	// 0.0782910683666
	// 0.3128509851974
}

func TestConst(t *testing.T) {
	r := rand.Const()()
	if r >= 0 && r < 1 {
		return
	}

	t.Fail()
}
