package distribution

import (
	"fmt"
	"math"
	"math/rand/v2"
)

// GammaDistribution represents a two-parameter Gamma distribution with range (0, +Inf).
type GammaDistribution struct {
	Generator rand.Source
	Alpha     float64
	Beta      float64
}

// Ensure GammaDistribution implements the Distribution interface.
var _ Distribution = (*GammaDistribution)(nil)

// NewGammaDistribution creates a GammaDistribution with the given generator, alpha, and beta.
func NewGammaDistribution(gen rand.Source, alpha, beta float64) (*GammaDistribution, error) {
	gd := &GammaDistribution{
		Generator: gen,
	}
	if !gd.SetParameters(alpha, beta, 0.0) {
		return nil, fmt.Errorf("given alpha (%f) and/or beta (%f) are invalid", alpha, beta)
	}
	return gd, nil
}

func (gd *GammaDistribution) GetTag() string {
	return "Gamma"
}

func (gd *GammaDistribution) GetParameterA() float64 {
	return gd.Alpha
}

func (gd *GammaDistribution) GetParameterB() float64 {
	return gd.Beta
}

func (gd *GammaDistribution) GetParameterC() float64 {
	return math.NaN()
}

func (gd *GammaDistribution) GetMaximum() float64 {
	return math.Inf(1)
}

func (gd *GammaDistribution) GetMinimum() float64 {
	return 0.0
}

func (gd *GammaDistribution) GetMean() float64 {
	return gd.Alpha / gd.Beta
}

func (gd *GammaDistribution) GetMedian() float64 {
	panic("Median is undefined for Gamma distribution")
}

func (gd *GammaDistribution) GetMode() []float64 {
	if gd.Alpha >= 1.0 {
		return []float64{(gd.Alpha - 1.0) / gd.Beta}
	}
	panic("Mode cannot be determined for the given parameters")
}

func (gd *GammaDistribution) GetVariance() float64 {
	return gd.Alpha / (gd.Beta * gd.Beta)
}

func (gd *GammaDistribution) SetParameters(a, b, c float64) bool {
	if a > 0.0 && b > 0.0 {
		gd.Alpha = a
		gd.Beta = b
		return true
	}
	return false
}

func (gd *GammaDistribution) NextDouble() float64 {
	return SampleGamma(gd.Generator, gd.Alpha, gd.Beta)
}

func (gd *GammaDistribution) Copy() Distribution {
	return &GammaDistribution{
		Generator: gd.Generator,
		Alpha:     gd.Alpha,
		Beta:      gd.Beta,
	}
}

func (gd *GammaDistribution) StringSerialize() string {
	return fmt.Sprintf("Gamma:%f:%f", gd.Alpha, gd.Beta)
}

func (gd *GammaDistribution) StringDeserialize(data string) (Distribution, error) {
	var alpha, beta float64
	_, err := fmt.Sscanf(data, "Gamma:%f:%f", &alpha, &beta)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize GammaDistribution: %w", err)
	}
	if !gd.SetParameters(alpha, beta, 0) {
		return nil, fmt.Errorf("invalid parameters in serialized data: alpha=%f, beta=%f", alpha, beta)
	}
	return gd, nil
}

// SampleGamma generates a gamma-distributed variate using Marsaglia and Tsang's method.
func SampleGamma(gen rand.Source, alpha, beta float64) float64 {
	oalpha := alpha
	if alpha < 1.0 {
		alpha += 1.0
	}
	a1 := alpha - 1.0/3.0
	a2 := 1.0 / math.Sqrt(9.0*a1)

	r := rand.New(gen)

	var u, v, x float64
	for {
		for {
			x = SampleNormal(gen, 0.0, 1.0)
			v = 1.0 + a2*x
			if v > 0.0 {
				break
			}
		}

		v = v * v * v
		u = r.Float64()
		x2 := x * x

		if u <= (1.0-0.331*x2*x2) || math.Log(u) <= (0.5*x2+a1*(1.0-v+math.Log(v))) {
			break
		}
	}

	const eps = 1.0 / (1 << 24) // 0x1p-24
	if math.Abs(alpha-oalpha) <= eps {
		return a1 * v / beta
	}

	// Equivalent to nextExclusiveDouble() in open interval (0, 1)
	uExclusive := (float64(r.Uint64()>>11)*1.1102230246251565e-16 + 5.551115123125782e-17)
	return math.Pow(uExclusive, 1.0/oalpha) * a1 * v / beta
}
