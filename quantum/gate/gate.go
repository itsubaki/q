package gate

import (
	"math"
	"math/cmplx"

	"github.com/itsubaki/q/math/matrix"
	"github.com/itsubaki/q/math/number"
)

// Theta returns 2 * pi / 2**k
func Theta(k int) float64 {
	return 2 * math.Pi / math.Pow(2, float64(k))
}

// New returns a new gate.
func New(v ...[]complex128) matrix.Matrix {
	g := make(matrix.Matrix, len(v))
	copy(g, v)
	return g
}

// Empty returns a list of empty gate.
func Empty(n ...int) []matrix.Matrix {
	if len(n) < 1 {
		return make([]matrix.Matrix, 0)
	}

	return make([]matrix.Matrix, n[0])
}

// U returns a unitary gate.
func U(theta, phi, lambda float64) matrix.Matrix {
	v := complex(theta/2, 0)
	return matrix.Matrix{
		[]complex128{cmplx.Cos(v), -1 * cmplx.Exp(complex(0, lambda)) * cmplx.Sin(v)},
		[]complex128{cmplx.Exp(complex(0, phi)) * cmplx.Sin(v), cmplx.Exp(complex(0, (phi+lambda))) * cmplx.Cos(v)},
	}
}

// I returns a identity gate.
func I(n ...int) matrix.Matrix {
	return matrix.TensorProductN(matrix.Matrix{
		[]complex128{1, 0},
		[]complex128{0, 1},
	}, n...)
}

// X returns a Pauli-X gate.
func X(n ...int) matrix.Matrix {
	return matrix.TensorProductN(matrix.Matrix{
		[]complex128{0, 1},
		[]complex128{1, 0},
	}, n...)
}

// Y returns a Pauli-Y gate.
func Y(n ...int) matrix.Matrix {
	return matrix.TensorProductN(matrix.Matrix{
		[]complex128{0, -1i},
		[]complex128{1i, 0},
	}, n...)
}

// Z returns a Pauli-Z gate.
func Z(n ...int) matrix.Matrix {
	return matrix.TensorProductN(matrix.Matrix{
		[]complex128{1, 0},
		[]complex128{0, -1},
	}, n...)
}

// H returns a Hadamard gate.
func H(n ...int) matrix.Matrix {
	v := complex(1/math.Sqrt2, 0)
	return matrix.TensorProductN(matrix.Matrix{
		[]complex128{v, v},
		[]complex128{v, -1 * v},
	}, n...)
}

// S returns a S gate.
func S(n ...int) matrix.Matrix {
	return matrix.TensorProductN(matrix.Matrix{
		[]complex128{1, 0},
		[]complex128{0, 1i},
	}, n...)
}

// T returns a T gate.
func T(n ...int) matrix.Matrix {
	v := cmplx.Exp(1i * math.Pi / 4)
	return matrix.TensorProductN(matrix.Matrix{
		[]complex128{1, 0},
		[]complex128{0, v},
	}, n...)
}

// R returns a rotation gate.
// R(Theta(k)) = [[1, 0], [0, exp(2 * pi * i / 2**k)]].
func R(theta float64) matrix.Matrix {
	e := cmplx.Exp(complex(0, theta))
	return matrix.Matrix{
		[]complex128{1, 0},
		[]complex128{0, e},
	}
}

// RX returns a rotation gate around the X axis.
func RX(theta float64) matrix.Matrix {
	v := complex(theta/2, 0)
	return matrix.Matrix{
		[]complex128{cmplx.Cos(v), -1i * cmplx.Sin(v)},
		[]complex128{-1i * cmplx.Sin(v), cmplx.Cos(v)},
	}
}

// RY returns a rotation gate around the Y axis.
func RY(theta float64) matrix.Matrix {
	v := complex(theta/2, 0)
	return matrix.Matrix{
		[]complex128{cmplx.Cos(v), -1 * cmplx.Sin(v)},
		[]complex128{cmplx.Sin(v), cmplx.Cos(v)},
	}
}

// RZ returns a rotation gate around the Z axis.
func RZ(theta float64) matrix.Matrix {
	v := complex(0, theta/2)
	return matrix.Matrix{
		[]complex128{cmplx.Exp(-1 * v), 0},
		[]complex128{0, cmplx.Exp(v)},
	}
}

// Controlled returns a controlled-u gate.
// u is a (2x2) unitary matrix and returns a (2**n x 2**n) matrix.
func Controlled(u matrix.Matrix, n int, c []int, t int) matrix.Matrix {
	var mask int
	for _, bit := range c {
		mask |= (1 << (n - 1 - bit))
	}

	s := (1 << n)
	g := I(n)
	for i := 0; i < s; i++ {
		if (i & mask) != mask {
			continue
		}

		for j := 0; j < s; j++ {
			if (j & mask) != mask {
				continue
			}

			// modify only the target qubit
			if (i & ^(1 << (n - 1 - t))) != (j & ^(1 << (n - 1 - t))) {
				continue
			}

			c := (j >> (n - 1 - t)) & 1
			r := (i >> (n - 1 - t)) & 1
			g[j][i] = u[c][r]
		}
	}

	return g
}

// C returns a controlled-u gate.
func C(u matrix.Matrix, n int, c int, t int) matrix.Matrix {
	return Controlled(u, n, []int{c}, t)
}

