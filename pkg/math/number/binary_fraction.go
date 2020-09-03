package number

import (
	"fmt"
	"math"
)

func BinaryFraction(q ...int) float64 {
	var d float64
	for i, m := range q {
		if m != 0 && m != 1 {
			panic(fmt.Sprintf("invalid input: %v", q))
		}

		if m == 0 {
			continue
		}

		d = d + math.Pow(0.5, float64(i+1))
	}

	return d
}
