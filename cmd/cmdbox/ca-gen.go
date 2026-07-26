package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/terefang/gommons/pkg/certs"
	"github.com/terefang/gommons/pkg/ctokener"
	"github.com/terefang/gommons/pkg/stemplate"
	"github.com/terefang/gommons/pkg/subcmd"
)

func init() {
	subcmd.Register(&GenCaCommand{})
}

type GenCaCommand struct {
	template string
	keyfile  string
	certfile string
	bits     int
	days     int
	useRsa   bool
	CaDn     string
	useEd    bool
}

func (r *GenCaCommand) Arguments(f *flag.FlagSet) {
	f.StringVar(&r.template, "T", "", "template with defaults for ca attributes")
	f.StringVar(&r.keyfile, "K", "-", "private key-file for ca")
	f.StringVar(&r.certfile, "C", "-", "cert-file for ca")
	f.StringVar(&r.CaDn, "dn", "", "dn for ca")
	f.IntVar(&r.bits, "bits", 2048, "minimum number of bits to generate certificates")
	f.IntVar(&r.days, "days", 365*5, "validity of certificates")
	f.BoolVar(&r.useRsa, "rsa", false, "use rsa style key")
	f.BoolVar(&r.useEd, "edsa", false, "use ed25519 style key")
}

func (r GenCaCommand) Info() (string, string) {
	return "gen-ca", `generate a root ca key and certificate`
}

func (r GenCaCommand) Execute(args []string) int {

	if len(args) > 0 {
		for _, arg := range args {
			if strings.HasSuffix(arg, ".crt") || strings.HasSuffix(arg, ".cert") {
				r.certfile = arg
			} else if strings.HasSuffix(arg, ".key") {
				r.keyfile = arg
			} else if strings.HasSuffix(arg, ".tmpl") || strings.HasSuffix(arg, ".template") {
				r.template = arg
			}
		}
	}

	if r.certfile == "" {
		r.Usage()
		return -1
	}

	if r.keyfile == "" {
		r.Usage()
		return -1
	}

	r.TryReadTemplate()

	if r.CaDn == "" {
		r.Usage()
		return -1
	}

	if r.useEd {
		certs.MakeEdRootCa(r.days, r.CaDn, r.keyfile, r.certfile)
	} else {
		certs.MakeRsaRootCa(r.bits, r.days, r.CaDn, r.keyfile, r.certfile)
	}
	return 0
}

func (r GenCaCommand) Usage() {
	fmt.Fprintln(os.Stderr, "usage: ca-generate [flags] ca.tmpl ca.key ca.crt")
	fmt.Fprintln(os.Stderr, "       ca-generate [flags] ca.key ca.crt")
}

func (r *GenCaCommand) TryReadTemplate() {
	if r.template != "" {
		_data, _err := ctokener.LdataParseFile(r.template)
		if _err == nil {
			_ctx, _ok := _data.(map[string]any)
			if _ok {
				r.CaDn = stemplate.BasicTemplateRendererDefault("{{DN}}", _ctx)
			}
		}
	}
}
