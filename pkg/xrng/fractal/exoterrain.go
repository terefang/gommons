package fractal

import (
	"math"

	"github.com/terefang/gommons/pkg/xrng"
	"github.com/terefang/gommons/pkg/xrng/noise"
)

type ExoTerrainFractal struct {
	source     xrng.NoiseSource
	octaves    int
	frequency  float64
	lacunarity float64
	H          float64
	offset     float64
	gain       float64
}

func (f *ExoTerrainFractal) SetH(H float64) {
	f.H = H
}

func (f *ExoTerrainFractal) SetGain(gain float64) {
	f.gain = gain
}

func (f *ExoTerrainFractal) SetFrequency(frequency float64) {
	f.frequency = frequency
}

func (f *ExoTerrainFractal) SetOctaves(octaves int) {
	f.octaves = octaves
}

func (f *ExoTerrainFractal) SetLacunarity(lacunarity float64) {
	f.lacunarity = lacunarity
}

func NewExoTerrainFractal(source xrng.NoiseSource) *ExoTerrainFractal {
	return &ExoTerrainFractal{source: source, octaves: xrng.BASE_octaves, lacunarity: xrng.BASE_lacunarity, H: xrng.BASE_H, frequency: xrng.BASE_frequency, gain: xrng.BASE_gain}
}

func (f *ExoTerrainFractal) Noise1D(x float64) float64 {
	if f.octaves <= 0 {
		return 0
	}

	x *= f.frequency

	value := f.source.Noise1D(x)
	value = noise.Clamp01(value)

	_pwr := f.gain
	_pwrHL := math.Pow(f.lacunarity, -f.H)

	_distor1 := f.source.Noise1D(x / 4.)
	_striation1 := _distor1 * 8.
	_noise1 := f.source.Noise1D(x + _striation1 + _distor1)
	for i := 1; i <= f.octaves; i++ {
		_distor2 := f.source.Noise1D(x / 8.)
		_striation2 := _distor2 * 8.
		_noise2 := f.source.Noise1D((x/2.)+_striation2+_distor2) * 1.5
		_roughness := f.source.Noise1D(x/6.) - .3
		_bumpdistort := f.source.Noise1D(x / .2)
		_bumpnoise := f.source.Noise1D((x / .5) + 2.*_bumpdistort)
		_noise1 += (_noise2*_noise2*_noise2*_noise2 + _roughness*_bumpnoise + f.offset) * _pwr

		x *= f.lacunarity
		_pwr /= -_pwrHL
	}
	return _noise1
}

func (f *ExoTerrainFractal) Noise2D(x, y float64) float64 {
	if f.octaves <= 0 {
		return 0
	}

	x *= f.frequency
	y *= f.frequency

	value := f.source.Noise2D(x, y)
	value = noise.Clamp01(value)

	_pwr := f.gain
	_pwrHL := math.Pow(f.lacunarity, -f.H)

	_distor1 := f.source.Noise2D(x/4., y/4.)
	_striation1 := _distor1 * 8.
	_noise1 := f.source.Noise2D(x, y+_striation1+_distor1)
	for i := 1; i <= f.octaves; i++ {
		_distor2 := f.source.Noise2D(x/8., y/8.)
		_striation2 := _distor2 * 8.
		_noise2 := f.source.Noise2D(x/2., (y/2.)+_striation2+_distor2) * 1.5
		_roughness := f.source.Noise2D(x/6., y/6.) - .3
		_bumpdistort := f.source.Noise2D(x/.2, y/.2)
		_bumpnoise := f.source.Noise2D(x/.5, (y/.5)+2.*_bumpdistort)
		_noise1 += (_noise2*_noise2*_noise2*_noise2 + _roughness*_bumpnoise + f.offset) * _pwr

		x *= f.lacunarity
		y *= f.lacunarity
		_pwr /= -_pwrHL
	}
	return _noise1
}

func (f *ExoTerrainFractal) Noise3D(x, y, z float64) float64 {
	if f.octaves <= 0 {
		return 0
	}

	x *= f.frequency
	y *= f.frequency
	z *= f.frequency

	value := f.source.Noise3D(x, y, z)
	value = noise.Clamp01(value)

	_pwr := f.gain
	_pwrHL := math.Pow(f.lacunarity, -f.H)

	_distor1 := f.source.Noise3D(x/4., y/4., z/4.)
	_striation1 := _distor1 * 8.
	_noise1 := f.source.Noise3D(x, y, z+_striation1+_distor1)
	for i := 1; i <= f.octaves; i++ {
		_distor2 := f.source.Noise3D(x/8., y/8., z/8.)
		_striation2 := _distor2 * 8.
		_noise2 := f.source.Noise3D(x/2., y/2., (z/2.)+_striation2+_distor2) * 1.5
		_roughness := f.source.Noise3D(x/6., y/6., z/6.) - .3
		_bumpdistort := f.source.Noise3D(x/.2, y/.2, z/.2)
		_bumpnoise := f.source.Noise3D(x/.5, y/.5, (z/.5)+2.*_bumpdistort)
		_noise1 += (_noise2*_noise2*_noise2*_noise2 + _roughness*_bumpnoise + f.offset) * _pwr

		x *= f.lacunarity
		y *= f.lacunarity
		z *= f.lacunarity
		_pwr /= -_pwrHL
	}
	return _noise1
}

