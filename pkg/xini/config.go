package xini

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func New() *IniConfig {
	ic := new(IniConfig)
	ic.options = DefaultIniOptions()
	ic.sections = make(sectionPropertyMap)

	return ic
}

func NewWithOptions(options *IniOptions) *IniConfig {
	ic := new(IniConfig)
	ic.options = options
	ic.sections = make(sectionPropertyMap)

	return ic
}

// NewIniConfig loads the INI file at path into a new IniConfig.
//
// NewIniConfig uses the options returned by DefaultIniOptions.
//
// NewIniConfig returns an error if the file cannot be opened or if
// the file cannot be parsed as an INI file.
func NewIniConfig(path string) (*IniConfig, error) {
	return NewIniConfigWithOptions(path, DefaultIniOptions())
}

// NewIniConfigWithOptions loads the INI file at path into a new IniConfig
// using the supplied parsing options.
//
// NewIniConfigWithOptions returns an error if the file cannot be opened
// or if the file cannot be parsed as an INI file.
func NewIniConfigWithOptions(path string, options *IniOptions) (*IniConfig, error) {
	if f, err := os.Open(path); err != nil {
		return nil, err
	} else {
		defer f.Close()

		return NewIniConfigFromFileWithOptions(f, options)
	}
}

// NewIniConfigFromFile loads an INI file from file into a new IniConfig.
//
// NewIniConfigFromFile uses the options returned by DefaultIniOptions.
// The caller is responsible for closing file.
//
// NewIniConfigFromFile returns an error if the supplied file cannot be read
// or if its contents cannot be parsed as an INI file.
func NewIniConfigFromFile(file *os.File) (*IniConfig, error) {
	return NewIniConfigFromFileWithOptions(file, DefaultIniOptions())
}

// NewIniConfigFromFileWithOptions loads an INI file from file into a new
// IniConfig using the supplied parsing options.
//
// The caller is responsible for closing file.
//
// NewIniConfigFromFileWithOptions returns an error if file or options is nil,
// if options.CommentStart is empty, if the supplied file cannot be read, or
// if its contents cannot be parsed as an INI file.
func NewIniConfigFromFileWithOptions(file *os.File, options *IniOptions) (*IniConfig, error) {
	if file == nil {
		return nil, errors.New("Nil file provided")
	}

	if options == nil {
		return nil, errors.New("Nil IniOptions provided")
	}

	ic := NewWithOptions(options)

	if err := ic.ParseFromFile(file); err != nil {
		return nil, err
	}

	return ic, nil
}

const GLOBAL_SECTION = ""

// DefaultIniOptions returns an IniOptions object populated with default values useful for working with most INI files.
//
// Default values are:
//
//	CommentStart 					[]rune{';', '#', '*', '!'}
//	StrictBoolTrue					"TRUE"
//	StrictBoolFalse					"FALSE"
//	EnclosingQuoteSymbols			[]rune{'\'','"'}
func DefaultIniOptions() *IniOptions {
	io := new(IniOptions)

	io.CommentStart = []rune{';', '#', '*', '!'}
	io.StrictBoolTrue = "TRUE"
	io.StrictBoolFalse = "FALSE"
	io.EnclosingQuoteSymbols = []rune{'\'', '"', '`'}

	return io
}

// IniOptions controls how an INI file is parsed and how parsed configuration
// values are subsequently interpreted.
//
// IniOptions is supplied when creating an IniConfig and allows applications to
// customize section and property name handling, comments, whitespace, blank
// lines, boolean conversion, quoting, and property assignment syntax.
type IniOptions struct {

	// CommentStart specifies the rune(s) that, when found at the start of a
	// line, identifies the line as a comment.
	CommentStart []rune

	// StrictBoolTrue specifies the value that represents true when
	// UseGoBoolRules is false.
	StrictBoolTrue string

	// StrictBoolFalse specifies the value that represents false when
	// UseGoBoolRules is false.
	StrictBoolFalse string

	// EnclosingQuoteSymbols specifies the quote characters recognized as
	// enclosing quotes when StripEnclosingQuotes is enabled.
	EnclosingQuoteSymbols []rune
}

type sectionPropertyMap map[string]map[string]*nilableString

// IniConfig provides access to configuration loaded from an INI file.
//
// IniConfig provides methods for checking whether sections and properties exist,
// retrieving raw property values, and converting property values to common Go
// types.
//
// The As... methods are convenience functions around the corresponding
// strconv.Parse... functions. The OrZero... methods suppress lookup and
// conversion errors and return the zero value of the requested type instead.
type IniConfig struct {
	sections sectionPropertyMap
	options  *IniOptions
}

// SectionExists reports whether a section with the specified name exists.
func (ic *IniConfig) SectionExists(sectionName string) bool {
	return ic.findSection(sectionName) != nil
}

