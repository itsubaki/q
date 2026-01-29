package number_test

import (
	"fmt"

	"github.com/itsubaki/q/math/number"
)

func ExampleMustParseInt() {
	fmt.Println(number.MustParseInt("101"))

	// Output:
	// 5
}
