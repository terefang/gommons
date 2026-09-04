package noise

import (
	"math/rand"

	"github.com/terefang/gommons/pkg/xmath"
)

// ValueNoise implements NoiseSource using Value Noise interpolation.
type ValueNoise struct {
	perm   [512]uint8
	values [256]float64
}

func NewValueNoise(seed int64) *ValueNoise {
	return NewValueNoiseWithSource(rand.NewSource(seed))
}

func NewValueNoiseWithSource(source rand.Source) *ValueNoise {
	vn := &ValueNoise{}

	var p [256]uint8
	for i := range p {
		p[i] = uint8(i)
	}

	// Deterministic LCG + Fisher-Yates shuffle
	rng := rand.New(source)
	for i := 255; i > 0; i-- {
		j := rng.Intn(i + 1)
		p[i], p[j] = p[j], p[i]
	}

	for i := 0; i < 512; i++ {
		vn.perm[i] = p[i&255]
	}

	// Precompute normalized continuous values [-1.0, 1.0] for the permutation table
	for i := 0; i < 256; i++ {
		vn.values[i] = (float64(i) / 127.5) - 1.0
	}

	return vn
}

func (vn *ValueNoise) eval(hash uint8) float64 {
	return vn.values[vn.perm[hash]]
}

// Noise1D returns 1D Value Noise scaled to [-1, 1].
func (vn *ValueNoise) Noise1D(x float64) float64 {
	X0 := xmath.FastFloor(x)
	X1 := X0 + 1

	fx := xmath.Scurve(x - float64(X0))

	v0 := vn.eval(uint8(X0))
	v1 := vn.eval(uint8(X1))

	return xmath.Lerp(fx, v0, v1)
}

// Noise2D returns 2D Value Noise scaled to [-1, 1].
func (vn *ValueNoise) Noise2D(x, y float64) float64 {
	X0, Y0 := xmath.FastFloor(x), xmath.FastFloor(y)
	X1, Y1 := X0+1, Y0+1

	fx := xmath.Scurve(x - float64(X0))
	fy := xmath.Scurve(y - float64(Y0))

	ix0, iy0 := uint8(X0), uint8(Y0)
	ix1, iy1 := uint8(X1), uint8(Y1)

	v00 := vn.eval(ix0 + vn.perm[iy0])
	v10 := vn.eval(ix1 + vn.perm[iy0])
	v01 := vn.eval(ix0 + vn.perm[iy1])
	v11 := vn.eval(ix1 + vn.perm[iy1])

	nx0 := xmath.Lerp(fx, v00, v10)
	nx1 := xmath.Lerp(fx, v01, v11)

	return xmath.Lerp(fy, nx0, nx1)
}

// Noise3D returns 3D Value Noise scaled to [-1, 1].
func (vn *ValueNoise) Noise3D(x, y, z float64) float64 {
	X0, Y0, Z0 := xmath.FastFloor(x), xmath.FastFloor(y), xmath.FastFloor(z)
	X1, Y1, Z1 := X0+1, Y0+1, Z0+1

	fx := xmath.Scurve(x - float64(X0))
	fy := xmath.Scurve(y - float64(Y0))
	fz := xmath.Scurve(z - float64(Z0))

	ix0, iy0, iz0 := uint8(X0), uint8(Y0), uint8(Z0)
	ix1, iy1, iz1 := uint8(X1), uint8(Y1), uint8(Z1)

	p00 := ix0 + vn.perm[iy0+vn.perm[iz0]]
	p10 := ix1 + vn.perm[iy0+vn.perm[iz0]]
	p01 := ix0 + vn.perm[iy1+vn.perm[iz0]]
	p11 := ix1 + vn.perm[iy1+vn.perm[iz0]]
	p00_1 := ix0 + vn.perm[iy0+vn.perm[iz1]]
	p10_1 := ix1 + vn.perm[iy0+vn.perm[iz1]]
	p01_1 := ix0 + vn.perm[iy1+vn.perm[iz1]]
	p11_1 := ix1 + vn.perm[iy1+vn.perm[iz1]]

	v000, v100 := vn.eval(p00), vn.eval(p10)
	v010, v110 := vn.eval(p01), vn.eval(p11)
	v001, v101 := vn.eval(p00_1), vn.eval(p10_1)
	v011, v111 := vn.eval(p01_1), vn.eval(p11_1)

	nx00 := xmath.Lerp(fx, v000, v100)
	nx10 := xmath.Lerp(fx, v010, v110)
	nx01 := xmath.Lerp(fx, v001, v101)
	nx11 := xmath.Lerp(fx, v011, v111)

	ny0 := xmath.Lerp(fy, nx00, nx10)
	ny1 := xmath.Lerp(fy, nx01, nx11)

	return xmath.Lerp(fz, ny0, ny1)
}

