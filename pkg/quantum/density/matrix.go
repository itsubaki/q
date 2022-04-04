package density

import (
	"fmt"
	"math"

	"github.com/itsubaki/q/pkg/math/matrix"
	"github.com/itsubaki/q/pkg/math/number"
	"github.com/itsubaki/q/pkg/quantum/gate"
	"github.com/itsubaki/q/pkg/quantum/qubit"
)

type Matrix struct {
	m matrix.Matrix
}

func New(v ...[]complex128) *Matrix {
	return &Matrix{matrix.New(v...)}
}

func (m *Matrix) Add(p float64, q *qubit.Qubit) *Matrix {
	if p < 0 || p > 1 {
		panic(fmt.Sprintf("p must be 0 <= p =< 1. p=%v", p))
	}

	n := q.Dimension()
	if len(m.m) < 1 {
		m.m = matrix.Zero(n)
	}

	if len(m.m) != n {
		panic(fmt.Sprintf("invalid dimension. m=%d n=%d", len(m.m), n))
	}

	op := q.OuterProduct(q).Mul(complex(p, 0))
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			m.m[i][j] = m.m[i][j] + op[i][j]
		}
	}

	return m
}

func (m *Matrix) Apply(u matrix.Matrix) *Matrix {
	m.m = u.Dagger().Apply(m.m).Apply(u)
	return m
}

func (m *Matrix) Measure(q *qubit.Qubit) complex128 {
	return m.m.Apply(q.OuterProduct(q)).Trace()
}

func (m *Matrix) ExpectedValue(u matrix.Matrix) complex128 {
	return m.m.Apply(u).Trace()
}

func (m *Matrix) Trace() complex128 {
	return m.m.Trace()
}

func (m *Matrix) PartialTrace(i int) (*Matrix, error) {
	n := number.Pow(2, m.NumberOfBit()-1)
	out := matrix.Zero(n)

	// TODO: Not Implemented
	return &Matrix{out}, fmt.Errorf("Not Implemented")
}

func (m *Matrix) Squared() *Matrix {
	c := m.m.Clone()
	return &Matrix{c.Apply(c)}
}

func (m *Matrix) NumberOfBit() int {
	mm, _ := m.m.Dimension()
	log := math.Log2(float64(mm))
	return int(log)
}

func (m *Matrix) Depolarizing(p float64) *Matrix {
	if p < 0 || p > 1 {
		panic(fmt.Sprintf("p must be 0 <= p =< 1. p=%v", p))
	}

	n := m.NumberOfBit()
	i := gate.I(n).Mul(complex(p/2, 0))
	r := m.m.Mul(complex(1-p, 0))

	return &Matrix{i.Add(r)}
}

func Flip(p float64, m matrix.Matrix) (matrix.Matrix, matrix.Matrix) {
	if p < 0 || p > 1 {
		panic(fmt.Sprintf("p must be 0 <= p =< 1. p=%v", p))
	}

	d, _ := m.Dimension()
	n := int(math.Log2(float64(d)))

	e0 := gate.I(n).Mul(complex(math.Sqrt(p), 0))
	e1 := m.Mul(complex(math.Sqrt(1-p), 0))
	return e0, e1
}

func BitFlip(p float64) (matrix.Matrix, matrix.Matrix) {
	return Flip(p, gate.X())
}

func PhaseFlip(p float64) (matrix.Matrix, matrix.Matrix) {
	return Flip(p, gate.Z())
}

func BitPhaseFlip(p float64) (matrix.Matrix, matrix.Matrix) {
	return Flip(p, gate.Y())
}
