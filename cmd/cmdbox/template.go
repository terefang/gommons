package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/terefang/gommons/pkg/ctokener"
	"github.com/terefang/gommons/pkg/stemplate"
	"github.com/terefang/gommons/pkg/subcmd"
)

func init() {
	subcmd.Register(&TemplateCommand{})
}

type TemplateCommand struct {
	template     string
	templateFile string
	varStart     string
	varEnd       string
	useJinja     bool
	useStd       bool
	dataFile     string
	outFile      string
	useBasic     bool
}

func (r *TemplateCommand) Arguments(f *flag.FlagSet) {
	f.BoolVar(&r.useJinja, "jinja", false, "Template is in jinja/gonja format (ignores -s, -e)")
	f.BoolVar(&r.useStd, "std", false, "Template is in go/text/template format (ignores -s, -e)")
	f.BoolVar(&r.useBasic, "basic", true, "Template is in basic format")
	f.StringVar(&r.template, "T", "", "Template to format")
	f.StringVar(&r.templateFile, "F", "", "Template-file to format")
	f.StringVar(&r.outFile, "o", "", "file to write to (instead of stdout).")
	f.StringVar(&r.dataFile, "D", "", "Data-file to read (.json, .ldata)")
	f.StringVar(&r.varStart, "s", "{{", "template variable start-delimiter")
	f.StringVar(&r.varEnd, "e", "}}", "template variable end-delimiter")
}

func (r TemplateCommand) Info() (string, string) {
	return "template", `Populate a template-file.

Details:

The basic use-case of this sub-command is to allow substituting
variables into template-files (jinja, std, and basic).

*jinja* - uses gonja as engine
*std* - uses go/text/template as engine
*basic* - uses simple k/V substitution as engine

In addition to the flags below you can also specify variable 
assignments as parameters like "key=value".

if the value ends in '.ldata' or '.json' it is treated as a file
and the whole context is assigned under the specified key.
`
}

func (r TemplateCommand) Execute(args []string) int {
	ctx := make(map[string]any)

	if r.dataFile != "" {
		if strings.HasSuffix(r.dataFile, ".ldata") || strings.HasSuffix(r.dataFile, ".json") {
			_d, _err := ctokener.LdataParseFile(r.dataFile)
			if _err != nil {
				panic(_err)
			}
			ctx["_data"] = _d
		} else {
			panic("Data-file is unknown format: " + filepath.Base(r.dataFile))
		}
	}

	for _, arg := range args {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) != 2 {
			continue
		}

		if strings.HasSuffix(parts[1], ".ldata") || strings.HasSuffix(r.dataFile, ".json") {
			_d, _err := ctokener.LdataParseFile(parts[1])
			if _err != nil {
				panic(_err)
			}
			ctx[parts[0]] = _d
		} else {
			ctx[parts[0]] = parts[1]
		}
	}

	ctx["_ctx"] = ctx

	if r.templateFile != "" {
		_t, _err := os.ReadFile(r.templateFile)
		if _err != nil {
			panic(_err)
		}
		r.template = string(_t)
	}

	if r.outFile != "" {
		if r.useJinja {
			os.WriteFile(r.outFile, []byte(stemplate.GonjaRenderer(r.template, ctx)), 0640)
		} else if r.useStd {
			os.WriteFile(r.outFile, []byte(stemplate.GoTemplateRenderer(r.template, ctx, nil)), 0640)
		} else {
			os.WriteFile(r.outFile, []byte(stemplate.BasicTemplateRenderer(r.template, ctx, r.varStart, r.varEnd)), 0640)
		}
	} else {
		if r.useJinja {
			fmt.Println(stemplate.GonjaRenderer(r.template, ctx))
		} else if r.useStd {
			fmt.Println(stemplate.GoTemplateRenderer(r.template, ctx, nil))
		} else {
			fmt.Println(stemplate.BasicTemplateRenderer(r.template, ctx, r.varStart, r.varEnd))
		}
	}
	return 0
}
