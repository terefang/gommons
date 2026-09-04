package noise

import (
	"math/rand"

	"github.com/terefang/gommons/pkg/xmath"
)

// SolidNoise implements multi-dimensional Perlin Solid Noise.
type SolidNoise struct {
	perm [512]uint8
}

// NewSolidNoise creates a new SolidNoise generator initialized with a seed.
func NewSolidNoise(seed int64) *SolidNoise {
	return NewSolidNoiseWithSource(rand.NewSource(seed))
}

func NewSolidNoiseWithSource(source rand.Source) *SolidNoise {
	var p [256]uint8
	for i, _ := range p {
		p[i] = uint8(i)
	}
	r := rand.New(source)
	for i := 255; i > 0; i-- {
		j := r.Intn(i + 1)
		p[i], p[j] = p[j], p[i]
	}

	sn := &SolidNoise{}
	for i := 0; i < 512; i++ {
		sn.perm[i] = p[i&0xff]
	}
	return sn
}

func grad1D(hash uint8, x float64) float64 {
	if hash&1 == 0 {
		return x
	}
	return -x
}

func grad2D(hash uint8, x, y float64) float64 {
	h := hash & 7
	u := x
	v := y
	if h >= 4 {
		u, v = y, x
	}
	if h&1 != 0 {
		u = -u
	}
	if h&2 != 0 {
		v = -v
	}
	return u + v
}

func grad3D(hash uint8, x, y, z float64) float64 {
	h := hash & 15
	u := x
	if h < 8 {
		u = x
	} else {
		u = y
	}
	v := y
	if h < 4 {
		v = y
	} else if h == 12 || h == 14 {
		v = x
	} else {
		v = z
	}
	res := 0.0
	if h&1 != 0 {
		res -= u
	} else {
		res += u
	}
	if h&2 != 0 {
		res -= v
	} else {
		res += v
	}
	return res
}

func grad4D(hash uint8, x, y, z, w float64) float64 {
	h := hash & 31
	var u, v, p float64

	if h < 24 {
		u = x
	} else {
		u = y
	}
	if h < 16 {
		v = y
	} else {
		v = z
	}
	if h < 8 {
		p = z
	} else {
		p = w
	}

	res := 0.0
	if h&1 != 0 {
		res -= u
	} else {
		res += u
	}
	if h&2 != 0 {
		res -= v
	} else {
		res += v
	}
	if h&4 != 0 {
		res -= p
	} else {
		res += p
	}
	return res
}

func grad5D(hash uint8, x, y, z, w, v float64) float64 {
	h := hash & 31
	res := 0.0
	if h&1 != 0 {
		res -= x
	} else {
		res += x
	}
	if h&2 != 0 {
		res -= y
	} else {
		res += y
	}
	if h&4 != 0 {
		res -= z
	} else {
		res += z
	}
	if h&8 != 0 {
		res -= w
	} else {
		res += w
	}
	if h&16 != 0 {
		res -= v
	} else {
		res += v
	}
	return res
}

func grad6D(hash uint8, x, y, z, w, v, u float64) float64 {
	h := hash & 63
	res := 0.0
	if h&1 != 0 {
		res -= x
	} else {
		res += x
	}
	if h&2 != 0 {
		res -= y
	} else {
		res += y
	}
	if h&4 != 0 {
		res -= z
	} else {
		res += z
	}
	if h&8 != 0 {
		res -= w
	} else {
		res += w
	}
	if h&16 != 0 {
		res -= v
	} else {
		res += v
	}
	if h&32 != 0 {
		res -= u
	} else {
		res += u
	}
	return res
}

// Noise1D returns 1D Solid Perlin Noise scaled to [-1, 1].
func (sn *SolidNoise) Noise1D(x float64) float64 {
	X0 := xmath.FastFloor(x)
	X1 := X0 + 1

	fx0 := x - float64(X0)
	fx1 := fx0 - 1.0

	u := xmath.Scurve(fx0)

	ix0 := uint8(X0)
	ix1 := uint8(X1)

	g0 := grad1D(sn.perm[ix0], fx0)
	g1 := grad1D(sn.perm[ix1], fx1)

	return xmath.Lerp(u, g0, g1)
}

