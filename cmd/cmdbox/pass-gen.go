package main

import (
	"flag"
	"fmt"

	"github.com/go-crypt/crypt/algorithm/argon2"
	"github.com/go-crypt/crypt/algorithm/bcrypt"
	"github.com/go-crypt/crypt/algorithm/md5crypt"
	"github.com/go-crypt/crypt/algorithm/pbkdf2"
	"github.com/go-crypt/crypt/algorithm/scrypt"
	"github.com/go-crypt/crypt/algorithm/sha1crypt"
	"github.com/go-crypt/crypt/algorithm/shacrypt"
	"github.com/terefang/gommons/pkg/subcmd"
	"github.com/terefang/gommons/pkg/util"
	"github.com/terefang/gommons/pkg/xcrypt"
	"github.com/terefang/gommons/pkg/xtui"
)

func init() {
	subcmd.Register(&passwordCommand{})
}

// Structure for our options and state.
type passwordCommand struct {

	// The length of the password to generate
	length int

	// Specials?
	specials bool

	// Digits?
	digits bool

	// Confusing characters?
	collide  bool
	doSha1   bool
	doMd5    bool
	doSha256 bool
	doSha512 bool
	doScrypt bool
	doBcrypt bool
	doArgon2 bool
	doPbkdf2 bool
	doCMd5   bool
	doAll    bool
	doApr1   bool
	doSsha   bool
	password string
	doPrompt bool
	listSets bool
	useSet   string
	useSalt  string
	doCisco8 bool
	doCisco9 bool
}

// Arguments adds per-command args to the object.
func (p *passwordCommand) Arguments(f *flag.FlagSet) {
	f.IntVar(&p.length, "length", 15, "The length of the password to generate")
	f.BoolVar(&p.specials, "specials", true, "Should we use special characters?")
	f.BoolVar(&p.digits, "digits", true, "Should we use digits?")
	f.BoolVar(&p.collide, "ambiguous", false, "Should we allow ambiguous characters (0OD1lI8B5S2Z)?")
	f.BoolVar(&p.listSets, "list", false, "list predefined symbol sets?")
	f.BoolVar(&p.doSha1, "sha1", false, "show sha1 passwd/xcrypt hash")
	f.BoolVar(&p.doSsha, "ssha", false, "show salted sha1 ldap hash")
	f.BoolVar(&p.doMd5, "md5", false, "show md5 passwd/xcrypt hash")
	f.BoolVar(&p.doApr1, "apr1", false, "show apr1 passwd/xcrypt hash")
	f.BoolVar(&p.doCMd5, "cisco-md5", false, "show md5 cisco-style hash")
	f.BoolVar(&p.doCisco9, "cisco9", false, "show $9$ cisco-style hash")
	f.BoolVar(&p.doCisco8, "cisco8", false, "show $8$ cisco-style hash")
	f.BoolVar(&p.doSha256, "sha256", false, "show sha256 passwd/xcrypt hash")
	f.BoolVar(&p.doSha512, "sha512", false, "show sha512 passwd/xcrypt hash")
	f.BoolVar(&p.doScrypt, "scrypt", false, "show scrypt passwd/xcrypt hash")
	f.BoolVar(&p.doBcrypt, "bcrypt", false, "show bcrypt passwd/xcrypt hash")
	f.BoolVar(&p.doArgon2, "argon2", false, "show argon2 passwd/xcrypt hash")
	f.BoolVar(&p.doPbkdf2, "pbkdf2", false, "show pbkdf2 passwd/xcrypt hash")
	f.BoolVar(&p.doAll, "all", false, "show all hashes")
	f.StringVar(&p.useSalt, "salt", "", "use given passwd/xcrypt salt")
	f.StringVar(&p.useSet, "set", "", "use predefined symbol set")
	f.StringVar(&p.password, "password", "", "dont generate password but use the one given")
	f.BoolVar(&p.doPrompt, "prompt", false, "prompt for a given password, more secure")
}

// Info returns the name of this subcommand.
func (p *passwordCommand) Info() (string, string) {
	return "gen-pass", `Generate a random password and passwd/xcrypt hashes.

Details:

This command generates a simple random password, by default being 12
characters long.  You can tweak the alphabet used via the command-line
flags if necessary.`
}

