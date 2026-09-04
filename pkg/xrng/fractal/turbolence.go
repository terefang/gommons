package fractal

import (
	"github.com/terefang/gommons/pkg/xrng"
	"github.com/terefang/gommons/pkg/xrng/noise"
)

type TurbolenceFractal struct {
	source     xrng.NoiseSource
	octaves    int
	frequency  float64
	lacunarity float64
	gain       float64
}

func (f *TurbolenceFractal) SetGain(gain float64) {
	f.gain = gain
}

func (f *TurbolenceFractal) SetFrequency(frequency float64) {
	f.frequency = frequency
}

func (f *TurbolenceFractal) SetOctaves(octaves int) {
	f.octaves = octaves
}

func (f *TurbolenceFractal) SetLacunarity(lacunarity float64) {
	f.lacunarity = lacunarity
}

func NewTurbolenceFractal(source xrng.NoiseSource) *TurbolenceFractal {
	return &TurbolenceFractal{source: source, octaves: xrng.BASE_octaves, lacunarity: xrng.BASE_H, frequency: xrng.BASE_frequency, gain: xrng.BASE_gain}
}

func (f *TurbolenceFractal) Noise1D(x float64) float64 {
	if f.octaves <= 0 {
		return 0
	}

	x *= f.frequency

	value := f.source.Noise1D(x)
	value = noise.Clamp01(value)

	_o := f.gain
	_l := f.lacunarity

	for i := 1; i <= f.octaves; i++ {
		value += f.source.Noise1D(x*_l) * _o
		_l *= f.lacunarity
		_o *= f.gain
	}
	return value
}

func (f *TurbolenceFractal) Noise2D(x, y float64) float64 {
	if f.octaves <= 0 {
		return 0
	}

	x *= f.frequency
	y *= f.frequency

	value := f.source.Noise2D(x, y)
	value = noise.Clamp01(value)

	_o := f.gain
	_l := f.lacunarity

	for i := 1; i <= f.octaves; i++ {
		value += f.source.Noise2D(x*_l, y*_l) * _o
		_l *= f.lacunarity
		_o *= f.gain
	}
	return value
}

func (f *TurbolenceFractal) Noise3D(x, y, z float64) float64 {
	if f.octaves <= 0 {
		return 0
	}

	x *= f.frequency
	y *= f.frequency
	z *= f.frequency

	value := f.source.Noise3D(x, y, z)
	value = noise.Clamp01(value)

	_o := f.gain
	_l := f.lacunarity

	for i := 1; i <= f.octaves; i++ {
		value += f.source.Noise3D(x*_l, y*_l, z*_l) * _o
		_l *= f.lacunarity
		_o *= f.gain
	}
	return value
}

func (f *TurbolenceFractal) Noise4D(x, y, z, w float64) float64 {
	if f.octaves <= 0 {
		return 0
	}

	x *= f.frequency
	y *= f.frequency
	z *= f.frequency
	w *= f.frequency

	value := f.source.Noise4D(x, y, z, w)
	value = noise.Clamp01(value)

	_o := f.gain
	_l := f.lacunarity

	for i := 1; i <= f.octaves; i++ {
		value += f.source.Noise4D(x*_l, y*_l, z*_l, w*_l) * _o
		_l *= f.lacunarity
		_o *= f.gain
	}
	return value
}

func (f *TurbolenceFractal) Noise5D(x, y, z, w, v float64) float64 {
	if f.octaves <= 0 {
		return 0
	}

	x *= f.frequency
	y *= f.frequency
	z *= f.frequency
	w *= f.frequency
	v *= f.frequency

	value := f.source.Noise5D(x, y, z, w, v)
	value = noise.Clamp01(value)

	_o := f.gain
	_l := f.lacunarity

	for i := 1; i <= f.octaves; i++ {
		value += f.source.Noise5D(x*_l, y*_l, z*_l, w*_l, v*_l) * _o
		_l *= f.lacunarity
		_o *= f.gain
	}
	return value
}

func (f *TurbolenceFractal) Noise6D(x, y, z, w, v, u float64) float64 {
	if f.octaves <= 0 {
		return 0
	}

	x *= f.frequency
	y *= f.frequency
	z *= f.frequency
	w *= f.frequency
	u *= f.frequency
	v *= f.frequency

	value := f.source.Noise6D(x, y, z, w, v, u)
	value = noise.Clamp01(value)

	_o := f.gain
	_l := f.lacunarity

	for i := 1; i <= f.octaves; i++ {
		value += f.source.Noise6D(x*_l, y*_l, z*_l, w*_l, v*_l, u*_l) * _o
		_l *= f.lacunarity
		_o *= f.gain
	}
	return value
}
