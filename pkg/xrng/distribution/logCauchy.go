package distribution

import (
	"fmt"
	"math"
	"math/rand/v2"
)

// LogCauchyDistribution represents a two-parameter distribution with range (0, +Inf).
type LogCauchyDistribution struct {
	Generator rand.Source
	Mu        float64 // Location parameter for underlying Cauchy distribution
	Sigma     float64 // Scale parameter (sigma > 0)
}

// Ensure LogCauchyDistribution implements the Distribution interface.
var _ Distribution = (*LogCauchyDistribution)(nil)

// NewLogCauchyDistribution creates a LogCauchyDistribution with the given generator, mu (location), and sigma (scale).
func NewLogCauchyDistribution(gen rand.Source, mu, sigma float64) (*LogCauchyDistribution, error) {
	lcd := &LogCauchyDistribution{
		Generator: gen,
	}
	if !lcd.SetParameters(mu, sigma, 0.0) {
		return nil, fmt.Errorf("given mu (%f) and/or sigma (%f) are invalid", mu, sigma)
	}
	return lcd, nil
}

func (lcd *LogCauchyDistribution) GetTag() string {
	return "LogCauchy"
}

func (lcd *LogCauchyDistribution) GetParameterA() float64 {
	return lcd.Mu
}

func (lcd *LogCauchyDistribution) GetParameterB() float64 {
	return lcd.Sigma
}

func (lcd *LogCauchyDistribution) GetParameterC() float64 {
	return math.NaN()
}

func (lcd *LogCauchyDistribution) GetMaximum() float64 {
	return math.Inf(1)
}

func (lcd *LogCauchyDistribution) GetMinimum() float64 {
	return 0.0
}

func (lcd *LogCauchyDistribution) GetMean() float64 {
	return math.Inf(1)
}

func (lcd *LogCauchyDistribution) GetMedian() float64 {
	return math.Exp(lcd.Mu)
}

func (lcd *LogCauchyDistribution) GetMode() []float64 {
	panic("Mode is undefined for LogCauchy distribution")
}

func (lcd *LogCauchyDistribution) GetVariance() float64 {
	return math.Inf(1)
}

func (lcd *LogCauchyDistribution) SetParameters(a, b, c float64) bool {
	if !math.IsNaN(a) && b > 0.0 {
		lcd.Mu = a
		lcd.Sigma = b
		return true
	}
	return false
}

func (lcd *LogCauchyDistribution) NextDouble() float64 {
	return SampleLogCauchy(lcd.Generator, lcd.Mu, lcd.Sigma)
}

func (lcd *LogCauchyDistribution) Copy() Distribution {
	return &LogCauchyDistribution{
		Generator: lcd.Generator,
		Mu:        lcd.Mu,
		Sigma:     lcd.Sigma,
	}
}

func (lcd *LogCauchyDistribution) StringSerialize() string {
	return fmt.Sprintf("LogCauchy:%f:%f", lcd.Mu, lcd.Sigma)
}

func (lcd *LogCauchyDistribution) StringDeserialize(data string) (Distribution, error) {
	var mu, sigma float64
	_, err := fmt.Sscanf(data, "LogCauchy:%f:%f", &mu, &sigma)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize LogCauchyDistribution: %w", err)
	}
	if !lcd.SetParameters(mu, sigma, 0) {
		return nil, fmt.Errorf("invalid parameters in serialized data: mu=%f, sigma=%f", mu, sigma)
	}
	return lcd, nil
}

// SampleLogCauchy generates a Log-Cauchy distributed variate by exponentiating a Cauchy sample.
func SampleLogCauchy(gen rand.Source, mu, sigma float64) float64 {
	return math.Exp(SampleCauchy(gen, mu, sigma))
}
