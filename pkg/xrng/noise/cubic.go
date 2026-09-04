package noise

import (
	"math/rand"

	"github.com/terefang/gommons/pkg/xmath"
)

// CubicNoise implements multi-dimensional Cubic Spline Noise.
type CubicNoise struct {
	perm [512]uint8
}

// NewCubicNoise creates a new CubicNoise generator initialized with a seed.
func NewCubicNoise(seed int64) *CubicNoise {
	return NewCubicNoiseWithSource(rand.NewSource(seed))
}

// NewCubicNoiseWithSource creates a new CubicNoise generator using a rand.Source.
func NewCubicNoiseWithSource(source rand.Source) *CubicNoise {
	var p [256]uint8
	for i := range p {
		p[i] = uint8(i)
	}
	r := rand.New(source)
	for i := 255; i > 0; i-- {
		j := r.Intn(i + 1)
		p[i], p[j] = p[j], p[i]
	}

	cn := &CubicNoise{}
	for i := 0; i < 512; i++ {
		cn.perm[i] = p[i&0xff]
	}
	return cn
}

// cubicInterpolate performs 1D cubic spline (Catmull-Rom) interpolation across 4 sample points.
func cubicInterpolate(p0, p1, p2, p3, t float64) float64 {
	a := -0.5*p0 + 1.5*p1 - 1.5*p2 + 0.5*p3
	b := p0 - 2.5*p1 + 2.0*p2 - 0.5*p3
	c := -0.5*p0 + 0.5*p2
	d := p1
	return a*t*t*t + b*t*t + c*t + d
}

// val maps a permutation hash to a normalized pseudo-random value in [-1, 1].
func val(hash uint8) float64 {
	return (float64(hash) / 127.5) - 1.0
}

// Noise1D returns 1D Cubic Noise scaled to [-1, 1].
func (cn *CubicNoise) Noise1D(x float64) float64 {
	xi := xmath.FastFloor(x)
	t := x - float64(xi)

	ix := uint8(xi)
	p0 := val(cn.perm[ix-1])
	p1 := val(cn.perm[ix])
	p2 := val(cn.perm[ix+1])
	p3 := val(cn.perm[ix+2])

	return cubicInterpolate(p0, p1, p2, p3, t)
}

// Noise2D returns 2D Cubic Noise scaled to [-1, 1].
func (cn *CubicNoise) Noise2D(x, y float64) float64 {
	xi, yi := xmath.FastFloor(x), xmath.FastFloor(y)
	tx, ty := x-float64(xi), y-float64(yi)

	ix, iy := uint8(xi), uint8(yi)

	var arrY [4]float64
	for j := 0; j < 4; j++ {
		py := iy + uint8(j-1)
		p0 := val(cn.perm[ix-1+cn.perm[py]])
		p1 := val(cn.perm[ix+cn.perm[py]])
		p2 := val(cn.perm[ix+1+cn.perm[py]])
		p3 := val(cn.perm[ix+2+cn.perm[py]])
		arrY[j] = cubicInterpolate(p0, p1, p2, p3, tx)
	}

	return cubicInterpolate(arrY[0], arrY[1], arrY[2], arrY[3], ty)
}

// Noise3D returns 3D Cubic Noise scaled to [-1, 1].
func (cn *CubicNoise) Noise3D(x, y, z float64) float64 {
	xi, yi, zi := xmath.FastFloor(x), xmath.FastFloor(y), xmath.FastFloor(z)
	tx, ty, tz := x-float64(xi), y-float64(yi), z-float64(zi)

	ix, iy, iz := uint8(xi), uint8(yi), uint8(zi)

	var arrZ [4]float64
	for k := 0; k < 4; k++ {
		pz := iz + uint8(k-1)
		var arrY [4]float64
		for j := 0; j < 4; j++ {
			py := iy + uint8(j-1)
			h := cn.perm[py+cn.perm[pz]]
			p0 := val(cn.perm[ix-1+h])
			p1 := val(cn.perm[ix+h])
			p2 := val(cn.perm[ix+1+h])
			p3 := val(cn.perm[ix+2+h])
			arrY[j] = cubicInterpolate(p0, p1, p2, p3, tx)
		}
		arrZ[k] = cubicInterpolate(arrY[0], arrY[1], arrY[2], arrY[3], ty)
	}

	return cubicInterpolate(arrZ[0], arrZ[1], arrZ[2], arrZ[3], tz)
}

