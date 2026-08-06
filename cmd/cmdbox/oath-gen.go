package main

import (
    "flag"
    "fmt"
    "os"

    "github.com/terefang/gommons/pkg/subcmd"
    "github.com/terefang/gommons/pkg/xbytes"
    "github.com/terefang/gommons/pkg/xcrypto"
)

func init() {
    subcmd.Register(&GenOathCommand{})
}

type GenOathCommand struct {
    outfile   string
    bytes     int
    seed      string
    doPrompt  bool
    useSHA256 bool
    useSHA512 bool
}

func (r *GenOathCommand) Arguments(f *flag.FlagSet) {
    f.StringVar(&r.outfile, "outfile", "-", "write to oath-file")
    f.StringVar(&r.seed, "seed", "", "key seed")
    f.BoolVar(&r.doPrompt, "prompt", false, "prompt for passphrase")
    f.IntVar(&r.bytes, "bytes", 10, "key length")
    f.BoolVar(&r.useSHA256, "sha256", false, "sha256 key")
    f.BoolVar(&r.useSHA512, "sha512", false, "sha256 key")
}

func (r GenOathCommand) Info() (string, string) {
    return "gen-oath", `generate oath/htop/top keys`
}

func (r *GenOathCommand) Execute(args []string) int {

    if r.outfile == "" {
        r.Usage()
        return -1
    }

    if r.seed == "" {
        r.seed = string(xcrypto.GenerateSalt(r.bytes))
    }

    _key := xcrypto.GenerateSecretKeyWithSalt(r.bytes, []byte(r.seed))
    _key = xbytes.ToBase32(_key)
    // _otp := xotp.NewConfig(string(_key))
    //    if r.useSHA256 {
    //        _otp.Hash = sha256.New
    //    }
    //    if r.useSHA512 {
    //        _otp.Hash = sha512.New
    //    }
    //    fmt.Println(_otp.ToUrl("totp", "issuer", "account").String())
    return 0
}

func (r GenOathCommand) Usage() {
    fmt.Fprintln(os.Stderr, "usage: gen-oath [flags]")
}
