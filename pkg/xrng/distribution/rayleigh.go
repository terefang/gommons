package distribution

import (
	"fmt"
	"math"
	"math/rand/v2"
)

// RayleighDistribution represents a one-parameter distribution with range (0, +Inf).
type RayleighDistribution struct {
	Generator rand.Source
	Sigma     float64 // Scale parameter (sigma > 0)
}

// Ensure RayleighDistribution implements the Distribution interface.
var _ Distribution = (*RayleighDistribution)(nil)

// NewRayleighDistribution creates a RayleighDistribution with the given generator and sigma.
func NewRayleighDistribution(gen rand.Source, sigma float64) (*RayleighDistribution, error) {
	rd := &RayleighDistribution{
		Generator: gen,
	}
	if !rd.SetParameters(sigma, 0.0, 0.0) {
		return nil, fmt.Errorf("given sigma (%f) is invalid", sigma)
	}
	return rd, nil
}

func (rd *RayleighDistribution) GetTag() string {
	return "Rayleigh"
}

func (rd *RayleighDistribution) GetParameterA() float64 {
	return rd.Sigma
}

func (rd *RayleighDistribution) GetParameterB() float64 {
	return math.NaN()
}

func (rd *RayleighDistribution) GetParameterC() float64 {
	return math.NaN()
}

func (rd *RayleighDistribution) GetMaximum() float64 {
	return math.Inf(1)
}

func (rd *RayleighDistribution) GetMinimum() float64 {
	return 0.0
}

func (rd *RayleighDistribution) GetMean() float64 {
	return rd.Sigma * math.Sqrt(0.5*math.Pi)
}

func (rd *RayleighDistribution) GetMedian() float64 {
	return rd.Sigma * math.Sqrt(math.Log(4.0))
}

func (rd *RayleighDistribution) GetMode() []float64 {
	return []float64{rd.Sigma}
}

func (rd *RayleighDistribution) GetVariance() float64 {
	return rd.Sigma * rd.Sigma * (4.0 - math.Pi) * 0.5
}

func (rd *RayleighDistribution) SetParameters(a, b, c float64) bool {
	if a > 0.0 {
		rd.Sigma = a
		return true
	}
	return false
}

func (rd *RayleighDistribution) NextDouble() float64 {
	return SampleRayleigh(rd.Generator, rd.Sigma)
}

func (rd *RayleighDistribution) Copy() Distribution {
	return &RayleighDistribution{
		Generator: rd.Generator,
		Sigma:     rd.Sigma,
	}
}

func (rd *RayleighDistribution) StringSerialize() string {
	return fmt.Sprintf("Rayleigh:%f", rd.Sigma)
}

func (rd *RayleighDistribution) StringDeserialize(data string) (Distribution, error) {
	var sigma float64
	_, err := fmt.Sscanf(data, "Rayleigh:%f", &sigma)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize RayleighDistribution: %w", err)
	}
	if !rd.SetParameters(sigma, 0, 0) {
		return nil, fmt.Errorf("invalid parameters in serialized data: sigma=%f", sigma)
	}
	return rd, nil
}

// SampleRayleigh generates a Rayleigh-distributed variate using two independent Gaussian samples.
func SampleRayleigh(gen rand.Source, sigma float64) float64 {
	r := rand.New(gen)
	n0 := r.NormFloat64() * sigma
	n1 := r.NormFloat64() * sigma
	return math.Sqrt(n0*n0 + n1*n1)
}
