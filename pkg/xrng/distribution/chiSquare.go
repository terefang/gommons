package distribution

import (
	"fmt"
	"math"
	"math/rand/v2"
)

// ChiSquareDistribution represents a one-parameter distribution with range (0, +Inf).
type ChiSquareDistribution struct {
	Generator rand.Source
	Alpha     int
}

// Ensure ChiSquareDistribution implements the Distribution interface.
var _ Distribution = (*ChiSquareDistribution)(nil)

// NewChiSquareDistribution creates a ChiSquareDistribution with the given generator and alpha (degrees of freedom).
func NewChiSquareDistribution(gen rand.Source, alpha int) (*ChiSquareDistribution, error) {
	cs := &ChiSquareDistribution{
		Generator: gen,
	}
	if !cs.SetParameters(float64(alpha), 0.0, 0.0) {
		return nil, fmt.Errorf("given alpha (%d) is invalid", alpha)
	}
	return cs, nil
}

func (cs *ChiSquareDistribution) GetTag() string {
	return "ChiSquare"
}

func (cs *ChiSquareDistribution) GetParameterA() float64 {
	return float64(cs.Alpha)
}

func (cs *ChiSquareDistribution) GetParameterB() float64 {
	return math.NaN()
}

func (cs *ChiSquareDistribution) GetParameterC() float64 {
	return math.NaN()
}

func (cs *ChiSquareDistribution) GetMaximum() float64 {
	return math.Inf(1)
}

func (cs *ChiSquareDistribution) GetMinimum() float64 {
	return 0.0
}

func (cs *ChiSquareDistribution) GetMean() float64 {
	return float64(cs.Alpha)
}

func (cs *ChiSquareDistribution) GetMedian() float64 {
	a := float64(cs.Alpha)
	term := 1.0 - (2.0 / (9.0 * a))
	return a * term * term * term
}

func (cs *ChiSquareDistribution) GetMode() []float64 {
	if cs.Alpha >= 2 {
		return []float64{float64(cs.Alpha - 2)}
	}
	panic("Mode cannot be determined for the given parameters")
}

func (cs *ChiSquareDistribution) GetVariance() float64 {
	return 2.0 * float64(cs.Alpha)
}

func (cs *ChiSquareDistribution) SetParameters(a, b, c float64) bool {
	if a >= 1.0 {
		cs.Alpha = int(a)
		return true
	}
	return false
}

func (cs *ChiSquareDistribution) NextDouble() float64 {
	return SampleChiSquare(cs.Generator, cs.Alpha)
}

func (cs *ChiSquareDistribution) Copy() Distribution {
	return &ChiSquareDistribution{
		Generator: cs.Generator,
		Alpha:     cs.Alpha,
	}
}

func (cs *ChiSquareDistribution) StringSerialize() string {
	return fmt.Sprintf("ChiSquare:%d", cs.Alpha)
}

func (cs *ChiSquareDistribution) StringDeserialize(data string) (Distribution, error) {
	var alpha int
	_, err := fmt.Sscanf(data, "ChiSquare:%d", &alpha)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize ChiSquareDistribution: %w", err)
	}
	if !cs.SetParameters(float64(alpha), 0, 0) {
		return nil, fmt.Errorf("invalid parameters in serialized data: alpha=%d", alpha)
	}
	return cs, nil
}

// SampleChiSquare generates a Chi-square distributed variate by summing squared standard normal variates.
func SampleChiSquare(gen rand.Source, alpha int) float64 {
	sum := 0.0
	for i := 0; i < alpha; i++ {
		g := SampleNormal(gen, 0.0, 1.0)
		sum += g * g
	}
	return sum
}