// ControlledNot returns a controlled-not gate.
func ControlledNot(n int, c []int, t int) matrix.Matrix {
	var mask int
	for _, bit := range c {
		mask |= (1 << (n - 1 - bit))
	}

	s := 1 << n
	perm := make([]int, s)
	for i := range s {
		perm[i] = i
		if (i & mask) == mask {
			perm[i] = i ^ (1 << (n - 1 - t))
		}
	}

	g, id := make(matrix.Matrix, s), I(n)
	for i, j := range perm {
		g[j] = id[i]
	}

	return g
}

// CNOT returns a controlled-not gate.
func CNOT(n, c, t int) matrix.Matrix {
	return ControlledNot(n, []int{c}, t)
}

// CCNOT returns a controlled-controlled-not gate.
func CCNOT(n, c0, c1, t int) matrix.Matrix {
	return ControlledNot(n, []int{c0, c1}, t)
}

// Toffoli returns a toffoli gate.
func Toffoli(n, c0, c1, t int) matrix.Matrix {
	return CCNOT(n, c0, c1, t)
}

// ControlledZ returns a controlled-z gate.
func ControlledZ(n int, c []int, t int) matrix.Matrix {
	var mask int
	for _, bit := range c {
		mask |= (1 << (n - 1 - bit))
	}

	g := I(n)
	for i := 0; i < (1 << n); i++ {
		if (i&mask) == mask && (i&(1<<(n-1-t))) != 0 {
			g[i][i] = -1 * g[i][i]
		}
	}

	return g
}

// CZ returns a controlled-z gate.
func CZ(n, c, t int) matrix.Matrix {
	return ControlledZ(n, []int{c}, t)
}

// ControlledS returns a controlled-s gate.
func ControlledS(n int, c []int, t int) matrix.Matrix {
	var mask int
	for _, bit := range c {
		mask |= (1 << (n - 1 - bit))
	}

	g := I(n)
	for i := 0; i < (1 << n); i++ {
		if (i&mask) == mask && (i&(1<<(n-1-t))) != 0 {
			g[i][i] = 1i * g[i][i]
		}
	}

	return g
}

// CS returns a controlled-s gate.
func CS(n, c, t int) matrix.Matrix {
	return ControlledS(n, []int{c}, t)
}

// ControlledR returns a controlled-r gate.
func ControlledR(theta float64, n int, c []int, t int) matrix.Matrix {
	// exp(i * theta)
	e := cmplx.Exp(complex(0, theta))

	var mask int
	for _, bit := range c {
		mask |= (1 << (n - 1 - bit))
	}

	g := I(n)
	for i := 0; i < (1 << n); i++ {
		if (i&mask) == mask && (i&(1<<(n-1-t))) != 0 {
			g[i][i] = e * g[i][i]
		}
	}

	return g
}

// CR returns a controlled-r gate.
func CR(theta float64, n, c, t int) matrix.Matrix {
	return ControlledR(theta, n, []int{c}, t)
}

// Swap returns a swap gate.
func Swap(n, c, t int) matrix.Matrix {
	return matrix.Apply(
		CNOT(n, c, t),
		CNOT(n, t, c),
		CNOT(n, c, t),
	)
}

// Fredkin returns a fredkin gate.
func Fredkin(n, c, t0, t1 int) matrix.Matrix {
	return matrix.Apply(
		CNOT(n, t0, t1),
		CCNOT(n, c, t1, t0),
		CNOT(n, t0, t1),
	)
}

// QFT returns a gate of Quantum Fourier Transform operation.
func QFT(n int) matrix.Matrix {
	g := I(n)

	for i := range n {
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

// ControlledModExp2 returns gate of controlled modular exponentiation operation.
// |j>|k> -> |j>|a**(2**j) * k mod N>.
// len(t) must be larger than log2(N).
func ControlledModExp2(n, a, j, N, c int, t []int) matrix.Matrix {
	m := I(n)
	r1len := len(t)
	a2jmodN := number.ModExp2(a, j, N)

	d, _ := m.Dimension()
	idx := make([]int, d)
	for i := range d {
		if (i>>(n-1-c))&1 == 0 {
			// control bit is 0, then do nothing
			idx[i] = i
			continue
		}

		// r1len bits of i
		mask := (1 << r1len) - 1
		k := i & mask
		if k > N-1 {
			idx[i] = i
			continue
		}

		// r0len bits of i + a2jkmodN bits
		a2jkmodN := a2jmodN * k % N
		idx[i] = (i >> r1len << r1len) | a2jkmodN
	}

	g := make(matrix.Matrix, d)
	for i, j := range idx {
		g[j] = m[i]
	}

	return g
}

// TensorProduct returns the tensor product of 'u' at specified indices over 'n' qubits.
func TensorProduct(u matrix.Matrix, n int, index []int) matrix.Matrix {
	idx := make(map[int]bool)
	for _, i := range index {
		idx[i] = true
	}

	g := I()
	if _, ok := idx[0]; ok {
		g = u
	}

	for i := 1; i < n; i++ {
		if _, ok := idx[i]; ok {
			g = g.TensorProduct(u)
			continue
		}

		g = g.TensorProduct(I())
	}

	return g
}
