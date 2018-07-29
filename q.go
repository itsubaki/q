package q

import "github.com/itsubaki/q/qubit"

type Q struct {
}

func New() *Q {
	return &Q{}
}

func (q *Q) Add(input ...*qubit.Qubit) {

}

func (q *Q) H(input ...*qubit.Qubit) {

}

func (q *Q) CNOT(controll *qubit.Qubit, target *qubit.Qubit) {

}

func (q *Q) Measure(input ...*qubit.Qubit) *qubit.Qubit {
	return &qubit.Qubit{}
}

func (q *Q) Probability() []qubit.Probability {
	return []qubit.Probability{}
}
