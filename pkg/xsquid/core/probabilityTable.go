package core

import (
	"math/rand/v2"
	"reflect"
	"slices"
)

// ProbabilityTable holds a probability table to determine weighted random outcomes, including support for nested tables.
type ProbabilityTable[T comparable] struct {
	Table      []T
	ExtraTable []*ProbabilityTable[T]
	Weights    []int
	Total      int
	RNG        rand.Source
}

// NewProbabilityTable creates an empty ProbabilityTable with a default random source.
func NewProbabilityTable[T comparable]() *ProbabilityTable[T] {
	return NewProbabilityTableWithRNG[T](rand.NewPCG(rand.Uint64(), rand.Uint64()))
}

// NewProbabilityTableWithRNG creates a ProbabilityTable with a specified rand.Source.
func NewProbabilityTableWithRNG[T comparable](rng rand.Source) *ProbabilityTable[T] {
	if rng == nil {
		rng = rand.NewPCG(rand.Uint64(), rand.Uint64())
	}
	return &ProbabilityTable[T]{
		Table:      make([]T, 0, 64),
		ExtraTable: make([]*ProbabilityTable[T], 0, 16),
		Weights:    make([]int, 0, 64),
		Total:      0,
		RNG:        rng,
	}
}

// Copy creates a deep copy of the ProbabilityTable.
// Note: T items are copied by value (not deep-copied if pointers/references).
func (pt *ProbabilityTable[T]) Copy() *ProbabilityTable[T] {
	if pt == nil {
		return nil
	}

	extraCopy := make([]*ProbabilityTable[T], len(pt.ExtraTable))
	for i, extra := range pt.ExtraTable {
		extraCopy[i] = extra.Copy()
	}

	return &ProbabilityTable[T]{
		Table:      slices.Clone(pt.Table),
		ExtraTable: extraCopy,
		Weights:    slices.Clone(pt.Weights),
		Total:      pt.Total,
		RNG:        pt.RNG,
	}
}

// Random returns a randomly selected item based on assigned weights.
// Returns the zero-value of T and false if the table is empty.
func (pt *ProbabilityTable[T]) Random() (T, bool) {
	var zero T
	if len(pt.Table) == 0 && len(pt.ExtraTable) == 0 || pt.Total <= 0 {
		return zero, false
	}

	r := rand.New(pt.RNG)
	index := r.IntN(pt.Total)
	sz := len(pt.Table)

	for i := 0; i < sz; i++ {
		index -= pt.Weights[i]
		if index < 0 {
			return pt.Table[i], true
		}
	}

	for i := 0; i < len(pt.ExtraTable); i++ {
		index -= pt.Weights[sz+i]
		if index < 0 {
			return pt.ExtraTable[i].Random()
		}
	}

	return zero, false
}

// Add adds or updates the weight for a given item in the table. Weight must be > 0.
func (pt *ProbabilityTable[T]) Add(item T, weight int) *ProbabilityTable[T] {
	if weight <= 0 {
		return pt
	}

	idx := slices.Index(pt.Table, item)
	if idx < 0 {
		pt.Table = append(pt.Table, item)
		// Insert weight right after table elements (before extraTable weights)
		insertIdx := len(pt.Table) - 1
		pt.Weights = slices.Insert(pt.Weights, insertIdx, weight)
		pt.Total += weight
	} else {
		oldWeight := pt.Weights[idx]
		newWeight := max(0, oldWeight+weight)
		pt.Weights[idx] = newWeight
		pt.Total += newWeight - oldWeight
	}
	return pt
}

// AddAll adds multiple item-weight pairs from a map.
func (pt *ProbabilityTable[T]) AddAll(itemsAndWeights map[T]int) *ProbabilityTable[T] {
	for item, weight := range itemsAndWeights {
		pt.Add(item, weight)
	}
	return pt
}

// Remove removes or reduces the weight of an item by the given amount.
func (pt *ProbabilityTable[T]) Remove(item T, weight int) bool {
	if weight <= 0 {
		return false
	}

	idx := slices.Index(pt.Table, item)
	if idx < 0 {
		return false
	}

	oldWeight := pt.Weights[idx]
	pt.Weights[idx] -= weight
	newWeight := pt.Weights[idx]

	if newWeight <= 0 {
		pt.Table = slices.Delete(pt.Table, idx, idx+1)
		pt.Weights = slices.Delete(pt.Weights, idx, idx+1)
	}

	removed := min(oldWeight, oldWeight-max(0, newWeight))
	pt.Total -= removed
	return true
}

