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

	p := Must(strconv.ParseInt(bin[0], 2, 0))
	v := float64(p)

	if len(bin) < 2 {
		return v, nil
	}

	for i, b := range bin[1] {
		if b == '0' {
			continue
		}

		v = v + math.Pow(0.5, float64(i+1))
	}

	return v, nil
}

func Must[T int | int64 | float64](v T, err error) T {
	if err != nil {
		panic(err)
	}

	return v
}
