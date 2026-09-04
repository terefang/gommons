package random

import (
	"fmt"
	"math/bits"
	"math/rand"
)

// AceRandom is a pseudo-random number generator with 5 64-bit states.
// It implements math/rand.Source, math/rand.Source64, and math/rand/v2.Source.
type AceRandom struct {
	StateA uint64
	StateB uint64
	StateC uint64
	StateD uint64
	StateE uint64
}

// Compile-time checks to ensure interface adherence.
var (
	_ rand.Source   = (*AceRandom)(nil)
	_ rand.Source64 = (*AceRandom)(nil)
)

// NewAceRandom creates an AceRandom instance with default zero states.
func NewAceRandom() *AceRandom {
	return &AceRandom{}
}

// NewAceRandomWithSeed initializes a new AceRandom using a single 64-bit seed.
func NewAceRandomWithSeed(seed int64) *AceRandom {
	r := &AceRandom{}
	r.Seed(seed)
	return r
}

// NewAceRandom2States constructs AceRandom using two explicit state values.
func NewAceRandom2States(stateA, stateB uint64) *AceRandom {
	return &AceRandom{
		StateA: stateA,
		StateB: stateB,
		StateC: stateA + stateB,
		StateD: stateA ^ stateB,
		StateE: stateB - stateA,
	}
}

// NewAceRandom3States constructs AceRandom using three explicit state values.
func NewAceRandom3States(stateA, stateB, stateC uint64) *AceRandom {
	return &AceRandom{
		StateA: stateA,
		StateB: stateB,
		StateC: stateC,
		StateD: stateA + stateC,
		StateE: stateB ^ stateC,
	}
}

// NewAceRandom4States constructs AceRandom using four explicit state values.
func NewAceRandom4States(stateA, stateB, stateC, stateD uint64) *AceRandom {
	return &AceRandom{
		StateA: stateA,
		StateB: stateB,
		StateC: stateC,
		StateD: stateD,
		StateE: (stateA + stateC) ^ (stateB + stateD),
	}
}

// NewAceRandom5States constructs AceRandom using five explicit state values.
func NewAceRandom5States(stateA, stateB, stateC, stateD, stateE uint64) *AceRandom {
	return &AceRandom{
		StateA: stateA,
		StateB: stateB,
		StateC: stateC,
		StateD: stateD,
		StateE: stateE,
	}
}

// Seed initializes all 5 states using the seed (satisfies math/rand.Source).
func (r *AceRandom) Seed(seed int64) {
	s := (uint64(seed) ^ 0x1C69B3F74AC4AE35) * 0x3C79AC492BA7B653
	r.StateA = s ^ ^uint64(0xC6BC279692B5C323)
	s ^= s >> 32
	r.StateB = s ^ 0xD3833E804F4C574B
	s *= 0xBEA225F9EB34556D
	s ^= s >> 29
	r.StateC = s ^ ^uint64(0xD3833E804F4C574B)
	s *= 0xBEA225F9EB34556D
	s ^= s >> 32
	r.StateD = s ^ 0xC6BC279692B5C323
	s *= 0xBEA225F9EB34556D
	s ^= s >> 29
	r.StateE = s
}

// Uint64 generates a unsigned 64-bit integer (satisfies math/rand.Source64 and math/rand/v2.Source).
func (r *AceRandom) Uint64() uint64 {
	fa, fb, fc, fd, fe := r.StateA, r.StateB, r.StateC, r.StateD, r.StateE
	r.StateA = fa + 0x9E3779B97F4A7C15
	r.StateB = fa ^ fe
	r.StateC = fb + fd
	r.StateD = bits.RotateLeft64(fc, 52)
	r.StateE = fb - fc
	return r.StateE
}

// Int63 returns a non-negative 63-bit integer (satisfies math/rand.Source).
func (r *AceRandom) Int63() int64 {
	return int64(r.Uint64() & 0x7FFFFFFFFFFFFFFF)
}

// NextLong is an alias for Uint64, returning signed int64 as in Java.
func (r *AceRandom) NextLong() int64 {
	return int64(r.Uint64())
}

// PreviousLong steps backwards one iteration and returns the previous output value.
func (r *AceRandom) PreviousLong() int64 {
	fb, fc, fd, fe := r.StateB, r.StateC, r.StateD, r.StateE
	r.StateA -= 0x9E3779B97F4A7C15
	r.StateC = bits.RotateLeft64(fd, -52) // equivalent to (fd >>> 52 | fd << 12)
	r.StateB = r.StateC + fe
	r.StateD = fc - r.StateB
	r.StateE = fb ^ r.StateA
	return int64(fe)
}

// NextBits returns the highest 'bitsCount' bits of the next long output.
func (r *AceRandom) NextBits(bitsCount uint) uint32 {
	return uint32(r.Uint64() >> (64 - bitsCount))
}

// Leap performs a large jump forward ($2^{48}$ steps) in the sequence.
func (r *AceRandom) Leap() int64 {
	fa, fb, fc, fd, fe := r.StateA, r.StateB, r.StateC, r.StateD, r.StateE
	r.StateA = fa + 0x7C15000000000000
	r.StateB = fa ^ fe
	r.StateC = fb + fd
	r.StateD = bits.RotateLeft64(fc, 52)
	r.StateE = fb - fc
	return int64(r.StateE)
}

// Copy creates a deep clone of the current generator state.
func (r *AceRandom) Copy() *AceRandom {
	return &AceRandom{
		StateA: r.StateA,
		StateB: r.StateB,
		StateC: r.StateC,
		StateD: r.StateD,
		StateE: r.StateE,
	}
}

// Equals checks if two AceRandom instances have identical state.
func (r *AceRandom) Equals(other *AceRandom) bool {
	if other == nil {
		return false
	}
	return r.StateA == other.StateA &&
		r.StateB == other.StateB &&
		r.StateC == other.StateC &&
		r.StateD == other.StateD &&
		r.StateE == other.StateE
}

func (r *AceRandom) String() string {
	return fmt.Sprintf("AceRandom{stateA=%dL, stateB=%dL, stateC=%dL, stateD=%dL, stateE=%dL}",
		int64(r.StateA), int64(r.StateB), int64(r.StateC), int64(r.StateD), int64(r.StateE))
}
