package main

import (
	"flag"
	"fmt"

	"github.com/terefang/gommons/pkg/subcmd"
	"github.com/terefang/gommons/pkg/xfnt"
	"github.com/terefang/gommons/pkg/xstrings"
)

type fntIndexCommand struct {
	verbose bool
	ini     bool
}

func (r fntIndexCommand) Info() (string, string) {
	return "font-index", "scan and print font meta info"
}

func (r fntIndexCommand) Execute(args []string) int {
	for _, arg := range args {
		fl, err := xfnt.ScanFontDirectory(arg)
		if err == nil {
			for _, f := range fl.List {
				if r.ini {
					fmt.Printf(`[%s]
Path="%s"
Full="%s"
Name="%s"
Variant="%s"
Style="%s"
Weight=%d
Stretch=%d

`, xstrings.Slugify(fmt.Sprintf("%s %s %d %d", f.Family, f.Style, f.Weight, f.Stretch)),
						f.RelPath, f.Name, f.Family, f.SubFamily, f.Style, f.Weight, f.Stretch)
				} else {
					f.Print()
				}
			}
		} else {
			fmt.Println(err)
		}
	}
	return 0
}

func (r *fntIndexCommand) Arguments(f *flag.FlagSet) {
	f.BoolVar(&r.ini, "ini", false, "ini output")
}
func init() {
	subcmd.Register(&fntIndexCommand{})
}
