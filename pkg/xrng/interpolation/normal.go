package interpolation

import "github.com/terefang/gommons/pkg/xrng"

type NormalizeSource struct {
	source xrng.NoiseSource
}

func NewNormalizeSource(source xrng.NoiseSource) *NormalizeSource {
	return &NormalizeSource{source: source}
}

func (n *NormalizeSource) Noise1D(x float64) float64 {
	return (n.source.Noise1D(x) + 1) * .5
}

func (n *NormalizeSource) Noise2D(x, y float64) float64 {
	return (n.source.Noise2D(x, y) + 1) * .5
}

func (n *NormalizeSource) Noise3D(x, y, z float64) float64 {
	return (n.source.Noise3D(x, y, z) + 1) * .5
}

func (n *NormalizeSource) Noise4D(x, y, z, w float64) float64 {
	return (n.source.Noise4D(x, y, z, w) + 1) * .5
}

func (n *NormalizeSource) Noise5D(x, y, z, w, v float64) float64 {
	return (n.source.Noise5D(x, y, z, w, v) + 1) * .5
}

func (n *NormalizeSource) Noise6D(x, y, z, w, v, u float64) float64 {
	return (n.source.Noise6D(x, y, z, w, v, u) + 1) * .5
}
