package noise

import (
	"math"
	"math/rand"

	"github.com/terefang/gommons/pkg/xmath"
)

// Package perlin provides coherent noise function over 1, 2 or 3 dimensions
// This code is go adaptation based on C implementation that can be found here:
// http://git.gnome.org/browse/gegl/tree/operations/common/perlin/perlin.c
// (original copyright Ken Perlin)

// General constants
const (
	B  = 0x100
	N  = 0x1000
	BM = 0xff
)

// OldPerlin is the noise generator
type OldPerlin struct {
	p  [B + B + 2]int32
	g6 [B + B + 2][6]float64
	g5 [B + B + 2][5]float64
	g4 [B + B + 2][4]float64
	g3 [B + B + 2][3]float64
	g2 [B + B + 2][2]float64
	g1 [B + B + 2]float64
}

// NewOldPerlin creates new Perlin noise generator
// In what follows "alpha" is the weight when the sum is formed.
// Typically it is 2, As this approaches 1 the function is noisier.
// "beta" is the harmonic scaling/spacing, typically 2, n is the
// number of iterations and seed is the math.rand seed value to use
func NewOldPerlin(seed int64) *OldPerlin {
	return NewOldPerlinWithSource(rand.NewSource(seed))
}

// NewOldPerlinWithSource creates new Perlin noise generator initialized
// by the source of pseudo-random int64 values
func NewOldPerlinWithSource(source rand.Source) *OldPerlin {
	var p OldPerlin
	var i, j int32

	r := rand.New(source)

	for i = 0; i < B; i++ {
		p.p[i] = i
	}

	for ; i > 0; i-- {
		j = r.Int31() % B
		p.p[i], p.p[j] = p.p[j], p.p[i]
	}

	for i = 0; i < B+2; i++ {
		p.p[B+i] = p.p[i]
	}

	p.init1D(r)
	p.init2D(r)
	p.init3D(r)
	p.init4D(r)
	p.init5D(r)
	p.init6D(r)
	return &p
}

func (p *OldPerlin) init1D(r *rand.Rand) {
	for i := 0; i < B; i++ {
		var x float64
		for {
			x = (float64(r.Intn(B+B)) - float64(B)) / float64(B)
			magSq := x * x
			if magSq <= 1.0 && magSq > 0.0 {
				break
			}
		}
		p.g1[i] = x
	}

	for i := 0; i < B+2; i++ {
		p.g1[B+i] = p.g1[i]
	}
}

func (p *OldPerlin) init2D(r *rand.Rand) {
	for i := 0; i < B; i++ {
		var x, y float64
		for {
			x = (float64(r.Intn(B+B)) - float64(B)) / float64(B)
			y = (float64(r.Intn(B+B)) - float64(B)) / float64(B)
			magSq := x*x + y*y
			if magSq <= 1.0 && magSq > 0.0 {
				break
			}
		}
		mag := math.Sqrt(x*x + y*y)
		p.g2[i] = [2]float64{x / mag, y / mag}
	}

	for i := 0; i < B+2; i++ {
		p.g2[B+i] = p.g2[i]
	}
}

func (p *OldPerlin) init3D(r *rand.Rand) {
	for i := 0; i < B; i++ {
		var x, y, z float64
		for {
			x = (float64(r.Intn(B+B)) - float64(B)) / float64(B)
			y = (float64(r.Intn(B+B)) - float64(B)) / float64(B)
			z = (float64(r.Intn(B+B)) - float64(B)) / float64(B)
			magSq := x*x + y*y + z*z
			if magSq <= 1.0 && magSq > 0.0 {
				break
			}
		}
		mag := math.Sqrt(x*x + y*y + z*z)
		p.g3[i] = [3]float64{x / mag, y / mag, z / mag}
	}

	for i := 0; i < B+2; i++ {
		p.g3[B+i] = p.g3[i]
	}
}

func (p *OldPerlin) init4D(r *rand.Rand) {
	for i := 0; i < B; i++ {
		var x, y, z, w4 float64
		for {
			x = (float64(r.Intn(B+B)) - float64(B)) / float64(B)
			y = (float64(r.Intn(B+B)) - float64(B)) / float64(B)
			z = (float64(r.Intn(B+B)) - float64(B)) / float64(B)
			w4 = (float64(r.Intn(B+B)) - float64(B)) / float64(B)
			magSq := x*x + y*y + z*z + w4*w4
			if magSq <= 1.0 && magSq > 0.0 {
				break
			}
		}
		mag := math.Sqrt(x*x + y*y + z*z + w4*w4)
		p.g4[i] = [4]float64{x / mag, y / mag, z / mag, w4 / mag}
	}

	for i := 0; i < B+2; i++ {
		p.g4[B+i] = p.g4[i]
	}
}

