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

func (q *Q) H(input ...*Qubit) {
	q.Apply(gate.H(), input...)
}

func (q *Q) X(input ...*Qubit) {
	q.Apply(gate.X(), input...)
}

func (q *Q) Y(input ...*Qubit) {
	q.Apply(gate.Y(), input...)
}

func (q *Q) Z(input ...*Qubit) {
	q.Apply(gate.Z(), input...)
}

func (q *Q) S(input ...*Qubit) {
	q.Apply(gate.S(), input...)
}

func (q *Q) T(input ...*Qubit) {
	q.Apply(gate.T(), input...)
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

func (q *Q) CNOT(controll *Qubit, target *Qubit) {
	bit := q.qubit.NumberOfBit()
	cnot := gate.CNOT(bit, controll.Index, target.Index)
	q.qubit.Apply(cnot)
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
