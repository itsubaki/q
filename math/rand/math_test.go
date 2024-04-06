package rand_test

import (
	"fmt"
	"testing"

	"github.com/itsubaki/q/math/rand"
)

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
	fmt.Printf("%.13f\n", rand.Const(1, 0)())
	fmt.Printf("%.13f\n", rand.Const(1, 1)())
	fmt.Printf("%.13f\n", rand.Const(1, 2)())

	// Output:
	// 0.9999275824803
	// 0.8856419373529
	// 0.3814775277115
	// 0.2384231908739
	// 0.2384231908739
	// 0.2384231908739
	// 0.8269781200925
	// 0.8353847703964
	// 0.2384231908739
	// 0.3402859786606
	// 0.6764556596678
}

func TestFloat64(t *testing.T) {
	r := rand.Float64()
	if r >= 0 && r < 1 {
		return
	}

	t.Fail()
}

func TestConst(t *testing.T) {
	r := rand.Const()()
	if r >= 0 && r < 1 {
		return
	}

	t.Fail()
}
