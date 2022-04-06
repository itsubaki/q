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

func (m *Matrix) Dimension() (int, int) {
	return len(m.m), len(m.m[0])
}

func (m *Matrix) NumberOfBit() int {
	mm, _ := m.m.Dimension()
	log := math.Log2(float64(mm))
	return int(log)
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

func (m *Matrix) Trace() float64 {
	return real(m.m.Trace())
}

func (m *Matrix) SquaredTrace() float64 {
	return real(m.m.Apply(m.m).Trace())
}

func (m *Matrix) PartialTrace(index int) *Matrix {
	n := m.NumberOfBit()
	f := fmt.Sprintf("%s%s%s", "%0", strconv.Itoa(n), "s")
	out := matrix.Zero(number.Pow(2, n-1))

	p, q := m.Dimension()
	for i := 0; i < p; i++ {
		ibin := fmt.Sprintf(f, strconv.FormatInt(int64(i), 2))
		k, kk := fmt.Sprintf("%s%s", ibin[:index], ibin[index+1:]), string(ibin[index])

		for j := 0; j < q; j++ {
			jbin := fmt.Sprintf(f, strconv.FormatInt(int64(j), 2))
			l, ll := fmt.Sprintf("%s%s", jbin[:index], jbin[index+1:]), string(jbin[index])

			if kk != ll {
				continue
			}

			r := number.Must(strconv.ParseInt(k, 2, 0))
			c := number.Must(strconv.ParseInt(l, 2, 0))

			out[r][c] = out[r][c] + m.m[i][j]

			// fmt.Printf("[%v][%v] = [%v][%v] + [%v][%v]\n", r, c, r, c, i, j)
			//
			// 4x4 explicit
			// index -> 0
			// out[0][0] = m.m[0][0] + m.m[2][2]
			// out[0][1] = m.m[0][1] + m.m[2][3]
			// out[1][0] = m.m[1][0] + m.m[3][2]
			// out[1][1] = m.m[1][1] + m.m[3][3]
			//
			// index -> 1
			// out[0][0] = m.m[0][0] + m.m[1][1]
			// out[0][1] = m.m[0][2] + m.m[1][3]
			// out[1][0] = m.m[2][0] + m.m[3][1]
			// out[1][1] = m.m[2][2] + m.m[3][3]
		}
	}

	return &Matrix{m: out}
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
