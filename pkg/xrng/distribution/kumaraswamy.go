package distribution

import (
	"fmt"
	"math"
	"math/rand/v2"
)

// KumaraswamyDistribution represents a two-parameter distribution bounded on [0, 1].
type KumaraswamyDistribution struct {
	Generator rand.Source
	Alpha     float64 // Stored internally as 1.0 / alpha
	Beta      float64 // Stored internally as 1.0 / beta
}

// Ensure KumaraswamyDistribution implements the Distribution interface.
var _ Distribution = (*KumaraswamyDistribution)(nil)

// NewKumaraswamyDistribution creates a KumaraswamyDistribution with the given generator, alpha, and beta.
func NewKumaraswamyDistribution(gen rand.Source, alpha, beta float64) (*KumaraswamyDistribution, error) {
	kd := &KumaraswamyDistribution{
		Generator: gen,
	}
	if !kd.SetParameters(alpha, beta, 0.0) {
		return nil, fmt.Errorf("given alpha (%f) and/or beta (%f) are invalid", alpha, beta)
	}
	return kd, nil
}

func (kd *KumaraswamyDistribution) GetTag() string {
	return "Kumaraswamy"
}

func (kd *KumaraswamyDistribution) GetParameterA() float64 {
	return 1.0 / kd.Alpha
}

func (kd *KumaraswamyDistribution) GetParameterB() float64 {
	return 1.0 / kd.Beta
}

func (kd *KumaraswamyDistribution) GetParameterC() float64 {
	return math.NaN()
}

func (kd *KumaraswamyDistribution) GetMaximum() float64 {
	return 1.0
}

func (kd *KumaraswamyDistribution) GetMinimum() float64 {
	return 0.0
}

func (kd *KumaraswamyDistribution) GetMean() float64 {
	b := 1.0 / kd.Beta
	a := 1.0 / kd.Alpha

	// MathTools.factorial(a) for real a is Gamma(a + 1)
	gammaA1, _ := math.Lgamma(a + 1.0)
	gammaB, _ := math.Lgamma(b)
	gammaAB1, _ := math.Lgamma(a + b + 1.0)

	logMean := gammaA1 + gammaB + math.Log(b) - gammaAB1
	return math.Exp(logMean)
}

func (kd *KumaraswamyDistribution) GetMedian() float64 {
	return math.Pow(1.0-math.Pow(2.0, -kd.Beta), kd.Alpha)
}

func (kd *KumaraswamyDistribution) GetMode() []float64 {
	origA := 1.0 / kd.Alpha
	origB := 1.0 / kd.Beta
	if origA >= 1.0 && origB >= 1.0 && !(origA == 1.0 && origB == 1.0) {
		mode := math.Pow((kd.Alpha-1.0)/(kd.Alpha*kd.Beta-1.0), kd.Alpha)
		return []float64{mode}
	}
	panic("Mode cannot be determined for the given parameters")
}

func (kd *KumaraswamyDistribution) GetVariance() float64 {
	panic("Variance is not supported for Kumaraswamy distribution")
}

func (kd *KumaraswamyDistribution) SetParameters(a, b, c float64) bool {
	if a > 0.0 && b > 0.0 {
		kd.Alpha = 1.0 / a
		kd.Beta = 1.0 / b
		return true
	}
	return false
}

func (kd *KumaraswamyDistribution) NextDouble() float64 {
	return SampleKumaraswamy(kd.Generator, kd.Alpha, kd.Beta)
}

func (kd *KumaraswamyDistribution) Copy() Distribution {
	return &KumaraswamyDistribution{
		Generator: kd.Generator,
		Alpha:     kd.Alpha,
		Beta:      kd.Beta,
	}
}

func (kd *KumaraswamyDistribution) StringSerialize() string {
	return fmt.Sprintf("Kumaraswamy:%f:%f", 1.0/kd.Alpha, 1.0/kd.Beta)
}

func (kd *KumaraswamyDistribution) StringDeserialize(data string) (Distribution, error) {
	var alpha, beta float64
	_, err := fmt.Sscanf(data, "Kumaraswamy:%f:%f", &alpha, &beta)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize KumaraswamyDistribution: %w", err)
	}
	if !kd.SetParameters(alpha, beta, 0) {
		return nil, fmt.Errorf("invalid parameters in serialized data: alpha=%f, beta=%f", alpha, beta)
	}
	return kd, nil
}

// SampleKumaraswamy generates a Kumaraswamy-distributed variate using inverse transform sampling.
func SampleKumaraswamy(gen rand.Source, inverseAlpha, inverseBeta float64) float64 {
	r := rand.New(gen)
	// Equivalent to nextExclusiveDouble() in open interval (0, 1)
	uExclusive := (float64(r.Uint64()>>11)*1.1102230246251565e-16 + 5.551115123125782e-17)

	return math.Pow(1.0-math.Pow(uExclusive, inverseBeta), inverseAlpha)
}
