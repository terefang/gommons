package random

import (
	"math/bits"
	"math/rand"
)

// Respite32Random implements a 32-bit pseudo-random number generator translated from Java's Respite32Random.
// It implements math/rand.Source, math/rand.Source64 (Go 1.8+), and math/rand/v2.Source.
type Respite32Random struct {
	StateA int32
	StateB int32
	StateC int32
}

// Compile-time checks for interface compliance.
var (
	_ rand.Source   = (*Respite32Random)(nil)
	_ rand.Source64 = (*Respite32Random)(nil)
)

// NewRespite32Random creates a generator initialized using a 64-bit seed.
func NewRespite32Random(seed int64) *Respite32Random {
	r := &Respite32Random{}
	r.Seed(seed)
	return r
}

// NewRespite32RandomWithStates creates a generator using three explicit state values.
func NewRespite32RandomWithStates(stateA, stateB, stateC int32) *Respite32Random {
	return &Respite32Random{
		StateA: stateA,
		StateB: stateB,
		StateC: stateC,
	}
}

// Seed initializes all 3 states using Speck cipher operations based on the given 64-bit seed.
func (r *Respite32Random) Seed(seed int64) {
	a := int32(seed)
	b := int32(uint64(seed) >> 32)
	c := int32(^uint64(seed) >> 16)

	for i := 0; i < 5; i++ {
		c++
		b = int32(bits.RotateLeft32(uint32(b), 24)) + a ^ c
		a = int32(bits.RotateLeft32(uint32(a), 3)) ^ b
	}
	r.StateA = a

	for i := 0; i < 5; i++ {
		c++
		b = int32(bits.RotateLeft32(uint32(b), 24)) + a ^ c
		a = int32(bits.RotateLeft32(uint32(a), 3)) ^ b
	}
	r.StateB = a

	for i := 0; i < 5; i++ {
		c++
		b = int32(bits.RotateLeft32(uint32(b), 24)) + a ^ c
		a = int32(bits.RotateLeft32(uint32(a), 3)) ^ b
	}
	r.StateC = a
}

// NextInt advances the internal state and returns the next pseudo-random 32-bit signed integer.
func (r *Respite32Random) NextInt() int32 {
	r.StateA += -0x6E1EF25B // 0x91E10DA5 as signed int32
	r.StateB += 0x6C8E9CF5 ^ int32(bits.LeadingZeros32(uint32(r.StateA)))
	r.StateC += 0x7FEB352D ^ int32(bits.LeadingZeros32(uint32(r.StateA&r.StateB)))

	a := r.StateA
	b := r.StateB
	c := r.StateC

	b = int32(bits.RotateLeft32(uint32(b), 24)) + a ^ c
	a = int32(bits.RotateLeft32(uint32(a), 3)) ^ b

	b = int32(bits.RotateLeft32(uint32(b), 24)) + a ^ c
	a = int32(bits.RotateLeft32(uint32(a), 3)) ^ b

	b = int32(bits.RotateLeft32(uint32(b), 24)) + a ^ c
	a = int32(bits.RotateLeft32(uint32(a), 3)) ^ b

	b = int32(bits.RotateLeft32(uint32(b), 24)) + a ^ c
	a = int32(bits.RotateLeft32(uint32(a), 3)) ^ b

	return a
}

// Uint64 returns a 64-bit pseudo-random value required by math/rand.Source64 and math/rand/v2.Source.
func (r *Respite32Random) Uint64() uint64 {
	hi := uint64(uint32(r.NextInt()))
	lo := uint64(uint32(r.NextInt()))
	return (hi << 32) | lo
}

// Int63 returns a non-negative 63-bit integer required by math/rand.Source.
func (r *Respite32Random) Int63() int64 {
	return int64(r.Uint64() & 0x7FFFFFFFFFFFFFFF)
}

// PreviousInt rolls back the state by one step and returns the 32-bit output generated at that prior state.
func (r *Respite32Random) PreviousInt() int32 {
	a := r.StateA
	b := r.StateB
	c := r.StateC

	b = int32(bits.RotateLeft32(uint32(b), 24)) + a ^ c
	a = int32(bits.RotateLeft32(uint32(a), 3)) ^ b

	b = int32(bits.RotateLeft32(uint32(b), 24)) + a ^ c
	a = int32(bits.RotateLeft32(uint32(a), 3)) ^ b

	b = int32(bits.RotateLeft32(uint32(b), 24)) + a ^ c
	a = int32(bits.RotateLeft32(uint32(a), 3)) ^ b

	b = int32(bits.RotateLeft32(uint32(b), 24)) + a ^ c
	res := int32(bits.RotateLeft32(uint32(a), 3)) ^ b

	r.StateA += -0x6E1EF25B // 0x91E10DA5 as signed int32
	r.StateB += 0x6C8E9CF5 ^ int32(bits.LeadingZeros32(uint32(a)))
	r.StateC += 0x7FEB352D ^ int32(bits.LeadingZeros32(uint32(a&b)))

	return res
}

// PreviousLong rolls back two steps and returns the 64-bit value that would have been generated.
func (r *Respite32Random) PreviousLong() int64 {
	lo := uint64(uint32(r.PreviousInt()))
	hi := uint64(uint32(r.PreviousInt()))
	return int64((hi << 32) | lo)
}

// Copy creates an exact duplicate of this generator's current state.
func (r *Respite32Random) Copy() *Respite32Random {
	return &Respite32Random{
		StateA: r.StateA,
		StateB: r.StateB,
		StateC: r.StateC,
	}
}