// RemoveItem completely removes an item from the top-level table.
func (pt *ProbabilityTable[T]) RemoveItem(item T) bool {
	return pt.Remove(item, pt.Weight(item))
}

// RemoveAll removes all specified items from the top-level table.
func (pt *ProbabilityTable[T]) RemoveAll(items []T) bool {
	changed := false
	for _, item := range items {
		changed = pt.RemoveItem(item) || changed
	}
	return changed
}

// AddNested adds a nested ProbabilityTable with a given weight.
func (pt *ProbabilityTable[T]) AddNested(table *ProbabilityTable[T], weight int) *ProbabilityTable[T] {
	if weight <= 0 || table == nil || pt.ContentEquals(table) || table.Total <= 0 {
		return pt
	}
	pt.Weights = append(pt.Weights, weight)
	pt.ExtraTable = append(pt.ExtraTable, table)
	pt.Total += weight
	return pt
}

// Weight returns the top-level weight of a specific item.
func (pt *ProbabilityTable[T]) Weight(item T) int {
	idx := slices.Index(pt.Table, item)
	if idx < 0 {
		return 0
	}
	return pt.Weights[idx]
}

// WeightNested returns the weight of a nested table.
func (pt *ProbabilityTable[T]) WeightNested(table *ProbabilityTable[T]) int {
	idx := -1
	for i, ext := range pt.ExtraTable {
		if ext == table {
			idx = i
			break
		}
	}
	if idx < 0 {
		return 0
	}
	return pt.Weights[len(pt.Table)+idx]
}

// Items returns a slice containing all unique items (including items inside nested tables).
func (pt *ProbabilityTable[T]) Items() []T {
	itemSet := make(map[T]struct{}, len(pt.Table))
	var result []T

	var collect func(p *ProbabilityTable[T])
	collect = func(p *ProbabilityTable[T]) {
		for _, item := range p.Table {
			if _, exists := itemSet[item]; !exists {
				itemSet[item] = struct{}{}
				result = append(result, item)
			}
		}
		for _, extra := range p.ExtraTable {
			collect(extra)
		}
	}

	collect(pt)
	return result
}

// SimpleItems returns a direct slice of items at the top-level table.
func (pt *ProbabilityTable[T]) SimpleItems() []T {
	return pt.Table
}

// Tables returns all nested tables.
func (pt *ProbabilityTable[T]) Tables() []*ProbabilityTable[T] {
	return pt.ExtraTable
}

// ContentEquals compares table structure and weights without checking RNG equality.
func (pt *ProbabilityTable[T]) ContentEquals(other *ProbabilityTable[T]) bool {
	if pt == other {
		return true
	}
	if pt == nil || other == nil {
		return false
	}
	if !slices.Equal(pt.Table, other.Table) || !slices.Equal(pt.Weights, other.Weights) {
		return false
	}
	if len(pt.ExtraTable) != len(other.ExtraTable) {
		return false
	}
	for i := range pt.ExtraTable {
		if !pt.ExtraTable[i].ContentEquals(other.ExtraTable[i]) {
			return false
		}
	}
	return true
}

// Equals checks for structural and content equality.
func (pt *ProbabilityTable[T]) Equals(other *ProbabilityTable[T]) bool {
	return pt.ContentEquals(other)
}

// HashCode returns a combined hash of the table, extra tables, and weights.
func (pt *ProbabilityTable[T]) HashCode() int {
	result := ptHashSlice(pt.Table)
	result = 421*result + ptHashNestedSlice(pt.ExtraTable)
	result = 83*result + ptHashSlice(pt.Weights)
	return result
}

func ptHashSlice[S ~[]E, E any](s S) int {
	h := 1
	for _, v := range s {
		h = 31*h + int(reflect.ValueOf(v).Uint())
	}
	return h
}

func ptHashNestedSlice[T comparable](tables []*ProbabilityTable[T]) int {
	h := 1
	for _, t := range tables {
		h = 31*h + t.HashCode()
	}
	return h
}
