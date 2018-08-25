package gate

import (
	"fmt"
	"math"
	"math/cmplx"
	"strconv"

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

func R(k int) matrix.Matrix {
	m := make(matrix.Matrix, 2)

	p := 2 * math.Pi / math.Pow(2, float64(k))
	e := cmplx.Exp(complex(0, p))

	m[0] = []complex128{1, 0}
	m[1] = []complex128{0, e}
	return m
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

func ControlledR(bit int, c []int, t, k int) matrix.Matrix {
	m := I([]int{bit}...)
	dim := len(m)

	p := 2 * math.Pi / math.Pow(2, float64(k))
	e := cmplx.Exp(complex(0, p))

	f := "%0" + strconv.Itoa(bit) + "s"
	for i := 0; i < dim; i++ {
		s := fmt.Sprintf(f, strconv.FormatInt(int64(i), 2))
		bits := []rune(s)

		// Apply R(k)
		apply := true
		for i := range c {
			if bits[c[i]] == '0' {
				apply = false
			}
		}

		if apply && bits[t] == '1' {
			m[i][i] = e * m[i][i]
		}
	}

	return m
}

func CR(bit, c, t, k int) matrix.Matrix {
	return ControlledR(bit, []int{c}, t, k)
}

func ControlledNot(bit int, c []int, t int) matrix.Matrix {
	m := I([]int{bit}...)
	dim := len(m)

	index := []int64{}
	f := "%0" + strconv.Itoa(bit) + "s"
	for i := 0; i < dim; i++ {
		s := fmt.Sprintf(f, strconv.FormatInt(int64(i), 2))
		bits := []rune(s)

		// Apply X
		flip := true
		for i := range c {
			if bits[c[i]] == '0' {
				flip = false
			}
		}

		if flip {
			if bits[t] == '1' {
				bits[t] = '0'
			} else if bits[t] == '0' {
				bits[t] = '1'
			}
		}

		v, err := strconv.ParseInt(string(bits), 2, 0)
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

func Toffoli() matrix.Matrix {
	return ControlledNot(3, []int{0, 1}, 2)
}

func CNOT(bit, c, t int) matrix.Matrix {
	return ControlledNot(bit, []int{c}, t)
}

func ControlledZ(bit int, c []int, t int) matrix.Matrix {
	m := I([]int{bit}...)
	dim := len(m)

	f := "%0" + strconv.Itoa(bit) + "s"
	for i := 0; i < dim; i++ {
		s := fmt.Sprintf(f, strconv.FormatInt(int64(i), 2))
		bits := []rune(s)

		// Apply Z
		pflip := true
		for i := range c {
			if bits[c[i]] == '0' {
				pflip = false
			}
		}

		if pflip && bits[t] == '1' {
			m[i][i] = complex(-1, 0) * m[i][i]
		}
	}

	return m
}

func CZ(bit, c, t int) matrix.Matrix {
	return ControlledZ(bit, []int{c}, t)
}

func ControlledS(bit int, c []int, t int) matrix.Matrix {
	m := I([]int{bit}...)
	dim := len(m)

	f := "%0" + strconv.Itoa(bit) + "s"
	for i := 0; i < dim; i++ {
		s := fmt.Sprintf(f, strconv.FormatInt(int64(i), 2))
		bits := []rune(s)

		// Apply S
		pflip := true
		for i := range c {
			if bits[c[i]] == '0' {
				pflip = false
			}
		}

		if pflip && bits[t] == '1' {
			m[i][i] = 1i * m[i][i]
		}
	}

	return m
}

func CS(bit, c, t int) matrix.Matrix {
	return ControlledS(bit, []int{c}, t)
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

func QFT() matrix.Matrix {
	m := make(matrix.Matrix, 8)
	o := []complex128{}
	for i := 0; i < 8; i++ {
		o = append(o, cmplx.Pow(cmplx.Sqrt(1i), complex(float64(i), 0)))
	}

	m[0] = []complex128{o[0], o[0], o[0], o[0], o[0], o[0], o[0], o[0]}
	m[1] = []complex128{o[0], o[1], o[2], o[3], o[4], o[5], o[6], o[7]}
	m[2] = []complex128{o[0], o[2], o[4], o[6], o[0], o[2], o[4], o[6]}
	m[3] = []complex128{o[0], o[3], o[6], o[1], o[4], o[7], o[2], o[5]}
	m[4] = []complex128{o[0], o[4], o[0], o[4], o[0], o[4], o[0], o[4]}
	m[5] = []complex128{o[0], o[5], o[2], o[7], o[4], o[1], o[6], o[3]}
	m[6] = []complex128{o[0], o[6], o[4], o[2], o[0], o[6], o[4], o[2]}
	m[7] = []complex128{o[0], o[7], o[6], o[5], o[4], o[3], o[2], o[1]}
	return m.Mul(complex(1/math.Sqrt(8), 0))
}
