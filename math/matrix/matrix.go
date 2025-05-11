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
func New(z ...[]complex128) Matrix {
	rows := len(z)
	var cols int
	if rows > 0 {
		cols = len(z[0])
	}

	data := make([]complex128, rows*cols)
	for i := range rows {
		copy(data[i*cols:(i+1)*cols], z[i])
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

// ZeroLike returns a zero matrix of same size as m.
func ZeroLike(m Matrix) Matrix {
	rows, cols := m.Dimension()
	return Matrix{
		Rows: rows,
		Cols: cols,
		Data: make([]complex128, rows*cols),
	}
}

// Identity returns an identity matrix.
func Identity(size int) Matrix {
	m := Zero(size, size)
	for i := range size {
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
func (m Matrix) Set(i, j int, z complex128) {
	m.Data[i*m.Cols+j] = z
}

// AddAt adds a value of matrix at (i,j).
func (m Matrix) AddAt(i, j int, z complex128) {
	m.Data[i*m.Cols+j] += z
}

// MulAt multiplies a value of matrix at (i,j).
func (m Matrix) MulAt(i, j int, z complex128) {
	m.Data[i*m.Cols+j] *= z
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

// Dimension returns a dimension of matrix.
func (m Matrix) Dimension() (rows int, cols int) {
	return m.Rows, m.Cols
}

// Transpose returns a transpose matrix.
func (m Matrix) Transpose() Matrix {
	out := ZeroLike(m)
	for i := range m.Rows {
		for j := range m.Cols {
			out.Set(j, i, m.At(i, j))
		}
	}

	return out
}

// Conjugate returns a conjugate matrix.
func (m Matrix) Conjugate() Matrix {
	out := ZeroLike(m)
	for i := range out.Data {
		out.Data[i] = cmplx.Conj(m.Data[i])
	}

	return out
}

// Dagger returns conjugate transpose matrix.
func (m Matrix) Dagger() Matrix {
	out := ZeroLike(m)
	for i := range m.Rows {
		for j := range m.Cols {
			out.Set(j, i, cmplx.Conj(m.At(i, j)))
		}
	}

	return out
}

// Equals returns true if m equals n.
// If eps is not given, epsilon.E13 is used.
func (m Matrix) Equals(n Matrix, eps ...float64) bool {
	p, q := m.Dimension()
	a, b := n.Dimension()

	if a != p || b != q {
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
	id := Identity(m.Rows)
	return mmd.Equals(id, epsilon.E13(eps...))
}

// Apply returns a matrix product of m and n.
// A.Apply(B) is BA.
func (m Matrix) Apply(n Matrix) Matrix {
	a, b := m.Dimension()
	_, p := n.Dimension()

	out := Zero(a, p)
	for i := range a {
		for k := range b {
			nik := n.At(i, k)
			for j := range p {
				out.AddAt(i, j, nik*m.At(k, j))
			}
		}
	}

	return out
}

// Mul returns a matrix of z*m.
func (m Matrix) Mul(z complex128) Matrix {
	out := ZeroLike(m)
	for i := range m.Data {
		out.Data[i] = m.Data[i] * z
	}

	return out
}

// Add returns a matrix of m+n.
func (m Matrix) Add(n Matrix) Matrix {
	out := ZeroLike(m)
	for i := range m.Data {
		out.Data[i] = m.Data[i] + n.Data[i]
	}

	return out
}

// Sub returns a matrix of m-n.
func (m Matrix) Sub(n Matrix) Matrix {
	out := ZeroLike(m)
	for i := range m.Data {
		out.Data[i] = m.Data[i] - n.Data[i]
	}

	return out
}

// Trace returns a trace of matrix.
func (m Matrix) Trace() complex128 {
	var z complex128
	for i := range m.Rows {
		z = z + m.At(i, i)
	}

	return z
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
	out := ZeroLike(m)
	copy(out.Data, m.Data)
	return out
}

// Inverse returns an inverse matrix of m.
func (m Matrix) Inverse() Matrix {
	p, q := m.Dimension()
	mm := m.Clone()

	out := Identity(p)
	for i := range p {
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

// TensorProduct returns a tensor product of m and n.
func (m Matrix) TensorProduct(n Matrix) Matrix {
	p, q := m.Dimension()
	a, b := n.Dimension()
	rows, cols := p*a, q*b

	data := make([]complex128, rows*cols)
	for i := range p {
		for j := range q {
			mij := m.At(i, j)
			for k := range a {
				for l := range b {
					row, col := i*a+k, j*b+l
					data[row*cols+col] = mij * n.At(k, l)
				}
			}
		}
	}

	return Matrix{
		Rows: rows,
		Cols: cols,
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
