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

func (m0 Matrix) Equals(m1 Matrix, eps ...float64) bool {
	m, n := m0.Dimension()
	p, q := m1.Dimension()

	if m != p {
		return false
	}

	if n != q {
		return false
	}

	e := Eps(eps...)
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if cmplx.Abs(m0[i][j]-m1[i][j]) > e {
				return false
			}
		}
	}

	return true
}

func (m0 Matrix) Dimension() (int, int) {
	return len(m0), len(m0[0])
}

func (m0 Matrix) Transpose() Matrix {
	p, q := m0.Dimension()

	m2 := Matrix{}
	for i := 0; i < p; i++ {
		v := make([]complex128, 0)
		for j := 0; j < q; j++ {
			v = append(v, m0[j][i])
		}
		m2 = append(m2, v)
	}

	return m2
}

func (m0 Matrix) Conjugate() Matrix {
	p, q := m0.Dimension()

	m2 := Matrix{}
	for i := 0; i < p; i++ {
		v := make([]complex128, 0)
		for j := 0; j < q; j++ {
			v = append(v, cmplx.Conj(m0[i][j]))
		}
		m2 = append(m2, v)
	}

	return m2
}

func (m0 Matrix) Dagger() Matrix {
	return m0.Transpose().Conjugate()
}

func (m0 Matrix) IsHermite(eps ...float64) bool {
	p, q := m0.Dimension()
	m := m0.Dagger()
	e := Eps(eps...)

	for i := 0; i < p; i++ {
		for j := 0; j < q; j++ {
			if cmplx.Abs(m0[i][j]-m[i][j]) > e {
				return false
			}
		}
	}

	return true
}

func (m0 Matrix) IsUnitary(eps ...float64) bool {
	p, q := m0.Dimension()
	m := m0.Apply(m0.Dagger())
	e := Eps(eps...)

	for i := 0; i < p; i++ {
		for j := 0; j < q; j++ {
			if i == j {
				if cmplx.Abs(m[i][j]-complex(1, 0)) > e {
					return false
				}
				continue
			}

			if cmplx.Abs(m[i][j]-complex(0, 0)) > e {
				return false
			}
		}
	}

	return true
}

func (m0 Matrix) Apply(m1 Matrix) Matrix {
	m, n := m1.Dimension()
	p, _ := m0.Dimension()

	m2 := Matrix{}
	for i := 0; i < m; i++ {
		v := make([]complex128, 0)
		for j := 0; j < n; j++ {
			c := complex(0, 0)
			for k := 0; k < p; k++ {
				c = c + m1[i][k]*m0[k][j]
			}
			v = append(v, c)
		}
		m2 = append(m2, v)
	}

	return m2
}

func (m0 Matrix) Mul(z complex128) Matrix {
	p, q := m0.Dimension()

	m := Matrix{}
	for i := 0; i < p; i++ {
		v := make([]complex128, 0)
		for j := 0; j < q; j++ {
			v = append(v, z*m0[i][j])
		}
		m = append(m, v)
	}

	return m
}

func (m0 Matrix) Add(m1 Matrix) Matrix {
	p, q := m0.Dimension()

	m := Matrix{}
	for i := 0; i < p; i++ {
		v := make([]complex128, 0)
		for j := 0; j < q; j++ {
			v = append(v, m0[i][j]+m1[i][j])
		}
		m = append(m, v)
	}

	return m
}

func (m0 Matrix) Sub(m1 Matrix) Matrix {
	p, q := m0.Dimension()

	m := Matrix{}
	for i := 0; i < p; i++ {
		v := make([]complex128, 0)
		for j := 0; j < q; j++ {
			v = append(v, m0[i][j]-m1[i][j])
		}
		m = append(m, v)
	}

	return m
}

func (m0 Matrix) Trace() complex128 {
	p, _ := m0.Dimension()
	var sum complex128
	for i := 0; i < p; i++ {
		sum = sum + m0[i][i]
	}

	return sum
}

func (m0 Matrix) Clone() Matrix {
	m, n := m0.Dimension()
	ret := Matrix{}
	for i := 0; i < m; i++ {
		v := make([]complex128, 0)
		for j := 0; j < n; j++ {
			v = append(v, m0[i][j])
		}
		ret = append(ret, v)
	}

	return ret
}

func (m0 Matrix) Inverse() Matrix {
	mat := m0.Clone()
	m, n := mat.Dimension()
	if m != n {
		panic(fmt.Sprintf("m=%d n=%d", m, n))
	}

	inv := Matrix{}
	for i := 0; i < m; i++ {
		v := make([]complex128, 0)
		for j := 0; j < n; j++ {
			if i == j {
				v = append(v, complex(1, 0))
				continue
			}
			v = append(v, complex(0, 0))
		}
		inv = append(inv, v)
	}

	for i := 0; i < m; i++ {
		c := 1 / mat[i][i]
		for j := 0; j < n; j++ {
			mat[i][j] = c * mat[i][j]
			inv[i][j] = c * inv[i][j]
		}
		for j := 0; j < n; j++ {
			if i == j {
				continue
			}
			c := mat[j][i]
			for k := 0; k < n; k++ {
				mat[j][k] = mat[j][k] - c*mat[i][k]
				inv[j][k] = inv[j][k] - c*inv[i][k]
			}
		}
	}

	return inv
}

func (m0 Matrix) TensorProduct(m1 Matrix) Matrix {
	m, n := m0.Dimension()
	p, q := m1.Dimension()

	tmp := make([]Matrix, 0)
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			tmp = append(tmp, m1.Mul(m0[i][j]))
		}
	}

	m2 := Matrix{}
	for l := 0; l < len(tmp); l = l + m {
		for j := 0; j < p; j++ {
			v := make([]complex128, 0)
			for i := l; i < l+m; i++ {
				for k := 0; k < q; k++ {
					v = append(v, tmp[i][j][k])
				}
			}
			m2 = append(m2, v)
		}
	}

	return m2
}

func TensorProductN(m Matrix, bit ...int) Matrix {
	if len(bit) < 1 {
		return m
	}

	m0 := m
	for i := 1; i < bit[0]; i++ {
		m0 = m0.TensorProduct(m)
	}

	return m0
}

func TensorProduct(m ...Matrix) Matrix {
	m0 := m[0]
	for i := 1; i < len(m); i++ {
		m0 = m0.TensorProduct(m[i])
	}

	return m0
}

func Commutator(m0, m1 Matrix) Matrix {
	m10 := m1.Apply(m0)
	m01 := m0.Apply(m1)

	return m10.Sub(m01)
}

func AntiCommutator(m0, m1 Matrix) Matrix {
	m10 := m1.Apply(m0)
	m01 := m0.Apply(m1)

	return m10.Add(m01)
}

func Eps(eps ...float64) float64 {
	if len(eps) > 0 {
		return eps[0]
	}

	return 0.0
}
