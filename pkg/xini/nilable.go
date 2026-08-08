package xini

// newNilableString creates a nilableString containing the supplied value.
//
// The returned nilableString is marked as set, including when v is the empty
// string.
func newNilableString(v string) *nilableString {
	ns := new(nilableString)
	ns.Set(v)

	return ns
}

// nilableString stores a string value together with an indication of whether
// the value has been explicitly set.
//
// This allows an explicitly set empty string to be distinguished from the
// zero value of nilableString.
type nilableString struct {
	val string
	set bool
}

// Set stores v as the contained value and marks the nilableString as set.
//
// An empty string is considered an explicitly set value.
func (ns *nilableString) Set(v string) {
	ns.val = v
	ns.set = true
}

// String returns the currently stored value.
//
// String returns the stored value regardless of whether it has been explicitly
// set. For an unset nilableString, String returns the empty string.
func (ns *nilableString) String() string {
	return ns.val
}

// IsSet reports whether a value has been explicitly set.
//
// IsSet returns true even when the explicitly set value is the empty string.
func (ns *nilableString) IsSet() bool {
	return ns.set
}
