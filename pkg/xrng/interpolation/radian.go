package interpolation

import (
	"math"

	"github.com/terefang/gommons/pkg/xmath"
	"github.com/terefang/gommons/pkg/xrng"
)

type RadianOutSource struct {
	source xrng.NoiseSource
}

func NewRadianOutSource(source xrng.NoiseSource) *RadianOutSource {
	return &RadianOutSource{source: source}
}

// clamp constrains v within [minVal, maxVal].
func clamp(v, minVal, maxVal float64) float64 {
	if v < minVal {
		return minVal
	}
	if v > maxVal {
		return maxVal
	}
	return v
}

// RadianLerp shapes factor t using a sine/radian S-curve before interpolating between a and b.
func RadianLerp(a, b, t float64) float64 {
	t = (clamp(t, -1.0, 1.0) - 0.5) * math.Pi
	t = (math.Sin(t) / 2.0) + 0.5
	return xmath.Lerp(t, a, b)
}

// RadianInterpolate maps a value v through the radian S-curve.
func RadianInterpolate(v float64) float64 {
	return RadianLerp(0.0, 1.0, v)
}

func (r *RadianOutSource) Noise1D(x float64) float64 {
	return RadianInterpolate(r.source.Noise1D(x))
}

func (r *RadianOutSource) Noise2D(x, y float64) float64 {
	return RadianInterpolate(r.source.Noise2D(x, y))
}

func (r *RadianOutSource) Noise3D(x, y, z float64) float64 {
	return RadianInterpolate(r.source.Noise3D(x, y, z))
}

func (r *RadianOutSource) Noise4D(x, y, z, w float64) float64 {
	return RadianInterpolate(r.source.Noise4D(x, y, z, w))
}

func (r *RadianOutSource) Noise5D(x, y, z, w, v float64) float64 {
	return RadianInterpolate(r.source.Noise5D(x, y, z, w, v))
}

func (r *RadianOutSource) Noise6D(x, y, z, w, v, u float64) float64 {
	return RadianInterpolate(r.source.Noise6D(x, y, z, w, v, u))
}

type RadianInSource struct {
	source xrng.NoiseSource
}

func NewRadianInSource(source xrng.NoiseSource) *RadianInSource {
	return &RadianInSource{source: source}
}

// RadianTransform applies the radian S-curve directly to an input coordinate.
func RadianTransform(v float64) float64 {
	if v < 0 {
		return -RadianLerp(0.0, 1.0, math.Abs(v))
	}
	return RadianLerp(0.0, 1.0, v)
}

func (r *RadianInSource) Noise1D(x float64) float64 {
	return r.source.Noise1D(RadianTransform(x))
}

func (r *RadianInSource) Noise2D(x, y float64) float64 {
	return r.source.Noise2D(
		RadianTransform(x),
		RadianTransform(y),
	)
}

func (r *RadianInSource) Noise3D(x, y, z float64) float64 {
	return r.source.Noise3D(
		RadianTransform(x),
		RadianTransform(y),
		RadianTransform(z),
	)
}

func (r *RadianInSource) Noise4D(x, y, z, w float64) float64 {
	return r.source.Noise4D(
		RadianTransform(x),
		RadianTransform(y),
		RadianTransform(z),
		RadianTransform(w),
	)
}

func (r *RadianInSource) Noise5D(x, y, z, w, v float64) float64 {
	return r.source.Noise5D(
		RadianTransform(x),
		RadianTransform(y),
		RadianTransform(z),
		RadianTransform(w),
		RadianTransform(v),
	)
}

func (r *RadianInSource) Noise6D(x, y, z, w, v, u float64) float64 {
	return r.source.Noise6D(
		RadianTransform(x),
		RadianTransform(y),
		RadianTransform(z),
		RadianTransform(w),
		RadianTransform(v),
		RadianTransform(u),
	)
}
