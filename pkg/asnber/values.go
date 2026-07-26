package asnber

type SnmpValue interface {
	GetType() BERType
	GetBytes() []byte
	GetEncoded() []byte
}

type SnmpOpaqueValue struct {
	Type  BERType
	Value []byte
}

func (s SnmpOpaqueValue) GetType() BERType {
	return s.Type
}

func (s SnmpOpaqueValue) GetBytes() []byte {
	return s.Value
}

func (s SnmpOpaqueValue) GetEncoded() []byte {
	return EncodeSnmp(s.Type, s.Value)
}

func EncodeSnmpValue(sv SnmpValue) []byte {
	return EncodeSnmp(sv.GetType(), sv.GetBytes())
}

func EncodeSnmp(t BERType, v []byte) []byte {
	_l := len(v)
	_ret := make([]byte, 1)
	_ret[0] = byte(t)
	_len := EncodeLength(_l)
	_ret = append(_ret, _len...)
	_ret = append(_ret, v...)
	return _ret
}

type SnmpOctetValue struct {
	Type  BERType
	Value []byte
}

func (s SnmpOctetValue) GetType() BERType {
	return s.Type
}

func (s SnmpOctetValue) GetBytes() []byte {
	return s.Value
}

func (s SnmpOctetValue) GetEncoded() []byte {
	return EncodeSnmp(s.Type, s.Value)
}

type SnmpStringValue struct {
	Type  BERType
	Value string
}

func (s SnmpStringValue) GetType() BERType {
	return s.Type
}

func (s SnmpStringValue) GetBytes() []byte {
	return []byte(s.Value)
}

func (s SnmpStringValue) GetEncoded() []byte {
	return EncodeSnmp(s.Type, s.GetBytes())
}

type SnmpIntValue struct {
	Type  BERType
	Value int
}

func (s SnmpIntValue) GetType() BERType {
	return s.Type
}

func (s SnmpIntValue) GetBytes() []byte {
	return EncodeInteger32(s.Value)
}

func (s SnmpIntValue) GetEncoded() []byte {
	return EncodeSnmp(s.Type, s.GetBytes())
}

type SnmpVarBind struct {
	Oid   string
	Value SnmpValue
}

func (s SnmpVarBind) GetType() BERType {
	return BERType(DerSequence)
}

func (s SnmpVarBind) GetBytes() []byte {
	_oid, _ := ParseOid(s.Oid)
	_b, _ := _oid.Encode()
	_ret := make([]byte, 0)
	_ret = append(_ret, EncodeSnmp(AsnObjectID, _b)...)
	_ret = append(_ret, EncodeSnmpValue(s.Value)...)
	return _ret
}

func (s SnmpVarBind) GetEncoded() []byte {
	return EncodeSnmp(s.GetType(), s.GetBytes())
}

type SnmpVarBinds struct {
	Type   BERType
	Values []SnmpVarBind
}

func (s SnmpVarBinds) GetType() BERType {
	return s.Type
}

func (s SnmpVarBinds) GetBytes() []byte {
	_ret := make([]byte, 0)
	for _, v := range s.Values {
		_ret = append(_ret, v.GetEncoded()...)
	}
	return _ret
}

func (s SnmpVarBinds) GetEncoded() []byte {
	return EncodeSnmp(s.GetType(), s.GetBytes())
}

type SnmpStruct struct {
	Type   BERType
	Values []SnmpValue
}

func (s SnmpStruct) GetType() BERType {
	return s.Type
}

func (s SnmpStruct) GetBytes() []byte {
	_ret := make([]byte, 0)
	for _, v := range s.Values {
		_ret = append(_ret, v.GetEncoded()...)
	}
	return _ret
}

func (s SnmpStruct) GetEncoded() []byte {
	return EncodeSnmp(s.GetType(), s.GetBytes())
}

type SnmpTrapValue struct {
	SourceAddress   string
	Version         int
	CommunityOrUser string
	Type            BERType
	Enterprise      int
	AgentAddress    string
	Generic         int
	Specific        int
	TimeStamp       int
	Varbinds        []SnmpVarBind
	Extras          []SnmpVarBind
}