// Noise4D returns 4D Value Noise scaled to [-1, 1].
func (vn *ValueNoise) Noise4D(x, y, z, w float64) float64 {
	X0, Y0, Z0, W0 := xmath.FastFloor(x), xmath.FastFloor(y), xmath.FastFloor(z), xmath.FastFloor(w)
	fx, fy, fz, fw := xmath.Scurve(x-float64(X0)), xmath.Scurve(y-float64(Y0)), xmath.Scurve(z-float64(Z0)), xmath.Scurve(w-float64(W0))

	ix0, iy0, iz0, iw0 := uint8(X0), uint8(Y0), uint8(Z0), uint8(W0)
	ix1, iy1, iz1, iw1 := uint8(X0+1), uint8(Y0+1), uint8(Z0+1), uint8(W0+1)

	var v [2][2][2][2]float64
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
				p := vn.perm[iy+vn.perm[iz+vn.perm[iw]]]
				v[dw][dz][dy][0] = vn.eval(ix0 + p)
				v[dw][dz][dy][1] = vn.eval(ix1 + p)
			}
		}
	}

	var wLerp [2][2][2]float64
	for dw := 0; dw < 2; dw++ {
		for dz := 0; dz < 2; dz++ {
			for dy := 0; dy < 2; dy++ {
				wLerp[dw][dz][dy] = xmath.Lerp(fx, v[dw][dz][dy][0], v[dw][dz][dy][1])
			}
		}
	}

	var zLerp [2][2]float64
	for dw := 0; dw < 2; dw++ {
		for dz := 0; dz < 2; dz++ {
			zLerp[dw][dz] = xmath.Lerp(fy, wLerp[dw][dz][0], wLerp[dw][dz][1])
		}
	}

	var yLerp [2]float64
	for dw := 0; dw < 2; dw++ {
		yLerp[dw] = xmath.Lerp(fz, zLerp[dw][0], zLerp[dw][1])
	}

	return xmath.Lerp(fw, yLerp[0], yLerp[1])
}

// Noise5D returns 5D Value Noise scaled to [-1, 1].
func (vn *ValueNoise) Noise5D(x, y, z, w, v float64) float64 {
	X0, Y0, Z0, W0, V0 := xmath.FastFloor(x), xmath.FastFloor(y), xmath.FastFloor(z), xmath.FastFloor(w), xmath.FastFloor(v)
	fx, fy, fz, fw, fv := xmath.Scurve(x-float64(X0)), xmath.Scurve(y-float64(Y0)), xmath.Scurve(z-float64(Z0)), xmath.Scurve(w-float64(W0)), xmath.Scurve(v-float64(V0))

	coords := [5][2]uint8{
		{uint8(X0), uint8(X0 + 1)},
		{uint8(Y0), uint8(Y0 + 1)},
		{uint8(Z0), uint8(Z0 + 1)},
		{uint8(W0), uint8(W0 + 1)},
		{uint8(V0), uint8(V0 + 1)},
	}

	var corners [32]float64
	for i := 0; i < 32; i++ {
		ix := coords[0][(i>>0)&1]
		iy := coords[1][(i>>1)&1]
		iz := coords[2][(i>>2)&1]
		iw := coords[3][(i>>3)&1]
		iv := coords[4][(i>>4)&1]

		hash := ix + vn.perm[iy+vn.perm[iz+vn.perm[iw+vn.perm[iv]]]]
		corners[i] = vn.eval(hash)
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

// Noise6D returns 6D Value Noise scaled to [-1, 1].
func (vn *ValueNoise) Noise6D(x, y, z, w, v, u float64) float64 {
	X0, Y0, Z0, W0, V0, U0 := xmath.FastFloor(x), xmath.FastFloor(y), xmath.FastFloor(z), xmath.FastFloor(w), xmath.FastFloor(v), xmath.FastFloor(u)
	fx, fy, fz, fw, fv, fu := xmath.Scurve(x-float64(X0)), xmath.Scurve(y-float64(Y0)), xmath.Scurve(z-float64(Z0)), xmath.Scurve(w-float64(W0)), xmath.Scurve(v-float64(V0)), xmath.Scurve(u-float64(U0))

	coords := [6][2]uint8{
		{uint8(X0), uint8(X0 + 1)},
		{uint8(Y0), uint8(Y0 + 1)},
		{uint8(Z0), uint8(Z0 + 1)},
		{uint8(W0), uint8(W0 + 1)},
		{uint8(V0), uint8(V0 + 1)},
		{uint8(U0), uint8(U0 + 1)},
	}

	var corners [64]float64
	for i := 0; i < 64; i++ {
		ix := coords[0][(i>>0)&1]
		iy := coords[1][(i>>1)&1]
		iz := coords[2][(i>>2)&1]
		iw := coords[3][(i>>3)&1]
		iv := coords[4][(i>>4)&1]
		iu := coords[5][(i>>5)&1]

		hash := ix + vn.perm[iy+vn.perm[iz+vn.perm[iw+vn.perm[iv+vn.perm[iu]]]]]
		corners[i] = vn.eval(hash)
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
