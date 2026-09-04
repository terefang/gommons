package distribution

import (
	"fmt"
	"math"
	"math/rand/v2"
)

// TriangularDistribution represents a three-parameter continuous distribution with range [Alpha, Beta].
type TriangularDistribution struct {
	Generator rand.Source
	Alpha     float64 // Lower limit (minimum)
	Beta      float64 // Upper limit (maximum)
	Gamma     float64 // Mode
}

// Ensure TriangularDistribution implements the Distribution interface.
var _ Distribution = (*TriangularDistribution)(nil)

// NewTriangularDistribution creates a TriangularDistribution with the given generator, alpha, beta, and gamma.
func NewTriangularDistribution(gen rand.Source, alpha, beta, gamma float64) (*TriangularDistribution, error) {
	td := &TriangularDistribution{
		Generator: gen,
	}
	if !td.SetParameters(alpha, beta, gamma) {
		return nil, fmt.Errorf("given alpha (%f), beta (%f), and/or gamma (%f) are invalid", alpha, beta, gamma)
	}
	return td, nil
}

func (td *TriangularDistribution) GetTag() string {
	return "Triangular"
}

func (td *TriangularDistribution) GetParameterA() float64 {
	return td.Alpha
}

func (td *TriangularDistribution) GetParameterB() float64 {
	return td.Beta
}

func (td *TriangularDistribution) GetParameterC() float64 {
	return td.Gamma
}

func (td *TriangularDistribution) GetMaximum() float64 {
	return td.Beta
}

func (td *TriangularDistribution) GetMinimum() float64 {
	return td.Alpha
}

func (td *TriangularDistribution) GetMean() float64 {
	return (td.Alpha + td.Beta + td.Gamma) / 3.0
}

func (td *TriangularDistribution) GetMedian() float64 {
	if td.Gamma >= (td.Beta-td.Alpha)*0.5 {
		return td.Alpha + (math.Sqrt((td.Beta-td.Alpha)*(td.Gamma-td.Alpha)) / math.Sqrt(2.0))
	}
	return td.Beta - (math.Sqrt((td.Beta-td.Alpha)*(td.Beta-td.Gamma)) / math.Sqrt(2.0))
}

func (td *TriangularDistribution) GetMode() []float64 {
	return []float64{td.Gamma}
}

func (td *TriangularDistribution) GetVariance() float64 {
	a, b, c := td.Alpha, td.Beta, td.Gamma
	return (a*a + b*b + c*c - a*b - a*c - b*c) / 18.0
}

func (td *TriangularDistribution) SetParameters(a, b, c float64) bool {
	if a < b && a <= c && c <= b {
		td.Alpha = a
		td.Beta = b
		td.Gamma = c
		return true
	}
	return false
}

func (td *TriangularDistribution) NextDouble() float64 {
	return SampleTriangular(td.Generator, td.Alpha, td.Beta, td.Gamma)
}

func (td *TriangularDistribution) Copy() Distribution {
	return &TriangularDistribution{
		Generator: td.Generator,
		Alpha:     td.Alpha,
		Beta:      td.Beta,
		Gamma:     td.Gamma,
	}
}

func (td *TriangularDistribution) StringSerialize() string {
	return fmt.Sprintf("Triangular:%f:%f:%f", td.Alpha, td.Beta, td.Gamma)
}

func (td *TriangularDistribution) StringDeserialize(data string) (Distribution, error) {
	var alpha, beta, gamma float64
	_, err := fmt.Sscanf(data, "Triangular:%f:%f:%f", &alpha, &beta, &gamma)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize TriangularDistribution: %w", err)
	}
	if !td.SetParameters(alpha, beta, gamma) {
		return nil, fmt.Errorf("invalid parameters in serialized data: alpha=%f, beta=%f, gamma=%f", alpha, beta, gamma)
	}
	return td, nil
}

// SampleTriangular generates a Triangular-distributed variate using inverse transform sampling.
func SampleTriangular(gen rand.Source, alpha, beta, gamma float64) float64 {
	r := rand.New(gen)
	helper1 := gamma - alpha
	helper2 := beta - alpha
	helper3 := math.Sqrt(helper1 * helper2)
	helper4 := math.Sqrt(beta - gamma)
	genNum := r.Float64() // Uniform sample in [0, 1)

	if genNum <= helper1/helper2 {
		return alpha + math.Sqrt(genNum)*helper3
	}
	return beta - math.Sqrt(genNum*helper2-helper1)*helper4
}
