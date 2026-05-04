package matrix

import (
	"iter"
	"math/cmplx"

	"github.com/itsubaki/q/math/epsilon"
)

// Matrix represents a matrix of complex128 values.
type Matrix struct {
	Rows int
	Cols int
	Data []complex128
}

// New returns a new matrix of complex128 values.
func New(z ...[]complex128) *Matrix {
	rows := len(z)
	var cols int
	if rows > 0 {
		cols = len(z[0])
	}

	data := make([]complex128, rows*cols)
	for i := range rows {
		copy(data[i*cols:(i+1)*cols], z[i])
	}

	return &Matrix{
		Rows: rows,
		Cols: cols,
		Data: data,
	}
}

// Zero returns a zero matrix.
func Zero(rows, cols int) *Matrix {
	return &Matrix{
		Rows: rows,
		Cols: cols,
		Data: make([]complex128, rows*cols),
	}
}

// ZeroLike returns a zero matrix of the same size as m.
func ZeroLike(m *Matrix) *Matrix {
	rows, cols := m.Dim()
	return &Matrix{
		Rows: rows,
		Cols: cols,
		Data: make([]complex128, rows*cols),
	}
}

// Identity returns an identity matrix.
func Identity(size int) *Matrix {
	m := Zero(size, size)
	for i := range size {
		m.Set(i, i, 1)
	}

	return m
}

// Clone returns a copy of m.
func (m *Matrix) Clone() *Matrix {
	out := ZeroLike(m)
	copy(out.Data, m.Data)
	return out
}

// At returns the value at (i, j).
func (m *Matrix) At(i, j int) complex128 {
	return m.Data[i*m.Cols+j]
}

// Row returns the row at i.
func (m *Matrix) Row(i int) []complex128 {
	return row(m.Data, m.Cols, i)
}

// Set sets the value at (i, j).
func (m *Matrix) Set(i, j int, z complex128) {
	m.Data[i*m.Cols+j] = z
}

// AddAt adds z to the value at (i, j).
func (m *Matrix) AddAt(i, j int, z complex128) {
	m.Data[i*m.Cols+j] += z
}

// SubAt subtracts z from the value at (i, j).
func (m *Matrix) SubAt(i, j int, z complex128) {
	m.Data[i*m.Cols+j] -= z
}

// MulAt multiplies the value at (i, j) by z.
func (m *Matrix) MulAt(i, j int, z complex128) {
	m.Data[i*m.Cols+j] *= z
}

// DivAt divides the value at (i, j) by z.
func (m *Matrix) DivAt(i, j int, z complex128) {
	m.Data[i*m.Cols+j] /= z
}

// Fdiag applies f to the diagonal elements of m.
func (m *Matrix) Fdiag(f func(v complex128) complex128) {
	for i := range min(m.Rows, m.Cols) {
		m.Set(i, i, f(m.At(i, i)))
	}
}

// Seq2 returns a sequence of rows.
func (m *Matrix) Seq2() iter.Seq2[int, []complex128] {
	return func(yield func(int, []complex128) bool) {
		for i := range m.Rows {
			if !yield(i, m.Row(i)) {
				return
			}
		}
	}
}

// Dim returns the dimensions of m.
func (m *Matrix) Dim() (rows int, cols int) {
	return m.Rows, m.Cols
}

// Transpose returns the transpose of m.
func (m *Matrix) Transpose() *Matrix {
	rows, cols := m.Dim()

	out := Zero(cols, rows)
	for i := range rows {
		for j := range cols {
			out.Set(j, i, m.At(i, j))
		}
	}

	return out
}

// Conjugate returns the complex conjugate of m.
func (m *Matrix) Conjugate() *Matrix {
	out := ZeroLike(m)
	for i := range out.Data {
		out.Data[i] = cmplx.Conj(m.Data[i])
	}

	return out
}

