package xcrypt

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const ciscoType7Key = "dsfd;kfoA,.iyewrkldJKDHSUBsgvca69834ncxv9873254k;fg87"

func Type7_decrypt(encoded string) (string, error) {
	if len(encoded) < 2 {
		return "", fmt.Errorf("encoded string too short")
	}

	// Extract salt (first two digits)
	salt, err := strconv.Atoi(encoded[:2])
	if err != nil {
		return "", err
	}

	// The rest is hex-encoded XORed data
	hexData := encoded[2:]
	decoded, err := hex.DecodeString(hexData)
	if err != nil {
		return "", err
	}

	var result []byte
	for i, b := range decoded {
		// XOR with the key at index (i + salt) % 53
		keyChar := ciscoType7Key[(i+salt)%len(ciscoType7Key)]
		result = append(result, b^keyChar)
	}

	return string(result), nil
}

func Type7_encrypt(plain string) string {
	// Seed random for salt generation
	rand.Seed(time.Now().UnixNano())

	// Generate random salt between 0 and 15
	salt := rand.Intn(16)

	return Type7_encrypt_salted(salt, plain)
}

func Type7_encrypt_salted(salt int, plain string) string {

	// Truncate password to 25 characters if necessary (Cisco IOS limit)
	//if len(plain) > 25 {
	//    plain = plain[:25]
	//}

	// Start with the zero-padded salt (e.g., "07")
	result := fmt.Sprintf("%02d", salt)

	// XOR each character with the key
	for i := 0; i < len(plain); i++ {
		// Calculate key index: (salt + i) % 53
		keyIndex := (salt + i) % len(ciscoType7Key)
		keyChar := ciscoType7Key[keyIndex]

		// XOR the plaintext byte with the key character
		xored := plain[i] ^ keyChar

		// Append as uppercase hex
		result += fmt.Sprintf("%02X", xored)
	}

	return result
}

func ValidateType7Credential(_given string, _encrypted string) (bool, error) {
	if strings.HasPrefix(_encrypted, "{type7}") {
		_res, err := Type7_decrypt(_encrypted[7:])
		if err != nil {
			return false, err
		}
		return _res == _given, nil
	}
	return false, errors.New("wrong mcf prefix: " + _encrypted)
}
