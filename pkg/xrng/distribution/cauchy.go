package distribution

import (
	"fmt"
	"math"
	"math/rand/v2"
)

// CauchyDistribution represents a two-parameter Cauchy distribution with infinite range (-Inf, +Inf).
type CauchyDistribution struct {
	Generator rand.Source
	Alpha     float64 // Location parameter
	Gamma     float64 // Scale parameter
}

// Ensure CauchyDistribution implements the Distribution interface.
var _ Distribution = (*CauchyDistribution)(nil)

// NewCauchyDistribution creates a CauchyDistribution with the given generator, alpha (location), and gamma (scale).
func NewCauchyDistribution(gen rand.Source, alpha, gamma float64) (*CauchyDistribution, error) {
	cd := &CauchyDistribution{
		Generator: gen,
	}
	if !cd.SetParameters(alpha, gamma, 0.0) {
		return nil, fmt.Errorf("given alpha (%f) and/or gamma (%f) are invalid", alpha, gamma)
	}
	return cd, nil
}

func (cd *CauchyDistribution) GetTag() string {
	return "Cauchy"
}

func (cd *CauchyDistribution) GetParameterA() float64 {
	return cd.Alpha
}

func (cd *CauchyDistribution) GetParameterB() float64 {
	return cd.Gamma
}

func (cd *CauchyDistribution) GetParameterC() float64 {
	return math.NaN()
}

func (cd *CauchyDistribution) GetMaximum() float64 {
	return math.Inf(1)
}

func (cd *CauchyDistribution) GetMinimum() float64 {
	return math.Inf(-1)
}

func (cd *CauchyDistribution) GetMean() float64 {
	panic("Mean is undefined for Cauchy distribution")
}

func (cd *CauchyDistribution) GetMedian() float64 {
	return cd.Alpha
}

func (cd *CauchyDistribution) GetMode() []float64 {
	return []float64{cd.Alpha}
}

func (cd *CauchyDistribution) GetVariance() float64 {
	panic("Variance is undefined for Cauchy distribution")
}

func (cd *CauchyDistribution) SetParameters(a, b, c float64) bool {
	if !math.IsNaN(a) && b > 0.0 {
		cd.Alpha = a
		cd.Gamma = b
		return true
	}
	return false
}

func (cd *CauchyDistribution) NextDouble() float64 {
	return SampleCauchy(cd.Generator, cd.Alpha, cd.Gamma)
}

func (cd *CauchyDistribution) Copy() Distribution {
	return &CauchyDistribution{
		Generator: cd.Generator,
		Alpha:     cd.Alpha,
		Gamma:     cd.Gamma,
	}
}

func (cd *CauchyDistribution) StringSerialize() string {
	return fmt.Sprintf("Cauchy:%f:%f", cd.Alpha, cd.Gamma)
}

func (cd *CauchyDistribution) StringDeserialize(data string) (Distribution, error) {
	var alpha, gamma float64
	_, err := fmt.Sscanf(data, "Cauchy:%f:%f", &alpha, &gamma)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize CauchyDistribution: %w", err)
	}
	if !cd.SetParameters(alpha, gamma, 0) {
		return nil, fmt.Errorf("invalid parameters in serialized data: alpha=%f, gamma=%f", alpha, gamma)
	}
	return cd, nil
}

// SampleCauchy generates a Cauchy-distributed variate using inverse transform sampling (tan(pi * (u - 0.5))).
func SampleCauchy(gen rand.Source, alpha, gamma float64) float64 {
	r := rand.New(gen)
	// Equivalent to nextExclusiveDouble() in open interval (0, 1)
	uExclusive := (float64(r.Uint64()>>11)*1.1102230246251565e-16 + 5.551115123125782e-17)

	return alpha + gamma*math.Tan(math.Pi*(uExclusive-0.5))
}
