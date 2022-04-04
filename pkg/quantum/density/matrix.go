package density

import (
	"fmt"
	"math"
	"strconv"

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

func (m *Matrix) Raw() matrix.Matrix {
	return m.m
}

func (m *Matrix) Dimension() int {
	return len(m.m)
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

func (m *Matrix) PartialTrace(index int) *Matrix {
	n, d := m.NumberOfBit(), m.Dimension()
	f := fmt.Sprintf("%s%s%s", "%0", strconv.Itoa(n), "s")

	out := matrix.Zero(number.Pow(2, n-1))
	for i := 0; i < d; i++ {
		ibits := fmt.Sprintf(f, strconv.FormatInt(int64(i), 2))

		for j := 0; j < d; j++ {
			jbits := fmt.Sprintf(f, strconv.FormatInt(int64(j), 2))
			if ibits[index] != jbits[index] {
				continue
			}

			v0 := number.Must(strconv.ParseInt(string(ibits[index:index+1]), 2, 0))
			v1 := number.Must(strconv.ParseInt(string(jbits[index:index+1]), 2, 0))

			out[v0][v1] = out[v0][v1] + m.m[i][j]
		}
	}

	return &Matrix{m: out}
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