func (p *OldPerlin) init5D(r *rand.Rand) {
	for i := 0; i < B; i++ {
		var x, y, z, w4, v5 float64
		for {
			x = (float64(r.Intn(B+B)) - float64(B)) / float64(B)
			y = (float64(r.Intn(B+B)) - float64(B)) / float64(B)
			z = (float64(r.Intn(B+B)) - float64(B)) / float64(B)
			w4 = (float64(r.Intn(B+B)) - float64(B)) / float64(B)
			v5 = (float64(r.Intn(B+B)) - float64(B)) / float64(B)
			magSq := x*x + y*y + z*z + w4*w4 + v5*v5
			if magSq <= 1.0 && magSq > 0.0 {
				break
			}
		}
		mag := math.Sqrt(x*x + y*y + z*z + w4*w4 + v5*v5)
		p.g5[i] = [5]float64{x / mag, y / mag, z / mag, w4 / mag, v5 / mag}
	}

	for i := 0; i < B+2; i++ {
		p.g5[B+i] = p.g5[i]
	}
}

func (p *OldPerlin) init6D(r *rand.Rand) {
	for i := 0; i < B; i++ {
		var x, y, z, w4, v5, u6 float64
		for {
			x = (float64(r.Intn(B+B)) - float64(B)) / float64(B)
			y = (float64(r.Intn(B+B)) - float64(B)) / float64(B)
			z = (float64(r.Intn(B+B)) - float64(B)) / float64(B)
			w4 = (float64(r.Intn(B+B)) - float64(B)) / float64(B)
			v5 = (float64(r.Intn(B+B)) - float64(B)) / float64(B)
			u6 = (float64(r.Intn(B+B)) - float64(B)) / float64(B)
			magSq := x*x + y*y + z*z + w4*w4 + v5*v5 + u6*u6
			if magSq <= 1.0 && magSq > 0.0 {
				break
			}
		}
		mag6D := math.Sqrt(x*x + y*y + z*z + w4*w4 + v5*v5 + u6*u6)
		p.g6[i] = [6]float64{x / mag6D, y / mag6D, z / mag6D, w4 / mag6D, v5 / mag6D, u6 / mag6D}
	}

	for i := 0; i < B+2; i++ {
		p.g6[B+i] = p.g6[i]
	}
}

func at2(rx, ry float64, q [2]float64) float64 {
	return rx*q[0] + ry*q[1]
}

func at3(rx, ry, rz float64, q [3]float64) float64 {
	return rx*q[0] + ry*q[1] + rz*q[2]
}

func (p *OldPerlin) noise1(arg float64) float64 {
	var vec [1]float64
	vec[0] = arg

	t := vec[0] + N
	bx0 := int32(t) & BM
	bx1 := (bx0 + 1) & BM
	rx0 := t - float64(int32(t))
	rx1 := rx0 - 1.

	sx := xmath.Scurve(rx0)
	u := rx0 * p.g1[p.p[bx0]]
	v := rx1 * p.g1[p.p[bx1]]

	return xmath.Lerp(sx, u, v)
}

func (p *OldPerlin) noise2(vec [2]float64) float64 {
	t := vec[0] + N
	bx0 := int32(t) & BM
	bx1 := (bx0 + 1) & BM
	rx0 := t - float64(int32(t))
	rx1 := rx0 - 1.

	t = vec[1] + N
	by0 := int32(t) & BM
	by1 := (by0 + 1) & BM
	ry0 := t - float64(int32(t))
	ry1 := ry0 - 1.

	i := p.p[bx0]
	j := p.p[bx1]

	b00 := p.p[i+by0]
	b10 := p.p[j+by0]
	b01 := p.p[i+by1]
	b11 := p.p[j+by1]

	sx := xmath.Scurve(rx0)
	sy := xmath.Scurve(ry0)

	q := p.g2[b00]
	u := at2(rx0, ry0, q)
	q = p.g2[b10]
	v := at2(rx1, ry0, q)
	a := xmath.Lerp(sx, u, v)

	q = p.g2[b01]
	u = at2(rx0, ry1, q)
	q = p.g2[b11]
	v = at2(rx1, ry1, q)
	b := xmath.Lerp(sx, u, v)

	return xmath.Lerp(sy, a, b)
}

