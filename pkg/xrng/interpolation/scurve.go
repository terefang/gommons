package interpolation

import (
	"github.com/terefang/gommons/pkg/xmath"
	"github.com/terefang/gommons/pkg/xrng"
)

// SCurveFilter wraps a NoiseSource and applies the S-curve transformation.
type SCurveFilter struct {
	Source xrng.NoiseSource
}

func NewSCurveFilter(source xrng.NoiseSource) *SCurveFilter {
	return &SCurveFilter{Source: source}
}

func (s *SCurveFilter) Noise1D(x float64) float64 {
	return xmath.SCurveContrast(s.Source.Noise1D(x))
}

func (s *SCurveFilter) Noise2D(x, y float64) float64 {
	return xmath.SCurveContrast(s.Source.Noise2D(x, y))
}

func (s *SCurveFilter) Noise3D(x, y, z float64) float64 {
	return xmath.SCurveContrast(s.Source.Noise3D(x, y, z))
}

func (s *SCurveFilter) Noise4D(x, y, z, w float64) float64 {
	return xmath.SCurveContrast(s.Source.Noise4D(x, y, z, w))
}

func (s *SCurveFilter) Noise5D(x, y, z, w, v float64) float64 {
	return xmath.SCurveContrast(s.Source.Noise5D(x, y, z, w, v))
}

func (s *SCurveFilter) Noise6D(x, y, z, w, v, u float64) float64 {
	return xmath.SCurveContrast(s.Source.Noise6D(x, y, z, w, v, u))
}
