package fractal

import (
	"math"

	"github.com/terefang/gommons/pkg/xrng"
)

type MusgraveHeteroTerrainFractal struct {
	source     xrng.NoiseSource
	octaves    int
	frequency  float64
	lacunarity float64
	H          float64
	gain       float64
}

func (f *MusgraveHeteroTerrainFractal) SetH(H float64) {
	f.H = H
}

func (f *MusgraveHeteroTerrainFractal) SetGain(gain float64) {
	f.gain = gain
}

func (f *MusgraveHeteroTerrainFractal) SetFrequency(frequency float64) {
	f.frequency = frequency
}

func (f *MusgraveHeteroTerrainFractal) SetOctaves(octaves int) {
	f.octaves = octaves
}

func (f *MusgraveHeteroTerrainFractal) SetLacunarity(lacunarity float64) {
	f.lacunarity = lacunarity
}

func NewMusgraveHeteroTerrainFractal(source xrng.NoiseSource) *MusgraveHeteroTerrainFractal {
	return &MusgraveHeteroTerrainFractal{source: source, octaves: xrng.BASE_octaves, lacunarity: xrng.BASE_lacunarity, H: xrng.BASE_H, frequency: xrng.BASE_frequency, gain: xrng.BASE_gain}
}

func (f *MusgraveHeteroTerrainFractal) Noise1D(x float64) float64 {
	if f.octaves <= 0 {
		return 0
	}

	x *= f.frequency

	value := 1.
	pwr := 1.
	pwHL := math.Pow(f.lacunarity, -f.H)

	for i := 0; i <= f.octaves; i++ {
		value += f.source.Noise1D(x) * pwr * value
		pwr *= pwHL

		x *= f.lacunarity
	}
	return value / math.Pow(pwHL, float64(f.octaves))
}

func (f *MusgraveHeteroTerrainFractal) Noise2D(x, y float64) float64 {
	if f.octaves <= 0 {
		return 0
	}

	x *= f.frequency
	y *= f.frequency

	value := 1.
	pwr := 1.
	pwHL := math.Pow(f.lacunarity, -f.H)

	for i := 0; i <= f.octaves; i++ {
		value += f.source.Noise2D(x, y) * pwr * value
		pwr *= pwHL

		x *= f.lacunarity
		y *= f.lacunarity
	}
	return value / math.Pow(pwHL, float64(f.octaves))
}

func (f *MusgraveHeteroTerrainFractal) Noise3D(x, y, z float64) float64 {
	if f.octaves <= 0 {
		return 0
	}

	x *= f.frequency
	y *= f.frequency
	z *= f.frequency

	value := 1.
	pwr := 1.
	pwHL := math.Pow(f.lacunarity, -f.H)

	for i := 0; i <= f.octaves; i++ {
		value += f.source.Noise3D(x, y, z) * pwr * value
		pwr *= pwHL

		x *= f.lacunarity
		y *= f.lacunarity
		z *= f.lacunarity
	}
	return value / math.Pow(pwHL, float64(f.octaves))
}

func (f *MusgraveHeteroTerrainFractal) Noise4D(x, y, z, w float64) float64 {
	if f.octaves <= 0 {
		return 0
	}

	x *= f.frequency
	y *= f.frequency
	z *= f.frequency
	w *= f.frequency

	value := 1.
	pwr := 1.
	pwHL := math.Pow(f.lacunarity, -f.H)

	for i := 0; i <= f.octaves; i++ {
		value += f.source.Noise4D(x, y, z, w) * pwr * value
		pwr *= pwHL

		x *= f.lacunarity
		y *= f.lacunarity
		z *= f.lacunarity
		w *= f.lacunarity
	}
	return value / math.Pow(pwHL, float64(f.octaves))
}

func (f *MusgraveHeteroTerrainFractal) Noise5D(x, y, z, w, v float64) float64 {
	if f.octaves <= 0 {
		return 0
	}

	x *= f.frequency
	y *= f.frequency
	z *= f.frequency
	w *= f.frequency
	v *= f.frequency

	value := 1.
	pwr := 1.
	pwHL := math.Pow(f.lacunarity, -f.H)

	for i := 0; i <= f.octaves; i++ {
		value += f.source.Noise5D(x, y, z, w, v) * pwr * value
		pwr *= pwHL

		x *= f.lacunarity
		y *= f.lacunarity
		z *= f.lacunarity
		v *= f.lacunarity
		w *= f.lacunarity
	}
	return value / math.Pow(pwHL, float64(f.octaves))
}

func (f *MusgraveHeteroTerrainFractal) Noise6D(x, y, z, w, v, u float64) float64 {
	if f.octaves <= 0 {
		return 0
	}

	x *= f.frequency
	y *= f.frequency
	z *= f.frequency
	w *= f.frequency
	u *= f.frequency
	v *= f.frequency

	value := 1.
	pwr := 1.
	pwHL := math.Pow(f.lacunarity, -f.H)

	for i := 0; i <= f.octaves; i++ {
		value += f.source.Noise6D(x, y, z, w, v, u) * pwr * value
		pwr *= pwHL

		x *= f.lacunarity
		y *= f.lacunarity
		z *= f.lacunarity
		u *= f.lacunarity
		v *= f.lacunarity
		w *= f.lacunarity
	}
	return value / math.Pow(pwHL, float64(f.octaves))
}
