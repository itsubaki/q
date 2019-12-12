package matrix

import (
	"fmt"
	"math/cmplx"
)

type Matrix [][]complex128

func New(v ...[]complex128) Matrix {
	m := make(Matrix, len(v))
	for i := 0; i < len(v); i++ {
		m[i] = v[i]
	}

	return m
}

func Zero(n int) Matrix {
	m := make(Matrix, n)
	for i := 0; i < n; i++ {
		m[i] = make([]complex128, 0)
		for j := 0; j < n; j++ {
			m[i] = append(m[i], complex(0, 0))
		}
	}

	return m
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

	e := Eps(eps...)
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

	t := Matrix{}
	for i := 0; i < p; i++ {
		v := make([]complex128, 0)
		for j := 0; j < q; j++ {
			v = append(v, m[j][i])
		}

		t = append(t, v)
	}

	return t
}

func (m Matrix) Conjugate() Matrix {
	p, q := m.Dimension()

	c := Matrix{}
	for i := 0; i < p; i++ {
		v := make([]complex128, 0)
		for j := 0; j < q; j++ {
			v = append(v, cmplx.Conj(m[i][j]))
		}

		c = append(c, v)
	}

	return c
}

func (m Matrix) Dagger() Matrix {
	return m.Transpose().Conjugate()
}

func (m Matrix) IsHermite(eps ...float64) bool {
	p, q := m.Dimension()
	d := m.Dagger()
	e := Eps(eps...)

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
	e := Eps(eps...)

	for i := 0; i < p; i++ {
		for j := 0; j < q; j++ {
			if i == j {
				if cmplx.Abs(d[i][j]-complex(1, 0)) > e {
					return false
				}
				continue
			}

			if cmplx.Abs(d[i][j]-complex(0, 0)) > e {
				return false
			}
		}
	}

	return true
}

func (m Matrix) Apply(n Matrix) Matrix {
	p, _ := m.Dimension()
	a, b := n.Dimension()

	apply := Matrix{}
	for i := 0; i < a; i++ {
		v := make([]complex128, 0)
		for j := 0; j < b; j++ {
			c := complex(0, 0)
			for k := 0; k < p; k++ {
				c = c + n[i][k]*m[k][j]
			}
			v = append(v, c)
		}

		apply = append(apply, v)
	}

	return apply
}

func (m Matrix) Mul(z complex128) Matrix {
	p, q := m.Dimension()

	mul := Matrix{}
	for i := 0; i < p; i++ {
		v := make([]complex128, 0)
		for j := 0; j < q; j++ {
			v = append(v, z*m[i][j])
		}

		mul = append(mul, v)
	}

	return mul
}

func (m Matrix) Add(n Matrix) Matrix {
	p, q := m.Dimension()

	add := Matrix{}
	for i := 0; i < p; i++ {
		v := make([]complex128, 0)
		for j := 0; j < q; j++ {
			v = append(v, m[i][j]+n[i][j])
		}

		add = append(add, v)
	}

	return add
}

func (m Matrix) Sub(n Matrix) Matrix {
	p, q := m.Dimension()

	sub := Matrix{}
	for i := 0; i < p; i++ {
		v := make([]complex128, 0)
		for j := 0; j < q; j++ {
			v = append(v, m[i][j]-n[i][j])
		}

		sub = append(sub, v)
	}

	return sub
}

func (m Matrix) Trace() complex128 {
	p, _ := m.Dimension()

	var sum complex128
	for i := 0; i < p; i++ {
		sum = sum + m[i][i]
	}

	return sum
}

func (m Matrix) Clone() Matrix {
	p, q := m.Dimension()

	c := Matrix{}
	for i := 0; i < p; i++ {
		v := make([]complex128, 0)
		for j := 0; j < q; j++ {
			v = append(v, m[i][j])
		}

		c = append(c, v)
	}

	return c
}

func (m Matrix) Inverse() Matrix {
	clone := m.Clone()
	p, q := clone.Dimension()
	if p != q {
		panic(fmt.Sprintf("dimension invalid. p=%d q=%d", p, q))
	}

	inv := Matrix{}
	for i := 0; i < p; i++ {
		v := make([]complex128, 0)
		for j := 0; j < q; j++ {
			if i == j {
				v = append(v, complex(1, 0))
				continue
			}
			v = append(v, complex(0, 0))
		}
		inv = append(inv, v)
	}

	for i := 0; i < p; i++ {
		c := 1 / clone[i][i]
		for j := 0; j < q; j++ {
			clone[i][j] = c * clone[i][j]
			inv[i][j] = c * inv[i][j]
		}
		for j := 0; j < q; j++ {
			if i == j {
				continue
			}

			c := clone[j][i]
			for k := 0; k < q; k++ {
				clone[j][k] = clone[j][k] - c*clone[i][k]
				inv[j][k] = inv[j][k] - c*inv[i][k]
			}
		}
	}

	return inv
}

func (m Matrix) TensorProduct(n Matrix) Matrix {
	p, q := m.Dimension()
	a, b := n.Dimension()

	t := Matrix{}
	for i := 0; i < p; i++ {
		for k := 0; k < a; k++ {
			r := make([]complex128, 0)
			for j := 0; j < q; j++ {
				for l := 0; l < b; l++ {
					r = append(r, m[i][j]*n[k][l])
				}
			}

			t = append(t, r)
		}
	}

	return t
}

func TensorProductN(m Matrix, bit ...int) Matrix {
	if len(bit) < 1 {
		return m
	}

	p := m
	for i := 1; i < bit[0]; i++ {
		p = p.TensorProduct(m)
	}

	return p
}

func TensorProduct(m ...Matrix) Matrix {
	p := m[0]
	for i := 1; i < len(m); i++ {
		p = p.TensorProduct(m[i])
	}

	return p
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

func Eps(eps ...float64) float64 {
	if len(eps) > 0 {
		return eps[0]
	}

	return 0.0
}
