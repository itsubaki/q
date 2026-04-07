package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"runtime/pprof"
	"sort"

	"github.com/itsubaki/q"
	"github.com/itsubaki/q/math/epsilon"
)

const (
	H       uint64 = 0
	T       uint64 = 1
	None    uint64 = 2
	HCancel int    = 2 // H**2 = I
	TCancel int    = 8 // T**8 = I
	maxLen  int    = 6 // 2**6 = 64
)

func main() {
	f, err := os.Create("cpu.prof")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	if err := pprof.StartCPUProfile(f); err != nil {
		panic(err)
	}
	defer pprof.StopCPUProfile()

	var n int
	var seq bool
	flag.IntVar(&n, "n", 16, "positive integer")
	flag.BoolVar(&seq, "seq", false, "print sequences")
	flag.Parse()

	if n > 1<<maxLen-1 || n < 1 {
		panic("n must be in [1, 63]")
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

// Bloch returns the polar(theta) and azimuthal(phi) angles of the state on the Bloch sphere.
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

func Sort(seqs []Seq) []Seq {
	sort.Slice(seqs, func(i, j int) bool {
		if seqs[i].length != seqs[j].length {
			return seqs[i].length < seqs[j].length
		}

		return seqs[i].bits < seqs[j].bits
	})

	return seqs
}

func GenerateSequences(n int) []Seq {
	var dfs func(bits uint64, length int, depth int)
	var seqs []Seq
	visited := make(map[uint64]struct{})

	dfs = func(bits uint64, length int, depth int) {
		if length > 0 {
			key := (bits << maxLen) | uint64(length)
			if _, ok := visited[key]; !ok {
				visited[key] = struct{}{}
				seqs = append(seqs, Seq{
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

	dfs(0, 0, 0)
	return Sort(seqs)
}

// Simplify returns a simplified sequence by applying cancellation rules (H**2 = I and T**8 = I).
func Simplify(bits uint64, length int) (uint64, int) {
	var out uint64
	var outLen int
	flush := func(bit uint64, cnt int) {
		var mod int
		switch bit {
		case H:
			mod = HCancel
		case T:
			mod = TCancel
		}

		k := cnt % mod
		for range k {
			out = (out << 1) | bit
			outLen++
		}
	}

	prev := None
	var count int
	for i := range length {
		b := (bits >> (length - 1 - i)) & 1
		if b == prev {
			count++
			continue
		}

		if count > 0 {
			flush(prev, count)
		}

		prev, count = b, 1
	}

	if count > 0 {
		flush(prev, count)
	}

	return out, outLen
}