// Section returns a view of the specified section.
//
// The returned IniSection provides the same property access and conversion
// methods as IniConfig, but all operations are scoped to the specified section.
//
// Section returns an error if the specified section does not exist.
func (ic *IniConfig) Section(sectionName string) (*IniSection, error) {
	if ic.SectionExists(sectionName) {
		is := new(IniSection)
		is.name = sectionName
		is.ic = ic

		return is, nil
	}

	return nil, errorf("Section %s does not exist", sectionName)
}

// PropertyExists reports whether the specified property exists in the specified section.
//
// PropertyExists returns false if the section does not exist.
func (ic *IniConfig) PropertyExists(sectionName, propertyName string) bool {
	propertyName = ic.normalise(propertyName)

	if foundSection := ic.findSection(sectionName); foundSection == nil {
		return false
	} else {
		return foundSection[propertyName] != nil
	}
}

// Get returns the raw string value of the specified property in the specified section.
//
// returns an error if the section or property does not exist.
func (ic *IniConfig) Get(sectionName, propertyName string) (string, error) {
	section := ic.findSection(sectionName)
	propertyName = ic.normalise(propertyName)

	if section == nil {
		return "", errorf("No such section %s", sectionName)
	}

	value, ok := section[propertyName]
	if !ok {
		return "", errorf("No such property [%s].%s", sectionName, propertyName)
	}

	if value != nil {
		return value.String(), nil
	}
	return "", errorf("No such property [%s].%s", sectionName, propertyName)
}

// OrZero returns the raw string value of the specified property.
//
// If the section or property does not exist, ValueOrZero returns the empty
// string.
func (ic *IniConfig) OrZero(sectionName, propertyName string) string {
	if v, err := ic.Get(sectionName, propertyName); err == nil {
		return v
	}

	return ""
}

// AsFloat64 returns the value of the specified property converted to a float64.
//
// AsFloat64 uses strconv.ParseFloat with a bit size of 64. It returns an
// error if the section or property does not exist or if the property value
// cannot be converted to a float64.
func (ic *IniConfig) AsFloat64(sectionName, propertyName string) (float64, error) {
	origSectionName := sectionName
	origPropName := propertyName

	sv, err := ic.Get(sectionName, propertyName)
	if err != nil {
		return 0, err
	}

	if v, err := strconv.ParseFloat(sv, 64); err == nil {
		return v, nil
	}

	return 0, errorf(
		"Unable to interpret [%s].%s (%s) as a float64.",
		origSectionName, origPropName, sv,
	)
}

// AsFloat64OrZero returns the value of the specified property as a float64.
//
// If the section or property does not exist, or if the property value cannot
// be converted to a float64, GetOrZeroAsFloat64 returns 0.
func (ic *IniConfig) AsFloat64OrZero(sectionName, propertyName string) float64 {
	if v, err := ic.AsFloat64(sectionName, propertyName); err == nil {
		return v
	}

	return 0
}

// AsInt64 returns the value of the specified property converted to an int64.
//
// AsInt64 uses strconv.ParseInt with base 10 and a bit size of 64. It
// returns an error if the section or property does not exist or if the property
// value cannot be converted to an int64.
func (ic *IniConfig) AsInt64(sectionName, propertyName string) (int64, error) {
	origSectionName := sectionName
	origPropName := propertyName

	sv, err := ic.Get(sectionName, propertyName)
	if err != nil {
		return 0, err
	}

	if v, err := strconv.ParseInt(sv, 10, 64); err == nil {
		return v, nil
	}

	return 0, errorf(
		"Unable to interpret [%s].%s (%s) as an int64.",
		origSectionName, origPropName, sv,
	)
}

// AsInt64OrZero returns the value of the specified property as an int64.
//
// If the section or property does not exist, or if the property value cannot
// be converted to an int64, GetOrZeroAsInt64 returns 0.
func (ic *IniConfig) AsInt64OrZero(sectionName, propertyName string) int64 {
	if v, err := ic.AsInt64(sectionName, propertyName); err == nil {
		return v
	}

	return 0
}

// AsUint64 returns the value of the specified property converted to a uint64.
//
// AsUint64 uses strconv.ParseUint with base 10 and a bit size of 64. It
// returns an error if the section or property does not exist or if the property
// value cannot be converted to a uint64.
func (ic *IniConfig) AsUint64(sectionName, propertyName string) (uint64, error) {
	origSectionName := sectionName
	origPropName := propertyName

	sv, err := ic.Get(sectionName, propertyName)
	if err != nil {
		return 0, err
	}

	if v, err := strconv.ParseUint(sv, 10, 64); err == nil {
		return v, nil
	}

	return 0, errorf(
		"Unable to interpret [%s].%s (%s) as a uint64.",
		origSectionName, origPropName, sv,
	)
}

