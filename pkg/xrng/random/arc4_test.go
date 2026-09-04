package random

import (
	"math/rand"
	"testing"
)

// TestArc4InterfaceCompliance ensures Arc4Random fulfills math/rand.Source and math/rand.Source64.
func TestArc4InterfaceCompliance(t *testing.T) {
	var _ rand.Source = (*Arc4Random)(nil)
	var _ rand.Source64 = (*Arc4Random)(nil)

	// Verify integration with standard library rand.Rand
	src := NewArc4Random(42, Arc4ModeVMPC)
	r := rand.New(src)

	val := r.Int63()
	if val < 0 {
		t.Errorf("expected non-negative Int63, got %d", val)
	}

	val64 := r.Uint64()
	if val64 == 0 {
		t.Log("Uint64 generated 0 (valid, but noting edge case)")
	}
}

// TestArc4SeedDeterminism verifies that seeding with the same value produces deterministic output.
func TestArc4SeedDeterminism(t *testing.T) {
	seed := int64(123456789)

	r1 := NewArc4Random(seed, Arc4ModeVMPC)
	r2 := NewArc4Random(seed, Arc4ModeVMPC)

	for i := 0; i < 100; i++ {
		v1 := r1.Uint64()
		v2 := r2.Uint64()
		if v1 != v2 {
			t.Fatalf("mismatch at iteration %d: %d != %d", i, v1, v2)
		}
	}
}

// TestArc4DifferentSeedsProducesDifferentValues verifies distinct seeds yield different sequences.
func TestArc4DifferentSeedsProducesDifferentValues(t *testing.T) {
	r1 := NewArc4Random(100, Arc4ModeVMPC)
	r2 := NewArc4Random(200, Arc4ModeVMPC)

	if r1.Uint64() == r2.Uint64() {
		t.Error("expected different seeds to yield different initial outputs")
	}
}

// TestArc4Modes verifies key scheduling and output behavior for each mode.
func TestArc4Modes(t *testing.T) {
	modes := []struct {
		name string
		mode Arc4Mode
	}{
		{"ModeVMPC", Arc4ModeVMPC},
		{"ModeLegacy", Arc4ModeLegacy},
		{"ModeExtended", Arc4ModeExtended},
	}

	key := []byte("secret_key_12345")

	for _, tt := range modes {
		t.Run(tt.name, func(t *testing.T) {
			r := NewArc4Random(0, tt.mode)

			if r.GetMode() != tt.mode {
				t.Fatalf("expected mode %v, got %v", tt.mode, r.GetMode())
			}

			r.Ksa(key)

			// Ensure generator produces distinct values upon consecutive calls
			out1 := r.Uint64()
			out2 := r.Uint64()

			if out1 == out2 {
				t.Errorf("consecutive Uint64 calls generated identical values: %d", out1)
			}
		})
	}
}

// TestArc4ModeOutputDifferences verifies that different modes with the same key generate distinct output.
func TestArc4ModeOutputDifferences(t *testing.T) {
	key := []byte("shared_key")

	rVMPC := NewArc4Random(0, Arc4ModeVMPC)
	rVMPC.Ksa(key)

	rLegacy := NewArc4Random(0, Arc4ModeLegacy)
	rLegacy.Ksa(key)

	rExtended := NewArc4Random(0, Arc4ModeExtended)
	rExtended.Ksa(key)

	vmpcVal := rVMPC.Uint64()
	legacyVal := rLegacy.Uint64()
	extendedVal := rExtended.Uint64()

	if vmpcVal == legacyVal || vmpcVal == extendedVal || legacyVal == extendedVal {
		t.Errorf("expected different mode outputs; got VMPC: %d, Legacy: %d, Extended: %d",
			vmpcVal, legacyVal, extendedVal)
	}
}

// TestArc4NewArc4RandomFromBytes tests initialization directly from a byte slice.
func TestArc4NewArc4RandomFromBytes(t *testing.T) {
	rawKey := []byte("sample_raw_bytes")

	r1 := NewArc4RandomFromKey(rawKey)
	r2 := NewArc4RandomFromKey(rawKey)

	for i := 0; i < 50; i++ {
		if r1.Uint64() != r2.Uint64() {
			t.Fatalf("mismatch on byte-initialized generators at index %d", i)
		}
	}
}

