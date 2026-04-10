package number

import "strconv"

// MustParseInt returns the integer value of a binary string.
// It panics if an error occurs.
func MustParseInt(binary string) int {
	return int(Must(strconv.ParseInt(binary, 2, 0)))
}
