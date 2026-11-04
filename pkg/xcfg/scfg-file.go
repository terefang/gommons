package xcfg

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// NewSection creates a new section under the given name and attaches it to the configuration.
func (c *Scfg) NewSection(name string) *ScfgSection {
	name = strings.TrimSpace(strings.ToLower(name))
	s := NewScfgSection()
	s.Name = name
	c.Sections[name] = s
	return s
}

// HasSection checks whether a section with the given name exists.
func (c *Scfg) HasSection(name string) bool {
	name = strings.TrimSpace(strings.ToLower(name))
	_, ok := c.Sections[name]
	return ok
}

// GetSection returns a section pointer by name and a boolean indicating if it exists.
func (c *Scfg) GetSection(name string) (*ScfgSection, bool) {
	name = strings.TrimSpace(strings.ToLower(name))
	s, ok := c.Sections[name]
	return s, ok
}

func (c *Scfg) SectionNames() []string {
	_ret := make([]string, 0)
	for k, _ := range c.Sections {
		_ret = append(_ret, k)
	}
	return _ret
}

// RemoveSection deletes a specific section by name from the configuration.
func (c *Scfg) RemoveSection(name string) {
	name = strings.TrimSpace(strings.ToLower(name))
	delete(c.Sections, name)
}

// RemoveAllSections clears all sections from the configuration.
func (c *Scfg) RemoveAllSections() {
	for name := range c.Sections {
		delete(c.Sections, name)
	}
}

func (c *Scfg) Keys(name string) []string {
	_ret := make([]string, 0)
	if !c.HasSection(name) {
		return _ret
	}
	s, _ := c.GetSection(name)
	for k, _ := range s.Properties {
		_ret = append(_ret, k)
	}
	return _ret
}

// Has checks whether a property exists within a specific section.
func (c *Scfg) Has(name, prop string) bool {
	if !c.HasSection(name) {
		return false
	}
	s, _ := c.GetSection(name)
	return s.Has(prop)
}

// Get returns the raw string value of a property in a given section, along with an existence boolean.
func (c *Scfg) Get(name, prop string) (string, bool) {
	if !c.HasSection(name) {
		return "", false
	}
	s, _ := c.GetSection(name)
	return s.Get(prop)
}

// AsStringOrEmpty returns the string value of a property, or "" if not found.
func (c *Scfg) AsStringOrEmpty(name, prop string) string {
	return c.AsStringOrDefault(name, prop, "")
}

// AsStringOrDefault returns the string value of a property, or _default if not found.
func (c *Scfg) AsStringOrDefault(name, prop, _default string) string {
	if !c.Has(name, prop) {
		return _default
	}
	_val, _ := c.Get(name, prop)
	return _val
}

// AsIntOrZero returns the int64 value of a property, or 0 if missing/invalid.
func (c *Scfg) AsIntOrZero(name, prop string) int64 {
	return c.AsIntOrDefault(name, prop, 0)
}

// AsIntOrDefault parses a property to an int64. Returns _default if missing or unparseable.
func (c *Scfg) AsIntOrDefault(name, prop string, _default int64) int64 {
	if !c.Has(name, prop) {
		return _default
	}
	_val, _ok := c.Get(name, prop)
	if !_ok {
		return _default
	}
	_int, _err := strconv.ParseInt(_val, 10, 64)
	if _err != nil {
		return _default
	}
	return _int
}

// AsUintOrZero returns the uint64 value of a property, or 0 if missing/invalid.
func (c *Scfg) AsUintOrZero(name, prop string) uint64 {
	return c.AsUintOrDefault(name, prop, 0)
}

// AsUintOrDefault parses a property to a uint64. Returns _default if missing or unparseable.
func (c *Scfg) AsUintOrDefault(name, prop string, _default uint64) uint64 {
	if !c.Has(name, prop) {
		return _default
	}
	_val, _ok := c.Get(name, prop)
	if !_ok {
		return _default
	}
	_int, _err := strconv.ParseUint(_val, 10, 64)
	if _err != nil {
		return _default
	}
	return _int
}

// AsBoolOrFalse parses a property as boolean, returning false as fallback.
func (c *Scfg) AsBoolOrFalse(name, prop string) bool {
	return c.AsBoolOrDefault(name, prop, false)
}

// AsBoolOrTrue parses a property as boolean, returning true as fallback.
func (c *Scfg) AsBoolOrTrue(name, prop string) bool {
	return c.AsBoolOrDefault(name, prop, true)
}

