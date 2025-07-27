package gate

import (
	"math"
	"math/cmplx"

	"github.com/itsubaki/q/math/matrix"
)

// Theta returns 2 * pi / 2**k
func Theta(k int) float64 {
	return 2 * math.Pi / math.Pow(2, float64(k))
}

// New returns a new gate.
func New(v ...[]complex128) *matrix.Matrix {
	return matrix.New(v...)
}

// U returns a unitary gate.
func U(theta, phi, lambda float64) *matrix.Matrix {
	v := complex(theta/2, 0)
	sin, cos := cmplx.Sin(v), cmplx.Cos(v)
	return matrix.New(
		[]complex128{cos, -1 * sin * cmplx.Exp(complex(0, lambda))},
		[]complex128{sin * cmplx.Exp(complex(0, phi)), cos * cmplx.Exp(complex(0, (phi+lambda)))},
	)
}

// I returns an identity gate.
func I(n ...int) *matrix.Matrix {
	return matrix.TensorProductN(matrix.New(
		[]complex128{1, 0},
		[]complex128{0, 1},
	), n...)
}

// X returns a Pauli-X gate.
func X(n ...int) *matrix.Matrix {
	return matrix.TensorProductN(matrix.New(
		[]complex128{0, 1},
		[]complex128{1, 0},
	), n...)
}

// Y returns a Pauli-Y gate.
func Y(n ...int) *matrix.Matrix {
	return matrix.TensorProductN(matrix.New(
		[]complex128{0, -1i},
		[]complex128{1i, 0},
	), n...)
}

// Z returns a Pauli-Z gate.
func Z(n ...int) *matrix.Matrix {
	return matrix.TensorProductN(matrix.New(
		[]complex128{1, 0},
		[]complex128{0, -1},
	), n...)
}

// H returns a Hadamard gate.
func H(n ...int) *matrix.Matrix {
	v := complex(1/math.Sqrt2, 0)
	return matrix.TensorProductN(matrix.New(
		[]complex128{v, v},
		[]complex128{v, -1 * v},
	), n...)
}

// S returns an S gate.
func S(n ...int) *matrix.Matrix {
	return matrix.TensorProductN(matrix.New(
		[]complex128{1, 0},
		[]complex128{0, 1i},
	), n...)
}

// T returns a T gate.
func T(n ...int) *matrix.Matrix {
	return matrix.TensorProductN(matrix.New(
		[]complex128{1, 0},
		[]complex128{0, cmplx.Exp(1i * math.Pi / 4)},
	), n...)
}

// R returns a rotation gate.
// R(Theta(k)) = [[1, 0], [0, exp(2 * pi * i / 2**k)]].
func R(theta float64) *matrix.Matrix {
	return matrix.New(
		[]complex128{1, 0},
		[]complex128{0, cmplx.Exp(complex(0, theta))},
	)
}

// RX returns a rotation gate around the X axis.
func RX(theta float64) *matrix.Matrix {
	v := complex(theta/2, 0)
	sin, cos := cmplx.Sin(v), cmplx.Cos(v)
	return matrix.New(
		[]complex128{cos, -1i * sin},
		[]complex128{-1i * sin, cos},
	)
}

// RY returns a rotation gate around the Y axis.
func RY(theta float64) *matrix.Matrix {
	v := complex(theta/2, 0)
	sin, cos := cmplx.Sin(v), cmplx.Cos(v)
	return matrix.New(
		[]complex128{cos, -1 * sin},
		[]complex128{sin, cos},
	)
}

// RZ returns a rotation gate around the Z axis.
func RZ(theta float64) *matrix.Matrix {
	v := complex(0, theta/2)
	return matrix.New(
		[]complex128{cmplx.Exp(-1 * v), 0},
		[]complex128{0, cmplx.Exp(v)},
	)
}

// C returns a controlled-u gate.
func C(u *matrix.Matrix, n int, c int, t int) *matrix.Matrix {
	return Controlled(u, n, []int{c}, t)
}

// CNOT returns a controlled-not gate.
func CNOT(n, c, t int) *matrix.Matrix {
	return ControlledNot(n, []int{c}, t)
}

// CCNOT returns a controlled-controlled-not gate.
func CCNOT(n, c0, c1, t int) *matrix.Matrix {
	return ControlledNot(n, []int{c0, c1}, t)
}

