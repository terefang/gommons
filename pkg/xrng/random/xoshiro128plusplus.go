package random

import (
	"math/bits"
	"math/rand"
)

// Xoshiro128PlusPlusRandom implements the xoshiro128++ 32-bit PRNG.
// It complies with math/rand.Source, math/rand.Source64, and math/rand/v2.Source.
type Xoshiro128PlusPlusRandom struct {
	StateA int32
	StateB int32
	StateC int32
	StateD int32
}

// Compile-time assertion for interface adherence.
var (
	_ rand.Source   = (*Xoshiro128PlusPlusRandom)(nil)
	_ rand.Source64 = (*Xoshiro128PlusPlusRandom)(nil)
)

// NewXoshiro128PlusPlusRandom constructs a generator initialized with a 64-bit seed.
func NewXoshiro128PlusPlusRandom(seed int64) *Xoshiro128PlusPlusRandom {
	x := &Xoshiro128PlusPlusRandom{}
	x.Seed(seed)
	return x
}

// NewXoshiro128PlusPlusRandomWithStates constructs a generator using explicit state values.
func NewXoshiro128PlusPlusRandomWithStates(a, b, c, d int32) *Xoshiro128PlusPlusRandom {
	if (a | b | c | d) == 0 {
		d = 1
	}
	return &Xoshiro128PlusPlusRandom{
		StateA: a,
		StateB: b,
		StateC: c,
		StateD: d,
	}
}

// Seed seeds all 4 internal 32-bit state variables using SplitMix64 operations.
func (x *Xoshiro128PlusPlusRandom) Seed(seed int64) {
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

	if (x.StateA | x.StateB | x.StateC | x.StateD) == 0 {
		x.StateD = 1
	}
}

// NextInt advances state and returns the next pseudo-random 32-bit signed integer.
func (x *Xoshiro128PlusPlusRandom) NextInt() int32 {
	res := x.StateA + x.StateD
	res = int32(bits.RotateLeft32(uint32(res), 7)) + x.StateA

	t := x.StateB << 9
	x.StateC ^= x.StateA
	x.StateD ^= x.StateB
	x.StateB ^= x.StateC
	x.StateA ^= x.StateD
	x.StateC ^= t
	x.StateD = int32(bits.RotateLeft32(uint32(x.StateD), 11))

	return res
}

// Uint64 generates a 64-bit random value by executing two 32-bit steps.
// It provides direct compatibility for math/rand.Source64 and math/rand/v2.Source.
func (x *Xoshiro128PlusPlusRandom) Uint64() uint64 {
	h := x.StateA + x.StateD
	h = int32(bits.RotateLeft32(uint32(h), 7)) + x.StateA

	l := x.StateC - x.StateB
	l = int32(bits.RotateLeft32(uint32(l), 13)) + x.StateC

	t := x.StateB << 9
	x.StateC ^= x.StateA
	x.StateD ^= x.StateB
	x.StateB ^= x.StateC
	x.StateA ^= x.StateD
	x.StateC ^= t
	x.StateD = int32(bits.RotateLeft32(uint32(x.StateD), 11))

	return (uint64(uint32(h)) << 32) | uint64(uint32(l))
}

// Int63 returns a non-negative 63-bit integer required by math/rand.Source.
func (x *Xoshiro128PlusPlusRandom) Int63() int64 {
	return int64(x.Uint64() & 0x7FFFFFFFFFFFFFFF)
}

// PreviousInt rewinds the PRNG state by one step and returns the prior output.
func (x *Xoshiro128PlusPlusRandom) PreviousInt() int32 {
	x.StateD = int32(bits.RotateLeft32(uint32(x.StateD), 21))
	pa := x.StateA ^ x.StateD
	x.StateA = pa

	x.StateC ^= x.StateB
	x.StateC ^= x.StateC << 9
	x.StateC ^= x.StateC << 18

	x.StateB ^= x.StateA
	x.StateC ^= x.StateB
	x.StateB ^= x.StateC
	pd := x.StateD ^ x.StateB
	x.StateD = pd

	pd = pa + pd
	pd = int32(bits.RotateLeft32(uint32(pd), 7)) + pa
	return pd
}

// PreviousLong rewinds the PRNG state by one step and returns the prior 64-bit combined value.
func (x *Xoshiro128PlusPlusRandom) PreviousLong() int64 {
	x.StateD = int32(bits.RotateLeft32(uint32(x.StateD), 21))
	pa := x.StateA ^ x.StateD
	x.StateA = pa

	x.StateC ^= x.StateB
	x.StateC ^= x.StateC << 9
	x.StateC ^= x.StateC << 18

	x.StateB ^= x.StateA
	pc := x.StateC ^ x.StateB
	x.StateC = pc
	pb := x.StateB ^ x.StateC
	x.StateB = pb
	pd := x.StateD ^ x.StateB
	x.StateD = pd

	pd = pa + pd
	pd = int32(bits.RotateLeft32(uint32(pd), 7)) + pa
	pb = pc - pb
	pb = int32(bits.RotateLeft32(uint32(pb), 13)) + pc

	return int64((uint64(uint32(pd)) << 32) | uint64(uint32(pb)))
}

// Leap advances the state by $2^{64}$ steps in a single call.
func (x *Xoshiro128PlusPlusRandom) Leap() int64 {
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

	s3 = s0 + s3
	s3 = int32(bits.RotateLeft32(uint32(s3), 7)) + s0
	s1 = s2 - s1
	s1 = int32(bits.RotateLeft32(uint32(s1), 13)) + s2

	return int64((uint64(uint32(s3)) << 32) | uint64(uint32(s1)))
}

// Copy creates an independent deep clone of the current generator state.
func (x *Xoshiro128PlusPlusRandom) Copy() *Xoshiro128PlusPlusRandom {
	return &Xoshiro128PlusPlusRandom{
		StateA: x.StateA,
		StateB: x.StateB,
		StateC: x.StateC,
		StateD: x.StateD,
	}
}
