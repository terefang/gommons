package distribution

// Distribution defines the interface for statistical random distributions.
type Distribution interface {
	// NextDouble generates a random double using this distribution.
	NextDouble() float64

	// GetMaximum returns the maximum possible value of distributed random numbers.
	GetMaximum() float64

	// GetMean returns the mean of distributed random numbers.
	GetMean() float64

	// GetMedian returns the median of distributed random numbers.
	GetMedian() float64

	// GetMinimum returns the minimum possible value of distributed random numbers.
	GetMinimum() float64

	// GetMode returns the mode(s) of distributed random numbers.
	GetMode() []float64

	// GetVariance returns the variance of distributed random numbers.
	GetVariance() float64

	// SetParameters validates and sets up to 3 parameters for this distribution.
	// Returns true if parameters are valid and set; false otherwise.
	SetParameters(a, b, c float64) bool

	// GetParameterA returns parameter "A" as a float64 (or NaN if unused).
	GetParameterA() float64

	// GetParameterB returns parameter "B" as a float64 (or NaN if unused).
	GetParameterB() float64

	// GetParameterC returns parameter "C" as a float64 (or NaN if unused).
	GetParameterC() float64

	// GetTag returns a unique identifier string for this type of distribution.
	GetTag() string

	// Copy returns a deep copy of this distribution.
	Copy() Distribution

	// StringSerialize serializes the current state of this distribution to a string.
	StringSerialize() string

	// StringDeserialize restores the state of this distribution from a serialized string.
	StringDeserialize(data string) (Distribution, error)
}
