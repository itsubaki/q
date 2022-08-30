package gate

import (
	"fmt"
	"math"
	"math/cmplx"
	"strconv"

	"github.com/itsubaki/q/math/matrix"
	"github.com/itsubaki/q/math/number"
)

func New(v ...[]complex128) matrix.Matrix {
	g := make(matrix.Matrix, len(v))
	copy(g, v)
	return g
}

func Empty(n ...int) []matrix.Matrix {
	if len(n) < 1 {
		return make([]matrix.Matrix, 0)
	}

	return make([]matrix.Matrix, n[0])
}

// Theta returns 2 * pi / 2**k
func Theta(k int) float64 {
	return 2 * math.Pi / math.Pow(2, float64(k))
}

func U(theta, phi, lambda float64) matrix.Matrix {
	v := complex(theta/2, 0)
	return matrix.Matrix{
		[]complex128{cmplx.Cos(v), -1 * cmplx.Exp(complex(0, lambda)) * cmplx.Sin(v)},
		[]complex128{cmplx.Exp(complex(0, phi)) * cmplx.Sin(v), cmplx.Exp(complex(0, (phi+lambda))) * cmplx.Cos(v)},
	}
}

func I(n ...int) matrix.Matrix {
	return matrix.TensorProductN(matrix.Matrix{
		[]complex128{1, 0},
		[]complex128{0, 1},
	}, n...)
}

func X(n ...int) matrix.Matrix {
	return matrix.TensorProductN(matrix.Matrix{
		[]complex128{0, 1},
		[]complex128{1, 0},
	}, n...)
}

func Y(n ...int) matrix.Matrix {
	return matrix.TensorProductN(matrix.Matrix{
		[]complex128{0, -1i},
		[]complex128{1i, 0},
	}, n...)
}

func Z(n ...int) matrix.Matrix {
	return matrix.TensorProductN(matrix.Matrix{
		[]complex128{1, 0},
		[]complex128{0, -1},
	}, n...)
}

func H(n ...int) matrix.Matrix {
	v := complex(1/math.Sqrt2, 0)
	return matrix.TensorProductN(matrix.Matrix{
		[]complex128{v, v},
		[]complex128{v, -1 * v},
	}, n...)
}

func S(n ...int) matrix.Matrix {
	return matrix.TensorProductN(matrix.Matrix{
		[]complex128{1, 0},
		[]complex128{0, 1i},
	}, n...)
}

func T(n ...int) matrix.Matrix {
	v := cmplx.Exp(1i * math.Pi / 4)
	return matrix.TensorProductN(matrix.Matrix{
		[]complex128{1, 0},
		[]complex128{0, v},
	}, n...)
}

func R(theta float64) matrix.Matrix {
	e := cmplx.Exp(complex(0, theta))
	return matrix.Matrix{
		[]complex128{1, 0},
		[]complex128{0, e},
	}
}

func RX(theta float64) matrix.Matrix {
	v := complex(theta/2, 0)
	return matrix.Matrix{
		[]complex128{cmplx.Cos(v), -1i * cmplx.Sin(v)},
		[]complex128{-1i * cmplx.Sin(v), cmplx.Cos(v)},
	}
}

func RY(theta float64) matrix.Matrix {
	v := complex(theta/2, 0)
	return matrix.Matrix{
		[]complex128{cmplx.Cos(v), -1 * cmplx.Sin(v)},
		[]complex128{cmplx.Sin(v), cmplx.Cos(v)},
	}
}

func RZ(theta float64) matrix.Matrix {
	v := complex(0, theta/2)
	return matrix.Matrix{
		[]complex128{cmplx.Exp(-1 * v), 0},
		[]complex128{0, cmplx.Exp(v)},
	}
}

func Controlled(u matrix.Matrix, n int, c []int, t int) matrix.Matrix {
	g := I([]int{n}...)
	f := fmt.Sprintf("%s%s%s", "%0", strconv.Itoa(n), "s")

	for i := range g {
		row := []rune(fmt.Sprintf(f, strconv.FormatInt(int64(i), 2)))

		found := false
		for _, j := range c {
			if row[j] == '0' {
				found = true
				break
			}
		}

		if found {
			continue
		}

		for j := range g[i] {
			col := []rune(fmt.Sprintf(f, strconv.FormatInt(int64(j), 2)))

			found := false
			for _, k := range c {
				if col[k] == '0' {
					found = true
					break
				}
			}

			if found {
				continue
			}

			diff := false
			for i := range row {
				if i == t {
					continue
				}

				if row[i] != col[i] {
					diff = true
					break
				}
			}

			if diff {
				continue
			}

			r := number.Must(strconv.Atoi(string(row[t])))
			c := number.Must(strconv.Atoi(string(col[t])))

			g[j][i] = u[c][r]
		}
	}

	return g
}

func C(u matrix.Matrix, n int, c int, t int) matrix.Matrix {
	return Controlled(u, n, []int{c}, t)
}

