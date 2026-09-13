package xcfg

// Package xcfg provides functions for parsing Caddy/Nginx-style configuration
// files (.cdy) into structured tree blocks.

import (
	"fmt"
	"os"

	"strings"

	"github.com/terefang/gommons/pkg/xstrings"
)

type CfgConfig struct {
	Configs []*CfgConfigBlock

	inBlock     bool
	inExtString bool
	extString   strings.Builder
}

// CfgConfigBlock represents a node in the configuration tree.
// It holds a single directive command, its arguments, and any child configuration blocks.
type CfgConfigBlock struct {
	// Cmd is the directive name.
	Cmd string

	// Args contains the list of whitespace-delimited arguments following the command.
	Args []string

	Ext string
	// SubConfigs holds nested child configuration blocks defined within braces ({ ... }).
	SubConfigs []*CfgConfigBlock
}

// ParseCfgFile reads the file at the specified file path and parses its content
// into a root CfgConfigBlock.
// Returns an error if the file cannot be read or contains invalid syntax.
func ParseCfgFile(fpath string) (*CfgConfig, error) {
	_buf, err := os.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	return ParseCfgString(string(_buf))
}

// ParseCfgString parses a raw configuration string into a root CfgConfigBlock.
// Iterates through lines sequentially and returns an error on syntax violations.
func ParseCfgString(_str string) (*CfgConfig, error) {
	cfg := new(CfgConfig)
	_lines := strings.Split(_str, "\n")
	_pos := 0
	for _, _line := range _lines {
		_pos++

		err := cfg.ParseLine(_pos, _line)

		if err != nil {
			return nil, err
		}
	}
	if cfg.inBlock {
		return nil, fmt.Errorf("unclosed block at line %d", _pos)
	}
	if cfg.inExtString {
		return nil, fmt.Errorf("unclosed multi-line string at line %d", _pos)
	}
	return cfg, nil
}

// ParseLine processes a single line of text from the configuration file.
// Handles comments, multiline string boundaries, block depth checks, and single commands.
func (c *CfgConfig) ParseLine(_pos int, _line string) error {
	_sline := strings.TrimSpace(_line)

	// Handle end of multiline extended string block
	if c.inExtString && _sline == "```" {
		c.inExtString = false
		return c.appendExtString()
	}

	// Capture lines inside multiline extended string block
	if c.inExtString {
		c.extString.WriteString(_line + "\n")
		return nil
	}

	// Skip empty/comment lines
	if c.isEmptyOrComment(_sline) {
		return nil
	}

	// Handle block close
	if _sline == "}" {
		c.inBlock = false
		return nil
	}

	// Handle block start
	if strings.HasSuffix(_sline, " {") || _sline == "{" {
		if c.inBlock {
			return fmt.Errorf("invalid depth at line %d", _pos)
		}
		slen := len(_sline)
		c.appendLine(strings.TrimSpace(_sline[:slen-2]))
		c.inBlock = true
		return nil
	}

	// Handle start of multiline string directive
	if strings.HasSuffix(_sline, "```") {
		slen := len(_sline)
		c.appendLine(strings.TrimSpace(_sline[:slen-3]))
		c.inExtString = true
		return nil
	}

	// Reject malformed block end statements
	if strings.HasSuffix(_sline, "}") {
		return fmt.Errorf("invalid block-end at line %d", _pos)
	}

	return c.appendLine(_sline)
}

func (c *CfgConfig) isEmptyOrComment(_sline string) bool {
	// Skip empty lines
	if _sline == "" {
		return true
	}

	// Skip comment lines
	if strings.HasPrefix(_sline, "//") ||
		strings.HasPrefix(_sline, "--") ||
		(_sline[0] == '#') ||
		(_sline[0] == ';') ||
		(_sline[0] == ':') ||
		(_sline[0] == '%') ||
		(_sline[0] == '!') ||
		(_sline[0] == '$') ||
		(_sline[0] == '*') {
		return true
	}
	return false
}

func (c *CfgConfig) appendLine(_sline string) error {
	current := new(CfgConfigBlock)
	parts := xstrings.SplitWsWithDefaultQuotes(_sline)
	current.Cmd = parts[0]
	if len(parts) > 1 {
		current.Args = parts[1:]
	}

	sclen := len(c.Configs)
	if c.inBlock && sclen > 0 {
		block := c.Configs[sclen-1]
		block.SubConfigs = append(block.SubConfigs, current)
	} else {
		c.Configs = append(c.Configs, current)
	}
	return nil
}

func (c *CfgConfig) appendExtString() error {
	sclen := len(c.Configs)
	block := c.Configs[sclen-1]
	blen := len(block.SubConfigs)
	if c.inBlock {
		block.SubConfigs[blen-1].Ext = c.extString.String()
	} else {
		block.Ext = c.extString.String()
	}
	c.extString.Reset()
	return nil
}
