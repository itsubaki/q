package number

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func ParseFloat(binary string) (float64, error) {
	for _, b := range binary {
		if b == '0' || b == '1' || b == '.' {
			continue
		}

		return 0, fmt.Errorf("invalid parameter. binary=%v", binary)
	}

	bin := strings.Split(binary, ".")
	if len(bin) > 2 {
		return 0, fmt.Errorf("invalid parameter. binary=%v", binary)
	}

	p, err := strconv.ParseInt(bin[0], 2, 0)
	if err != nil {
		return 0, fmt.Errorf("parse int. binary=%v", binary)
	}

	f := float64(p)
	if len(bin) < 2 {
		return f, nil
	}

	for i, b := range bin[1] {
		if b == '0' {
			continue
		}

		f = f + math.Pow(0.5, float64(i+1))
	}

	return f, nil
}
