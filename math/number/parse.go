package number

import "strconv"

// MustParseInt returns int from binary string.
// It panics if error occurs.
func MustParseInt(binary string) int {
	return int(Must(strconv.ParseInt(binary, 2, 0)))
}