// TestArc4GetContext verifies that GetContext returns a valid 256-element array copy.
func TestArc4GetContext(t *testing.T) {
	r := NewArc4Random(99, Arc4ModeVMPC)
	ctx1 := r.GetContext()

	if len(ctx1) != 256 {
		t.Fatalf("expected context length 256, got %d", len(ctx1))
	}

	// Mutate returned context slice to ensure deep copy protection
	ctx1[0] = -999
	ctx2 := r.GetContext()

	if ctx2[0] == -999 {
		t.Error("GetContext returned a reference instead of a copy")
	}
}

// TestArc4EmptyKsa handles zero-length byte inputs gracefully without panicking.
func TestArc4EmptyKsa(t *testing.T) {
	r := NewArc4Random(0, Arc4ModeVMPC)
	ctxBefore := r.GetContext()

	r.Ksa([]byte{})
	ctxAfter := r.GetContext()

	for i := range ctxBefore {
		if ctxBefore[i] != ctxAfter[i] {
			t.Errorf("state changed on empty Ksa at index %d", i)
		}
	}
}

// BenchmarkArc4Uint64 measures Uint64 output generation performance.
func BenchmarkArc4Uint64(b *testing.B) {
	r := NewArc4Random(42, Arc4ModeVMPC)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = r.Uint64()
	}
}

// TestArc4Vectors evaluates the generator output against known golden test vectors.
func TestArc4Vectors(t *testing.T) {
	type testVector struct {
		name     string
		mode     Arc4Mode
		seed     int64
		rawKey   []byte
		expected []uint64 // Fill in the expected sequential Uint64() outputs
	}

	// Placeholders: Replace expected values with your actual test vector outputs
	vectors := []testVector{
		{
			name:   "Legacy Mode - Seeded",
			mode:   Arc4ModeLegacy,
			seed:   42,
			rawKey: nil,
			expected: []uint64{
				0x871B8D29C81C2A10,
			},
		},
		{
			name:   "Legacy Mode - Key#1",
			mode:   Arc4ModeLegacy,
			seed:   0,
			rawKey: []byte("Key"),
			expected: []uint64{
				0xEB9F7781B734CA72,
			},
		},
		{
			name:   "Legacy Mode - Key#2",
			mode:   Arc4ModeLegacy,
			seed:   0,
			rawKey: []byte("Wiki"),
			expected: []uint64{
				0x6044DB6D41B7E8E7,
			},
		},
		{
			name:   "Legacy Mode - Key#3",
			mode:   Arc4ModeLegacy,
			seed:   0,
			rawKey: []byte("Secret"),
			expected: []uint64{
				0x04D46B053CA87B59,
			},
		},
		{
			name:   "VMPC Mode - Seeded",
			mode:   Arc4ModeVMPC,
			seed:   12345,
			rawKey: nil,
			expected: []uint64{
				0xEA48F1FC7C8ADDE0,
			},
		},
		{
			name:   "VMPC Mode - Keyed",
			mode:   Arc4ModeVMPC,
			seed:   0,
			rawKey: []byte("TestKey123"),
			expected: []uint64{
				0xDB89599B1F4DF0AD,
			},
		},
		{
			name:   "Extended Mode - Seeded",
			mode:   Arc4ModeExtended,
			seed:   42,
			rawKey: nil,
			expected: []uint64{
				0xF12C5F6D2E9AB632,
			},
		},
		{
			name:   "Extended Mode - Keyed",
			mode:   Arc4ModeExtended,
			seed:   0,
			rawKey: []byte("ExtendedKey789"),
			expected: []uint64{
				0x75F791B3AC5AFA39,
				0x2A78A21A7C75C3D3,
			},
		},
	}

	for _, tv := range vectors {
		t.Run(tv.name, func(t *testing.T) {
			if len(tv.expected) == 0 {
				t.Skip("Placeholder test vector empty - fill in expected slice to enable")
			}

			var r *Arc4Random
			if len(tv.rawKey) > 0 {
				r = NewArc4Random(0, tv.mode)
				r.Ksa(tv.rawKey)
			} else {
				r = NewArc4Random(tv.seed, tv.mode)
			}

			for i, exp := range tv.expected {
				got := r.Uint64()
				if got != exp {
					t.Errorf("step %d: expected 0x%016X, got 0x%016X", i, exp, got)
				}
			}
		})
	}
}
