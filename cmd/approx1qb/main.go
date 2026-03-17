package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"sort"

	"github.com/itsubaki/q"
	"github.com/itsubaki/q/math/epsilon"
)

const (
	H uint64 = 0
	T uint64 = 1
)

func main() {
	var n int
	var seq bool
	flag.IntVar(&n, "n", 16, "positive integer")
	flag.BoolVar(&seq, "seq", false, "print sequences")
	flag.Parse()

	if n >= 64 {
		panic("n must be < 64")
	}

	seqs := GenerateSequences(n)
	if seq {
		for _, s := range seqs {
			fmt.Println(s)
		}

		return
	}

	w := csv.NewWriter(os.Stdout)
	defer w.Flush()

	for _, s := range seqs {
		qsim := q.New()
		qb := qsim.Zero()

		for i := range s.length {
			switch (s.bits >> (s.length - 1 - i)) & 1 {
			case H:
				qsim.H(qb)
			case T:
				qsim.T(qb)
			}
		}

		amp := qsim.Amplitude()
		alpha, beta := amp[0], amp[1]
		theta, phi := Bloch(alpha, beta)

		if err := w.Write([]string{
			fmt.Sprintf("%.6f", theta),
			fmt.Sprintf("%.6f", phi),
		}); err != nil {
			panic(err)
		}
	}
}

type Seq struct {
	bits   uint64
	length int
}

func (s Seq) String() string {
	str := make([]string, 0, s.length)
	for i := range s.length {
		switch (s.bits >> (s.length - 1 - i)) & 1 {
		case H:
			str = append(str, "H")
		case T:
			str = append(str, "T")
		}
	}

	return fmt.Sprint(str)
}

func GenerateSequences(n int) []Seq {
	visited := make(map[uint64]struct{})
	var result []Seq

	var dfs func(bits uint64, length int, depth int)
	dfs = func(bits uint64, length int, depth int) {
		if length > 0 {
			key := (bits << 6) | uint64(length)
			if _, ok := visited[key]; !ok {
				visited[key] = struct{}{}
				result = append(result, Seq{
					bits,
					length,
				})
			}
		}

		if depth == n {
			return
		}

		for _, g := range []uint64{H, T} {
			nbits, nlen := (bits<<1)|g, length+1
			nbits, nlen = Simplify(nbits, nlen)
			dfs(nbits, nlen, depth+1)
		}
	}

	// depth first search
	dfs(0, 0, 0)

	// sort
	sort.Slice(result, func(i, j int) bool {
		if result[i].length != result[j].length {
			return result[i].length < result[j].length
		}

		return result[i].bits < result[j].bits
	})

	return result
}

func Simplify(bits uint64, length int) (uint64, int) {
	for {
		prevBits, prevLen := bits, length
		bits, length = reduce(bits, length, H, 2) // HH = I
		bits, length = reduce(bits, length, T, 4) // TTTT = I
		if bits == prevBits && length == prevLen {
			return bits, length
		}
	}
}

func reduce(bits uint64, length int, target uint64, mod int) (uint64, int) {
	var out uint64
	var outLen, count int
	for i := range length {
		b := (bits >> (length - 1 - i)) & 1
		if b == target {
			count++
			continue
		}

		for range count % mod {
			out = (out << 1) | target
			outLen++
		}

		out = (out << 1) | b
		outLen++

		// reset
		count = 0
	}

	for range count % mod {
		out = (out << 1) | target
		outLen++
	}

	return out, outLen
}

func Bloch(alpha, beta complex128) (float64, float64) {
	if epsilon.IsZeroF64(cmplx.Abs(alpha)) {
		return math.Pi, 0
	}

	if epsilon.IsZeroF64(cmplx.Abs(beta)) {
		return 0, 0
	}

	theta := 2 * math.Acos(min(1, cmplx.Abs(alpha)))
	phi := cmplx.Phase(beta) - cmplx.Phase(alpha)
	phi = math.Mod(phi+2*math.Pi, 2*math.Pi) // phi is in [0, 2π)
	return theta, phi
}
