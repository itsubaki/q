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
	g := make(matrix.Matrix, len(v))
	for i := 0; i < len(v); i++ {
		g[i] = v[i]
	}

	return g
}

func Empty(n ...int) []matrix.Matrix {
	if len(n) < 1 {
		return make([]matrix.Matrix, 0)
	}

	return make([]matrix.Matrix, n[0])
}

func U(alpha, beta, gamma, delta float64) matrix.Matrix {
	return matrix.Apply(
		RZ(beta),
		RY(gamma),
		RZ(delta),
	).Mul(cmplx.Exp(complex(0, alpha)))
}

func RX(theta float64) matrix.Matrix {
	v := complex(theta/2, 0)

	g := make(matrix.Matrix, 2)
	g[0] = []complex128{cmplx.Cos(v), complex(0, -1) * cmplx.Sin(v)}
	g[1] = []complex128{complex(0, -1) * cmplx.Sin(v), cmplx.Cos(v)}
	return g
}

func RY(theta float64) matrix.Matrix {
	v := complex(theta/2, 0)

	g := make(matrix.Matrix, 2)
	g[0] = []complex128{cmplx.Cos(v), -1 * cmplx.Sin(v)}
	g[1] = []complex128{cmplx.Sin(v), cmplx.Cos(v)}
	return g
}

func RZ(theta float64) matrix.Matrix {
	v := complex(0, 1*theta/2)

	g := make(matrix.Matrix, 2)
	g[0] = []complex128{cmplx.Exp(cmplx.Conj(v)), 0}
	g[1] = []complex128{0, cmplx.Exp(v)}
	return g
}

func R(k int) matrix.Matrix {
	p := 2 * math.Pi / math.Pow(2, float64(k))
	e := cmplx.Exp(complex(0, p))

	g := make(matrix.Matrix, 2)
	g[0] = []complex128{1, 0}
	g[1] = []complex128{0, e}
	return g
}

func I(bit ...int) matrix.Matrix {
	g := make(matrix.Matrix, 2)
	g[0] = []complex128{1, 0}
	g[1] = []complex128{0, 1}
	return matrix.TensorProductN(g, bit...)
}

func X(bit ...int) matrix.Matrix {
	g := make(matrix.Matrix, 2)
	g[0] = []complex128{0, 1}
	g[1] = []complex128{1, 0}
	return matrix.TensorProductN(g, bit...)
}

func Y(bit ...int) matrix.Matrix {
	g := make(matrix.Matrix, 2)
	g[0] = []complex128{0, -1i}
	g[1] = []complex128{1i, 0}
	return matrix.TensorProductN(g, bit...)
}

func Z(bit ...int) matrix.Matrix {
	g := make(matrix.Matrix, 2)
	g[0] = []complex128{1, 0}
	g[1] = []complex128{0, -1}
	return matrix.TensorProductN(g, bit...)
}

func H(bit ...int) matrix.Matrix {
	v := complex(1/math.Sqrt2, 0)

	g := make(matrix.Matrix, 2)
	g[0] = []complex128{v, v}
	g[1] = []complex128{v, -1 * v}
	return matrix.TensorProductN(g, bit...)
}

func S(bit ...int) matrix.Matrix {
	g := make(matrix.Matrix, 2)
	g[0] = []complex128{1, 0}
	g[1] = []complex128{0, 1i}
	return matrix.TensorProductN(g, bit...)
}

func T(bit ...int) matrix.Matrix {
	v := cmplx.Exp(complex(0, 1) * math.Pi / 4)

	g := make(matrix.Matrix, 2)
	g[0] = []complex128{1, 0}
	g[1] = []complex128{0, v}
	return matrix.TensorProductN(g, bit...)
}

