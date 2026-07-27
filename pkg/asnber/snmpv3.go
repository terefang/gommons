package asnber

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"crypto/sha1"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"strings"
	"time"
	"unicode"
)

const (
	SNMPv3 SNMPVersion = 3
)

// //--go:generate stringer -type=SNMPError
type SNMPError uint8 // SNMPError is the type for standard SNMP errors.

// SNMP Errors
const (
	NoError             SNMPError = iota // No error occurred. This code is also used in all request PDUs, since they have no error status to report.
	TooBig                               // The size of the Response-PDU would be too large to transport.
	NoSuchName                           // The name of a requested object was not found.
	BadValue                             // A value in the request didn't match the structure that the recipient of the request had for the object. For example, an object in the request was specified with an incorrect length or type.
	ReadOnly                             // An attempt was made to set a variable that has an Access value indicating that it is read-only.
	GenErr                               // An error occurred other than one indicated by a more specific error code in this table.
	NoAccess                             // Access was denied to the object for security reasons.
	WrongType                            // The object type in a variable binding is incorrect for the object.
	WrongLength                          // A variable binding specifies a length incorrect for the object.
	WrongEncoding                        // A variable binding specifies an encoding incorrect for the object.
	WrongValue                           // The value given in a variable binding is not possible for the object.
	NoCreation                           // A specified variable does not exist and cannot be created.
	InconsistentValue                    // A variable binding specifies a value that could be held by the variable but cannot be assigned to it at this time.
	ResourceUnavailable                  // An attempt to set a variable required a resource that is not available.
	CommitFailed                         // An attempt to set a particular variable failed.
	UndoFailed                           // An attempt to set a particular variable as part of a group of variables failed, and the attempt to then undo the setting of other variables was not successful.
	AuthorizationError                   // A problem occurred in authorization.
	NotWritable                          // The variable cannot be written or created.
	InconsistentName                     // The name in a variable binding specifies a variable that does not exist.
)

// SnmpV3MsgFlags contains various message flags to describe Authentication, Privacy, and whether a report PDU must be sent.
type SnmpV3MsgFlags uint8

// Possible values of SnmpV3MsgFlags
const (
	NoAuthNoPriv   SnmpV3MsgFlags = 0x0 // No authentication, and no privacy
	AuthNoPriv     SnmpV3MsgFlags = 0x1 // Authentication and no privacy
	AuthPriv       SnmpV3MsgFlags = 0x3 // Authentication and privacy
	Reportable     SnmpV3MsgFlags = 0x4 // Report PDU must be sent.
	AuthPrivReport SnmpV3MsgFlags = 0x7 //Authentication and privacy + report PDU
)

// SnmpV3SecurityModel describes the security model used by a SnmpV3 connection
type SnmpV3SecurityModel int

// UserSecurityModel is the only SnmpV3SecurityModel currently implemented.
const (
	UserSecurityModel SnmpV3SecurityModel = 0x03
)

type V3user struct {
	User    string
	AuthAlg string //MD5 or SHA1
	AuthPwd string
	PrivAlg string //AES or DES
	PrivPwd string
}

// The object type that lets you do SNMP requests.
type SnmpHostv3 struct {
	SnmpHost
	//SNMP V3 variables
	User         string
	AuthAlg      string //MD5 or SHA1
	AuthPwd      string
	PrivAlg      string //AES or DES
	PrivPwd      string
	engineID     string
	MessageFlags SnmpV3MsgFlags
	//V3 temp variables
	AuthKey     string
	PrivKey     string
	engineBoots int32
	engineTime  int32
	desIV       uint32
	aesIV       int64
	Trapusers   []V3user
}

const (
	maxMsgSize int    = 65500
	SNMP_AES   string = "AES"
	SNMP_DES   string = "DES"
	SNMP_SHA1  string = "SHA1"
	SNMP_MD5   string = "MD5"
)

// NewSnmpHostv3 creates a new SnmpHostv3 object. Opens a udp connection to the device that will be used for the SNMP packets.
func NewSnmpHostv3(target, community string, version SNMPVersion, timeout time.Duration, retries int) (*SnmpHostv3, error) {
	targetPort := fmt.Sprintf("%s:161", target)
	conn, err := net.DialTimeout("udp", targetPort, timeout)
	if err != nil {
		return nil, fmt.Errorf(`error connecting to ("udp", "%s") : %s`, targetPort, err)
	}
	return &SnmpHostv3{SnmpHost: SnmpHost{
		Target:    target,
		Community: community,
		Version:   version,
		timeout:   timeout,
		retries:   retries,
		conn:      conn,
	},
	}, nil
}

