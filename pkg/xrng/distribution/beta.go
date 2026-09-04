package distribution

import (
	"fmt"
	"math"
	"math/rand/v2"
)

// BetaDistribution represents a two-parameter Beta distribution with range (0, 1).
type BetaDistribution struct {
	Generator rand.Source
	Alpha     float64
	Beta      float64
}

// Ensure BetaDistribution implements the Distribution interface.
var _ Distribution = (*BetaDistribution)(nil)

// NewBetaDistribution creates a BetaDistribution with the given generator, alpha, and beta.
func NewBetaDistribution(gen rand.Source, alpha, beta float64) (*BetaDistribution, error) {
	bd := &BetaDistribution{
		Generator: gen,
	}
	if !bd.SetParameters(alpha, beta, 0.0) {
		return nil, fmt.Errorf("given alpha (%f) and/or beta (%f) are invalid", alpha, beta)
	}
	return bd, nil
}

func (bd *BetaDistribution) GetTag() string {
	return "Beta"
}

func (bd *BetaDistribution) GetParameterA() float64 {
	return bd.Alpha
}

func (bd *BetaDistribution) GetParameterB() float64 {
	return bd.Beta
}

func (bd *BetaDistribution) GetParameterC() float64 {
	return math.NaN()
}

func (bd *BetaDistribution) GetMaximum() float64 {
	return 1.0
}

func (bd *BetaDistribution) GetMinimum() float64 {
	return 0.0
}

func (bd *BetaDistribution) GetMean() float64 {
	return bd.Alpha / (bd.Alpha + bd.Beta)
}

func (bd *BetaDistribution) GetMedian() float64 {
	panic("Median is undefined for Beta distribution")
}

func (bd *BetaDistribution) GetMode() []float64 {
	const eps = 1.0 / (1 << 24) // 0x1p-24

	if bd.Alpha > 1 && bd.Beta > 1 {
		return []float64{(bd.Alpha - 1.0) / (bd.Alpha + bd.Beta - 2.0)}
	}
	if bd.Alpha < 1 && bd.Beta < 1 {
		return []float64{0.0, 1.0}
	}
	if (bd.Alpha < 1 && bd.Beta >= 1) || (math.Abs(bd.Alpha-1.0) <= eps && bd.Beta > 1) {
		return []float64{0.0}
	}
	if (bd.Alpha >= 1 && bd.Beta < 1) || (bd.Alpha > 1 && math.Abs(bd.Beta-1.0) <= eps) {
		return []float64{1.0}
	}

	panic("Mode cannot be determined for the given parameters")
}

func (bd *BetaDistribution) GetVariance() float64 {
	sum := bd.Alpha + bd.Beta
	return (bd.Alpha * bd.Beta) / (sum * sum * (sum + 1.0))
}

func (bd *BetaDistribution) SetParameters(a, b, c float64) bool {
	if a > 0.0 && b > 0.0 {
		bd.Alpha = a
		bd.Beta = b
		return true
	}
	return false
}

func (bd *BetaDistribution) NextDouble() float64 {
	return SampleBeta(bd.Generator, bd.Alpha, bd.Beta)
}

func (bd *BetaDistribution) Copy() Distribution {
	// Assumes generator type supports deep copy or can be cloned as needed.
	return &BetaDistribution{
		Generator: bd.Generator,
		Alpha:     bd.Alpha,
		Beta:      bd.Beta,
	}
}

func (bd *BetaDistribution) StringSerialize() string {
	return fmt.Sprintf("Beta:%f:%f", bd.Alpha, bd.Beta)
}

func (bd *BetaDistribution) StringDeserialize(data string) (Distribution, error) {
	var a, b float64
	_, err := fmt.Sscanf(data, "Beta:%f:%f", &a, &b)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize BetaDistribution: %w", err)
	}
	if !bd.SetParameters(a, b, 0) {
		return nil, fmt.Errorf("invalid parameters in serialized data: alpha=%f, beta=%f", a, b)
	}
	return bd, nil
}

// SampleBeta generates a beta-distributed variate using gamma random variables.
func SampleBeta(gen rand.Source, alpha, beta float64) float64 {
	x := SampleGamma(gen, alpha, 1.0)
	var t float64
	for {
		t = x + SampleGamma(gen, beta, 1.0)
		if math.Abs(t) > (1.0 / (1 << 66)) { // 0x1p-66 non-zero check
			break
		}
	}
	return x / t
}
