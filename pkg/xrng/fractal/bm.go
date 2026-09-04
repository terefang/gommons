package fractal

import "github.com/terefang/gommons/pkg/xrng"

type BrownianMotionFractal struct {
	source     xrng.NoiseSource
	octaves    int
	frequency  float64
	lacunarity float64
	gain       float64
}

func (f *BrownianMotionFractal) SetGain(gain float64) {
	f.gain = gain
}

func (f *BrownianMotionFractal) SetFrequency(frequency float64) {
	f.frequency = frequency
}

func (f *BrownianMotionFractal) SetOctaves(octaves int) {
	f.octaves = octaves
}

func (f *BrownianMotionFractal) SetLacunarity(lacunarity float64) {
	f.lacunarity = lacunarity
}

func NewBrownianMotionFractal(source xrng.NoiseSource) *BrownianMotionFractal {
	return &BrownianMotionFractal{source: source, octaves: xrng.BASE_octaves, lacunarity: xrng.BASE_H, frequency: xrng.BASE_frequency, gain: xrng.BASE_gain}
}

func (f *BrownianMotionFractal) Noise1D(x float64) float64 {
	if f.octaves <= 0 {
		return 0
	}

	x *= f.frequency

	sum := f.source.Noise1D(x)
	amp := 1.
	ampFractal := 1.

	for i := 0; i < f.octaves; i++ {
		x *= f.lacunarity

		amp *= f.gain
		ampFractal += amp
		sum += f.source.Noise1D(x) * amp
	}
	return sum / ampFractal
}

func (f *BrownianMotionFractal) Noise2D(x, y float64) float64 {
	if f.octaves <= 0 {
		return 0
	}

	x *= f.frequency
	y *= f.frequency

	sum := f.source.Noise2D(x, y)
	amp := 1.
	ampFractal := 1.

	for i := 0; i < f.octaves; i++ {
		x *= f.lacunarity
		y *= f.lacunarity

		amp *= f.gain
		ampFractal += amp
		sum += f.source.Noise2D(x, y) * amp
	}
	return sum / ampFractal
}

func (f *BrownianMotionFractal) Noise3D(x, y, z float64) float64 {
	if f.octaves <= 0 {
		return 0
	}

	x *= f.frequency
	y *= f.frequency
	z *= f.frequency

	sum := f.source.Noise3D(x, y, z)
	amp := 1.
	ampFractal := 1.

	for i := 0; i < f.octaves; i++ {
		x *= f.lacunarity
		y *= f.lacunarity
		z *= f.lacunarity

		amp *= f.gain
		ampFractal += amp
		sum += f.source.Noise3D(x, y, z) * amp
	}
	return sum / ampFractal
}

func (f *BrownianMotionFractal) Noise4D(x, y, z, w float64) float64 {
	if f.octaves <= 0 {
		return 0
	}

	x *= f.frequency
	y *= f.frequency
	z *= f.frequency
	w *= f.frequency

	sum := f.source.Noise4D(x, y, z, w)
	amp := 1.
	ampFractal := 1.

	for i := 0; i < f.octaves; i++ {
		x *= f.lacunarity
		y *= f.lacunarity
		z *= f.lacunarity
		w *= f.lacunarity

		amp *= f.gain
		ampFractal += amp
		sum += f.source.Noise4D(x, y, z, w) * amp
	}
	return sum / ampFractal
}

func (f *BrownianMotionFractal) Noise5D(x, y, z, w, v float64) float64 {
	if f.octaves <= 0 {
		return 0
	}

	x *= f.frequency
	y *= f.frequency
	z *= f.frequency
	w *= f.frequency
	v *= f.frequency

	sum := f.source.Noise5D(x, y, z, w, v)
	amp := 1.
	ampFractal := 1.

	for i := 0; i < f.octaves; i++ {
		x *= f.lacunarity
		y *= f.lacunarity
		z *= f.lacunarity
		w *= f.lacunarity
		v *= f.lacunarity

		amp *= f.gain
		ampFractal += amp
		sum += f.source.Noise5D(x, y, z, w, v) * amp
	}
	return sum / ampFractal
}

func (f *BrownianMotionFractal) Noise6D(x, y, z, w, v, u float64) float64 {
	if f.octaves <= 0 {
		return 0
	}

	x *= f.frequency
	y *= f.frequency
	z *= f.frequency
	w *= f.frequency
	u *= f.frequency
	v *= f.frequency

	sum := f.source.Noise6D(x, y, z, w, v, u)
	amp := 1.
	ampFractal := 1.

	for i := 0; i < f.octaves; i++ {
		x *= f.lacunarity
		y *= f.lacunarity
		z *= f.lacunarity
		w *= f.lacunarity
		u *= f.lacunarity
		v *= f.lacunarity

		amp *= f.gain
		ampFractal += amp
		sum += f.source.Noise6D(x, y, z, w, v, u) * amp
	}
	return sum / ampFractal
}