func NewSnmpHostv3v3(w *SnmpHostv3, timeout time.Duration, retries int) (*SnmpHostv3, error) {
	if w.MessageFlags != NoAuthNoPriv && w.MessageFlags != AuthPrivReport {
		return nil, fmt.Errorf(`currently only NoAuthNoPriv(0x00) and AuthPrivReport(0x07) message flags are implemented`)
	}

	if w.MessageFlags == AuthPrivReport {
		if w.AuthAlg != SNMP_MD5 && w.AuthAlg != SNMP_SHA1 {
			return nil, fmt.Errorf(`invalid auth algorithm %s, needs SHA1 or MD5`, w.AuthAlg)
		}
		if w.PrivAlg != SNMP_AES && w.PrivAlg != SNMP_DES {
			return nil, fmt.Errorf(`invalid priv algorithm %s, needs AES or DES`, w.PrivAlg)
		}
	}

	targetPort := fmt.Sprintf("%s:161", w.Target)
	conn, err := net.DialTimeout("udp", targetPort, timeout)
	if err != nil {
		return nil, fmt.Errorf(`error connecting to ("udp", "%s") : %s`, targetPort, err)
	}
	return &SnmpHostv3{
		SnmpHost: SnmpHost{
			Target:  w.Target,
			Version: SNMPv3,
			timeout: timeout,
			retries: retries,
			conn:    conn,
		},
		User:         w.User,
		AuthAlg:      w.AuthAlg,
		AuthPwd:      w.AuthPwd,
		AuthKey:      w.AuthKey,
		PrivAlg:      w.AuthAlg,
		PrivPwd:      w.PrivPwd,
		PrivKey:      w.PrivKey,
		MessageFlags: w.MessageFlags,
	}, nil
}

/*
	NewSnmpHostv3OnConn creates a new SnmpHostv3 object from an existing net.Conn.

It does not check if the provided target is valid.
*/
func NewSnmpHostv3OnConn(target, community string, version SNMPVersion, timeout time.Duration, retries int, conn net.Conn) *SnmpHostv3 {
	return &SnmpHostv3{
		SnmpHost: SnmpHost{
			Target:  target,
			Version: SNMPv3,
			timeout: timeout,
			retries: retries,
			conn:    conn,
		},
	}
}

/*
IsStringAsciiPrintable checks if the given string is ASCII and is
printable form. Returns boolean value
*/
func IsStringAsciiPrintable(s string) bool {
	for _, r := range s {
		if r > unicode.MaxASCII || !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}

func passwordToKey(password string, engineID string, hashAlg string) string {
	h := sha1.New()
	if hashAlg == "MD5" {
		h = md5.New()
	}

	count := 0
	pLen := len(password)
	repeat := 1048576 / pLen
	remain := 1048576 % pLen
	for count < repeat {
		_, _ = io.WriteString(h, password)
		count++
	}
	if remain > 0 {
		_, _ = io.WriteString(h, password[:remain])
	}
	ku := string(h.Sum(nil))

	h.Reset()
	_, _ = io.WriteString(h, ku)
	_, _ = io.WriteString(h, engineID)
	_, _ = io.WriteString(h, ku)
	localKey := h.Sum(nil)

	return string(localKey)
}

// Generate a valid SNMP request ID.
func getRandomRequestID() int {
	return int(rand.Int31())
}

/*
SNMP V3 requires a discover packet being sent before a request being sent,

	so that agent's engineID and other parameters can be automatically detected
*/
func (w *SnmpHostv3) Discover() error {
	msgID := getRandomRequestID()
	requestID := getRandomRequestID()
	v3Header, _ := EncodeSequence([]interface{}{DerSequence, "", 0, 0, "", "", ""})
	flags := string([]byte{4})
	USM := 0x03
	req, err := EncodeSequence([]interface{}{
		DerSequence, int(w.Version),
		[]interface{}{DerSequence, msgID, maxMsgSize, flags, USM},
		string(v3Header),
		[]interface{}{DerSequence, "", "",
			[]interface{}{SnmpGetRequest, requestID, 0, 0, []interface{}{DerSequence}}}})
	if err != nil {
		fmt.Printf("Error encoding in discover:%v\n", err)
		panic(err)
	}

	response := make([]byte, bufSize)
	numRead, err := poll(w.conn, req, response, w.retries, w.timeout)
	if err != nil {
		return err
	}

	decodedResponse, err := DecodeSequence(response[:numRead])
	if err != nil {
		fmt.Printf("Error decoding discover:%v\n", err)
		panic(err)
	}

	//This helps in recovering from unknown panic situations in reading the packet data
	// Mostly errors for missing packet data
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovering from panic in Discover() for %v: %v \n", w.Target, r)
		}
	}()

	v3HeaderStr := decodedResponse[3].(string)
	v3HeaderDecoded, err := DecodeSequence([]byte(v3HeaderStr))
	if err != nil {
		fmt.Printf("Error 2 decoding:%v\n", err)
		return err
	}

	w.engineID = v3HeaderDecoded[1].(string)
	w.engineBoots = int32(v3HeaderDecoded[2].(int))
	w.engineTime = int32(v3HeaderDecoded[3].(int))
	w.aesIV = rand.Int63()
	w.desIV = rand.Uint32()

	//keys
	if w.AuthKey == "" && w.MessageFlags != NoAuthNoPriv {
		w.AuthKey = passwordToKey(w.AuthPwd, w.engineID, w.AuthAlg)
	}

	if w.PrivKey == "" && w.MessageFlags != NoAuthNoPriv {
		privKey := passwordToKey(w.PrivPwd, w.engineID, w.AuthAlg)
		w.PrivKey = string(([]byte(privKey))[0:16])
	}

	return nil
}

