package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/itsubaki/q"
	"github.com/itsubaki/q/pkg/math/number"
	"github.com/itsubaki/q/pkg/math/rand"
)

// go run main.go --N 15 --t 4 --shot 10
func main() {
	var N, t, shot int
	flag.IntVar(&N, "N", 15, "")
	flag.IntVar(&t, "t", 4, "")
	flag.IntVar(&shot, "shot", 10, "")
	flag.Parse()

	if number.IsEven(N) {
		fmt.Printf("N=%v, p=%v, q=%v\n", N, 2, N/2)
		return
	}

	a := rand.Coprime(N)
	fmt.Printf("N=%v, a=%v, t=%v\n\n", N, a, t)

	qsim := q.New()
	r0 := qsim.ZeroWith(t)
	r1 := qsim.ZeroLog2(N)

	qsim.X(r1[len(r1)-1])
	print("initial state", qsim, r0, r1)

	qsim.H(r0...)
	print("create superposition", qsim, r0, r1)

	qsim.CModExp2(N, a, r0, r1)
	print("apply controlled-U", qsim, r0, r1)

	qsim.InvQFT(r0...)
	print("apply inverse QFT", qsim, r0, r1)

	qsim.Measure(r1...)
	print("measure reg1", qsim, r0, r1)

	for i := 0; i < shot; i++ {
		m := qsim.Clone().MeasureAsBinary(r0...)
		d := number.BinaryFraction(m)
		_, s, r := number.ContinuedFraction(d)

		ar2 := number.Pow(a, r/2)
		if number.IsOdd(r) || ar2%N == -1 {
			fmt.Printf("  i=%2d: N=%d, a=%d. s/r=%2d/%2d (%v=%.3f).\n", i, N, a, s, r, m, d)
			continue
		}

		p0 := number.GCD(ar2-1, N)
		p1 := number.GCD(ar2+1, N)

		found := " "
		for _, p := range []int{p0, p1} {
			if 1 < p && p < N && N%p == 0 {
				found = "*"
				break
			}
		}

		fmt.Printf("%s i=%2d: N=%d, a=%d. s/r=%2d/%2d (%v=%.3f). p=%v, q=%v.\n", found, i, N, a, s, r, m, d, p0, p1)
	}
}

func print(desc string, qsim *q.Q, reg ...[]q.Qubit) {
	fmt.Println(desc)

	max := number.Max(qsim.Probability())
	for _, s := range qsim.State(reg...) {
		p := strings.Repeat("*", int(s.Probability/max*32))
		fmt.Printf("%s: %s\n", s, p)
	}

	fmt.Println()
}
