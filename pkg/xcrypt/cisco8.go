package xcrypt

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

const (
	CISCO_TYPE_SLEN = 14
	CISCO_TYPE_LEN  = 32
	CISCO_TYPE8_C   = 20000
)

func GenerateCisco8CryptWithSalt(given, salt string) string {
	if salt == "" || !strings.HasPrefix(salt, "$8$") || len(salt) < CISCO_TYPE_SLEN+4 {
		salt = "$8$" + GenerateRadixSalt(CISCO_TYPE_SLEN) + "$none"
	}

	_key := GenerateKey([]byte(given), []byte(salt[3:CISCO_TYPE_SLEN+3]), CISCO_TYPE8_C, CISCO_TYPE_LEN, sha256.New)
	_rad := GenerateRadixEncoding(_key)
	return fmt.Sprintf("$8$%s$%s", salt[3:CISCO_TYPE_SLEN+3], _rad)
}
