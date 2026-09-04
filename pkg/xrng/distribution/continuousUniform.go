package distribution

import (
	"fmt"
	"math"
	"math/rand/v2"
)

// ContinuousUniformDistribution represents a two-parameter continuous uniform distribution with range [Alpha, Beta].
type ContinuousUniformDistribution struct {
	Generator rand.Source
	Alpha     float64
	Beta      float64
}

// Ensure ContinuousUniformDistribution implements the Distribution interface.
var _ Distribution = (*ContinuousUniformDistribution)(nil)

// NewContinuousUniformDistribution creates a ContinuousUniformDistribution with the given generator, alpha, and beta.
func NewContinuousUniformDistribution(gen rand.Source, alpha, beta float64) (*ContinuousUniformDistribution, error) {
	cud := &ContinuousUniformDistribution{
		Generator: gen,
	}
	if !cud.SetParameters(alpha, beta, 0.0) {
		return nil, fmt.Errorf("given alpha (%f) and/or beta (%f) are invalid", alpha, beta)
	}
	return cud, nil
}

func (cud *ContinuousUniformDistribution) GetTag() string {
	return "ContinuousUniform"
}

func (cud *ContinuousUniformDistribution) GetParameterA() float64 {
	return cud.Alpha
}

func (cud *ContinuousUniformDistribution) GetParameterB() float64 {
	return cud.Beta
}

func (cud *ContinuousUniformDistribution) GetParameterC() float64 {
	return math.NaN()
}

func (cud *ContinuousUniformDistribution) GetMaximum() float64 {
	return cud.Beta
}

func (cud *ContinuousUniformDistribution) GetMinimum() float64 {
	return cud.Alpha
}

func (cud *ContinuousUniformDistribution) GetMean() float64 {
	return (cud.Alpha + cud.Beta) * 0.5
}

func (cud *ContinuousUniformDistribution) GetMedian() float64 {
	return (cud.Alpha + cud.Beta) * 0.5
}

func (cud *ContinuousUniformDistribution) GetMode() []float64 {
	panic("Mode is undefined for ContinuousUniform distribution")
}

func (cud *ContinuousUniformDistribution) GetVariance() float64 {
	diff := cud.Beta - cud.Alpha
	return (diff * diff) / 12.0
}

func (cud *ContinuousUniformDistribution) SetParameters(a, b, c float64) bool {
	if a <= b {
		cud.Alpha = a
		cud.Beta = b
		return true
	}
	return false
}

func (cud *ContinuousUniformDistribution) NextDouble() float64 {
	return SampleContinuousUniform(cud.Generator, cud.Alpha, cud.Beta)
}

func (cud *ContinuousUniformDistribution) Copy() Distribution {
	return &ContinuousUniformDistribution{
		Generator: cud.Generator,
		Alpha:     cud.Alpha,
		Beta:      cud.Beta,
	}
}

func (cud *ContinuousUniformDistribution) StringSerialize() string {
	return fmt.Sprintf("ContinuousUniform:%f:%f", cud.Alpha, cud.Beta)
}

func (cud *ContinuousUniformDistribution) StringDeserialize(data string) (Distribution, error) {
	var alpha, beta float64
	_, err := fmt.Sscanf(data, "ContinuousUniform:%f:%f", &alpha, &beta)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize ContinuousUniformDistribution: %w", err)
	}
	if !cud.SetParameters(alpha, beta, 0) {
		return nil, fmt.Errorf("invalid parameters in serialized data: alpha=%f, beta=%f", alpha, beta)
	}
	return cud, nil
}

// SampleContinuousUniform generates a continuously uniform random variable in the closed interval [alpha, beta].
func SampleContinuousUniform(gen rand.Source, alpha, beta float64) float64 {
	r := rand.New(gen)
	// Generates a uniform float64 in closed interval [0.0, 1.0]
	uInclusive := float64(r.Uint64()>>11) * (1.0 / 9007199254740991.0)
	return alpha + uInclusive*(beta-alpha)
}
