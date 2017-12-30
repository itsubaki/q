package q

import (
	"testing"

	"github.com/itsubaki/q/gate"
	"github.com/itsubaki/q/matrix"
	"github.com/itsubaki/q/qubit"
)

func TestGrover3qubit(t *testing.T) {
	x := matrix.TensorProduct(gate.X(), gate.I(3))
	oracle := x.Apply(gate.CNOT(4)).Apply(x)

	h4 := matrix.TensorProduct(gate.H(3), gate.H())
	x3 := matrix.TensorProduct(gate.X(3), gate.I())
	cz := matrix.TensorProduct(gate.CZ(3), gate.I())
	h3 := matrix.TensorProduct(gate.H(3), gate.I())
	amp := h4.Apply(x3).Apply(cz).Apply(x3).Apply(h3)

	q := qubit.TensorProduct(qubit.Zero(3), qubit.One())
	q.Apply(gate.H(4)).Apply(oracle).Apply(amp)

	// q3 is always |1>
	q3 := q.MeasureAt(3)
	if !q3.IsOne() {
		t.Error(q3)
	}

	p := q.Probability()
	for i, pp := range p {
		// |011>|1>
		if i == 7 {
			if pp-0.78125 > 1e-13 {
				t.Error(q.Probability())
			}
			continue
		}

		if pp-0.03125 > 1e-13 {
			t.Error(q.Probability())
		}
	}

}

func TestGrover2qubit(t *testing.T) {
	oracle := gate.CZ(2)

	h2 := gate.H(2)
	x2 := gate.X(2)
	amp := h2.Apply(x2).Apply(gate.CZ(2)).Apply(x2).Apply(h2)

	qc := h2.Apply(oracle).Apply(amp)
	q := qubit.Zero(2).Apply(qc)

	q.Measure()
	if q.Probability()[3]-1 > 1e-13 {
		t.Error(q.Probability())
	}
}

func TestQuantumTeleportation(t *testing.T) {
	g0 := matrix.TensorProduct(gate.H(), gate.I())
	g1 := gate.CNOT()
	bell := qubit.Zero(2).Apply(g0).Apply(g1)

	phi := qubit.New(1, 2)
	phi.TensorProduct(bell)

	g2 := matrix.TensorProduct(gate.CNOT(), gate.I())
	g3 := matrix.TensorProduct(gate.H(), gate.I(2))
	phi.Apply(g2).Apply(g3)

	mz := phi.MeasureAt(0)
	mx := phi.MeasureAt(1)

	if mz.IsOne() {
		z := matrix.TensorProduct(gate.I(2), gate.Z())
		phi.Apply(z)
	}

	if mx.IsOne() {
		x := matrix.TensorProduct(gate.I(2), gate.X())
		phi.Apply(x)
	}

	var test = []struct {
		zero int
		one  int
		zval qubit.Probability
		oval qubit.Probability
		eps  qubit.Probability
		mz   *qubit.Qubit
		mx   *qubit.Qubit
	}{
		{0, 1, 0.2, 0.8, 1e-13, qubit.Zero(), qubit.Zero()},
		{2, 3, 0.2, 0.8, 1e-13, qubit.Zero(), qubit.One()},
		{4, 5, 0.2, 0.8, 1e-13, qubit.One(), qubit.Zero()},
		{6, 7, 0.2, 0.8, 1e-13, qubit.One(), qubit.One()},
	}

	p := phi.Probability()
	for _, tt := range test {
		if p[tt.zero] == 0 {
			continue
		}

		if p[tt.zero]-tt.zval > tt.eps {
			t.Error(p)
		}
		if p[tt.one]-tt.oval > tt.eps {
			t.Error(p)
		}

		if !mz.Equals(tt.mz) {
			t.Error(p)
		}

		if !mx.Equals(tt.mx) {
			t.Error(p)
		}

		if qubit.Sum(p)-1 > tt.eps {
			t.Error(p)
		}
	}
}

func TestQuantumTeleportationPattern2(t *testing.T) {
	g0 := matrix.TensorProduct(gate.H(), gate.I())
	g1 := gate.CNOT()
	bell := qubit.Zero(2).Apply(g0).Apply(g1)

	phi := qubit.New(1, 2)
	phi.TensorProduct(bell)

	g2 := matrix.TensorProduct(gate.CNOT(), gate.I())
	g3 := matrix.TensorProduct(gate.H(), gate.I(2))
	g4 := matrix.TensorProduct(gate.I(), gate.CNOT())
	g5 := gate.CZc1t3()

	phi.Apply(g2).Apply(g3).Apply(g4).Apply(g5)

	mz := phi.MeasureAt(0)
	mx := phi.MeasureAt(1)

	var test = []struct {
		zero int
		one  int
		zval qubit.Probability
		oval qubit.Probability
		eps  qubit.Probability
		mz   *qubit.Qubit
		mx   *qubit.Qubit
	}{
		{0, 1, 0.2, 0.8, 1e-13, qubit.Zero(), qubit.Zero()},
		{2, 3, 0.2, 0.8, 1e-13, qubit.Zero(), qubit.One()},
		{4, 5, 0.2, 0.8, 1e-13, qubit.One(), qubit.Zero()},
		{6, 7, 0.2, 0.8, 1e-13, qubit.One(), qubit.One()},
	}

	p := phi.Probability()
	for _, tt := range test {
		if p[tt.zero] == 0 {
			continue
		}

		if p[tt.zero]-tt.zval > tt.eps {
			t.Error(p)
		}
		if p[tt.one]-tt.oval > tt.eps {
			t.Error(p)
		}

		if !mz.Equals(tt.mz) {
			t.Error(p)
		}

		if !mx.Equals(tt.mx) {
			t.Error(p)
		}

		if qubit.Sum(p)-1 > tt.eps {
			t.Error(p)
		}
	}
}
