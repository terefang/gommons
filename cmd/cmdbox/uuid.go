package main

import (
	"flag"
	"fmt"

	"github.com/terefang/gommons/pkg/subcmd"
	"github.com/terefang/gommons/pkg/xstrings"
)

type UUIDCommand struct {
	useBase32 bool
}

func (r UUIDCommand) Info() (string, string) {
	return "uuid", "print next uuidv7"
}

func (r UUIDCommand) Execute(args []string) int {
	if r.useBase32 {
		fmt.Println(xstrings.Uuidv7EncodedString())
	} else {
		fmt.Println(xstrings.Uuidv7String())
	}
	return 0
}

func (r *UUIDCommand) Arguments(f *flag.FlagSet) {
	f.BoolVar(&r.useBase32, "base32", false, "use base32 encoding")
}
func init() {
	subcmd.Register(&UUIDCommand{})
}
