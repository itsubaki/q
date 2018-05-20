package gate

import (
	"fmt"
	"math"
	"math/cmplx"
	"strconv"
	"strings"

	"github.com/itsubaki/q/matrix"
)

func U(alpha, beta, gamma, delta float64) matrix.Matrix {
	m0 := make(matrix.Matrix, 2)
	v0 := complex(0, 1*beta/2)
	m0[0] = []complex128{cmplx.Exp(cmplx.Conj(v0)), 0}
	m0[1] = []complex128{0, cmplx.Exp(v0)}

	m1 := make(matrix.Matrix, 2)
	v1 := complex(gamma/2, 0)
	m1[0] = []complex128{cmplx.Cos(v1), -1 * cmplx.Sin(v1)}
	m1[1] = []complex128{cmplx.Sin(v1), cmplx.Cos(v1)}

	m2 := make(matrix.Matrix, 2)
	v2 := complex(0, 1*delta/2)
	m2[0] = []complex128{cmplx.Exp(cmplx.Conj(v2)), 0}
	m2[1] = []complex128{0, cmplx.Exp(v2)}

	u := m0.Apply(m1).Apply(m2)
	return u.Mul(cmplx.Exp(complex(0, alpha)))
}

func I(bit ...int) matrix.Matrix {
	m := make(matrix.Matrix, 2)
	m[0] = []complex128{1, 0}
	m[1] = []complex128{0, 1}
	return matrix.Tensor(m, bit...)
}

func X(bit ...int) matrix.Matrix {
	m := make(matrix.Matrix, 2)
	m[0] = []complex128{0, 1}
	m[1] = []complex128{1, 0}
	return matrix.Tensor(m, bit...)
}

func Y(bit ...int) matrix.Matrix {
	m := make(matrix.Matrix, 2)
	m[0] = []complex128{0, -1i}
	m[1] = []complex128{1i, 0}
	return matrix.Tensor(m, bit...)
}

func Z(bit ...int) matrix.Matrix {
	m := make(matrix.Matrix, 2)
	m[0] = []complex128{1, 0}
	m[1] = []complex128{0, -1}
	return matrix.Tensor(m, bit...)
}

func H(bit ...int) matrix.Matrix {
	m := make(matrix.Matrix, 2)
	v := complex(1/math.Sqrt2, 0)
	m[0] = []complex128{v, v}
	m[1] = []complex128{v, -1 * v}
	return matrix.Tensor(m, bit...)
}

func S(bit ...int) matrix.Matrix {
	m := make(matrix.Matrix, 2)
	m[0] = []complex128{1, 0}
	m[1] = []complex128{0, 1i}
	return matrix.Tensor(m, bit...)
}

func T(bit ...int) matrix.Matrix {
	m := make(matrix.Matrix, 2)
	v := cmplx.Exp(complex(0, 1) * math.Pi / 4)
	m[0] = []complex128{1, 0}
	m[1] = []complex128{0, v}
	return matrix.Tensor(m, bit...)
}

func Swap(bit ...int) matrix.Matrix {
	if len(bit) < 1 {
		bit = []int{2}
	}

	m := I(bit...)
	dim := len(m)

	for i := 1; i < dim/2; i++ {
		m[i], m[dim-1-i] = m[dim-1-i], m[i]
	}

	return m
}

// CZ(2) -> Controlled-Z
// CZ(3) -> Contrlled-Controlled-Z
func CZ(bit ...int) matrix.Matrix {
	if len(bit) < 1 {
		bit = []int{2}
	}

	m := I(bit...)
	dim := len(m)

	m[dim-1][dim-1] = -1
	return m
}

func CS(bit ...int) matrix.Matrix {
	if len(bit) < 1 {
		bit = []int{2}
	}

	m := I(bit...)
	dim := len(m)

	m[dim-1][dim-1] = 1i
	return m
}

