package trace

import (
	"math"
)

var identity = NewM44()

type M44 struct {
	m [4][4]float64
}

func NewM44() *M44 {
	return &M44{[4][4]float64{{1.0, 0.0, 0.0, 0.0},
		{0.0, 1.0, 0.0, 0.0},
		{0.0, 0.0, 1.0, 0.0},
		{0.0, 0.0, 0.0, 1.0}}}
}

func (m *M44) inverse() *M44 {
	s := NewM44()
	t := *m

	// Forward elimination
	for i := 0; i < 3; i++ {
		pivot := i

		pivotsize := math.Abs(t.m[i][i])

		for j := i + 1; j < 4; j++ {
			tmp := math.Abs(t.m[j][i])

			if tmp > pivotsize {
				pivot = j
				pivotsize = tmp
			}
		}

		if pivotsize == 0 {
			// singular matrix
			return NewM44()
		}

		if pivot != i {
			for j := 0; j < 4; j++ {
				tmp := t.m[i][j]
				t.m[i][j] = t.m[pivot][j]
				t.m[pivot][j] = tmp

				tmp = s.m[i][j]
				s.m[i][j] = s.m[pivot][j]
				s.m[pivot][j] = tmp
			}
		}

		for j := i + 1; j < 4; j++ {
			f := t.m[j][i] / t.m[i][i]
			for k := 0; k < 4; k++ {
				t.m[j][k] = t.m[j][k] - f*t.m[i][k]
				s.m[j][k] = s.m[j][k] - f*s.m[i][k]
			}
		}
	}

	// Backward substitution
	for i := 3; i >= 0; i-- {
		f := t.m[i][i]
		if f == 0 {
			// singular matrix
			return NewM44()
		}

		for j := 0; j < 4; j++ {
			t.m[i][j] = t.m[i][j] / f
			s.m[i][j] = s.m[i][j] / f
		}

		for j := 0; j < i; j++ {
			f = t.m[j][i]

			for k := 0; k < 4; k++ {
				t.m[j][k] = t.m[j][k] - f*t.m[i][k]
				s.m[j][k] = s.m[j][k] - f*s.m[i][k]
			}
		}
	}

	return s
}

func (m *M44) transformPt(p *Pt) *Pt {
	x := p.X*m.m[0][0] + p.Y*m.m[1][0] + p.Z*m.m[2][0] + m.m[3][0]
	y := p.X*m.m[0][1] + p.Y*m.m[1][1] + p.Z*m.m[2][1] + m.m[3][1]
	z := p.X*m.m[0][2] + p.Y*m.m[1][2] + p.Z*m.m[2][2] + m.m[3][2]
	w := p.X*m.m[0][3] + p.Y*m.m[1][3] + p.Z*m.m[2][3] + m.m[3][3]

	pt := &Pt{x / w, y / w, z / w}
	return pt
}

func (m *M44) rotateV3(v *V3) *V3 {
	x := v.X*m.m[0][0] + v.Y*m.m[1][0] + v.Z*m.m[2][0]
	y := v.X*m.m[0][1] + v.Y*m.m[1][1] + v.Z*m.m[2][1]
	z := v.X*m.m[0][2] + v.Y*m.m[1][2] + v.Z*m.m[2][2]

	return &V3{x, y, z}
}

func (m *M44) Translate(v *V3) *M44 {
	m.m[3][0] = m.m[3][0] + v.X
	m.m[3][1] = m.m[3][1] + v.Y
	m.m[3][2] = m.m[3][2] + v.Z
	return m
}

func (m *M44) Scale(v *V3) *M44 {
	m.m[0][0] = m.m[0][0] * v.X
	m.m[1][1] = m.m[1][1] * v.Y
	m.m[2][2] = m.m[2][2] * v.Z
	return m
}
