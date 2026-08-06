package xcrypt

import (
    "github.com/go-crypt/crypt/algorithm/md5crypt"
)

func GenerateCrypt1(_plain string) (string, error) {
    _cry, _err := md5crypt.New()
    if _err != nil {
        return "", _err
    }
    _dgst, _err := _cry.Hash(_plain)
    if _err != nil {
        return "", _err
    }
    return _dgst.String(), nil
}

func GenerateCrypt1WithSalt(_plain string, _salt string) (string, error) {
    _cry, _err := md5crypt.New()
    if _err != nil {
        return "", _err
    }
    _dgst, _err := _cry.HashWithSalt(_plain, []byte(_salt))
    if _err != nil {
        return "", _err
    }
    return _dgst.String(), nil
}
