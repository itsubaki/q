package q

import (
	"math"

	"github.com/itsubaki/q/pkg/math/matrix"
	"github.com/itsubaki/q/pkg/math/rand"
	"github.com/itsubaki/q/pkg/quantum/gate"
	"github.com/itsubaki/q/pkg/quantum/qubit"
)

type Qubit int

func (q Qubit) Index() int {
	return int(q)
}

func Index(qb ...Qubit) []int {
	index := make([]int, 0)
	for i := range qb {
		index = append(index, qb[i].Index())
	}

	return index
}

type Q struct {
	qb   *qubit.Qubit
	Seed []int
	Rand func(seed ...int) float64
}

func New() *Q {
	return &Q{
		qb:   nil,
		Rand: rand.Crypto,
	}
}

func (q *Q) New(v ...complex128) Qubit {
	if q.qb == nil {
		q.qb = qubit.New(v...)
		q.qb.Seed = q.Seed
		q.qb.Rand = q.Rand
		return Qubit(0)
	}

	q.qb.TensorProduct(qubit.New(v...))
	index := q.NumberOfBit() - 1
	return Qubit(index)
}

func (q *Q) NewOf(binary string) []Qubit {
	qb := make([]Qubit, 0)
	for _, b := range binary {
		if b == '0' {
			qb = append(qb, q.Zero())
			continue
		}

		qb = append(qb, q.One())
	}

	return qb
}

func (q *Q) Zero() Qubit {
	return q.New(1, 0)
}

func (q *Q) One() Qubit {
	return q.New(0, 1)
}

func (q *Q) ZeroWith(n int) []Qubit {
	qb := make([]Qubit, 0)
	for i := 0; i < n; i++ {
		qb = append(qb, q.Zero())
	}

	return qb
}

func (q *Q) OneWith(n int) []Qubit {
	qb := make([]Qubit, 0)
	for i := 0; i < n; i++ {
		qb = append(qb, q.One())
	}

	return qb
}

func (q *Q) ZeroLog2(N int) []Qubit {
	n := int(math.Log2(float64(N))) + 1
	return q.ZeroWith(n)
}

func (q *Q) NumberOfBit() int {
	return q.qb.NumberOfBit()
}

func (q *Q) Amplitude() []complex128 {
	return q.qb.Amplitude()
}

func (q *Q) Probability() []float64 {
	return q.qb.Probability()
}

func (q *Q) Reset(qb ...Qubit) {
	for i := range qb {
		if q.Measure(qb[i]).IsOne() {
			q.X(qb[i])
		}
	}
}

func (q *Q) U(theta, phi, lambda float64, qb ...Qubit) *Q {
	return q.Apply(gate.U(theta, phi, lambda), qb...)
}

func (q *Q) I(qb ...Qubit) *Q {
	return q.Apply(gate.I(), qb...)
}

func (q *Q) X(qb ...Qubit) *Q {
	return q.Apply(gate.X(), qb...)
}

func (q *Q) Y(qb ...Qubit) *Q {
	return q.Apply(gate.Y(), qb...)
}

func (q *Q) Z(qb ...Qubit) *Q {
	return q.Apply(gate.Z(), qb...)
}

func (q *Q) H(qb ...Qubit) *Q {
	return q.Apply(gate.H(), qb...)
}

func (q *Q) S(qb ...Qubit) *Q {
	return q.Apply(gate.S(), qb...)
}

func (q *Q) T(qb ...Qubit) *Q {
	return q.Apply(gate.T(), qb...)
}

func (q *Q) R(theta float64, qb ...Qubit) *Q {
	return q.Apply(gate.R(theta), qb...)
}

func (q *Q) RX(theta float64, qb ...Qubit) *Q {
	return q.Apply(gate.RX(theta), qb...)
}

func (q *Q) RY(theta float64, qb ...Qubit) *Q {
	return q.Apply(gate.RY(theta), qb...)
}

func (q *Q) RZ(theta float64, qb ...Qubit) *Q {
	return q.Apply(gate.RZ(theta), qb...)
}

func (q *Q) Apply(m matrix.Matrix, qb ...Qubit) *Q {
	if len(qb) < 1 {
		q.qb.Apply(m)
		return q
	}

	index := make(map[int]bool)
	for _, i := range Index(qb...) {
		index[i] = true
	}

	g := gate.I()
	if _, ok := index[0]; ok {
		g = m
	}

	for i := 1; i < q.NumberOfBit(); i++ {
		if _, ok := index[i]; ok {
			g = g.TensorProduct(m)
			continue
		}

		g = g.TensorProduct(gate.I())
	}

	q.qb.Apply(g)
	return q
}

func (q *Q) Controlled(m matrix.Matrix, control []Qubit, target Qubit) *Q {
	n := q.NumberOfBit()
	g := gate.Controlled(m, n, Index(control...), target.Index())
	q.qb.Apply(g)
	return q
}

func (q *Q) C(m matrix.Matrix, control, target Qubit) *Q {
	return q.Controlled(m, []Qubit{control}, target)
}

