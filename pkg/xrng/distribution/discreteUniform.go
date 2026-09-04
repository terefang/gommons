package distribution

import (
	"fmt"
	"math"
	"math/rand/v2"
)

// DiscreteUniformDistribution represents a two-parameter discrete uniform distribution with integer range [Alpha, Beta].
type DiscreteUniformDistribution struct {
	Generator rand.Source
	Alpha     int
	Beta      int
}

// Ensure DiscreteUniformDistribution implements the Distribution interface.
var _ Distribution = (*DiscreteUniformDistribution)(nil)

// NewDiscreteUniformDistribution creates a DiscreteUniformDistribution with the given generator, alpha, and beta.
func NewDiscreteUniformDistribution(gen rand.Source, alpha, beta int) (*DiscreteUniformDistribution, error) {
	dud := &DiscreteUniformDistribution{
		Generator: gen,
	}
	if !dud.SetParameters(float64(alpha), float64(beta), 0.0) {
		return nil, fmt.Errorf("given alpha (%d) and/or beta (%d) are invalid", alpha, beta)
	}
	return dud, nil
}

func (dud *DiscreteUniformDistribution) GetTag() string {
	return "DiscreteUniform"
}

func (dud *DiscreteUniformDistribution) GetParameterA() float64 {
	return float64(dud.Alpha)
}

func (dud *DiscreteUniformDistribution) GetParameterB() float64 {
	return float64(dud.Beta)
}

func (dud *DiscreteUniformDistribution) GetParameterC() float64 {
	return math.NaN()
}

func (dud *DiscreteUniformDistribution) GetMaximum() float64 {
	return float64(dud.Beta)
}

func (dud *DiscreteUniformDistribution) GetMinimum() float64 {
	return float64(dud.Alpha)
}

func (dud *DiscreteUniformDistribution) GetMean() float64 {
	return float64(dud.Alpha+dud.Beta) * 0.5
}

func (dud *DiscreteUniformDistribution) GetMedian() float64 {
	return float64(dud.Alpha+dud.Beta) * 0.5
}

func (dud *DiscreteUniformDistribution) GetMode() []float64 {
	panic("Mode is undefined for DiscreteUniform distribution")
}

func (dud *DiscreteUniformDistribution) GetVariance() float64 {
	diff := float64(dud.Beta - dud.Alpha + 1)
	return ((diff * diff) - 1.0) / 12.0
}

func (dud *DiscreteUniformDistribution) SetParameters(a, b, c float64) bool {
	if int(a) <= int(b) {
		dud.Alpha = int(a)
		dud.Beta = int(b)
		return true
	}
	return false
}

func (dud *DiscreteUniformDistribution) NextDouble() float64 {
	return SampleDiscreteUniform(dud.Generator, dud.Alpha, dud.Beta)
}

func (dud *DiscreteUniformDistribution) Copy() Distribution {
	return &DiscreteUniformDistribution{
		Generator: dud.Generator,
		Alpha:     dud.Alpha,
		Beta:      dud.Beta,
	}
}

func (dud *DiscreteUniformDistribution) StringSerialize() string {
	return fmt.Sprintf("DiscreteUniform:%d:%d", dud.Alpha, dud.Beta)
}

func (dud *DiscreteUniformDistribution) StringDeserialize(data string) (Distribution, error) {
	var alpha, beta int
	_, err := fmt.Sscanf(data, "DiscreteUniform:%d:%d", &alpha, &beta)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize DiscreteUniformDistribution: %w", err)
	}
	if !dud.SetParameters(float64(alpha), float64(beta), 0) {
		return nil, fmt.Errorf("invalid parameters in serialized data: alpha=%d, beta=%d", alpha, beta)
	}
	return dud, nil
}

// SampleDiscreteUniform generates a discrete uniform integer variate in the range [alpha, beta], returned as a float64.
func SampleDiscreteUniform(gen rand.Source, alpha, beta int) float64 {
	r := rand.New(gen)
	// Uniformly samples in integer range [alpha, beta]
	n := uint64(beta - alpha + 1)
	return float64(int64(alpha) + int64(r.Uint64N(n)))
}