func (f *ExoTerrainFractal) Noise4D(x, y, z, w float64) float64 {
	if f.octaves <= 0 {
		return 0
	}

	x *= f.frequency
	y *= f.frequency
	z *= f.frequency
	w *= f.frequency

	value := f.source.Noise4D(x, y, z, w)
	value = noise.Clamp01(value)

	_pwr := f.gain
	_pwrHL := math.Pow(f.lacunarity, -f.H)

	_distor1 := f.source.Noise4D(x/4., y/4., z/4., w/4.)
	_striation1 := _distor1 * 8.
	_noise1 := f.source.Noise4D(x, y, z, w+_striation1+_distor1)
	for i := 1; i <= f.octaves; i++ {
		_distor2 := f.source.Noise4D(x/8., y/8., z/8., w/8.)
		_striation2 := _distor2 * 8.
		_noise2 := f.source.Noise4D(x/2., y/2., z/2., (w/2.)+_striation2+_distor2) * 1.5
		_roughness := f.source.Noise4D(x/6., y/6., z/6., w/6.) - .3
		_bumpdistort := f.source.Noise4D(x/.2, y/.2, z/.2, w/.2)
		_bumpnoise := f.source.Noise4D(x/.5, y/.5, z/.5, (w/.5)+2.*_bumpdistort)
		_noise1 += (_noise2*_noise2*_noise2*_noise2 + _roughness*_bumpnoise + f.offset) * _pwr

		x *= f.lacunarity
		y *= f.lacunarity
		z *= f.lacunarity
		_pwr /= -_pwrHL
	}
	return _noise1
}

func (f *ExoTerrainFractal) Noise5D(x, y, z, w, v float64) float64 {
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

	_pwr := f.gain
	_pwrHL := math.Pow(f.lacunarity, -f.H)

	_distor1 := f.source.Noise5D(x/4., y/4., z/4., w/4., v/4.)
	_striation1 := _distor1 * 8.
	_noise1 := f.source.Noise5D(x, y, z, w, v+_striation1+_distor1)
	for i := 1; i <= f.octaves; i++ {
		_distor2 := f.source.Noise5D(x/8., y/8., z/8., w/8., v/8.)
		_striation2 := _distor2 * 8.
		_noise2 := f.source.Noise5D(x/2., y/2., z/2., w/2., (v/2.)+_striation2+_distor2) * 1.5
		_roughness := f.source.Noise5D(x/6., y/6., z/6., w/6., v/6.) - .3
		_bumpdistort := f.source.Noise5D(x/.2, y/.2, z/.2, w/.2, v/.2)
		_bumpnoise := f.source.Noise5D(x/.5, y/.5, z/.5, w/.5, (v/.5)+2.*_bumpdistort)
		_noise1 += (_noise2*_noise2*_noise2*_noise2 + _roughness*_bumpnoise + f.offset) * _pwr

		x *= f.lacunarity
		y *= f.lacunarity
		z *= f.lacunarity
		v *= f.lacunarity
		w *= f.lacunarity
		_pwr /= -_pwrHL
	}
	return _noise1
}

func (f *ExoTerrainFractal) Noise6D(x, y, z, w, v, u float64) float64 {
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

	_pwr := f.gain
	_pwrHL := math.Pow(f.lacunarity, -f.H)

	_distor1 := f.source.Noise6D(x/4., y/4., z/4., w/4., v/4., u/4.)
	_striation1 := _distor1 * 8.
	_noise1 := f.source.Noise6D(x, y, z, w, v, u+_striation1+_distor1)
	for i := 1; i <= f.octaves; i++ {
		_distor2 := f.source.Noise6D(x/8., y/8., z/8., w/8., v/8., u/8.)
		_striation2 := _distor2 * 8.
		_noise2 := f.source.Noise6D(x/2., y/2., z/2., w/2., v/2., (u/2.)+_striation2+_distor2) * 1.5
		_roughness := f.source.Noise6D(x/6., y/6., z/6., w/6., v/6., u/6.) - .3
		_bumpdistort := f.source.Noise6D(x/.2, y/.2, z/.2, w/.2, v/.2, u/.2)
		_bumpnoise := f.source.Noise6D(x/.5, y/.5, z/.5, w/.5, v/.5, (u/.5)+2.*_bumpdistort)
		_noise1 += (_noise2*_noise2*_noise2*_noise2 + _roughness*_bumpnoise + f.offset) * _pwr

		x *= f.lacunarity
		y *= f.lacunarity
		z *= f.lacunarity
		u *= f.lacunarity
		v *= f.lacunarity
		w *= f.lacunarity
		_pwr /= -_pwrHL
	}
	return _noise1
}