func ControlledNot(bit int, c []int, t int) matrix.Matrix {
	m := I([]int{bit}...)
	d, _ := m.Dimension()
	f := fmt.Sprintf("%s%s%s", "%0", strconv.Itoa(bit), "s")

	index := make([]int64, 0)
	for i := 0; i < d; i++ {
		bits := []rune(fmt.Sprintf(f, strconv.FormatInt(int64(i), 2)))

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

	g := make(matrix.Matrix, d)
	for i, ii := range index {
		g[i] = m[ii]
	}

	return g
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
	g := I([]int{bit}...)
	d, _ := g.Dimension()
	f := fmt.Sprintf("%s%s%s", "%0", strconv.Itoa(bit), "s")

	for i := 0; i < d; i++ {
		bits := []rune(fmt.Sprintf(f, strconv.FormatInt(int64(i), 2)))

		// Apply Z
		apply := true
		for i := range c {
			if bits[c[i]] == '0' {
				apply = false
				break
			}
		}

		if apply && bits[t] == '1' {
			g[i][i] = complex(-1, 0) * g[i][i]
		}
	}

	return g
}

func CZ(bit, c, t int) matrix.Matrix {
	return ControlledZ(bit, []int{c}, t)
}

func ControlledS(bit int, c []int, t int) matrix.Matrix {
	g := I([]int{bit}...)
	d, _ := g.Dimension()
	f := fmt.Sprintf("%s%s%s", "%0", strconv.Itoa(bit), "s")

	for i := 0; i < d; i++ {
		bits := []rune(fmt.Sprintf(f, strconv.FormatInt(int64(i), 2)))

		// Apply S
		apply := true
		for i := range c {
			if bits[c[i]] == '0' {
				apply = false
				break
			}
		}

		if apply && bits[t] == '1' {
			g[i][i] = 1i * g[i][i]
		}
	}

	return g
}

func CS(bit, c, t int) matrix.Matrix {
	return ControlledS(bit, []int{c}, t)
}

func ControlledR(bit int, c []int, t, k int) matrix.Matrix {
	g := I([]int{bit}...)
	d, _ := g.Dimension()
	f := fmt.Sprintf("%s%s%s", "%0", strconv.Itoa(bit), "s")

	p := 2 * math.Pi / math.Pow(2, float64(k))
	e := cmplx.Exp(complex(0, p))

	for i := 0; i < d; i++ {
		bits := []rune(fmt.Sprintf(f, strconv.FormatInt(int64(i), 2)))

		// Apply R(k)
		apply := true
		for i := range c {
			if bits[c[i]] == '0' {
				apply = false
				break
			}
		}

		if apply && bits[t] == '1' {
			g[i][i] = e * g[i][i]
		}
	}

	return g
}

func CR(bit, c, t, k int) matrix.Matrix {
	return ControlledR(bit, []int{c}, t, k)
}

func Swap(bit, c, t int) matrix.Matrix {
	return matrix.Apply(
		CNOT(bit, c, t),
		CNOT(bit, t, c),
		CNOT(bit, c, t),
	)
}

// Fredkin returns unitary matrix of Controlled-Swap operation.
func Fredkin(bit, c, t0, t1 int) matrix.Matrix {
	return matrix.Apply(
		CNOT(bit, t0, t1),
		CCNOT(bit, c, t1, t0),
		CNOT(bit, t0, t1),
	)
}

// QFT returns unitary matrix of Quantum Fourier Transform operation.
func QFT(bit int) matrix.Matrix {
	g := I(bit)

	for i := 0; i < bit; i++ {
		h := make([]matrix.Matrix, 0)
		for j := 0; j < bit; j++ {
			if i == j {
				h = append(h, H())
				continue
			}
			h = append(h, I())
		}
		g = g.Apply(matrix.TensorProduct(h...))

		k := 2
		for j := i + 1; j < bit; j++ {
			g = g.Apply(CR(bit, j, i, k))
			k++
		}
	}

	return g
}

// CModExp2 returns unitary matrix of controlled modular exponentiation operation. |j>|k> -> |j>|a**(2**j) * k mod N>
func CModExp2(bit, a, j, N, c int, t []int) matrix.Matrix {
	min := int(math.Log2(float64(N))) + 1
	if len(t) < min {
		panic(fmt.Sprintf("invalid parameter. len(target)=%v < log2(%d)=%v", len(t), N, min))
	}

	m := I([]int{bit}...)
	d, _ := m.Dimension()

	r0len := bit - len(t)
	r1len := len(t)
	a2jmodN := number.ModExp2(a, j, N)

	bf := fmt.Sprintf("%s%s%s", "%0", strconv.Itoa(bit), "s")
	tf := fmt.Sprintf("%s%s%s", "%0", strconv.Itoa(r1len), "s")

	index := make([]int64, 0)
	for i := 0; i < d; i++ {
		bits := []rune(fmt.Sprintf(bf, strconv.FormatInt(int64(i), 2)))

		if bits[c] == '1' {
			r0bits, r1bits := bits[:r0len], bits[r0len:]

			k, err := strconv.ParseInt(string(r1bits), 2, 0)
			if err != nil {
				panic(fmt.Sprintf("parse int: %v", err))
			}

			if int(k) < N {
				a2jkmodN := (a2jmodN * int(k)) % N
				a2jkmodNs := fmt.Sprintf(tf, strconv.FormatInt(int64(a2jkmodN), 2))
				bits = append(r0bits, []rune(a2jkmodNs)...)

				// fmt.Printf("%v: %v=%2v -> %2v=%s -> ", string(bits[:r0len]), string(bits[r0len:]), k, a2jkmodN, a2jkmodNs)
				// fmt.Println(string(bits))
			}
		}

		v, err := strconv.ParseInt(string(bits), 2, 0)
		if err != nil {
			panic(fmt.Sprintf("parse int: %v", err))
		}

		index = append(index, v)
	}

	g := make(matrix.Matrix, d)
	for i, ii := range index {
		g[i] = m[ii]
	}

	return g.Transpose()
}
