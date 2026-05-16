package observable_test

import (
	"testing"

	"github.com/itsubaki/q/math/matrix"
	"github.com/itsubaki/q/quantum/gate"
	"github.com/itsubaki/q/quantum/observable"
	"github.com/itsubaki/q/quantum/qubit"
)

func TestPauli(t *testing.T) {
	cases := []struct {
		s string
	}{
		{"I"},
		{"X"},
		{"Y"},
		{"Z"},
		{"IX"},
		{"IY"},
		{"IZ"},
		{"XIZ"},
		{"YIZ"},
		{"ZIZ"},
	}

	for _, c := range cases {
		if !observable.Pauli(c.s).IsHermitian() {
			t.Errorf("observable.Pauli(%s) is not Hermitian", c.s)
		}
	}
}

func TestI(t *testing.T) {
	cases := []struct {
		n    []int
		want *matrix.Matrix
	}{
		{
			n:    []int{},
			want: gate.I(),
		},
		{
			n:    []int{1},
			want: gate.I(1),
		},
		{
			n:    []int{2},
			want: gate.I(2),
		},
	}

	for _, c := range cases {
		got := observable.I(c.n...)
		if !got.Equal(c.want) {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}

func TestX(t *testing.T) {
	cases := []struct {
		n    []int
		want *matrix.Matrix
	}{
		{
			n:    []int{},
			want: gate.X(),
		},
		{
			n:    []int{1},
			want: gate.X(1),
		},
		{
			n:    []int{2},
			want: gate.X(2),
		},
	}

	for _, c := range cases {
		got := observable.X(c.n...)
		if !got.Equal(c.want) {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}

func TestY(t *testing.T) {
	cases := []struct {
		n    []int
		want *matrix.Matrix
	}{
		{
			n:    []int{},
			want: gate.Y(),
		},
		{
			n:    []int{1},
			want: gate.Y(1),
		},
		{
			n:    []int{2},
			want: gate.Y(2),
		},
	}

	for _, c := range cases {
		got := observable.Y(c.n...)
		if !got.Equal(c.want) {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}

func TestZ(t *testing.T) {
	cases := []struct {
		n    []int
		want *matrix.Matrix
	}{
		{
			n:    []int{},
			want: gate.Z(),
		},
		{
			n:    []int{1},
			want: gate.Z(1),
		},
		{
			n:    []int{2},
			want: gate.Z(2),
		},
	}

	for _, c := range cases {
		got := observable.Z(c.n...)
		if !got.Equal(c.want) {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}

func TestProjector(t *testing.T) {
	cases := []struct {
		in   []*qubit.Qubit
		want *matrix.Matrix
	}{
		{
			in: []*qubit.Qubit{qubit.Zero()},
			want: gate.New(
				[]complex128{1, 0},
				[]complex128{0, 0},
			),
		},
		{
			in: []*qubit.Qubit{qubit.One()},
			want: gate.New(
				[]complex128{0, 0},
				[]complex128{0, 1},
			),
		},
		{
			in: []*qubit.Qubit{
				qubit.Zero(),
				qubit.Zero(),
			},
			want: gate.New(
				[]complex128{1, 0, 0, 0},
				[]complex128{0, 0, 0, 0},
				[]complex128{0, 0, 0, 0},
				[]complex128{0, 0, 0, 0},
			),
		},
	}

	for _, c := range cases {
		got := observable.Projector(c.in...)
		if !got.Equal(c.want) {
			t.Errorf("got=%v, want=%v", got, c.want)
		}
	}
}
