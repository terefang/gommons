package random

import (
	"math/bits"
)

// Bear32Random implements a pseudo-random number generator equivalent to Java's Bear32Random.
// It implements math/rand.Source and math/rand.Source64 (Go 1.8+), as well as math/rand/v2.Source.
type Bear32Random struct {
	StateA int32
	StateB int32
	StateC int32
	StateD int32
}

// NewBear32Random creates a generator seeded with state values derived from a 64-bit seed.
func NewBear32Random(seed int64) *Bear32Random {
	b := &Bear32Random{}
	b.Seed(seed)
	return b
}

// NewBear32RandomWithStates creates a generator using explicit state values.
func NewBear32RandomWithStates(a, b, c, d int32) *Bear32Random {
	return &Bear32Random{
		StateA: a,
		StateB: b,
		StateC: c,
		StateD: d,
	}
}

// Seed initializes all 4 states using a single 64-bit seed.
func (b *Bear32Random) Seed(seed int64) {
	uSeed := uint64(seed)

	a := int32(uSeed) ^ -0x24B0F46F          // 0xDB4F0B91
	bState := int32(uSeed>>16) ^ -0x441FA9CD // 0xBBE05633
	c := int32(uSeed>>32) ^ -0x5F0D138B      // 0xA0F2EC75
	d := int32(uSeed>>48) ^ -0x761E7D7B      // 0x89E18285

	a = (a ^ int32(uint32(a)>>16)) * 0x21f0aaad
	a = (a ^ int32(uint32(a)>>15)) * 0x735a2d97
	b.StateA = a ^ int32(uint32(a)>>15)

	bState = (bState ^ int32(uint32(bState)>>16)) * 0x21f0aaad
	bState = (bState ^ int32(uint32(bState)>>15)) * 0x735a2d97
	b.StateB = bState ^ int32(uint32(bState)>>15)

	c = (c ^ int32(uint32(c)>>16)) * 0x21f0aaad
	c = (c ^ int32(uint32(c)>>15)) * 0x735a2d97
	b.StateC = c ^ int32(uint32(c)>>15)

	d = (d ^ int32(uint32(d)>>16)) * 0x21f0aaad
	d = (d ^ int32(uint32(d)>>15)) * 0x735a2d97
	b.StateD = d ^ int32(uint32(d)>>15)
}

// NextInt generates the next pseudo-random 32-bit signed integer.
func (b *Bear32Random) NextInt() int32 {
	b.StateA += -0x61C88647 // 0x9E3779B9
	a := b.StateA

	b.StateB += a ^ int32(bits.LeadingZeros32(uint32(a)))
	sb := b.StateB

	a &= sb
	b.StateC += sb ^ int32(bits.LeadingZeros32(uint32(a)))
	sc := b.StateC

	a &= sc
	b.StateD += sc ^ int32(bits.LeadingZeros32(uint32(a)))
	sd := b.StateD

	a = sd + int32(bits.RotateLeft32(uint32(a), 13))
	a *= 0x2C1B3C6D
	a = (a ^ int32(uint32(a)>>12)) * 0x297A2D39
	a ^= int32(uint32(a) >> 15)

	return a
}

// Int63 returns a non-negative pseudo-random 63-bit integer required by math/rand.Source.
func (b *Bear32Random) Int63() int64 {
	return int64(b.Uint64() & 0x7FFFFFFFFFFFFFFF)
}

// Uint64 returns a 64-bit pseudo-random value required by math/rand.Source64 and math/rand/v2.Source.
func (b *Bear32Random) Uint64() uint64 {
	hi := uint64(uint32(b.NextInt()))
	lo := uint64(uint32(b.NextInt()))
	return (hi << 32) | lo
}

// PreviousInt steps back one state and returns the prior 32-bit output.
func (b *Bear32Random) PreviousInt() int32 {
	a := b.StateA
	sb := b.StateB
	sc := b.StateC
	sd := b.StateD

	m := a & sb & sc
	m = sd + int32(bits.RotateLeft32(uint32(m), 13))
	m *= 0x2C1B3C6D
	m = (m ^ int32(uint32(m)>>12)) * 0x297A2D39
	m ^= int32(uint32(m) >> 15)

	b.StateA = a - -0x61C88647
	b.StateB = sb - (a ^ int32(bits.LeadingZeros32(uint32(a))))

	a &= sb
	b.StateC = sc - (sb ^ int32(bits.LeadingZeros32(uint32(a))))

	a &= sc
	b.StateD = sd - (sc ^ int32(bits.LeadingZeros32(uint32(a))))

	return m
}

// PreviousLong steps back two states and returns the prior 64-bit value.
func (b *Bear32Random) PreviousLong() int64 {
	lo := uint64(uint32(b.PreviousInt()))
	hi := uint64(uint32(b.PreviousInt()))
	return int64((hi << 32) | lo)
}

// Copy creates a distinct duplicate of this generator's state.
func (b *Bear32Random) Copy() *Bear32Random {
	return &Bear32Random{
		StateA: b.StateA,
		StateB: b.StateB,
		StateC: b.StateC,
		StateD: b.StateD,
	}
}
