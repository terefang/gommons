package stemplate

import (
	"strings"

	"github.com/terefang/gommons/pkg/xstrings"
)

func BasicTemplateRenderer(template string, input map[string]any, _start string, _end string) string {
	_sb := strings.Builder{}
	_lt := len(template)
	_ls := len(_start)
	_le := len(_end)
	_last := 0
	for {
		_ofs := xstrings.IndexOf(template, _start, _last)
		if _ofs > -1 {
			_sb.WriteString(template[_last:_ofs])
			_last = _ls + _ofs
			_ofs = xstrings.IndexOf(template, _end, _last)
			if _ofs > 0 {
				_token := template[_last:_ofs]
				_repl, _ok := input[_token]
				if !_ok {
					_repl = template[_last-_ls : _ofs+_le]
				}
				_str, _ok := _repl.(string)
				if !_ok {
					//_sb.WriteString(stringify.String(_repl, nil))
					_sb.WriteString(xstrings.StringifyNoErrorWithLevel(_repl, 3))
				} else {
					_sb.WriteString(_str)
				}
				_last = _le + _ofs
			} else {
				break
			}
		} else {
			break
		}
	}

	if _last < _lt {
		_sb.WriteString(template[_last:])
	}

	return _sb.String()
}

func BasicTemplateRendererDefault(template string, input map[string]any) string {
	return BasicTemplateRenderer(template, input, "{{", "}}")
}