func (p *OldPerlin) noise3(vec [3]float64) float64 {
	t := vec[0] + N
	bx0 := int32(t) & BM
	bx1 := (bx0 + 1) & BM
	rx0 := t - float64(int32(t))
	rx1 := rx0 - 1.

	t = vec[1] + N
	by0 := int32(t) & BM
	by1 := (by0 + 1) & BM
	ry0 := t - float64(int32(t))
	ry1 := ry0 - 1.

	t = vec[2] + N
	bz0 := int32(t) & BM
	bz1 := (bz0 + 1) & BM
	rz0 := t - float64(int32(t))
	rz1 := rz0 - 1.

	i := p.p[bx0]
	j := p.p[bx1]

	b00 := p.p[i+by0]
	b10 := p.p[j+by0]
	b01 := p.p[i+by1]
	b11 := p.p[j+by1]

	t = xmath.Scurve(rx0)
	sy := xmath.Scurve(ry0)
	sz := xmath.Scurve(rz0)

	q := p.g3[b00+bz0]
	u := at3(rx0, ry0, rz0, q)
	q = p.g3[b10+bz0]
	v := at3(rx1, ry0, rz0, q)
	a := xmath.Lerp(t, u, v)

	q = p.g3[b01+bz0]
	u = at3(rx0, ry1, rz0, q)
	q = p.g3[b11+bz0]
	v = at3(rx1, ry1, rz0, q)
	b := xmath.Lerp(t, u, v)

	c := xmath.Lerp(sy, a, b)

	q = p.g3[b00+bz1]
	u = at3(rx0, ry0, rz1, q)
	q = p.g3[b10+bz1]
	v = at3(rx1, ry0, rz1, q)
	a = xmath.Lerp(t, u, v)

	q = p.g3[b01+bz1]
	u = at3(rx0, ry1, rz1, q)
	q = p.g3[b11+bz1]
	v = at3(rx1, ry1, rz1, q)
	b = xmath.Lerp(t, u, v)

	d := xmath.Lerp(sy, a, b)

	return xmath.Lerp(sz, c, d)
}

func at4(rx, ry, rz, rw float64, q [4]float64) float64 {
	return rx*q[0] + ry*q[1] + rz*q[2] + rw*q[3]
}

