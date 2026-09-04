package fractal

import (
	"math"

	"github.com/terefang/gommons/pkg/xrng"
)

type BillowFractal struct {
	source     xrng.NoiseSource
	octaves    int
	frequency  float64
	lacunarity float64
	gain       float64
}

func (f *BillowFractal) SetGain(gain float64) {
	f.gain = gain
}

func (f *BillowFractal) SetFrequency(frequency float64) {
	f.frequency = frequency
}

func (f *BillowFractal) SetOctaves(octaves int) {
	f.octaves = octaves
}

func (f *BillowFractal) SetLacunarity(lacunarity float64) {
	f.lacunarity = lacunarity
}

func NewBillowFractal(source xrng.NoiseSource) *BillowFractal {
	return &BillowFractal{source: source, octaves: xrng.BASE_octaves, lacunarity: xrng.BASE_H, frequency: xrng.BASE_frequency, gain: xrng.BASE_gain}
}

func (f *BillowFractal) Noise1D(x float64) float64 {
	if f.octaves <= 0 {
		return 0
	}

	x *= f.frequency

	sum := math.Abs(f.source.Noise1D(x)*2) - 1
	amp := 1.
	ampFractal := 1.

	for i := 1; i < f.octaves; i++ {
		x *= f.lacunarity

		amp *= f.gain
		ampFractal += amp
		sum += (math.Abs(f.source.Noise1D(x)*2) - 1) * amp
	}
	return sum / ampFractal
}

func (f *BillowFractal) Noise2D(x, y float64) float64 {
	if f.octaves <= 0 {
		return 0
	}

	x *= f.frequency
	y *= f.frequency

	sum := math.Abs(f.source.Noise2D(x, y)*2) - 1
	amp := 1.
	ampFractal := 1.

	for i := 1; i < f.octaves; i++ {
		x *= f.lacunarity
		y *= f.lacunarity

		amp *= f.gain
		ampFractal += amp
		sum += (math.Abs(f.source.Noise2D(x, y)*2) - 1) * amp
	}
	return sum / ampFractal
}

func (f *BillowFractal) Noise3D(x, y, z float64) float64 {
	if f.octaves <= 0 {
		return 0
	}

	x *= f.frequency
	y *= f.frequency
	z *= f.frequency

	sum := math.Abs(f.source.Noise3D(x, y, z)*2) - 1
	amp := 1.
	ampFractal := 1.

	for i := 1; i < f.octaves; i++ {
		x *= f.lacunarity
		y *= f.lacunarity
		z *= f.lacunarity

		amp *= f.gain
		ampFractal += amp
		sum += (math.Abs(f.source.Noise3D(x, y, z)*2) - 1) * amp
	}
	return sum / ampFractal
}

func (f *BillowFractal) Noise4D(x, y, z, w float64) float64 {
	if f.octaves <= 0 {
		return 0
	}

	x *= f.frequency
	y *= f.frequency
	z *= f.frequency
	w *= f.frequency

	sum := math.Abs(f.source.Noise4D(x, y, z, w)*2) - 1
	amp := 1.
	ampFractal := 1.

	for i := 1; i < f.octaves; i++ {
		x *= f.lacunarity
		y *= f.lacunarity
		z *= f.lacunarity
		w *= f.lacunarity

		amp *= f.gain
		ampFractal += amp
		sum += (math.Abs(f.source.Noise4D(x, y, z, w)*2) - 1) * amp
	}
	return sum / ampFractal
}

func (f *BillowFractal) Noise5D(x, y, z, w, v float64) float64 {
	if f.octaves <= 0 {
		return 0
	}

	x *= f.frequency
	y *= f.frequency
	z *= f.frequency
	w *= f.frequency
	v *= f.frequency

	sum := math.Abs(f.source.Noise5D(x, y, z, w, v)*2) - 1
	amp := 1.
	ampFractal := 1.

	for i := 1; i < f.octaves; i++ {
		x *= f.lacunarity
		y *= f.lacunarity
		z *= f.lacunarity
		w *= f.lacunarity
		v *= f.lacunarity

		amp *= f.gain
		ampFractal += amp
		sum += (math.Abs(f.source.Noise5D(x, y, z, w, v)*2) - 1) * amp
	}
	return sum / ampFractal
}

func (f *BillowFractal) Noise6D(x, y, z, w, v, u float64) float64 {
	if f.octaves <= 0 {
		return 0
	}

	x *= f.frequency
	y *= f.frequency
	z *= f.frequency
	w *= f.frequency
	u *= f.frequency
	v *= f.frequency

	sum := math.Abs(f.source.Noise6D(x, y, z, w, v, u)*2) - 1
	amp := 1.
	ampFractal := 1.

	for i := 1; i < f.octaves; i++ {
		x *= f.lacunarity
		y *= f.lacunarity
		z *= f.lacunarity
		w *= f.lacunarity
		u *= f.lacunarity
		v *= f.lacunarity

		amp *= f.gain
		ampFractal += amp
		sum += (math.Abs(f.source.Noise6D(x, y, z, w, v, u)*2) - 1) * amp
	}
	return sum / ampFractal
}
