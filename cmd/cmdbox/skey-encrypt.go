package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/terefang/gommons/pkg/subcmd"
	"github.com/terefang/gommons/pkg/xcrypto"
	"github.com/terefang/gommons/pkg/xtui"
)

func init() {
	subcmd.Register(&GenSKeyEncryptCommand{})
}

type GenSKeyEncryptCommand struct {
	keyfile   string
	password  string
	doDecrypt bool
	doPem     bool
	doPrompt  bool
	useCipher string
	useMode   string
	usePad    string
}

func (r *GenSKeyEncryptCommand) Arguments(f *flag.FlagSet) {
	f.StringVar(&r.keyfile, "key", "-", "key-file")
	f.StringVar(&r.password, "password", "", "secret key passphrase")
	f.StringVar(&r.useCipher, "cipher", "KHC-SHA1", "TODO")
	f.StringVar(&r.useMode, "mode", "CTR", "TODO")
	f.StringVar(&r.usePad, "pad", "PKCS7", "TODO")
	f.BoolVar(&r.doDecrypt, "decrypt", false, "decrypt instead")
	f.BoolVar(&r.doPrompt, "prompt", false, "prompt for passphrase")
	f.BoolVar(&r.doPem, "pem", false, "use pem content wrapper")
}

func (r GenSKeyEncryptCommand) Info() (string, string) {
	return "skey-encrypt", `encrypts file to file with secret key`
}

func (r GenSKeyEncryptCommand) Execute(args []string) int {

	if r.keyfile == "" || len(args) != 2 {
		r.Usage()
		return -1
	}

	if r.doPrompt && r.password == "" {
		_pass, _err := xtui.ReadSecretVerifyString("Enter Passphrase: ", "Re-Enter Passphrase: ")
		if _err != nil {
			panic(_err)
		}
		r.password = _pass
	}

	_key := xcrypto.CrypterSetupFromFileWithPassphraseSeed(0, r.password, nil, r.keyfile, false)

	if r.doDecrypt {
		_err := xcrypto.CrypterDecryptFile(_key, args[0], args[1], r.doPem)
		if _err != nil {
			panic(_err)
		}
	} else {
		_err := xcrypto.CrypterEncryptFile(_key, args[0], args[1], r.doPem)
		if _err != nil {
			panic(_err)
		}
	}
	return 0
}

func (r GenSKeyEncryptCommand) Usage() {
	fmt.Fprintln(os.Stderr, "usage: skey-encrypt [flags] plainfile cryptedfile")
}
