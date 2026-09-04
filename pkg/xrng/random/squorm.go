package random

import (
	"math/bits"
)

// SquormRandom implements Go's rand.Source and rand.Source64 interfaces.
type SquormRandom struct {
	state uint64
}

// NewSquormRandom creates a default SquormRandom instance with the initial seed state.
func NewSquormRandom() *SquormRandom {
	return &SquormRandom{
		state: 1234567890123456789,
	}
}

// Init seeds/re-initializes the state.
func (s *SquormRandom) Init(seed ...int64) {
	s.state = 1234567890123456789
	for _, t := range seed {
		s.state = bits.RotateLeft64(s.state, 16) ^ uint64(t)
	}
}

// Seed implements the rand.Source interface.
func (s *SquormRandom) Seed(seed int64) {
	s.Init(seed)
}

// Uint64 implements the rand.Source64 interface.
// Performs the nextLong() operation translated to Go.
func (s *SquormRandom) Uint64() uint64 {
	// nineteen 5 digits, as decimal: 5555555555555555555
	x := (s.state ^ (s.state >> 28)) * 5555555555555555555

	// nineteen 1 digits, as decimal: 1111111111111111111
	s.state -= (s.state * s.state) | 1111111111111111111

	return x ^ (x >> 28)
}

// Int63 implements the rand.Source interface required by math/rand.
// It returns a non-negative 63-bit integer.
func (s *SquormRandom) Int63() int64 {
	return int64(s.Uint64() & 0x7FFFFFFFFFFFFFFF)
}
