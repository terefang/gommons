package main

import (
    "flag"
    "fmt"
    "os"

    "github.com/terefang/gommons/pkg/subcmd"
    "github.com/terefang/gommons/pkg/xotp"
)

func init() {
    subcmd.Register(&CalcOathCommand{})
}

type CalcOathCommand struct {
    uri     string
    window  int
    doSteam bool
    urifile string
}

func (r *CalcOathCommand) Arguments(f *flag.FlagSet) {
    f.StringVar(&r.uri, "uri", "", "token url")
    f.StringVar(&r.urifile, "urifile", "", "token file")
    f.IntVar(&r.window, "window", 0, "token window")
    f.BoolVar(&r.doSteam, "steam", false, "steam token")
}

func (r CalcOathCommand) Info() (string, string) {
    return "calc-oath", `calculate oath/hotp/totp token`
}

func (r *CalcOathCommand) Execute(args []string) int {
    if r.urifile != "" {
        _str, err := os.ReadFile(r.urifile)
        if err != nil {
            panic(err)
        }
        r.uri = string(_str)
    }

    _otp, _ := xotp.FromURL(r.uri)

    if r.doSteam {
        _otp.SymbolSet = xotp.SteamSymbolSet
        _otp.Digits = 5
    }

    if r.window == 0 {
        _code, _ := _otp.TOTP()
        fmt.Println(_code)
    } else {
        _codes, _ := _otp.TOTPWithWindow(r.window)
        for _, code := range _codes {
            fmt.Println(code)
        }
    }
    return 0
}

func (r CalcOathCommand) Usage() {
    fmt.Fprintln(os.Stderr, "usage: calc-oath [flags]")
}