// Noise2D returns 2D Solid Perlin Noise scaled to [-1, 1].
func (sn *SolidNoise) Noise2D(x, y float64) float64 {
	X0, Y0 := xmath.FastFloor(x), xmath.FastFloor(y)
	X1, Y1 := X0+1, Y0+1

	fx0, fy0 := x-float64(X0), y-float64(Y0)
	fx1, fy1 := fx0-1.0, fy0-1.0

	u := xmath.Scurve(fx0)
	v := xmath.Scurve(fy0)

	ix0, iy0 := uint8(X0), uint8(Y0)
	ix1, iy1 := uint8(X1), uint8(Y1)

	g00 := grad2D(sn.perm[ix0+sn.perm[iy0]], fx0, fy0)
	g10 := grad2D(sn.perm[ix1+sn.perm[iy0]], fx1, fy0)
	g01 := grad2D(sn.perm[ix0+sn.perm[iy1]], fx0, fy1)
	g11 := grad2D(sn.perm[ix1+sn.perm[iy1]], fx1, fy1)

	nx0 := xmath.Lerp(u, g00, g10)
	nx1 := xmath.Lerp(u, g01, g11)

	return xmath.Lerp(v, nx0, nx1)
}

// Noise3D returns 3D Solid Perlin Noise scaled to [-1, 1].
func (sn *SolidNoise) Noise3D(x, y, z float64) float64 {
	X0, Y0, Z0 := xmath.FastFloor(x), xmath.FastFloor(y), xmath.FastFloor(z)
	X1, Y1, Z1 := X0+1, Y0+1, Z0+1

	fx0, fy0, fz0 := x-float64(X0), y-float64(Y0), z-float64(Z0)
	fx1, fy1, fz1 := fx0-1.0, fy0-1.0, fz0-1.0

	u, v, w := xmath.Scurve(fx0), xmath.Scurve(fy0), xmath.Scurve(fz0)

	ix0, iy0, iz0 := uint8(X0), uint8(Y0), uint8(Z0)
	ix1, iy1, iz1 := uint8(X1), uint8(Y1), uint8(Z1)

	p00 := ix0 + sn.perm[iy0+sn.perm[iz0]]
	p10 := ix1 + sn.perm[iy0+sn.perm[iz0]]
	p01 := ix0 + sn.perm[iy1+sn.perm[iz0]]
	p11 := ix1 + sn.perm[iy1+sn.perm[iz0]]
	p001 := ix0 + sn.perm[iy0+sn.perm[iz1]]
	p101 := ix1 + sn.perm[iy0+sn.perm[iz1]]
	p011 := ix0 + sn.perm[iy1+sn.perm[iz1]]
	p111 := ix1 + sn.perm[iy1+sn.perm[iz1]]

	g000, g100 := grad3D(sn.perm[p00], fx0, fy0, fz0), grad3D(sn.perm[p10], fx1, fy0, fz0)
	g010, g110 := grad3D(sn.perm[p01], fx0, fy1, fz0), grad3D(sn.perm[p11], fx1, fy1, fz0)
	g001, g101 := grad3D(sn.perm[p001], fx0, fy0, fz1), grad3D(sn.perm[p101], fx1, fy0, fz1)
	g011, g111 := grad3D(sn.perm[p011], fx0, fy1, fz1), grad3D(sn.perm[p111], fx1, fy1, fz1)

	nx00 := xmath.Lerp(u, g000, g100)
	nx10 := xmath.Lerp(u, g010, g110)
	nx01 := xmath.Lerp(u, g001, g101)
	nx11 := xmath.Lerp(u, g011, g111)

	ny0 := xmath.Lerp(v, nx00, nx10)
	ny1 := xmath.Lerp(v, nx01, nx11)

	return xmath.Lerp(w, ny0, ny1)
}

