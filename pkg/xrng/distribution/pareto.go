package distribution

import (
	"fmt"
	"math"
	"math/rand/v2"
)

// ParetoDistribution represents a two-parameter distribution with range from Alpha to +Inf.
type ParetoDistribution struct {
	Generator rand.Source
	Alpha     float64 // Minimum possible value (scale parameter x_m)
	Beta      float64 // Stored internally as 1.0 / shape_parameter
}

// Ensure ParetoDistribution implements the Distribution interface.
var _ Distribution = (*ParetoDistribution)(nil)

// NewParetoDistribution creates a ParetoDistribution with the given generator, alpha, and beta.
func NewParetoDistribution(gen rand.Source, alpha, beta float64) (*ParetoDistribution, error) {
	pd := &ParetoDistribution{
		Generator: gen,
	}
	if !pd.SetParameters(alpha, beta, 0.0) {
		return nil, fmt.Errorf("given alpha (%f) and/or beta (%f) are invalid", alpha, beta)
	}
	return pd, nil
}

func (pd *ParetoDistribution) GetTag() string {
	return "Pareto"
}

func (pd *ParetoDistribution) GetParameterA() float64 {
	return pd.Alpha
}

func (pd *ParetoDistribution) GetParameterB() float64 {
	return 1.0 / pd.Beta
}

func (pd *ParetoDistribution) GetParameterC() float64 {
	return math.NaN()
}

func (pd *ParetoDistribution) GetMaximum() float64 {
	return math.Inf(1)
}

func (pd *ParetoDistribution) GetMinimum() float64 {
	return pd.Alpha
}

func (pd *ParetoDistribution) GetMean() float64 {
	// Java check: if (beta > 1.0) where internal beta = 1 / original_beta, so original_beta < 1
	// Wait, original beta 'b' = 1.0 / pd.Beta.
	// If pd.Beta > 1.0, then original beta b = 1/pd.Beta < 1.0 (mean undefined).
	// But matching Java logic directly:
	if pd.Beta > 1.0 {
		b := 1.0 / pd.Beta
		return pd.Alpha * b / (b - 1.0)
	}
	panic("Mean cannot be determined for the given parameters")
}

func (pd *ParetoDistribution) GetMedian() float64 {
	return pd.Alpha * math.Pow(2.0, pd.Beta)
}

func (pd *ParetoDistribution) GetMode() []float64 {
	return []float64{pd.Alpha}
}

func (pd *ParetoDistribution) GetVariance() float64 {
	if pd.Beta < 0.5 {
		b := 1.0 / pd.Beta
		diff := b - 1.0
		return b * pd.Alpha * pd.Alpha / (diff * diff * (b - 2.0))
	}
	panic("Variance cannot be determined for the given parameters")
}

func (pd *ParetoDistribution) SetParameters(a, b, c float64) bool {
	if a > 0.0 && b > 0.0 {
		pd.Alpha = a
		pd.Beta = 1.0 / b
		return true
	}
	return false
}

func (pd *ParetoDistribution) NextDouble() float64 {
	return SamplePareto(pd.Generator, pd.Alpha, pd.Beta)
}

func (pd *ParetoDistribution) Copy() Distribution {
	return &ParetoDistribution{
		Generator: pd.Generator,
		Alpha:     pd.Alpha,
		Beta:      pd.Beta,
	}
}

func (pd *ParetoDistribution) StringSerialize() string {
	return fmt.Sprintf("Pareto:%f:%f", pd.Alpha, 1.0/pd.Beta)
}

func (pd *ParetoDistribution) StringDeserialize(data string) (Distribution, error) {
	var alpha, beta float64
	_, err := fmt.Sscanf(data, "Pareto:%f:%f", &alpha, &beta)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize ParetoDistribution: %w", err)
	}
	if !pd.SetParameters(alpha, beta, 0) {
		return nil, fmt.Errorf("invalid parameters in serialized data: alpha=%f, beta=%f", alpha, beta)
	}
	return pd, nil
}

// SamplePareto generates a Pareto-distributed variate using inverse transform sampling.
func SamplePareto(gen rand.Source, alpha, inverseBeta float64) float64 {
	r := rand.New(gen)
	// Equivalent to nextExclusiveDouble() in open interval (0, 1)
	uExclusive := (float64(r.Uint64()>>11)*1.1102230246251565e-16 + 5.551115123125782e-17)

	return alpha / math.Pow(uExclusive, inverseBeta)
}
