package distribution

import (
	"fmt"
	"math"
	"math/rand/v2"
)

// WeibullDistribution represents a two-parameter distribution with range [0, +Inf).
type WeibullDistribution struct {
	Generator rand.Source
	Alpha     float64 // Shape parameter k (alpha > 0)
	Lambda    float64 // Scale parameter lambda (lambda > 0)
}

// Ensure WeibullDistribution implements the Distribution interface.
var _ Distribution = (*WeibullDistribution)(nil)

// NewWeibullDistribution creates a WeibullDistribution with the given generator, alpha, and lambda.
func NewWeibullDistribution(gen rand.Source, alpha, lambda float64) (*WeibullDistribution, error) {
	wd := &WeibullDistribution{
		Generator: gen,
	}
	if !wd.SetParameters(alpha, lambda, 0.0) {
		return nil, fmt.Errorf("given alpha (%f) and/or lambda (%f) are invalid", alpha, lambda)
	}
	return wd, nil
}

func (wd *WeibullDistribution) GetTag() string {
	return "Weibull"
}

func (wd *WeibullDistribution) GetParameterA() float64 {
	return wd.Alpha
}

func (wd *WeibullDistribution) GetParameterB() float64 {
	return wd.Lambda
}

func (wd *WeibullDistribution) GetParameterC() float64 {
	return math.NaN()
}

func (wd *WeibullDistribution) GetMaximum() float64 {
	return math.Inf(1)
}

func (wd *WeibullDistribution) GetMinimum() float64 {
	return 0.0
}

func (wd *WeibullDistribution) GetMean() float64 {
	return wd.Lambda * math.Gamma(1.0+1.0/wd.Alpha)
}

func (wd *WeibullDistribution) GetMedian() float64 {
	return wd.Lambda * math.Pow(math.Ln2, 1.0/wd.Alpha)
}

func (wd *WeibullDistribution) GetMode() []float64 {
	if wd.Alpha >= 1.0 {
		return []float64{wd.Lambda * math.Pow(1.0-1.0/wd.Alpha, 1.0/wd.Alpha)}
	}
	panic("Mode cannot be determined for the given parameters")
}

func (wd *WeibullDistribution) GetVariance() float64 {
	mean := wd.GetMean()
	return (wd.Lambda*wd.Lambda)*math.Gamma(1.0+2.0/wd.Alpha) - (mean * mean)
}

func (wd *WeibullDistribution) SetParameters(a, b, c float64) bool {
	if a > 0.0 && b > 0.0 {
		wd.Alpha = a
		wd.Lambda = b
		return true
	}
	return false
}

func (wd *WeibullDistribution) NextDouble() float64 {
	return SampleWeibull(wd.Generator, wd.Alpha, wd.Lambda)
}

func (wd *WeibullDistribution) Copy() Distribution {
	return &WeibullDistribution{
		Generator: wd.Generator,
		Alpha:     wd.Alpha,
		Lambda:    wd.Lambda,
	}
}

func (wd *WeibullDistribution) StringSerialize() string {
	return fmt.Sprintf("Weibull:%f:%f", wd.Alpha, wd.Lambda)
}

func (wd *WeibullDistribution) StringDeserialize(data string) (Distribution, error) {
	var alpha, lambda float64
	_, err := fmt.Sscanf(data, "Weibull:%f:%f", &alpha, &lambda)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize WeibullDistribution: %w", err)
	}
	if !wd.SetParameters(alpha, lambda, 0) {
		return nil, fmt.Errorf("invalid parameters in serialized data: alpha=%f, lambda=%f", alpha, lambda)
	}
	return wd, nil
}

// SampleWeibull generates a Weibull-distributed variate using inverse transform sampling.
func SampleWeibull(gen rand.Source, alpha, lambda float64) float64 {
	r := rand.New(gen)
	// Equivalent to nextExclusiveDouble() in open interval (0, 1)
	uExclusive := (float64(r.Uint64()>>11)*1.1102230246251565e-16 + 5.551115123125782e-17)

	return lambda * math.Pow(-math.Log(uExclusive), 1.0/alpha)
}
