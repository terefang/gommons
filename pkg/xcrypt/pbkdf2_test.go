package xcrypt

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"testing"
)

type pbkdf2testvector struct {
	P  []byte
	S  []byte
	C  int
	L  int
	DK string
	H  func() hash.Hash
}

// RFC 6070 - PKCS #5 PBKDF2 Test Vectors
var pbkdf2sha1vectors = []pbkdf2testvector{
	{[]byte("password"), []byte("salt"), 1, 20, "0c60c80f961f0e71f3a9b524af6012062fe037a6", sha1.New},
	{[]byte("password"), []byte("salt"), 2, 20, "ea6c014dc72d6f8ccd1ed92ace1d41f0d8de8957", sha1.New},
	{[]byte("password"), []byte("salt"), 4096, 20, "4b007901b765489abead49d926f721d065a429c1", sha1.New},
	{[]byte("password"), []byte("salt"), 16777216, 20, "eefe3d61cd4da4e4e9945b3d6ba2158c2634e984", sha1.New},
	{[]byte("passwordPASSWORDpassword"), []byte("saltSALTsaltSALTsaltSALTsaltSALTsalt"), 4096, 25, "3d2eec4fe41c849b80c8d83662c0e44a8b291a964cf2f07038", sha1.New},
	{[]byte("pass\000word"), []byte("sa\000lt"), 4096, 16, "56fa6aa75548099dcc37d7f03425e0c3", sha1.New},
}

// sha256
var pbkdf2sha256vectors = []pbkdf2testvector{
	{[]byte("password"), []byte("salt"), 1, 32, "120fb6cffcf8b32c43e7225256c4f837a86548c92ccc35480805987cb70be17b", sha256.New},
}

func Test_GenerateKey_SHA1(t *testing.T) {
	for _, v := range pbkdf2sha1vectors {
		t.Run(fmt.Sprintf("%s-%s-%v-%v-SHA1", v.P, v.S, v.C, v.L), func(t *testing.T) {
			_dk := GenerateKey(v.P, v.S, v.C, v.L, v.H)
			_dkh := hex.EncodeToString(_dk)
			if _dkh != v.DK {
				t.Errorf("%s != %s", _dkh, v.DK)
			}
		})
	}
}

func Test_GenerateKey_SHA256(t *testing.T) {
	for _, v := range pbkdf2sha256vectors {
		t.Run(fmt.Sprintf("%s-%s-%v-%v-SHA256", v.P, v.S, v.C, v.L), func(t *testing.T) {
			_dk := GenerateKey(v.P, v.S, v.C, v.L, v.H)
			_dkh := hex.EncodeToString(_dk)
			if _dkh != v.DK {
				t.Errorf("%s != %s", _dkh, v.DK)
			}
		})
	}
}

/*

PBKDF2 HMAC-SHA256 Test Vectors

Input:
  P = "password" (8 octets)
  S = "salt" (4 octets)
  c = 2
  dkLen = 32

Output:
  DK = ae 4d 0c 95 af 6b 46 d3
       2d 0a df f9 28 f0 6d d0
       2a 30 3f 8e f3 c2 51 df
       d6 e2 d8 5a 95 47 4c 43 (32 octets)


Input:
  P = "password" (8 octets)
  S = "salt" (4 octets)
  c = 4096
  dkLen = 32

Output:
  DK = c5 e4 78 d5 92 88 c8 41
       aa 53 0d b6 84 5c 4c 8d
       96 28 93 a0 01 ce 4e 11
       a4 96 38 73 aa 98 13 4a (32 octets)


Input:
  P = "password" (8 octets)
  S = "salt" (4 octets)
  c = 16777216
  dkLen = 32

Output:
  DK = cf 81 c6 6f e8 cf c0 4d
       1f 31 ec b6 5d ab 40 89
       f7 f1 79 e8 9b 3b 0b cb
       17 ad 10 e3 ac 6e ba 46 (32 octets)


Input:
  P = "passwordPASSWORDpassword" (24 octets)
  S = "saltSALTsaltSALTsaltSALTsaltSALTsalt" (36 octets)
  c = 4096
  dkLen = 40

Output:
  DK = 34 8c 89 db cb d3 2b 2f
       32 d8 14 b8 11 6e 84 cf
       2b 17 34 7e bc 18 00 18
       1c 4e 2a 1f b8 dd 53 e1
       c6 35 51 8c 7d ac 47 e9 (40 octets)


Input:
  P = "pass\0word" (9 octets)
  S = "sa\0lt" (5 octets)
  c = 4096
  dkLen = 16

Output:
  DK = 89 b6 9d 05 16 f8 29 89
       3c 69 62 26 65 0a 86 87 (16 octets)

*/
