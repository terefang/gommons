package random

import (
	"math/bits"
	"math/rand"
)

// Xoshiro160RoadroxoRandom implements a 160-bit state pseudo-random number generator.
// It complies with math/rand.Source, math/rand.Source64, and math/rand/v2.Source.
type Xoshiro160RoadroxoRandom struct {
	StateA int32
	StateB int32
	StateC int32
	StateD int32
	StateE int32
}

// Compile-time checks for interface compliance.
var (
	_ rand.Source   = (*Xoshiro160RoadroxoRandom)(nil)
	_ rand.Source64 = (*Xoshiro160RoadroxoRandom)(nil)
)

// NewXoshiro160RoadroxoRandom constructs a generator initialized with a 64-bit seed.
func NewXoshiro160RoadroxoRandom(seed int64) *Xoshiro160RoadroxoRandom {
	x := &Xoshiro160RoadroxoRandom{}
	x.Seed(seed)
	return x
}

// NewXoshiro160RoadroxoRandomWithStates constructs a generator using explicit state values.
func NewXoshiro160RoadroxoRandomWithStates(a, b, c, d, e int32) *Xoshiro160RoadroxoRandom {
	if (a | b | c | d) == 0 {
		d = 1
	}
	return &Xoshiro160RoadroxoRandom{
		StateA: a,
		StateB: b,
		StateC: c,
		StateD: d,
		StateE: e,
	}
}

// Seed initializes all 5 states using SplitMix64 based on the given 64-bit seed.
func (x *Xoshiro160RoadroxoRandom) Seed(seed int64) {
	uSeed := uint64(seed)

	s := uSeed
	s ^= s >> 27
	s *= 0x3C79AC492BA7B653
	s ^= s >> 33
	s *= 0x1C69B3F74AC4AE35
	s ^= s >> 27
	x.StateA = int32(s)
	x.StateB = int32(s >> 32)

	s = uSeed + 0x9E3779B97F4A7C15
	s ^= s >> 27
	s *= 0x3C79AC492BA7B653
	s ^= s >> 33
	s *= 0x1C69B3F74AC4AE35
	s ^= s >> 27
	x.StateC = int32(s)
	x.StateD = int32(s >> 32)
	x.StateE = int32(uSeed ^ (uSeed >> 32))

	if (x.StateA | x.StateB | x.StateC | x.StateD) == 0 {
		x.StateD = 1
	}
}

// SeedInt32 initializes all 5 states using 32-bit math operations based on an int32 seed.
func (x *Xoshiro160RoadroxoRandom) SeedInt32(seed int32) {
	imul := func(a, b int32) int32 {
		return int32(uint32(a) * uint32(b))
	}

	a := seed ^ -0x24B0F46F                                       // 0xDB4F0B91 as signed int32
	b := int32(bits.RotateLeft32(uint32(seed), 8)) ^ -0x441FA9CD  // 0xBBE05633
	c := int32(bits.RotateLeft32(uint32(seed), 16)) ^ -0x5F0D138B // 0xA0F2EC75
	d := int32(bits.RotateLeft32(uint32(seed), 24)) ^ -0x761E7D7B // 0x89E18285

	a = imul(a^(a>>16), 0x21f0aaad)
	a = imul(a^(a>>15), 0x735a2d97)
	x.StateA = a ^ (a >> 15)

	b = imul(b^(b>>16), 0x21f0aaad)
	b = imul(b^(b>>15), 0x735a2d97)
	x.StateB = b ^ (b >> 15)

	c = imul(c^(c>>16), 0x21f0aaad)
	c = imul(c^(c>>15), 0x735a2d97)
	x.StateC = c ^ (c >> 15)

	d = imul(d^(d>>16), 0x21f0aaad)
	d = imul(d^(d>>15), 0x735a2d97)
	x.StateD = d ^ (d >> 15)

	x.StateE = seed ^ (seed >> 16)

	if (x.StateA | x.StateB | x.StateC | x.StateD) == 0 {
		x.StateD = 1
	}
}

// NextInt advances state and returns the next pseudo-random 32-bit signed integer.
func (x *Xoshiro160RoadroxoRandom) NextInt() int32 {
	res := int32(bits.RotateLeft32(uint32(x.StateE), 23)) ^ (int32(bits.RotateLeft32(uint32(x.StateA), 14)) + x.StateB)

	t := x.StateB << 9
	x.StateE += -0x3CA9B16B ^ x.StateD // 0xC3564E95 as signed int32
	x.StateC ^= x.StateA
	x.StateD ^= x.StateB
	x.StateB ^= x.StateC
	x.StateA ^= x.StateD
	x.StateC ^= t
	x.StateD = int32(bits.RotateLeft32(uint32(x.StateD), 11))

	return res
}

