package xmath

import "math/bits"

// LowestOneBit32 returns an int32 value with at most a single one-bit,
// in the position of the lowest-order ("rightmost") one-bit in the specified value.
func LowestOneBit32(num int32) int32 {
	return num & -num
}

// LowestOneBit64 returns an int64 value with at most a single one-bit,
// in the position of the lowest-order ("rightmost") one-bit in the specified value.
func LowestOneBit64(num int64) int64 {
	return num & ^(num - 1) // is equivalent to num & -num
}

// IMul performs C-style/JavaScript 32-bit signed integer multiplication.
// In Go, multiplying two int32 values natively wraps on overflow.
func IMul(left, right int32) int32 {
	return left * right
}

// CountLeadingZeros32 returns the number of zero bits preceding the highest-order
// ("leftmost") one-bit in the 32-bit integer.
func CountLeadingZeros32(n int32) int {
	return bits.LeadingZeros32(uint32(n))
}

// CountTrailingZeros32 returns the number of zero bits following the lowest-order
// ("rightmost") one-bit in the 32-bit integer.
func CountTrailingZeros32(n int32) int {
	return bits.TrailingZeros32(uint32(n))
}

// CountLeadingZeros64 returns the number of zero bits preceding the highest-order
// ("leftmost") one-bit in the 64-bit integer.
func CountLeadingZeros64(n int64) int {
	return bits.LeadingZeros64(uint64(n))
}

// CountTrailingZeros64 returns the number of zero bits following the lowest-order
// ("rightmost") one-bit in the 64-bit integer.
func CountTrailingZeros64(n int64) int {
	return bits.TrailingZeros64(uint64(n))
}
