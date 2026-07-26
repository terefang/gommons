package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/terefang/gommons/pkg/certs"
	"github.com/terefang/gommons/pkg/subcmd"
)

func init() {
	subcmd.Register(&GenKeyCommand{})
}

type GenKeyCommand struct {
	keyfile string
	bits    int
	useRsa  bool
	useEd   bool
	useXd   bool
}

func (r *GenKeyCommand) Arguments(f *flag.FlagSet) {
	f.StringVar(&r.keyfile, "K", "-", "key-file")
	f.IntVar(&r.bits, "bits", 2048, "minimum number of bits to generate key")
	f.BoolVar(&r.useRsa, "rsa", false, "use rsa style key")
	f.BoolVar(&r.useEd, "ed25519", false, "use ed25519 style key")
	f.BoolVar(&r.useXd, "x25519", false, "use x25519 style key")
}

func (r GenKeyCommand) Info() (string, string) {
	return "gen-key", `generate keys in pem format`
}

func (r GenKeyCommand) Execute(args []string) int {

	if len(args) > 0 {
		for _, arg := range args {
			if strings.HasSuffix(arg, ".key") {
				r.keyfile = arg
			}
		}
	}

	if r.keyfile == "" {
		r.Usage()
		return -1
	}

	if r.useXd {
		certs.MakeXdKeyFile(r.keyfile)
	} else if r.useEd {
		certs.MakeEdKeyFile(r.keyfile)
	} else {
		certs.MakeRsaKeyFile(r.bits, r.keyfile)
	}
	return 0
}

func (r GenKeyCommand) Usage() {
	fmt.Fprintln(os.Stderr, "usage: key-generate [flags] (private.key)")
}
