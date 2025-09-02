package number_test

import (
	"errors"
	"testing"

	"github.com/itsubaki/q/math/number"
)

var ErrSomethingWentWrong = errors.New("something went wrong")

func TestMustPanic(t *testing.T) {
	defer func() {
		if rec := recover(); rec != nil {
			err, ok := rec.(error)
			if !ok {
				t.Fail()
			}

			if !errors.Is(err, ErrSomethingWentWrong) {
				t.Fail()
			}
		}
	}()

	number.Must(-1, ErrSomethingWentWrong)
	t.Fail()
}
