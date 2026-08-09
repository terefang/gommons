package xini

import "testing"

func Test_File1(t *testing.T) {
	_cfg, err := NewIniConfig("./testfiles/file1.ini")
	if err != nil {
		t.Error(err)
	}
	if !_cfg.PropertyExists("section", "propertyName") {
		t.Error("property not exists")
	}
	if _str, err := _cfg.Get("section", "propertyName"); err == nil {
		if _str != "value" {
			t.Error("property not equal")
		}
	} else {
		t.Error("property lookup error", err)
	}

	if _b, err := _cfg.AsBool("section", "boolProperty"); err == nil {
		if _b != true {
			t.Error("property not equal")
		}
	} else {
		t.Error("property lookup error", err)
	}
}

func Test_Global(t *testing.T) {
	_cfg, err := NewIniConfig("./testfiles/global.ini")
	if err != nil {
		t.Error(err)
	}
	if !_cfg.PropertyExists(GLOBAL_SECTION, "propertyName") {
		t.Error("property not exists")
	}
	if _str, err := _cfg.Get(GLOBAL_SECTION, "propertyName"); err == nil {
		if _str != "value" {
			t.Error("property not equal")
		}
	} else {
		t.Error("property lookup error", err)
	}

	if _b, err := _cfg.AsBool(GLOBAL_SECTION, "boolProperty"); err == nil {
		if _b != true {
			t.Error("property not equal")
		}
	} else {
		t.Error("property lookup error", err)
	}
}

type parseFromStringStruct struct {
	A string `ini:"av"`
	B int64
	C int  `ini:"cv"`
	D bool `ini:"dv"`
	E uint64
	F uint `ini:"fv"`
	G float32
	H float64 `ini:"hv"`
	I []int
}

func TestIniConfig_ParseFromString(t *testing.T) {
	_cfg := New()
	err := _cfg.ParseFromString(`
av = "test-av"
b = -2
cv = -3
dv = FALSE
e = 5
fv = 6
g = 7.5
hv = 8.5
i = 1,2,3,4,5
`)
	if err != nil {
		t.Error(err)
	}

	_s := &parseFromStringStruct{}

	err = _cfg.Unmarshal(GLOBAL_SECTION, _s)
	if err != nil {
		t.Fatal(err)
	}

	if _s.A != "test-av" {
		t.Errorf("%s != %s", _s.A, "test-av")
	}
	if _s.B != -2 {
		t.Errorf("%v != %v", _s.B, -2)
	}
	if _s.C != -3 {
		t.Errorf("%v != %v", _s.C, -3)
	}
	if _s.D != false {
		t.Errorf("%v != %v", _s.D, false)
	}
	if _s.E != 5 {
		t.Errorf("%v != %v", _s.E, 5)
	}
	if _s.F != 6 {
		t.Errorf("%v != %v", _s.F, 6)
	}
	if _s.G != 7.5 {
		t.Errorf("%v != %v", _s.G, 7.5)
	}
	if _s.H != 8.5 {
		t.Errorf("%v != %v", _s.H, 8.5)
	}
	//t.Errorf("%v", _s.I)
}
