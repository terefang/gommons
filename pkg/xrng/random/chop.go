package random

import (
	"math/bits"
	"math/rand"
)

// ChopRandom implements a 32-bit pseudo-random number generator translated from Java's ChopRandom.
// It implements math/rand.Source, math/rand.Source64, and math/rand/v2.Source.
type ChopRandom struct {
	StateA int32
	StateB int32
	StateC int32
	StateD int32
}

// Compile-time checks for interface compliance
var (
	_ rand.Source   = (*ChopRandom)(nil)
	_ rand.Source64 = (*ChopRandom)(nil)
)

// NewChopRandom creates a generator seeded using a 64-bit value.
func NewChopRandom(seed int64) *ChopRandom {
	c := &ChopRandom{}
	c.Seed(seed)
	return c
}

// NewChopRandomWithStates creates a generator using four explicit state values.
func NewChopRandomWithStates(a, b, c, d int32) *ChopRandom {
	return &ChopRandom{
		StateA: a,
		StateB: b,
		StateC: c,
		StateD: d,
	}
}

// Seed initializes all 4 states using SplitMix64-style mixing on the 64-bit seed.
func (c *ChopRandom) Seed(seed int64) {
	uSeed := uint64(seed)

	mix := func() int32 {
		uSeed += 0x9E3779B97F4A7C15
		x := uSeed
		x = (x ^ (x >> 27)) * 0x3C79AC492BA7B653
		x = (x ^ (x >> 33)) * 0x1C69B3F74AC4AE35
		return int32(x ^ (x >> 27))
	}

	c.StateA = mix()
	c.StateB = mix()
	c.StateC = mix()
	c.StateD = mix()
}

// NextInt generates the next pseudo-random 32-bit signed integer.
func (c *ChopRandom) NextInt() int32 {
	fa := c.StateA
	fb := c.StateB
	fc := c.StateC
	fd := c.StateD

	sa := fb ^ fc
	c.StateA = int32(bits.RotateLeft32(uint32(sa), 26))

	sb := fc ^ fd
	c.StateB = int32(bits.RotateLeft32(uint32(sb), 11))

	c.StateC = fa ^ (fb + fc)
	c.StateD = fd + -0x524A4E9B // 0xADB5B165

	return fc
}

// Uint64 returns a 64-bit pseudo-random value required by math/rand.Source64 and math/rand/v2.Source.
func (c *ChopRandom) Uint64() uint64 {
	fa := c.StateA
	fb := c.StateB
	fc := c.StateC
	fd := c.StateD

	ga := fb ^ fc
	ga = int32(bits.RotateLeft32(uint32(ga), 26))

	gb := fc ^ fd
	gb = int32(bits.RotateLeft32(uint32(gb), 11))

	gc := fa ^ (fb + fc)
	gd := fd + -0x524A4E9B // 0xADB5B165

	fa = gb ^ gc
	c.StateA = int32(bits.RotateLeft32(uint32(fa), 26))

	fb = gc ^ gd
	c.StateB = int32(bits.RotateLeft32(uint32(fb), 11))

	c.StateC = ga ^ (gb + gc)
	c.StateD = fd + 0x5B6B62CA

	hi := uint64(uint32(fc))
	lo := uint64(uint32(gc))
	return (hi << 32) ^ lo
}

// Int63 returns a non-negative 63-bit integer required by math/rand.Source.
func (c *ChopRandom) Int63() int64 {
	return int64(c.Uint64() & 0x7FFFFFFFFFFFFFFF)
}

// PreviousInt steps back one state transition and returns the prior 32-bit output.
func (c *ChopRandom) PreviousInt() int32 {
	ga := c.StateA
	gb := c.StateB
	gc := c.StateC

	c.StateD = c.StateD - -0x524A4E9B // -0xADB5B165
	c.StateC = int32(bits.RotateLeft32(uint32(gb), 21)) ^ c.StateD
	c.StateB = int32(bits.RotateLeft32(uint32(ga), 6)) ^ c.StateC
	c.StateA = gc ^ (c.StateB + c.StateC)

	return c.StateC
}

// PreviousLong steps back two state transitions and returns the prior 64-bit output.
func (c *ChopRandom) PreviousLong() int64 {
	fa := c.StateA
	fb := c.StateB
	fc := c.StateC

	gc := int32(bits.RotateLeft32(uint32(fb), 21)) ^ (c.StateD - -0x524A4E9B)
	gb := int32(bits.RotateLeft32(uint32(fa), 6)) ^ gc
	ga := fc ^ (gb + gc)

	c.StateD = c.StateD - 0x5B6B62CA
	c.StateC = int32(bits.RotateLeft32(uint32(gb), 21)) ^ c.StateD
	c.StateB = int32(bits.RotateLeft32(uint32(ga), 6)) ^ c.StateC
	c.StateA = gc ^ (c.StateB + c.StateC)

	hi := uint64(uint32(c.StateC))
	lo := uint64(uint32(gc))
	return int64((hi << 32) ^ lo)
}

// Copy creates an exact duplicate of this generator's current state.
func (c *ChopRandom) Copy() *ChopRandom {
	return &ChopRandom{
		StateA: c.StateA,
		StateB: c.StateB,
		StateC: c.StateC,
		StateD: c.StateD,
	}
}
