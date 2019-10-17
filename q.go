package q

import (
	"math"

	"github.com/itsubaki/q/pkg/circuit/gate"
	"github.com/itsubaki/q/pkg/circuit/qubit"
	"github.com/itsubaki/q/pkg/math/matrix"
)

type Q struct {
	qubit *qubit.Qubit
}

type Qubit struct {
	Index int
}

func Index(input []*Qubit) []int {
	index := make([]int, 0)
	for i := range input {
		index = append(index, input[i].Index)
	}

	return index
}

func New() *Q {
	return &Q{}
}

func (q *Q) New(z ...complex128) *Qubit {
	if q.qubit == nil {
		q.qubit = qubit.New(z...)
		return &Qubit{Index: 0}
	}

	q.qubit.TensorProduct(qubit.New(z...))
	index := q.NumberOfBit() - 1
	return &Qubit{Index: index}
}

func (q *Q) Zero() *Qubit {
	return q.New(1, 0)
}

func (q *Q) One() *Qubit {
	return q.New(0, 1)
}

func (q *Q) Probability() []float64 {
	return q.qubit.Probability()
}

func (q *Q) Measure(input ...*Qubit) *qubit.Qubit {
	if len(input) < 1 {
		return q.qubit.Measure()
	}

	return q.qubit.Measure(input[0].Index)
}

func (q *Q) NumberOfBit() int {
	return q.qubit.NumberOfBit()
}

func (q *Q) Clone() *Q {
	return &Q{qubit: q.qubit.Clone()}
}

func (q *Q) H(input ...*Qubit) *Q {
	return q.Apply(gate.H(), input...)
}

func (q *Q) X(input ...*Qubit) *Q {
	return q.Apply(gate.X(), input...)
}

func (q *Q) Y(input ...*Qubit) *Q {
	return q.Apply(gate.Y(), input...)
}

func (q *Q) Z(input ...*Qubit) *Q {
	return q.Apply(gate.Z(), input...)
}

func (q *Q) S(input ...*Qubit) *Q {
	return q.Apply(gate.S(), input...)
}

func (q *Q) T(input ...*Qubit) *Q {
	return q.Apply(gate.T(), input...)
}

func (q *Q) Apply(mat matrix.Matrix, input ...*Qubit) *Q {
	index := make([]int, 0)
	for i := range input {
		index = append(index, input[i].Index)
	}

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

	q.qubit.Apply(g)
	return q
}

func (q *Q) ControlledR(control []*Qubit, target *Qubit, k int) *Q {
	cr := gate.ControlledR(q.NumberOfBit(), Index(control), target.Index, k)
	q.qubit.Apply(cr)
	return q
}

func (q *Q) CR(control *Qubit, target *Qubit, k int) *Q {
	return q.ControlledR([]*Qubit{control}, target, k)
}

func (q *Q) ControlledZ(control []*Qubit, target *Qubit) *Q {
	cnot := gate.ControlledZ(q.NumberOfBit(), Index(control), target.Index)
	q.qubit.Apply(cnot)
	return q
}

func (q *Q) CZ(control *Qubit, target *Qubit) *Q {
	return q.ControlledZ([]*Qubit{control}, target)
}

func (q *Q) ControlledNot(control []*Qubit, target *Qubit) *Q {
	cnot := gate.ControlledNot(q.NumberOfBit(), Index(control), target.Index)
	q.qubit.Apply(cnot)
	return q
}

func (q *Q) CNOT(control *Qubit, target *Qubit) *Q {
	return q.ControlledNot([]*Qubit{control}, target)
}

func (q *Q) QFT() *Q {
	q.qubit.Apply(gate.QFT(q.NumberOfBit()))
	return q
}

func (q *Q) InverseQFT() *Q {
	q.qubit.Apply(gate.QFT(q.NumberOfBit()).Dagger())
	return q
}

func (q *Q) ConditionX(condition bool, input ...*Qubit) *Q {
	if condition {
		return q.X(input...)
	}

	return q
}

func (q *Q) ConditionZ(condition bool, input ...*Qubit) *Q {
	if condition {
		return q.Z(input...)
	}

	return q
}

func (q *Q) Swap(q0, q1 *Qubit) *Q {
	swap := gate.Swap(q.NumberOfBit(), q0.Index, q1.Index)
	q.qubit.Apply(swap)
	return q
}

func (q *Q) Estimate(input *Qubit, loop ...int) *qubit.Qubit {
	limit := 1000
	if len(loop) > 0 {
		limit = loop[0]
	}

	c := []int{0, 0}
	for i := 0; i < limit; i++ {
		m := q.Clone().Measure(input)

		if m.IsZero() {
			c[0]++
		} else {
			c[1]++
		}
	}

	z := complex(math.Sqrt(float64(c[0])/float64(limit)), 0)
	o := complex(math.Sqrt(float64(c[1])/float64(limit)), 0)

	return qubit.New(z, o)
}
