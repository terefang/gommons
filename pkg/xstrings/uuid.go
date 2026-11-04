package xstrings

import (
	"encoding/base32"
	"fmt"
	"math/big"
	"math/rand"
	"time"
)

func Uuidv7() [16]byte {
	// random bytes
	var value [16]byte
	rand.Read(value[:])

	// current timestamp in ms
	timestamp := big.NewInt(time.Now().UnixMilli())

	// timestamp
	timestamp.FillBytes(value[0:6])

	// version and variant
	value[6] = (value[6] & 0x0F) | 0x70
	value[8] = (value[8] & 0x3F) | 0x80

	return value
}

func Uuidv7String() string {
	// random bytes
	value := Uuidv7()

	return fmt.Sprintf("%02x%02x%02x%02x-%02x%02x-%02x%02x-%02x%02x-%02x%02x%02x%02x%02x%02x",
		value[0], value[1], value[2], value[3],
		value[4], value[5], value[6], value[7],
		value[8], value[9], value[10], value[11],
		value[12], value[13], value[14], value[15])
}

func Uuidv7EncodedString() string {
	// random bytes
	value := Uuidv7()
	enc := base32.NewEncoding("23456789ABCDEFGHJKLMNPQRSTUVWXYZ").WithPadding(-1)

	return fmt.Sprintf("%s-%s-%s-%s",
		enc.EncodeToString(value[0:6]),
		enc.EncodeToString(value[6:9]),
		enc.EncodeToString(value[9:12]),
		enc.EncodeToString(value[12:16]))
}
