package q

import (
	"math"

	"github.com/itsubaki/q/circuit/gate"
	"github.com/itsubaki/q/circuit/qubit"
	"github.com/itsubaki/q/math/matrix"
)

type Q struct {
	qubit *qubit.Qubit
}

type Qubit struct {
	Index int
}

func index(input []*Qubit) []int {
	index := []int{}
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
	index := q.qubit.NumberOfBit() - 1
	return &Qubit{Index: index}
}

func (q *Q) Zero() *Qubit {
	return q.New(1, 0)
}

func (q *Q) One() *Qubit {
	return q.New(0, 1)
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
	index := []int{}
	for i := range input {
		index = append(index, input[i].Index)
	}

	g := gate.I()
	if index[0] == 0 {
		g = mat
	}

	for i := 1; i < q.qubit.NumberOfBit(); i++ {
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

func (q *Q) ControlledR(controll []*Qubit, target *Qubit, k int) *Q {
	bit := q.qubit.NumberOfBit()
	cr := gate.ControlledR(bit, index(controll), target.Index, k)

	q.qubit.Apply(cr)
	return q
}

func (q *Q) CR(controll *Qubit, target *Qubit, k int) *Q {
	return q.ControlledR([]*Qubit{controll}, target, k)
}

func (q *Q) ControlledZ(controll []*Qubit, target *Qubit) *Q {
	bit := q.qubit.NumberOfBit()
	cnot := gate.ControlledZ(bit, index(controll), target.Index)

	q.qubit.Apply(cnot)
	return q
}

func (q *Q) CZ(controll *Qubit, target *Qubit) *Q {
	return q.ControlledZ([]*Qubit{controll}, target)
}

func (q *Q) ControlledNot(controll []*Qubit, target *Qubit) *Q {
	bit := q.qubit.NumberOfBit()
	cnot := gate.ControlledNot(bit, index(controll), target.Index)

	q.qubit.Apply(cnot)
	return q
}

func (q *Q) CNOT(controll *Qubit, target *Qubit) *Q {
	return q.ControlledNot([]*Qubit{controll}, target)
}

func (q *Q) QFT() *Q {
	bit := q.qubit.NumberOfBit()
	q.qubit.Apply(gate.QFT(bit))
	return q
}

func (q *Q) InverseQFT() *Q {
	bit := q.qubit.NumberOfBit()
	q.qubit.Apply(gate.QFT(bit).Dagger())
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
	bit := q.qubit.NumberOfBit()
	swap := gate.Swap(bit, q0.Index, q1.Index)
	q.qubit.Apply(swap)
	return q
}

func (q *Q) Measure(input ...*Qubit) *qubit.Qubit {
	if len(input) < 1 {
		return q.qubit.Measure()
	}
	return q.qubit.Measure(input[0].Index)
}

func (q *Q) Probability() []float64 {
	return q.qubit.Probability()
}

func (q *Q) Estimate(input *Qubit, loop ...int) *qubit.Qubit {
	limit := 1000
	if len(loop) > 0 {
		limit = loop[0]
	}

	c := []int{0, 0}
	for i := 0; i < limit; i++ {
		clone := q.qubit.Clone()
		m := clone.Measure(input.Index)

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
