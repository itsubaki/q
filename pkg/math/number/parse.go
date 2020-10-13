package number

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func ParseFloat(binary string) (float64, error) {
	bin := strings.Split(binary, ".")
	if len(bin) > 2 {
		return 0, fmt.Errorf("invalid parameter. binary=%v", binary)
	}

	bi, err := strconv.ParseInt(bin[0], 2, 0)
	if err != nil {
		return 0, fmt.Errorf("parse int. binary=%v", binary)
	}

	d := float64(bi)
	if len(bin) < 2 {
		return d, nil
	}

	for i, b := range bin[1] {
		if b != '0' && b != '1' {
			return 0, fmt.Errorf("invalid parameter. binary=%v", binary)
		}

		if b == '0' {
			continue
		}

		d = d + math.Pow(0.5, float64(i+1))
	}

	return d, nil
}
