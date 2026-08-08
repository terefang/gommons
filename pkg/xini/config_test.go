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
