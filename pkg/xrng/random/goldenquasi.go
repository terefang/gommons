package random

import "math/rand"

const goldenRatioIncrement uint64 = 0x9E3779B97F4A7C15

type GoldenQuasiRandom struct {
	State uint64
}

var (
	_ rand.Source   = (*GoldenQuasiRandom)(nil)
	_ rand.Source64 = (*GoldenQuasiRandom)(nil)
)

func NewGoldenQuasiRandom(seed int64) *GoldenQuasiRandom {
	return &GoldenQuasiRandom{State: uint64(seed)}
}

// Seed implements rand.Source.
func (g *GoldenQuasiRandom) Seed(seed int64) {
	g.State = uint64(seed)
}

// Uint64 implements rand.Source64.
func (g *GoldenQuasiRandom) Uint64() uint64 {
	g.State += goldenRatioIncrement
	return g.State
}

// Int63 implements rand.Source.
func (g *GoldenQuasiRandom) Int63() int64 {
	return int64(g.Uint64() & 0x7fffffffffffffff)
}

// NextDouble returns a float64 in the half-open range (0.0, 1.0).
func (g *GoldenQuasiRandom) NextDouble() float64 {
	g.State += goldenRatioIncrement
	return float64(g.State>>11)*1.1102230246251565e-16 + 5.551115123125782e-17
}
