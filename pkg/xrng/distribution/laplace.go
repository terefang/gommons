package distribution

import (
	"fmt"
	"math"
	"math/rand/v2"
)

// LaplaceDistribution represents a two-parameter distribution with infinite range (-Inf, +Inf).
type LaplaceDistribution struct {
	Generator rand.Source
	Alpha     float64 // Scale parameter (b > 0)
	Mu        float64 // Location parameter
}

// Ensure LaplaceDistribution implements the Distribution interface.
var _ Distribution = (*LaplaceDistribution)(nil)

// NewLaplaceDistribution creates a LaplaceDistribution with the given generator, alpha (scale), and mu (location).
func NewLaplaceDistribution(gen rand.Source, alpha, mu float64) (*LaplaceDistribution, error) {
	ld := &LaplaceDistribution{
		Generator: gen,
	}
	if !ld.SetParameters(alpha, mu, 0.0) {
		return nil, fmt.Errorf("given alpha (%f) and/or mu (%f) are invalid", alpha, mu)
	}
	return ld, nil
}

func (ld *LaplaceDistribution) GetTag() string {
	return "Laplace"
}

func (ld *LaplaceDistribution) GetParameterA() float64 {
	return ld.Alpha
}

func (ld *LaplaceDistribution) GetParameterB() float64 {
	return ld.Mu
}

func (ld *LaplaceDistribution) GetParameterC() float64 {
	return math.NaN()
}

func (ld *LaplaceDistribution) GetMaximum() float64 {
	return math.Inf(1)
}

func (ld *LaplaceDistribution) GetMinimum() float64 {
	return math.Inf(-1)
}

func (ld *LaplaceDistribution) GetMean() float64 {
	return ld.Mu
}

func (ld *LaplaceDistribution) GetMedian() float64 {
	return ld.Mu
}

func (ld *LaplaceDistribution) GetMode() []float64 {
	return []float64{ld.Mu}
}

func (ld *LaplaceDistribution) GetVariance() float64 {
	return 2.0 * ld.Alpha * ld.Alpha
}

func (ld *LaplaceDistribution) SetParameters(a, b, c float64) bool {
	if a > 0.0 && !math.IsNaN(b) {
		ld.Alpha = a
		ld.Mu = b
		return true
	}
	return false
}

func (ld *LaplaceDistribution) NextDouble() float64 {
	return SampleLaplace(ld.Generator, ld.Alpha, ld.Mu)
}

func (ld *LaplaceDistribution) Copy() Distribution {
	return &LaplaceDistribution{
		Generator: ld.Generator,
		Alpha:     ld.Alpha,
		Mu:        ld.Mu,
	}
}

func (ld *LaplaceDistribution) StringSerialize() string {
	return fmt.Sprintf("Laplace:%f:%f", ld.Alpha, ld.Mu)
}

func (ld *LaplaceDistribution) StringDeserialize(data string) (Distribution, error) {
	var alpha, mu float64
	_, err := fmt.Sscanf(data, "Laplace:%f:%f", &alpha, &mu)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize LaplaceDistribution: %w", err)
	}
	if !ld.SetParameters(alpha, mu, 0) {
		return nil, fmt.Errorf("invalid parameters in serialized data: alpha=%f, mu=%f", alpha, mu)
	}
	return ld, nil
}

// SampleLaplace generates a Laplace-distributed variate using inverse transform sampling.
func SampleLaplace(gen rand.Source, alpha, mu float64) float64 {
	r := rand.New(gen)
	// Equivalent to nextExclusiveDouble() in open interval (0, 1)
	uExclusive := (float64(r.Uint64()>>11)*1.1102230246251565e-16 + 5.551115123125782e-17)

	randVal := 0.5 - uExclusive

	var tmp float64
	// Check if randVal is close to 0 using tolerance 0x1p-66 (approx 1.8425e-20)
	if math.Abs(randVal) <= 1.8425031308300624e-20 {
		tmp = math.Inf(-1)
	} else {
		tmp = math.Log(2.0 * math.Abs(randVal))
	}

	var sign float64
	if randVal > 0 {
		sign = 1.0
	} else if randVal < 0 {
		sign = -1.0
	} else {
		sign = 0.0
	}

	return mu - alpha*sign*tmp
}
