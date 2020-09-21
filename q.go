package q

import (
	"fmt"
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

func Index(qb ...Qubit) []int {
	index := make([]int, 0)
	for i := range qb {
		index = append(index, qb[i].Index())
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
		Rand:     rand.Crypto,
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

func (q *Q) Measure(qb Qubit) *qubit.Qubit {
	return q.internal.Measure(qb.Index())
}

func (q *Q) MeasureAsInt(qb ...Qubit) int64 {
	b := q.BinaryString(qb...)
	i, err := strconv.ParseInt(b, 2, 0)
	if err != nil {
		panic(err)
	}

	return i
}

func (q *Q) MeasureAsBinary(qb ...Qubit) []int {
	b := make([]int, 0)
	for _, i := range qb {
		b = append(b, q.Measure(i).Int())
	}

	return b
}

func (q *Q) BinaryString(qb ...Qubit) string {
	var sb strings.Builder
	for _, i := range qb {
		if q.Measure(i).IsZero() {
			sb.WriteString("0")
			continue
		}

		sb.WriteString("1")
	}

	return sb.String()
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

func (q *Q) Apply(mat matrix.Matrix, qb ...Qubit) *Q {
	if len(qb) < 1 {
		q.internal.Apply(mat)
		return q
	}

	index := Index(qb...)

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

func (q *Q) ControlledModExp2(N, a, j int, c Qubit, t []Qubit) *Q {
	n := q.NumberOfBit()
	g := gate.CModExp2(n, N, a, j, c.Index(), Index(t...))
	q.internal.Apply(g)
	return q
}

func (q *Q) CModExp2(N, a int, c []Qubit, t []Qubit) *Q {
	for j := 0; j < len(c); j++ {
		q.ControlledModExp2(N, a, j, c[j], t)
	}

	return q
}

func (q *Q) Swap(qb ...Qubit) *Q {
	n := q.NumberOfBit()
	l := len(qb)

	for i := 0; i < l/2; i++ {
		q0, q1 := qb[i], qb[(l-1)-i]
		swap := gate.Swap(n, q0.Index(), q1.Index())
		q.internal.Apply(swap)
	}

	return q
}

func (q *Q) QFT(qb ...Qubit) *Q {
	l := len(qb)
	for i := 0; i < l; i++ {
		q.H(qb[i])

		k := 2
		for j := i + 1; j < l; j++ {
			q.CR(qb[j], qb[i], k)
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
			q.CR(qb[j], qb[i], k)
			k--
		}

		q.H(qb[i])
	}

	return q
}

func (q *Q) InvQFT(qb ...Qubit) *Q {
	return q.InverseQFT(qb...)
}

func (q *Q) Estimate(qb Qubit, loop ...int) *qubit.Qubit {
	c0, c1, limit := 0, 0, 1000
	if len(loop) > 0 {
		limit = loop[0]
	}

	for i := 0; i < limit; i++ {
		m := q.Clone().Measure(qb)

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

func (q *Q) StringRegister(reg ...[]Qubit) string {
	return q.StringRegisterWith(" ", reg...)
}

func (q *Q) StringRegisterln(reg ...[]Qubit) string {
	return q.StringRegisterWith("\n", reg...)
}

func (q *Q) StringRegisterWith(delimiter string, reg ...[]Qubit) string {
	if len(reg) == 0 {
		return ""
	}

	var sb strings.Builder
	binf := fmt.Sprintf("%s%s%s", "%0", strconv.Itoa(q.NumberOfBit()), "s")
	for i, a := range q.Amplitude() {
		if a == 0 {
			continue
		}
		if math.Abs(real(a)) < 1e-13 {
			a = complex(0, imag(a))
		}
		if math.Abs(imag(a)) < 1e-13 {
			a = complex(real(a), 0)
		}
		sb.WriteString(fmt.Sprintf("%.2v", a))

		bin := fmt.Sprintf(binf, strconv.FormatInt(int64(i), 2))
		for _, r := range reg {
			var rbin strings.Builder
			for _, qb := range r {
				idx := qb.Index()
				rbin.WriteString(bin[idx : idx+1])
			}

			rstr := rbin.String()
			rint, err := strconv.ParseInt(rstr, 2, 0)
			if err != nil {
				panic(fmt.Sprintf("parse int bin=%s, rstr=%s", bin, rstr))
			}

			sb.WriteString(fmt.Sprintf("|%d>", rint))
		}
		sb.WriteString(delimiter)
	}

	return sb.String()
}
