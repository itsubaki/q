package qubit

import (
	"fmt"
	"math"
	"math/cmplx"

	"github.com/itsubaki/q/pkg/math/matrix"
	"github.com/itsubaki/q/pkg/math/number"
	"github.com/itsubaki/q/pkg/math/rand"
	"github.com/itsubaki/q/pkg/math/vector"
)

type Qubit struct {
	vector vector.Vector
	Seed   []int64
	Rand   func(seed ...int64) float64
}

func New(z ...complex128) *Qubit {
	q := &Qubit{
		vector: vector.New(z...),
		Rand:   rand.Crypto,
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

func (q *Qubit) InnerProduct(qb *Qubit) complex128 {
	return q.vector.InnerProduct(qb.vector)
}

func (q *Qubit) OuterProduct(qb *Qubit) matrix.Matrix {
	return q.vector.OuterProduct(qb.vector)
}

func (q *Qubit) Dimension() int {
	return q.vector.Dimension()
}

func (q *Qubit) Clone() *Qubit {
	return &Qubit{
		vector: q.vector.Clone(),
		Seed:   q.Seed,
		Rand:   q.Rand,
	}
}

func (q *Qubit) Fidelity(qb *Qubit) float64 {
	p0 := qb.Probability()
	p1 := q.Probability()

	var sum float64
	for i := 0; i < len(p0); i++ {
		sum = sum + math.Sqrt(p0[i]*p1[i])
	}

	return sum
}

func (q *Qubit) TraceDistance(qb *Qubit) float64 {
	p0 := qb.Probability()
	p1 := q.Probability()

	var sum float64
	for i := 0; i < len(p0); i++ {
		sum = sum + math.Abs(p0[i]-p1[i])
	}

	return sum / 2
}

func (q *Qubit) Equals(qb *Qubit, eps ...float64) bool {
	return q.vector.Equals(qb.vector, eps...)
}

func (q *Qubit) TensorProduct(qb *Qubit) *Qubit {
	q.vector = q.vector.TensorProduct(qb.vector)
	return q
}

func (q *Qubit) Apply(m ...matrix.Matrix) *Qubit {
	for _, mm := range m {
		q.vector = q.vector.Apply(mm)
	}

	return q
}

func (q *Qubit) Normalize() *Qubit {
	sum := number.Sum(q.Probability())
	z := 1 / math.Sqrt(sum)
	q.vector = q.vector.Mul(complex(z, 0))
	return q
}

func (q *Qubit) Amplitude() []complex128 {
	return q.vector.Complex()
}

func (q *Qubit) Probability() []float64 {
	p := make([]float64, 0)
	for _, a := range q.Amplitude() {
		p = append(p, math.Pow(cmplx.Abs(a), 2))
	}

	return p
}

func (q *Qubit) Measure(bit int) *Qubit {
	index, p := q.ProbabilityZeroAt(bit)
	r := q.Rand(q.Seed...)

	sum := number.Sum(p)
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
	dim := q.Dimension()
	den := int(math.Pow(2, float64(bit+1)))
	div := dim / den

	p := q.Probability()
	index, prob := make([]int, 0), make([]float64, 0)
	for i := 0; i < dim; i++ {
		index, prob = append(index, i), append(prob, p[i])

		if len(p) == dim/2 {
			break
		}

		if (i+1)%div == 0 {
			i = i + div
		}
	}

	return index, prob
}

func (q *Qubit) ProbabilityOneAt(bit int) ([]int, []float64) {
	z, _ := q.ProbabilityZeroAt(bit)

	one := make([]int, 0)
	for i := range q.vector {
		found := false
		for _, zi := range z {
			if i == zi {
				found = true
				break
			}
		}

		if !found {
			one = append(one, i)
		}
	}

	p := q.Probability()
	index, prob := make([]int, 0), make([]float64, 0)
	for _, i := range one {
		index, prob = append(index, i), append(prob, p[i])
	}

	return index, prob
}

func (q *Qubit) String() string {
	return fmt.Sprintf("%v", q.vector)
}

func TensorProduct(qb ...*Qubit) *Qubit {
	q := qb[0]
	for i := 1; i < len(qb); i++ {
		q = q.TensorProduct(qb[i])
	}

	return q
}
