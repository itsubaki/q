package vector

import (
	"fmt"
	"math/cmplx"

	"github.com/itsubaki/q/math/epsilon"
	"github.com/itsubaki/q/math/matrix"
)

// Vector is a vector of complex128.
type Vector struct {
	Data []complex128
}

// New returns a new vector of complex128.
func New(z ...complex128) *Vector {
	return &Vector{
		Data: z,
	}
}

// Zero returns a vector of length n with all elements zero.
func Zero(n int) *Vector {
	return &Vector{
		Data: make([]complex128, n),
	}
}

// Clone returns a clone of vector.
func (v *Vector) Clone() *Vector {
	data := make([]complex128, len(v.Data))
	copy(data, v.Data)
	return &Vector{
		Data: data,
	}
}

// Dual returns a dual vector.
func (v *Vector) Dual() *Vector {
	data := make([]complex128, len(v.Data))
	for i := range v.Data {
		data[i] = cmplx.Conj(v.Data[i])
	}

	return &Vector{
		Data: data,
	}
}

// Add returns a vector of v+w.
func (v *Vector) Add(w *Vector) *Vector {
	data := make([]complex128, len(v.Data))
	for i := range v.Data {
		data[i] = v.Data[i] + w.Data[i]
	}

	return &Vector{
		Data: data,
	}
}

// Mul returns a vector of z*v.
func (v *Vector) Mul(z complex128) *Vector {
	data := make([]complex128, len(v.Data))
	for i := range v.Data {
		data[i] = z * v.Data[i]
	}

	return &Vector{
		Data: data,
	}
}

// TensorProduct returns the tensor product of v and w.
func (v *Vector) TensorProduct(w *Vector) *Vector {
	p, q := len(v.Data), len(w.Data)

	data := make([]complex128, p*q)
	for i := range p {
		for j := range q {
			data[i*q+j] = v.Data[i] * w.Data[j]
		}
	}

	return &Vector{
		Data: data,
	}
}

// InnerProduct returns the inner product of v and w.
func (v *Vector) InnerProduct(w *Vector) complex128 {
	dual := w.Dual()

	var z complex128
	for i := range v.Data {
		z = z + v.Data[i]*dual.Data[i]
	}

	return z
}

// OuterProduct returns the outer product of v and w.
func (v *Vector) OuterProduct(w *Vector) *matrix.Matrix {
	rows, cols := len(v.Data), len(w.Data)
	dual := w.Dual()

	data := make([]complex128, rows*cols)
	for i := range v.Data {
		for j := range dual.Data {
			data[i*cols+j] = v.Data[i] * dual.Data[j]
		}
	}

	return &matrix.Matrix{
		Rows: rows,
		Cols: cols,
		Data: data,
	}
}

// IsOrthogonal returns true if v and w are orthogonal.
func (v *Vector) IsOrthogonal(w *Vector) bool {
	return v.InnerProduct(w) == 0
}

// Norm returns a norm of vector.
func (v *Vector) Norm() complex128 {
	return cmplx.Sqrt(v.InnerProduct(v))
}

// IsUnit returns true if v is unit vector.
func (v *Vector) IsUnit() bool {
	return v.Norm() == 1
}

// Apply returns a matrix product of v and m.
// v.Apply(X) -> X|v>
func (v *Vector) Apply(m *matrix.Matrix) *Vector {
	p, q := m.Dimension()

	data := make([]complex128, p)
	for i := range p {
		for j := range q {
			data[i] += m.At(i, j) * v.Data[j]
		}
	}

	return &Vector{
		Data: data,
	}
}

// Equals returns true if v and w are equal.
// If eps is not given, epsilon.E13 is used.
func (v *Vector) Equals(w *Vector, eps ...float64) bool {
	if len(v.Data) != len(w.Data) {
		return false
	}

	e := epsilon.E13(eps...)
	for i := range v.Data {
		if cmplx.Abs(v.Data[i]-w.Data[i]) > e {
			return false
		}
	}

	return true
}

// Real returns a slice of real part.
func (v *Vector) Real() []float64 {
	data := make([]float64, len(v.Data))
	for i := range v.Data {
		data[i] = real(v.Data[i])
	}

	return data
}

// Imag returns a slice of imaginary part.
func (v *Vector) Imag() []float64 {
	data := make([]float64, len(v.Data))
	for i := range v.Data {
		data[i] = imag(v.Data[i])
	}

	return data
}

// String returns the string representation of v.
func (v *Vector) String() string {
	return fmt.Sprintf("%v", v.Data)
}

func TensorProductN(v *Vector, n ...int) *Vector {
	if len(n) < 1 {
		return v
	}

	list := make([]*Vector, n[0])
	for i := range n[0] {
		list[i] = v
	}

	out := list[0]
	for i := 1; i < len(list); i++ {
		out = out.TensorProduct(list[i])
	}

	return out
}
