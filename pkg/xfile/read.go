package xfile

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
	"unicode"

	"github.com/pkg/errors"
	"github.com/terefang/gommons/pkg/utfbom"
)

// ReadAll returns a slice of bytes with the content of the given reader.
func ReadAll(r io.Reader) ([]byte, error) {
	b, err := io.ReadAll(r)
	return b, errors.Wrap(err, "error reading data")
}

// ReadString reads one line from the given io.Reader.
func ReadString(r io.Reader) (string, error) {
	br := bufio.NewReader(r)
	str, err := br.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", errors.Wrap(err, "error reading string")
	}
	return strings.TrimSpace(str), nil
}

// ReadPasswordFromFile reads and returns the password from the given filename.
// The contents of the file will be trimmed at the right.
func ReadPasswordFromFile(filename string) ([]byte, error) {
	password, err := os.ReadFile(filename) // #nosec G703 -- file intended to be provided by user
	if err != nil {
		return nil, errors.New(err.Error() + ": " + filename)
	}
	password = bytes.TrimRightFunc(password, unicode.IsSpace)
	return password, nil
}

// ReadStringPasswordFromFile reads and returns the password from the given filename.
// The contents of the file will be trimmed at the right.
func ReadStringPasswordFromFile(filename string) (string, error) {
	b, err := ReadPasswordFromFile(filename)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ReadFile returns the contents of the file identified by name. It reads from
// STDIN if name is a hyphen ("-").
func ReadFile(name string) (b []byte, err error) {
	if name == "-" {
		name = "/dev/stdin"
		b, err = io.ReadAll(os.Stdin)
	} else {
		var contents []byte
		contents, err = os.ReadFile(name)
		if err != nil {
			return nil, errors.New(err.Error() + ": " + name)
		}
		b, err = io.ReadAll(utfbom.SkipOnly(bytes.NewReader(contents)))
	}
	if err != nil {
		return nil, errors.New(err.Error() + ": " + name)
	}
	return b, nil
}
