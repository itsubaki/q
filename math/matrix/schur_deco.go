package matrix

// Schur performs the Schur decomposition of matrix a.
// It returns Q (unitary) and T (upper triangular) such that A = Q * T * Q^dagger.
func Schur(a *Matrix, iter int, eps ...float64) (q *Matrix, t *Matrix) {
	qt, ak := Identity(a.Rows), a.Clone()

	for range iter {
		qk, rk := QR(ak)
		ak = rk.MatMul(qk)
		qt = qt.MatMul(qk)

		if ak.IsUpperTriangular(eps...) {
			break
		}
	}

	return qt, ak
}
