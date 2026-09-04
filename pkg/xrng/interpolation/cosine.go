package interpolation

import (
	"math"

	"github.com/terefang/gommons/pkg/xrng"
)

type CosineOutSource struct {
	source xrng.NoiseSource
}

func NewCosineOutSource(source xrng.NoiseSource) *CosineOutSource {
	return &CosineOutSource{source: source}
}

// cosineLerp interpolates between a and b using a cosine S-curve on factor t in [0, 1].
func CosineLerp(a, b, t float64) float64 {
	ft := t * math.Pi
	f := math.Cos(ft) * 0.5
	return (a * (1.0 - f)) + (b * f)
}

// CosineInterpolate maps a value v through a smooth cosine curve.
func CosineInterpolate(v float64) float64 {
	return math.Cos(v * math.Pi)
}

func (c *CosineOutSource) Noise1D(x float64) float64 {
	return CosineInterpolate(c.source.Noise1D(x))
}

func (c *CosineOutSource) Noise2D(x, y float64) float64 {
	return CosineInterpolate(c.source.Noise2D(x, y))
}

func (c *CosineOutSource) Noise3D(x, y, z float64) float64 {
	return CosineInterpolate(c.source.Noise3D(x, y, z))
}

func (c *CosineOutSource) Noise4D(x, y, z, w float64) float64 {
	return CosineInterpolate(c.source.Noise4D(x, y, z, w))
}

func (c *CosineOutSource) Noise5D(x, y, z, w, v float64) float64 {
	return CosineInterpolate(c.source.Noise5D(x, y, z, w, v))
}

func (c *CosineOutSource) Noise6D(x, y, z, w, v, u float64) float64 {
	return CosineInterpolate(c.source.Noise6D(x, y, z, w, v, u))
}
