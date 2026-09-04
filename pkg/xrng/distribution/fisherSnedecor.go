package distribution

import (
	"fmt"
	"math"
	"math/rand/v2"
)

// FisherSnedecorDistribution represents a two-parameter F-distribution with range (0, +Inf).
type FisherSnedecorDistribution struct {
	Generator rand.Source
	Alpha     float64 // d1 degrees of freedom
	Beta      float64 // d2 degrees of freedom
}

// Ensure FisherSnedecorDistribution implements the Distribution interface.
var _ Distribution = (*FisherSnedecorDistribution)(nil)

// NewFisherSnedecorDistribution creates a FisherSnedecorDistribution with the given generator, alpha, and beta.
func NewFisherSnedecorDistribution(gen rand.Source, alpha, beta float64) (*FisherSnedecorDistribution, error) {
	fd := &FisherSnedecorDistribution{
		Generator: gen,
	}
	if !fd.SetParameters(alpha, beta, 0.0) {
		return nil, fmt.Errorf("given alpha (%f) and/or beta (%f) are invalid", alpha, beta)
	}
	return fd, nil
}

func (fd *FisherSnedecorDistribution) GetTag() string {
	return "FisherSnedecor"
}

func (fd *FisherSnedecorDistribution) GetParameterA() float64 {
	return fd.Alpha
}

func (fd *FisherSnedecorDistribution) GetParameterB() float64 {
	return fd.Beta
}

func (fd *FisherSnedecorDistribution) GetParameterC() float64 {
	return math.NaN()
}

func (fd *FisherSnedecorDistribution) GetMaximum() float64 {
	return math.Inf(1)
}

func (fd *FisherSnedecorDistribution) GetMinimum() float64 {
	return 0.0
}

func (fd *FisherSnedecorDistribution) GetMean() float64 {
	if fd.Beta > 2.0 {
		return fd.Beta / (fd.Beta - 2.0)
	}
	panic("Mean cannot be determined for the given parameters")
}

func (fd *FisherSnedecorDistribution) GetMedian() float64 {
	panic("Median is undefined for FisherSnedecor distribution")
}

func (fd *FisherSnedecorDistribution) GetMode() []float64 {
	if fd.Alpha > 2.0 {
		mode := ((fd.Alpha - 2.0) / fd.Alpha) * (fd.Beta / (fd.Beta + 2.0))
		return []float64{mode}
	}
	panic("Mode cannot be determined for the given parameters")
}

func (fd *FisherSnedecorDistribution) GetVariance() float64 {
	if fd.Beta > 4.0 {
		betaSq := fd.Beta * fd.Beta
		betaMinusTwoSq := (fd.Beta - 2.0) * (fd.Beta - 2.0)
		return (2.0 * betaSq * (fd.Alpha + fd.Beta - 2.0)) / (fd.Alpha * betaMinusTwoSq * (fd.Beta - 4.0))
	}
	panic("Variance cannot be determined for the given parameters")
}

func (fd *FisherSnedecorDistribution) SetParameters(a, b, c float64) bool {
	if a > 0.0 && b > 0.0 {
		fd.Alpha = a
		fd.Beta = b
		return true
	}
	return false
}

func (fd *FisherSnedecorDistribution) NextDouble() float64 {
	return SampleFisherSnedecor(fd.Generator, fd.Alpha, fd.Beta)
}

func (fd *FisherSnedecorDistribution) Copy() Distribution {
	return &FisherSnedecorDistribution{
		Generator: fd.Generator,
		Alpha:     fd.Alpha,
		Beta:      fd.Beta,
	}
}

func (fd *FisherSnedecorDistribution) StringSerialize() string {
	return fmt.Sprintf("FisherSnedecor:%f:%f", fd.Alpha, fd.Beta)
}

func (fd *FisherSnedecorDistribution) StringDeserialize(data string) (Distribution, error) {
	var alpha, beta float64
	_, err := fmt.Sscanf(data, "FisherSnedecor:%f:%f", &alpha, &beta)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize FisherSnedecorDistribution: %w", err)
	}
	if !fd.SetParameters(alpha, beta, 0) {
		return nil, fmt.Errorf("invalid parameters in serialized data: alpha=%f, beta=%f", alpha, beta)
	}
	return fd, nil
}

// SampleFisherSnedecor generates a Fisher-Snedecor F-distributed variate using a Beta distribution sample.
func SampleFisherSnedecor(gen rand.Source, alpha, beta float64) float64 {
	x := SampleBeta(gen, alpha*0.5, beta*0.5)
	return (beta * x) / (alpha * (1.0 - x))
}
