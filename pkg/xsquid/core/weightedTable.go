package core

import (
	"errors"
	"math/rand/v2"
	"slices"
)

// WeightedTable implements Vose's Alias Method to look up random indices based on positive float weights.
// Construction runs in O(n) time, and subsequent random index selection runs in O(1) constant time.
type WeightedTable struct {
	mixed []int32
	size  int
}

// NewWeightedTable creates a WeightedTable from a variadic slice of probability weights.
func NewWeightedTable(probabilities ...float32) (*WeightedTable, error) {
	return NewWeightedTableSlice(probabilities, 0, len(probabilities))
}

// NewWeightedTableSlice constructs a WeightedTable with a subset slice of probabilities defined by offset and length.
func NewWeightedTableSlice(probabilities []float32, offset, length int) (*WeightedTable, error) {
	if offset < 0 || length <= 0 || offset+length > len(probabilities) {
		return nil, errors.New("invalid offset or length for probabilities slice")
	}

	size := length
	if size == 0 {
		return nil, errors.New("array 'probabilities' given to WeightedTable must be nonempty")
	}

	wt := &WeightedTable{
		size:  size,
		mixed: make([]int32, size<<1),
	}

	var sum float32 = 0.0
	probs := make([]float32, size)

	for i, idx := offset, 0; idx < size; i, idx = i+1, idx+1 {
		if probabilities[i] <= 0 {
			continue
		}
		probs[idx] = probabilities[i]
		sum += probabilities[i]
	}

	if sum <= 0 {
		return nil, errors.New("at least one probability must be positive")
	}

	average := sum / float32(size)
	invAverage := 1.0 / average

	// Create two worklists (stacks)
	small := make([]int, 0, size)
	large := make([]int, 0, size)

	for i := 0; i < size; i++ {
		if probs[i] >= average {
			large = append(large, i)
		} else {
			small = append(small, i)
		}
	}

	for len(small) > 0 && len(large) > 0 {
		// Pop worklist elements
		less := small[len(small)-1]
		small = small[:len(small)-1]
		less2 := less << 1

		more := large[len(large)-1]
		large = large[:len(large)-1]

		wt.mixed[less2] = int32(0x7FFFFFFF * (probs[less] * invAverage))
		wt.mixed[less2|1] = int32(more)

		probs[more] += probs[less] - average

		if probs[more] >= average {
			large = append(large, more)
		} else {
			small = append(small, more)
		}
	}

	for len(small) > 0 {
		less := small[len(small)-1]
		small = small[:len(small)-1]
		wt.mixed[less<<1] = 0x7FFFFFFF
	}

	for len(large) > 0 {
		more := large[len(large)-1]
		large = large[:len(large)-1]
		wt.mixed[more<<1] = 0x7FFFFFFF
	}

	return wt, nil
}

// Copy creates an exact duplicate of the WeightedTable without sharing underlying slice backing array state.
func (wt *WeightedTable) Copy() *WeightedTable {
	if wt == nil {
		return nil
	}
	mixedCopy := make([]int32, len(wt.mixed))
	copy(mixedCopy, wt.mixed)

	return &WeightedTable{
		mixed: mixedCopy,
		size:  wt.size,
	}
}

// Size returns the number of index choices in the weighted table.
func (wt *WeightedTable) Size() int {
	return wt.size
}

// Random returns a random index [0, Size) determined deterministically by the given state variable using MX3 unary hashing.
func (wt *WeightedTable) Random(state uint64) int {
	// MX3 Unary Hash Algorithm
	state = wtRandomize3(state)

	column := int((uint64(wt.size) * (state & 0xFFFFFFFF)) >> 32)
	probThreshold := wt.mixed[column<<1]

	if int32(state>>33) <= probThreshold {
		return column
	}
	return int(wt.mixed[(column<<1)|1])
}

// RandomFromSource returns a random index [0, Size) using a Go standard rand.Source.
func (wt *WeightedTable) RandomFromSource(gen rand.Source) int {
	r := rand.New(gen)
	state := r.Uint64()

	column := int((uint64(wt.size) * (state & 0xFFFFFFFF)) >> 32)
	probThreshold := wt.mixed[column<<1]

	if int32(state>>33) <= probThreshold {
		return column
	}
	return int(wt.mixed[(column<<1)|1])
}

// Equals checks for structural and internal content equality with another WeightedTable.
func (wt *WeightedTable) Equals(other *WeightedTable) bool {
	if wt == other {
		return true
	}
	if wt == nil || other == nil {
		return false
	}
	if wt.size != other.size {
		return false
	}
	return slices.Equal(wt.mixed, other.mixed)
}

// HashCode returns a 64-bit hash code of the table content matching original Hasher behavior.
func (wt *WeightedTable) HashCode() uint64 {
	h := uint64(wt.size) ^ 0xFEDCBA9876543210
	for _, v := range wt.mixed {
		h ^= uint64(v)
		h *= 0xBF58476D1CE4E5B9
		h ^= h >> 32
	}
	return h
}

// wtRandomize3 implements MX3 unary hash function.
func wtRandomize3(x uint64) uint64 {
	x ^= x >> 32
	x *= 0xbea225f9eb34556d
	x ^= x >> 29
	x *= 0xbea225f9eb34556d
	x ^= x >> 32
	x *= 0xbea225f9eb34556d
	x ^= x >> 29
	return x
}
