package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strings"

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
	collide    bool
	doSha1     bool
	doMd5      bool
	doSha256   bool
	doSha512   bool
	doScrypt   bool
	doBcrypt   bool
	doArgon2   bool
	useUpper   string
	useLower   string
	useDigit   string
	useSpecial string
	doPbkdf2   bool
	altAlgo    bool
	safer      bool
	wsafe      bool
	password   string
	doPrompt   bool
	doCMd5     bool
	doAll      bool
	doApr1     bool
	doSsha     bool
}

// Arguments adds per-command args to the object.
func (p *passwordCommand) Arguments(f *flag.FlagSet) {
	f.IntVar(&p.length, "length", 15, "The length of the password to generate")
	f.BoolVar(&p.specials, "specials", true, "Should we use special characters?")
	f.BoolVar(&p.digits, "digits", true, "Should we use digits?")
	f.BoolVar(&p.collide, "ambiguous", false, "Should we allow ambiguous characters (0OD1lI8B5S2Z)?")
	f.BoolVar(&p.safer, "safer", false, "Should we use safer characters?")
	f.BoolVar(&p.wsafe, "wordsafe", false, "Should we use word-safe characters?")
	f.BoolVar(&p.altAlgo, "alt", false, "Should we use alternative password mixer?")
	f.BoolVar(&p.doSha1, "sha1", false, "show sha1 passwd/xcrypt hash")
	f.BoolVar(&p.doSsha, "ssha", false, "show salted sha1 ldap hash")
	f.BoolVar(&p.doMd5, "md5", false, "show md5 passwd/xcrypt hash")
	f.BoolVar(&p.doApr1, "apr1", false, "show apr1 passwd/xcrypt hash")
	f.BoolVar(&p.doCMd5, "cisco-md5", false, "show md5 cisco-style hash")
	f.BoolVar(&p.doSha256, "sha256", false, "show sha256 passwd/xcrypt hash")
	f.BoolVar(&p.doSha512, "sha512", false, "show sha512 passwd/xcrypt hash")
	f.BoolVar(&p.doScrypt, "scrypt", false, "show scrypt passwd/xcrypt hash")
	f.BoolVar(&p.doBcrypt, "bcrypt", false, "show bcrypt passwd/xcrypt hash")
	f.BoolVar(&p.doArgon2, "argon2", false, "show argon2 passwd/xcrypt hash")
	f.BoolVar(&p.doPbkdf2, "pbkdf2", false, "show pbkdf2 passwd/xcrypt hash")
	f.BoolVar(&p.doAll, "all", false, "show all hashes")
	f.StringVar(&p.useUpper, "U", " ", "Set upper case letters (overrides -ambiguous)")
	f.StringVar(&p.useLower, "L", " ", "Set lower case letters (overrides -ambiguous)")
	f.StringVar(&p.useDigit, "D", " ", "Set digit letters (overrides -ambiguous, -digits)")
	f.StringVar(&p.useSpecial, "S", " ", "Set digit letters (overrides -ambiguous, -specials)")
	f.StringVar(&p.password, "password", "", "dont generate password but use the one given")
	f.BoolVar(&p.doPrompt, "prompt", false, "prompt for a given password, more secure")
}

// Info returns the name of this subcommand.
func (p *passwordCommand) Info() (string, string) {
	return "gen-pass", `Generate a random password and hashes.

Details:

This command generates a simple random password, by default being 12
characters long.  You can tweak the alphabet used via the command-line
flags if necessary.`
}

