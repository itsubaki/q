package number

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrInvalidParameter = errors.New("invalid parameter")

// ParseFloat returns float64 from binary string.
func ParseFloat(binary string) (float64, error) {
	if binary == "" || strings.Count(binary, ".") > 1 {
		return 0, fmt.Errorf("binary=%q: %w", binary, ErrInvalidParameter)
	}

	for _, b := range binary {
		if b != '0' && b != '1' && b != '.' {
			return 0, fmt.Errorf("binary=%q: %w", binary, ErrInvalidParameter)
		}
	}

	// integer part
	bin := strings.Split(binary, ".")
	p, err := strconv.ParseUint(bin[0], 2, 64)
	if err != nil {
		return 0, fmt.Errorf("parse binary=%q: %v: %w", binary, err, ErrInvalidParameter)
	}

	if len(bin) == 1 {
		return float64(p), nil
	}

	// fractional part
	frac, f := 0.0, 0.5
	for _, b := range bin[1] {
		if b == '1' {
			frac += f
		}

		f *= 0.5
	}

	return float64(p) + frac, nil
}
