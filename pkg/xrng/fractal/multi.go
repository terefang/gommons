package fractal

import (
	"github.com/terefang/gommons/pkg/xrng"
)

type MultiFractal struct {
	source     xrng.NoiseSource
	octaves    int
	frequency  float64
	lacunarity float64
	gain       float64
}

func (f *MultiFractal) SetGain(gain float64) {
	f.gain = gain
}

func (f *MultiFractal) SetFrequency(frequency float64) {
	f.frequency = frequency
}

func (f *MultiFractal) SetOctaves(octaves int) {
	f.octaves = octaves
}

func (f *MultiFractal) SetLacunarity(lacunarity float64) {
	f.lacunarity = lacunarity
}

func NewMultiFractal(source xrng.NoiseSource) *MultiFractal {
	return &MultiFractal{source: source, octaves: xrng.BASE_octaves, lacunarity: xrng.BASE_lacunarity, frequency: xrng.BASE_frequency, gain: xrng.BASE_gain}
}

func (f *MultiFractal) Noise1D(x float64) float64 {
	if f.octaves <= 0 {
		return 0
	}

	x *= f.frequency

	sum := 0.
	correction := 0.

	for i := 0; i <= f.octaves; i++ {
		correction += f.gain
		sum += f.source.Noise1D(x) * f.gain

		x *= f.lacunarity
	}
	return sum / correction
}

func (f *MultiFractal) Noise2D(x, y float64) float64 {
	if f.octaves <= 0 {
		return 0
	}

	x *= f.frequency
	y *= f.frequency

	sum := 0.
	correction := 0.

	for i := 0; i <= f.octaves; i++ {
		correction += f.gain
		sum += f.source.Noise2D(x, y) * f.gain

		x *= f.lacunarity
		y *= f.lacunarity
	}
	return sum / correction
}

func (f *MultiFractal) Noise3D(x, y, z float64) float64 {
	if f.octaves <= 0 {
		return 0
	}

	x *= f.frequency
	y *= f.frequency
	z *= f.frequency

	sum := 0.
	correction := 0.

	for i := 0; i <= f.octaves; i++ {
		correction += f.gain
		sum += f.source.Noise3D(x, y, z) * f.gain

		x *= f.lacunarity
		y *= f.lacunarity
		z *= f.lacunarity
	}
	return sum / correction
}

func (f *MultiFractal) Noise4D(x, y, z, w float64) float64 {
	if f.octaves <= 0 {
		return 0
	}

	x *= f.frequency
	y *= f.frequency
	z *= f.frequency
	w *= f.frequency

	sum := 0.
	correction := 0.

	for i := 0; i <= f.octaves; i++ {
		correction += f.gain
		sum += f.source.Noise4D(x, y, z, w) * f.gain

		x *= f.lacunarity
		y *= f.lacunarity
		z *= f.lacunarity
		w *= f.lacunarity
	}
	return sum / correction
}

func (f *MultiFractal) Noise5D(x, y, z, w, v float64) float64 {
	if f.octaves <= 0 {
		return 0
	}

	x *= f.frequency
	y *= f.frequency
	z *= f.frequency
	w *= f.frequency
	v *= f.frequency

	sum := 0.
	correction := 0.

	for i := 0; i <= f.octaves; i++ {
		correction += f.gain
		sum += f.source.Noise5D(x, y, z, w, v) * f.gain

		x *= f.lacunarity
		y *= f.lacunarity
		z *= f.lacunarity
		v *= f.lacunarity
		w *= f.lacunarity
	}
	return sum / correction
}

func (f *MultiFractal) Noise6D(x, y, z, w, v, u float64) float64 {
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
	correction := 0.

	for i := 0; i <= f.octaves; i++ {
		correction += f.gain
		sum += f.source.Noise6D(x, y, z, w, v, u) * f.gain

		x *= f.lacunarity
		y *= f.lacunarity
		z *= f.lacunarity
		u *= f.lacunarity
		v *= f.lacunarity
		w *= f.lacunarity
	}
	return sum / correction
}
