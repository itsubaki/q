package number

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// ParseFloat returns float64 from binary string.
func ParseFloat(binary string) (float64, error) {
	for _, b := range binary {
		if b == '0' || b == '1' || b == '.' {
			continue
		}

		return 0, fmt.Errorf("invalid parameter. binary=%v", binary)
	}

	if !strings.Contains(binary, ".") {
		p := Must(strconv.ParseInt(binary, 2, 0))
		return float64(p), nil
	}

	bin := strings.Split(binary, ".")
	if len(bin) != 2 {
		return 0, fmt.Errorf("invalid parameter. binary=%v", binary)
	}

	p := Must(strconv.ParseInt(bin[0], 2, 0))
	v := float64(p)
	for i, b := range bin[1] {
		if b == '0' {
			continue
		}

		v = v + math.Pow(0.5, float64(i+1))
	}

	return v, nil
}
