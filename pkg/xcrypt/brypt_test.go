package xcrypt

import (
	"testing"
)

var testVector = []string{
	"123456", "$2z$14$H5sUgN8scsoSoyO9Hn6L1OZnt7zUQzB6xCi7R6U6Vf61Q/01PXDCO",
}

func Test_Bcrypt(t *testing.T) {
	for i := 0; i < len(testVector); i += 2 {
		_enc := GenerateBcrypt2Z(testVector[i])
		_ok, _err := VerifyBcrypt2Z(testVector[i], _enc)
		if _err != nil {
			t.Error(_err)
		}
		if !_ok {
			t.Error("bcrypt roundtrip failed")
		}
		_ok, _err = VerifyBcrypt2Z(testVector[i], _enc)
		if _err != nil {
			t.Error(_err)
		}
		if !_ok {
			t.Error("bcrypt testvector failed")
		}
	}
}
