package main

import (
    "flag"
    "fmt"
    "os"

    qrcode "github.com/skip2/go-qrcode"
    "github.com/terefang/gommons/pkg/subcmd"
    "github.com/terefang/gommons/pkg/xbytes"
    "github.com/terefang/gommons/pkg/xcrypto"
    "github.com/terefang/gommons/pkg/xotp"
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
    useRAW    bool
    useOATH   bool
    useQrcode bool
    issuer    string
    account   string
    digits    int
    period    int
}

func (r *GenOathCommand) Arguments(f *flag.FlagSet) {
    f.StringVar(&r.outfile, "outfile", "-", "write to oath-file")
    f.StringVar(&r.seed, "seed", "", "key seed")
    f.StringVar(&r.issuer, "issuer", "", "issuer tag")
    f.StringVar(&r.account, "account", "", "account tag")
    f.BoolVar(&r.doPrompt, "prompt", false, "prompt for passphrase")
    f.IntVar(&r.bytes, "bytes", 10, "key length")
    f.IntVar(&r.digits, "digits", 6, "token digits")
    f.IntVar(&r.period, "period", 30, "token period")
    f.BoolVar(&r.useRAW, "raw", false, "raw format")
    f.BoolVar(&r.useQrcode, "qrcode", false, "qr-code format")
    f.BoolVar(&r.useOATH, "oath", false, "oath format")
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
    if r.useRAW {
        fmt.Println(string(xbytes.ToBase32(_key)))
    } else {
        // create the token
        _otp, err := xotp.From(_key, r.digits, xotp.DefaultAlgorithm)
        if err != nil {
            panic(err)
        }
        _otp.Period = r.period
        if r.issuer != "" {
            _otp.Issuer = r.issuer
        }
        if r.account != "" {
            _otp.Account = r.account
        }
        if r.useSHA256 {
            _otp.Algorithm = xotp.AlgorithmSHA256
        }
        if r.useSHA512 {
            _otp.Algorithm = xotp.AlgorithmSHA512
        }
        if r.useOATH {
            // to url, encode in b64 with prefix "{OATH}"
            fmt.Printf("{OATH}%s\n", string(xbytes.ToBase64([]byte(_otp.ToURL()))))
        } else {
            fmt.Println(_otp.ToURL())
        }

        if r.useQrcode {
            _qr, err := qrcode.New(_otp.ToURL(), qrcode.Low)
            if err != nil {
                panic(err)
            }
            fmt.Println(_qr.ToString(false))
        }
    }
    return 0
}

func (r GenOathCommand) Usage() {
    fmt.Fprintln(os.Stderr, "usage: gen-oath [flags]")
}
