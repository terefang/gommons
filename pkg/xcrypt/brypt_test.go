package xcrypt

import (
	"testing"
)

var testVector = []string{
	"123456", "$2z$14$H5sUgN8scsoSoyO9Hn6L1OZnt7zUQzB6xCi7R6U6Vf61Q/01PXDCO",
	//$2$ (Original OpenBSD):
	"password", `$2$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy`,
	//$2a$ (Standard UTF-8 / Null-terminated):
	"password", `$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy`,
	//$2b$ (OpenBSD length-fix standard):
	"password", `$2b$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy`,
	//$2y$ (PHP bugfix standard):
	"password", `$2y$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy`,
	//$2x$ (Broken Sign-Extension Hash):
	`àèìòù`, `$2x$10$OU4dM67.Px0VRzAc29.JLu3lJkS3wve7VvPz5dC9Z43W8L30qH9G`,
	//$2y$ (Fixed Output for the same input):
	`àèìòù`, `$2y$10$OU4dM67.Px0VRzAc29.JLuv.x/S6G5qC.0zWfK0PjT0v5a7J/h6e`,
	//$2z$ alias for $bcrypt-sha256$
	"password", `$2z$12$n79VH.0Q2TMWmt3Oqt9uku$Kq4Noyk3094Y2QlB8NdRT8SvGiI4ft2`,
}

func Test_Bcrypt(t *testing.T) {
	for i := 1; i < len(testVector); i += 2 {
		_enc := GenerateBcrypt2Z(testVector[i-1])
		_ok, _err := VerifyBcrypt(testVector[i-1], _enc)
		if _err != nil {
			t.Error(_err)
		}
		if !_ok {
			t.Error("bcrypt roundtrip failed")
		}
		_ok, _err = VerifyBcrypt(testVector[i-1], testVector[i])
		if _err != nil {
			t.Error(_err)
		}
		if !_ok {
			t.Errorf("bcrypt testvector failed %s", testVector[i])
		}
	}
}
