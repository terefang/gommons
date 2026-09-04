package distribution

import (
	"fmt"
	"math"
	"math/rand/v2"
)

// PoissonDistribution represents a one-parameter discrete distribution with integer range [0, +Inf).
type PoissonDistribution struct {
	Generator rand.Source
	Lambda    float64 // Rate parameter (lambda > 0)
}

// Ensure PoissonDistribution implements the Distribution interface.
var _ Distribution = (*PoissonDistribution)(nil)

// NewPoissonDistribution creates a PoissonDistribution with the given generator and lambda.
func NewPoissonDistribution(gen rand.Source, lambda float64) (*PoissonDistribution, error) {
	pd := &PoissonDistribution{
		Generator: gen,
	}
	if !pd.SetParameters(lambda, 0.0, 0.0) {
		return nil, fmt.Errorf("given lambda (%f) is invalid", lambda)
	}
	return pd, nil
}

func (pd *PoissonDistribution) GetTag() string {
	return "Poisson"
}

func (pd *PoissonDistribution) GetParameterA() float64 {
	return pd.Lambda
}

func (pd *PoissonDistribution) GetParameterB() float64 {
	return math.NaN()
}

func (pd *PoissonDistribution) GetParameterC() float64 {
	return math.NaN()
}

func (pd *PoissonDistribution) GetMaximum() float64 {
	return math.Inf(1)
}

func (pd *PoissonDistribution) GetMinimum() float64 {
	return 0.0
}

func (pd *PoissonDistribution) GetMean() float64 {
	return pd.Lambda
}

func (pd *PoissonDistribution) GetMedian() float64 {
	panic("Median is undefined for Poisson distribution")
}

func (pd *PoissonDistribution) GetMode() []float64 {
	floorLambda := math.Floor(pd.Lambda)
	// Tolerance check corresponding to MathTools.isEqual(lambda, Math.floor(lambda), 0x1p-24)
	if math.Abs(pd.Lambda-floorLambda) <= 5.9604644775390625e-08 {
		return []float64{pd.Lambda - 1.0, pd.Lambda}
	}
	return []float64{floorLambda}
}

func (pd *PoissonDistribution) GetVariance() float64 {
	return pd.Lambda
}

func (pd *PoissonDistribution) SetParameters(a, b, c float64) bool {
	if a > 0.0 {
		pd.Lambda = a
		return true
	}
	return false
}

func (pd *PoissonDistribution) NextDouble() float64 {
	return SamplePoisson(pd.Generator, pd.Lambda)
}

func (pd *PoissonDistribution) Copy() Distribution {
	return &PoissonDistribution{
		Generator: pd.Generator,
		Lambda:    pd.Lambda,
	}
}

func (pd *PoissonDistribution) StringSerialize() string {
	return fmt.Sprintf("Poisson:%f", pd.Lambda)
}

func (pd *PoissonDistribution) StringDeserialize(data string) (Distribution, error) {
	var lambda float64
	_, err := fmt.Sscanf(data, "Poisson:%f", &lambda)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize PoissonDistribution: %w", err)
	}
	if !pd.SetParameters(lambda, 0, 0) {
		return nil, fmt.Errorf("invalid parameters in serialized data: lambda=%f", lambda)
	}
	return pd, nil
}

// SamplePoisson generates a Poisson-distributed variate using Knuth's algorithm / inverse transform sampling.
func SamplePoisson(gen rand.Source, lambda float64) float64 {
	r := rand.New(gen)
	x := 0.0
	p := math.Exp(-lambda)
	s := p
	u := r.Float64() // Uniform sample in [0, 1)

	for u > s {
		x++
		p *= lambda / x
		s += p
	}
	return x
}
