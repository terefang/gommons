package distribution

import (
	"fmt"
	"math"
	"math/rand/v2"
)

// BinomialDistribution represents a two-parameter discrete distribution with integer range [0, Beta].
type BinomialDistribution struct {
	Generator rand.Source
	Alpha     float64
	Beta      int
}

// Ensure BinomialDistribution implements the Distribution interface.
var _ Distribution = (*BinomialDistribution)(nil)

// NewBinomialDistribution creates a BinomialDistribution with the given generator, alpha, and beta.
func NewBinomialDistribution(gen rand.Source, alpha float64, beta int) (*BinomialDistribution, error) {
	bd := &BinomialDistribution{
		Generator: gen,
	}
	if !bd.SetParameters(alpha, float64(beta), 0.0) {
		return nil, fmt.Errorf("given alpha (%f) and/or beta (%d) are invalid", alpha, beta)
	}
	return bd, nil
}

func (bd *BinomialDistribution) GetTag() string {
	return "Binomial"
}

func (bd *BinomialDistribution) GetParameterA() float64 {
	return bd.Alpha
}

func (bd *BinomialDistribution) GetParameterB() float64 {
	return float64(bd.Beta)
}

func (bd *BinomialDistribution) GetParameterC() float64 {
	return math.NaN()
}

func (bd *BinomialDistribution) GetMaximum() float64 {
	return float64(bd.Beta)
}

func (bd *BinomialDistribution) GetMinimum() float64 {
	return 0.0
}

func (bd *BinomialDistribution) GetMean() float64 {
	return float64(bd.Beta) * bd.Alpha
}

func (bd *BinomialDistribution) GetMedian() float64 {
	panic("Median is undefined for Binomial distribution")
}

func (bd *BinomialDistribution) GetMode() []float64 {
	return []float64{math.Floor(bd.Alpha * float64(bd.Beta+1))}
}

func (bd *BinomialDistribution) GetVariance() float64 {
	return bd.Alpha * (1.0 - bd.Alpha) * float64(bd.Beta)
}

func (bd *BinomialDistribution) SetParameters(a, b, c float64) bool {
	if a >= 0.0 && a <= 1.0 && int(b) >= 0 {
		bd.Alpha = a
		bd.Beta = int(b)
		return true
	}
	return false
}

func (bd *BinomialDistribution) NextDouble() float64 {
	return SampleBinomial(bd.Generator, bd.Alpha, bd.Beta)
}

func (bd *BinomialDistribution) Copy() Distribution {
	return &BinomialDistribution{
		Generator: bd.Generator,
		Alpha:     bd.Alpha,
		Beta:      bd.Beta,
	}
}

func (bd *BinomialDistribution) StringSerialize() string {
	return fmt.Sprintf("Binomial:%f:%d", bd.Alpha, bd.Beta)
}

func (bd *BinomialDistribution) StringDeserialize(data string) (Distribution, error) {
	var alpha float64
	var beta int
	_, err := fmt.Sscanf(data, "Binomial:%f:%d", &alpha, &beta)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize BinomialDistribution: %w", err)
	}
	if !bd.SetParameters(alpha, float64(beta), 0) {
		return nil, fmt.Errorf("invalid parameters in serialized data: alpha=%f, beta=%d", alpha, beta)
	}
	return bd, nil
}

// SampleBinomial generates a binomially distributed variate counting successes in `beta` Bernoulli trials.
func SampleBinomial(gen rand.Source, alpha float64, beta int) float64 {
	r := rand.New(gen)
	successes := 0.0
	for i := 0; i < beta; i++ {
		// Equivalent to nextExclusiveDouble() in open interval (0, 1)
		uExclusive := (float64(r.Uint64()>>11)*1.1102230246251565e-16 + 5.551115123125782e-17)
		if uExclusive < alpha {
			successes++
		}
	}
	return successes
}