// AsBoolOrDefault parses a boolean value.
// Case-insensitive true strings: "true", "t", "yes", "y", "on", "1".
// Case-insensitive false strings: "false", "f", "no", "n", "off", "0".
// Returns _default if missing or non-matching.
func (c *Scfg) AsBoolOrDefault(name, prop string, _default bool) bool {
	if !c.Has(name, prop) {
		return _default
	}
	_val, _ := c.Get(name, prop)
	_val = strings.ToLower(_val)
	switch _val {
	case "true", "t", "yes", "y", "on", "1":
		return true
	case "false", "f", "no", "n", "off", "0":
		return false
	}
	return _default
}

// Add sets a property under a given section. Creates the section automatically if missing.
func (c *Scfg) Add(name, prop, value string) {
	if !c.HasSection(name) {
		c.NewSection(name)
	}
	s, _ := c.GetSection(name)
	s.Add(prop, value)
}

// ParseFile opens a file at path and parses its contents into the Scfg configuration.
func (c *Scfg) ParseFile(path string) error {
	if f, err := os.Open(path); err != nil {
		return err
	} else {
		defer f.Close()
		return c.Parse(f)
	}
}

// Parse reads lines from an open file handle, parsing sections, properties, line continuations, and comments.
func (c *Scfg) Parse(content *os.File) error {
	s := bufio.NewScanner(content)

	section := SCFG_GLOBAL_SECTION

	lineNumber := 0

	for s.Scan() {

		lineNumber++

		l := strings.TrimSpace(s.Text())
		lineLength := len(l)

		if lineLength == 0 {
			continue
		}

		commentFound := false
		for _, c := range []rune(SCFG_COMMENT_START) {
			if rune(l[0]) == c {
				commentFound = true
				break
			}
		}
		if commentFound {
			continue
		}

		// read line continuation
		for l[lineLength-1] == '\\' {
			l = l[:lineLength-1]
			ok := s.Scan()
			if ok {
				lineNumber++
				l += strings.TrimSpace(s.Text())
				lineLength = len(l)
			} else {
				break
			}
		}

		if l[0] == '[' && l[lineLength-1] == ']' {
			section = strings.TrimSpace(l[1 : lineLength-1])
			c.NewSection(section)
			continue
		}

		parts := strings.SplitN(l, "=", 2)
		plen := len(parts)
		parts[0] = strings.TrimSpace(parts[0])
		if plen == 1 {
			// single key is boolean
			c.Add(section, parts[0], SCFG_BOOL_TRUE)
		} else {
			parts[1] = strings.TrimSpace(parts[1])
			parts[1] = stripQuotes(parts[1])
			c.Add(section, parts[0], parts[1])
		}
	}
	return nil
}

func stripQuotes(value string) string {
	vLen := len(value)

	if vLen < 2 {
		return value
	}

	for _, r := range []rune(SCFG_QUOTES) {
		if rune(value[0]) == r && rune(value[vLen-1]) == r {
			return value[1 : vLen-1]
		}
	}

	return value
}