func EncryptDESCBC(dst, src, key, iv []byte) error {
	desBlockEncrypter, err := des.NewCipher(key)
	if err != nil {
		return err
	}
	desEncrypter := cipher.NewCBCEncrypter(desBlockEncrypter, iv)
	desEncrypter.CryptBlocks(dst, src)
	return nil
}

func DecryptDESCBC(dst, src, key, iv []byte) error {
	desBlockEncrypter, err := des.NewCipher(key)
	if err != nil {
		return err
	}
	desDecrypter := cipher.NewCBCDecrypter(desBlockEncrypter, iv)
	desDecrypter.CryptBlocks(dst, src)
	return nil
}

func EncryptAESCFB(dst, src, key, iv []byte) error {
	aesBlockEncrypter, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	aesEncrypter := cipher.NewCFBEncrypter(aesBlockEncrypter, iv)
	aesEncrypter.XORKeyStream(dst, src)
	return nil
}

func DecryptAESCFB(dst, src, key, iv []byte) error {
	aesBlockDecrypter, err := aes.NewCipher(key)
	if err != nil {
		return nil
	}
	aesDecrypter := cipher.NewCFBDecrypter(aesBlockDecrypter, iv)
	aesDecrypter.XORKeyStream(dst, src)
	return nil
}

func strXor(s1, s2 string) string {
	if len(s1) != len(s2) {
		panic("strXor called with two strings of different length\n")
	}
	n := len(s1)
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = s1[i] ^ s2[i]
	}
	return string(b)
}

func (w SnmpHostv3) auth(wholeMsg string) string {
	//Auth
	padLen := 64 - len(w.AuthKey)
	eAuthKey := w.AuthKey + strings.Repeat("\x00", padLen)
	ipad := strings.Repeat("\x36", 64)
	opad := strings.Repeat("\x5C", 64)
	k1 := strXor(eAuthKey, ipad)
	k2 := strXor(eAuthKey, opad)
	h := sha1.New()
	if w.AuthAlg == "MD5" {
		h = md5.New()
	}
	_, _ = io.WriteString(h, k1+wholeMsg)
	tmp1 := string(h.Sum(nil))
	h.Reset()
	_, _ = io.WriteString(h, k2+tmp1)
	msgAuthParam := string(h.Sum(nil)[:12])
	return msgAuthParam
}

