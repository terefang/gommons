package xarray

// Remove removes the first matching value from a slice.
func Remove[T comparable](slice []T, value T) []T {
	for i, v := range slice {
		if v == value {
			// Append elements before index i with elements after index i
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

// RemoveAll removes every instance of a value from a slice in-place.
func RemoveAll[T comparable](slice []T, value T) []T {
	out := slice[:0] // Reuse underlying array memory
	for _, v := range slice {
		if v != value {
			out = append(out, v)
		}
	}
	return out
}

// RemoveUnordered removes the first instance of value in-place without preserving order.
func RemoveUnordered[T comparable](slice []T, value T) []T {
	for i, v := range slice {
		if v == value {
			lastIdx := len(slice) - 1
			slice[i] = slice[lastIdx]
			return slice[:lastIdx]
		}
	}
	return slice
}

// RemoveAllUnordered removes all instances of value in-place without preserving order.
func RemoveAllUnordered[T comparable](slice []T, value T) []T {
	length := len(slice)
	for i := 0; i < length; {
		if slice[i] == value {
			length--
			slice[i] = slice[length]
			continue
		}
		i++
	}
	return slice[:length]
}
