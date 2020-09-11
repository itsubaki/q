package gate

import (
	"fmt"
	"math"
	"math/cmplx"
	"strconv"

	"github.com/itsubaki/q/pkg/math/matrix"
	"github.com/itsubaki/q/pkg/math/number"
)

func New(v ...[]complex128) matrix.Matrix {
	out := make(matrix.Matrix, len(v))
	for i := 0; i < len(v); i++ {
		out[i] = v[i]
	}

	return out
}

func NewSlice(l ...int) []matrix.Matrix {
	if len(l) < 1 {
		return make([]matrix.Matrix, 0)
	}

	return make([]matrix.Matrix, l[0])
}

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
	p := 2 * math.Pi / math.Pow(2, float64(k))
	e := cmplx.Exp(complex(0, p))

	out := make(matrix.Matrix, 2)
	out[0] = []complex128{1, 0}
	out[1] = []complex128{0, e}
	return out
}

func I(bit ...int) matrix.Matrix {
	m := make(matrix.Matrix, 2)
	m[0] = []complex128{1, 0}
	m[1] = []complex128{0, 1}
	return matrix.TensorProductN(m, bit...)
}

func X(bit ...int) matrix.Matrix {
	m := make(matrix.Matrix, 2)
	m[0] = []complex128{0, 1}
	m[1] = []complex128{1, 0}
	return matrix.TensorProductN(m, bit...)
}

func Y(bit ...int) matrix.Matrix {
	m := make(matrix.Matrix, 2)
	m[0] = []complex128{0, -1i}
	m[1] = []complex128{1i, 0}
	return matrix.TensorProductN(m, bit...)
}

func Z(bit ...int) matrix.Matrix {
	m := make(matrix.Matrix, 2)
	m[0] = []complex128{1, 0}
	m[1] = []complex128{0, -1}
	return matrix.TensorProductN(m, bit...)
}

func H(bit ...int) matrix.Matrix {
	v := complex(1/math.Sqrt2, 0)

	m := make(matrix.Matrix, 2)
	m[0] = []complex128{v, v}
	m[1] = []complex128{v, -1 * v}
	return matrix.TensorProductN(m, bit...)
}

func S(bit ...int) matrix.Matrix {
	m := make(matrix.Matrix, 2)
	m[0] = []complex128{1, 0}
	m[1] = []complex128{0, 1i}
	return matrix.TensorProductN(m, bit...)
}

func T(bit ...int) matrix.Matrix {
	v := cmplx.Exp(complex(0, 1) * math.Pi / 4)

	m := make(matrix.Matrix, 2)
	m[0] = []complex128{1, 0}
	m[1] = []complex128{0, v}
	return matrix.TensorProductN(m, bit...)
}

func ControlledNot(bit int, c []int, t int) matrix.Matrix {
	mat := I([]int{bit}...)
	dim, _ := mat.Dimension()

	index := make([]int64, 0)
	f := "%0" + strconv.Itoa(bit) + "s"
	for i := 0; i < dim; i++ {
		s := fmt.Sprintf(f, strconv.FormatInt(int64(i), 2))
		bits := []rune(s)

		// Apply X
		apply := true
		for i := range c {
			if bits[c[i]] == '0' {
				apply = false
				break
			}
		}

		if apply {
			if bits[t] == '0' {
				bits[t] = '1'
			} else {
				bits[t] = '0'
			}
		}

		v, err := strconv.ParseInt(string(bits), 2, 0)
		if err != nil {
			panic(fmt.Sprintf("parse int: %v", err))
		}

		index = append(index, v)
	}

	out := make(matrix.Matrix, dim)
	for i, ii := range index {
		out[i] = mat[ii]
	}

	return out
}

func CNOT(bit, c, t int) matrix.Matrix {
	return ControlledNot(bit, []int{c}, t)
}

func CCNOT(bit, c0, c1, t int) matrix.Matrix {
	return ControlledNot(bit, []int{c0, c1}, t)
}

func Toffoli(bit, c0, c1, t int) matrix.Matrix {
	return CCNOT(bit, c0, c1, t)
}

func ControlledZ(bit int, c []int, t int) matrix.Matrix {
	out := I([]int{bit}...)
	dim, _ := out.Dimension()

	f := "%0" + strconv.Itoa(bit) + "s"
	for i := 0; i < dim; i++ {
		s := fmt.Sprintf(f, strconv.FormatInt(int64(i), 2))
		bits := []rune(s)

		// Apply Z
		apply := true
		for i := range c {
			if bits[c[i]] == '0' {
				apply = false
				break
			}
		}

		if apply && bits[t] == '1' {
			out[i][i] = complex(-1, 0) * out[i][i]
		}
	}

	return out
}

