package q

import (
	"math"

	"github.com/itsubaki/q/gate"
	"github.com/itsubaki/q/matrix"
	"github.com/itsubaki/q/qubit"
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

func (q *Q) QFT(qb ...*Qubit) *Q {
	dim := len(qb)

	for i := 0; i < dim; i++ {
		q.H(qb[i])

		k := 2
		for j := i + 1; j < dim; j++ {
			q.CR(qb[j], qb[i], k)
			k++
		}
	}

	for i := 0; i < dim/2; i++ {
		q.Swap(qb[i], qb[dim-1-i])
	}

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
	q.CNOT(q0, q1)
	q.CNOT(q1, q0)
	q.CNOT(q0, q1)
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
	limit := 10000
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
