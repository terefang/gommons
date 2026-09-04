package distribution

import (
	"fmt"
	"math"
	"math/rand/v2"
)

// NormalDistribution represents a two-parameter Gaussian distribution with infinite range.
type NormalDistribution struct {
	Generator rand.Source
	Mu        float64
	Sigma     float64
}

// Ensure NormalDistribution implements the Distribution interface.
var _ Distribution = (*NormalDistribution)(nil)

// NewNormalDistribution creates a NormalDistribution with the given generator, mu, and sigma.
func NewNormalDistribution(gen rand.Source, mu, sigma float64) (*NormalDistribution, error) {
	nd := &NormalDistribution{
		Generator: gen,
	}
	if !nd.SetParameters(mu, sigma, 0.0) {
		return nil, fmt.Errorf("given mu (%f) and/or sigma (%f) are invalid", mu, sigma)
	}
	return nd, nil
}

func (nd *NormalDistribution) GetTag() string {
	return "Normal"
}

func (nd *NormalDistribution) GetParameterA() float64 {
	return nd.Mu
}

func (nd *NormalDistribution) GetParameterB() float64 {
	return nd.Sigma
}

func (nd *NormalDistribution) GetParameterC() float64 {
	return math.NaN()
}

func (nd *NormalDistribution) GetMaximum() float64 {
	return math.Inf(1)
}

func (nd *NormalDistribution) GetMinimum() float64 {
	return math.Inf(-1)
}

func (nd *NormalDistribution) GetMean() float64 {
	return nd.Mu
}

func (nd *NormalDistribution) GetMedian() float64 {
	return nd.Mu
}

func (nd *NormalDistribution) GetMode() []float64 {
	return []float64{nd.Mu}
}

func (nd *NormalDistribution) GetVariance() float64 {
	return nd.Sigma * nd.Sigma
}

func (nd *NormalDistribution) SetParameters(a, b, c float64) bool {
	if !math.IsNaN(a) && b > 0.0 {
		nd.Mu = a
		nd.Sigma = b
		return true
	}
	return false
}

func (nd *NormalDistribution) NextDouble() float64 {
	return SampleNormal(nd.Generator, nd.Mu, nd.Sigma)
}

func (nd *NormalDistribution) Copy() Distribution {
	return &NormalDistribution{
		Generator: nd.Generator,
		Mu:        nd.Mu,
		Sigma:     nd.Sigma,
	}
}

func (nd *NormalDistribution) StringSerialize() string {
	return fmt.Sprintf("Normal:%f:%f", nd.Mu, nd.Sigma)
}

func (nd *NormalDistribution) StringDeserialize(data string) (Distribution, error) {
	var mu, sigma float64
	_, err := fmt.Sscanf(data, "Normal:%f:%f", &mu, &sigma)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize NormalDistribution: %w", err)
	}
	if !nd.SetParameters(mu, sigma, 0) {
		return nil, fmt.Errorf("invalid parameters in serialized data: mu=%f, sigma=%f", mu, sigma)
	}
	return nd, nil
}

// SampleNormal generates a normally distributed variate using the given generator, mu, and sigma.
func SampleNormal(gen rand.Source, mu, sigma float64) float64 {
	// Uses standard Box-Muller transformation or standard library distribution adapter
	r := rand.New(gen)
	return mu + r.NormFloat64()*sigma
}
