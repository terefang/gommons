package util

import (
	"encoding/binary"
	"io"
	"os"

	"github.com/andrewchambers/goline"
	"github.com/pkg/errors"
	"github.com/terefang/gommons/pkg/xfile"
)

func ReadBeInt32(r io.Reader) (val int32) {
	_ = binary.Read(r, binary.BigEndian, &val)
	//if err != nil && err != io.EOF {
	//		f.err = err
	//	}
	return
}

func ReadBeInt16(r io.Reader) (val int16) {
	_ = binary.Read(r, binary.BigEndian, &val)
	//if err != nil && err != io.EOF {
	//		f.err = err
	//	}
	return
}

func ReadBeInt64(r io.Reader) (val int64) {
	_ = binary.Read(r, binary.BigEndian, &val)
	//if err != nil && err != io.EOF {
	//		f.err = err
	//	}
	return
}

func ReadLeInt32(r io.Reader) (val int32) {
	_ = binary.Read(r, binary.LittleEndian, &val)
	//if err != nil && err != io.EOF {
	//		f.err = err
	//	}
	return
}

func ReadLeInt16(r io.Reader) (val int16) {
	_ = binary.Read(r, binary.LittleEndian, &val)
	//if err != nil && err != io.EOF {
	//		f.err = err
	//	}
	return
}

func ReadLeInt64(r io.Reader) (val int64) {
	_ = binary.Read(r, binary.LittleEndian, &val)
	//if err != nil && err != io.EOF {
	//		f.err = err
	//	}
	return
}

func ReadByte(r io.Reader) (val byte) {
	_ = binary.Read(r, binary.BigEndian, &val)
	//if err != nil {
	//		f.err = err
	//	}
	return
}

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
