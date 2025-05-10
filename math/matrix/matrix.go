package matrix

import (
	"iter"
	"math/cmplx"

	"github.com/itsubaki/q/math/epsilon"
)

// Matrix is a matrix of complex128.
type Matrix struct {
	Data [][]complex128
}

// New returns a new matrix of complex128.
func New(v ...[]complex128) Matrix {
	out := make([][]complex128, len(v))
	copy(out, v)

	return Matrix{
		Data: out,
	}
}

// Zero returns a zero matrix.
func Zero(n, m int) Matrix {
	out := make([][]complex128, n)
	for i := range n {
		out[i] = make([]complex128, m)
	}

	return Matrix{
		Data: out,
	}
}

// Identity returns an identity matrix.
func Identity(n, m int) Matrix {
	out := Zero(n, m)
	for i := range n {
		out.Data[i][i] = 1
	}

	return out
}

func (m Matrix) Seq2() iter.Seq2[int, []complex128] {
	return func(yield func(int, []complex128) bool) {
		for i := range m.Data {
			if !yield(i, m.Data[i]) {
				return
			}
		}
	}
}

// Equals returns true if m equals n.
// If eps is not given, epsilon.E13 is used.
func (m Matrix) Equals(n Matrix, eps ...float64) bool {
	p, q := m.Dimension()
	a, b := n.Dimension()

	if a != p {
		return false
	}

	if b != q {
		return false
	}

	e := epsilon.E13(eps...)
	for i := range p {
		for j := range q {
			if cmplx.Abs(m.Data[i][j]-n.Data[i][j]) > e {
				return false
			}
		}
	}

	return true
}

// Dimension returns a dimension of matrix.
func (m Matrix) Dimension() (int, int) {
	return len(m.Data), len(m.Data[0])
}

// Transpose returns a transpose matrix.
func (m Matrix) Transpose() Matrix {
	p, q := m.Dimension()

	out := Zero(p, q)
	for i := range q {
		for j := range p {
			out.Data[i][j] = m.Data[j][i]
		}
	}

	return out
}

// Conjugate returns a conjugate matrix.
func (m Matrix) Conjugate() Matrix {
	p, q := m.Dimension()

	out := Zero(p, q)
	for i := range p {
		for j := range q {
			out.Data[i][j] = cmplx.Conj(m.Data[i][j])
		}
	}

	return out
}

// Dagger returns conjugate transpose matrix.
func (m Matrix) Dagger() Matrix {
	p, q := m.Dimension()

	out := Zero(p, q)
	for i := range p {
		for j := range q {
			out.Data[i][j] = cmplx.Conj(m.Data[j][i])
		}
	}

	return out
}

// IsSquare returns true if m is square matrix.
func (m Matrix) IsSquare() bool {
	p, q := m.Dimension()
	return p == q
}

// IsHermitian returns true if m is hermitian matrix.
func (m Matrix) IsHermite(eps ...float64) bool {
	if !m.IsSquare() {
		return false
	}

	return m.Equals(m.Dagger(), epsilon.E13(eps...))
}

// IsUnitary returns true if m is unitary matrix.
func (m Matrix) IsUnitary(eps ...float64) bool {
	if !m.IsSquare() {
		return false
	}

	p, q := m.Dimension()
	return m.Apply(m.Dagger()).Equals(Identity(p, q), epsilon.E13(eps...))
}

// Apply returns a matrix product of m and n.
// A.Apply(B) is BA.
func (m Matrix) Apply(n Matrix) Matrix {
	a, b := m.Dimension()
	_, p := n.Dimension()

	out := Zero(a, p)
	for i := range a {
		for j := range p {
			var c complex128
			for k := range b {
				c = c + n.Data[i][k]*m.Data[k][j]
			}

			out.Data[i][j] = c
		}
	}

	return out
}

// Mul returns a matrix of z*m.
func (m Matrix) Mul(z complex128) Matrix {
	p, q := m.Dimension()

	out := Zero(p, q)
	for i := range p {
		for j := range q {
			out.Data[i][j] = z * m.Data[i][j]
		}
	}

	return out
}

