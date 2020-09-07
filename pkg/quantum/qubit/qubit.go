package qubit

import (
	"fmt"
	"math"
	"math/cmplx"

	"github.com/itsubaki/q/pkg/math/matrix"
	"github.com/itsubaki/q/pkg/math/rand"
	"github.com/itsubaki/q/pkg/math/vector"
)

type Qubit struct {
	vector vector.Vector
	Rand   func(seed ...int64) float64
}

func New(z ...complex128) *Qubit {
	v := vector.New()
	for _, zi := range z {
		v = append(v, zi)
	}

	q := &Qubit{
		vector: v,
		Rand:   rand.Math,
	}

	q.Normalize()
	return q
}

func Zero(bit ...int) *Qubit {
	v := vector.TensorProductN(vector.Vector{1, 0}, bit...)
	return New(v.Complex()...)
}

func One(bit ...int) *Qubit {
	v := vector.TensorProductN(vector.Vector{0, 1}, bit...)
	return New(v.Complex()...)
}

func (q *Qubit) NumberOfBit() int {
	dim := float64(q.Dimension())
	log := math.Log2(dim)
	return int(log)
}

func (q *Qubit) IsZero(eps ...float64) bool {
	return q.Equals(Zero(), eps...)
}

func (q *Qubit) IsOne(eps ...float64) bool {
	return q.Equals(One(), eps...)
}

func (q *Qubit) InnerProduct(q0 *Qubit) complex128 {
	return q.vector.InnerProduct(q0.vector)
}

func (q *Qubit) OuterProduct(q0 *Qubit) matrix.Matrix {
	return q.vector.OuterProduct(q0.vector)
}

func (q *Qubit) Dimension() int {
	return q.vector.Dimension()
}

func (q *Qubit) Clone() *Qubit {
	return &Qubit{
		vector: q.vector.Clone(),
		Rand:   q.Rand,
	}
}

func (q *Qubit) Fidelity(q0 *Qubit) float64 {
	p0 := q0.Probability()
	p1 := q.Probability()

	var sum float64
	for i := 0; i < len(p0); i++ {
		sum = sum + math.Sqrt(p0[i]*p1[i])
	}

	return sum
}

func (q *Qubit) TraceDistance(q0 *Qubit) float64 {
	p0 := q0.Probability()
	p1 := q.Probability()

	var sum float64
	for i := 0; i < len(p0); i++ {
		sum = sum + math.Abs(p0[i]-p1[i])
	}

	return sum / 2
}

func (q *Qubit) Equals(q0 *Qubit, eps ...float64) bool {
	return q.vector.Equals(q0.vector, eps...)
}

func (q *Qubit) TensorProduct(q0 *Qubit) *Qubit {
	q.vector = q.vector.TensorProduct(q0.vector)
	return q
}

func (q *Qubit) Apply(m matrix.Matrix) *Qubit {
	q.vector = q.vector.Apply(m)
	return q
}

func (q *Qubit) Normalize() *Qubit {
	var sum float64
	for _, amp := range q.vector {
		sum = sum + math.Pow(cmplx.Abs(amp), 2)
	}

	z := 1 / math.Sqrt(sum)
	q.vector = q.vector.Mul(complex(z, 0))

	return q
}

func (q *Qubit) Amplitude() []complex128 {
	a := make([]complex128, 0)
	for _, amp := range q.vector {
		a = append(a, amp)
	}

	return a
}

func (q *Qubit) Probability() []float64 {
	list := make([]float64, 0)
	for _, amp := range q.vector {
		p := math.Pow(cmplx.Abs(amp), 2)
		list = append(list, p)
	}

	return list
}

func (q *Qubit) Measure(bit int) *Qubit {
	index, p := q.ProbabilityZeroAt(bit)
	r := q.Rand()

	var sum float64
	for _, pp := range p {
		sum = sum + pp
	}

	if r > sum {
		for _, i := range index {
			q.vector[i] = complex(0, 0)
		}

		q.Normalize()
		return One()
	}

	one := make([]int, 0)
	for i := range q.vector {
		found := false
		for _, ix := range index {
			if i == ix {
				found = true
				break
			}
		}

		if !found {
			one = append(one, i)
		}
	}

	for _, i := range one {
		q.vector[i] = complex(0, 0)
	}

	q.Normalize()
	return Zero()
}

func (q *Qubit) ProbabilityZeroAt(bit int) ([]int, []float64) {
	p := make([]float64, 0)
	index := make([]int, 0)

	dim := q.Dimension()
	den := int(math.Pow(2, float64(bit+1)))
	div := dim / den

	prob := q.Probability()
	for i := 0; i < dim; i++ {
		p = append(p, prob[i])
		index = append(index, i)

		if len(p) == dim/2 {
			break
		}

		if (i+1)%div == 0 {
			i = i + div
		}
	}

	return index, p
}

func (q *Qubit) ProbabilityOneAt(bit int) ([]int, []float64) {
	p := make([]float64, 0)
	index := make([]int, 0)

	zi, _ := q.ProbabilityZeroAt(bit)
	one := make([]int, 0)
	for i := range q.vector {
		found := false
		for _, zii := range zi {
			if i == zii {
				found = true
				break
			}
		}

		if !found {
			one = append(one, i)
		}
	}

	for _, i := range one {
		p = append(p, q.Probability()[i])
		index = append(index, i)
	}

	return index, p
}

func (q *Qubit) Int() int {
	n := q.NumberOfBit()
	if n != 1 {
		panic(fmt.Sprintf("invalid number of bit=%d", n))
	}

	if q.IsZero() {
		return 0
	}

	if q.IsOne() {
		return 1
	}

	panic("invalid state")
}

func TensorProduct(q ...*Qubit) *Qubit {
	q1 := q[0]
	for i := 1; i < len(q); i++ {
		q1 = q1.TensorProduct(q[i])
	}

	return q1
}
