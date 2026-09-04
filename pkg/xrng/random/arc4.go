package random

import (
	"encoding/binary"
	"math/rand"
)

// Arc4Mode represents the PRNG key-scheduling and generation variant.
type Arc4Mode int

const (
	Arc4ModeLegacy Arc4Mode = iota
	Arc4ModeVMPC
	Arc4ModeExtended
)

type Arc4Random struct {
	ctx  [256]int
	a    int
	b    int
	mode Arc4Mode
}

var (
	_ rand.Source   = (*Arc4Random)(nil)
	_ rand.Source64 = (*Arc4Random)(nil)
)

func NewArc4Legacy(s int64) *Arc4Random {
	return NewArc4Random(s, Arc4ModeLegacy)
}

func NewArc4Vmpc(s int64) *Arc4Random {
	return NewArc4Random(s, Arc4ModeVMPC)
}

func NewArc4Extended(s int64) *Arc4Random {
	return NewArc4Random(s, Arc4ModeExtended)
}

func NewArc4LegacyFromKey(s []byte) *Arc4Random {
	r := NewArc4Random(0, Arc4ModeLegacy)
	r.Ksa(s)
	return r
}

func NewArc4VmpcFromKey(s []byte) *Arc4Random {
	r := NewArc4Random(0, Arc4ModeVMPC)
	r.Ksa(s)
	return r
}

func NewArc4ExtendedFromKey(s []byte) *Arc4Random {
	r := NewArc4Random(0, Arc4ModeExtended)
	r.Ksa(s)
	return r
}

// NewArc4Random constructs an Arc4Random generator initialized in Arc4ModeVMPC.
func NewArc4Random(s int64, m Arc4Mode) *Arc4Random {
	r := &Arc4Random{mode: m}
	r.Seed(s)
	return r
}

// NewArc4RandomFromKey constructs an Arc4Random generator from a byte slice.
func NewArc4RandomFromKey(s []byte) *Arc4Random {
	r := NewArc4Random(0, Arc4ModeVMPC)
	r.Ksa(s)
	return r
}

func (r *Arc4Random) GetMode() Arc4Mode {
	return r.mode
}

func (r *Arc4Random) SetMode(m Arc4Mode) {
	r.mode = m
}

func (r *Arc4Random) GetContext() []int {
	ctxCopy := make([]int, 256)
	copy(ctxCopy, r.ctx[:])
	return ctxCopy
}

// Seed implements rand.Source.
func (r *Arc4Random) Seed(seed int64) {
	for i := 0; i < 256; i++ {
		r.ctx[i] = i
	}
	r.a = 0
	r.b = 0

	if seed == 0 {
		return
	}

	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf[:], uint64(seed))

	r.Ksa(buf)
}

func (r *Arc4Random) Ksa(key []byte) {
	keylen := len(key)
	if keylen == 0 {
		return
	}

	il := 256 // Arc4ModeLegacy
	switch r.mode {
	case Arc4ModeLegacy:
		il = 256
	case Arc4ModeVMPC:
		il = 768
	case Arc4ModeExtended:
		il = 4096
	}

	j := 0
	for i := 0; i < il; i++ {
		j = (j + r.ctx[i&0xff] + int(key[i%keylen])) & 0xff
		r.ctx[i&0xff], r.ctx[j] = r.ctx[j], r.ctx[i&0xff]
	}
	r.a = 0
	r.b = 0
	if r.mode == Arc4ModeExtended {
		for i := 0; i < il; i++ {
			r.Uint64()
		}
	}
}

// Uint64 implements rand.Source64.
func (r *Arc4Random) Uint64() uint64 {
	next8 := func() int {
		if r.mode == Arc4ModeVMPC {
			r.b = r.ctx[(r.b+r.ctx[r.a&0xff])&0xff]
			ret := r.ctx[(r.ctx[r.ctx[r.b&0xff]&0xff]+1)&0xff]
			r.ctx[r.b&0xff], r.ctx[r.a&0xff] = r.ctx[r.a&0xff], r.ctx[r.b&0xff]
			r.a = (r.a + 1) & 0xff
			return ret
		}

		// Standard ARC4 output generation for ModeLegacy and ModeExtended
		r.a = (r.a + 1) & 0xff
		r.b = (r.b + r.ctx[r.a]) & 0xff
		r.ctx[r.b], r.ctx[r.a] = r.ctx[r.a], r.ctx[r.b]
		return r.ctx[(r.ctx[r.b]+r.ctx[r.a])&0xff]
	}

	var next64 uint64 = 0
	for i := 0; i < 8; i++ {
		next64 = (next64 << 8) | uint64(next8()&0xff)
	}
	return next64
}

// Int63 implements rand.Source.
func (r *Arc4Random) Int63() int64 {
	return int64(r.Uint64() & 0x7fffffffffffffff)
}
