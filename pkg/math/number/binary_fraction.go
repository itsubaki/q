package number

import (
	"fmt"
	"math"
)

func BinaryFraction(binary string) float64 {
	var d float64
	for i, b := range binary {
		if b != '0' && b != '1' {
			panic(fmt.Sprintf("invalid input: %v", binary))
		}

		if b == '0' {
			continue
		}

		d = d + math.Pow(0.5, float64(i+1))
	}

	return d
}
