package q

import (
	"fmt"
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
	internal *qubit.Qubit
	Seed     []int64
	Rand     func(seed ...int64) float64
}

func New() *Q {
	return &Q{
		internal: nil,
		Rand:     rand.Crypto,
	}
}

func (q *Q) New(v ...complex128) Qubit {
	if q.internal == nil {
		q.internal = qubit.New(v...)
		q.internal.Seed = q.Seed
		q.internal.Rand = q.Rand
		return Qubit(0)
	}

	q.internal.TensorProduct(qubit.New(v...))
	index := q.NumberOfBit() - 1
	return Qubit(index)
}

func (q *Q) Zero() Qubit {
	return q.New(1, 0)
}

func (q *Q) One() Qubit {
	return q.New(0, 1)
}

func (q *Q) ZeroWith(n int) []Qubit {
	r := make([]Qubit, 0)
	for i := 0; i < n; i++ {
		r = append(r, q.Zero())
	}

	return r
}

func (q *Q) OneWith(n int) []Qubit {
	r := make([]Qubit, 0)
	for i := 0; i < n; i++ {
		r = append(r, q.One())
	}

	return r
}

func (q *Q) ZeroLog2(N int) []Qubit {
	n := int(math.Log2(float64(N))) + 1
	return q.ZeroWith(n)
}

func (q *Q) NumberOfBit() int {
	return q.internal.NumberOfBit()
}

func (q *Q) Amplitude() []complex128 {
	return q.internal.Amplitude()
}

func (q *Q) Probability() []float64 {
	return q.internal.Probability()
}

func (q *Q) Reset(qb ...Qubit) {
	for i := range qb {
		if q.Measure(qb[i]).IsOne() {
			q.X(qb[i])
		}
	}
}

func (q *Q) I(qb ...Qubit) *Q {
	return q.Apply(gate.I(), qb...)
}

func (q *Q) H(qb ...Qubit) *Q {
	return q.Apply(gate.H(), qb...)
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

func (q *Q) S(qb ...Qubit) *Q {
	return q.Apply(gate.S(), qb...)
}

func (q *Q) T(qb ...Qubit) *Q {
	return q.Apply(gate.T(), qb...)
}

func (q *Q) U(theta, phi, lambda float64, qb ...Qubit) *Q {
	return q.Apply(gate.U(theta, phi, lambda), qb...)
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
		q.internal.Apply(m)
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

	q.internal.Apply(g)
	return q
}

func (q *Q) C(m matrix.Matrix, control, target Qubit) *Q {
	n := q.NumberOfBit()
	g := gate.C(m, n, control.Index(), target.Index())
	q.internal.Apply(g)
	return q
}

func (q *Q) ControlledR(k int, control []Qubit, target Qubit) *Q {
	n := q.NumberOfBit()
	g := gate.ControlledR(k, n, Index(control...), target.Index())
	q.internal.Apply(g)
	return q
}

func (q *Q) CR(k int, control, target Qubit) *Q {
	return q.ControlledR(k, []Qubit{control}, target)
}

func (q *Q) InvCR(k int, control, target Qubit) *Q {
	n := q.NumberOfBit()
	g := gate.ControlledR(k, n, []int{control.Index()}, target.Index()).Dagger()
	q.internal.Apply(g)
	return q
}

func (q *Q) ControlledZ(control []Qubit, target Qubit) *Q {
	n := q.NumberOfBit()
	g := gate.ControlledZ(n, Index(control...), target.Index())
	q.internal.Apply(g)
	return q
}

func (q *Q) CZ(control, target Qubit) *Q {
	return q.ControlledZ([]Qubit{control}, target)
}

func (q *Q) CCZ(control0, control1, target Qubit) *Q {
	return q.ControlledZ([]Qubit{control0, control1}, target)
}

func (q *Q) ControlledNot(control []Qubit, target Qubit) *Q {
	n := q.NumberOfBit()
	g := gate.ControlledNot(n, Index(control...), target.Index())
	q.internal.Apply(g)
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

func (q *Q) ConditionX(condition bool, qb ...Qubit) *Q {
	if condition {
		return q.X(qb...)
	}

	return q
}

func (q *Q) ConditionZ(condition bool, qb ...Qubit) *Q {
	if condition {
		return q.Z(qb...)
	}

	return q
}

func (q *Q) ControlledModExp2(a, j, N int, control Qubit, target []Qubit) *Q {
	n := q.NumberOfBit()
	g := gate.CModExp2(n, a, j, N, control.Index(), Index(target...))
	q.internal.Apply(g)
	return q
}

func (q *Q) CModExp2(a, N int, control []Qubit, target []Qubit) *Q {
	for i := 0; i < len(control); i++ {
		q.ControlledModExp2(a, i, N, control[i], target)
	}

	return q
}

func (q *Q) Swap(qb ...Qubit) *Q {
	n := q.NumberOfBit()
	l := len(qb)

	for i := 0; i < l/2; i++ {
		q0, q1 := qb[i], qb[(l-1)-i]
		g := gate.Swap(n, q0.Index(), q1.Index())
		q.internal.Apply(g)
	}

	return q
}

func (q *Q) QFT(qb ...Qubit) *Q {
	l := len(qb)
	for i := 0; i < l; i++ {
		q.H(qb[i])

		k := 2
		for j := i + 1; j < l; j++ {
			q.CR(k, qb[j], qb[i])
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
			q.InvCR(k, qb[j], qb[i])
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
			m = append(m, q.internal.Measure(i))
		}

		return qubit.TensorProduct(m...)
	}

	m := make([]*qubit.Qubit, 0)
	for i := 0; i < len(qb); i++ {
		m = append(m, q.internal.Measure(qb[i].Index()))
	}

	return qubit.TensorProduct(m...)
}

func (q *Q) Clone() *Q {
	if q.internal == nil {
		return &Q{
			internal: nil,
			Seed:     q.Seed,
			Rand:     q.Rand,
		}
	}

	return &Q{
		internal: q.internal.Clone(),
		Seed:     q.internal.Seed,
		Rand:     q.internal.Rand,
	}
}

func (q *Q) Raw() *qubit.Qubit {
	return q.internal
}

func (q *Q) String() string {
	return q.internal.String()
}

func (q *Q) State(reg ...interface{}) []qubit.State {
	index := make([][]int, 0)
	for _, r := range reg {
		switch r := r.(type) {
		case Qubit:
			index = append(index, []int{r.Index()})
		case []Qubit:
			index = append(index, Index(r...))
		default:
			panic(fmt.Sprintf("invalid type %T", r))
		}
	}

	return q.internal.State(index...)
}
