package matrix

import (
	"fmt"
	"math/cmplx"
)

type Matrix [][]complex128

func New(v ...[]complex128) Matrix {
	out := make(Matrix, len(v))
	for i := 0; i < len(v); i++ {
		out[i] = v[i]
	}

	return out
}

func Zero(n int) Matrix {
	out := make(Matrix, n)
	for i := 0; i < n; i++ {
		out[i] = make([]complex128, 0)
		for j := 0; j < n; j++ {
			out[i] = append(out[i], 0)
		}
	}

	return out
}

func (m Matrix) Equals(n Matrix, eps ...float64) bool {
	p, q := m.Dimension()
	a, b := n.Dimension()

	if a != p {
		return false
	}

	if b != q {
		return false
	}

	e := epsilon(eps...)
	for i := 0; i < p; i++ {
		for j := 0; j < q; j++ {
			if cmplx.Abs(m[i][j]-n[i][j]) > e {
				return false
			}
		}
	}

	return true
}

func (m Matrix) Dimension() (int, int) {
	return len(m), len(m[0])
}

func (m Matrix) Transpose() Matrix {
	p, q := m.Dimension()

	out := Matrix{}
	for i := 0; i < p; i++ {
		v := make([]complex128, 0)
		for j := 0; j < q; j++ {
			v = append(v, m[j][i])
		}

		out = append(out, v)
	}

	return out
}

func (m Matrix) Conjugate() Matrix {
	p, q := m.Dimension()

	out := Matrix{}
	for i := 0; i < p; i++ {
		v := make([]complex128, 0)
		for j := 0; j < q; j++ {
			v = append(v, cmplx.Conj(m[i][j]))
		}

		out = append(out, v)
	}

	return out
}

func (m Matrix) Dagger() Matrix {
	return m.Transpose().Conjugate()
}

func (m Matrix) IsHermite(eps ...float64) bool {
	p, q := m.Dimension()
	d := m.Dagger()

	e := epsilon(eps...)
	for i := 0; i < p; i++ {
		for j := 0; j < q; j++ {
			if cmplx.Abs(m[i][j]-d[i][j]) > e {
				return false
			}
		}
	}

	return true
}

func (m Matrix) IsUnitary(eps ...float64) bool {
	p, q := m.Dimension()
	d := m.Apply(m.Dagger())

	e := epsilon(eps...)
	for i := 0; i < p; i++ {
		for j := 0; j < q; j++ {
			if i == j {
				if cmplx.Abs(d[i][j]-1) > e {
					return false
				}
				continue
			}

			if cmplx.Abs(d[i][j]) > e {
				return false
			}
		}
	}

	return true
}

func (m Matrix) Apply(n Matrix) Matrix {
	p, _ := m.Dimension()
	a, b := n.Dimension()

	out := Matrix{}
	for i := 0; i < a; i++ {
		v := make([]complex128, 0)
		for j := 0; j < b; j++ {
			c := complex(0, 0)
			for k := 0; k < p; k++ {
				c = c + n[i][k]*m[k][j]
			}
			v = append(v, c)
		}

		out = append(out, v)
	}

	return out
}

func (m Matrix) Mul(z complex128) Matrix {
	p, q := m.Dimension()

	out := Matrix{}
	for i := 0; i < p; i++ {
		v := make([]complex128, 0)
		for j := 0; j < q; j++ {
			v = append(v, z*m[i][j])
		}

		out = append(out, v)
	}

	return out
}

func (m Matrix) Add(n Matrix) Matrix {
	p, q := m.Dimension()

	out := Matrix{}
	for i := 0; i < p; i++ {
		v := make([]complex128, 0)
		for j := 0; j < q; j++ {
			v = append(v, m[i][j]+n[i][j])
		}

		out = append(out, v)
	}

	return out
}

func (m Matrix) Sub(n Matrix) Matrix {
	p, q := m.Dimension()

	out := Matrix{}
	for i := 0; i < p; i++ {
		v := make([]complex128, 0)
		for j := 0; j < q; j++ {
			v = append(v, m[i][j]-n[i][j])
		}

		out = append(out, v)
	}

	return out
}

func (m Matrix) Trace() complex128 {
	p, _ := m.Dimension()

	var sum complex128
	for i := 0; i < p; i++ {
		sum = sum + m[i][i]
	}

	return sum
}

func (m Matrix) Real() [][]float64 {
	out := make([][]float64, 0)
	for i, r := range m {
		out = append(out, make([]float64, 0))
		for j := range r {
			out[i] = append(out[i], real(m[i][j]))
		}
	}

	return out
}

func (m Matrix) Imag() [][]float64 {
	out := make([][]float64, 0)
	for i, r := range m {
		out = append(out, make([]float64, 0))
		for j := range r {
			out[i] = append(out[i], imag(m[i][j]))
		}
	}

	return out
}

func (m Matrix) Clone() Matrix {
	p, q := m.Dimension()

	out := Matrix{}
	for i := 0; i < p; i++ {
		v := make([]complex128, 0)
		for j := 0; j < q; j++ {
			v = append(v, m[i][j])
		}

		out = append(out, v)
	}

	return out
}

func (m Matrix) Inverse() Matrix {
	clone := m.Clone()
	p, q := clone.Dimension()
	if p != q {
		panic(fmt.Sprintf("invalid dimension. p=%d q=%d", p, q))
	}

	out := Matrix{}
	for i := 0; i < p; i++ {
		v := make([]complex128, 0)
		for j := 0; j < q; j++ {
			if i == j {
				v = append(v, 1)
				continue
			}
			v = append(v, 0)
		}
		out = append(out, v)
	}

	for i := 0; i < p; i++ {
		c := 1 / clone[i][i]
		for j := 0; j < q; j++ {
			clone[i][j] = c * clone[i][j]
			out[i][j] = c * out[i][j]
		}
		for j := 0; j < q; j++ {
			if i == j {
				continue
			}

			c := clone[j][i]
			for k := 0; k < q; k++ {
				clone[j][k] = clone[j][k] - c*clone[i][k]
				out[j][k] = out[j][k] - c*out[i][k]
			}
		}
	}

	return out
}

func (m Matrix) TensorProduct(n Matrix) Matrix {
	p, q := m.Dimension()
	a, b := n.Dimension()

	out := Matrix{}
	for i := 0; i < p; i++ {
		for k := 0; k < a; k++ {
			r := make([]complex128, 0)
			for j := 0; j < q; j++ {
				for l := 0; l < b; l++ {
					r = append(r, m[i][j]*n[k][l])
				}
			}

			out = append(out, r)
		}
	}

	return out
}

func Apply(m ...Matrix) Matrix {
	out := m[0]
	for i := 1; i < len(m); i++ {
		out = out.Apply(m[i])
	}

	return out
}

func TensorProductN(m Matrix, n ...int) Matrix {
	if len(n) < 1 {
		return m
	}

	list := make([]Matrix, 0)
	for i := 0; i < n[0]; i++ {
		list = append(list, m)
	}

	return TensorProduct(list...)
}

func TensorProduct(m ...Matrix) Matrix {
	out := m[0]
	for i := 1; i < len(m); i++ {
		out = out.TensorProduct(m[i])
	}

	return out
}

func Commutator(m, n Matrix) Matrix {
	mn := n.Apply(m)
	nm := m.Apply(n)
	return mn.Sub(nm)
}

func AntiCommutator(m, n Matrix) Matrix {
	mn := n.Apply(m)
	nm := m.Apply(n)
	return mn.Add(nm)
}

func epsilon(eps ...float64) float64 {
	if len(eps) > 0 {
		return eps[0]
	}

	return 1e-13
}
