package number_test

import (
	"fmt"
	"testing"

	"github.com/itsubaki/q/math/number"
)

func TestMustPanic(t *testing.T) {
	defer func() {
		if rec := recover(); rec != nil {
			err, ok := rec.(error)
			if !ok {
				t.Fail()
			}

			if err.Error() != "something went wrong" {
				t.Fail()
			}
		}
	}()

	number.Must(-1, fmt.Errorf("something went wrong"))
	t.Fail()
}
