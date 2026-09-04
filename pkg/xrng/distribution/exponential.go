package distribution

import (
	"fmt"
	"math"
	"math/rand/v2"
)

// ExponentialDistribution represents a one-parameter distribution with range (0, +Inf).
type ExponentialDistribution struct {
	Generator rand.Source
	Lambda    float64 // Stored internally as 1.0 / lambda (mean)
}

// Ensure ExponentialDistribution implements the Distribution interface.
var _ Distribution = (*ExponentialDistribution)(nil)

// NewExponentialDistribution creates an ExponentialDistribution with the given generator and rate parameter (lambda).
func NewExponentialDistribution(gen rand.Source, lambda float64) (*ExponentialDistribution, error) {
	ed := &ExponentialDistribution{
		Generator: gen,
	}
	if !ed.SetParameters(lambda, 0.0, 0.0) {
		return nil, fmt.Errorf("given lambda (%f) is invalid", lambda)
	}
	return ed, nil
}

func (ed *ExponentialDistribution) GetTag() string {
	return "Exponential"
}

func (ed *ExponentialDistribution) GetParameterA() float64 {
	return 1.0 / ed.Lambda
}

func (ed *ExponentialDistribution) GetParameterB() float64 {
	return math.NaN()
}

func (ed *ExponentialDistribution) GetParameterC() float64 {
	return math.NaN()
}

func (ed *ExponentialDistribution) GetMaximum() float64 {
	return math.Inf(1)
}

func (ed *ExponentialDistribution) GetMinimum() float64 {
	return 0.0
}

func (ed *ExponentialDistribution) GetMean() float64 {
	return ed.Lambda
}

func (ed *ExponentialDistribution) GetMedian() float64 {
	return ed.Lambda * math.Ln2
}

func (ed *ExponentialDistribution) GetMode() []float64 {
	return []float64{0.0}
}

func (ed *ExponentialDistribution) GetVariance() float64 {
	return ed.Lambda * ed.Lambda
}

func (ed *ExponentialDistribution) SetParameters(a, b, c float64) bool {
	if a > 0.0 {
		ed.Lambda = 1.0 / a
		return true
	}
	return false
}

func (ed *ExponentialDistribution) NextDouble() float64 {
	return SampleExponential(ed.Generator, ed.Lambda)
}

func (ed *ExponentialDistribution) Copy() Distribution {
	return &ExponentialDistribution{
		Generator: ed.Generator,
		Lambda:    ed.Lambda,
	}
}

func (ed *ExponentialDistribution) StringSerialize() string {
	return fmt.Sprintf("Exponential:%f", 1.0/ed.Lambda)
}

func (ed *ExponentialDistribution) StringDeserialize(data string) (Distribution, error) {
	var lambda float64
	_, err := fmt.Sscanf(data, "Exponential:%f", &lambda)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize ExponentialDistribution: %w", err)
	}
	if !ed.SetParameters(lambda, 0, 0) {
		return nil, fmt.Errorf("invalid parameters in serialized data: lambda=%f", lambda)
	}
	return ed, nil
}

// SampleExponential generates an Exponential-distributed variate using inverse transform sampling.
func SampleExponential(gen rand.Source, inverseLambda float64) float64 {
	r := rand.New(gen)
	// Equivalent to nextExclusiveDouble() in open interval (0, 1)
	uExclusive := (float64(r.Uint64()>>11)*1.1102230246251565e-16 + 5.551115123125782e-17)

	return -math.Log(uExclusive) * inverseLambda
}
