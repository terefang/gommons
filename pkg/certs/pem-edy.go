package certs

import (
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"fmt"
	"os"

	"go.step.sm/crypto/x25519"
)

// id-X25519    OBJECT IDENTIFIER ::= { 1 3 101 110 }
// var privKey pkcs8
//
//	curvePrivateKey, err := asn1.Marshal(k.Seed())
//	if err != nil {
//		return nil, fmt.Errorf("x509: failed to marshal private key: %v", err)
//	}
//	privKey.PrivateKey = curvePrivateKey
//	return asn1.Marshal(privKey)
type pkcs8 struct {
	Version    int
	Algo       pkix.AlgorithmIdentifier
	PrivateKey []byte
}

func EncodeX25519PrivateKey(privateKey []byte) (string, error) {
	return EncodePKCS8PrivateKey(pkix.AlgorithmIdentifier{Algorithm: asn1.ObjectIdentifier{1, 3, 101, 110}}, privateKey)
}

func EncodePKCS8PrivateKey(algo pkix.AlgorithmIdentifier, privateKey []byte) (string, error) {
	var privKey pkcs8
	privKey.Algo = algo // pkix.AlgorithmIdentifier{ Algorithm: oidPublicKeyEd25519, }
	curvePrivateKey, err := asn1.Marshal(privateKey)
	if err != nil {
		return "", fmt.Errorf("pkcs8: failed to marshal private key: %v", err)
	}
	privKey.PrivateKey = curvePrivateKey
	_b, _ := asn1.Marshal(privKey)
	privatePem := &pem.Block{Type: "PRIVATE KEY", Bytes: _b}
	privateBytes := pem.EncodeToMemory(privatePem)
	return string(privateBytes), nil
}

func EncodeEdPrivateToPem(privateKey crypto.PrivateKey) string {
	_b, _ := x509.MarshalPKCS8PrivateKey(privateKey)
	privatePem := &pem.Block{Type: "PRIVATE KEY", Bytes: _b}
	privateBytes := pem.EncodeToMemory(privatePem)
	return string(privateBytes)
}

func EncodeEdPublicToPem(pKey crypto.PublicKey) string {
	_b, _ := x509.MarshalPKIXPublicKey(pKey)
	privatePem := &pem.Block{Type: "PUBLIC KEY", Bytes: _b}
	privateBytes := pem.EncodeToMemory(privatePem)
	return string(privateBytes)
}

func MakeEdKeyFile(f string) {
	_, caKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	WriteEdPrivateToPem(&caKey, f)
}

func MakeXdKeyFile(f string) {
	_, caKey, err := x25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	WriteXdPrivateToPem(&caKey, f)
}

func WriteXdPrivateToPem(pkey crypto.PrivateKey, f string) error {
	switch k := pkey.(type) {
	case *x25519.PrivateKey:
		_b := []byte(*k)
		_s, _ := EncodeX25519PrivateKey(_b)
		if f == "-" {
			fmt.Println(_s)
			return nil
		}
		return os.WriteFile(f, []byte(_s), os.FileMode(0600))
	}
	return errors.New("pkcs8: unknown algorithm")
}

func WriteEdPrivateToPem(pkey crypto.PrivateKey, f string) error {
	if f == "-" {
		fmt.Println(EncodeEdPrivateToPem(pkey))
		return nil
	}
	return os.WriteFile(f, []byte(EncodeEdPrivateToPem(pkey)), os.FileMode(0600))
}

func DecodeEdPrivateFromPem(data []byte) (ed25519.PublicKey, ed25519.PrivateKey, error) {
	_block, _ := pem.Decode(data)
	if _block == nil {
		return nil, nil, errors.New("decode pem error")
	}
	_key, _err := x509.ParsePKCS8PrivateKey(_block.Bytes)
	if _err != nil {
		return nil, nil, _err
	}
	_pkey := _key.(ed25519.PrivateKey)
	_key = ed25519.NewKeyFromSeed(_pkey.Seed())
	_pkey = _key.(ed25519.PrivateKey)
	_plkey := _pkey.Public().(ed25519.PublicKey)
	return _plkey, _pkey, nil
}

func ReadEdPrivateFromPem(f string) (ed25519.PublicKey, ed25519.PrivateKey, error) {
	_data, _err := os.ReadFile(f)
	if _err != nil {
		return nil, nil, _err
	}
	return DecodeEdPrivateFromPem(_data)
}
