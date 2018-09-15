package qubit

import (
	"math"
	"math/cmplx"
	"math/rand"
	"time"

	"github.com/itsubaki/q/math/matrix"
	v "github.com/itsubaki/q/math/vector"
)

type Qubit struct {
	v v.Vector
}

func New(z ...complex128) *Qubit {
	v := v.Vector{}
	for _, zi := range z {
		v = append(v, zi)
	}
	q := &Qubit{v}
	q.Normalize()
	return q
}

func Zero(bit ...int) *Qubit {
	return &Qubit{v.TensorProductN(v.Vector{1, 0}, bit...)}
}

func One(bit ...int) *Qubit {
	return &Qubit{v.TensorProductN(v.Vector{0, 1}, bit...)}
}

func (q *Qubit) NumberOfBit() int {
	dim := float64(q.v.Dimension())
	log := math.Log2(dim)
	return int(log)
}

func (q *Qubit) IsZero(eps ...float64) bool {
	return q.Equals(Zero(), eps...)
}

func (q *Qubit) IsOne(eps ...float64) bool {
	return q.Equals(One(), eps...)
}

func (q *Qubit) Clone() *Qubit {
	return &Qubit{q.v.Clone()}
}

func (q *Qubit) Fidelity(q0 *Qubit) float64 {
	p0 := q0.Probability()
	p1 := q.Probability()

	var sum float64
	for i := 0; i < len(p0); i++ {
		sum = sum + math.Sqrt(float64(p0[i])*float64(p1[i]))
	}

	return sum
}

func (q *Qubit) TraceDistance(q0 *Qubit) float64 {
	p0 := q0.Probability()
	p1 := q.Probability()

	var sum float64
	for i := 0; i < len(p0); i++ {
		sum = sum + math.Abs(float64(p0[i]-p1[i]))
	}

	return sum / 2
}

func (q *Qubit) Equals(q0 *Qubit, eps ...float64) bool {
	return q.v.Equals(q0.v, eps...)
}

func (q *Qubit) TensorProduct(q0 *Qubit) *Qubit {
	q.v = q.v.TensorProduct(q0.v)
	return q
}

func (q *Qubit) Apply(m matrix.Matrix) *Qubit {
	q.v = q.v.Apply(m)
	return q
}

func (q *Qubit) Normalize() *Qubit {
	var sum float64
	for _, amp := range q.v {
		sum = sum + math.Pow(cmplx.Abs(amp), 2)
	}
	z := 1 / math.Sqrt(sum)
	q.v = q.v.Mul(complex(z, 0))
	return q
}

func (q *Qubit) Amplitude() []complex128 {
	a := []complex128{}
	for _, amp := range q.v {
		a = append(a, amp)
	}
	return a
}

func (q *Qubit) Probability() []float64 {
	list := []float64{}
	for _, amp := range q.v {
		p := math.Pow(cmplx.Abs(amp), 2)
		list = append(list, p)
	}
	return list
}

func (q *Qubit) Measure(bit ...int) *Qubit {
	if len(bit) > 0 {
		return q.MeasureAt(bit[0])
	}

	rand.Seed(time.Now().UnixNano())
	r := rand.Float64()

	plist := q.Probability()
	var sum float64
	for i, p := range plist {
		if sum <= r && r < sum+p {
			q.v = v.NewZero(len(q.v))
			q.v[i] = 1
			break
		}
		sum = sum + p
	}

	return q
}

func (q *Qubit) ProbabilityZeroAt(bit int) ([]int, []float64) {
	p := []float64{}
	index := []int{}

	dim := q.v.Dimension()
	den := int(math.Pow(2, float64(bit+1)))
	div := dim / den

	for i := 0; i < dim; i++ {
		p = append(p, q.Probability()[i])
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
	p := []float64{}
	index := []int{}

	zi, _ := q.ProbabilityZeroAt(bit)
	one := []int{}
	for i, _ := range q.v {
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

func (q *Qubit) MeasureAt(bit int) *Qubit {
	index, p := q.ProbabilityZeroAt(bit)

	rand.Seed(time.Now().UnixNano())
	r := rand.Float64()

	var sum float64
	for _, pp := range p {
		sum = sum + pp
	}

	if r > sum {
		for _, i := range index {
			q.v[i] = complex(0, 0)
		}

		q.Normalize()
		return One()
	}

	one := []int{}
	for i, _ := range q.v {
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
		q.v[i] = complex(0, 0)
	}

	q.Normalize()
	return Zero()
}

func TensorProduct(q ...*Qubit) *Qubit {
	q1 := q[0]
	for i := 1; i < len(q); i++ {
		q1 = q1.TensorProduct(q[i])
	}
	return q1
}
