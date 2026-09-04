package interpolation

import (
	"math"

	"github.com/terefang/gommons/pkg/xrng"
)

type HermiteOutSource struct {
	source xrng.NoiseSource
}

func NewHermiteOutSource(source xrng.NoiseSource) *HermiteOutSource {
	return &HermiteOutSource{source: source}
}

// HermiteLerp interpolates between a and b using a smoothstep (cubic Hermite) curve on factor t in [0, 1].
func HermiteLerp(a, b, t float64) float64 {
	return a + (t*t*(3.-2.*t))*(b-a)
}

// HermiteInterpolate maps a value v through a smooth Hermite curve.
func HermiteInterpolate(v float64) float64 {
	return HermiteLerp(0.0, 1.0, v)
}

func (h *HermiteOutSource) Noise1D(x float64) float64 {
	return HermiteInterpolate(h.source.Noise1D(x))
}

func (h *HermiteOutSource) Noise2D(x, y float64) float64 {
	return HermiteInterpolate(h.source.Noise2D(x, y))
}

func (h *HermiteOutSource) Noise3D(x, y, z float64) float64 {
	return HermiteInterpolate(h.source.Noise3D(x, y, z))
}

func (h *HermiteOutSource) Noise4D(x, y, z, w float64) float64 {
	return HermiteInterpolate(h.source.Noise4D(x, y, z, w))
}

func (h *HermiteOutSource) Noise5D(x, y, z, w, v float64) float64 {
	return HermiteInterpolate(h.source.Noise5D(x, y, z, w, v))
}

func (h *HermiteOutSource) Noise6D(x, y, z, w, v, u float64) float64 {
	return HermiteInterpolate(h.source.Noise6D(x, y, z, w, v, u))
}

type HermiteInSource struct {
	source xrng.NoiseSource
}

func NewHermiteInSource(source xrng.NoiseSource) *HermiteInSource {
	return &HermiteInSource{source: source}
}

// HermiteTransform maps an input coordinate v through a smoothstep (cubic Hermite) curve.
func HermiteTransform(v float64) float64 {
	if v < 0 {
		return -HermiteLerp(0.0, 1.0, math.Abs(v))
	}
	return HermiteLerp(0.0, 1.0, v)
}

func (h *HermiteInSource) Noise1D(x float64) float64 {
	return h.source.Noise1D(HermiteTransform(x))
}

func (h *HermiteInSource) Noise2D(x, y float64) float64 {
	return h.source.Noise2D(
		HermiteTransform(x),
		HermiteTransform(y),
	)
}

func (h *HermiteInSource) Noise3D(x, y, z float64) float64 {
	return h.source.Noise3D(
		HermiteTransform(x),
		HermiteTransform(y),
		HermiteTransform(z),
	)
}

func (h *HermiteInSource) Noise4D(x, y, z, w float64) float64 {
	return h.source.Noise4D(
		HermiteTransform(x),
		HermiteTransform(y),
		HermiteTransform(z),
		HermiteTransform(w),
	)
}

func (h *HermiteInSource) Noise5D(x, y, z, w, v float64) float64 {
	return h.source.Noise5D(
		HermiteTransform(x),
		HermiteTransform(y),
		HermiteTransform(z),
		HermiteTransform(w),
		HermiteTransform(v),
	)
}

func (h *HermiteInSource) Noise6D(x, y, z, w, v, u float64) float64 {
	return h.source.Noise6D(
		HermiteTransform(x),
		HermiteTransform(y),
		HermiteTransform(z),
		HermiteTransform(w),
		HermiteTransform(v),
		HermiteTransform(u),
	)
}
