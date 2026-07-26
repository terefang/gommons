package stemplate

import (
	"strings"
	"text/template"
)

func GoTemplateRenderer(_template string, _input map[string]any, _funcMap map[string]any) string {

	_templateBuild := template.New("_template")
	_xfuncMap := make(map[string]any)
	_xfuncMap["env"] = gtthEnv
	_xfuncMap["split"] = gtthSplit
	if _funcMap != nil {
		for _k, _v := range _funcMap {
			_xfuncMap[_k] = _v
		}
	}
	_templateBuild = _templateBuild.Funcs(_xfuncMap)
	// Parse the file
	_t := template.Must(_templateBuild.Parse(_template))
	// Render
	_sb := NewioStringWriter()
	_err := _t.Execute(_sb, _input)
	if _err != nil {
		return _err.Error()
	}
	return _sb.ToString()
}

func NewioStringWriter() *ioStringWriter {
	iosw := &ioStringWriter{}
	iosw.sb = new(strings.Builder)
	return iosw
}

type ioStringWriter struct {
	sb *strings.Builder
}

func (i ioStringWriter) Write(p []byte) (n int, err error) {
	_l := 0
	for _, b := range p {
		_err := i.sb.WriteByte(b)
		if _err != nil {
			return _l, _err
		}
		_l++
	}
	return _l, nil
}

func (i ioStringWriter) ToString() string {
	return i.sb.String()
}