// Add returns a matrix of m+n.
func (m Matrix) Add(n Matrix) Matrix {
	p, q := m.Dimension()

	out := Zero(p, q)
	for i := range p {
		for j := range q {
			out.Data[i][j] = m.Data[i][j] + n.Data[i][j]
		}
	}

	return out
}

// Sub returns a matrix of m-n.
func (m Matrix) Sub(n Matrix) Matrix {
	p, q := m.Dimension()

	out := Zero(p, q)
	for i := range p {
		for j := range q {
			out.Data[i][j] = m.Data[i][j] - n.Data[i][j]
		}
	}

	return out
}

// Trace returns a trace of matrix.
func (m Matrix) Trace() complex128 {
	p, _ := m.Dimension()

	var sum complex128
	for i := range p {
		sum = sum + m.Data[i][i]
	}

	return sum
}

// Real returns a real part of matrix.
func (m Matrix) Real() [][]float64 {
	out := make([][]float64, len(m.Data))
	for i, r := range m.Data {
		out[i] = make([]float64, len(m.Data[i]))
		for j := range r {
			out[i][j] = real(m.Data[i][j])
		}
	}

	return out
}

// Imag returns an imaginary part of matrix.
func (m Matrix) Imag() [][]float64 {
	out := make([][]float64, len(m.Data))
	for i, r := range m.Data {
		out[i] = make([]float64, len(m.Data[i]))
		for j := range r {
			out[i][j] = imag(m.Data[i][j])
		}
	}

	return out
}

// Clone returns a clone of matrix.
func (m Matrix) Clone() Matrix {
	p, q := m.Dimension()

	out := Zero(p, q)
	for i := range p {
		for j := range q {
			out.Data[i][j] = m.Data[i][j]
		}
	}

	return out
}

// Inverse returns an inverse matrix of m.
func (m Matrix) Inverse() Matrix {
	p, q := m.Dimension()
	mm := m.Clone()

	out := Identity(p, q)
	for i := range p {
		c := 1 / mm.Data[i][i]
		for j := range q {
			mm.Data[i][j] = c * mm.Data[i][j]
			out.Data[i][j] = c * out.Data[i][j]
		}

		for j := range q {
			if i == j {
				continue
			}

			c := mm.Data[j][i]
			for k := range q {
				mm.Data[j][k] = mm.Data[j][k] - c*mm.Data[i][k]
				out.Data[j][k] = out.Data[j][k] - c*out.Data[i][k]
			}
		}
	}

	return out
}

// TensorProduct returns a tensor product of m and n.
func (m Matrix) TensorProduct(n Matrix) Matrix {
	p, q := m.Dimension()
	a, b := n.Dimension()

	data := make([][]complex128, 0, p*a)
	for i := range p {
		for k := 0; k < a; k++ {
			r := make([]complex128, 0, q*b)
			for j := range q {
				for l := 0; l < b; l++ {
					r = append(r, m.Data[i][j]*n.Data[k][l])
				}
			}

			data = append(data, r)
		}
	}

	return Matrix{
		Data: data,
	}
}

// Apply returns a matrix product of m1, m2, ..., mn.
// Apply(A, B, C, D, ...) is ...DCBA.
func Apply(m ...Matrix) Matrix {
	out := m[0]
	for i := 1; i < len(m); i++ {
		out = out.Apply(m[i])
	}

	return out
}

func ApplyN(m Matrix, n ...int) Matrix {
	if len(n) < 1 {
		return m
	}

	list := make([]Matrix, n[0])
	for i := range n[0] {
		list[i] = m
	}

	return Apply(list...)
}

func TensorProductN(m Matrix, n ...int) Matrix {
	if len(n) < 1 {
		return m
	}

	list := make([]Matrix, n[0])
	for i := range n[0] {
		list[i] = m
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

// Commutator returns a matrix of [m,n].
func Commutator(m, n Matrix) Matrix {
	mn := n.Apply(m)
	nm := m.Apply(n)
	return mn.Sub(nm)
}

// AntiCommutator returns a matrix of {m,n}.
func AntiCommutator(m, n Matrix) Matrix {
	mn := n.Apply(m)
	nm := m.Apply(n)
	return mn.Add(nm)
}
