package decomp_test

import (
	"fmt"
	"math/cmplx"

	"github.com/itsubaki/q/math/decomp"
	"github.com/itsubaki/q/math/matrix"
)

func ExampleParlett() {
	t := matrix.New(
		[]complex128{2, 2, 3},
		[]complex128{0, 2, 5},
		[]complex128{0, 0, 3},
	)
	t2 := matrix.MatMul(t, t)

	pow := func(p complex128) (decomp.ParlettF, decomp.ParlettF) {
		return func(z complex128) complex128 {
				return cmplx.Pow(z, p)
			}, func(z complex128) complex128 {
				return p * cmplx.Pow(z, p-1)
			}
	}

	f, df := pow(2)
	a := decomp.Parlett(t, f, df)
	fmt.Println(a.Equals(t2))

	// Output:
	// true
}
