package matrix

// QRFunc is a function type that performs QR decomposition on a given matrix.
// It returns Q (orthonormal columns) and R (upper triangular) such that A = Q * R.
type QRFunc func(a *Matrix, eps ...float64) (q *Matrix, r *Matrix)

var (
	_ QRFunc = QR
	_ QRFunc = QRHH
)

// Schur performs the Schur decomposition of matrix a using iterative QR decomposition.
// It returns Q (unitary) and T (upper triangular) such that A = Q * T * Q^dagger.
func Schur(a *Matrix, qr QRFunc, iter int, eps ...float64) (q *Matrix, t *Matrix) {
	qt, ak := Identity(a.Rows), a.Clone()

	for range iter {
		qk, rk := qr(ak, eps...)
		ak = rk.MatMul(qk)
		qt = qt.MatMul(qk)

		if ak.IsUpperTriangular(eps...) {
			break
		}
	}

	return qt, ak
}