func (p *OldPerlin) noise4(vec [4]float64) float64 {
	// 1. Calculate Grid Coordinates and Fractional Distances for X
	t := vec[0] + N
	bx0 := int32(t) & BM
	bx1 := (bx0 + 1) & BM
	rx0 := t - float64(int32(t))
	rx1 := rx0 - 1.

	// 2. Calculate Grid Coordinates and Fractional Distances for Y
	t = vec[1] + N
	by0 := int32(t) & BM
	by1 := (by0 + 1) & BM
	ry0 := t - float64(int32(t))
	ry1 := ry0 - 1.

	// 3. Calculate Grid Coordinates and Fractional Distances for Z
	t = vec[2] + N
	bz0 := int32(t) & BM
	bz1 := (bz0 + 1) & BM
	rz0 := t - float64(int32(t))
	rz1 := rz0 - 1.

	// 4. Calculate Grid Coordinates and Fractional Distances for W
	t = vec[3] + N
	bw0 := int32(t) & BM
	bw1 := (bw0 + 1) & BM
	rw0 := t - float64(int32(t))
	rw1 := rw0 - 1.

	// 5. Permutation Table Lookups for X and Y
	i := p.p[bx0]
	j := p.p[bx1]

	b00 := p.p[i+by0]
	b10 := p.p[j+by0]
	b01 := p.p[i+by1]
	b11 := p.p[j+by1]

	// 6. S-Curve Smoothing Factors
	sx := xmath.Scurve(rx0)
	sy := xmath.Scurve(ry0)
	sz := xmath.Scurve(rz0)
	sw := xmath.Scurve(rw0)

	// --- W = bw0 Hyperplane ---

	// Z = bz0
	q := p.g4[p.p[b00+bz0]+bw0]
	u := at4(rx0, ry0, rz0, rw0, q)
	q = p.g4[p.p[b10+bz0]+bw0]
	v := at4(rx1, ry0, rz0, rw0, q)
	a := xmath.Lerp(sx, u, v)

	q = p.g4[p.p[b01+bz0]+bw0]
	u = at4(rx0, ry1, rz0, rw0, q)
	q = p.g4[p.p[b11+bz0]+bw0]
	v = at4(rx1, ry1, rz0, rw0, q)
	b := xmath.Lerp(sx, u, v)

	c0 := xmath.Lerp(sy, a, b)

	// Z = bz1
	q = p.g4[p.p[b00+bz1]+bw0]
	u = at4(rx0, ry0, rz1, rw0, q)
	q = p.g4[p.p[b10+bz1]+bw0]
	v = at4(rx1, ry0, rz1, rw0, q)
	a = xmath.Lerp(sx, u, v)

	q = p.g4[p.p[b01+bz1]+bw0]
	u = at4(rx0, ry1, rz1, rw0, q)
	q = p.g4[p.p[b11+bz1]+bw0]
	v = at4(rx1, ry1, rz1, rw0, q)
	b = xmath.Lerp(sx, u, v)

	d0 := xmath.Lerp(sy, a, b)

	// Interpolate across Z for W0
	w0 := xmath.Lerp(sz, c0, d0)

	// --- W = bw1 Hyperplane ---

	// Z = bz0
	q = p.g4[p.p[b00+bz0]+bw1]
	u = at4(rx0, ry0, rz0, rw1, q)
	q = p.g4[p.p[b10+bz0]+bw1]
	v = at4(rx1, ry0, rz0, rw1, q)
	a = xmath.Lerp(sx, u, v)

	q = p.g4[p.p[b01+bz0]+bw1]
	u = at4(rx0, ry1, rz0, rw1, q)
	q = p.g4[p.p[b11+bz0]+bw1]
	v = at4(rx1, ry1, rz0, rw1, q)
	b = xmath.Lerp(sx, u, v)

	c1 := xmath.Lerp(sy, a, b)

	// Z = bz1
	q = p.g4[p.p[b00+bz1]+bw1]
	u = at4(rx0, ry0, rz1, rw1, q)
	q = p.g4[p.p[b10+bz1]+bw1]
	v = at4(rx1, ry0, rz1, rw1, q)
	a = xmath.Lerp(sx, u, v)

	q = p.g4[p.p[b01+bz1]+bw1]
	u = at4(rx0, ry1, rz1, rw1, q)
	q = p.g4[p.p[b11+bz1]+bw1]
	v = at4(rx1, ry1, rz1, rw1, q)
	b = xmath.Lerp(sx, u, v)

	d1 := xmath.Lerp(sy, a, b)

	// Interpolate across Z for W1
	w1 := xmath.Lerp(sz, c1, d1)

	// 7. Final Interpolation along W
	return xmath.Lerp(sw, w0, w1)
}

func at5(rx, ry, rz, rw, rv float64, q [5]float64) float64 {
	return rx*q[0] + ry*q[1] + rz*q[2] + rw*q[3] + rv*q[4]
}

