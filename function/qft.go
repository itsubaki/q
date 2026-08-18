package function

import "github.com/itsubaki/q"

// QFT applies the quantum Fourier transform.
func QFT(qsim *q.Q, qb ...q.Qubit) {
	for i := range qb {
		qsim.H(qb[i])

		k := 2
		for j := i + 1; j < len(qb); j++ {
			qsim.CR(q.Theta(k), qb[i], qb[j])
			k++
		}
	}
}

// InvQFT applies the inverse quantum Fourier transform.
func InvQFT(qsim *q.Q, qb ...q.Qubit) {
	n := len(qb)
	for i := n - 1; i >= 0; i-- {
		k := n - i

		for j := n - 1; j > i; j-- {
			qsim.CR(-q.Theta(k), qb[j], qb[i])
			k--
		}

		qsim.H(qb[i])
	}
}