// Execute is invoked if the user specifies `make-password` as the subcommand.
func (p *passwordCommand) Execute(args []string) int {

	// Alphabets we use for generation
	//
	// Notice that some items are removed as "ambiguous":
	//
	//     0O1lI8B5S2ZD
	//
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

	// "Safer" = [ "ACDEFGHJKLMNPQRTUVWXY" "abcdefghjkmnpqrtuvwxy" "34679" "%&*+" ]
	if p.safer {
		digits = "34679"
		upper = "ACDEFGHJKLMNPQRTUVWXY"
		lower = "abcdefghjkmnpqrtuvwxy"
		specials = "%&*+"
	}

	// "Word Safe" = [ "RVWXcfghjmpqrvwx" "23456789" "CFGHJMPQ" "" ]
	if p.wsafe {
		digits = "23456789"
		upper = "RVWXcfghjmpqrvwx"
		lower = "CFGHJMPQ"
		specials = ""
	}

	// check if use set symbols explicitely
	if p.useUpper != " " {
		upper = p.useUpper
	}
	if p.useLower != " " {
		lower = p.useLower
	}
	if p.useDigit != " " {
		digits = p.useDigit
		p.digits = true
	}
	if p.useSpecial != " " {
		specials = p.useSpecial
		p.specials = true
	}

	// Make a buffer and fill it with all characters
	buf := make([]byte, p.length)
	if !p.doPrompt && (p.password == "") {
		if p.altAlgo {
			_ll := len(lower)
			_lu := len(upper)
			_ld := len(digits)
			_ls := len(specials)
			for i := 0; i < p.length; {
				buf[i] = lower[rand.Intn(_ll)]
				i++
				if i == p.length {
					break
				}

				buf[i] = upper[rand.Intn(_lu)]
				i++
				if i == p.length {
					break
				}

				if p.digits && _ld > 0 {
					buf[i] = digits[rand.Intn(_ld)]
					i++
					if i == p.length {
						break
					}
				}

				if p.specials && _ls > 0 {
					buf[i] = specials[rand.Intn(_ls)]
					i++
					if i == p.length {
						break
					}
				}
			}

			// Shuffle
			rand.Shuffle(len(buf)-4, func(i, j int) {
				buf[i+4], buf[j+4] = buf[j+4], buf[i+4]
			})
		} else {
			all := upper + lower
			// Extend our alphabet if we should
			if p.digits {
				all = all + digits
			}
			if p.specials {
				all = all + specials
			}

			for i := 0; i < p.length; i++ {
				buf[i] = all[rand.Intn(len(all))]
			}

			buf[0] = lower[rand.Intn(len(lower))]
			buf[3] = upper[rand.Intn(len(upper))]

			// Add a digit if we should.
			//
			// We might already have them present, because our `all`
			// alphabet was used already.  But this ensures we have at
			// least one digit present.
			if p.digits {
				buf[2] = digits[rand.Intn(len(digits))]
			}

			// Add a special-character if we should.
			//
			// We might already have them present, because our `all`
			// alphabet was used already.  But this ensures we have at
			// least one special-character present.
			if p.specials {
				buf[1] = specials[rand.Intn(len(specials))]
			}

			// Shuffle
			rand.Shuffle(len(buf), func(i, j int) {
				buf[i], buf[j] = buf[j], buf[i]
			})

		}
		fmt.Printf("%s\n", buf)
	}

	if p.password != "" {
		buf = []byte(p.password)
	}

	if p.doPrompt {
		_pass1, _err := util.ReadInputString("Enter Password: ", true)
		if _err != nil {
			panic(_err)
		}
		_pass2, _err := util.ReadInputString("Re-Enter Password: ", true)
		if _err != nil {
			panic(_err)
		}

		if strings.Compare(_pass1, _pass2) == 0 {
			buf = []byte(_pass1)
		} else {
			panic("passwords dont match")
		}
	}

	if p.doAll || p.doCMd5 {
		_cry, _ := md5crypt.New(md5crypt.WithSaltLength(4))
		_dgst, _ := _cry.Hash(string(buf))
		fmt.Printf("Cisco-MD5: %s\n", _dgst.String())
	}
	if p.doAll || p.doMd5 {
		_cry, _ := md5crypt.New()
		_dgst, _ := _cry.Hash(string(buf))
		fmt.Printf("MD5: %s\n", _dgst.String())
	}
	if p.doAll || p.doApr1 {
		_cry, _ := md5crypt.New()
		_dgst, _ := _cry.Hash(string(buf))
		fmt.Printf("APR1: $apr1%s\n", _dgst.String()[2:])
	}
	if p.doAll || p.doSha1 {
		_cry, _ := sha1crypt.New(sha1crypt.WithIterations(1000))
		_dgst, _ := _cry.Hash(string(buf))
		fmt.Printf("SHA1: %s\n", _dgst.String())
		_cry, _ = sha1crypt.New(sha1crypt.WithSaltLength(16))
		_dgst, _ = _cry.Hash(string(buf))
		fmt.Printf("SHA1: %s\n", _dgst.String())
	}
	if p.doAll || p.doSsha {
		_dgst, _ := xcrypt.SshaCrypt(string(buf))
		fmt.Printf("SSHA: %s\n", _dgst)
	}
	if p.doAll || p.doSha256 {
		_cry, _ := shacrypt.New(shacrypt.WithSHA256(), shacrypt.WithIterations(5000))
		_dgst, _ := _cry.Hash(string(buf))
		fmt.Printf("SHA256: %s\n", _dgst.String())
	}
	if p.doAll || p.doSha512 {
		_cry, _ := shacrypt.New(shacrypt.WithSHA512(), shacrypt.WithIterations(5000))
		_dgst, _ := _cry.Hash(string(buf))
		fmt.Printf("SHA512: %s\n", _dgst.String())
	}
	if p.doAll || p.doScrypt {
		_cry, _ := scrypt.New()
		_dgst, _ := _cry.Hash(string(buf))
		fmt.Printf("SCRYPT: %s\n", _dgst.String())
	}
	if p.doAll || p.doBcrypt {
		_cry, _ := bcrypt.New(bcrypt.WithCost(14))
		_dgst, _ := _cry.Hash(string(buf))
		fmt.Printf("BCRYPT: %s\n", _dgst.String())
		_cry, _ = bcrypt.NewSHA256(bcrypt.WithCost(14))
		_dgst, _ = _cry.Hash(string(buf))
		fmt.Printf("BCRYPT: %s\n", _dgst.String())
	}
	if p.doAll || p.doPbkdf2 {
		_cry, _ := pbkdf2.NewSHA1()
		_dgst, _ := _cry.Hash(string(buf))
		fmt.Printf("PBKDF2: %s\n", _dgst.String())
		_cry, _ = pbkdf2.NewSHA256()
		_dgst, _ = _cry.Hash(string(buf))
		fmt.Printf("PBKDF2: %s\n", _dgst.String())
	}
	if p.doAll || p.doArgon2 {
		_cry, _ := argon2.New()
		_dgst, _ := _cry.Hash(string(buf))
		fmt.Printf("ARGON2: %s\n", _dgst.String())
	}

	return 0
}