// Uint64 generates a 64-bit random value required by math/rand.Source64 and math/rand/v2.Source.
func (x *Xoshiro160RoadroxoRandom) Uint64() uint64 {
	hi := int32(bits.RotateLeft32(uint32(x.StateE), 23)) ^ (int32(bits.RotateLeft32(uint32(x.StateA), 14)) + x.StateB)
	lo := int32(bits.RotateLeft32(uint32(x.StateC), 19)) ^ (int32(bits.RotateLeft32(uint32(x.StateE), 7)) + x.StateD)

	t := x.StateB << 9
	x.StateE += -0x3CA9B16B ^ x.StateD
	x.StateC ^= x.StateA
	x.StateD ^= x.StateB
	x.StateB ^= x.StateC
	x.StateA ^= x.StateD
	x.StateC ^= t
	x.StateD = int32(bits.RotateLeft32(uint32(x.StateD), 11))

	return (uint64(uint32(hi)) << 32) ^ uint64(uint32(lo))
}

// Int63 returns a non-negative 63-bit integer required by math/rand.Source.
func (x *Xoshiro160RoadroxoRandom) Int63() int64 {
	return int64(x.Uint64() & 0x7FFFFFFFFFFFFFFF)
}

// PreviousInt rewinds the PRNG state by one step and returns the prior 32-bit output.
func (x *Xoshiro160RoadroxoRandom) PreviousInt() int32 {
	x.StateD = int32(bits.RotateLeft32(uint32(x.StateD), 21))
	x.StateA ^= x.StateD

	x.StateC ^= x.StateB
	x.StateC ^= x.StateC << 9
	x.StateC ^= x.StateC << 18

	x.StateB ^= x.StateA
	x.StateC ^= x.StateB
	x.StateB ^= x.StateC
	x.StateD ^= x.StateB

	x.StateE -= -0x3CA9B16B ^ x.StateD

	return int32(bits.RotateLeft32(uint32(x.StateE), 23)) ^ (int32(bits.RotateLeft32(uint32(x.StateA), 14)) + x.StateB)
}

// PreviousLong rewinds the PRNG state by one step and returns the prior 64-bit output.
func (x *Xoshiro160RoadroxoRandom) PreviousLong() int64 {
	x.StateD = int32(bits.RotateLeft32(uint32(x.StateD), 21))
	x.StateA ^= x.StateD

	x.StateC ^= x.StateB
	x.StateC ^= x.StateC << 9
	x.StateC ^= x.StateC << 18

	x.StateB ^= x.StateA
	x.StateC ^= x.StateB
	x.StateB ^= x.StateC
	x.StateD ^= x.StateB

	x.StateE -= -0x3CA9B16B ^ x.StateD

	hi := int32(bits.RotateLeft32(uint32(x.StateE), 23)) ^ (int32(bits.RotateLeft32(uint32(x.StateA), 14)) + x.StateB)
	lo := int32(bits.RotateLeft32(uint32(x.StateC), 19)) ^ (int32(bits.RotateLeft32(uint32(x.StateE), 7)) + x.StateD)

	return int64((uint64(uint32(hi)) << 32) ^ uint64(uint32(lo)))
}

// Leap jumps $2^{64}$ steps in the sequence and returns what Uint64() would output at that target state.
func (x *Xoshiro160RoadroxoRandom) Leap() uint64 {
	var s0, s1, s2, s3 int32
	jumpTable := [4]uint32{0x8764000b, 0xf542d2d3, 0x6fa035c3, 0x77f2db5b}

	for _, jump := range jumpTable {
		for b := 0; b < 32; b++ {
			if (jump & (1 << b)) != 0 {
				s0 ^= x.StateA
				s1 ^= x.StateB
				s2 ^= x.StateC
				s3 ^= x.StateD
			}
			t := x.StateB << 9
			x.StateE += -0x3CA9B16B ^ x.StateD
			x.StateC ^= x.StateA
			x.StateD ^= x.StateB
			x.StateB ^= x.StateC
			x.StateA ^= x.StateD
			x.StateC ^= t
			x.StateD = int32(bits.RotateLeft32(uint32(x.StateD), 11))
		}
	}

	x.StateA = s0
	x.StateB = s1
	x.StateC = s2
	x.StateD = s3

	s3 = int32(bits.RotateLeft32(uint32(s3), 21))
	s0 ^= s3
	s2 ^= s1
	s2 ^= s2 << 9
	s2 ^= s2 << 18
	s1 ^= s0
	s2 ^= s1
	s1 ^= s2
	s3 ^= s1

	s4 := x.StateE - (-0x3CA9B16B ^ s3)

	hi := int32(bits.RotateLeft32(uint32(s4), 23)) ^ (int32(bits.RotateLeft32(uint32(s0), 14)) + s1)
	lo := int32(bits.RotateLeft32(uint32(s2), 19)) ^ (int32(bits.RotateLeft32(uint32(s4), 7)) + s3)

	return (uint64(uint32(hi)) << 32) ^ uint64(uint32(lo))
}

// Copy creates an exact clone of the current generator state.
func (x *Xoshiro160RoadroxoRandom) Copy() *Xoshiro160RoadroxoRandom {
	return &Xoshiro160RoadroxoRandom{
		StateA: x.StateA,
		StateB: x.StateB,
		StateC: x.StateC,
		StateD: x.StateD,
		StateE: x.StateE,
	}
}