func (w SnmpHostv3) encrypt(payload string) (string, string) {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.BigEndian, w.engineBoots)
	if w.PrivAlg == SNMP_AES {
		buf2 := new(bytes.Buffer)
		_ = binary.Write(buf2, binary.BigEndian, w.engineTime)
		buf3 := new(bytes.Buffer)
		w.aesIV += 1
		_ = binary.Write(buf3, binary.BigEndian, w.aesIV)
		privParam := string(buf3.Bytes())
		iv := string(buf.Bytes()) + string(buf2.Bytes()) + privParam

		// AES Encrypt
		encrypted := make([]byte, len(payload))
		err := EncryptAESCFB(encrypted, []byte(payload), []byte(w.PrivKey), []byte(iv))
		if err != nil {
			panic(err)
		}
		return string(encrypted), privParam
	} else {
		desKey := w.PrivKey[:8]
		preIV := w.PrivKey[8:16]
		buf2 := new(bytes.Buffer)
		w.desIV += 1
		_ = binary.Write(buf2, binary.BigEndian, w.desIV)
		privParam := string(buf.Bytes()) + string(buf2.Bytes())
		iv := strXor(preIV, privParam)

		//DES Encrypt
		plen := len(payload)
		//padding
		if (plen % 8) != 0 {
			payload = payload + strings.Repeat("\x00", 8-(plen%8))
		}
		encrypted := make([]byte, len(payload))
		_ = EncryptDESCBC(encrypted, []byte(payload), []byte(desKey), []byte(iv))
		return string(encrypted), privParam
	}
}

func (w SnmpHostv3) decrypt(payload, privParam string) (string, error) {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.BigEndian, w.engineBoots)

	if w.PrivAlg == SNMP_AES {
		buf2 := new(bytes.Buffer)
		_ = binary.Write(buf2, binary.BigEndian, w.engineTime)
		iv := string(buf.Bytes()) + string(buf2.Bytes()) + privParam

		// Decrypt
		decrypted := make([]byte, len(payload))
		err := DecryptAESCFB(decrypted, []byte(payload), []byte(w.PrivKey), []byte(iv))
		if err != nil {
			return "", err
		}
		return string(decrypted), nil
	} else {
		desKey := w.PrivKey[:8]
		preIV := w.PrivKey[8:16]
		iv := strXor(preIV, privParam)

		//DES Decrypt
		pLen := len(payload)
		if (pLen % 8) != 0 {
			return "", errors.New("DES encrypted payload is not multiple of 8 bytes")
		}
		decrypted := make([]byte, len(payload))
		_ = DecryptDESCBC(decrypted, []byte(payload), []byte(desKey), []byte(iv))
		return string(decrypted), nil
	}

}

// GetNext issues a GETNEXT SNMP request.
func (w *SnmpHostv3) GetNextV3(oid Oid) (string, interface{}, error) {
	return w.doGetV3(oid, SnmpGetNextRequest)
}

// Get
func (w *SnmpHostv3) GetV3(oid Oid) (interface{}, error) {
	_, val, err := w.doGetV3(oid, SnmpGetRequest)
	return val, err
}

