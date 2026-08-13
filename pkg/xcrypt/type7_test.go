package xcrypt

import (
	"fmt"
	"testing"
)

var _test = []string{
	"0822455D0A16", "cisco",
	"05080F1C2243", "cisco",
	"13061B13181F", "class",
	"044B0A151C36435C0D", "password",
	"1511021F0725082829202612121405140E445D0901075A564E41010107020006005E0D515F01", "ciscoClassPassword1234567890123456789",
}

func TestType7_decrypt(t *testing.T) {
	for i := 0; i < len(_test); i += 2 {
		_res, _ := Type7_decrypt(_test[i])
		if _res != _test[i+1] {
			t.Errorf("%s %s != %s", _test[i], _res, _test[i+1])
		}
	}
}

func TestType7_encrypt(t *testing.T) {
	fmt.Println(Type7_encrypt_salted(8, "cisco"))
	fmt.Println(Type7_encrypt_salted(5, "cisco"))
	fmt.Println(Type7_encrypt_salted(13, "class"))
	fmt.Println(Type7_encrypt_salted(4, "password"))
	fmt.Println(Type7_encrypt_salted(15, "ciscoClassPassword1234567890123456789"))
}
