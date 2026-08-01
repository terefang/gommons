package xcrypt

import (
	"fmt"
	"strings"

	"github.com/go-crypt/crypt/algorithm/scrypt"
)

const (
	CISCO_TYPE9_N  = 16384
	CISCO_TYPE9_LN = 14
	CISCO_TYPE9_R  = 1
	CISCO_TYPE9_P  = 1
)

func GenerateCisco9CryptWithSalt(given, salt string) string {
	if salt == "" || !strings.HasPrefix(salt, "$9$") || len(salt) < CISCO_TYPE_SLEN+4 {
		salt = "$9$" + GenerateRadixSalt(CISCO_TYPE_SLEN) + "$none"
	}

	_cry, _ := scrypt.NewScrypt(scrypt.WithR(CISCO_TYPE9_R), scrypt.WithP(CISCO_TYPE9_P), scrypt.WithS(CISCO_TYPE_SLEN), scrypt.WithK(CISCO_TYPE_LEN), scrypt.WithLN(CISCO_TYPE9_LN))
	_key, _ := _cry.HashWithSalt(given, []byte(salt[3:CISCO_TYPE_SLEN+3]))
	_rad := GenerateRadixEncoding(_key.Key())
	return fmt.Sprintf("$9$%s$%s", salt[3:CISCO_TYPE_SLEN+3], _rad)
}