func (p *OldPerlin) noise5(vec [5]float64) float64 {
	t := vec[0] + N
	bx0 := int32(t) & BM
	bx1 := (bx0 + 1) & BM
	rx0 := t - float64(int32(t))
	rx1 := rx0 - 1.

	t = vec[1] + N
	by0 := int32(t) & BM
	by1 := (by0 + 1) & BM
	ry0 := t - float64(int32(t))
	ry1 := ry0 - 1.

	t = vec[2] + N
	bz0 := int32(t) & BM
	bz1 := (bz0 + 1) & BM
	rz0 := t - float64(int32(t))
	rz1 := rz0 - 1.

	t = vec[3] + N
	bw0 := int32(t) & BM
	bw1 := (bw0 + 1) & BM
	rw0 := t - float64(int32(t))
	rw1 := rw0 - 1.

	t = vec[4] + N
	bv0 := int32(t) & BM
	bv1 := (bv0 + 1) & BM
	rv0 := t - float64(int32(t))
	rv1 := rv0 - 1.

	i := p.p[bx0]
	j := p.p[bx1]

	b00 := p.p[i+by0]
	b10 := p.p[j+by0]
	b01 := p.p[i+by1]
	b11 := p.p[j+by1]

	sx := xmath.Scurve(rx0)
	sy := xmath.Scurve(ry0)
	sz := xmath.Scurve(rz0)
	sw := xmath.Scurve(rw0)
	sv := xmath.Scurve(rv0)

	// Helper closure to compute 4D noise slice at a specific V coordinate index and offset
	eval4D := func(bw int32, rw float64, bv int32, rv float64) float64 {
		// Z = bz0
		q := p.g5[p.p[p.p[b00+bz0]+bw]+bv]
		u := at5(rx0, ry0, rz0, rw, rv, q)
		q = p.g5[p.p[p.p[b10+bz0]+bw]+bv]
		v := at5(rx1, ry0, rz0, rw, rv, q)
		a := xmath.Lerp(sx, u, v)

		q = p.g5[p.p[p.p[b01+bz0]+bw]+bv]
		u = at5(rx0, ry1, rz0, rw, rv, q)
		q = p.g5[p.p[p.p[b11+bz0]+bw]+bv]
		v = at5(rx1, ry1, rz0, rw, rv, q)
		b := xmath.Lerp(sx, u, v)

		c0 := xmath.Lerp(sy, a, b)

		// Z = bz1
		q = p.g5[p.p[p.p[b00+bz1]+bw]+bv]
		u = at5(rx0, ry0, rz1, rw, rv, q)
		q = p.g5[p.p[p.p[b10+bz1]+bw]+bv]
		v = at5(rx1, ry0, rz1, rw, rv, q)
		a = xmath.Lerp(sx, u, v)

		q = p.g5[p.p[p.p[b01+bz1]+bw]+bv]
		u = at5(rx0, ry1, rz1, rw, rv, q)
		q = p.g5[p.p[p.p[b11+bz1]+bw]+bv]
		v = at5(rx1, ry1, rz1, rw, rv, q)
		b = xmath.Lerp(sx, u, v)

		d0 := xmath.Lerp(sy, a, b)

		return xmath.Lerp(sz, c0, d0)
	}

	// --- V = bv0 Hypercube ---
	v0_w0 := eval4D(bw0, rw0, bv0, rv0)
	v0_w1 := eval4D(bw1, rw1, bv0, rv0)
	v0 := xmath.Lerp(sw, v0_w0, v0_w1)

	// --- V = bv1 Hypercube ---
	v1_w0 := eval4D(bw0, rw0, bv1, rv1)
	v1_w1 := eval4D(bw1, rw1, bv1, rv1)
	v1 := xmath.Lerp(sw, v1_w0, v1_w1)

	// Final interpolation across V dimension
	return xmath.Lerp(sv, v0, v1)
}

// 6D Gradient Dot-Product Helper
func at6(rx, ry, rz, rw, rv, ru float64, q [6]float64) float64 {
	return rx*q[0] + ry*q[1] + rz*q[2] + rw*q[3] + rv*q[4] + ru*q[5]
}

