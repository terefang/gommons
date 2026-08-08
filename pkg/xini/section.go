package xini

// IniSection provides access to the properties of a single section in an IniConfig.
//
// An IniSection is obtained by calling Section on an IniConfig. All property
// operations performed through an IniSection are scoped to the section
// represented by the IniSection.
type IniSection struct {
	name string
	ic   *IniConfig
}

// Name returns the name of the section.
func (is *IniSection) Name() string {
	return is.name
}

// PropertyExists reports whether the specified property exists in the section.
func (is *IniSection) PropertyExists(propertyName string) bool {
	return is.ic.PropertyExists(is.name, propertyName)
}

// Value returns the value of the specified property.
//
// If the property does not exist or cannot be read, Value returns an error.
func (is *IniSection) Get(propertyName string) (string, error) {
	return is.ic.Get(is.name, propertyName)
}

// OrZero returns the value of the specified property.
//
// If the property does not exist or cannot be read, returns the
// zero value for a string.
func (is *IniSection) OrZero(propertyName string) string {
	return is.ic.OrZero(is.name, propertyName)
}

// AsFloat64 returns the value of the specified property as a float64.
//
// If the property does not exist or cannot be converted to a float64, AsFloat64
// returns an error.
func (is *IniSection) AsFloat64(propertyName string) (float64, error) {
	return is.ic.AsFloat64(is.name, propertyName)
}

// AsFloat64OrZero returns the value of the specified property as a float64.
//
// If the property does not exist or cannot be converted to a float64,
// AsFloat64OrZero returns the zero value for float64.
func (is *IniSection) AsFloat64OrZero(propertyName string) float64 {
	return is.ic.AsFloat64OrZero(is.name, propertyName)
}

// AsInt64 returns the value of the specified property as an int64.
//
// If the property does not exist or cannot be converted to an int64, AsInt64
// returns an error.
func (is *IniSection) AsInt64(propertyName string) (int64, error) {
	return is.ic.AsInt64(is.name, propertyName)
}

// AsInt64OrZero returns the value of the specified property as an int64.
//
// If the property does not exist or cannot be converted to an int64,
// AsInt64OrZero returns the zero value for int64.
func (is *IniSection) AsInt64OrZero(propertyName string) int64 {
	return is.ic.AsInt64OrZero(is.name, propertyName)
}

// AsUint64 returns the value of the specified property as a uint64.
//
// If the property does not exist or cannot be converted to a uint64, AsUint64
// returns an error.
func (is *IniSection) AsUint64(propertyName string) (uint64, error) {
	return is.ic.AsUint64(is.name, propertyName)
}

// AsUint64OrZero returns the value of the specified property as a uint64.
//
// If the property does not exist or cannot be converted to a uint64,
// AsUint64OrZero returns the zero value for uint64.
func (is *IniSection) AsUint64OrZero(propertyName string) uint64 {
	return is.ic.AsUint64OrZero(is.name, propertyName)
}

// AsBool returns the value of the specified property as a bool.
//
// If the property does not exist or cannot be converted to a bool, AsBool
// returns an error.
func (is *IniSection) AsBool(propertyName string) (bool, error) {
	return is.ic.AsBool(is.name, propertyName)
}

// AsBoolOrZero returns the value of the specified property as a bool.
//
// If the property does not exist or cannot be converted to a bool,
// AsBoolOrZero returns the zero value for bool.
func (is *IniSection) AsBoolOrZero(propertyName string) bool {
	return is.ic.AsBoolOrZero(is.name, propertyName)
}

// Add adds or updates a property in the section.
//
// The property is stored under the section represented by the IniSection.
func (is *IniSection) Add(propertyName string, value string) {
	is.ic.Add(is.name, propertyName, value)
}
