package fractal

import (
	"math"

	"github.com/terefang/gommons/pkg/xrng"
)

type RidgedMultiFractal struct {
	source     xrng.NoiseSource
	octaves    int
	frequency  float64
	lacunarity float64
	gain       float64
}

func (f *RidgedMultiFractal) SetGain(gain float64) {
	f.gain = gain
}

func (f *RidgedMultiFractal) SetFrequency(frequency float64) {
	f.frequency = frequency
}

func (f *RidgedMultiFractal) SetOctaves(octaves int) {
	f.octaves = octaves
}

func (f *RidgedMultiFractal) SetLacunarity(lacunarity float64) {
	f.lacunarity = lacunarity
}

func NewRidgedMultiFractal(source xrng.NoiseSource) *RidgedMultiFractal {
	return &RidgedMultiFractal{source: source, octaves: xrng.BASE_octaves, lacunarity: xrng.BASE_lacunarity, frequency: xrng.BASE_frequency, gain: xrng.BASE_gain}
}

func (f *RidgedMultiFractal) Noise1D(x float64) float64 {
	if f.octaves <= 0 {
		return 0
	}

	x *= f.frequency

	sum := 0.
	exp := 1 / f.gain
	correction := 0.

	for i := 0; i <= f.octaves; i++ {
		spike := math.Abs(f.source.Noise1D(x))
		exp *= 0.5
		correction += exp
		sum += spike * exp

		x *= f.lacunarity
	}
	return (sum * 2) / (correction - 1)
}

func (f *RidgedMultiFractal) Noise2D(x, y float64) float64 {
	if f.octaves <= 0 {
		return 0
	}

	x *= f.frequency
	y *= f.frequency

	sum := 0.
	exp := 1 / f.gain
	correction := 0.

	for i := 0; i <= f.octaves; i++ {
		spike := math.Abs(f.source.Noise2D(x, y))
		exp *= 0.5
		correction += exp
		sum += spike * exp

		x *= f.lacunarity
		y *= f.lacunarity
	}
	return (sum * 2) / (correction - 1)
}

func (f *RidgedMultiFractal) Noise3D(x, y, z float64) float64 {
	if f.octaves <= 0 {
		return 0
	}

	x *= f.frequency
	y *= f.frequency
	z *= f.frequency

	sum := 0.
	exp := 1 / f.gain
	correction := 0.

	for i := 0; i <= f.octaves; i++ {
		spike := math.Abs(f.source.Noise3D(x, y, z))
		exp *= 0.5
		correction += exp
		sum += spike * exp

		x *= f.lacunarity
		y *= f.lacunarity
		z *= f.lacunarity
	}
	return (sum * 2) / (correction - 1)
}

func (f *RidgedMultiFractal) Noise4D(x, y, z, w float64) float64 {
	if f.octaves <= 0 {
		return 0
	}

	x *= f.frequency
	y *= f.frequency
	z *= f.frequency
	w *= f.frequency

	sum := 0.
	exp := 1 / f.gain
	correction := 0.

	for i := 0; i <= f.octaves; i++ {
		spike := math.Abs(f.source.Noise4D(x, y, z, w))
		exp *= 0.5
		correction += exp
		sum += spike * exp

		x *= f.lacunarity
		y *= f.lacunarity
		z *= f.lacunarity
		w *= f.lacunarity
	}
	return (sum * 2) / (correction - 1)
}

func (f *RidgedMultiFractal) Noise5D(x, y, z, w, v float64) float64 {
	if f.octaves <= 0 {
		return 0
	}

	x *= f.frequency
	y *= f.frequency
	z *= f.frequency
	w *= f.frequency
	v *= f.frequency

	sum := 0.
	exp := 1 / f.gain
	correction := 0.

	for i := 0; i <= f.octaves; i++ {
		spike := math.Abs(f.source.Noise5D(x, y, z, w, v))
		exp *= 0.5
		correction += exp
		sum += spike * exp

		x *= f.lacunarity
		y *= f.lacunarity
		z *= f.lacunarity
		v *= f.lacunarity
		w *= f.lacunarity
	}
	return (sum * 2) / (correction - 1)
}

func (f *RidgedMultiFractal) Noise6D(x, y, z, w, v, u float64) float64 {
	if f.octaves <= 0 {
		return 0
	}

	x *= f.frequency
	y *= f.frequency
	z *= f.frequency
	w *= f.frequency
	u *= f.frequency
	v *= f.frequency

	sum := 0.
	exp := 1 / f.gain
	correction := 0.

	for i := 0; i <= f.octaves; i++ {
		spike := math.Abs(f.source.Noise6D(x, y, z, w, v, u))
		exp *= 0.5
		correction += exp
		sum += spike * exp

		x *= f.lacunarity
		y *= f.lacunarity
		z *= f.lacunarity
		u *= f.lacunarity
		v *= f.lacunarity
		w *= f.lacunarity
	}
	return (sum * 2) / (correction - 1)
}
