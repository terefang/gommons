package random

import (
	"math/bits"
	"math/rand"
)

// Jsf32Random implements Bob Jenkins' Small Fast Generator (32-bit PRNG).
// It implements math/rand.Source, math/rand.Source64 (Go 1.8+), and math/rand/v2.Source.
type Jsf32Random struct {
	StateA int32
	StateB int32
	StateC int32
	StateD int32
}

// Interface compliance checks
var (
	_ rand.Source   = (*Jsf32Random)(nil)
	_ rand.Source64 = (*Jsf32Random)(nil)
)

// NewJsf32Random creates a generator seeded with a 64-bit value.
func NewJsf32Random(seed int64) *Jsf32Random {
	j := &Jsf32Random{}
	j.Seed(seed)
	return j
}

// NewJsf32RandomWithStates creates a generator using four explicit state values.
func NewJsf32RandomWithStates(a, b, c, d int32) *Jsf32Random {
	return &Jsf32Random{
		StateA: a,
		StateB: b,
		StateC: c,
		StateD: d,
	}
}

// Seed initializes all 4 states using the Java implementation's seeding strategy.
func (j *Jsf32Random) Seed(seed int64) {
	j.StateA = -0x0E15A113 // 0xF1EA5EED
	y := int32(seed)

	if seed == int64(y) {
		j.StateB = y
		j.StateC = y
		j.StateD = y
		for i := 0; i < 20; i++ {
			j.NextInt()
		}
		return
	}

	uSeed := uint64(seed)

	x := uSeed
	x ^= x >> 27
	x *= 0x3C79AC492BA7B653
	x ^= x >> 33
	x *= 0x1C69B3F74AC4AE35
	j.StateB = int32(x ^ (x >> 27))

	x = uSeed + 0x9E3779B97F4A7C15
	x ^= x >> 27
	x *= 0x3C79AC492BA7B653
	x ^= x >> 33
	x *= 0x1C69B3F74AC4AE35
	x ^= x >> 27

	j.StateC = int32(x)
	j.StateD = int32(x >> 32)
}

// NextInt advances state and returns the next pseudo-random 32-bit signed integer.
func (j *Jsf32Random) NextInt() int32 {
	e := j.StateA - int32(bits.RotateLeft32(uint32(j.StateB), 27))
	j.StateA = j.StateB ^ int32(bits.RotateLeft32(uint32(j.StateC), 17))
	j.StateB = j.StateC + j.StateD
	j.StateC = j.StateD + e
	j.StateD = e + j.StateA
	return j.StateD
}

// Uint64 returns a 64-bit pseudo-random value required by math/rand.Source64 and math/rand/v2.Source.
func (j *Jsf32Random) Uint64() uint64 {
	e1 := j.StateA - int32(bits.RotateLeft32(uint32(j.StateB), 27))
	j.StateA = j.StateB ^ int32(bits.RotateLeft32(uint32(j.StateC), 17))
	j.StateB = j.StateC + j.StateD
	j.StateC = j.StateD + e1
	j.StateD = e1 + j.StateA
	h := j.StateD

	e2 := j.StateA - int32(bits.RotateLeft32(uint32(j.StateB), 27))
	j.StateA = j.StateB ^ int32(bits.RotateLeft32(uint32(j.StateC), 17))
	j.StateB = j.StateC + j.StateD
	j.StateC = j.StateD + e2
	j.StateD = e2 + j.StateA
	l := j.StateD

	return (uint64(uint32(h)) << 32) | uint64(uint32(l))
}

// Int63 returns a non-negative 63-bit integer required by math/rand.Source.
func (j *Jsf32Random) Int63() int64 {
	return int64(j.Uint64() & 0x7FFFFFFFFFFFFFFF)
}

// PreviousInt steps back one state transition and returns the prior stateD value.
func (j *Jsf32Random) PreviousInt() int32 {
	l := j.StateD
	e := j.StateD - j.StateA
	j.StateD = j.StateC - e
	j.StateC = j.StateB - j.StateD
	j.StateB = j.StateA ^ int32(bits.RotateLeft32(uint32(j.StateC), 17))
	j.StateA = e + int32(bits.RotateLeft32(uint32(j.StateB), 27))
	return l
}

// PreviousLong steps back two state transitions and returns the prior 64-bit output.
func (j *Jsf32Random) PreviousLong() int64 {
	l := j.StateD
	e1 := j.StateD - j.StateA
	h := j.StateC - e1
	j.StateD = h
	j.StateC = j.StateB - j.StateD
	j.StateB = j.StateA ^ int32(bits.RotateLeft32(uint32(j.StateC), 17))
	j.StateA = e1 + int32(bits.RotateLeft32(uint32(j.StateB), 27))

	e2 := j.StateD - j.StateA
	j.StateD = j.StateC - e2
	j.StateC = j.StateB - j.StateD
	j.StateB = j.StateA ^ int32(bits.RotateLeft32(uint32(j.StateC), 17))
	j.StateA = e2 + int32(bits.RotateLeft32(uint32(j.StateB), 27))

	return int64((uint64(uint32(h)) << 32) | uint64(uint32(l)))
}

// Copy creates an exact duplicate of this generator's current state.
func (j *Jsf32Random) Copy() *Jsf32Random {
	return &Jsf32Random{
		StateA: j.StateA,
		StateB: j.StateB,
		StateC: j.StateC,
		StateD: j.StateD,
	}
}
