package distribution

import (
	"fmt"
	"math"
	"math/rand/v2"
)

// ZipfianDistribution represents a discrete two-parameter distribution with integer range from 1 to Alpha.
type ZipfianDistribution struct {
	Generator   rand.Source
	Alpha       float64 // Maximum integer bound (stored as float64 to match interface, cast from long)
	Skew        float64 // Skew parameter (0 <= Skew < 1)
	Zeta        float64 // Precalculated generalized harmonic number sum
	ZetaTwoSkew float64 // Internal optimization parameter: 1 + 0.5^skew
}

// Ensure ZipfianDistribution implements the Distribution interface.
var _ Distribution = (*ZipfianDistribution)(nil)

// NewZipfianDistribution creates a ZipfianDistribution with the given generator, alpha, and skew.
func NewZipfianDistribution(gen rand.Source, alpha int64, skew float64) (*ZipfianDistribution, error) {
	zd := &ZipfianDistribution{
		Generator: gen,
	}
	if !zd.SetParameters(float64(alpha), skew, 0.0) {
		return nil, fmt.Errorf("given alpha (%d) and/or skew (%f) are invalid", alpha, skew)
	}
	return zd, nil
}

// NewZipfianDistributionWithZeta creates a ZipfianDistribution with precalculated zeta for performance optimization.
func NewZipfianDistributionWithZeta(gen rand.Source, alpha int64, skew, zeta float64) (*ZipfianDistribution, error) {
	zd := &ZipfianDistribution{
		Generator: gen,
	}
	if !zd.SetParameters(float64(alpha), skew, -1.0) {
		return nil, fmt.Errorf("given alpha (%d) and/or skew (%f) are invalid", alpha, skew)
	}
	zd.Zeta = zeta
	return zd, nil
}

func (zd *ZipfianDistribution) GetTag() string {
	return "Zipfian"
}

func (zd *ZipfianDistribution) GetParameterA() float64 {
	return zd.Alpha
}

func (zd *ZipfianDistribution) GetParameterB() float64 {
	return zd.Skew
}

func (zd *ZipfianDistribution) GetParameterC() float64 {
	return math.NaN()
}

func (zd *ZipfianDistribution) GetMaximum() float64 {
	return zd.Alpha
}

func (zd *ZipfianDistribution) GetMinimum() float64 {
	return 1.0
}

func (zd *ZipfianDistribution) GetMean() float64 {
	if zd.Skew > 1.0 {
		return Harmonic(int64(zd.Alpha), zd.Skew-1.0) / zd.Zeta
	}
	panic("Mean cannot be determined for the given parameters")
}

func (zd *ZipfianDistribution) GetMedian() float64 {
	panic("Median cannot be determined")
}

func (zd *ZipfianDistribution) GetMode() []float64 {
	return []float64{1.0}
}

func (zd *ZipfianDistribution) GetVariance() float64 {
	if zd.Skew > 2.0 {
		mean := zd.GetMean()
		return Harmonic(int64(zd.Alpha), zd.Skew-2.0)/zd.Zeta - (mean * mean)
	}
	panic("Variance cannot be determined for the given parameters")
}

func (zd *ZipfianDistribution) SetParameters(a, b, c float64) bool {
	if a >= 1.0 && b >= 0.0 && b < 1.0 {
		zd.Alpha = float64(int64(a))
		zd.Skew = b
		if math.IsNaN(c) || c >= 0.0 {
			zd.Zeta = Harmonic(int64(zd.Alpha), zd.Skew)
		}
		zd.ZetaTwoSkew = 1.0 + math.Pow(0.5, zd.Skew)
		return true
	}
	return false
}

func (zd *ZipfianDistribution) NextDouble() float64 {
	return SampleZipfian(zd.Generator, zd.Alpha, zd.Skew, zd.Zeta, zd.ZetaTwoSkew)
}

func (zd *ZipfianDistribution) Copy() Distribution {
	return &ZipfianDistribution{
		Generator:   zd.Generator,
		Alpha:       zd.Alpha,
		Skew:        zd.Skew,
		Zeta:        zd.Zeta,
		ZetaTwoSkew: zd.ZetaTwoSkew,
	}
}

func (zd *ZipfianDistribution) StringSerialize() string {
	return fmt.Sprintf("Zipfian:%d:%f:%f", int64(zd.Alpha), zd.Skew, zd.Zeta)
}

func (zd *ZipfianDistribution) StringDeserialize(data string) (Distribution, error) {
	var alpha int64
	var skew, zeta float64
	_, err := fmt.Sscanf(data, "Zipfian:%d:%f:%f", &alpha, &skew, &zeta)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize ZipfianDistribution: %w", err)
	}
	if !zd.SetParameters(float64(alpha), skew, -1.0) {
		return nil, fmt.Errorf("invalid parameters in serialized data: alpha=%d, skew=%f", alpha, skew)
	}
	zd.Zeta = zeta
	return zd, nil
}

// Harmonic computes the Nth generalized harmonic number with the given skew parameter.
func Harmonic(limit int64, skew float64) float64 {
	result := 1.0
	for i := int64(2); i <= limit; i++ {
		result += math.Pow(1.0/float64(i), skew)
	}
	return result
}

// SampleZipfian generates a Zipfian-distributed variate using an approximation algorithm.
func SampleZipfian(gen rand.Source, alpha, skew, zeta, zetaTwoSkew float64) float64 {
	r := rand.New(gen)
	over := 1.0 / (1.0 - skew)
	eta := (1.0 - math.Pow(2.0/alpha, 1.0-skew)) / (1.0 - zetaTwoSkew/zeta)

	// Equivalent to nextExclusiveDouble() in open interval (0, 1)
	uExclusive := (float64(r.Uint64()>>11)*1.1102230246251565e-16 + 5.551115123125782e-17)

	uz := uExclusive * zeta
	if uz < 1.0 {
		return 1.0
	}
	if uz < zetaTwoSkew {
		return 2.0
	}
	return 1.0 + (alpha * math.Pow(eta*uExclusive-eta+1.0, over))
}
