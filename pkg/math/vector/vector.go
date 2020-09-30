package vector

import (
	"fmt"
	"math/cmplx"

	"github.com/itsubaki/q/pkg/math/matrix"
)

type Vector []complex128

func New(z ...complex128) Vector {
	out := Vector{}
	for _, zi := range z {
		out = append(out, zi)
	}

	return out
}

func Zero(n int) Vector {
	out := Vector{}
	for i := 0; i < n; i++ {
		out = append(out, 0)
	}

	return out
}

func (v Vector) Complex() []complex128 {
	return []complex128(v)
}

func (v Vector) Clone() Vector {
	clone := Vector{}
	for i := 0; i < len(v); i++ {
		clone = append(clone, v[i])
	}

	return clone
}

func (v Vector) Dual() Vector {
	out := Vector{}
	for i := 0; i < len(v); i++ {
		out = append(out, cmplx.Conj(v[i]))
	}

	return out
}

func (v Vector) Add(v1 Vector) Vector {
	out := Vector{}
	for i := 0; i < len(v); i++ {
		out = append(out, v[i]+v1[i])
	}

	return out
}

func (v Vector) Mul(z complex128) Vector {
	out := Vector{}
	for i := range v {
		out = append(out, z*v[i])
	}

	return out
}

func (v Vector) TensorProduct(v1 Vector) Vector {
	out := Vector{}
	for i := 0; i < len(v); i++ {
		for j := 0; j < len(v1); j++ {
			out = append(out, v[i]*v1[j])
		}
	}

	return out
}

func (v Vector) InnerProduct(v1 Vector) complex128 {
	dual := v1.Dual()

	out := complex(0, 0)
	for i := 0; i < len(v); i++ {
		out = out + v[i]*dual[i]
	}

	return out
}

func (v Vector) OuterProduct(v1 Vector) matrix.Matrix {
	dual := v1.Dual()

	out := matrix.Matrix{}
	for i := 0; i < len(v); i++ {
		vv := make([]complex128, 0)
		for j := 0; j < len(dual); j++ {
			vv = append(vv, v[i]*dual[j])
		}

		out = append(out, vv)
	}

	return out
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

func (v Vector) Apply(m matrix.Matrix) Vector {
	p, q := m.Dimension()

	if q != len(v) {
		panic(fmt.Sprintf("invalid dimension. p=%d q=%d len(v)=%d", p, q, len(v)))
	}

	out := Vector{}
	for i := 0; i < p; i++ {
		tmp := complex(0, 0)
		for j := 0; j < q; j++ {
			tmp = tmp + m[i][j]*v[j]
		}

		out = append(out, tmp)
	}

	return out
}

func (v Vector) Equals(v1 Vector, eps ...float64) bool {
	if len(v) != len(v1) {
		return false
	}

	e := 1e-13
	if len(eps) > 0 {
		e = eps[0]
	}

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

func (v Vector) Real() []float64 {
	out := make([]float64, 0)
	for i := range v {
		out = append(out, real(v[i]))
	}

	return out
}

func (v Vector) Imag() []float64 {
	out := make([]float64, 0)
	for i := range v {
		out = append(out, imag(v[i]))
	}

	return out
}

func TensorProductN(v Vector, bit ...int) Vector {
	if len(bit) < 1 {
		return v
	}

	list := make([]Vector, 0)
	for i := 0; i < bit[0]; i++ {
		list = append(list, v)
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
