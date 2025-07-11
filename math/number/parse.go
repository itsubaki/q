package number

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

var ErrInvalidParameter = errors.New("invalid parameter")

// ParseFloat returns float64 from binary string.
func ParseFloat(binary string) (float64, error) {
	for _, b := range binary {
		if b == '0' || b == '1' || b == '.' {
			continue
		}

		return 0, ErrInvalidParameter
	}

	bin := strings.Split(binary, ".")
	p, err := strconv.ParseInt(bin[0], 2, 0)
	if err != nil {
		return 0, fmt.Errorf("parse int: %w", err)
	}

	v := float64(p)
	if len(bin) == 1 {
		return v, nil
	}

	if len(bin) != 2 {
		return 0, ErrInvalidParameter
	}

	for i, b := range bin[1] {
		if b == '0' {
			continue
		}

		v = v + math.Pow(0.5, float64(i+1))
	}

	return v, nil
}