// Dagger returns the conjugate transpose of m.
func (m *Matrix) Dagger() *Matrix {
	out := ZeroLike(m)
	for i := range m.Rows {
		for j := range m.Cols {
			out.Set(j, i, cmplx.Conj(m.At(i, j)))
		}
	}

	return out
}

// Equal returns true if m equals n.
func (m *Matrix) Equal(n *Matrix, tol ...float64) bool {
	p, q := m.Dim()
	a, b := n.Dim()

	if a != p || b != q {
		return false
	}

	for i := range m.Data {
		if !epsilon.IsClose(m.Data[i], n.Data[i], tol...) {
			return false
		}
	}

	return true
}

// IsSquare returns true if m is a square matrix.
func (m *Matrix) IsSquare() bool {
	return m.Rows == m.Cols
}

// IsIdentity returns true if m is the identity matrix.
func (m *Matrix) IsIdentity(tol ...float64) bool {
	return m.IsSquare() && m.Equal(Identity(m.Rows), tol...)
}

// IsHermitian returns true if m is a Hermitian matrix.
func (m *Matrix) IsHermitian(tol ...float64) bool {
	return m.IsSquare() && m.Equal(m.Dagger(), tol...)
}

// IsUnitary returns true if m is a unitary matrix.
func (m *Matrix) IsUnitary(tol ...float64) bool {
	return m.IsSquare() && m.MatMul(m.Dagger()).Equal(Identity(m.Rows), tol...)
}

// IsDiagonal returns true if m is a diagonal matrix.
func (m *Matrix) IsDiagonal(tol ...float64) bool {
	for i := range m.Rows {
		for j := range m.Cols {
			if i == j {
				continue
			}

			if !epsilon.IsZero(m.At(i, j), tol...) {
				return false
			}
		}
	}

	return true
}

// Apply returns a matrix product of m and n.
// A.Apply(B) is BA.
// For example, to compute XHZ|v>, you can write v.Apply(Z).Apply(H).Apply(X).
func (m *Matrix) Apply(n *Matrix) *Matrix {
	return n.MatMul(m)
}

// MatMul returns the matrix product of m and n.
// A.MatMul(B) is AB.
func (m *Matrix) MatMul(n *Matrix) *Matrix {
	a, b := m.Dim()
	_, p := n.Dim()

	out := Zero(a, p)
	for i := range a {
		for k := range b {
			mik := m.Data[i*b+k]
			for j := range p {
				out.Data[i*p+j] += mik * n.Data[k*p+j]
			}
		}
	}

	return out
}

// Mul returns z * m.
func (m *Matrix) Mul(z complex128) *Matrix {
	out := ZeroLike(m)
	for i := range m.Data {
		out.Data[i] = m.Data[i] * z
	}

	return out
}

// Add returns m + n.
func (m *Matrix) Add(n *Matrix) *Matrix {
	out := ZeroLike(m)
	for i := range m.Data {
		out.Data[i] = m.Data[i] + n.Data[i]
	}

	return out
}

// Sub returns m - n.
func (m *Matrix) Sub(n *Matrix) *Matrix {
	out := ZeroLike(m)
	for i := range m.Data {
		out.Data[i] = m.Data[i] - n.Data[i]
	}

	return out
}

// Trace returns the trace of m.
func (m *Matrix) Trace() complex128 {
	var z complex128
	for i := range m.Rows {
		z = z + m.At(i, i)
	}

	return z
}

