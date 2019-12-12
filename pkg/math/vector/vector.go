package vector

import (
	"math/cmplx"

	"github.com/itsubaki/q/pkg/math/matrix"
)

type Vector []complex128

func New(z ...complex128) Vector {
	v := Vector{}
	for _, zi := range z {
		v = append(v, zi)
	}

	return v
}

func NewZero(n int) Vector {
	v := Vector{}
	for i := 0; i < n; i++ {
		v = append(v, 0)
	}

	return v
}

func (v Vector) Clone() Vector {
	clone := Vector{}
	for i := 0; i < len(v); i++ {
		clone = append(clone, v[i])
	}

	return clone
}

func (v Vector) Dual() Vector {
	dual := Vector{}
	for i := 0; i < len(v); i++ {
		dual = append(dual, cmplx.Conj(v[i]))
	}

	return dual
}

func (v Vector) Add(v1 Vector) Vector {
	add := Vector{}
	for i := 0; i < len(v); i++ {
		add = append(add, v[i]+v1[i])
	}

	return add
}

func (v Vector) Mul(z complex128) Vector {
	mul := Vector{}
	for i := range v {
		mul = append(mul, z*v[i])
	}

	return mul
}

func (v Vector) TensorProduct(v1 Vector) Vector {
	p := Vector{}
	for i := 0; i < len(v); i++ {
		for j := 0; j < len(v1); j++ {
			p = append(p, v[i]*v1[j])
		}
	}

	return p
}

// <v1|v0>
func (v Vector) InnerProduct(v1 Vector) complex128 {
	dual := v1.Dual()

	p := complex(0, 0)
	for i := 0; i < len(v); i++ {
		p = p + v[i]*dual[i]
	}

	return p
}

// |v0><v1|
func (v Vector) OuterProduct(v1 Vector) matrix.Matrix {
	dual := v1.Dual()

	p := matrix.Matrix{}
	for i := 0; i < len(v); i++ {
		vv := make([]complex128, 0)
		for j := 0; j < len(dual); j++ {
			vv = append(vv, v[i]*dual[j])
		}

		p = append(p, vv)
	}

	return p
}

func (v Vector) IsOrthogonal(v1 Vector) bool {
	return v.InnerProduct(v1) == complex(0, 0)
}

func (v Vector) Norm() complex128 {
	return cmplx.Sqrt(v.InnerProduct(v))
}

func (v Vector) IsUnit() bool {
	return v.Norm() == complex(1, 0)
}

func (v Vector) Apply(mat matrix.Matrix) Vector {
	apply := Vector{}

	m, _ := mat.Dimension()
	for i := 0; i < m; i++ {
		tmp := complex(0, 0)
		for j := 0; j < len(v); j++ {
			tmp = tmp + mat[i][j]*v[j]
		}

		apply = append(apply, tmp)
	}

	return apply
}

func (v Vector) Equals(v1 Vector, eps ...float64) bool {
	if len(v) != len(v1) {
		return false
	}

	e := matrix.Eps(eps...)
	for i := 0; i < len(v); i++ {
		if cmplx.Abs(v[i]-v1[i]) > e {
			return false
		}
	}

	return true
}

func (v Vector) Dimension() int {
	return len(v)
}

func TensorProductN(v Vector, bit ...int) Vector {
	if len(bit) < 1 {
		return v
	}

	p := v
	for i := 1; i < bit[0]; i++ {
		p = p.TensorProduct(v)
	}

	return p
}

func TensorProduct(v ...Vector) Vector {
	p := v[0]
	for i := 1; i < len(v); i++ {
		p = p.TensorProduct(v[i])
	}

	return p
}