// Noise4D returns 4D Solid Perlin Noise scaled to [-1, 1].
func (sn *SolidNoise) Noise4D(x, y, z, w float64) float64 {
	X0, Y0, Z0, W0 := xmath.FastFloor(x), xmath.FastFloor(y), xmath.FastFloor(z), xmath.FastFloor(w)

	fx0, fy0, fz0, fw0 := x-float64(X0), y-float64(Y0), z-float64(Z0), w-float64(W0)
	fx1, fy1, fz1, fw1 := fx0-1.0, fy0-1.0, fz0-1.0, fw0-1.0

	fu, fv, fw, ft := xmath.Scurve(fx0), xmath.Scurve(fy0), xmath.Scurve(fz0), xmath.Scurve(fw0)

	ix0, iy0, iz0, iw0 := uint8(X0), uint8(Y0), uint8(Z0), uint8(W0)
	ix1, iy1, iz1, iw1 := uint8(X0+1), uint8(Y0+1), uint8(Z0+1), uint8(W0+1)

	var g [2][2][2][2]float64
	coordsX := [2]float64{fx0, fx1}
	coordsY := [2]float64{fy0, fy1}
	coordsZ := [2]float64{fz0, fz1}
	coordsW := [2]float64{fw0, fw1}

	for dw := 0; dw < 2; dw++ {
		iw := iw0
		if dw == 1 {
			iw = iw1
		}
		for dz := 0; dz < 2; dz++ {
			iz := iz0
			if dz == 1 {
				iz = iz1
			}
			for dy := 0; dy < 2; dy++ {
				iy := iy0
				if dy == 1 {
					iy = iy1
				}
				p := sn.perm[iy+sn.perm[iz+sn.perm[iw]]]
				g[dw][dz][dy][0] = grad4D(sn.perm[ix0+p], coordsX[0], coordsY[dy], coordsZ[dz], coordsW[dw])
				g[dw][dz][dy][1] = grad4D(sn.perm[ix1+p], coordsX[1], coordsY[dy], coordsZ[dz], coordsW[dw])
			}
		}
	}

	var wLerp [2][2][2]float64
	for dw := 0; dw < 2; dw++ {
		for dz := 0; dz < 2; dz++ {
			for dy := 0; dy < 2; dy++ {
				wLerp[dw][dz][dy] = xmath.Lerp(fu, g[dw][dz][dy][0], g[dw][dz][dy][1])
			}
		}
	}

	var zLerp [2][2]float64
	for dw := 0; dw < 2; dw++ {
		for dz := 0; dz < 2; dz++ {
			zLerp[dw][dz] = xmath.Lerp(fv, wLerp[dw][dz][0], wLerp[dw][dz][1])
		}
	}

	var yLerp [2]float64
	for dw := 0; dw < 2; dw++ {
		yLerp[dw] = xmath.Lerp(fw, zLerp[dw][0], zLerp[dw][1])
	}

	return xmath.Lerp(ft, yLerp[0], yLerp[1])
}

