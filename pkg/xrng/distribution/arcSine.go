package distribution

import (
	"fmt"
	"math"
	"math/rand/v2"
)

// ArcsineDistribution represents a two-parameter distribution with range (Alpha, Beta) exclusive.
type ArcsineDistribution struct {
	Generator rand.Source
	Alpha     float64
	Beta      float64
}

// Ensure ArcsineDistribution implements the Distribution interface.
var _ Distribution = (*ArcsineDistribution)(nil)

func NewArcsineDistributionDefault(gen rand.Source) (*ArcsineDistribution, error) {
	return NewArcsineDistribution(gen, 0., 1.)
}

// NewArcsineDistribution creates an ArcsineDistribution with the given generator, alpha, and beta.
func NewArcsineDistribution(gen rand.Source, alpha, beta float64) (*ArcsineDistribution, error) {
	ad := &ArcsineDistribution{
		Generator: gen,
	}
	if !ad.SetParameters(alpha, beta, 0.0) {
		return nil, fmt.Errorf("given alpha (%f) and/or beta (%f) are invalid", alpha, beta)
	}
	return ad, nil
}

func (ad *ArcsineDistribution) GetTag() string {
	return "Arcsine"
}

func (ad *ArcsineDistribution) GetParameterA() float64 {
	return ad.Alpha
}

func (ad *ArcsineDistribution) GetParameterB() float64 {
	return ad.Beta
}

func (ad *ArcsineDistribution) GetParameterC() float64 {
	return math.NaN()
}

func (ad *ArcsineDistribution) GetMaximum() float64 {
	return ad.Beta
}

func (ad *ArcsineDistribution) GetMinimum() float64 {
	return ad.Alpha
}

func (ad *ArcsineDistribution) GetMean() float64 {
	return 0.5 * (ad.Alpha + ad.Beta)
}

func (ad *ArcsineDistribution) GetMedian() float64 {
	return 0.5 * (ad.Alpha + ad.Beta)
}

func (ad *ArcsineDistribution) GetMode() []float64 {
	return []float64{ad.Alpha, ad.Beta}
}

func (ad *ArcsineDistribution) GetVariance() float64 {
	diff := ad.Beta - ad.Alpha
	return diff * diff * 0.125
}

func (ad *ArcsineDistribution) SetParameters(a, b, c float64) bool {
	if a < b {
		ad.Alpha = a
		ad.Beta = b
		return true
	}
	return false
}

func (ad *ArcsineDistribution) NextDouble() float64 {
	return SampleArcsine(ad.Generator, ad.Alpha, ad.Beta)
}

func (ad *ArcsineDistribution) Copy() Distribution {
	return &ArcsineDistribution{
		Generator: ad.Generator,
		Alpha:     ad.Alpha,
		Beta:      ad.Beta,
	}
}

func (ad *ArcsineDistribution) StringSerialize() string {
	return fmt.Sprintf("Arcsine:%f:%f", ad.Alpha, ad.Beta)
}

func (ad *ArcsineDistribution) StringDeserialize(data string) (Distribution, error) {
	var alpha, beta float64
	_, err := fmt.Sscanf(data, "Arcsine:%f:%f", &alpha, &beta)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize ArcsineDistribution: %w", err)
	}
	if !ad.SetParameters(alpha, beta, 0) {
		return nil, fmt.Errorf("invalid parameters in serialized data: alpha=%f, beta=%f", alpha, beta)
	}
	return ad, nil
}

// SampleArcsine generates an Arcsine-distributed random variable bounded by alpha and beta.
func SampleArcsine(gen rand.Source, alpha, beta float64) float64 {
	r := rand.New(gen)
	// Equivalent to nextExclusiveDouble() in open interval (0, 1)
	uExclusive := (float64(r.Uint64()>>11)*1.1102230246251565e-16 + 5.551115123125782e-17)

	s := sinQuarterTurns(uExclusive)
	return alpha + (beta-alpha)*s*s
}

func sinQuarterTurns(quarterTurns float64) float64 {
	ceil := int64(math.Ceil(quarterTurns)) & -2
	quarterTurns -= float64(ceil)
	x2 := quarterTurns * quarterTurns
	x3 := quarterTurns * x2
	return (((11.0*quarterTurns - 3.0*x3) / (7.0 + x2)) * float64(1-(ceil&2)))
}
