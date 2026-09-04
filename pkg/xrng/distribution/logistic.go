package distribution

import (
	"fmt"
	"math"
	"math/rand/v2"
)

// LogisticDistribution represents a two-parameter continuous distribution with infinite range (-Inf, +Inf).
type LogisticDistribution struct {
	Generator rand.Source
	Mu        float64 // Location parameter
	Sigma     float64 // Scale parameter (sigma > 0)
}

// Ensure LogisticDistribution implements the Distribution interface.
var _ Distribution = (*LogisticDistribution)(nil)

// NewLogisticDistribution creates a LogisticDistribution with the given generator, mu (location), and sigma (scale).
func NewLogisticDistribution(gen rand.Source, mu, sigma float64) (*LogisticDistribution, error) {
	ld := &LogisticDistribution{
		Generator: gen,
	}
	if !ld.SetParameters(mu, sigma, 0.0) {
		return nil, fmt.Errorf("given mu (%f) and/or sigma (%f) are invalid", mu, sigma)
	}
	return ld, nil
}

func (ld *LogisticDistribution) GetTag() string {
	return "Logistic"
}

func (ld *LogisticDistribution) GetParameterA() float64 {
	return ld.Mu
}

func (ld *LogisticDistribution) GetParameterB() float64 {
	return ld.Sigma
}

func (ld *LogisticDistribution) GetParameterC() float64 {
	return math.NaN()
}

func (ld *LogisticDistribution) GetMaximum() float64 {
	return math.Inf(1)
}

func (ld *LogisticDistribution) GetMinimum() float64 {
	return math.Inf(-1)
}

func (ld *LogisticDistribution) GetMean() float64 {
	return ld.Mu
}

func (ld *LogisticDistribution) GetMedian() float64 {
	return ld.Mu
}

func (ld *LogisticDistribution) GetMode() []float64 {
	return []float64{ld.Mu}
}

func (ld *LogisticDistribution) GetVariance() float64 {
	return ld.Sigma * ld.Sigma * math.Pi * math.Pi / 3.0
}

func (ld *LogisticDistribution) SetParameters(a, b, c float64) bool {
	if !math.IsNaN(a) && b > 0.0 {
		ld.Mu = a
		ld.Sigma = b
		return true
	}
	return false
}

func (ld *LogisticDistribution) NextDouble() float64 {
	return SampleLogistic(ld.Generator, ld.Mu, ld.Sigma)
}

func (ld *LogisticDistribution) Copy() Distribution {
	return &LogisticDistribution{
		Generator: ld.Generator,
		Mu:        ld.Mu,
		Sigma:     ld.Sigma,
	}
}

func (ld *LogisticDistribution) StringSerialize() string {
	return fmt.Sprintf("Logistic:%f:%f", ld.Mu, ld.Sigma)
}

func (ld *LogisticDistribution) StringDeserialize(data string) (Distribution, error) {
	var mu, sigma float64
	_, err := fmt.Sscanf(data, "Logistic:%f:%f", &mu, &sigma)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize LogisticDistribution: %w", err)
	}
	if !ld.SetParameters(mu, sigma, 0) {
		return nil, fmt.Errorf("invalid parameters in serialized data: mu=%f, sigma=%f", mu, sigma)
	}
	return ld, nil
}

// SampleLogistic generates a Logistic-distributed variate using inverse transform sampling (quantile function).
func SampleLogistic(gen rand.Source, mu, sigma float64) float64 {
	r := rand.New(gen)
	// Equivalent to nextExclusiveDouble() in open interval (0, 1)
	uExclusive := (float64(r.Uint64()>>11)*1.1102230246251565e-16 + 5.551115123125782e-17)

	return mu + sigma*math.Log(uExclusive/(1.0-uExclusive))
}
