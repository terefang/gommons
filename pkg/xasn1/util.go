package xasn1

// MarshalLengthBytes returns the ASN1 encoded bytes for the length 'l'
//
// There are two forms: short (for lengths between 0 and 127), and long definite (for lengths between 0 and 2^1008 -1).
//
// Short form: One octet. Bit 8 has value "0" and bits 7-1 give the length.
//
// Long form: Two to 127 octets. Bit 8 of first octet has value "1" and bits 7-1 give the number of additional length octets. Second and following octets give the length, base 256, most significant digit first.
func MarshalLengthBytes(l int) []byte {
	if l <= 127 {
		return []byte{byte(l)}
	}
	var b []byte
	p := 1
	for i := 1; i < 127; {
		b = append([]byte{byte((l % (p * 256)) / p)}, b...)
		p = p * 256
		l = l - l%p
		if l <= 0 {
			break
		}
	}
	return append([]byte{byte(128 + len(b))}, b...)
}

// GetLengthFromASN returns the length of a slice of ASN1 encoded bytes from the ASN1 length header it contains.
func GetLengthFromASN(b []byte) int {
	if int(b[1]) <= 127 {
		return int(b[1])
	}
	// The bytes that indicate the length
	lb := b[2 : 2+int(b[1])-128]
	base := 1
	l := 0
	for i := len(lb) - 1; i >= 0; i-- {
		l += int(lb[i]) * base
		base = base * 256
	}
	return l
}

// GetNumberBytesInLengthHeader returns the number of bytes in the ASn1 header that indicate the length.
func GetNumberBytesInLengthHeader(b []byte) int {
	if int(b[1]) <= 127 {
		return 1
	}
	// The bytes that indicate the length
	return 1 + int(b[1]) - 128
}

// AddASNAppTag adds an ASN1 encoding application tag value to the raw bytes provided.
func AddASNAppTag(b []byte, tag int) []byte {
	r := RawValue{
		Class:      ClassApplication,
		IsCompound: true,
		Tag:        tag,
		Bytes:      b,
	}
	ab, _ := Marshal(r)
	return ab
}
