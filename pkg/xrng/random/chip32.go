package random

import (
	"math/bits"
	"math/rand"
)

// Chip32Random implements a 32-bit pseudo-random number generator translated from Java's Chip32Random.
// It implements math/rand.Source, math/rand.Source64, and math/rand/v2.Source.
type Chip32Random struct {
	StateA int32
	StateB int32
	StateC int32
	StateD int32
}

// Compile-time assertions for interface compliance
var (
	_ rand.Source   = (*Chip32Random)(nil)
	_ rand.Source64 = (*Chip32Random)(nil)
)

// NewChip32Random creates a generator using a 64-bit seed.
func NewChip32Random(seed int64) *Chip32Random {
	c := &Chip32Random{}
	c.Seed(seed)
	return c
}

// NewChip32RandomWithIntSeed creates a generator using a 32-bit seed.
func NewChip32RandomWithIntSeed(seed int32) *Chip32Random {
	c := &Chip32Random{}
	c.SeedInt(seed)
	return c
}

// NewChip32RandomWithStates creates a generator using four explicit state values.
func NewChip32RandomWithStates(a, b, c, d int32) *Chip32Random {
	return &Chip32Random{
		StateA: a,
		StateB: b,
		StateC: c,
		StateD: d,
	}
}

// Seed initializes all 4 states using a 64-bit seed.
func (c *Chip32Random) Seed(seed int64) {
	uSeed := uint64(seed)

	a := int32(uSeed) ^ -0x24B0F46F          // 0xDB4F0B91
	b := int32(uSeed>>11) ^ -0x441FA9CD      // 0xBBE05633
	cState := int32(uSeed>>21) ^ -0x5F0D138B // 0xA0F2EC75
	d := int32(uSeed>>32) ^ -0x761E7D7B      // 0x89E18285

	a = (a ^ int32(uint32(a)>>16)) * 0x21f0aaad
	a = (a ^ int32(uint32(a)>>15)) * 0x735a2d97
	c.StateA = a ^ int32(uint32(a)>>15)

	b = (b ^ int32(uint32(b)>>16)) * 0x21f0aaad
	b = (b ^ int32(uint32(b)>>15)) * 0x735a2d97
	c.StateB = b ^ int32(uint32(b)>>15)

	cState = (cState ^ int32(uint32(cState)>>16)) * 0x21f0aaad
	cState = (cState ^ int32(uint32(cState)>>15)) * 0x735a2d97
	c.StateC = cState ^ int32(uint32(cState)>>15)

	d = (d ^ int32(uint32(d)>>16)) * 0x21f0aaad
	d = (d ^ int32(uint32(d)>>15)) * 0x735a2d97
	c.StateD = d ^ int32(uint32(d)>>15)
}

// SeedInt initializes all 4 states using a 32-bit seed.
func (c *Chip32Random) SeedInt(seed int32) {
	uSeed := uint32(seed)

	a := seed ^ -0x24B0F46F                                     // 0xDB4F0B91
	b := int32(bits.RotateLeft32(uSeed, 8)) ^ -0x441FA9CD       // 0xBBE05633
	cState := int32(bits.RotateLeft32(uSeed, 16)) ^ -0x5F0D138B // 0xA0F2EC75
	d := int32(bits.RotateLeft32(uSeed, 24)) ^ -0x761E7D7B      // 0x89E18285

	a = (a ^ int32(uint32(a)>>16)) * 0x21f0aaad
	a = (a ^ int32(uint32(a)>>15)) * 0x735a2d97
	c.StateA = a ^ int32(uint32(a)>>15)

	b = (b ^ int32(uint32(b)>>16)) * 0x21f0aaad
	b = (b ^ int32(uint32(b)>>15)) * 0x735a2d97
	c.StateB = b ^ int32(uint32(b)>>15)

	cState = (cState ^ int32(uint32(cState)>>16)) * 0x21f0aaad
	cState = (cState ^ int32(uint32(cState)>>15)) * 0x735a2d97
	c.StateC = cState ^ int32(uint32(cState)>>15)

	d = (d ^ int32(uint32(d)>>16)) * 0x21f0aaad
	d = (d ^ int32(uint32(d)>>15)) * 0x735a2d97
	c.StateD = d ^ int32(uint32(d)>>15)
}

