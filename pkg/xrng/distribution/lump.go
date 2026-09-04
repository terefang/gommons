package distribution

import (
	"fmt"
	"math"
	"math/rand/v2"
)

// LumpDistribution represents a two-parameter distribution with range [0, 1).
type LumpDistribution struct {
	Generator rand.Source
	Alpha     float64 // Affects skewness (lower = lower values, higher = higher values)
	Beta      float64 // Affects centrality vs. extremity
}

// Ensure LumpDistribution implements the Distribution interface.
var _ Distribution = (*LumpDistribution)(nil)

// NewLumpDistribution creates a LumpDistribution with the given generator, alpha, and beta.
func NewLumpDistribution(gen rand.Source, alpha, beta float64) (*LumpDistribution, error) {
	ld := &LumpDistribution{
		Generator: gen,
	}
	if !ld.SetParameters(alpha, beta, 0.0) {
		return nil, fmt.Errorf("given alpha (%f) and/or beta (%f) are invalid", alpha, beta)
	}
	return ld, nil
}

func (ld *LumpDistribution) GetTag() string {
	return "Lump"
}

func (ld *LumpDistribution) GetParameterA() float64 {
	return ld.Alpha
}

func (ld *LumpDistribution) GetParameterB() float64 {
	return ld.Beta
}

func (ld *LumpDistribution) GetParameterC() float64 {
	return math.NaN()
}

func (ld *LumpDistribution) GetMaximum() float64 {
	return 1.0
}

func (ld *LumpDistribution) GetMinimum() float64 {
	return 0.0
}

func (ld *LumpDistribution) GetMean() float64 {
	panic("Mean is not supported for Lump distribution")
}

func (ld *LumpDistribution) GetMedian() float64 {
	panic("Median is not supported for Lump distribution")
}

func (ld *LumpDistribution) GetMode() []float64 {
	panic("Mode is not supported for Lump distribution")
}

func (ld *LumpDistribution) GetVariance() float64 {
	panic("Variance is not supported for Lump distribution")
}

func (ld *LumpDistribution) SetParameters(a, b, c float64) bool {
	// Equivalent to a == a && b == b (checking for !math.IsNaN)
	if !math.IsNaN(a) && !math.IsNaN(b) {
		ld.Alpha = a
		ld.Beta = b
		return true
	}
	return false
}

func (ld *LumpDistribution) NextDouble() float64 {
	return SampleLump(ld.Generator, ld.Alpha, ld.Beta)
}

func (ld *LumpDistribution) Copy() Distribution {
	return &LumpDistribution{
		Generator: ld.Generator,
		Alpha:     ld.Alpha,
		Beta:      ld.Beta,
	}
}

func (ld *LumpDistribution) StringSerialize() string {
	return fmt.Sprintf("Lump:%f:%f", ld.Alpha, ld.Beta)
}

func (ld *LumpDistribution) StringDeserialize(data string) (Distribution, error) {
	var alpha, beta float64
	_, err := fmt.Sscanf(data, "Lump:%f:%f", &alpha, &beta)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize LumpDistribution: %w", err)
	}
	if !ld.SetParameters(alpha, beta, 0) {
		return nil, fmt.Errorf("invalid parameters in serialized data: alpha=%f, beta=%f", alpha, beta)
	}
	return ld, nil
}

// SampleLump generates a Lump-distributed variate in [0, 1) by mapping the atan2 angle (in turns) of two Gaussian samples.
func SampleLump(gen rand.Source, alpha, beta float64) float64 {
	r := rand.New(gen)
	y := r.NormFloat64() - alpha
	x := r.NormFloat64() + beta

	// Math.atan2 returns angle in radians [-Pi, Pi]
	// Convert to turns in [0, 1) matching digital's TrigTools.atan2Turns
	turns := math.Atan2(y, x) / (2.0 * math.Pi)
	if turns < 0 {
		turns += 1.0
	}
	return turns
}