func (q *Q) ControlledNot(control []Qubit, target Qubit) *Q {
	n := q.NumberOfBit()
	g := gate.ControlledNot(n, Index(control...), target.Index())
	q.qb.Apply(g)
	return q
}

func (q *Q) CNOT(control, target Qubit) *Q {
	return q.ControlledNot([]Qubit{control}, target)
}

func (q *Q) CCNOT(control0, control1, target Qubit) *Q {
	return q.ControlledNot([]Qubit{control0, control1}, target)
}

func (q *Q) CCCNOT(control0, control1, control2, target Qubit) *Q {
	return q.ControlledNot([]Qubit{control0, control1, control2}, target)
}

func (q *Q) Toffoli(control0, control1, target Qubit) *Q {
	return q.CCNOT(control0, control1, target)
}

func (q *Q) ControlledZ(control []Qubit, target Qubit) *Q {
	n := q.NumberOfBit()
	g := gate.ControlledZ(n, Index(control...), target.Index())
	q.qb.Apply(g)
	return q
}

func (q *Q) CZ(control, target Qubit) *Q {
	return q.ControlledZ([]Qubit{control}, target)
}

func (q *Q) CCZ(control0, control1, target Qubit) *Q {
	return q.ControlledZ([]Qubit{control0, control1}, target)
}

func (q *Q) ControlledR(theta float64, control []Qubit, target Qubit) *Q {
	n := q.NumberOfBit()
	g := gate.ControlledR(theta, n, Index(control...), target.Index())
	q.qb.Apply(g)
	return q
}

func (q *Q) CR(theta float64, control, target Qubit) *Q {
	return q.ControlledR(theta, []Qubit{control}, target)
}

func (q *Q) ControlledModExp2(a, j, N int, control Qubit, target []Qubit) *Q {
	n := q.NumberOfBit()
	g := gate.ControlledModExp2(n, a, j, N, control.Index(), Index(target...))
	q.qb.Apply(g)
	return q
}

func (q *Q) CModExp2(a, N int, control []Qubit, target []Qubit) *Q {
	for i := 0; i < len(control); i++ {
		q.ControlledModExp2(a, i, N, control[i], target)
	}

	return q
}

func (q *Q) Condition(condition bool, m matrix.Matrix, qb ...Qubit) *Q {
	if condition {
		return q.Apply(m, qb...)
	}

	return q
}

func (q *Q) ConditionX(condition bool, qb ...Qubit) *Q {
	return q.Condition(condition, gate.X(), qb...)
}

func (q *Q) ConditionZ(condition bool, qb ...Qubit) *Q {
	return q.Condition(condition, gate.Z(), qb...)
}

func (q *Q) Swap(qb ...Qubit) *Q {
	n := q.NumberOfBit()
	l := len(qb)

	for i := 0; i < l/2; i++ {
		q0, q1 := qb[i], qb[(l-1)-i]
		g := gate.Swap(n, q0.Index(), q1.Index())
		q.qb.Apply(g)
	}

	return q
}

func (q *Q) QFT(qb ...Qubit) *Q {
	l := len(qb)
	for i := 0; i < l; i++ {
		q.H(qb[i])

		k := 2
		for j := i + 1; j < l; j++ {
			q.CR(gate.Theta(k), qb[j], qb[i])
			k++
		}
	}

	return q
}

func (q *Q) InverseQFT(qb ...Qubit) *Q {
	l := len(qb)
	for i := l - 1; i > -1; i-- {
		k := l - i
		for j := l - 1; j > i; j-- {
			q.CR(-1*gate.Theta(k), qb[j], qb[i])
			k--
		}

		q.H(qb[i])
	}

	return q
}

func (q *Q) InvQFT(qb ...Qubit) *Q {
	return q.InverseQFT(qb...)
}

func (q *Q) Measure(qb ...Qubit) *qubit.Qubit {
	if len(qb) < 1 {
		m := make([]*qubit.Qubit, 0)
		for i := 0; i < q.NumberOfBit(); i++ {
			m = append(m, q.qb.Measure(i))
		}

		return qubit.TensorProduct(m...)
	}

	m := make([]*qubit.Qubit, 0)
	for i := 0; i < len(qb); i++ {
		m = append(m, q.qb.Measure(qb[i].Index()))
	}

	return qubit.TensorProduct(m...)
}

func (q *Q) Clone() *Q {
	if q.qb == nil {
		return &Q{
			qb:   nil,
			Seed: q.Seed,
			Rand: q.Rand,
		}
	}

	return &Q{
		qb:   q.qb.Clone(),
		Seed: q.qb.Seed,
		Rand: q.qb.Rand,
	}
}

func (q *Q) Raw() *qubit.Qubit {
	return q.qb
}

func (q *Q) String() string {
	return q.qb.String()
}

func (q *Q) State(reg ...any) []qubit.State {
	index := make([][]int, 0)
	for _, r := range reg {
		switch r := r.(type) {
		case Qubit:
			index = append(index, []int{r.Index()})
		case []Qubit:
			index = append(index, Index(r...))
		}
	}

	return q.qb.State(index...)
}
