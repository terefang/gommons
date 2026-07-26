package certs

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"math/big"
	"time"
)

func MakeRsaRootCa(bits int, days int, dn string, keyfile string, certfile string) {
	caKey, err := ReadRsaPrivateFromPem(keyfile)
	if err != nil {
		caKey, err = rsa.GenerateKey(rand.Reader, bits)
	}
	_ckl := caKey.Size()
	if _ckl < (bits / 8) {
		caKey, err = rsa.GenerateKey(rand.Reader, bits)
	}
	err = WriteRsaPrivateToPem(caKey, keyfile)
	if err != nil {
		panic(err)
	}

	algo := x509.SHA256WithRSA //for rsa
	notBefore := time.Now().Add(-12 * time.Hour)
	validFor := time.Duration(days+1) * 24 * time.Hour
	notAfter := notBefore.Add(validFor)
	subj, err := ParseSimpleDN(dn)
	if err != nil {
		panic(err)
	}
	serialNumber := big.NewInt(1)
	template := &x509.Certificate{
		Subject:            *subj,
		SerialNumber:       serialNumber,
		PublicKeyAlgorithm: x509.RSA,
		//PublicKey:             caKey.PublicKey,
		SignatureAlgorithm:    algo,
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, template, template, &caKey.PublicKey, caKey)
	if err != nil {
		panic(err)
	}
	err = WriteCertBytesToPem(certBytes, certfile)
	if err != nil {
		panic(err)
	}
}
