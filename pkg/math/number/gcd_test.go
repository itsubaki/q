package number

import "testing"

func TestGCD(t *testing.T) {
	if GCD(15, 7) != 1 {
		t.Fail()
	}
}
