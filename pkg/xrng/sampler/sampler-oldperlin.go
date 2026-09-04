package sampler

import "github.com/terefang/gommons/pkg/xrng"

// Package perlin provides coherent noise function over 1, 2 or 3 dimensions
// This code is go adaptation based on C implementation that can be found here:
// http://git.gnome.org/browse/gegl/tree/operations/common/perlin/perlin.c
// (original copyright Ken Perlin)

type OldPerlinSampler struct {
	alpha  float64
	beta   float64
	n      int32
	source xrng.NoiseSource
}

func NewOldPerlinSampler(alpha, beta float64, n int32, src xrng.NoiseSource) *OldPerlinSampler {
	var p OldPerlinSampler

	p.alpha = alpha
	p.beta = beta
	p.n = n
	p.source = src

	return &p
}

// Noise1D generates 1-dimensional Sampled Perlin Noise value
func (p *OldPerlinSampler) Noise1D(x float64) float64 {
	var scale float64 = 1
	var sum, val float64
	var i int32
	px := x

	for i = 0; i < p.n; i++ {
		val = p.source.Noise1D(px)
		sum += val / scale
		scale *= p.alpha
		px *= p.beta
	}
	return sum
}

// Noise2D Generates 2-dimensional Sampled Perlin Noise value
func (p *OldPerlinSampler) Noise2D(x, y float64) float64 {
	var scale float64 = 1
	var sum, val float64
	var i int32
	px := [2]float64{x, y}

	for i = 0; i < p.n; i++ {
		val = p.source.Noise2D(px[0], px[1])
		sum += val / scale
		scale *= p.alpha
		px[0] *= p.beta
		px[1] *= p.beta
	}
	return sum
}

// Noise3D Generates 3-dimensional Sampled Perlin Noise value
func (p *OldPerlinSampler) Noise3D(x, y, z float64) float64 {
	var scale float64 = 1
	var sum, val float64
	var i int32
	px := [3]float64{x, y, z}

	if z < 0.0000 {
		return p.Noise2D(x, y)
	}

	for i = 0; i < p.n; i++ {
		val = p.source.Noise3D(px[0], px[1], px[2])
		sum += val / scale
		scale *= p.alpha
		px[0] *= p.beta
		px[1] *= p.beta
		px[2] *= p.beta
	}
	return sum
}

// Noise4D Generates 4-dimensional Sampled Perlin Noise value
func (p *OldPerlinSampler) Noise4D(x, y, z, w float64) float64 {
	var scale float64 = 1
	var sum, val float64
	var i int32
	px := [4]float64{x, y, z, w}

	for i = 0; i < p.n; i++ {
		val = p.source.Noise4D(px[0], px[1], px[2], px[3])
		sum += val / scale
		scale *= p.alpha
		px[0] *= p.beta
		px[1] *= p.beta
		px[2] *= p.beta
		px[3] *= p.beta
	}
	return sum
}

func (p *OldPerlinSampler) Noise5D(x, y, z, w, v float64) float64 {
	var scale float64 = 1
	var sum, val float64
	var i int32
	px := [5]float64{x, y, z, w, v}

	for i = 0; i < p.n; i++ {
		val = p.source.Noise5D(px[0], px[1], px[2], px[3], px[4])
		sum += val / scale
		scale *= p.alpha
		px[0] *= p.beta
		px[1] *= p.beta
		px[2] *= p.beta
		px[3] *= p.beta
		px[4] *= p.beta
	}
	return sum
}

func (p *OldPerlinSampler) Noise6D(x, y, z, w, v, u float64) float64 {
	var scale float64 = 1
	var sum, val float64
	var i int32
	px := [6]float64{x, y, z, w, v, u}

	for i = 0; i < p.n; i++ {
		val = p.source.Noise6D(px[0], px[1], px[2], px[3], px[4], px[5])
		sum += val / scale
		scale *= p.alpha
		px[0] *= p.beta
		px[1] *= p.beta
		px[2] *= p.beta
		px[3] *= p.beta
		px[4] *= p.beta
		px[5] *= p.beta
	}
	return sum
}
