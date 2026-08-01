package certs

import (
    "crypto/rand"
    "crypto/x509"
    "crypto/x509/pkix"
    "encoding/asn1"
    "encoding/pem"
    "log"
    "os"
)

func MakeRsaCrqFile(keyfile string, dn string, san []string, server bool, client bool, crqfile string) error {
    subj, err := ParseSimpleDN(dn)
    if err != nil {
        panic(err)
    }

    ckey, err := ReadRsaPrivateFromPem(keyfile)
    if err != nil {
        panic(err)
    }

    //[]x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth}
    template := &x509.CertificateRequest{
        Subject:            *subj,
        DNSNames:           san,
        SignatureAlgorithm: x509.SHA256WithRSA,
    }

    keyUsage := x509.KeyUsage(x509.KeyUsageDigitalSignature)
    extKeyUsage, err := marshalKeyUsage(keyUsage)
    if err != nil {
        log.Fatal(err)
    }

    // add this extension to csr like below
    template.ExtraExtensions = []pkix.Extension{extKeyUsage}
    //ExtKeyUsageServerAuth // serverAuth
    //ExtKeyUsageClientAuth // clientAuth

    creq, err := x509.CreateCertificateRequest(rand.Reader, template, ckey)
    if err != nil {
        panic(err)
    }

    return WriteRsaCrqToPem(creq, crqfile)
}

func WriteRsaCrqToPem(crq []byte, crqfile string) error {
    _buf := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: crq})
    return os.WriteFile(crqfile, _buf, 0600)
}

// func marshalExtKeyUsage(bserver, bclient bool) (pkix.Extension, error) {
//	ext := pkix.Extension{Id: x509.Ex.ExtKeyUsage.()ExtExtensionExtendedKeyUsage}
//
//	oids := make([]asn1.ObjectIdentifier, 0)
//	oids = append(oids, x509.ExtKeyUsageClientAuth.OID())
//	oids = append(oids, x509.ExtKeyUsageServerAuth.OID())
//
//	ext.Value, err = asn1.Marshal(oids)
//	return ext, err
//}

func marshalKeyUsage(ku x509.KeyUsage) (pkix.Extension, error) {
    ext := pkix.Extension{Id: asn1.ObjectIdentifier{2, 5, 29, 15}, Critical: true}

    var a [2]byte
    a[0] = reverseBitsInAByte(byte(ku))
    a[1] = reverseBitsInAByte(byte(ku >> 8))

    l := 1
    if a[1] != 0 {
        l = 2
    }

    bitString := a[:l]
    var err error

    ext.Value, err = asn1.Marshal(asn1.BitString{Bytes: bitString, BitLength: asn1BitLength(bitString)})
    if err != nil {
        return ext, err
    }
    return ext, nil
}

func reverseBitsInAByte(in byte) byte {
    b1 := in>>4 | in<<4
    b2 := b1>>2&0x33 | b1<<2&0xcc
    b3 := b2>>1&0x55 | b2<<1&0xaa
    return b3
}

func asn1BitLength(bitString []byte) int {
    bitLen := len(bitString) * 8

    for i := range bitString {
        b := bitString[len(bitString)-i-1]

        for bit := uint(0); bit < 8; bit++ {
            if (b>>bit)&1 == 1 {
                return bitLen
            }
            bitLen--
        }
    }

    return 0
}
