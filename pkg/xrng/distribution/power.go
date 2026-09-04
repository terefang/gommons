package distribution

import (
	"fmt"
	"math"
	"math/rand/v2"
)

// PowerDistribution represents a two-parameter continuous distribution with range [0, 1.0/beta].
type PowerDistribution struct {
	Generator rand.Source
	Alpha     float64 // Stored internally as 1.0 / alpha_param
	Beta      float64 // Stored internally as 1.0 / beta_param
}

// Ensure PowerDistribution implements the Distribution interface.
var _ Distribution = (*PowerDistribution)(nil)

// NewPowerDistribution creates a PowerDistribution with the given generator, alpha, and beta.
func NewPowerDistribution(gen rand.Source, alpha, beta float64) (*PowerDistribution, error) {
	pd := &PowerDistribution{
		Generator: gen,
	}
	if !pd.SetParameters(alpha, beta, 0.0) {
		return nil, fmt.Errorf("given alpha (%f) and/or beta (%f) are invalid", alpha, beta)
	}
	return pd, nil
}

func (pd *PowerDistribution) GetTag() string {
	return "Power"
}

func (pd *PowerDistribution) GetParameterA() float64 {
	return 1.0 / pd.Alpha
}

func (pd *PowerDistribution) GetParameterB() float64 {
	return 1.0 / pd.Beta
}

func (pd *PowerDistribution) GetParameterC() float64 {
	return math.NaN()
}

func (pd *PowerDistribution) GetMaximum() float64 {
	return pd.Beta
}

func (pd *PowerDistribution) GetMinimum() float64 {
	return 0.0
}

func (pd *PowerDistribution) GetMean() float64 {
	a := 1.0 / pd.Alpha
	return a * pd.Beta / (a + 1.0)
}

func (pd *PowerDistribution) GetMedian() float64 {
	panic("Median is undefined for Power distribution")
}

func (pd *PowerDistribution) GetMode() []float64 {
	if pd.Alpha > 1.0 {
		return []float64{pd.Beta}
	}
	if pd.Alpha < 1.0 {
		return []float64{0.0}
	}
	panic("Mode cannot be determined for the given parameters")
}

func (pd *PowerDistribution) GetVariance() float64 {
	a := 1.0 / pd.Alpha
	denom := (a + 1.0) * (a + 1.0)
	return a * pd.Beta * pd.Beta / denom / (a + 2.0)
}

func (pd *PowerDistribution) SetParameters(a, b, c float64) bool {
	if a > 0.0 && b > 0.0 {
		pd.Alpha = 1.0 / a
		pd.Beta = 1.0 / b
		return true
	}
	return false
}

func (pd *PowerDistribution) NextDouble() float64 {
	return SamplePower(pd.Generator, pd.Alpha, pd.Beta)
}

func (pd *PowerDistribution) Copy() Distribution {
	return &PowerDistribution{
		Generator: pd.Generator,
		Alpha:     pd.Alpha,
		Beta:      pd.Beta,
	}
}

func (pd *PowerDistribution) StringSerialize() string {
	return fmt.Sprintf("Power:%f:%f", 1.0/pd.Alpha, 1.0/pd.Beta)
}

func (pd *PowerDistribution) StringDeserialize(data string) (Distribution, error) {
	var alpha, beta float64
	_, err := fmt.Sscanf(data, "Power:%f:%f", &alpha, &beta)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize PowerDistribution: %w", err)
	}
	if !pd.SetParameters(alpha, beta, 0) {
		return nil, fmt.Errorf("invalid parameters in serialized data: alpha=%f, beta=%f", alpha, beta)
	}
	return pd, nil
}

// SamplePower generates a Power-distributed variate using inverse transform sampling.
func SamplePower(gen rand.Source, inverseAlpha, inverseBeta float64) float64 {
	r := rand.New(gen)
	// Equivalent to nextExclusiveDouble() in open interval (0, 1)
	uExclusive := (float64(r.Uint64()>>11)*1.1102230246251565e-16 + 5.551115123125782e-17)

	return math.Pow(uExclusive, inverseAlpha) * inverseBeta
}
