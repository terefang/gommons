package xcrypt

import "testing"

func TestVerifyCrypt1(t *testing.T) {
	_ok, _err := VerifyCrypt1("1234", `$1$7/KhLnYr$H3K/62BMf3tIDAnx4nXx40`)
	if _err != nil {
		t.Errorf("VerifyCrypt1() error = %v", _err)
	} else if _ok != true {
		t.Errorf("VerifyCrypt1() error = %v", _ok)
	}
}

func TestVerifyApr1(t *testing.T) {
	_ok := VerifyApr1("1234", `$apr1$8eNN9cX7$W2/Xl9Ju7kUKk.lyuhXXp1`)
	if _ok != true {
		t.Errorf("VerifyApr1() error = %v", _ok)
	}
}
