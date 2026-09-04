package random

import "math/rand"

const (
	lehmer_magic1 = 0xe120fc15
	lehmer_magic2 = 0x4a39b70d
	lehmer_magic3 = 0x12fad5c9
)

type Lehmer32Random struct {
	pstate uint64
}

var (
	_ rand.Source   = (*Lehmer32Random)(nil)
	_ rand.Source64 = (*Lehmer32Random)(nil)
)

// NewLehmer32Random constructs a Lehmer32Random generator initialized with a seed.
func NewLehmer32Random(seed int64) *Lehmer32Random {
	l := &Lehmer32Random{}
	l.Seed(seed)
	return l
}

// Seed initializes the PRNG state (implements rand.Source).
func (l *Lehmer32Random) Seed(seed int64) {
	l.pstate = uint64(seed)
}

// NextInt generates a 32-bit pseudo-random integer.
func (l *Lehmer32Random) NextInt() uint32 {
	l.pstate = (l.pstate + lehmer_magic1) & 0xffffffff
	l.pstate = ((l.pstate >> 32) ^ l.pstate) & 0xffffffff

	tmp := (l.pstate * lehmer_magic2) & 0xffffffff
	tmp = ((tmp >> 32) ^ tmp) & 0xffffffff
	tmp = (tmp * lehmer_magic3) & 0xffffffff
	tmp = ((tmp >> 32) ^ tmp) & 0xffffffff

	return uint32(tmp)
}

// Uint64 generates a 64-bit integer (implements rand.Source64).
func (l *Lehmer32Random) Uint64() uint64 {
	hi := uint64(l.NextInt())
	lo := uint64(l.NextInt())
	return (hi << 32) | lo
}

// Int63 generates a non-negative 63-bit integer (implements rand.Source).
func (l *Lehmer32Random) Int63() int64 {
	return int64(l.Uint64() & 0x7fffffffffffffff)
}