// NextInt generates the next pseudo-random 32-bit signed integer.
func (c *Chip32Random) NextInt() int32 {
	fa := c.StateA
	fb := c.StateB
	fc := c.StateC
	fd := c.StateD

	c.StateA = fb + fc
	c.StateB = fd ^ fa
	c.StateC = int32(bits.RotateLeft32(uint32(fb), 11))
	c.StateD = fd + -0x61C88647 // 0x9E3779B9

	rotA := int32(bits.RotateLeft32(uint32(fa), 14))
	rotB := int32(bits.RotateLeft32(uint32(fb), 23))
	return rotA ^ (rotB + fc)
}

// Uint64 returns a 64-bit pseudo-random value required by math/rand.Source64.
func (c *Chip32Random) Uint64() uint64 {
	fa := c.StateA
	fb := c.StateB
	fc := c.StateC
	fd := c.StateD

	hiRotA := int32(bits.RotateLeft32(uint32(fa), 14))
	hiRotB := int32(bits.RotateLeft32(uint32(fb), 23))
	hi := hiRotA ^ (hiRotB + fc)

	ga := fb + fc
	gb := fa ^ fd
	gc := int32(bits.RotateLeft32(uint32(fb), 11))
	gd := fd + -0x61C88647

	loRotA := int32(bits.RotateLeft32(uint32(ga), 14))
	loRotB := int32(bits.RotateLeft32(uint32(gb), 23))
	lo := loRotA ^ (loRotB + gc)

	c.StateA = gb + gc
	c.StateB = ga ^ gd
	c.StateC = int32(bits.RotateLeft32(uint32(gb), 11))
	c.StateD = gd + -0x61C88647

	return (uint64(uint32(hi)) << 32) ^ uint64(uint32(lo))
}

// Int63 returns a non-negative 63-bit integer required by math/rand.Source.
func (c *Chip32Random) Int63() int64 {
	return int64(c.Uint64() & 0x7FFFFFFFFFFFFFFF)
}

// PreviousInt steps back one state and returns the prior 32-bit output.
func (c *Chip32Random) PreviousInt() int32 {
	ga := c.StateA
	gb := c.StateB
	gc := c.StateC
	gd := c.StateD

	c.StateD = gd - -0x61C88647
	c.StateA = gb ^ c.StateD
	c.StateB = int32(bits.RotateLeft32(uint32(gc), 21)) // right-rotate 11 is left-rotate 21
	c.StateC = ga - c.StateB

	rotA := int32(bits.RotateLeft32(uint32(c.StateA), 14))
	rotB := int32(bits.RotateLeft32(uint32(c.StateB), 23))
	return rotA ^ (rotB + c.StateC)
}

// PreviousLong steps back two states and returns the prior 64-bit output.
func (c *Chip32Random) PreviousLong() int64 {
	ga := c.StateA
	gb := c.StateB
	gc := c.StateC
	gd := c.StateD

	fd := gd - -0x61C88647
	fa := gb ^ fd
	fb := int32(bits.RotateLeft32(uint32(gc), 21))
	fc := ga - fb

	loRotA := int32(bits.RotateLeft32(uint32(fa), 14))
	loRotB := int32(bits.RotateLeft32(uint32(fb), 23))
	lo := loRotA ^ (loRotB + fc)

	c.StateD = fd - -0x61C88647
	c.StateA = fb ^ c.StateD
	c.StateB = int32(bits.RotateLeft32(uint32(fc), 21))
	c.StateC = fa - c.StateB

	hiRotA := int32(bits.RotateLeft32(uint32(c.StateA), 14))
	hiRotB := int32(bits.RotateLeft32(uint32(c.StateB), 23))
	hi := hiRotA ^ (hiRotB + c.StateC)

	return int64((uint64(uint32(hi)) << 32) ^ uint64(uint32(lo)))
}

// Copy creates an exact duplicate of this generator.
func (c *Chip32Random) Copy() *Chip32Random {
	return &Chip32Random{
		StateA: c.StateA,
		StateB: c.StateB,
		StateC: c.StateC,
		StateD: c.StateD,
	}
}
