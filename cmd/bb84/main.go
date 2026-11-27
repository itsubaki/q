package main

import (
	"flag"
	"fmt"
	"math/rand"

	"github.com/itsubaki/q"
)

func bb84(n int, eve bool) (int, float64) {
	var sifted, errors int
	for range n {
		qsim := q.New()
		q0 := qsim.Zero()

		// Alice: randomly choose a bit and a basis.
		a := rand.Intn(2)
		abasis := rand.Intn(2)

		// Encode: 0 = Z basis, 1 = X basis.
		if a == 1 {
			qsim.X(q0)
		}

		if abasis == 1 {
			qsim.H(q0)
		}

		// Eve: intercept-resend attack.
		if eve {
			ebasis := rand.Intn(2)
			if ebasis == 1 {
				qsim.H(q0)
			}

			qsim.Measure(q0)
			if ebasis == 1 {
				qsim.H(q0)
			}
		}

		// Bob: randomly choose a measurement basis.
		bbasis := rand.Intn(2)
		if bbasis == 1 {
			qsim.H(q0)
		}

		var b int
		if qsim.Measure(q0).IsOne() {
			b = 1
		}

		// Sifting: discard mismatched bases.
		if abasis != bbasis {
			continue
		}

		sifted++
		if a != b {
			errors++
		}
	}

	if sifted == 0 {
		return 0, 0
	}

	return sifted, float64(errors) / float64(sifted)
}

func main() {
	var n int
	var eve bool
	flag.IntVar(&n, "n", 10000, "number of qubits")
	flag.BoolVar(&eve, "eve", false, "enable Eve")
	flag.Parse()

	sifted, qber := bb84(n, eve)
	fmt.Println("sifted:", sifted)
	fmt.Println("qber  :", qber)
}
