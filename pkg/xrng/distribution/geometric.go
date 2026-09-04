package distribution

import (
	"fmt"
	"math"
	"math/rand/v2"
)

// GeometricDistribution represents a one-parameter discrete distribution with integer range from 1 to +Inf.
type GeometricDistribution struct {
	Generator rand.Source
	Alpha     float64 // Success probability parameter (p)
}

// Ensure GeometricDistribution implements the Distribution interface.
var _ Distribution = (*GeometricDistribution)(nil)

// NewGeometricDistribution creates a GeometricDistribution with the given generator and success probability (alpha).
func NewGeometricDistribution(gen rand.Source, alpha float64) (*GeometricDistribution, error) {
	gd := &GeometricDistribution{
		Generator: gen,
	}
	if !gd.SetParameters(alpha, 0.0, 0.0) {
		return nil, fmt.Errorf("given alpha (%f) is invalid", alpha)
	}
	return gd, nil
}

func (gd *GeometricDistribution) GetTag() string {
	return "Geometric"
}

func (gd *GeometricDistribution) GetParameterA() float64 {
	return gd.Alpha
}

func (gd *GeometricDistribution) GetParameterB() float64 {
	return math.NaN()
}

func (gd *GeometricDistribution) GetParameterC() float64 {
	return math.NaN()
}

func (gd *GeometricDistribution) GetMaximum() float64 {
	return math.Inf(1)
}

func (gd *GeometricDistribution) GetMinimum() float64 {
	return 1.0
}

func (gd *GeometricDistribution) GetMean() float64 {
	return 1.0 / gd.Alpha
}

func (gd *GeometricDistribution) GetMedian() float64 {
	panic("Median is undefined for Geometric distribution")
}

func (gd *GeometricDistribution) GetMode() []float64 {
	return []float64{1.0}
}

func (gd *GeometricDistribution) GetVariance() float64 {
	return (1.0 - gd.Alpha) / (gd.Alpha * gd.Alpha)
}

func (gd *GeometricDistribution) SetParameters(a, b, c float64) bool {
	if a > 0.0 && a <= 1.0 {
		gd.Alpha = a
		return true
	}
	return false
}

func (gd *GeometricDistribution) NextDouble() float64 {
	return SampleGeometric(gd.Generator, gd.Alpha)
}

func (gd *GeometricDistribution) Copy() Distribution {
	return &GeometricDistribution{
		Generator: gd.Generator,
		Alpha:     gd.Alpha,
	}
}

func (gd *GeometricDistribution) StringSerialize() string {
	return fmt.Sprintf("Geometric:%f", gd.Alpha)
}

func (gd *GeometricDistribution) StringDeserialize(data string) (Distribution, error) {
	var alpha float64
	_, err := fmt.Sscanf(data, "Geometric:%f", &alpha)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize GeometricDistribution: %w", err)
	}
	if !gd.SetParameters(alpha, 0, 0) {
		return nil, fmt.Errorf("invalid parameters in serialized data: alpha=%f", alpha)
	}
	return gd, nil
}

// SampleGeometric generates a geometrically distributed trial count until the first success using direct simulation.
func SampleGeometric(gen rand.Source, alpha float64) float64 {
	r := rand.New(gen)
	samples := 1.0
	for r.Float64() >= alpha {
		samples++
	}
	return samples
}
