package distribution

import (
	"fmt"
	"math"
	"math/rand/v2"
)

// LogNormalDistribution represents a two-parameter distribution with range (0, +Inf).
type LogNormalDistribution struct {
	Generator rand.Source
	Mu        float64 // Location parameter (mean of the underlying normal distribution)
	Sigma     float64 // Scale parameter (standard deviation of the underlying normal distribution, > 0)
}

// Ensure LogNormalDistribution implements the Distribution interface.
var _ Distribution = (*LogNormalDistribution)(nil)

// NewLogNormalDistribution creates a LogNormalDistribution with the given generator, mu, and sigma.
func NewLogNormalDistribution(gen rand.Source, mu, sigma float64) (*LogNormalDistribution, error) {
	lnd := &LogNormalDistribution{
		Generator: gen,
	}
	if !lnd.SetParameters(mu, sigma, 0.0) {
		return nil, fmt.Errorf("given mu (%f) and/or sigma (%f) are invalid", mu, sigma)
	}
	return lnd, nil
}

func (lnd *LogNormalDistribution) GetTag() string {
	return "LogNormal"
}

func (lnd *LogNormalDistribution) GetParameterA() float64 {
	return lnd.Mu
}

func (lnd *LogNormalDistribution) GetParameterB() float64 {
	return lnd.Sigma
}

func (lnd *LogNormalDistribution) GetParameterC() float64 {
	return math.NaN()
}

func (lnd *LogNormalDistribution) GetMaximum() float64 {
	return math.Inf(1)
}

func (lnd *LogNormalDistribution) GetMinimum() float64 {
	return 0.0
}

func (lnd *LogNormalDistribution) GetMean() float64 {
	return math.Exp(lnd.Mu + 0.5*lnd.Sigma*lnd.Sigma)
}

func (lnd *LogNormalDistribution) GetMedian() float64 {
	return math.Exp(lnd.Mu)
}

func (lnd *LogNormalDistribution) GetMode() []float64 {
	return []float64{math.Exp(lnd.Mu - lnd.Sigma*lnd.Sigma)}
}

func (lnd *LogNormalDistribution) GetVariance() float64 {
	sigmaSq := lnd.Sigma * lnd.Sigma
	return (math.Exp(sigmaSq) - 1.0) * math.Exp(2.0*lnd.Mu+sigmaSq)
}

func (lnd *LogNormalDistribution) SetParameters(a, b, c float64) bool {
	if !math.IsNaN(a) && b > 0.0 {
		lnd.Mu = a
		lnd.Sigma = b
		return true
	}
	return false
}

func (lnd *LogNormalDistribution) NextDouble() float64 {
	return SampleLogNormal(lnd.Generator, lnd.Mu, lnd.Sigma)
}

func (lnd *LogNormalDistribution) Copy() Distribution {
	return &LogNormalDistribution{
		Generator: lnd.Generator,
		Mu:        lnd.Mu,
		Sigma:     lnd.Sigma,
	}
}

func (lnd *LogNormalDistribution) StringSerialize() string {
	return fmt.Sprintf("LogNormal:%f:%f", lnd.Mu, lnd.Sigma)
}

func (lnd *LogNormalDistribution) StringDeserialize(data string) (Distribution, error) {
	var mu, sigma float64
	_, err := fmt.Sscanf(data, "LogNormal:%f:%f", &mu, &sigma)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize LogNormalDistribution: %w", err)
	}
	if !lnd.SetParameters(mu, sigma, 0) {
		return nil, fmt.Errorf("invalid parameters in serialized data: mu=%f, sigma=%f", mu, sigma)
	}
	return lnd, nil
}

// SampleLogNormal generates a Log-Normal distributed variate by exponentiating a Gaussian (normal) sample.
func SampleLogNormal(gen rand.Source, mu, sigma float64) float64 {
	r := rand.New(gen)
	gaussianSample := mu + r.NormFloat64()*sigma
	return math.Exp(gaussianSample)
}
