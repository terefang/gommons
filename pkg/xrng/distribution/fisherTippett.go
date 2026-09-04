package distribution

import (
	"fmt"
	"math"
	"math/rand/v2"
)

// FisherTippettDistribution represents a two-parameter Gumbel / Type-I Extreme Value distribution with infinite range (-Inf, +Inf).
type FisherTippettDistribution struct {
	Generator rand.Source
	Alpha     float64 // Scale parameter
	Mu        float64 // Location parameter
}

// Ensure FisherTippettDistribution implements the Distribution interface.
var _ Distribution = (*FisherTippettDistribution)(nil)

// NewFisherTippettDistribution creates a FisherTippettDistribution with the given generator, alpha (scale), and mu (location).
func NewFisherTippettDistribution(gen rand.Source, alpha, mu float64) (*FisherTippettDistribution, error) {
	ftd := &FisherTippettDistribution{
		Generator: gen,
	}
	if !ftd.SetParameters(alpha, mu, 0.0) {
		return nil, fmt.Errorf("given alpha (%f) and/or mu (%f) are invalid", alpha, mu)
	}
	return ftd, nil
}

func (ftd *FisherTippettDistribution) GetTag() string {
	return "FisherTippett"
}

func (ftd *FisherTippettDistribution) GetParameterA() float64 {
	return ftd.Alpha
}

func (ftd *FisherTippettDistribution) GetParameterB() float64 {
	return ftd.Mu
}

func (ftd *FisherTippettDistribution) GetParameterC() float64 {
	return math.NaN()
}

func (ftd *FisherTippettDistribution) GetMaximum() float64 {
	return math.Inf(1)
}

func (ftd *FisherTippettDistribution) GetMinimum() float64 {
	return math.Inf(-1)
}

func (ftd *FisherTippettDistribution) GetMean() float64 {
	// mu + alpha * Euler-Mascheroni constant
	return ftd.Mu + ftd.Alpha*0.577215664901532860606512090082402431042159335
}

func (ftd *FisherTippettDistribution) GetMedian() float64 {
	// mu - alpha * ln(ln(2))
	return ftd.Mu - ftd.Alpha*-0.36651292058166435
}

func (ftd *FisherTippettDistribution) GetMode() []float64 {
	return []float64{ftd.Mu}
}

func (ftd *FisherTippettDistribution) GetVariance() float64 {
	return ((math.Pi * math.Pi) / 6.0) * ftd.Alpha * ftd.Alpha
}

func (ftd *FisherTippettDistribution) SetParameters(a, b, c float64) bool {
	if a > 0.0 && !math.IsNaN(b) {
		ftd.Alpha = a
		ftd.Mu = b
		return true
	}
	return false
}

func (ftd *FisherTippettDistribution) NextDouble() float64 {
	return SampleFisherTippett(ftd.Generator, ftd.Alpha, ftd.Mu)
}

func (ftd *FisherTippettDistribution) Copy() Distribution {
	return &FisherTippettDistribution{
		Generator: ftd.Generator,
		Alpha:     ftd.Alpha,
		Mu:        ftd.Mu,
	}
}

func (ftd *FisherTippettDistribution) StringSerialize() string {
	return fmt.Sprintf("FisherTippett:%f:%f", ftd.Alpha, ftd.Mu)
}

func (ftd *FisherTippettDistribution) StringDeserialize(data string) (Distribution, error) {
	var alpha, mu float64
	_, err := fmt.Sscanf(data, "FisherTippett:%f:%f", &alpha, &mu)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize FisherTippettDistribution: %w", err)
	}
	if !ftd.SetParameters(alpha, mu, 0) {
		return nil, fmt.Errorf("invalid parameters in serialized data: alpha=%f, mu=%f", alpha, mu)
	}
	return ftd, nil
}

// SampleFisherTippett generates a Fisher-Tippett (Gumbel) variate using inverse transform sampling.
func SampleFisherTippett(gen rand.Source, alpha, mu float64) float64 {
	r := rand.New(gen)
	// Equivalent to nextExclusiveDouble() in open interval (0, 1)
	uExclusive := (float64(r.Uint64()>>11)*1.1102230246251565e-16 + 5.551115123125782e-17)

	return mu - alpha*math.Log(-math.Log(1.0-uExclusive))
}
