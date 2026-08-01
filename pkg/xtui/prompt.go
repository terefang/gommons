package xtui

import (
	"os"

	"github.com/andrewchambers/goline"
	"github.com/pkg/errors"
	"github.com/terefang/gommons/pkg/xfile"
)

// ReadInput from stdin if something is detected or ask the user for an input
// using the given prompt.
func ReadInput(prompt string, mask bool) ([]byte, error) {
	st, err := os.Stdin.Stat()
	if err != nil {
		return nil, errors.Wrap(err, "error reading data")
	}

	if st.Size() == 0 && st.Mode()&os.ModeNamedPipe == 0 {
		//	return ui.PromptPassword(prompt)
		if mask {
			goline.MaskModeEnable()
		}
		_line, _err := goline.ReadLine(prompt)
		if mask {
			goline.MaskModeDisable()
		}
		if _err != nil {
			return nil, _err
		}
		return []byte(_line), nil
	}

	return xfile.ReadAll(os.Stdin)
}

func ReadInputString(prompt string, mask bool) (string, error) {
	_str, _err := ReadInput(prompt, mask)
	if _err != nil {
		return "", _err
	}
	return string(_str), nil
}
