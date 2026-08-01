package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/terefang/gommons/pkg/subcmd"
	"github.com/terefang/gommons/pkg/xcrypt"
	"github.com/terefang/gommons/pkg/xcrypto"
	"github.com/terefang/gommons/pkg/xtui"
)

func init() {
	subcmd.Register(&GenSKeyCommand{})
}

type GenSKeyCommand struct {
	keyfile  string
	bytes    int
	seed     string
	password string
	doPrompt bool
	doRandom bool
	doUuid   bool
}

func (r *GenSKeyCommand) Arguments(f *flag.FlagSet) {
	f.StringVar(&r.keyfile, "key", "-", "key-file")
	f.StringVar(&r.seed, "seed", "", "secret key seed")
	f.StringVar(&r.password, "password", "", "secret key passphrase")
	f.BoolVar(&r.doPrompt, "prompt", false, "prompt for passphrase")
	f.BoolVar(&r.doRandom, "random", false, "random passphrase")
	f.BoolVar(&r.doUuid, "uuid", false, "uuid passphrase")
	f.IntVar(&r.bytes, "bytes", 256, "minimum number of bytes to generate key")
}

func (r GenSKeyCommand) Info() (string, string) {
	return "gen-skey", `generate secret keys in pem format`
}

func (r GenSKeyCommand) Execute(args []string) int {

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

	if r.doUuid {
		_uuid, _ := uuid.NewV7()
		r.password = _uuid.String()
		fmt.Println(r.password)
	} else if r.doRandom {
		r.password = xcrypt.GeneratePassword(32)
		fmt.Println(r.password)
	}

	if r.doPrompt && r.password == "" {
		_pass, _err := xtui.ReadSecretVerifyString("Enter Passphrase: ", "Re-Enter Passphrase: ")
		if _err != nil {
			panic(_err)
		}
		r.password = _pass
	}

	if r.seed != "" && r.password != "" {
		xcrypto.CrypterSetupSecretKeyFileWithPassphraseSeed(r.bytes, r.password, []byte(r.seed), r.keyfile)
	} else if r.password != "" {
		xcrypto.CrypterSetupSecretKeyFileWithPassphraseSeed(r.bytes, r.password, nil, r.keyfile)
	} else if r.seed != "" {
		xcrypto.CrypterSetupSecretKeyFileWithPassphraseSeed(r.bytes, "", []byte(r.seed), r.keyfile)
	} else {
		xcrypto.CrypterSetupSecretKeyFileWithPassphraseSeed(r.bytes, "", nil, r.keyfile)
	}
	return 0
}

func (r GenSKeyCommand) Usage() {
	fmt.Fprintln(os.Stderr, "usage: gen-skey [flags] (secret.key)")
}