// Real returns the real part of m.
func (m *Matrix) Real() [][]float64 {
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

// Imag returns the imaginary part of m.
func (m *Matrix) Imag() [][]float64 {
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

// Inverse returns the inverse of m.
// It supports only square, non-singular matrices.
func (m *Matrix) Inverse(tol ...float64) *Matrix {
	p, q := m.Dim()
	mm := m.Clone()

	out := Identity(p)
	for i := range p {
		if epsilon.IsZero(mm.At(i, i), tol...) {
			// swap rows
			for r := i + 1; r < p; r++ {
				if epsilon.IsZero(mm.At(r, i), tol...) {
					continue
				}

				mm = mm.Swap(i, r)
				out = out.Swap(i, r)
				break
			}
		}

		c := 1 / mm.At(i, i)
		for j := range q {
			mm.MulAt(i, j, c)
			out.MulAt(i, j, c)
		}

		for j := range q {
			if i == j {
				continue
			}

			c := mm.At(j, i)
			for k := range q {
				mm.AddAt(j, k, -c*mm.At(i, k))
				out.AddAt(j, k, -c*out.At(i, k))
			}
		}
	}

	return out
}

// Swap returns a copy of m with the i-th and j-th rows swapped.
func (m *Matrix) Swap(i, j int) *Matrix {
	data := make([]complex128, len(m.Data))
	copy(data, m.Data)

	i0, i1 := i*m.Cols, (i+1)*m.Cols
	j0, j1 := j*m.Cols, (j+1)*m.Cols

	tmp := make([]complex128, m.Cols)
	copy(tmp, data[i0:i1])
	copy(data[i0:i1], data[j0:j1])
	copy(data[j0:j1], tmp)

	return &Matrix{
		Rows: m.Rows,
		Cols: m.Cols,
		Data: data,
	}
}

// TensorProduct returns a tensor product of m and n.
func (m *Matrix) TensorProduct(n *Matrix) *Matrix {
	p, q := m.Dim()
	a, b := n.Dim()
	rows, cols := p*a, q*b

	data := make([]complex128, rows*cols)
	for i := range p {
		for j := range q {
			mij := m.Data[i*q+j]
			for k := range a {
				for l := range b {
					row, col := i*a+k, j*b+l
					data[row*cols+col] = mij * n.Data[k*b+l]
				}
			}
		}
	}

	return &Matrix{
		Rows: rows,
		Cols: cols,
		Data: data,
	}
}

// MatMul returns a matrix product of m1, m2, ..., mn.
// MatMul(A, B, C, D, ...) is ABCD....
func MatMul(m ...*Matrix) *Matrix {
	out := m[0]
	for i := 1; i < len(m); i++ {
		out = out.MatMul(m[i])
	}

	return out
}

// Apply returns a matrix product of m1, m2, ..., mn.
// Apply(A, B, C, D, ...) is ...DCBA.
func Apply(m ...*Matrix) *Matrix {
	out := m[0]
	for i := 1; i < len(m); i++ {
		out = out.Apply(m[i])
	}

	return out
}

// ApplyN returns the result of applying m n times.
// If n is 0, it returns the identity matrix.
func ApplyN(m *Matrix, n int) *Matrix {
	if n == 0 {
		return Identity(m.Rows)
	}

	list := make([]*Matrix, n)
	for i := range n {
		list[i] = m
	}

	return Apply(list...)
}

// TensorProduct returns a tensor product of m1, m2, ..., mn.
func TensorProduct(m ...*Matrix) *Matrix {
	out := m[0]
	for i := 1; i < len(m); i++ {
		out = out.TensorProduct(m[i])
	}

	return out
}

// TensorProductN returns the n-fold tensor product of m with itself.
func TensorProductN(m *Matrix, n ...int) *Matrix {
	if len(n) < 1 {
		return m
	}

	list := make([]*Matrix, n[0])
	for i := range n[0] {
		list[i] = m
	}

	return TensorProduct(list...)
}

// Commutator returns [m, n].
func Commutator(m, n *Matrix) *Matrix {
	mn, nm := MatMul(m, n), MatMul(n, m)
	return mn.Sub(nm)
}

// AntiCommutator returns {m, n}.
func AntiCommutator(m, n *Matrix) *Matrix {
	mn, nm := MatMul(m, n), MatMul(n, m)
	return mn.Add(nm)
}

func row[T any](arr []T, cols, i int) []T {
	return arr[i*cols : (i+1)*cols]
}
