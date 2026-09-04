package distribution

import (
	"fmt"
	"math"
	"math/rand/v2"
)

// KnobDistribution represents a three-parameter distribution that interpolates
// between a continuous uniform distribution and a normal distribution.
type KnobDistribution struct {
	Generator rand.Source
	Mu        float64 // Center parameter
	Sigma     float64 // Scale parameter
	Iota      float64 // Interpolation parameter [0.0, 1.0]
}

// Ensure KnobDistribution implements the Distribution interface.
var _ Distribution = (*KnobDistribution)(nil)

// NewKnobDistribution creates a KnobDistribution with the given generator, mu, sigma, and iota.
func NewKnobDistribution(gen rand.Source, mu, sigma, iota float64) (*KnobDistribution, error) {
	kd := &KnobDistribution{
		Generator: gen,
	}
	if !kd.SetParameters(mu, sigma, iota) {
		return nil, fmt.Errorf("given mu (%f), sigma (%f), and/or iota (%f) are invalid", mu, sigma, iota)
	}
	return kd, nil
}

func (kd *KnobDistribution) GetTag() string {
	return "Knob"
}

func (kd *KnobDistribution) GetParameterA() float64 {
	return kd.Mu
}

func (kd *KnobDistribution) GetParameterB() float64 {
	return kd.Sigma
}

func (kd *KnobDistribution) GetParameterC() float64 {
	return kd.Iota
}

func (kd *KnobDistribution) GetMaximum() float64 {
	return math.Inf(1)
}

func (kd *KnobDistribution) GetMinimum() float64 {
	return math.Inf(-1)
}

func (kd *KnobDistribution) GetMean() float64 {
	return kd.Mu
}

func (kd *KnobDistribution) GetMedian() float64 {
	return kd.Mu
}

func (kd *KnobDistribution) GetMode() []float64 {
	return []float64{kd.Mu}
}

func (kd *KnobDistribution) GetVariance() float64 {
	panic("Variance is undefined for Knob distribution")
}

func (kd *KnobDistribution) SetParameters(a, b, c float64) bool {
	if !math.IsNaN(a) && b > 0.0 && c >= 0.0 && c <= 1.0 {
		kd.Mu = a
		kd.Sigma = b
		kd.Iota = c
		return true
	}
	return false
}

func (kd *KnobDistribution) NextDouble() float64 {
	return SampleKnob(kd.Generator, kd.Mu, kd.Sigma, kd.Iota)
}

func (kd *KnobDistribution) Copy() Distribution {
	return &KnobDistribution{
		Generator: kd.Generator,
		Mu:        kd.Mu,
		Sigma:     kd.Sigma,
		Iota:      kd.Iota,
	}
}

func (kd *KnobDistribution) StringSerialize() string {
	return fmt.Sprintf("Knob:%f:%f:%f", kd.Mu, kd.Sigma, kd.Iota)
}

func (kd *KnobDistribution) StringDeserialize(data string) (Distribution, error) {
	var mu, sigma, iota float64
	_, err := fmt.Sscanf(data, "Knob:%f:%f:%f", &mu, &sigma, &iota)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize KnobDistribution: %w", err)
	}
	if !kd.SetParameters(mu, sigma, iota) {
		return nil, fmt.Errorf("invalid parameters in serialized data: mu=%f, sigma=%f, iota=%f", mu, sigma, iota)
	}
	return kd, nil
}

// SampleKnob generates a variate by linearly interpolating between a uniform
// distribution [-sigma + mu, sigma + mu] and a normal distribution N(mu, sigma^2).
func SampleKnob(gen rand.Source, mu, sigma, iota float64) float64 {
	r := rand.New(gen)

	// Uniform sample in [-sigma, sigma] + mu
	uniformSample := (r.Float64()*(2.0*sigma) - sigma) + mu

	// Gaussian sample N(mu, sigma)
	gaussianSample := mu + r.NormFloat64()*sigma

	// Linear interpolation: (1 - iota)*uniformSample + iota*gaussianSample
	return uniformSample + iota*(gaussianSample-uniformSample)
}
