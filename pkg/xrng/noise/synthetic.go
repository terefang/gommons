package noise

import (
	"github.com/terefang/gommons/pkg/xrng"
	"github.com/terefang/gommons/pkg/xrng/interpolation"
)

func NewHoneyNoise(sn, vn xrng.NoiseSource) xrng.NoiseSource {
	return interpolation.NewCosineOutSource(interpolation.NewAverageNoise(sn, vn))
}

func NewHoneyNoiseDefaults(seed int64) xrng.NoiseSource {
	return NewHoneyNoise(NewSimplex(seed), NewValueNoise(seed+1337))
}
