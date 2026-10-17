package channel_test

import (
	"math"
	"testing"

	"github.com/itsubaki/q/quantum/channel"
)

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

		pauli := channel.Pauli(pX, pY, pZ, 0)(1)
		if !pauli.IsTracePreserving() {
			t.Errorf("pX=%v pY=%v pZ=%v", pX, pY, pZ)
		}
	})
}
