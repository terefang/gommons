package certs

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
)

func EncodeRsaPrivateToPem(privateKey *rsa.PrivateKey) string {
	privatePem := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}
	privateBytes := pem.EncodeToMemory(privatePem)
	return string(privateBytes)
}

func WriteRsaPrivateToPem(pkey *rsa.PrivateKey, f string) error {
	if f == "-" {
		fmt.Println(EncodeRsaPrivateToPem(pkey))
		return nil
	}
	return os.WriteFile(f, []byte(EncodeRsaPrivateToPem(pkey)), os.FileMode(0600))
}

func DecodeRsaPrivateFromPem(data []byte) (*rsa.PrivateKey, error) {
	_block, _ := pem.Decode(data)
	if _block == nil {
		return nil, errors.New("decode pem error")
	}
	_key, _err := x509.ParsePKCS1PrivateKey(_block.Bytes)
	if _err != nil {
		return nil, _err
	}
	return _key, nil
}

func ReadRsaPrivateFromPem(f string) (*rsa.PrivateKey, error) {
	_data, _err := os.ReadFile(f)
	if _err != nil {
		return nil, _err
	}
	return DecodeRsaPrivateFromPem(_data)
}

func WriteCertBytesToPem(ber []byte, f string) error {
	if f == "-" {
		fmt.Println(EncodeCertBytesToPem(ber))
		return nil
	}
	return os.WriteFile(f, []byte(EncodeCertBytesToPem(ber)), os.FileMode(0600))
}

func EncodeCertBytesToPem(ber []byte) string {
	privatePem := &pem.Block{Type: "CERTIFICATE", Bytes: ber}
	privateBytes := pem.EncodeToMemory(privatePem)
	return string(privateBytes)
}

func MakeRsaKeyFile(bits int, f string) {
	caKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		panic(err)
	}
	WriteRsaPrivateToPem(caKey, f)
}
