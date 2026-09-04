package interpolation

import (
	"math"

	"github.com/terefang/gommons/pkg/xrng"
)

type QuinticOutSource struct {
	source xrng.NoiseSource
}

func NewQuinticOutSource(source xrng.NoiseSource) *QuinticOutSource {
	return &QuinticOutSource{source: source}
}

// QuinticLerp interpolates between a and b using Ken Perlin's C2 quintic curve on factor t in [0, 1].
func QuinticLerp(a, b, t float64) float64 {
	return a + (t*t*t*(t*(t*6.0-15.0)+9.999998))*(b-a)
}

// QuinticInterpolate maps a value v through a smooth quintic curve.
func QuinticInterpolate(v float64) float64 {
	return QuinticLerp(0.0, 1.0, v)
}

func (q *QuinticOutSource) Noise1D(x float64) float64 {
	return QuinticInterpolate(q.source.Noise1D(x))
}

func (q *QuinticOutSource) Noise2D(x, y float64) float64 {
	return QuinticInterpolate(q.source.Noise2D(x, y))
}

func (q *QuinticOutSource) Noise3D(x, y, z float64) float64 {
	return QuinticInterpolate(q.source.Noise3D(x, y, z))
}

func (q *QuinticOutSource) Noise4D(x, y, z, w float64) float64 {
	return QuinticInterpolate(q.source.Noise4D(x, y, z, w))
}

func (q *QuinticOutSource) Noise5D(x, y, z, w, v float64) float64 {
	return QuinticInterpolate(q.source.Noise5D(x, y, z, w, v))
}

func (q *QuinticOutSource) Noise6D(x, y, z, w, v, u float64) float64 {
	return QuinticInterpolate(q.source.Noise6D(x, y, z, w, v, u))
}

type QuinticInSource struct {
	source xrng.NoiseSource
}

func NewQuinticInSource(source xrng.NoiseSource) *QuinticInSource {
	return &QuinticInSource{source: source}
}

// QuinticTransform maps an input coordinate v through Ken Perlin's quintic S-curve.
func QuinticTransform(v float64) float64 {
	if v < 0 {
		return -QuinticLerp(0.0, 1.0, math.Abs(v))
	}
	return QuinticLerp(0.0, 1.0, v)
}

func (q *QuinticInSource) Noise1D(x float64) float64 {
	return q.source.Noise1D(QuinticTransform(x))
}

func (q *QuinticInSource) Noise2D(x, y float64) float64 {
	return q.source.Noise2D(
		QuinticTransform(x),
		QuinticTransform(y),
	)
}

func (q *QuinticInSource) Noise3D(x, y, z float64) float64 {
	return q.source.Noise3D(
		QuinticTransform(x),
		QuinticTransform(y),
		QuinticTransform(z),
	)
}

func (q *QuinticInSource) Noise4D(x, y, z, w float64) float64 {
	return q.source.Noise4D(
		QuinticTransform(x),
		QuinticTransform(y),
		QuinticTransform(z),
		QuinticTransform(w),
	)
}

func (q *QuinticInSource) Noise5D(x, y, z, w, v float64) float64 {
	return q.source.Noise5D(
		QuinticTransform(x),
		QuinticTransform(y),
		QuinticTransform(z),
		QuinticTransform(w),
		QuinticTransform(v),
	)
}

func (q *QuinticInSource) Noise6D(x, y, z, w, v, u float64) float64 {
	return q.source.Noise6D(
		QuinticTransform(x),
		QuinticTransform(y),
		QuinticTransform(z),
		QuinticTransform(w),
		QuinticTransform(v),
		QuinticTransform(u),
	)
}