// SetV3 sends an SNMP V3 set request to change the value associated with an oid.
func (w *SnmpHostv3) SetV3(oid Oid, value interface{}) (interface{}, error) {
	msgID := getRandomRequestID()
	requestID := getRandomRequestID()
	req, err := EncodeSequence(
		[]interface{}{DerSequence, w.engineID, "",
			[]interface{}{SnmpSetRequest, requestID, 0, 0,
				[]interface{}{DerSequence,
					[]interface{}{DerSequence, oid, value}}}})
	if err != nil {
		return nil, fmt.Errorf("error creating sequence: %w", err)
	}

	encrypted, privParam := w.encrypt(string(req))

	v3Header, err := EncodeSequence([]interface{}{DerSequence, w.engineID,
		int(w.engineBoots), int(w.engineTime), w.User, strings.Repeat("\x00", 12), privParam})
	if err != nil {
		return nil, fmt.Errorf("error creating v3 header: %w", err)
	}

	flags := string([]byte{7})
	USM := 0x03
	packet, err := EncodeSequence([]interface{}{
		DerSequence, int(w.Version),
		[]interface{}{DerSequence, msgID, maxMsgSize, flags, USM},
		string(v3Header),
		encrypted})
	if err != nil {
		return nil, fmt.Errorf("error assembling packet: %w", err)
	}
	authParam := w.auth(string(packet))
	finalPacket := strings.Replace(string(packet), strings.Repeat("\x00", 12), authParam, 1)

	response := make([]byte, bufSize)
	numRead, err := poll(w.conn, []byte(finalPacket), response, w.retries, w.timeout)
	if err != nil {
		return nil, fmt.Errorf("error with request: %w", err)
	}

	decodedResponse, err := DecodeSequence(response[:numRead])
	if err != nil {
		return nil, fmt.Errorf("error decoding request: %w", err)
	}

	v3HeaderStr := decodedResponse[3].(string)
	v3HeaderDecoded, err := DecodeSequence([]byte(v3HeaderStr))
	if err != nil {
		return nil, fmt.Errorf("error decoding header: %w", err)
	}

	w.engineID = v3HeaderDecoded[1].(string)
	w.engineBoots = int32(v3HeaderDecoded[2].(int))
	w.engineTime = int32(v3HeaderDecoded[3].(int))
	// skip checking authParam for now
	respAuthParam := v3HeaderDecoded[5].(string)
	respPrivParam := v3HeaderDecoded[6].(string)

	if len(respAuthParam) == 0 || len(respPrivParam) == 0 {
		return nil, fmt.Errorf("response is not encrypted")
	}

	encryptedResp := decodedResponse[4].(string)
	plainResp, err := w.decrypt(encryptedResp, respPrivParam)
	if err != nil {
		return nil, fmt.Errorf("error decrypt response: %w", err)
	}

	pduDecoded, err := DecodeSequence([]byte(plainResp))
	if err != nil {
		return nil, fmt.Errorf("error decoding pdu: %w", err)
	}

	// Find the varbinds
	respPacket := pduDecoded[3].([]interface{})
	if err := respPacket[2].(int); err != 0 {
		return nil, fmt.Errorf("error in setting snmp value: %s", SNMPError(err))
	}

	varbinds := respPacket[4].([]interface{})
	result := varbinds[1].([]interface{})[2]

	return result, nil
}

// SetMultipleV3 packages multiple SNMP set requests together in a single call
func (w *SnmpHostv3) SetMultipleV3(oidList map[string]interface{}) (map[string]interface{}, error) {
	msgID := getRandomRequestID()
	requestID := getRandomRequestID()

	oids := []interface{}{DerSequence}
	for oid, value := range oidList {
		oids = append(oids, []interface{}{DerSequence, oid, value})
	}

	req, err := EncodeSequence(
		[]interface{}{DerSequence, w.engineID, "",
			[]interface{}{SnmpSetRequest, requestID, 0, 0, oids}})
	if err != nil {
		return nil, fmt.Errorf("error creating sequence: %w", err)
	}

	encrypted, privParam := w.encrypt(string(req))

	v3Header, err := EncodeSequence([]interface{}{DerSequence, w.engineID,
		int(w.engineBoots), int(w.engineTime), w.User, strings.Repeat("\x00", 12), privParam})
	if err != nil {
		return nil, fmt.Errorf("error creating v3 header: %w", err)
	}

	flags := string([]byte{7})
	USM := 0x03
	packet, err := EncodeSequence([]interface{}{
		DerSequence, int(w.Version),
		[]interface{}{DerSequence, msgID, maxMsgSize, flags, USM},
		string(v3Header),
		encrypted})
	if err != nil {
		return nil, fmt.Errorf("error assembling packet: %w", err)
	}
	authParam := w.auth(string(packet))
	finalPacket := strings.Replace(string(packet), strings.Repeat("\x00", 12), authParam, 1)

	response := make([]byte, bufSize)
	numRead, err := poll(w.conn, []byte(finalPacket), response, w.retries, w.timeout)
	if err != nil {
		return nil, fmt.Errorf("error with request: %w", err)
	}

	decodedResponse, err := DecodeSequence(response[:numRead])
	if err != nil {
		return nil, fmt.Errorf("error decoding request: %w", err)
	}

	v3HeaderStr := decodedResponse[3].(string)
	v3HeaderDecoded, err := DecodeSequence([]byte(v3HeaderStr))
	if err != nil {
		return nil, fmt.Errorf("error decoding header: %w", err)
	}

	w.engineID = v3HeaderDecoded[1].(string)
	w.engineBoots = int32(v3HeaderDecoded[2].(int))
	w.engineTime = int32(v3HeaderDecoded[3].(int))
	// skip checking authParam for now
	respAuthParam := v3HeaderDecoded[5].(string)
	respPrivParam := v3HeaderDecoded[6].(string)

	if len(respAuthParam) == 0 || len(respPrivParam) == 0 {
		return nil, fmt.Errorf("response is not encrypted")
	}

	encryptedResp := decodedResponse[4].(string)
	plainResp, err := w.decrypt(encryptedResp, respPrivParam)
	if err != nil {
		return nil, fmt.Errorf("error decrypt response: %w", err)
	}

	pduDecoded, err := DecodeSequence([]byte(plainResp))
	if err != nil {
		return nil, fmt.Errorf("error decoding pdu: %w", err)
	}

	// Check if sets failed
	respPacket := pduDecoded[3].([]interface{})
	if err := respPacket[2].(int); err != int(NoError) {
		return nil, fmt.Errorf("error in setting snmp value: %w", errors.New(SNMPError(err).String()))
	}

	result := make(map[string]interface{})
	varbinds := respPacket[4].([]interface{})
	for _, v := range varbinds[1:] {
		o := v.([]interface{})[1].(Oid).String()
		value := v.([]interface{})[2]
		result[o] = value
	}

	return result, nil
}

