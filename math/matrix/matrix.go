package matrix

import (
	"iter"
	"math/cmplx"

	"github.com/itsubaki/q/math/epsilon"
)

// Matrix is a matrix of complex128.
type Matrix struct {
	Rows int
	Cols int
	Data []complex128
}

// New returns a new matrix of complex128.
func New(v ...[]complex128) Matrix {
	rows := len(v)
	var cols int
	if rows > 0 {
		cols = len(v[0])
	}

	data := make([]complex128, 0, rows*cols)
	for i := range rows {
		data = append(data, v[i]...)
	}

	return Matrix{
		Rows: rows,
		Cols: cols,
		Data: data,
	}
}

// Zero returns a zero matrix.
func Zero(rows, cols int) Matrix {
	return Matrix{
		Rows: rows,
		Cols: cols,
		Data: make([]complex128, rows*cols),
	}
}

// Identity returns an identity matrix.
func Identity(rows, cols int) Matrix {
	m := Zero(rows, cols)
	for i := range rows {
		m.Set(i, i, 1)
	}

	return m
}

// At returns a value of matrix at (i,j).
func (m Matrix) At(i, j int) complex128 {
	return m.Data[i*m.Cols+j]
}

// Row returns a row of matrix at (i).
func (m Matrix) Row(i int) []complex128 {
	return row(m.Data, m.Cols, i)
}

// Set sets a value of matrix at (i,j).
func (m Matrix) Set(i, j int, v complex128) {
	m.Data[i*m.Cols+j] = v
}

// Seq2 returns a sequence of rows.
func (m Matrix) Seq2() iter.Seq2[int, []complex128] {
	return func(yield func(int, []complex128) bool) {
		for i := range m.Rows {
			if !yield(i, m.Row(i)) {
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
	for i := range m.Data {
		if cmplx.Abs(m.Data[i]-n.Data[i]) > e {
			return false
		}
	}

	return true
}

// Dimension returns a dimension of matrix.
func (m Matrix) Dimension() (rows int, cols int) {
	return m.Rows, m.Cols
}

// Transpose returns a transpose matrix.
func (m Matrix) Transpose() Matrix {
	p, q := m.Dimension()

	out := Zero(p, q)
	for i := range p {
		for j := range q {
			out.Set(i, j, m.At(j, i))
		}
	}

	return out
}

// Conjugate returns a conjugate matrix.
func (m Matrix) Conjugate() Matrix {
	p, q := m.Dimension()

	out := Zero(p, q)
	for i := range out.Data {
		out.Data[i] = cmplx.Conj(m.Data[i])
	}

	return out
}

// Dagger returns conjugate transpose matrix.
func (m Matrix) Dagger() Matrix {
	p, q := m.Dimension()

	out := Zero(p, q)
	for i := range p {
		for j := range q {
			out.Set(i, j, cmplx.Conj(m.At(j, i)))
		}
	}

	return out
}

// IsSquare returns true if m is square matrix.
func (m Matrix) IsSquare() bool {
	return m.Rows == m.Cols
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

	mmd := m.Apply(m.Dagger())
	id := Identity(m.Dimension())
	return mmd.Equals(id, epsilon.E13(eps...))
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
				c = c + m.At(k, j)*n.At(i, k)
			}

			out.Set(i, j, c)
		}
	}

	return out
}

// Mul returns a matrix of z*m.
func (m Matrix) Mul(z complex128) Matrix {
	out := Zero(m.Dimension())
	for i := range m.Data {
		out.Data[i] = m.Data[i] * z
	}

	return out
}

// Add returns a matrix of m+n.
func (m Matrix) Add(n Matrix) Matrix {
	out := Zero(m.Dimension())
	for i := range m.Data {
		out.Data[i] = m.Data[i] + n.Data[i]
	}

	return out
}

// Sub returns a matrix of m-n.
func (m Matrix) Sub(n Matrix) Matrix {
	out := Zero(m.Dimension())
	for i := range m.Data {
		out.Data[i] = m.Data[i] - n.Data[i]
	}

	return out
}

// Trace returns a trace of matrix.
func (m Matrix) Trace() complex128 {
	var sum complex128
	for i := range m.Rows {
		sum = sum + m.At(i, i)
	}

	return sum
}

// Real returns a real part of matrix.
func (m Matrix) Real() [][]float64 {
	data := make([]float64, len(m.Data))
	for i := range m.Data {
		data[i] = real(m.Data[i])
	}

	out := make([][]float64, m.Rows)
	for i := range m.Rows {
		out[i] = row(data, m.Cols, i)
	}

	return out
}

// Imag returns an imaginary part of matrix.
func (m Matrix) Imag() [][]float64 {
	data := make([]float64, len(m.Data))
	for i := range m.Data {
		data[i] = imag(m.Data[i])
	}

	out := make([][]float64, m.Rows)
	for i := range m.Rows {
		out[i] = row(data, m.Cols, i)
	}

	return out
}

// Clone returns a clone of matrix.
func (m Matrix) Clone() Matrix {
	out := Zero(m.Dimension())
	copy(out.Data, m.Data)
	return out
}

// Inverse returns an inverse matrix of m.
func (m Matrix) Inverse() Matrix {
	p, q := m.Dimension()
	mm := m.Clone()

	out := Identity(p, q)
	for i := range p {
		c := 1 / mm.At(i, i)
		for j := range q {
			mm.Set(i, j, c*mm.At(i, j))
			out.Set(i, j, c*out.At(i, j))
		}

		for j := range q {
			if i == j {
				continue
			}

			c := mm.At(j, i)
			for k := range q {
				mm.Set(j, k, mm.At(j, k)-c*mm.At(i, k))
				out.Set(j, k, out.At(j, k)-c*out.At(i, k))
			}
		}
	}

	return out
}

// TensorProduct returns a tensor product of m and n.
func (m Matrix) TensorProduct(n Matrix) Matrix {
	p, q := m.Dimension()
	a, b := n.Dimension()

	data := make([]complex128, 0, p*a*q*b)
	for i := range p {
		for k := range a {
			for j := range q {
				for l := range b {
					data = append(data, m.At(i, j)*n.At(k, l))
				}
			}
		}
	}

	return Matrix{
		Rows: p * a,
		Cols: q * b,
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
	mn, nm := n.Apply(m), m.Apply(n)
	return mn.Sub(nm)
}

// AntiCommutator returns a matrix of {m,n}.
func AntiCommutator(m, n Matrix) Matrix {
	mn, nm := n.Apply(m), m.Apply(n)
	return mn.Add(nm)
}

func row[T any](arr []T, cols, i int) []T {
	return arr[i*cols : (i+1)*cols]
}
