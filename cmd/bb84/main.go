package main

import (
	"flag"
	"fmt"
	"math/rand"

	"github.com/itsubaki/q"
)

type Basis int

const (
	Z Basis = iota
	X
)

const Unknown = -1

type Party struct {
	Name  string
	Basis Basis
	Bit   int
}

func NewParty(name string) *Party {
	return &Party{
		Name:  name,
		Basis: Basis(rand.Intn(2)),
		Bit:   Unknown,
	}
}

// Encode prepares a qubit in the party's bit and basis.
func (p *Party) Encode() (*q.Q, q.Qubit) {
	qsim := q.New()
	qb := qsim.Zero()

	if p.Bit == 1 {
		qsim.X(qb)
	}

	if p.Basis == X {
		qsim.H(qb)
	}

	return qsim, qb
}

// Measure obtains a bit by measuring in the party's basis.
func (p *Party) Measure(qsim *q.Q, qb q.Qubit) {
	if p.Basis == X {
		qsim.H(qb)
	}

	p.Bit = 0
	if qsim.Measure(qb).IsOne() {
		p.Bit = 1
	}
}

// Intercept measures and resends a qubit.
func (p *Party) Intercept(qsim *q.Q, qb q.Qubit) (*q.Q, q.Qubit) {
	p.Measure(qsim, qb)
	return p.Encode()
}

// bb84 simulates the BB84 quantum key distribution protocol
// with an optional intercept-resend attack.
// It returns the sifted key length and quantum bit error rate (QBER).
func bb84(n int, attack bool) (int, float64) {
	var sifted, errors int
	for range n {
		// Alice: encode and send.
		alice := NewParty("Alice")
		alice.Bit = rand.Intn(2)
		qsim, qb := alice.Encode()

		if attack {
			// Eve: intercept, measure, and resend.
			eve := NewParty("Eve")
			qsim, qb = eve.Intercept(qsim, qb)
		}

		// Bob: receive and measure.
		bob := NewParty("Bob")
		bob.Measure(qsim, qb)

		// Sifting: publicly compare bases.
		if alice.Basis != bob.Basis {
			continue
		}

		sifted++
		if alice.Bit != bob.Bit {
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
	var attack bool
	flag.IntVar(&n, "n", 10000, "number of transmitted qubits")
	flag.BoolVar(&attack, "attack", false, "enable intercept-resend attack")
	flag.Parse()

	sifted, qber := bb84(n, attack)
	fmt.Println("sifted:", sifted)
	fmt.Printf("qber  : %.2f%%\n", qber*100)
}