func CZ(bit, c, t int) matrix.Matrix {
	return ControlledZ(bit, []int{c}, t)
}

func ControlledS(bit int, c []int, t int) matrix.Matrix {
	out := I([]int{bit}...)
	dim, _ := out.Dimension()

	f := "%0" + strconv.Itoa(bit) + "s"
	for i := 0; i < dim; i++ {
		s := fmt.Sprintf(f, strconv.FormatInt(int64(i), 2))
		bits := []rune(s)

		// Apply S
		apply := true
		for i := range c {
			if bits[c[i]] == '0' {
				apply = false
				break
			}
		}

		if apply && bits[t] == '1' {
			out[i][i] = 1i * out[i][i]
		}
	}

	return out
}

func CS(bit, c, t int) matrix.Matrix {
	return ControlledS(bit, []int{c}, t)
}

func ControlledR(bit int, c []int, t, k int) matrix.Matrix {
	out := I([]int{bit}...)
	dim, _ := out.Dimension()

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
				break
			}
		}

		if apply && bits[t] == '1' {
			out[i][i] = e * out[i][i]
		}
	}

	return out
}

func CR(bit, c, t, k int) matrix.Matrix {
	return ControlledR(bit, []int{c}, t, k)
}

func Swap(bit, c, t int) matrix.Matrix {
	g0 := CNOT(bit, c, t)
	g1 := CNOT(bit, t, c)
	g2 := CNOT(bit, c, t)
	return g0.Apply(g1).Apply(g2)
}

// Fredkin returns unitary matrix of Controlled-Swap operation
func Fredkin(bit, c, t0, t1 int) matrix.Matrix {
	g0 := CNOT(bit, t0, t1)
	g1 := CCNOT(bit, c, t1, t0)
	g2 := CNOT(bit, t0, t1)
	return g0.Apply(g1).Apply(g2)
}

// QFT returns unitary matrix of Quantum Fourier Transform
func QFT(bit int) matrix.Matrix {
	out := I(bit)

	for i := 0; i < bit; i++ {
		h := make([]matrix.Matrix, 0)
		for j := 0; j < bit; j++ {
			if i == j {
				h = append(h, H())
				continue
			}
			h = append(h, I())
		}
		out = out.Apply(matrix.TensorProduct(h...))

		k := 2
		for j := i + 1; j < bit; j++ {
			out = out.Apply(CR(bit, j, i, k))
			k++
		}
	}

	return out
}

// CModExp returns unitary matrix of Controlled-U^(2^j)|k> -> |a^(2^j)*k mod N> operation
func CModExp(bit, N, a, j, control int, target []int) matrix.Matrix {
	min := int(math.Log2(float64(N))) + 1
	if bit < min {
		panic(fmt.Sprintf("invalid parameter. bit=%v < log2(%d)=%v", bit, N, min))
	}

	mat := I([]int{bit}...)
	dim, _ := mat.Dimension()
	r0len := bit - len(target)
	r1len := len(target)

	bf := "%0" + strconv.Itoa(bit) + "s"
	tf := "%0" + strconv.Itoa(r1len) + "s"

	index := make([]int64, 0)
	for i := 0; i < dim; i++ {
		s := fmt.Sprintf(bf, strconv.FormatInt(int64(i), 2))
		bits := []rune(s)

		if bits[control] == '1' {
			r0bits, r1bits := bits[:r0len], bits[r0len:]

			k, err := strconv.ParseInt(string(r1bits), 2, 0)
			if err != nil {
				panic(fmt.Sprintf("parse int: %v", err))
			}

			if k < int64(N) {
				a2jkmodN := (number.ModExp2(a, j, N) * k) % int64(N)
				a2jkmodNs := fmt.Sprintf(tf, strconv.FormatInt(a2jkmodN, 2))
				bits = append(r0bits, []rune(a2jkmodNs)...)
			}
		}

		v, err := strconv.ParseInt(string(bits), 2, 0)
		if err != nil {
			panic(fmt.Sprintf("parse int: %v", err))
		}

		index = append(index, v)
	}

	out := make(matrix.Matrix, dim)
	for i, ii := range index {
		out[i] = mat[ii]
	}

	return out
}
