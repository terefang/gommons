package distribution

import (
	"fmt"
	"math"
	"math/rand/v2"
)

// BetaPrimeDistribution represents a two-parameter distribution with range [0, +Inf).
type BetaPrimeDistribution struct {
	Generator rand.Source
	Alpha     float64
	Beta      float64
}

// Ensure BetaPrimeDistribution implements the Distribution interface.
var _ Distribution = (*BetaPrimeDistribution)(nil)

// NewBetaPrimeDistribution creates a BetaPrimeDistribution with the given generator, alpha, and beta.
func NewBetaPrimeDistribution(gen rand.Source, alpha, beta float64) (*BetaPrimeDistribution, error) {
	bpd := &BetaPrimeDistribution{
		Generator: gen,
	}
	if !bpd.SetParameters(alpha, beta, 0.0) {
		return nil, fmt.Errorf("given alpha (%f) and/or beta (%f) are invalid", alpha, beta)
	}
	return bpd, nil
}

func (bpd *BetaPrimeDistribution) GetTag() string {
	return "BetaPrime"
}

func (bpd *BetaPrimeDistribution) GetParameterA() float64 {
	return bpd.Alpha
}

func (bpd *BetaPrimeDistribution) GetParameterB() float64 {
	return bpd.Beta
}

func (bpd *BetaPrimeDistribution) GetParameterC() float64 {
	return math.NaN()
}

func (bpd *BetaPrimeDistribution) GetMaximum() float64 {
	return math.Inf(1)
}

func (bpd *BetaPrimeDistribution) GetMinimum() float64 {
	return 0.0
}

func (bpd *BetaPrimeDistribution) GetMean() float64 {
	return bpd.Alpha / (bpd.Beta - 1.0)
}

func (bpd *BetaPrimeDistribution) GetMedian() float64 {
	panic("Median is undefined for BetaPrime distribution")
}

func (bpd *BetaPrimeDistribution) GetMode() []float64 {
	return []float64{(bpd.Alpha - 1.0) / (bpd.Beta + 1.0)}
}

func (bpd *BetaPrimeDistribution) GetVariance() float64 {
	if bpd.Beta >= 2.0 {
		denom := bpd.Beta - 1.0
		return (bpd.Alpha * (bpd.Alpha + bpd.Beta - 1.0)) / ((denom * denom) * (bpd.Beta - 2.0))
	}
	panic("Variance cannot be determined for the given parameters")
}

func (bpd *BetaPrimeDistribution) SetParameters(a, b, c float64) bool {
	if a > 1.0 && b > 2.0 {
		bpd.Alpha = a
		bpd.Beta = b
		return true
	}
	return false
}

func (bpd *BetaPrimeDistribution) NextDouble() float64 {
	return SampleBetaPrime(bpd.Generator, bpd.Alpha, bpd.Beta)
}

func (bpd *BetaPrimeDistribution) Copy() Distribution {
	return &BetaPrimeDistribution{
		Generator: bpd.Generator,
		Alpha:     bpd.Alpha,
		Beta:      bpd.Beta,
	}
}

func (bpd *BetaPrimeDistribution) StringSerialize() string {
	return fmt.Sprintf("BetaPrime:%f:%f", bpd.Alpha, bpd.Beta)
}

func (bpd *BetaPrimeDistribution) StringDeserialize(data string) (Distribution, error) {
	var alpha, beta float64
	_, err := fmt.Sscanf(data, "BetaPrime:%f:%f", &alpha, &beta)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize BetaPrimeDistribution: %w", err)
	}
	if !bpd.SetParameters(alpha, beta, 0) {
		return nil, fmt.Errorf("invalid parameters in serialized data: alpha=%f, beta=%f", alpha, beta)
	}
	return bpd, nil
}

// SampleBetaPrime generates a Beta-prime distributed variate using a transformed Beta variate.
func SampleBetaPrime(gen rand.Source, alpha, beta float64) float64 {
	variate := SampleBeta(gen, alpha, beta)
	rev := 1.0 - variate

	const eps = 1.0 / (1 << 66) // 0x1p-66 non-zero check
	if math.Abs(rev) <= eps {
		return math.Inf(1)
	}
	return variate / rev
}