// Execute is invoked if the user specifies `make-password` as the subcommand.
func (p *passwordCommand) Execute(args []string) int {
	if p.listSets {
		for k, _ := range xcrypt.PasswordSymbolSets {
			fmt.Println(k)
		}
		return 0
	}
	// Alphabets we use for generation
	// Notice that some items are removed as "ambiguous":
	//     0O1lI8B5S2ZD
	digits := "34679"
	specials := "+-*/%&="
	upper := "ACEFGHJKLMNPQRTUVWXY"
	lower := "abcdefghijkmnopqrstuvwxyz"

	// Reinstate the missing characters, if we need to.
	if p.collide {
		digits = "0123456789"
		upper = "ABCDEFGHIJKLMNOPQRSTUVWXY"
		lower = "abcdefghijklmnopqrstuvwxy"
		specials = "~=&+%^*/()[]{}/!@#$?|"
	}

	symbolset := make([]string, 0)
	symbolset = append(symbolset, upper, lower)
	if p.digits {
		symbolset = append(symbolset, digits)
	}
	if p.specials {
		symbolset = append(symbolset, specials)
	}

	if p.useSet != "" {
		sset, ok := xcrypt.PasswordSymbolSets[p.useSet]
		util.PanicIf(!ok, "invalid symbol set: %s", p.useSet)
		symbolset = sset
	}

	if !p.doPrompt && (p.password == "") {
		p.password = xcrypt.GeneratePasswordWithSym(symbolset, p.length)
		fmt.Printf("%s\n", p.password)
	}

	var buf []byte
	if p.password != "" {
		buf = []byte(p.password)
	}

	if p.doPrompt {
		_pass, _err := xtui.ReadSecretVerifyString("Enter Password: ","Re-Enter Password: ")
		if _err != nil {
			panic(_err)
		}
		p.password = _pass
		buf = []byte(_pass)
	}

	if p.doAll || p.doCMd5 {
		_cry, _ := md5crypt.New(md5crypt.WithSaltLength(4))
		_dgst, _ := _cry.Hash(string(buf))
		if p.useSalt != "" {
			_dgst, _ = _cry.HashWithSalt(string(buf), []byte(p.useSalt))
		}
		fmt.Printf("Cisco-MD5: %s\n", _dgst.String())
	}
	if p.doAll || p.doCisco8 {
		_dgst := xcrypt.GenerateCisco8CryptWithSalt(string(buf), p.useSalt)
		fmt.Printf("Cisco-$8$: %s\n", _dgst)
	}
	if p.doAll || p.doCisco9 {
		_dgst := xcrypt.GenerateCisco9CryptWithSalt(string(buf), p.useSalt)
		fmt.Printf("Cisco-$9$: %s\n", _dgst)
	}
	if p.doAll || p.doMd5 {
		_cry, _ := md5crypt.New()
		_dgst, _ := _cry.Hash(string(buf))
		if p.useSalt != "" {
			_dgst, _ = _cry.HashWithSalt(string(buf), []byte(p.useSalt))
		}
		fmt.Printf("MD5: %s\n", _dgst.String())
	}
	if p.doAll || p.doApr1 {
		_cry, _ := md5crypt.New()
		_dgst, _ := _cry.Hash(string(buf))
		if p.useSalt != "" {
			_dgst, _ = _cry.HashWithSalt(string(buf), []byte(p.useSalt))
		}
		fmt.Printf("APR1: $apr1%s\n", _dgst.String()[2:])
	}
	if p.doAll || p.doSha1 {
		_cry, _ := sha1crypt.New(sha1crypt.WithIterations(1000))
		_dgst, _ := _cry.Hash(string(buf))
		fmt.Printf("SHA1: %s\n", _dgst.String())
		_cry, _ = sha1crypt.New(sha1crypt.WithSaltLength(16))
		_dgst, _ = _cry.Hash(string(buf))
		if p.useSalt != "" {
			_dgst, _ = _cry.HashWithSalt(string(buf), []byte(p.useSalt))
		}
		fmt.Printf("SHA1: %s\n", _dgst.String())
	}
	if p.doAll || p.doSsha {
		_dgst, _ := xcrypt.SshaCrypt(string(buf))
		fmt.Printf("SSHA: %s\n", _dgst)
	}
	if p.doAll || p.doSha256 {
		_cry, _ := shacrypt.New(shacrypt.WithSHA256(), shacrypt.WithIterations(5000))
		_dgst, _ := _cry.Hash(string(buf))
		if p.useSalt != "" {
			_dgst, _ = _cry.HashWithSalt(string(buf), []byte(p.useSalt))
		}
		fmt.Printf("SHA256: %s\n", _dgst.String())
	}
	if p.doAll || p.doSha512 {
		_cry, _ := shacrypt.New(shacrypt.WithSHA512(), shacrypt.WithIterations(5000))
		_dgst, _ := _cry.Hash(string(buf))
		if p.useSalt != "" {
			_dgst, _ = _cry.HashWithSalt(string(buf), []byte(p.useSalt))
		}
		fmt.Printf("SHA512: %s\n", _dgst.String())
	}
	if p.doAll || p.doScrypt {
		_cry, _ := scrypt.New()
		_dgst, _ := _cry.Hash(string(buf))
		if p.useSalt != "" {
			_dgst, _ = _cry.HashWithSalt(string(buf), []byte(p.useSalt))
		}
		fmt.Printf("SCRYPT: %s\n", _dgst.String())
	}
	if p.doAll || p.doBcrypt {
		_cry, _ := bcrypt.New(bcrypt.WithCost(14))
		_dgst, _ := _cry.Hash(string(buf))
		if p.useSalt != "" {
			_dgst, _ = _cry.HashWithSalt(string(buf), []byte(p.useSalt))
		}
		fmt.Printf("BCRYPT: %s\n", _dgst.String())
		_cry, _ = bcrypt.NewSHA256(bcrypt.WithCost(14))
		_dgst, _ = _cry.Hash(string(buf))
		if p.useSalt != "" {
			_dgst, _ = _cry.HashWithSalt(string(buf), []byte(p.useSalt))
		}
		fmt.Printf("BCRYPT: %s\n", _dgst.String())
	}
	if p.doAll || p.doPbkdf2 {
		_cry, _ := pbkdf2.NewSHA1()
		_dgst, _ := _cry.Hash(string(buf))
		if p.useSalt != "" {
			_dgst, _ = _cry.HashWithSalt(string(buf), []byte(p.useSalt))
		}
		fmt.Printf("PBKDF2: %s\n", _dgst.String())
		_cry, _ = pbkdf2.NewSHA256()
		_dgst, _ = _cry.Hash(string(buf))
		if p.useSalt != "" {
			_dgst, _ = _cry.HashWithSalt(string(buf), []byte(p.useSalt))
		}
		fmt.Printf("PBKDF2: %s\n", _dgst.String())
	}
	if p.doAll || p.doArgon2 {
		_cry, _ := argon2.New()
		_dgst, _ := _cry.Hash(string(buf))
		if p.useSalt != "" {
			_dgst, _ = _cry.HashWithSalt(string(buf), []byte(p.useSalt))
		}
		fmt.Printf("ARGON2: %s\n", _dgst.String())
	}

	return 0
}
