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
		qb := qsim.Zero()

		// Alice: randomly choose a bit and a basis.
		aBit := rand.Intn(2)
		aBasis := rand.Intn(2)

		// Encode: 0 = Z basis, 1 = X basis.
		if aBit == 1 {
			qsim.X(qb)
		}

		if aBasis == 1 {
			qsim.H(qb)
		}

		// Eve: intercept-resend attack.
		if eve {
			eBasis := rand.Intn(2)
			if eBasis == 1 {
				qsim.H(qb)
			}

			qsim.Measure(qb)
			if eBasis == 1 {
				qsim.H(qb)
			}
		}

		// Bob: randomly choose a measurement basis.
		bBasis := rand.Intn(2)
		if bBasis == 1 {
			qsim.H(qb)
		}

		var bBit int
		if qsim.Measure(qb).IsOne() {
			bBit = 1
		}

		// Sifting: discard mismatched bases.
		if aBasis != bBasis {
			continue
		}

		sifted++
		if aBit != bBit {
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
	fmt.Printf("qber  : %.2f%%\n", qber*100)
}
