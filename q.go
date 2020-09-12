package q

import (
	"math"
	"strconv"
	"strings"

	"github.com/itsubaki/q/pkg/math/matrix"
	"github.com/itsubaki/q/pkg/math/rand"
	"github.com/itsubaki/q/pkg/quantum/gate"
	"github.com/itsubaki/q/pkg/quantum/qubit"
)

type Qubit int

func (q Qubit) Index() int {
	return int(q)
}

func Index(input ...Qubit) []int {
	index := make([]int, 0)
	for i := range input {
		index = append(index, input[i].Index())
	}

	return index
}

type Q struct {
	internal *qubit.Qubit
	Rand     func(seed ...int64) float64
}

func New() *Q {
	return &Q{
		internal: nil,
		Rand:     rand.Math,
	}
}

func (q *Q) UseCryptoRand() {
	q.Rand = rand.Crypto

	if q.internal != nil {
		q.internal.Rand = rand.Crypto
	}
}

func (q *Q) New(z ...complex128) Qubit {
	if q.internal == nil {
		q.internal = qubit.New(z...)
		q.internal.Rand = q.Rand
		return Qubit(0)
	}

	q.internal.TensorProduct(qubit.New(z...))
	index := q.NumberOfBit() - 1
	return Qubit(index)
}

func (q *Q) Zero() Qubit {
	return q.New(1, 0)
}

func (q *Q) One() Qubit {
	return q.New(0, 1)
}

func (q *Q) ZeroWith(bit int) []Qubit {
	r := make([]Qubit, 0)
	for i := 0; i < bit; i++ {
		r = append(r, q.Zero())
	}

	return r
}

func (q *Q) OneWith(bit int) []Qubit {
	r := make([]Qubit, 0)
	for i := 0; i < bit; i++ {
		r = append(r, q.One())
	}

	return r
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

func (q *Q) Measure(input Qubit) *qubit.Qubit {
	return q.internal.Measure(input.Index())
}

func (q *Q) MeasureAsInt(input ...Qubit) int64 {
	b := q.BinaryString(input...)
	i, err := strconv.ParseInt(b, 2, 0)
	if err != nil {
		panic(err)
	}

	return i
}

func (q *Q) MeasureAsBinary(input ...Qubit) []int {
	b := make([]int, 0)
	for _, i := range input {
		b = append(b, q.Measure(i).Int())
	}

	return b
}

func (q *Q) BinaryString(input ...Qubit) string {
	var sb strings.Builder
	for _, i := range input {
		if q.Measure(i).IsZero() {
			sb.WriteString("0")
			continue
		}

		sb.WriteString("1")
	}

	return sb.String()
}

func (q *Q) I(input ...Qubit) *Q {
	return q.Apply(gate.I(), input...)
}

func (q *Q) H(input ...Qubit) *Q {
	return q.Apply(gate.H(), input...)
}

func (q *Q) X(input ...Qubit) *Q {
	return q.Apply(gate.X(), input...)
}

func (q *Q) Y(input ...Qubit) *Q {
	return q.Apply(gate.Y(), input...)
}

func (q *Q) Z(input ...Qubit) *Q {
	return q.Apply(gate.Z(), input...)
}

func (q *Q) S(input ...Qubit) *Q {
	return q.Apply(gate.S(), input...)
}

func (q *Q) T(input ...Qubit) *Q {
	return q.Apply(gate.T(), input...)
}

func (q *Q) Apply(mat matrix.Matrix, input ...Qubit) *Q {
	if len(input) < 1 {
		q.internal.Apply(mat)
		return q
	}

	index := Index(input...)

	g := gate.I()
	if index[0] == 0 {
		g = mat
	}

	for i := 1; i < q.NumberOfBit(); i++ {
		found := false
		for j := range index {
			if i == index[j] {
				found = true
				break
			}
		}

		if found {
			g = g.TensorProduct(mat)
			continue
		}

		g = g.TensorProduct(gate.I())
	}

	q.internal.Apply(g)
	return q
}

func (q *Q) ControlledR(control []Qubit, target Qubit, k int) *Q {
	cr := gate.ControlledR(q.NumberOfBit(), Index(control...), target.Index(), k)
	q.internal.Apply(cr)
	return q
}

func (q *Q) CR(control, target Qubit, k int) *Q {
	return q.ControlledR([]Qubit{control}, target, k)
}

func (q *Q) ControlledZ(control []Qubit, target Qubit) *Q {
	cnot := gate.ControlledZ(q.NumberOfBit(), Index(control...), target.Index())
	q.internal.Apply(cnot)
	return q
}

func (q *Q) CZ(control, target Qubit) *Q {
	return q.ControlledZ([]Qubit{control}, target)
}

func (q *Q) CCZ(control0, control1, target Qubit) *Q {
	return q.ControlledZ([]Qubit{control0, control1}, target)
}

func (q *Q) ControlledNot(control []Qubit, target Qubit) *Q {
	cnot := gate.ControlledNot(q.NumberOfBit(), Index(control...), target.Index())
	q.internal.Apply(cnot)
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

func (q *Q) ConditionX(condition bool, input ...Qubit) *Q {
	if condition {
		return q.X(input...)
	}

	return q
}

func (q *Q) ConditionZ(condition bool, input ...Qubit) *Q {
	if condition {
		return q.Z(input...)
	}

	return q
}

func (q *Q) CModExp2(a, j, N int, control Qubit, target ...Qubit) *Q {
	n := q.NumberOfBit()
	g := gate.CModExp2(n, a, j, N, control.Index(), Index(target...))
	q.internal.Apply(g)
	return q
}

func (q *Q) Swap(input ...Qubit) *Q {
	n := q.NumberOfBit()
	l := len(input)

	for i := 0; i < l/2; i++ {
		q0, q1 := input[i], input[(l-1)-i]
		swap := gate.Swap(n, q0.Index(), q1.Index())
		q.internal.Apply(swap)
	}

	return q
}

func (q *Q) QFT(input ...Qubit) *Q {
	l := len(input)
	for i := 0; i < l; i++ {
		q.H(input[i])

		k := 2
		for j := i + 1; j < l; j++ {
			q.CR(input[j], input[i], k)
			k++
		}
	}

	return q
}

func (q *Q) InverseQFT(input ...Qubit) *Q {
	l := len(input)
	for i := l - 1; i > -1; i-- {
		k := l - i
		for j := l - 1; j > i; j-- {
			q.CR(input[j], input[i], k)
			k--
		}

		q.H(input[i])
	}

	return q
}

func (q *Q) Estimate(input Qubit, loop ...int) *qubit.Qubit {
	c0, c1, limit := 0, 0, 1000
	if len(loop) > 0 {
		limit = loop[0]
	}

	for i := 0; i < limit; i++ {
		m := q.Clone().Measure(input)

		if m.IsZero() {
			c0++
			continue
		}

		c1++
	}

	z := math.Sqrt(float64(c0) / float64(limit))
	o := math.Sqrt(float64(c1) / float64(limit))

	return qubit.New(complex(z, 0), complex(o, 0))
}

func (q *Q) Clone() *Q {
	if q.internal == nil {
		return &Q{
			internal: nil,
			Rand:     q.Rand,
		}
	}

	return &Q{
		internal: q.internal.Clone(),
		Rand:     q.internal.Rand,
	}
}

func (q *Q) String() string {
	return q.internal.String()
}
