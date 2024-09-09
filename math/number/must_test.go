package number_test

import (
	"errors"
	"testing"

	"github.com/itsubaki/q/math/number"
)

var ErrSomtingWentWrong = errors.New("something went wrong")

func TestMustPanic(t *testing.T) {
	defer func() {
		if rec := recover(); rec != nil {
			err, ok := rec.(error)
			if !ok {
				t.Fail()
			}

			if !errors.Is(err, ErrSomtingWentWrong) {
				t.Fail()
			}
		}
	}()

	number.Must(-1, ErrSomtingWentWrong)
	t.Fail()
}