func (c *Scfg) Unmarshal(sectionName string, holder interface{}) error {
	sectionName = strings.TrimSpace(strings.ToLower(sectionName))
	_type := reflect.TypeOf(holder)

	if _type.Kind() != reflect.Ptr {
		return fmt.Errorf("argument must be a pointer, not %s", _type.Kind())
	}

	_type = _type.Elem()

	if _type.Kind() != reflect.Struct {
		return fmt.Errorf("argument must be a struct, not %s", _type.Kind())
	}

	_value := reflect.ValueOf(holder).Elem()

	for i := 0; i < _type.NumField(); i++ {
		_f := _type.Field(i)

		_vmap, ok := _f.Tag.Lookup("ini")
		if !ok {
			_vmap = _f.Name
		}

		_vmaps := strings.Split(_vmap, ",")
		_vf := _value.Field(i)
		for _, _vm := range _vmaps {
			if c.Has(sectionName, _vm) {
				switch _vf.Kind() {
				case reflect.Bool:
					_vf.SetBool(c.AsBoolOrFalse(sectionName, _vm))
				case reflect.Float32:
					_v, _ := strconv.ParseFloat(strings.TrimSpace(c.AsStringOrEmpty(sectionName, _vm)), 32)
					_vf.SetFloat(_v)
				case reflect.Float64:
					_v, _ := strconv.ParseFloat(strings.TrimSpace(c.AsStringOrEmpty(sectionName, _vm)), 64)
					_vf.SetFloat(_v)
				case reflect.Int, reflect.Int64:
					_vf.SetInt(c.AsIntOrZero(sectionName, _vm))
				case reflect.Uint, reflect.Uint64:
					_vf.SetUint(c.AsUintOrZero(sectionName, _vm))
				case reflect.String:
					_vf.SetString(c.AsStringOrEmpty(sectionName, _vm))
				case reflect.Slice:
					_slice := strings.Split(c.AsStringOrEmpty(sectionName, _vm), ",")
					_sl := len(_slice)
					switch _f.Type.String() {
					case "[]string":
						_vf.Set(reflect.ValueOf(_slice))
					case "[]int":
						_sv := make([]int, _sl)
						for j := 0; j < _sl; j++ {
							_sv[j], _ = strconv.Atoi(strings.TrimSpace(_slice[j]))
						}
						_vf.Set(reflect.ValueOf(_sv))
					case "[]int8":
						_sv := make([]int8, _sl)
						for j := 0; j < _sl; j++ {
							_v, _ := strconv.ParseInt(strings.TrimSpace(_slice[j]), 10, 8)
							_sv[j] = int8(_v)
						}
						_vf.Set(reflect.ValueOf(_sv))
					case "[]int16":
						_sv := make([]int16, _sl)
						for j := 0; j < _sl; j++ {
							_v, _ := strconv.ParseInt(strings.TrimSpace(_slice[j]), 10, 16)
							_sv[j] = int16(_v)
						}
						_vf.Set(reflect.ValueOf(_sv))
					case "[]int32":
						_sv := make([]int32, _sl)
						for j := 0; j < _sl; j++ {
							_v, _ := strconv.ParseInt(strings.TrimSpace(_slice[j]), 10, 32)
							_sv[j] = int32(_v)
						}
						_vf.Set(reflect.ValueOf(_sv))
					case "[]int64":
						_sv := make([]int64, _sl)
						for j := 0; j < _sl; j++ {
							_sv[j], _ = strconv.ParseInt(strings.TrimSpace(_slice[j]), 10, 64)
						}
						_vf.Set(reflect.ValueOf(_sv))
					case "[]uint":
						_sv := make([]uint, _sl)
						for j := 0; j < _sl; j++ {
							_v, _ := strconv.ParseUint(strings.TrimSpace(_slice[j]), 10, 32)
							_sv[j] = uint(_v)
						}
						_vf.Set(reflect.ValueOf(_sv))
					case "[]uint8":
						_sv := make([]uint8, _sl)
						for j := 0; j < _sl; j++ {
							_v, _ := strconv.ParseUint(strings.TrimSpace(_slice[j]), 10, 8)
							_sv[j] = uint8(_v)
						}
						_vf.Set(reflect.ValueOf(_sv))
					case "[]uint16":
						_sv := make([]uint16, _sl)
						for j := 0; j < _sl; j++ {
							_v, _ := strconv.ParseUint(strings.TrimSpace(_slice[j]), 10, 16)
							_sv[j] = uint16(_v)
						}
						_vf.Set(reflect.ValueOf(_sv))
					case "[]uint32":
						_sv := make([]uint32, _sl)
						for j := 0; j < _sl; j++ {
							_v, _ := strconv.ParseUint(strings.TrimSpace(_slice[j]), 10, 32)
							_sv[j] = uint32(_v)
						}
						_vf.Set(reflect.ValueOf(_sv))
					case "[]uint64":
						_sv := make([]uint64, _sl)
						for j := 0; j < _sl; j++ {
							_sv[j], _ = strconv.ParseUint(strings.TrimSpace(_slice[j]), 10, 64)
						}
						_vf.Set(reflect.ValueOf(_sv))
					case "[]float32":
						_sv := make([]float32, _sl)
						for j := 0; j < _sl; j++ {
							_v, _ := strconv.ParseFloat(strings.TrimSpace(_slice[j]), 32)
							_sv[j] = float32(_v)
						}
						_vf.Set(reflect.ValueOf(_sv))
					case "[]float64":
						_sv := make([]float64, _sl)
						for j := 0; j < _sl; j++ {
							_sv[j], _ = strconv.ParseFloat(strings.TrimSpace(_slice[j]), 64)
						}
						_vf.Set(reflect.ValueOf(_sv))
					}
				default: // IGNORE
					return fmt.Errorf("unknown type: %s on field '%s'", _vf.Kind(), _f.Name)
				}
				break
			}
		}
	}
	return nil
}