// CNOT(3) -> Toffoli (Controlled-Controlled-NOT)
func CNOT(bit ...int) matrix.Matrix {
	if len(bit) < 1 {
		bit = []int{2}
	}

	m := I(bit...)
	dim := len(m)

	m[dim-1], m[dim-2] = m[dim-2], m[dim-1]
	return m
}

func Fredkin() matrix.Matrix {
	m := make(matrix.Matrix, 8)
	m[0] = []complex128{1, 0, 0, 0, 0, 0, 0, 0}
	m[1] = []complex128{0, 1, 0, 0, 0, 0, 0, 0}
	m[2] = []complex128{0, 0, 1, 0, 0, 0, 0, 0}
	m[3] = []complex128{0, 0, 0, 1, 0, 0, 0, 0}
	m[4] = []complex128{0, 0, 0, 0, 1, 0, 0, 0}
	m[5] = []complex128{0, 0, 0, 0, 0, 0, 1, 0}
	m[6] = []complex128{0, 0, 0, 0, 0, 1, 0, 0}
	m[7] = []complex128{0, 0, 0, 0, 0, 0, 0, 1}
	return m
}

func ControlledNot(bit, c, t int) matrix.Matrix {
	m := I([]int{bit}...)
	dim := len(m)

	index := []int64{}
	f := "%0" + strconv.Itoa(bit) + "s"
	for i := 0; i < dim; i++ {
		s := fmt.Sprintf(f, strconv.FormatInt(int64(i), 2))

		bits := []string{}
		for j := 0; j < bit; j++ {
			bits = append(bits, s[j:j+1])
		}

		// Apply X
		if bits[c] == "1" {
			if bits[t] == "1" {
				bits[t] = "0"
			} else if bits[t] == "0" {
				bits[t] = "1"
			}
		}

		v, err := strconv.ParseInt(strings.Join(bits, ""), 2, 0)
		if err != nil {
			panic(err)
		}

		index = append(index, v)
	}

	cnot := make(matrix.Matrix, dim)
	for i, ii := range index {
		cnot[i] = m[ii]
	}

	return cnot
}

func ControlledZ(bit, c, t int) matrix.Matrix {
	m := I([]int{bit}...)
	dim := len(m)

	f := "%0" + strconv.Itoa(bit) + "s"
	for i := 0; i < dim; i++ {
		s := fmt.Sprintf(f, strconv.FormatInt(int64(i), 2))

		bits := []string{}
		for j := 0; j < bit; j++ {
			bits = append(bits, s[j:j+1])
		}

		// Apply Z
		if bits[c] == "1" && bits[t] == "1" {
			for j := 0; j < dim; j++ {
				m[i][j] = complex(-1, 0) * m[i][j]
			}
		}
	}

	return m
}

func QFT() matrix.Matrix {
	m := make(matrix.Matrix, 8)

	o0 := complex(1, 0)
	o1 := cmplx.Sqrt(1i)
	o2 := cmplx.Pow(o1, 2)
	o3 := cmplx.Pow(o1, 3)
	o4 := cmplx.Pow(o1, 4)
	o5 := cmplx.Pow(o1, 5)
	o6 := cmplx.Pow(o1, 6)
	o7 := cmplx.Pow(o1, 7)

	m[0] = []complex128{o0, o0, o0, o0, o0, o0, o0, o0}
	m[1] = []complex128{o0, o1, o2, o3, o4, o5, o6, o7}
	m[2] = []complex128{o0, o2, o4, o6, o0, o2, o4, o6}
	m[3] = []complex128{o0, o3, o6, o1, o4, o7, o2, o5}
	m[4] = []complex128{o0, o4, o0, o4, o0, o4, o0, o4}
	m[5] = []complex128{o0, o5, o2, o7, o4, o1, o6, o3}
	m[6] = []complex128{o0, o6, o4, o2, o0, o6, o4, o2}
	m[7] = []complex128{o0, o7, o6, o5, o4, o3, o2, o1}
	return m.Mul(complex(1/math.Sqrt(8), 0))
}
