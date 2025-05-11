package vector

import (
	"math/cmplx"

	"github.com/itsubaki/q/math/epsilon"
	"github.com/itsubaki/q/math/matrix"
)

// Vector is a vector of complex128.
type Vector []complex128

// New returns a new vector of complex128.
func New(z ...complex128) Vector {
	out := make(Vector, len(z))
	copy(out, z)
	return out
}

// Zero returns a vector of length n with all elements zero.
func Zero(n int) Vector {
	return make(Vector, n)
}

// Complex returns a slice of complex128.
func (v Vector) Complex() []complex128 {
	return []complex128(v)
}

// Clone returns a clone of vector.
func (v Vector) Clone() Vector {
	return New(v...)
}

// Dual returns a dual vector.
func (v Vector) Dual() Vector {
	out := make(Vector, len(v))
	for i := range v {
		out[i] = cmplx.Conj(v[i])
	}

	return out
}

// Add returns a vector of v+w.
func (v Vector) Add(w Vector) Vector {
	out := make(Vector, len(v))
	for i := range v {
		out[i] = v[i] + w[i]
	}

	return out
}

// Mul returns a vector of z*v.
func (v Vector) Mul(z complex128) Vector {
	out := make(Vector, len(v))
	for i := range v {
		out[i] = z * v[i]
	}

	return out
}

// TensorProduct returns the tensor product of v and w.
func (v Vector) TensorProduct(w Vector) Vector {
	p, q := len(v), len(w)

	out := make(Vector, p*q)
	for i := range p {
		for j := range q {
			out[i*q+j] = v[i] * w[j]
		}
	}

	return out
}

// InnerProduct returns the inner product of v and w.
func (v Vector) InnerProduct(w Vector) complex128 {
	dual := w.Dual()

	var out complex128
	for i := range v {
		out = out + v[i]*dual[i]
	}

	return out
}

// OuterProduct returns the outer product of v and w.
func (v Vector) OuterProduct(w Vector) matrix.Matrix {
	rows, cols := len(v), len(w)
	dual := w.Dual()

	data := make([]complex128, rows*cols)
	for i := range v {
		for j := range dual {
			data[i*cols+j] = v[i] * dual[j]
		}
	}

	return matrix.Matrix{
		Rows: rows,
		Cols: cols,
		Data: data,
	}
}

// IsOrthogonal returns true if v and w are orthogonal.
func (v Vector) IsOrthogonal(w Vector) bool {
	return v.InnerProduct(w) == 0
}

// Norm returns a norm of vector.
func (v Vector) Norm() complex128 {
	return cmplx.Sqrt(v.InnerProduct(v))
}

// IsUnit returns true if v is unit vector.
func (v Vector) IsUnit() bool {
	return v.Norm() == 1
}

// Apply returns a matrix product of v and m.
func (v Vector) Apply(m matrix.Matrix) Vector {
	p, q := m.Dimension()

	out := make(Vector, p)
	for i := range p {
		var c complex128
		for j := range q {
			c = c + m.At(i, j)*v[j]
		}

		out[i] = c
	}

	return out
}

// Equals returns true if v and w are equal.
// If eps is not given, epsilon.E13 is used.
func (v Vector) Equals(w Vector, eps ...float64) bool {
	if len(v) != len(w) {
		return false
	}

	e := epsilon.E13(eps...)
	for i := range v {
		if cmplx.Abs(v[i]-w[i]) > e {
			return false
		}
	}

	return true
}

// Dimension returns a dimension of vector.
func (v Vector) Dimension() int {
	return len(v)
}

// Real returns a slice of real part.
func (v Vector) Real() []float64 {
	out := make([]float64, len(v))
	for i := range v {
		out[i] = real(v[i])
	}

	return out
}

// Imag returns a slice of imaginary part.
func (v Vector) Imag() []float64 {
	out := make([]float64, len(v))
	for i := range v {
		out[i] = imag(v[i])
	}

	return out
}

func TensorProductN(v Vector, n ...int) Vector {
	if len(n) < 1 {
		return v
	}

	list := make([]Vector, n[0])
	for i := range n[0] {
		list[i] = v
	}

	return TensorProduct(list...)
}

func TensorProduct(v ...Vector) Vector {
	out := v[0]
	for i := 1; i < len(v); i++ {
		out = out.TensorProduct(v[i])
	}

	return out
}
