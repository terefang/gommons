package interpolation

import (
	"math"

	"github.com/terefang/gommons/pkg/xrng"
)

type SineOutSource struct {
	source xrng.NoiseSource
}

func NewSineOutSource(source xrng.NoiseSource) *SineOutSource {
	return &SineOutSource{source: source}
}

// cosineLerp interpolates between a and b using a sine S-curve on factor t in [0, 1].
func SineLerp(a, b, t float64) float64 {
	ft := t * math.Pi
	f := (1.0 - math.Sin(ft)) * 0.5
	return (a * (1.0 - f)) + (b * f)
}

// SineInterpolate maps a value v through a smooth cosine curve.
func SineInterpolate(v float64) float64 {
	return math.Sin(v * math.Pi)
}

func (c *SineOutSource) Noise1D(x float64) float64 {
	return SineInterpolate(c.source.Noise1D(x))
}

func (c *SineOutSource) Noise2D(x, y float64) float64 {
	return SineInterpolate(c.source.Noise2D(x, y))
}

func (c *SineOutSource) Noise3D(x, y, z float64) float64 {
	return SineInterpolate(c.source.Noise3D(x, y, z))
}

func (c *SineOutSource) Noise4D(x, y, z, w float64) float64 {
	return SineInterpolate(c.source.Noise4D(x, y, z, w))
}

func (c *SineOutSource) Noise5D(x, y, z, w, v float64) float64 {
	return SineInterpolate(c.source.Noise5D(x, y, z, w, v))
}

func (c *SineOutSource) Noise6D(x, y, z, w, v, u float64) float64 {
	return SineInterpolate(c.source.Noise6D(x, y, z, w, v, u))
}