// Noise4D returns 4D Cubic Noise scaled to [-1, 1].
func (cn *CubicNoise) Noise4D(x, y, z, w float64) float64 {
	xi, yi, zi, wi := xmath.FastFloor(x), xmath.FastFloor(y), xmath.FastFloor(z), xmath.FastFloor(w)
	tx, ty, tz, tw := x-float64(xi), y-float64(yi), z-float64(zi), w-float64(wi)

	ix, iy, iz, iw := uint8(xi), uint8(yi), uint8(zi), uint8(wi)

	var arrW [4]float64
	for l := 0; l < 4; l++ {
		pw := iw + uint8(l-1)
		var arrZ [4]float64
		for k := 0; k < 4; k++ {
			pz := iz + uint8(k-1)
			var arrY [4]float64
			for j := 0; j < 4; j++ {
				py := iy + uint8(j-1)
				h := cn.perm[py+cn.perm[pz+cn.perm[pw]]]
				p0 := val(cn.perm[ix-1+h])
				p1 := val(cn.perm[ix+h])
				p2 := val(cn.perm[ix+1+h])
				p3 := val(cn.perm[ix+2+h])
				arrY[j] = cubicInterpolate(p0, p1, p2, p3, tx)
			}
			arrZ[k] = cubicInterpolate(arrY[0], arrY[1], arrY[2], arrY[3], ty)
		}
		arrW[l] = cubicInterpolate(arrZ[0], arrZ[1], arrZ[2], arrZ[3], tz)
	}

	return cubicInterpolate(arrW[0], arrW[1], arrW[2], arrW[3], tw)
}

// Noise5D returns 5D Cubic Noise scaled to [-1, 1].
func (cn *CubicNoise) Noise5D(x, y, z, w, v float64) float64 {
	xi, yi, zi, wi, vi := xmath.FastFloor(x), xmath.FastFloor(y), xmath.FastFloor(z), xmath.FastFloor(w), xmath.FastFloor(v)
	tx, ty, tz, tw, tv := x-float64(xi), y-float64(yi), z-float64(zi), w-float64(wi), v-float64(vi)

	ix, iy, iz, iw, iv := uint8(xi), uint8(yi), uint8(zi), uint8(wi), uint8(vi)

	var arrV [4]float64
	for m := 0; m < 4; m++ {
		pv := iv + uint8(m-1)
		var arrW [4]float64
		for l := 0; l < 4; l++ {
			pw := iw + uint8(l-1)
			var arrZ [4]float64
			for k := 0; k < 4; k++ {
				pz := iz + uint8(k-1)
				var arrY [4]float64
				for j := 0; j < 4; j++ {
					py := iy + uint8(j-1)
					h := cn.perm[py+cn.perm[pz+cn.perm[pw+cn.perm[pv]]]]
					p0 := val(cn.perm[ix-1+h])
					p1 := val(cn.perm[ix+h])
					p2 := val(cn.perm[ix+1+h])
					p3 := val(cn.perm[ix+2+h])
					arrY[j] = cubicInterpolate(p0, p1, p2, p3, tx)
				}
				arrZ[k] = cubicInterpolate(arrY[0], arrY[1], arrY[2], arrY[3], ty)
			}
			arrW[l] = cubicInterpolate(arrZ[0], arrZ[1], arrZ[2], arrZ[3], tz)
		}
		arrV[m] = cubicInterpolate(arrW[0], arrW[1], arrW[2], arrW[3], tw)
	}

	return cubicInterpolate(arrV[0], arrV[1], arrV[2], arrV[3], tv)
}

// Noise6D returns 6D Cubic Noise scaled to [-1, 1].
func (cn *CubicNoise) Noise6D(x, y, z, w, v, u float64) float64 {
	xi, yi, zi, wi, vi, ui := xmath.FastFloor(x), xmath.FastFloor(y), xmath.FastFloor(z), xmath.FastFloor(w), xmath.FastFloor(v), xmath.FastFloor(u)
	tx, ty, tz, tw, tv, tu := x-float64(xi), y-float64(yi), z-float64(zi), w-float64(wi), v-float64(vi), u-float64(ui)

	ix, iy, iz, iw, iv, iu := uint8(xi), uint8(yi), uint8(zi), uint8(wi), uint8(vi), uint8(ui)

	var arrU [4]float64
	for n := 0; n < 4; n++ {
		pu := iu + uint8(n-1)
		var arrV [4]float64
		for m := 0; m < 4; m++ {
			pv := iv + uint8(m-1)
			var arrW [4]float64
			for l := 0; l < 4; l++ {
				pw := iw + uint8(l-1)
				var arrZ [4]float64
				for k := 0; k < 4; k++ {
					pz := iz + uint8(k-1)
					var arrY [4]float64
					for j := 0; j < 4; j++ {
						py := iy + uint8(j-1)
						h := cn.perm[py+cn.perm[pz+cn.perm[pw+cn.perm[pv+cn.perm[pu]]]]]
						p0 := val(cn.perm[ix-1+h])
						p1 := val(cn.perm[ix+h])
						p2 := val(cn.perm[ix+1+h])
						p3 := val(cn.perm[ix+2+h])
						arrY[j] = cubicInterpolate(p0, p1, p2, p3, tx)
					}
					arrZ[k] = cubicInterpolate(arrY[0], arrY[1], arrY[2], arrY[3], ty)
				}
				arrW[l] = cubicInterpolate(arrZ[0], arrZ[1], arrZ[2], arrZ[3], tz)
			}
			arrV[m] = cubicInterpolate(arrW[0], arrW[1], arrW[2], arrW[3], tw)
		}
		arrU[n] = cubicInterpolate(arrV[0], arrV[1], arrV[2], arrV[3], tv)
	}

	return cubicInterpolate(arrU[0], arrU[1], arrU[2], arrU[3], tu)
}