// AsUint64OrZero returns the value of the specified property as a uint64.
//
// If the section or property does not exist, or if the property value cannot
// be converted to a uint64, returns 0.
func (ic *IniConfig) AsUint64OrZero(sectionName, propertyName string) uint64 {
	if v, err := ic.AsUint64(sectionName, propertyName); err == nil {
		return v
	}

	return 0
}

// AsBool returns the value of the specified property converted to a bool.
//
// Boolean conversion is controlled by the IniOptions supplied when creating
// the IniConfig.
//
// If UseGoBoolRules is true, AsBool uses strconv.ParseBool and accepts
// the boolean representations supported by Go.
//
// If UseGoBoolRules is false, the property value must match StrictBoolTrue
// or StrictBoolFalse. Matching is case-insensitive when StrictBoolCaseSensitive
// is false.
//
// AsBool returns an error if the section or property does not exist or
// if the property value cannot be interpreted as a boolean according to the
// configured rules.
func (ic *IniConfig) AsBool(sectionName, propertyName string) (bool, error) {
	sv, err := ic.Get(sectionName, propertyName)

	origSv := sv
	options := ic.options

	if err != nil {
		return false, err
	}

	if bv, err := strconv.ParseBool(sv); err == nil {
		return bv, nil
	}

	strictTrue := options.StrictBoolTrue
	strictFalse := options.StrictBoolFalse

	strictTrue = strings.ToUpper(strictTrue)
	strictFalse = strings.ToUpper(strictFalse)
	sv = strings.ToUpper(sv)

	if sv == strictTrue {
		return true, nil
	}

	if sv == strictFalse {
		return false, nil
	}

	return false, errorf(
		"Value of [%s].%s (%s) could not be matched to %s or %s",
		sectionName, propertyName, origSv,
		options.StrictBoolTrue, options.StrictBoolFalse,
	)
}

// AsBoolOrZero returns the value of the specified property as a bool.
//
// If the section or property does not exist, or if the property value cannot
// be interpreted as a bool, returns false.
func (ic *IniConfig) AsBoolOrZero(sectionName, propertyName string) bool {
	if v, err := ic.AsBool(sectionName, propertyName); err == nil {
		return v
	}

	return false
}

// Add stores a property in the specified section.
//
// If the section does not exist, it is created. If the property already exists,
// its value is replaced with the supplied value.
func (ic *IniConfig) Add(section, propertyName string, value string) {
	section = ic.normalise(section)
	propertyName = ic.normalise(propertyName)

	storedSection := ic.sections[section]

	if storedSection == nil {
		storedSection = make(map[string]*nilableString)
		ic.sections[section] = storedSection
	}

	storedSection[propertyName] = newNilableString(value)
}

func (ic *IniConfig) stripQuotes(value string) string {

	options := ic.options

	vLength := len(value)

	if vLength < 2 {
		return value
	}

	for _, r := range options.EnclosingQuoteSymbols {

		if rune(value[0]) == r && rune(value[vLength-1]) == r {

			stripped := value[1 : vLength-1]

			return stripped
		}

	}

	return value
}

func (ic *IniConfig) findSection(sectionName string) map[string]*nilableString {
	sectionName = ic.normalise(sectionName)

	return ic.sections[sectionName]
}

func (ic *IniConfig) normalise(s string) string {
	return strings.ToLower(s)
}

func errorf(template string, args ...interface{}) error {
	m := fmt.Sprintf(template, args...)

	return errors.New(m)
}

func (ic *IniConfig) ParseFromFile(cf *os.File) error {
	s := bufio.NewScanner(cf)
	return ic.Parse(s)
}

// Parse scans the supplied scanner line by line according to the rules defined in the IniOptions
func (ic *IniConfig) Parse(s *bufio.Scanner) error {

	section := GLOBAL_SECTION

	lineNumber := 0

	for s.Scan() {

		lineNumber++

		l := strings.TrimSpace(s.Text())
		lineLength := len(l)

		if lineLength == 0 {
			continue
		}

		for _, c := range ic.options.CommentStart {
			if rune(l[0]) == c {
				continue
			}
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
			section = l[1 : lineLength-1]
			continue
		}

		parts := strings.SplitN(l, "=", 2)
		plen := len(parts)
		parts[0] = strings.TrimSpace(parts[0])
		if plen == 1 {
			// single key is boolean
			ic.Add(section, parts[0], "TRUE")
		} else {
			parts[1] = strings.TrimSpace(parts[1])
			parts[1] = ic.stripQuotes(parts[1])
			ic.Add(section, parts[0], parts[1])
		}
	}

	return nil
}