func (w *SnmpHostv3) marshalV3(req []interface{}) (string, error) {
	var finalPacket string
	msgID := getRandomRequestID()
	flags := w.MessageFlags

	header := []interface{}{DerSequence, msgID, maxMsgSize, string(flags), int(UserSecurityModel)}

	switch flags {
	case NoAuthNoPriv:
		v3Header, _ := EncodeSequence([]interface{}{DerSequence, w.engineID,
			int(w.engineBoots), int(w.engineTime), w.User, "", ""})

		packet, err := EncodeSequence([]interface{}{
			DerSequence, int(w.Version), header,
			string(v3Header), req})
		if err != nil {
			return "", err
		}

		finalPacket = string(packet)
	case AuthPrivReport:
		reqEncoded, err := EncodeSequence(req)
		if err != nil {
			return "", err
		}

		encrypted, privParam := w.encrypt(string(reqEncoded))

		v3Header, err := EncodeSequence([]interface{}{DerSequence, w.engineID,
			int(w.engineBoots), int(w.engineTime), w.User, strings.Repeat("\x00", 12), privParam})
		if err != nil {
			return "", err
		}

		packet, err := EncodeSequence([]interface{}{
			DerSequence, int(w.Version), header,
			string(v3Header),
			encrypted})
		if err != nil {
			return "", err
		}

		authParam := w.auth(string(packet))
		finalPacket = strings.Replace(string(packet), strings.Repeat("\x00", 12), authParam, 1)
	default:
		return "", fmt.Errorf("incorrect message flag: %s", string(flags))
	}

	return finalPacket, nil
}

// A function does both GetNext and Get for SNMP V3
func (w *SnmpHostv3) doGetV3(oid Oid, request BERType) (string, interface{}, error) {
	requestID := getRandomRequestID()
	req := []interface{}{DerSequence, w.engineID, "",
		[]interface{}{request, requestID, 0, 0,
			[]interface{}{DerSequence,
				[]interface{}{DerSequence, oid, nil}}}}

	// Function to apply the right level of security parameters and PDU packet
	finalPacket, err := w.marshalV3(req)
	if err != nil {
		return "", nil, err
	}

	response := make([]byte, bufSize)
	numRead, err := poll(w.conn, []byte(finalPacket), response, w.retries, w.timeout)
	if err != nil {
		return "", nil, err
	}

	decodedResponse, err := DecodeSequence(response[:numRead])
	if err != nil {
		return "", nil, err
	}

	pduResponse, err := w.unMarshalV3(decodedResponse)
	if err != nil {
		return "", nil, err
	}

	// Find the varbinds
	respPacket := pduResponse[3].([]interface{})
	varbinds := respPacket[4].([]interface{})
	result := varbinds[1].([]interface{})

	resultOid := result[1].(string)
	resultVal := result[2]

	// Check if the value is string and printable. To distinguish HEX-String from normal string
	if res, ok := resultVal.(string); ok && !IsStringAsciiPrintable(resultVal.(string)) {
		return resultOid, fmt.Sprintf("%x", res), nil
	}
	return resultOid, resultVal, nil
}