func ControlledNot(n int, c []int, t int) matrix.Matrix {
	m := I([]int{n}...)
	d, _ := m.Dimension()
	f := fmt.Sprintf("%s%s%s", "%0", strconv.Itoa(n), "s")

	index := make([]int64, 0, d)
	for i := 0; i < d; i++ {
		bits := []rune(fmt.Sprintf(f, strconv.FormatInt(int64(i), 2)))

		apply := true
		for _, j := range c {
			if bits[j] == '0' {
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

		v := number.Must(strconv.ParseInt(string(bits), 2, 0))
		index = append(index, v)
	}

	g := make(matrix.Matrix, d)
	for i, ii := range index {
		g[ii] = m[i]
	}

	return g
}

func CNOT(n, c, t int) matrix.Matrix {
	return ControlledNot(n, []int{c}, t)
}

func CCNOT(n, c0, c1, t int) matrix.Matrix {
	return ControlledNot(n, []int{c0, c1}, t)
}

func Toffoli(n, c0, c1, t int) matrix.Matrix {
	return CCNOT(n, c0, c1, t)
}

func ControlledZ(n int, c []int, t int) matrix.Matrix {
	g := I([]int{n}...)
	d, _ := g.Dimension()
	f := fmt.Sprintf("%s%s%s", "%0", strconv.Itoa(n), "s")

	for i := 0; i < d; i++ {
		bits := []rune(fmt.Sprintf(f, strconv.FormatInt(int64(i), 2)))

		apply := true
		for _, j := range c {
			if bits[j] == '0' {
				apply = false
				break
			}
		}

		if apply && bits[t] == '1' {
			g[i][i] = -1 * g[i][i]
		}
	}

	return g
}

func CZ(n, c, t int) matrix.Matrix {
	return ControlledZ(n, []int{c}, t)
}

func ControlledS(n int, c []int, t int) matrix.Matrix {
	g := I([]int{n}...)
	d, _ := g.Dimension()
	f := fmt.Sprintf("%s%s%s", "%0", strconv.Itoa(n), "s")

	for i := 0; i < d; i++ {
		bits := []rune(fmt.Sprintf(f, strconv.FormatInt(int64(i), 2)))

		apply := true
		for _, j := range c {
			if bits[j] == '0' {
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

func CS(n, c, t int) matrix.Matrix {
	return ControlledS(n, []int{c}, t)
}

func ControlledR(theta float64, n int, c []int, t int) matrix.Matrix {
	g := I([]int{n}...)
	d, _ := g.Dimension()
	f := fmt.Sprintf("%s%s%s", "%0", strconv.Itoa(n), "s")

	// exp(i * theta)
	e := cmplx.Exp(complex(0, theta))

	for i := 0; i < d; i++ {
		bits := []rune(fmt.Sprintf(f, strconv.FormatInt(int64(i), 2)))

		// Apply R(k)
		apply := true
		for _, j := range c {
			if bits[j] == '0' {
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

func CR(theta float64, n, c, t int) matrix.Matrix {
	return ControlledR(theta, n, []int{c}, t)
}

// Swap returns unitary matrix of Swap operation.
func Swap(n, c, t int) matrix.Matrix {
	return matrix.Apply(
		CNOT(n, c, t),
		CNOT(n, t, c),
		CNOT(n, c, t),
	)
}

// Fredkin returns unitary matrix of Controlled-Swap operation.
func Fredkin(n, c, t0, t1 int) matrix.Matrix {
	return matrix.Apply(
		CNOT(n, t0, t1),
		CCNOT(n, c, t1, t0),
		CNOT(n, t0, t1),
	)
}

// QFT returns unitary matrix of Quantum Fourier Transform operation.
func QFT(n int) matrix.Matrix {
	g := I(n)

	for i := 0; i < n; i++ {
		h := make([]matrix.Matrix, 0)
		for j := 0; j < n; j++ {
			if i == j {
				h = append(h, H())
				continue
			}
			h = append(h, I())
		}
		g = g.Apply(matrix.TensorProduct(h...))

		k := 2
		for j := i + 1; j < n; j++ {
			g = g.Apply(CR(Theta(k), n, j, i))
			k++
		}
	}

	return g
}

// ControlledModExp2 returns unitary matrix of controlled modular exponentiation operation. |j>|k> -> |j>|a**(2**j) * k mod N>.
// len(t) must be larger than log2(N).
func ControlledModExp2(n, a, j, N, c int, t []int) matrix.Matrix {
	m := I([]int{n}...)
	d, _ := m.Dimension()

	r0len, r1len := n-len(t), len(t)
	a2jmodN := number.ModExp2(a, j, N)

	bf := fmt.Sprintf("%s%s%s", "%0", strconv.Itoa(n), "s")
	tf := fmt.Sprintf("%s%s%s", "%0", strconv.Itoa(r1len), "s")

	index := make([]int64, 0, d)
	for i := 0; i < d; i++ {
		bits := []rune(fmt.Sprintf(bf, strconv.FormatInt(int64(i), 2)))
		if bits[c] == '0' {
			index = append(index, int64(i))
			continue
		}

		k := number.Must(strconv.ParseInt(string(bits[r0len:]), 2, 0))
		if k > int64(N-1) {
			index = append(index, int64(i))
			continue
		}

		a2jkmodN := (int64(a2jmodN) * k) % int64(N)
		a2jkmodNs := fmt.Sprintf(tf, strconv.FormatInt(a2jkmodN, 2))
		newbits := append(bits[:r0len], []rune(a2jkmodNs)...)

		v := number.Must(strconv.ParseInt(string(newbits), 2, 0))
		index = append(index, v)
	}

	g := make(matrix.Matrix, d)
	for i, ii := range index {
		g[ii] = m[i]
	}

	return g
}
