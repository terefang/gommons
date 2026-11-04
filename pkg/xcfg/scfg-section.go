package xcfg

import "strings"

// Add sets or updates a property key-value pair within the section.
func (cs *ScfgSection) Add(property string, value string) {
	property = strings.TrimSpace(strings.ToLower(property))
	cs.Properties[property] = value
}

// Has checks if a given property key exists within the section.
func (cs *ScfgSection) Has(property string) bool {
	property = strings.TrimSpace(strings.ToLower(property))
	_, ok := cs.Properties[property]
	return ok
}

// Get retrieves a property value by key from the section.
// Returns the string value and a boolean indicating whether the key was present.
func (cs *ScfgSection) Get(property string) (string, bool) {
	property = strings.TrimSpace(strings.ToLower(property))
	value, ok := cs.Properties[property]
	return value, ok
}

// GetDefault retrieves a property value by key from the section, returning _default if key is missing.
func (cs *ScfgSection) GetDefault(property, _default string) string {
	property = strings.TrimSpace(strings.ToLower(property))
	value, ok := cs.Properties[property]
	if ok {
		return value
	}
	return _default
}