// CZ returns a controlled-z gate.
func CZ(n, c, t int) *matrix.Matrix {
	return ControlledZ(n, []int{c}, t)
}

// CS returns a controlled-s gate.
func CS(n, c, t int) *matrix.Matrix {
	return ControlledS(n, []int{c}, t)
}

// CR returns a controlled-r gate.
func CR(theta float64, n, c, t int) *matrix.Matrix {
	return ControlledR(theta, n, []int{c}, t)
}

// Controlled returns a controlled-u gate.
// u is a (2x2) unitary matrix and returns a (2**n x 2**n) matrix.
func Controlled(u *matrix.Matrix, n int, c []int, t int) *matrix.Matrix {
	var mask int
	for _, b := range c {
		mask |= (1 << (n - 1 - b))
	}

	s := (1 << n)
	g := I(n)
	for i := range s {
		if (i & mask) != mask {
			continue
		}

		for j := range s {
			if (j & mask) != mask {
				continue
			}

			// modify only the target qubit
			if (i & ^(1 << (n - 1 - t))) != (j & ^(1 << (n - 1 - t))) {
				continue
			}

			c := (j >> (n - 1 - t)) & 1
			r := (i >> (n - 1 - t)) & 1
			g.Set(i, j, u.At(c, r))
		}
	}

	return g
}

// ControlledNot returns a controlled-not gate.
func ControlledNot(n int, c []int, t int) *matrix.Matrix {
	var mask int
	for _, b := range c {
		mask |= (1 << (n - 1 - b))
	}

	s := 1 << n
	perm := make([]int, s)
	for i := range s {
		perm[i] = i
		if (i & mask) == mask {
			perm[i] = i ^ (1 << (n - 1 - t))
		}
	}

	data, id := make([][]complex128, s), I(n)
	for i, j := range perm {
		data[j] = id.Row(i)
	}

	return matrix.New(data...)
}

// ControlledZ returns a controlled-z gate.
func ControlledZ(n int, c []int, t int) *matrix.Matrix {
	var mask int
	for _, b := range c {
		mask |= (1 << (n - 1 - b))
	}

	g := I(n)
	for i := range 1 << n {
		if (i&mask) == mask && (i&(1<<(n-1-t))) != 0 {
			g.MulAt(i, i, -1)
		}
	}

	return g
}

// ControlledS returns a controlled-s gate.
func ControlledS(n int, c []int, t int) *matrix.Matrix {
	var mask int
	for _, b := range c {
		mask |= (1 << (n - 1 - b))
	}

	g := I(n)
	for i := range 1 << n {
		if (i&mask) == mask && (i&(1<<(n-1-t))) != 0 {
			g.MulAt(i, i, 1i)
		}
	}

	return g
}

// ControlledR returns a controlled-r gate.
func ControlledR(theta float64, n int, c []int, t int) *matrix.Matrix {
	// exp(i * theta)
	e := cmplx.Exp(complex(0, theta))

	var mask int
	for _, b := range c {
		mask |= (1 << (n - 1 - b))
	}

	g := I(n)
	for i := range 1 << n {
		if (i&mask) == mask && (i&(1<<(n-1-t))) != 0 {
			g.MulAt(i, i, e)
		}
	}

	return g
}

// Swap returns a swap gate.
func Swap(n, c, t int) *matrix.Matrix {
	return matrix.Apply(
		CNOT(n, c, t),
		CNOT(n, t, c),
		CNOT(n, c, t),
	)
}

// QFT returns a gate of Quantum Fourier Transform operation.
func QFT(n int) *matrix.Matrix {
	g := I(n)

	for i := range n {
		h := make([]*matrix.Matrix, n)
		for j := range n {
			if i == j {
				h[j] = H()
				continue
			}

			h[j] = I()
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

// TensorProduct returns the tensor product of 'u' at specified indices over 'n' qubits.
func TensorProduct(u *matrix.Matrix, n int, idx []int) *matrix.Matrix {
	target := make(map[int]bool)
	for _, i := range idx {
		target[i] = true
	}

	g := I()
	if _, ok := target[0]; ok {
		g = u
	}

	for i := 1; i < n; i++ {
		if _, ok := target[i]; ok {
			g = g.TensorProduct(u)
			continue
		}

		g = g.TensorProduct(I())
	}

	return g
}
