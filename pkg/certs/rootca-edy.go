package certs

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"math/big"
	"time"
)

func MakeEdRootCa(days int, dn string, keyfile string, certfile string) {
	caPub, caKey, err := ReadEdPrivateFromPem(keyfile)
	if err != nil {
		caPub, caKey, err = ed25519.GenerateKey(rand.Reader)
	}
	err = WriteEdPrivateToPem(&caKey, keyfile)
	if err != nil {
		panic(err)
	}

	algo := x509.PureEd25519
	notBefore := time.Now().Add(-12 * time.Hour)
	validFor := time.Duration(days+1) * 24 * time.Hour
	notAfter := notBefore.Add(validFor)
	subj, err := ParseSimpleDN(dn)
	if err != nil {
		panic(err)
	}
	serialNumber := big.NewInt(1)
	template := &x509.Certificate{
		Subject:               *subj,
		SerialNumber:          serialNumber,
		PublicKeyAlgorithm:    x509.Ed25519,
		SignatureAlgorithm:    algo,
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, template, template, caPub, caKey)
	if err != nil {
		panic(err)
	}
	err = WriteCertBytesToPem(certBytes, certfile)
	if err != nil {
		panic(err)
	}
}
