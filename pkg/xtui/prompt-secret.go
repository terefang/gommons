package xtui

import (
	"errors"
	"strings"
)

func ReadSecretString(prompt string) (string, error) {
	_pass1, _err := ReadInputString(prompt, true)
	if _err != nil {
		return "", _err
	}
	return _pass1, nil
}

func ReadSecretVerifyString(prompt1 string, prompt2 string) (string, error) {
	// take1
	_pass1, _err := ReadInputString(prompt1, true)
	if _err != nil {
		return "", _err
	}
	// take2
	_pass2, _err := ReadInputString(prompt2, true)
	if _err != nil {
		return "", _err
	}
	// verify
	if strings.Compare(_pass1, _pass2) == 0 {
		return _pass1, nil
	} else {
		return "", errors.New("passwords dont match")
	}
}