// A function that does GetMultiple for SNMP V3
func (w *SnmpHostv3) GetMultipleV3(oids []Oid) (map[string]interface{}, error) {
	requestID := getRandomRequestID()

	varbinds := []interface{}{DerSequence}

	for _, oid := range oids {
		varbinds = append(varbinds, []interface{}{DerSequence, oid, nil})
	}

	req := []interface{}{DerSequence, w.engineID, "",
		[]interface{}{SnmpGetRequest, requestID, 0, 0, varbinds}}

	// Function to apply the right level of security parameters and PDU packet
	finalPacket, err := w.marshalV3(req)
	if err != nil {
		return nil, err
	}

	response := make([]byte, bufSize)
	numRead, err := poll(w.conn, []byte(finalPacket), response, w.retries, w.timeout)
	if err != nil {
		return nil, err
	}

	decodedResponse, err := DecodeSequence(response[:numRead])
	if err != nil {
		return nil, err
	}

	pduResponse, err := w.unMarshalV3(decodedResponse)
	if err != nil {
		return nil, err
	}

	// Find the varbinds
	respPacket := pduResponse[3].([]interface{})
	respVarbinds := respPacket[4].([]interface{})

	result := make(map[string]interface{})
	for _, v := range respVarbinds[1:] { // First element is just a sequence
		oid := v.([]interface{})[1].(string)
		value := v.([]interface{})[2]
		if value == nil {
			result[oid] = map[string]interface{}{
				"value": nil,
				"error": v.([]interface{})[3],
			}
		} else {
			// Check if the value is string and printable. To distinguish HEX-String from normal string
			if res, ok := value.(string); ok && !IsStringAsciiPrintable(value.(string)) {
				result[oid] = map[string]interface{}{
					"value": fmt.Sprintf("%x", res),
					"error": nil,
				}
			} else {
				result[oid] = map[string]interface{}{
					"value": value,
					"error": nil,
				}
			}
		}
	}

	return result, nil
}

func (w *SnmpHostv3) unMarshalV3(decodedResponse []interface{}) ([]interface{}, error) {
	v3HeaderStr := decodedResponse[3].(string)
	v3HeaderDecoded, err := DecodeSequence([]byte(v3HeaderStr))
	if err != nil {
		return nil, err
	}

	w.engineID = v3HeaderDecoded[1].(string)
	w.engineBoots = int32(v3HeaderDecoded[2].(int))
	w.engineTime = int32(v3HeaderDecoded[3].(int))
	// skip checking authParam for now
	respAuthParam := v3HeaderDecoded[5].(string)
	respPrivParam := v3HeaderDecoded[6].(string)

	if (len(respAuthParam) == 0 || len(respPrivParam) == 0) && w.MessageFlags == AuthPrivReport {
		return nil, fmt.Errorf("response is not encrypted")
	}
	var pduResponse []interface{}

	if w.MessageFlags == AuthPrivReport {
		encryptedResp := decodedResponse[4].(string)
		plainResp, err := w.decrypt(encryptedResp, respPrivParam)

		pduDecoded, err := DecodeSequence([]byte(plainResp))
		if err != nil {
			return nil, err
		}
		pduResponse = pduDecoded
	} else {
		pduResponse = decodedResponse[4].([]interface{})
	}

	return pduResponse, nil
}

// GetNext issues a GETNEXT SNMP request.
func (w SnmpHostv3) GetNext(oid Oid) (string, interface{}, error) {
	requestID := getRandomRequestID()
	req, err := EncodeSequence([]interface{}{DerSequence, int(w.Version), w.Community,
		[]interface{}{SnmpGetNextRequest, requestID, 0, 0,
			[]interface{}{DerSequence,
				[]interface{}{DerSequence, oid, nil}}}})
	if err != nil {
		return "", nil, err
	}

	response := make([]byte, bufSize)
	numRead, err := poll(w.conn, req, response, w.retries, w.timeout)
	if err != nil {
		return "", nil, err
	}

	decodedResponse, err := DecodeSequence(response[:numRead])
	if err != nil {
		return "", nil, err
	}

	//This helps in recovering from unknown panic situations in reading the packet data
	// Mostly errors for missing packet data
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovering from panic in GetNext() for %v: %v \n", w.Target, r)
		}
	}()

	// Find the varbinds
	respPacket := decodedResponse[3].([]interface{})
	varbinds := respPacket[4].([]interface{})
	result := varbinds[1].([]interface{})

	resultOid := result[1].(string)
	resultVal := result[2]

	return resultOid, resultVal, nil
}

