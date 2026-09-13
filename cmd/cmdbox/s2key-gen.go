package main

import (
	"crypto/sha1"
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
	subcmd.Register(&GenS2kKeyCommand{})
}

type GenS2kKeyCommand struct {
	keyfile  string
	bytes    int
	seed     string
	password string
	doPrompt bool
	doRandom bool
	doUuid   bool
}

func (r *GenS2kKeyCommand) Arguments(f *flag.FlagSet) {
	f.StringVar(&r.keyfile, "keyfile", "-", "key-file")
	f.StringVar(&r.seed, "seed", "", "secret key seed")
	f.StringVar(&r.password, "password", "", "secret key passphrase")
	f.BoolVar(&r.doPrompt, "prompt", false, "prompt for passphrase")
	f.BoolVar(&r.doRandom, "random", false, "random passphrase")
	f.BoolVar(&r.doUuid, "uuid", false, "uuid passphrase")
	f.IntVar(&r.bytes, "bytes", 256, "minimum number of bytes to generate key")
}

func (r GenS2kKeyCommand) Info() (string, string) {
	return "gen-s2k", `generate secret keys in s2k format`
}

func (r GenS2kKeyCommand) Execute(args []string) int {

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

	var bytes []byte
	hashf := sha1.New
	var err error
	if r.seed != "" && r.password != "" {
		bytes, err = xcrypto.S2kIteratedAndSalted([]byte(r.password), []byte(r.seed), 0x96, r.bytes, hashf)
	} else if r.seed != "" {
		bytes, err = xcrypto.S2kSalted([]byte(r.password), []byte(r.seed), r.bytes, hashf)
	} else {
		bytes, err = xcrypto.S2kSimple([]byte(r.password), r.bytes, hashf)
	}
	if err != nil {
		panic(err)
	}
	xcrypto.WriteSecretToHex(bytes, r.keyfile)
	return 0
}

func (r GenS2kKeyCommand) Usage() {
	fmt.Fprintln(os.Stderr, "usage: gen-s2k [flags] (secret.key)")
}
