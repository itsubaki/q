package density_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/itsubaki/q/math/matrix"
	"github.com/itsubaki/q/quantum/density"
)

func ExamplePauli() {
	pauli := density.Pauli(0.3, 0.3, 0.3, 0)(1)

	sum := matrix.Zero(2, 2)
	for _, k := range pauli.Kraus {
		sum = sum.Add(matrix.MatMul(k.Dagger(), k))
	}

	fmt.Println(sum.Equal(matrix.Identity(2)))

	// Output:
	// true
}

func FuzzPauli(f *testing.F) {
	f.Add(0.0, 0.0, 0.0)
	f.Add(0.1, 0.2, 0.3)
	f.Add(0.3, 0.3, 0.3)

	f.Fuzz(func(t *testing.T, pX, pY, pZ float64) {
		if math.IsNaN(pX) || math.IsNaN(pY) || math.IsNaN(pZ) {
			return
		}

		if math.IsInf(pX, 0) || math.IsInf(pY, 0) || math.IsInf(pZ, 0) {
			return
		}

		if pX < 0 || pY < 0 || pZ < 0 || pX+pY+pZ > 1 {
			return
		}

		pauli := density.Pauli(pX, pY, pZ, 0)(1)

		sum := matrix.Zero(2, 2)
		for _, k := range pauli.Kraus {
			sum = sum.Add(matrix.MatMul(k.Dagger(), k))
		}

		if !sum.Equal(matrix.Identity(2)) {
			t.Errorf("pX=%v pY=%v pZ=%v", pX, pY, pZ)
		}
	})
}
