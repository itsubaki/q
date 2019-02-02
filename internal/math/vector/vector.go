package vector

import (
	"math/cmplx"

	"github.com/itsubaki/q/internal/math/matrix"
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

func (v0 Vector) Clone() Vector {
	clone := Vector{}
	for i := 0; i < len(v0); i++ {
		clone = append(clone, v0[i])
	}
	return clone
}

func (v0 Vector) Dual() Vector {
	dual := Vector{}
	for i := 0; i < len(v0); i++ {
		dual = append(dual, cmplx.Conj(v0[i]))
	}
	return dual
}

func (v0 Vector) Add(v1 Vector) Vector {
	v2 := Vector{}
	for i := 0; i < len(v0); i++ {
		v2 = append(v2, v0[i]+v1[i])
	}
	return v2
}

func (v0 Vector) Mul(z complex128) Vector {
	v2 := Vector{}
	for i := range v0 {
		v2 = append(v2, z*v0[i])
	}
	return v2
}

func (v0 Vector) TensorProduct(v1 Vector) Vector {
	v2 := Vector{}
	for i := 0; i < len(v0); i++ {
		for j := 0; j < len(v1); j++ {
			v2 = append(v2, v0[i]*v1[j])
		}
	}
	return v2
}

func (v0 Vector) InnerProduct(v1 Vector) complex128 {
	p := complex(0, 0)

	dual := v1.Dual()
	for i := 0; i < len(v0); i++ {
		p = p + v0[i]*dual[i]
	}

	return p
}

func (v0 Vector) IsOrthogonal(v1 Vector) bool {
	p := v0.InnerProduct(v1)
	if p == complex(0, 0) {
		return true
	}
	return false
}

func (v0 Vector) Norm() complex128 {
	return cmplx.Sqrt(v0.InnerProduct(v0))
}

func (v0 Vector) IsUnit() bool {
	n := v0.Norm()
	if n == complex(1, 0) {
		return true
	}
	return false
}

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

func (v0 Vector) Dimension() int {
	return len(v0)
}

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

func TensorProduct(v0 ...Vector) Vector {
	v1 := v0[0]
	for i := 1; i < len(v0); i++ {
		v1 = v1.TensorProduct(v0[i])
	}
	return v1
}
