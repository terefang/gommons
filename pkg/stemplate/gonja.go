package stemplate

import (
	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/exec"
)

func GonjaRenderer(template string, input map[string]any) string {
	_tmpl, err := gonja.FromString(template)
	if err != nil {
		panic(err)
	}

	data := exec.EmptyContext()
	for k, v := range input {
		data.Set(k, v)
	}
	_str, err := _tmpl.ExecuteToString(data)
	if err != nil {
		panic(err)
	}
	return _str
}
