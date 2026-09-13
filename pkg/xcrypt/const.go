package xcrypt

import (
	"encoding/base64"
	"math/rand"
	"strings"
)

const RADIX64 = "./0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const MCFRADIX = "./ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

var McfRadixEncoding = base64.NewEncoding(MCFRADIX)

func McfRadixEncode(src []byte) []byte {
	n := McfRadixEncoding.EncodedLen(len(src))
	dst := make([]byte, n)
	McfRadixEncoding.Encode(dst, src)
	for dst[n-1] == '=' {
		n--
	}
	return dst[:n]
}

func McfRadixDecode(src []byte) ([]byte, error) {
	numOfEquals := 4 - (len(src) % 4)
	for i := 0; i < numOfEquals; i++ {
		src = append(src, '=')
	}

	dst := make([]byte, McfRadixEncoding.DecodedLen(len(src)))
	n, err := McfRadixEncoding.Decode(dst, src)
	if err != nil {
		return nil, err
	}
	return dst[:n], nil
}

func GenerateRadix64Encoding(b []byte) string {
	_sb := strings.Builder{}
	_len := len(b)
	_mod := _len % 3
	_cnt := (_len - _mod)
	for i := 0; i < _cnt; i += 3 {
		_c := b[i]
		_b := b[i+1]
		_a := b[i+2]
		_v := (int(_c) << 16) | (int(_b) << 8) | int(_a)
		_sb.WriteString(Radix64Encode(_v, 4))
	}
	if _mod == 2 {
		_c := b[_len-2]
		_b := b[_len-1]
		_v := (int(_c) << 16) | (int(_b) << 8)
		_sb.WriteString(Radix64Encode(_v, 3))
	}
	if _mod == 1 {
		_c := b[_len-2]
		_v := int(_c) << 16
		_sb.WriteString(Radix64Encode(_v, 2))
	}
	return _sb.String()
}

func Radix64Encode(v, n int) string {
	_sb := strings.Builder{}
	for i := 0; i < n; i++ {
		_sb.WriteByte(RADIX64[(v>>18)&0x3f])
		v = v << 6
	}
	return _sb.String()
}

func GenerateRadix64Salt(slen int) string {
	_rl := len(RADIX64)
	_ret := make([]byte, slen)
	for i := 0; i < slen; i++ {
		_ret[i] = RADIX64[rand.Int()%_rl]
	}
	return string(_ret)
}
