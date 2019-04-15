package vector

import (
	"math/cmplx"

	"github.com/axamon/q/matrix"
)

// Vector type manages complex number slice in vector
type Vector []complex128

// New creates a new vector with dimensions equal to the length of complex
// numbers passed as arguments.
func New(z ...complex128) Vector {
	v := Vector{}
	for _, zi := range z {
		v = append(v, zi)
	}
	return v
}

// NewZero creates a new vector with all components set to zero.
func NewZero(n int) Vector {
	v := Vector{}
	for i := 0; i < n; i++ {
		v = append(v, 0)
	}
	return v
}

// Clone clones the vector.
func (v0 Vector) Clone() Vector {
	clone := Vector{}
	for i := 0; i < len(v0); i++ {
		clone = append(clone, v0[i])
	}
	return clone
}

// Dual returns the complex conjugate of the vector.
func (v0 Vector) Dual() Vector {
	dual := Vector{}
	for i := 0; i < len(v0); i++ {
		dual = append(dual, cmplx.Conj(v0[i]))
	}
	return dual
}

// Add adds the vector to the first one.
func (v0 Vector) Add(v1 Vector) Vector {
	v2 := Vector{}
	for i := 0; i < len(v0); i++ {
		v2 = append(v2, v0[i]+v1[i])
	}
	return v2
}

// Mul multiplies the vector by the complex number.
func (v0 Vector) Mul(z complex128) Vector {
	v2 := Vector{}
	for i := range v0 {
		v2 = append(v2, z*v0[i])
	}
	return v2
}

// TensorProduct returns the vector that results from applying
// the tensor product between the two vectors.
func (v0 Vector) TensorProduct(v1 Vector) Vector {
	v2 := Vector{}
	for i := 0; i < len(v0); i++ {
		for j := 0; j < len(v1); j++ {
			v2 = append(v2, v0[i]*v1[j])
		}
	}
	return v2
}

// InnerProduct returns the complex number that results from applying
// the inner product between the two vectors.
func (v0 Vector) InnerProduct(v1 Vector) complex128 {
	p := complex(0, 0)

	dual := v1.Dual()
	for i := 0; i < len(v0); i++ {
		p = p + v0[i]*dual[i]
	}

	return p
}

// OuterProduct returns the matrix that results from applying
// the outer product between the two vectors.
func (v0 Vector) OuterProduct(v1 Vector) matrix.Matrix {
	m := matrix.Matrix{}
	for i := 0; i < len(v0); i++ {
		v := []complex128{}
		for j := 0; j < len(v1); j++ {
			v = append(v, v0[i]*v1[j])
		}
		m = append(m, v)
	}

	return m
}

// IsOrthogonal returns true if the two vectors are othogonal to
// each other: their inner product is zero.
func (v0 Vector) IsOrthogonal(v1 Vector) bool {
	return v0.InnerProduct(v1) == complex(0, 0)
}

// Norm returns the normalized value of the vector elements.
func (v0 Vector) Norm() complex128 {
	return cmplx.Sqrt(v0.InnerProduct(v0))
}

// IsUnit returns true if the normalized value of the vector elements is 1.
func (v0 Vector) IsUnit() bool {
	return v0.Norm() == complex(1, 0)
}

// Apply returns the vector resulting from multiplying the vector by the matrix.
func (v0 Vector) Apply(mat matrix.Matrix) Vector {
	v := Vector{}

	m, _ := mat.Dimension()
	for i := 0; i < m; i++ {
		tmp := complex(0, 0)
		for j := 0; j < len(v0); j++ {
			tmp = tmp + mat[i][j]*v0[j]
		}
		v = append(v, tmp)
	}

	return v
}

// Equals returns true if the two vectors are the same.
func (v0 Vector) Equals(v1 Vector, eps ...float64) bool {
	if len(v0) != len(v1) {
		return false
	}

	e := matrix.Eps(eps...)
	for i := 0; i < len(v0); i++ {
		if cmplx.Abs(v0[i]-v1[i]) > e {
			return false
		}

	}
	return true
}

// Dimension returns the dimension of the vector
func (v0 Vector) Dimension() int {
	return len(v0)
}

// TensorProductN returns the vector resulting from applying the
// tensor product beetween the vector and itself N times.
func TensorProductN(v0 Vector, bit ...int) Vector {
	if len(bit) < 1 {
		return v0
	}

	v1 := v0
	for i := 1; i < bit[0]; i++ {
		v1 = v1.TensorProduct(v0)
	}
	return v1
}

// TensorProduct returns the vector resulting from multipling the
// all vectors passed as arguments together.
func TensorProduct(v0 ...Vector) Vector {
	v1 := v0[0]
	for i := 1; i < len(v0); i++ {
		v1 = v1.TensorProduct(v0[i])
	}
	return v1
}
