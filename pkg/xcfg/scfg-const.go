package xcfg

// scfg provides a lightweight parser and memory representation for INI-style
// configuration files with support for sections, key-value properties, line continuations,
// quoted values, and boolean shorthand syntax.

// Global constants defining default configuration parsing settings.
const (
	// SCFG_GLOBAL_SECTION defines the key for properties defined before any explicit [section] header.
	SCFG_GLOBAL_SECTION = ""
	// SCFG_COMMENT_START contains characters that denote the beginning of a comment line.
	SCFG_COMMENT_START = ";#*!%"
	// SCFG_QUOTES contains supported quote characters for string values.
	SCFG_QUOTES = "'\"`"
	// SCFG_BOOL_TRUE is the internal string value saved for true boolean flags.
	SCFG_BOOL_TRUE = "TRUE"
	// SCFG_BOOL_FALSE is the internal string value saved for false boolean flags.
	SCFG_BOOL_FALSE = "FALSE"
)

// Scfg represents a full configuration structure containing named sections.
type Scfg struct {
	Sections map[string]*ScfgSection
}

// ScfgSection holds key-value property pairs for a single configuration section.
type ScfgSection struct {
	Name       string
	Properties map[string]string
}

// NewScfgSection initializes and returns a new empty ScfgSection.
func NewScfgSection() *ScfgSection {
	return &ScfgSection{Properties: make(map[string]string)}
}

// NewScfg initializes and returns a new empty Scfg configuration object.
func NewScfg() *Scfg {
	return &Scfg{Sections: make(map[string]*ScfgSection)}
}
