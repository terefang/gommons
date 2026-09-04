package sampler

import "github.com/terefang/gommons/pkg/xrng"

type SamplerFractal struct {
	source        xrng.NoiseSource
	octaves       int
	frequency     float64
	lacunarity    float64
	gain          float64
	fractalSpiral bool
}

func NewSamplerFractal(source xrng.NoiseSource) *SamplerFractal {
	return &SamplerFractal{source: source, octaves: xrng.BASE_octaves, frequency: xrng.BASE_frequency, lacunarity: xrng.BASE_lacunarity, gain: xrng.BASE_gain, fractalSpiral: false}
}

func (s *SamplerFractal) SetOctaves(octaves int) {
	s.octaves = octaves
}

func (s *SamplerFractal) SetFrequency(frequency float64) {
	s.frequency = frequency
}

func (s *SamplerFractal) SetLacunarity(lacunarity float64) {
	s.lacunarity = lacunarity
}

func (s *SamplerFractal) SetGain(gain float64) {
	s.gain = gain
}

func (s *SamplerFractal) SetFractalSpiral(fractalSpiral bool) {
	s.fractalSpiral = fractalSpiral
}

func rotateX2D(x, y float64) float64 {
	return x*+0.6088885514347261 + y*-0.7943553508622062
}

func rotateY2D(x, y float64) float64 {
	return x*+0.7943553508622062 + y*+0.6088885514347261
}

func (s SamplerFractal) Noise1D(x float64) float64 {
	x *= s.frequency
	y := 0.

	amp := 1.0
	adj := 1.0
	sum := 0.0

	for i := 0; i < s.octaves; i++ {

		sum += s.source.Noise1D(x) * amp
		amp *= s.gain

		if s.fractalSpiral {
			x2 := rotateX2D(x, y)
			y2 := rotateY2D(x, y)
			x, y = x2, y2
		}

		x /= s.lacunarity
		y /= s.lacunarity

		adj += amp
	}

	return sum / adj
}

func (s SamplerFractal) Noise2D(x, y float64) float64 {
	x *= s.frequency
	y *= s.frequency

	amp := 1.0
	adj := 1.0
	sum := 0.0

	for i := 0; i < s.octaves; i++ {

		sum += s.source.Noise2D(x, y) * amp
		amp *= s.gain

		if s.fractalSpiral {
			x2 := rotateX2D(x, y)
			y2 := rotateY2D(x, y)
			x, y = x2, y2
		}

		x /= s.lacunarity
		y /= s.lacunarity

		adj += amp
	}

	return sum / adj
}

func (s SamplerFractal) Noise3D(x, y, z float64) float64 {
	x *= s.frequency
	y *= s.frequency
	z *= s.frequency

	amp := 1.0
	adj := 1.0
	sum := 0.0

	for i := 0; i < s.octaves; i++ {

		sum += s.source.Noise3D(x, y, z) * amp
		amp *= s.gain

		if s.fractalSpiral {
			x2 := rotateX2D(x, y)
			y2 := rotateY2D(x, y)
			x, y = x2, y2
		}

		x /= s.lacunarity
		y /= s.lacunarity
		z /= s.lacunarity

		adj += amp
	}

	return sum / adj
}

func (s SamplerFractal) Noise4D(x, y, z, w float64) float64 {
	x *= s.frequency
	y *= s.frequency
	z *= s.frequency
	w *= s.frequency

	amp := 1.0
	adj := 1.0
	sum := 0.0

	for i := 0; i < s.octaves; i++ {

		sum += s.source.Noise4D(x, y, z, w) * amp
		amp *= s.gain

		if s.fractalSpiral {
			x2 := rotateX2D(x, y)
			y2 := rotateY2D(x, y)
			x, y = x2, y2
		}

		x /= s.lacunarity
		y /= s.lacunarity
		z /= s.lacunarity
		w /= s.lacunarity

		adj += amp
	}

	return sum / adj
}

func (s SamplerFractal) Noise5D(x, y, z, w, v float64) float64 {
	x *= s.frequency
	y *= s.frequency
	z *= s.frequency
	w *= s.frequency
	v *= s.frequency

	amp := 1.0
	adj := 1.0
	sum := 0.0

	for i := 0; i < s.octaves; i++ {

		sum += s.source.Noise5D(x, y, z, w, v) * amp
		amp *= s.gain

		if s.fractalSpiral {
			x2 := rotateX2D(x, y)
			y2 := rotateY2D(x, y)
			x, y = x2, y2
		}

		x /= s.lacunarity
		y /= s.lacunarity
		z /= s.lacunarity
		w /= s.lacunarity
		v /= s.lacunarity

		adj += amp
	}

	return sum / adj
}

func (s SamplerFractal) Noise6D(x, y, z, w, v, u float64) float64 {
	x *= s.frequency
	y *= s.frequency
	z *= s.frequency
	w *= s.frequency
	v *= s.frequency
	u *= s.frequency

	amp := 1.0
	adj := 1.0
	sum := 0.0

	for i := 0; i < s.octaves; i++ {

		sum += s.source.Noise6D(x, y, z, w, v, u) * amp
		amp *= s.gain

		if s.fractalSpiral {
			x2 := rotateX2D(x, y)
			y2 := rotateY2D(x, y)
			x, y = x2, y2
		}

		x /= s.lacunarity
		y /= s.lacunarity
		z /= s.lacunarity
		w /= s.lacunarity
		v /= s.lacunarity
		u /= s.lacunarity

		adj += amp
	}

	return sum / adj
}
