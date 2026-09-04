package distribution

import (
	"fmt"
	"math"
	"math/rand/v2"
)

// ErlangDistribution represents a two-parameter Erlang distribution with range [0, +Inf).
type ErlangDistribution struct {
	Generator rand.Source
	Alpha     int
	Lambda    float64
}

// Ensure ErlangDistribution implements the Distribution interface.
var _ Distribution = (*ErlangDistribution)(nil)

// NewErlangDistribution creates an ErlangDistribution with the given generator, alpha (shape parameter, k), and lambda (rate parameter).
func NewErlangDistribution(gen rand.Source, alpha int, lambda float64) (*ErlangDistribution, error) {
	ed := &ErlangDistribution{
		Generator: gen,
	}
	if !ed.SetParameters(float64(alpha), lambda, 0.0) {
		return nil, fmt.Errorf("given alpha (%d) and/or lambda (%f) are invalid", alpha, lambda)
	}
	return ed, nil
}

func (ed *ErlangDistribution) GetTag() string {
	return "Erlang"
}

func (ed *ErlangDistribution) GetParameterA() float64 {
	return float64(ed.Alpha)
}

func (ed *ErlangDistribution) GetParameterB() float64 {
	return ed.Lambda
}

func (ed *ErlangDistribution) GetParameterC() float64 {
	return math.NaN()
}

func (ed *ErlangDistribution) GetMaximum() float64 {
	return math.Inf(1)
}

func (ed *ErlangDistribution) GetMinimum() float64 {
	return 0.0
}

func (ed *ErlangDistribution) GetMean() float64 {
	return float64(ed.Alpha) / ed.Lambda
}

func (ed *ErlangDistribution) GetMedian() float64 {
	panic("Median is undefined for Erlang distribution")
}

func (ed *ErlangDistribution) GetMode() []float64 {
	return []float64{(float64(ed.Alpha) - 1.0) / ed.Lambda}
}

func (ed *ErlangDistribution) GetVariance() float64 {
	return float64(ed.Alpha) / (ed.Lambda * ed.Lambda)
}

func (ed *ErlangDistribution) SetParameters(a, b, c float64) bool {
	if a >= 1.0 && b > 0.0 {
		ed.Alpha = int(a)
		ed.Lambda = b
		return true
	}
	return false
}

func (ed *ErlangDistribution) NextDouble() float64 {
	return SampleErlang(ed.Generator, ed.Alpha, ed.Lambda)
}

func (ed *ErlangDistribution) Copy() Distribution {
	return &ErlangDistribution{
		Generator: ed.Generator,
		Alpha:     ed.Alpha,
		Lambda:    ed.Lambda,
	}
}

func (ed *ErlangDistribution) StringSerialize() string {
	return fmt.Sprintf("Erlang:%d:%f", ed.Alpha, ed.Lambda)
}

func (ed *ErlangDistribution) StringDeserialize(data string) (Distribution, error) {
	var alpha int
	var lambda float64
	_, err := fmt.Sscanf(data, "Erlang:%d:%f", &alpha, &lambda)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize ErlangDistribution: %w", err)
	}
	if !ed.SetParameters(float64(alpha), lambda, 0) {
		return nil, fmt.Errorf("invalid parameters in serialized data: alpha=%d, lambda=%f", alpha, lambda)
	}
	return ed, nil
}

// SampleErlang generates an Erlang-distributed variate using Marsaglia and Tsang's rejection method.
func SampleErlang(gen rand.Source, alpha int, lambda float64) float64 {
	if math.IsInf(lambda, 1) {
		return float64(alpha)
	}

	d := float64(alpha) - (1.0 / 3.0)
	c := 1.0 / math.Sqrt(9.0*d)

	r := rand.New(gen)

	for {
		x := r.NormFloat64()
		v := 1.0 + (c * x)
		for v <= 0.0 {
			x = r.NormFloat64()
			v = 1.0 + (c * x)
		}

		v = v * v * v
		// Generate standard uniform sample u in open interval (0, 1)
		u := (float64(r.Uint64()>>11)*1.1102230246251565e-16 + 5.551115123125782e-17)
		xSq := x * x

		if u < 1.0-(0.0331*xSq*xSq) {
			return d * v / lambda
		}

		if math.Log(u) < (0.5*xSq)+(d*(1.0-v+math.Log(v))) {
			return d * v / lambda
		}
	}
}
