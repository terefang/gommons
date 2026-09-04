package fractal

import "github.com/terefang/gommons/pkg/xrng"

type FbmFractal struct {
	source      xrng.NoiseSource
	octaves     int
	lacunarity  float64
	persistence float64
}

func (f *FbmFractal) SetOctaves(octaves int) {
	f.octaves = octaves
}

func (f *FbmFractal) SetLacunarity(lacunarity float64) {
	f.lacunarity = lacunarity
}

func (f *FbmFractal) SetPersistence(persistence float64) {
	f.persistence = persistence
}

func NewFbmFractal(source xrng.NoiseSource) *FbmFractal {
	return &FbmFractal{source: source, octaves: xrng.BASE_octaves, lacunarity: xrng.BASE_H, persistence: 1. / xrng.BASE_gain}
}

func (f *FbmFractal) Noise1D(x float64) float64 {
	if f.octaves <= 0 {
		return 0
	}

	value := 0.0
	amplitude := 1.0
	frequency := 1.0
	maxAmplitude := 0.0

	for i := 0; i < f.octaves; i++ {
		value += f.source.Noise1D(x*frequency) * amplitude
		maxAmplitude += amplitude
		frequency *= f.lacunarity
		amplitude *= f.persistence
	}

	if maxAmplitude == 0 {
		return 0
	}
	return value / maxAmplitude
}

func (f *FbmFractal) Noise2D(x, y float64) float64 {
	if f.octaves <= 0 {
		return 0
	}

	value := 0.0
	amplitude := 1.0
	frequency := 1.0
	maxAmplitude := 0.0

	for i := 0; i < f.octaves; i++ {
		value += f.source.Noise2D(x*frequency, y*frequency) * amplitude
		maxAmplitude += amplitude
		frequency *= f.lacunarity
		amplitude *= f.persistence
	}

	if maxAmplitude == 0 {
		return 0
	}
	return value / maxAmplitude
}

func (f *FbmFractal) Noise3D(x, y, z float64) float64 {
	if f.octaves <= 0 {
		return 0
	}

	value := 0.0
	amplitude := 1.0
	frequency := 1.0
	maxAmplitude := 0.0

	for i := 0; i < f.octaves; i++ {
		value += f.source.Noise3D(x*frequency, y*frequency, z*frequency) * amplitude
		maxAmplitude += amplitude
		frequency *= f.lacunarity
		amplitude *= f.persistence
	}

	if maxAmplitude == 0 {
		return 0
	}
	return value / maxAmplitude
}

func (f *FbmFractal) Noise4D(x, y, z, w float64) float64 {
	if f.octaves <= 0 {
		return 0
	}

	value := 0.0
	amplitude := 1.0
	frequency := 1.0
	maxAmplitude := 0.0

	for i := 0; i < f.octaves; i++ {
		value += f.source.Noise4D(x*frequency, y*frequency, z*frequency, w*frequency) * amplitude
		maxAmplitude += amplitude
		frequency *= f.lacunarity
		amplitude *= f.persistence
	}

	if maxAmplitude == 0 {
		return 0
	}
	return value / maxAmplitude
}

func (f *FbmFractal) Noise5D(x, y, z, w, v float64) float64 {
	if f.octaves <= 0 {
		return 0
	}

	value := 0.0
	amplitude := 1.0
	frequency := 1.0
	maxAmplitude := 0.0

	for i := 0; i < f.octaves; i++ {
		value += f.source.Noise5D(x*frequency, y*frequency, z*frequency, w*frequency, v*frequency) * amplitude
		maxAmplitude += amplitude
		frequency *= f.lacunarity
		amplitude *= f.persistence
	}

	if maxAmplitude == 0 {
		return 0
	}
	return value / maxAmplitude
}

func (f *FbmFractal) Noise6D(x, y, z, w, v, u float64) float64 {
	if f.octaves <= 0 {
		return 0
	}

	value := 0.0
	amplitude := 1.0
	frequency := 1.0
	maxAmplitude := 0.0

	for i := 0; i < f.octaves; i++ {
		value += f.source.Noise6D(x*frequency, y*frequency, z*frequency, w*frequency, v*frequency, u*frequency) * amplitude
		maxAmplitude += amplitude
		frequency *= f.lacunarity
		amplitude *= f.persistence
	}

	if maxAmplitude == 0 {
		return 0
	}
	return value / maxAmplitude
}