// ParseTrapV3 parses a received SNMP trap and returns a map of oid to objects
func (w SnmpHostv3) ParseTrapV3(response []byte) (map[string]interface{}, error) {
	decodedResponse, err := DecodeSequence(response)
	if err != nil {
		return nil, err
	}

	// Fetch the varbinds out of the packet.
	snmpVer := decodedResponse[1].(int)
	// we might have received a v1/v2c trap
	if snmpVer == int(SNMPv1) {
		return w.ParseTrapV1(response)
	} else if snmpVer == int(SNMPv2c) {
		return w.ParseTrapV2(response)
	}
	v3HeaderStr := decodedResponse[3].(string)
	v3HeaderDecoded, err := DecodeSequence([]byte(v3HeaderStr))
	if err != nil {
		fmt.Printf("Error 2 decoding:%v\n", err)
		return nil, err
	}

	w.engineID = v3HeaderDecoded[1].(string)
	w.engineBoots = int32(v3HeaderDecoded[2].(int))
	w.engineTime = int32(v3HeaderDecoded[3].(int))
	w.User = v3HeaderDecoded[4].(string)
	respAuthParam := v3HeaderDecoded[5].(string)
	respPrivParam := v3HeaderDecoded[6].(string)

	if len(respAuthParam) == 0 || len(respPrivParam) == 0 {
		return nil, errors.New("response is not encrypted")
	}
	if len(w.Trapusers) == 0 {
		return nil, errors.New("no SNMP V3 trap user configured")
	}

	foundUser := false
	for _, v3user := range w.Trapusers {
		if v3user.User == w.User {
			w.AuthAlg = v3user.AuthAlg
			w.PrivAlg = v3user.PrivAlg
			w.AuthPwd = v3user.AuthPwd
			w.PrivPwd = v3user.PrivPwd
			foundUser = true
			break
		}
	}
	if !foundUser {
		return nil, errors.New("no matching user found")
	}

	//keys
	if w.AuthKey == "" {
		w.AuthKey = passwordToKey(w.AuthPwd, w.engineID, w.AuthAlg)
	}

	if w.PrivKey == "" {
		privKey := passwordToKey(w.PrivPwd, w.engineID, w.AuthAlg)
		w.PrivKey = string(([]byte(privKey))[0:16])
	}

	encryptedResp := decodedResponse[4].(string)
	plainResp, err := w.decrypt(encryptedResp, respPrivParam)
	if err != nil {

	}

	pduDecoded, err := DecodeSequence([]byte(plainResp))
	if err != nil {
		return nil, err
	}
	decodedResponse = pduDecoded

	respPacket := decodedResponse[3].([]interface{})
	var varbinds []interface{}
	if snmpVer == 1 {
		varbinds = respPacket[6].([]interface{})
	} else {
		varbinds = respPacket[4].([]interface{})
	}

	result := make(map[string]interface{})
	for i := 1; i < len(varbinds); i++ {
		oid := varbinds[i].([]interface{})[1].(Oid).String()
		val := varbinds[i].([]interface{})[2]
		result[oid] = val
	}
	fmt.Printf("\n")

	return result, nil
}

// ParseTrapV2 parses a received SNMP trap and returns a map of oid to objects
func (w SnmpHost) ParseTrapV2(response []byte) ([]interface{}, error) {
	decodedResponse, err := DecodeSequence(response)
	if err != nil {
		return nil, err
	}

	// Fetch the vars out of the packet.
	snmpVer := decodedResponse[1].(int)
	// we might have received a v1/v2c trap
	if snmpVer != int(SNMPv2c) {
		// ERROR
	}

	v12community := decodedResponse[2].(string)

	v12pdu := decodedResponse[3].([]interface{})

	plainResp := decodedResponse[4].(string)

	pduDecoded, err := DecodeSequence([]byte(plainResp))
	if err != nil {
		return nil, err
	}
	decodedResponse = pduDecoded

	respPacket := decodedResponse[3].([]interface{})
	var varbinds []interface{}
	if snmpVer == 1 {
		varbinds = respPacket[6].([]interface{})
	} else {
		varbinds = respPacket[4].([]interface{})
	}

	result := make([]interface{}, 0)
	for i := 0; i < len(varbinds); i++ {
		result = append(result, varbinds[i])
	}

	return result, nil
}