// noise6 evaluates 6D Perlin noise over 64 hypercube corners
func (p *OldPerlin) noise6(vec [6]float64) float64 {
	t := vec[0] + N
	bx0 := int32(t) & BM
	bx1 := (bx0 + 1) & BM
	rx0 := t - float64(int32(t))
	rx1 := rx0 - 1.

	t = vec[1] + N
	by0 := int32(t) & BM
	by1 := (by0 + 1) & BM
	ry0 := t - float64(int32(t))
	ry1 := ry0 - 1.

	t = vec[2] + N
	bz0 := int32(t) & BM
	bz1 := (bz0 + 1) & BM
	rz0 := t - float64(int32(t))
	rz1 := rz0 - 1.

	t = vec[3] + N
	bw0 := int32(t) & BM
	bw1 := (bw0 + 1) & BM
	rw0 := t - float64(int32(t))
	rw1 := rw0 - 1.

	t = vec[4] + N
	bv0 := int32(t) & BM
	bv1 := (bv0 + 1) & BM
	rv0 := t - float64(int32(t))
	rv1 := rv0 - 1.

	t = vec[5] + N
	bu0 := int32(t) & BM
	bu1 := (bu0 + 1) & BM
	ru0 := t - float64(int32(t))
	ru1 := ru0 - 1.

	i := p.p[bx0]
	j := p.p[bx1]

	b00 := p.p[i+by0]
	b10 := p.p[j+by0]
	b01 := p.p[i+by1]
	b11 := p.p[j+by1]

	sx := xmath.Scurve(rx0)
	sy := xmath.Scurve(ry0)
	sz := xmath.Scurve(rz0)
	sw := xmath.Scurve(rw0)
	sv := xmath.Scurve(rv0)
	su := xmath.Scurve(ru0)

	// Helper closure to evaluate a 4D plane slice given specific (W, V, U) grid coords
	eval4DSlice := func(bw int32, rw float64, bv int32, rv float64, bu int32, ru float64) float64 {
		// Z = bz0
		q := p.g6[p.p[p.p[p.p[b00+bz0]+bw]+bv]+bu]
		uVal := at6(rx0, ry0, rz0, rw, rv, ru, q)
		q = p.g6[p.p[p.p[p.p[b10+bz0]+bw]+bv]+bu]
		vVal := at6(rx1, ry0, rz0, rw, rv, ru, q)
		a := xmath.Lerp(sx, uVal, vVal)

		q = p.g6[p.p[p.p[p.p[b01+bz0]+bw]+bv]+bu]
		uVal = at6(rx0, ry1, rz0, rw, rv, ru, q)
		q = p.g6[p.p[p.p[p.p[b11+bz0]+bw]+bv]+bu]
		vVal = at6(rx1, ry1, rz0, rw, rv, ru, q)
		b := xmath.Lerp(sx, uVal, vVal)

		c0 := xmath.Lerp(sy, a, b)

		// Z = bz1
		q = p.g6[p.p[p.p[p.p[b00+bz1]+bw]+bv]+bu]
		uVal = at6(rx0, ry0, rz1, rw, rv, ru, q)
		q = p.g6[p.p[p.p[p.p[b10+bz1]+bw]+bv]+bu]
		vVal = at6(rx1, ry0, rz1, rw, rv, ru, q)
		a = xmath.Lerp(sx, uVal, vVal)

		q = p.g6[p.p[p.p[p.p[b01+bz1]+bw]+bv]+bu]
		uVal = at6(rx0, ry1, rz1, rw, rv, ru, q)
		q = p.g6[p.p[p.p[p.p[b11+bz1]+bw]+bv]+bu]
		vVal = at6(rx1, ry1, rz1, rw, rv, ru, q)
		b = xmath.Lerp(sx, uVal, vVal)

		d0 := xmath.Lerp(sy, a, b)

		return xmath.Lerp(sz, c0, d0)
	}

	// Helper to evaluate a 5D hypercube slice given a specific U coordinate
	eval5DSlice := func(bu int32, ru float64) float64 {
		v0_w0 := eval4DSlice(bw0, rw0, bv0, rv0, bu, ru)
		v0_w1 := eval4DSlice(bw1, rw1, bv0, rv0, bu, ru)
		v0 := xmath.Lerp(sw, v0_w0, v0_w1)

		v1_w0 := eval4DSlice(bw0, rw0, bv1, rv1, bu, ru)
		v1_w1 := eval4DSlice(bw1, rw1, bv1, rv1, bu, ru)
		v1 := xmath.Lerp(sw, v1_w0, v1_w1)

		return xmath.Lerp(sv, v0, v1)
	}

	// Interpolate across U dimension
	u0 := eval5DSlice(bu0, ru0)
	u1 := eval5DSlice(bu1, ru1)

	return xmath.Lerp(su, u0, u1)
}

// Noise1D generates 1-dimensional Perlin Noise value
func (p *OldPerlin) Noise1D(x float64) float64 {
	return p.noise1(x)
}

// Noise2D Generates 2-dimensional Perlin Noise value
func (p *OldPerlin) Noise2D(x, y float64) float64 {
	return p.noise2([2]float64{x, y})
}

// Noise3D Generates 3-dimensional Perlin Noise value
func (p *OldPerlin) Noise3D(x, y, z float64) float64 {
	return p.noise3([3]float64{x, y, z})
}

// Noise4D Generates 4-dimensional Perlin Noise value
func (p *OldPerlin) Noise4D(x, y, z, w float64) float64 {
	return p.noise4([4]float64{x, y, z, w})
}

// Noise5D Generates 5-dimensional Perlin Noise value
func (p *OldPerlin) Noise5D(x, y, z, w, v float64) float64 {
	return p.noise5([5]float64{x, y, z, w, v})
}

// Noise6D Generates 6-dimensional Perlin Noise value
func (p *OldPerlin) Noise6D(x, y, z, w, v, u float64) float64 {
	return p.noise6([6]float64{x, y, z, w, v, u})
}