// Noise5D returns 5D Solid Perlin Noise scaled to [-1, 1].
func (sn *SolidNoise) Noise5D(x, y, z, w, v float64) float64 {
	X0, Y0, Z0, W0, V0 := xmath.FastFloor(x), xmath.FastFloor(y), xmath.FastFloor(z), xmath.FastFloor(w), xmath.FastFloor(v)

	fx0, fy0, fz0, fw0, fv0 := x-float64(X0), y-float64(Y0), z-float64(Z0), w-float64(W0), v-float64(V0)

	fx := xmath.Scurve(fx0)
	fy := xmath.Scurve(fy0)
	fz := xmath.Scurve(fz0)
	fw := xmath.Scurve(fw0)
	fv := xmath.Scurve(fv0)

	coords := [5][2]uint8{
		{uint8(X0), uint8(X0 + 1)},
		{uint8(Y0), uint8(Y0 + 1)},
		{uint8(Z0), uint8(Z0 + 1)},
		{uint8(W0), uint8(W0 + 1)},
		{uint8(V0), uint8(V0 + 1)},
	}

	vecs := [5][2]float64{
		{fx0, fx0 - 1.0},
		{fy0, fy0 - 1.0},
		{fz0, fz0 - 1.0},
		{fw0, fw0 - 1.0},
		{fv0, fv0 - 1.0},
	}

	var corners [32]float64
	for i := 0; i < 32; i++ {
		bx, by, bz, bw, bv := (i>>0)&1, (i>>1)&1, (i>>2)&1, (i>>3)&1, (i>>4)&1

		hash := coords[0][bx] + sn.perm[coords[1][by]+sn.perm[coords[2][bz]+sn.perm[coords[3][bw]+sn.perm[coords[4][bv]]]]]
		corners[i] = grad5D(sn.perm[hash], vecs[0][bx], vecs[1][by], vecs[2][bz], vecs[3][bw], vecs[4][bv])
	}

	factors := [5]float64{fx, fy, fz, fw, fv}
	for step := 0; step < 5; step++ {
		f := factors[step]
		bound := 1 << (5 - step)
		for i := 0; i < bound; i += 2 {
			corners[i/2] = xmath.Lerp(f, corners[i], corners[i+1])
		}
	}

	return corners[0]
}

// Noise6D returns 6D Solid Perlin Noise scaled to [-1, 1].
func (sn *SolidNoise) Noise6D(x, y, z, w, v, u float64) float64 {
	X0, Y0, Z0, W0, V0, U0 := xmath.FastFloor(x), xmath.FastFloor(y), xmath.FastFloor(z), xmath.FastFloor(w), xmath.FastFloor(v), xmath.FastFloor(u)

	fx0, fy0, fz0, fw0, fv0, fu0 := x-float64(X0), y-float64(Y0), z-float64(Z0), w-float64(W0), v-float64(V0), u-float64(U0)

	fx := xmath.Scurve(fx0)
	fy := xmath.Scurve(fy0)
	fz := xmath.Scurve(fz0)
	fw := xmath.Scurve(fw0)
	fv := xmath.Scurve(fv0)
	fu := xmath.Scurve(fu0)

	coords := [6][2]uint8{
		{uint8(X0), uint8(X0 + 1)},
		{uint8(Y0), uint8(Y0 + 1)},
		{uint8(Z0), uint8(Z0 + 1)},
		{uint8(W0), uint8(W0 + 1)},
		{uint8(V0), uint8(V0 + 1)},
		{uint8(U0), uint8(U0 + 1)},
	}

	vecs := [6][2]float64{
		{fx0, fx0 - 1.0},
		{fy0, fy0 - 1.0},
		{fz0, fz0 - 1.0},
		{fw0, fw0 - 1.0},
		{fv0, fv0 - 1.0},
		{fu0, fu0 - 1.0},
	}

	var corners [64]float64
	for i := 0; i < 64; i++ {
		bx, by, bz, bw, bv, bu := (i>>0)&1, (i>>1)&1, (i>>2)&1, (i>>3)&1, (i>>4)&1, (i>>5)&1

		hash := coords[0][bx] + sn.perm[coords[1][by]+sn.perm[coords[2][bz]+sn.perm[coords[3][bw]+sn.perm[coords[4][bv]+sn.perm[coords[5][bu]]]]]]
		corners[i] = grad6D(sn.perm[hash], vecs[0][bx], vecs[1][by], vecs[2][bz], vecs[3][bw], vecs[4][bv], vecs[5][bu])
	}

	factors := [6]float64{fx, fy, fz, fw, fv, fu}
	for step := 0; step < 6; step++ {
		f := factors[step]
		bound := 1 << (6 - step)
		for i := 0; i < bound; i += 2 {
			corners[i/2] = xmath.Lerp(f, corners[i], corners[i+1])
		}
	}

	return corners[0]
}
