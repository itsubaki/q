package q

import (
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
	q.Apply(gate.H(), input...)
	return q
}

func (q *Q) X(input ...*Qubit) *Q {
	q.Apply(gate.X(), input...)
	return q
}

func (q *Q) Y(input ...*Qubit) *Q {
	q.Apply(gate.Y(), input...)
	return q
}

func (q *Q) Z(input ...*Qubit) *Q {
	q.Apply(gate.Z(), input...)
	return q
}

func (q *Q) S(input ...*Qubit) *Q {
	q.Apply(gate.S(), input...)
	return q
}

func (q *Q) T(input ...*Qubit) *Q {
	q.Apply(gate.T(), input...)
	return q
}

func (q *Q) Apply(mat matrix.Matrix, input ...*Qubit) {
	g := gate.I()
	if input[0].Index == 0 {
		g = mat
	}

	index := []int{}
	for i := range input {
		index = append(index, input[i].Index)
	}

	for i := 1; i < q.qubit.NumberOfBit(); i++ {
		for j := range index {
			if i == index[j] {
				g = g.TensorProduct(mat)
				continue
			}
			g = g.TensorProduct(gate.I())
		}
	}

	q.qubit.Apply(g)
}

func (q *Q) CNOT(controll *Qubit, target *Qubit) *Q {
	bit := q.qubit.NumberOfBit()
	cnot := gate.CNOT(bit, controll.Index, target.Index)
	q.qubit.Apply(cnot)

	return q
}

func (q *Q) ConditionX(condition bool, input ...*Qubit) *Q {
	if condition {
		q.X(input...)
		return q
	}

	return q
}

func (q *Q) ConditionZ(condition bool, input ...*Qubit) *Q {
	if condition {
		q.Z(input...)
		return q
	}

	return q
}

func (q *Q) Measure(input ...*Qubit) *qubit.Qubit {
	if len(input) < 1 {
		return q.qubit.Measure()
	}

	return q.qubit.MeasureAt(input[0].Index)
}

func (q *Q) Probability() []qubit.Probability {
	return q.qubit.Probability()
}
