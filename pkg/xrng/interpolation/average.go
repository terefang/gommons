package interpolation

import (
	"github.com/terefang/gommons/pkg/xrng"
)

// AverageNoise aggregates multiple NoiseSources and returns their arithmetic mean.
type AverageNoise struct {
	sources []xrng.NoiseSource
}

// NewAverageNoise creates a new AverageNoise compositor.
// Returns an error if fewer than two sources are provided.
func NewAverageNoise(sources ...xrng.NoiseSource) *AverageNoise {
	return &AverageNoise{sources: sources}
}

// Noise1D evaluates and averages all sources at (x).
func (a *AverageNoise) Noise1D(x float64) float64 {
	var sum float64
	for _, source := range a.sources {
		sum += source.Noise1D(x)
	}
	return sum / float64(len(a.sources))
}

// Noise2D evaluates and averages all sources at (x, y).
func (a *AverageNoise) Noise2D(x, y float64) float64 {
	var sum float64
	for _, source := range a.sources {
		sum += source.Noise2D(x, y)
	}
	return sum / float64(len(a.sources))
}

// Noise3D evaluates and averages all sources at (x, y, z).
func (a *AverageNoise) Noise3D(x, y, z float64) float64 {
	var sum float64
	for _, source := range a.sources {
		sum += source.Noise3D(x, y, z)
	}
	return sum / float64(len(a.sources))
}

// Noise4D evaluates and averages all sources at (x, y, z, w).
func (a *AverageNoise) Noise4D(x, y, z, w float64) float64 {
	var sum float64
	for _, source := range a.sources {
		sum += source.Noise4D(x, y, z, w)
	}
	return sum / float64(len(a.sources))
}

func (a *AverageNoise) Noise5D(x, y, z, w, v float64) float64 {
	var sum float64
	for _, source := range a.sources {
		sum += source.Noise5D(x, y, z, w, v)
	}
	return sum / float64(len(a.sources))
}

func (a *AverageNoise) Noise6D(x, y, z, w, v, u float64) float64 {
	var sum float64
	for _, source := range a.sources {
		sum += source.Noise6D(x, y, z, w, v, u)
	}
	return sum / float64(len(a.sources))
}
