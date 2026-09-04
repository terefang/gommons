package distribution

import (
	"fmt"
	"math"
	"math/rand/v2"
)

// ChiDistribution represents a one-parameter Chi distribution with range (0, +Inf).
type ChiDistribution struct {
	Generator rand.Source
	Alpha     int
}

// Ensure ChiDistribution implements the Distribution interface.
var _ Distribution = (*ChiDistribution)(nil)

// NewChiDistribution creates a ChiDistribution with the given generator and alpha (degrees of freedom).
func NewChiDistribution(gen rand.Source, alpha int) (*ChiDistribution, error) {
	cd := &ChiDistribution{
		Generator: gen,
	}
	if !cd.SetParameters(float64(alpha), 0.0, 0.0) {
		return nil, fmt.Errorf("given alpha (%d) is invalid", alpha)
	}
	return cd, nil
}

func (cd *ChiDistribution) GetTag() string {
	return "Chi"
}

func (cd *ChiDistribution) GetParameterA() float64 {
	return float64(cd.Alpha)
}

func (cd *ChiDistribution) GetParameterB() float64 {
	return math.NaN()
}

func (cd *ChiDistribution) GetParameterC() float64 {
	return math.NaN()
}

func (cd *ChiDistribution) GetMaximum() float64 {
	return math.Inf(1)
}

func (cd *ChiDistribution) GetMinimum() float64 {
	return 0.0
}

func (cd *ChiDistribution) GetMean() float64 {
	a := float64(cd.Alpha)
	gamma1, _ := math.Lgamma((a + 1.0) * 0.5)
	gamma2, _ := math.Lgamma(a * 0.5)
	return math.Sqrt(2.0) * math.Exp(gamma1-gamma2)
}

func (cd *ChiDistribution) GetMedian() float64 {
	panic("Median is undefined for Chi distribution")
}

func (cd *ChiDistribution) GetMode() []float64 {
	return []float64{math.Sqrt(float64(cd.Alpha - 1))}
}

func (cd *ChiDistribution) GetVariance() float64 {
	mean := cd.GetMean()
	return float64(cd.Alpha) - (mean * mean)
}

func (cd *ChiDistribution) SetParameters(a, b, c float64) bool {
	if a >= 1.0 {
		cd.Alpha = int(a)
		return true
	}
	return false
}

func (cd *ChiDistribution) NextDouble() float64 {
	return SampleChi(cd.Generator, cd.Alpha)
}

func (cd *ChiDistribution) Copy() Distribution {
	return &ChiDistribution{
		Generator: cd.Generator,
		Alpha:     cd.Alpha,
	}
}

func (cd *ChiDistribution) StringSerialize() string {
	return fmt.Sprintf("Chi:%d", cd.Alpha)
}

func (cd *ChiDistribution) StringDeserialize(data string) (Distribution, error) {
	var alpha int
	_, err := fmt.Sscanf(data, "Chi:%d", &alpha)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize ChiDistribution: %w", err)
	}
	if !cd.SetParameters(float64(alpha), 0, 0) {
		return nil, fmt.Errorf("invalid parameters in serialized data: alpha=%d", alpha)
	}
	return cd, nil
}

// SampleChi generates a Chi-distributed variate using the square root of a Chi-square sample.
func SampleChi(gen rand.Source, alpha int) float64 {
	return math.Sqrt(SampleChiSquare(gen, alpha))
}
