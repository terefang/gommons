package xrng

const BASE_H = 1.987654321
const BASE_octaves = 8.
const BASE_frequency = 1.3125
const BASE_lacunarity = 1. / BASE_H
const BASE_gain = 0.54321

type NoiseSource interface {
	Noise1D(x float64) float64
	Noise2D(x, y float64) float64
	Noise3D(x, y, z float64) float64
	Noise4D(x, y, z, w float64) float64
	Noise5D(x, y, z, w, v float64) float64
	Noise6D(x, y, z, w, v, u float64) float64
}
