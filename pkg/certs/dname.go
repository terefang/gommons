package certs

import (
	"crypto/x509/pkix"
	"encoding/asn1"
	"strings"
)

var oids = map[string]asn1.ObjectIdentifier{
	"businesscategory":           {2, 5, 4, 15},
	"c":                          {2, 5, 4, 6},
	"cn":                         {2, 5, 4, 3},
	"dc":                         {0, 9, 2342, 19200300, 100, 1, 25},
	"description":                {2, 5, 4, 13},
	"destinationindicator":       {2, 5, 4, 27},
	"distinguishedName":          {2, 5, 4, 49},
	"dnqualifier":                {2, 5, 4, 46},
	"emailaddress":               {1, 2, 840, 113549, 1, 9, 1},
	"enhancedsearchguide":        {2, 5, 4, 47},
	"facsimiletelephonenumber":   {2, 5, 4, 23},
	"generationqualifier":        {2, 5, 4, 44},
	"givenname":                  {2, 5, 4, 42},
	"houseidentifier":            {2, 5, 4, 51},
	"initials":                   {2, 5, 4, 43},
	"internationalisdnnumber":    {2, 5, 4, 25},
	"l":                          {2, 5, 4, 7},
	"member":                     {2, 5, 4, 31},
	"name":                       {2, 5, 4, 41},
	"o":                          {2, 5, 4, 10},
	"ou":                         {2, 5, 4, 11},
	"owner":                      {2, 5, 4, 32},
	"physicaldeliveryofficename": {2, 5, 4, 19},
	"postaladdress":              {2, 5, 4, 16},
	"postalcode":                 {2, 5, 4, 17},
	"postOfficebox":              {2, 5, 4, 18},
	"preferreddeliverymethod":    {2, 5, 4, 28},
	"registeredaddress":          {2, 5, 4, 26},
	"roleoccupant":               {2, 5, 4, 33},
	"searchguide":                {2, 5, 4, 14},
	"seealso":                    {2, 5, 4, 34},
	"serialnumber":               {2, 5, 4, 5},
	"sn":                         {2, 5, 4, 4},
	"st":                         {2, 5, 4, 8},
	"street":                     {2, 5, 4, 9},
	"telephonenumber":            {2, 5, 4, 20},
	"teletexterminalidentifier":  {2, 5, 4, 22},
	"telexnumber":                {2, 5, 4, 21},
	"title":                      {2, 5, 4, 12},
	"uid":                        {0, 9, 2342, 19200300, 100, 1, 1},
	"uniquemember":               {2, 5, 4, 50},
	"userpassword":               {2, 5, 4, 35},
	"x121address":                {2, 5, 4, 24},
}

// Aliases for the attribute names.
// e.g. from https://learn.microsoft.com/en-us/windows/win32/seccrypto/name-properties
var alias = map[string]string{
	"e": "emailaddress",
}

// ParseSimpleDN returns a distinguishedName or an error.
// The function respects somewhat https://tools.ietf.org/html/rfc4514
func ParseSimpleDN(str string) (*pkix.Name, error) {
	dn := make(pkix.RelativeDistinguishedNameSET, 0)

	_parts := strings.Split(str, ",")
	for _, _part := range _parts {
		_attr := strings.SplitN(_part, "=", 2)
		_attr[0] = strings.ToLower(_attr[0])
		if aliasName, ok := alias[_attr[0]]; ok {
			_attr[0] = aliasName
		}
		rdn := new(pkix.AttributeTypeAndValue)
		rdn.Type = oids[_attr[0]]
		rdn.Value = _attr[1]
		dn = append(dn, *rdn)
	}

	var sequence = pkix.RDNSequence{dn}
	var name pkix.Name
	//name.CommonName
	name.FillFromRDNSequence(&sequence)
	err := fillExtraNames(&sequence, &name)
	if err != nil {
		return nil, err
	}
	return &name, nil
}

// Fill in ExtraNames with RDNs with OID prefix other than 2.5.4
func fillExtraNames(rdns *pkix.RDNSequence, name *pkix.Name) error {
	for _, rdn := range *rdns {
		if len(rdn) == 0 {
			continue
		}

		for _, atv := range rdn {
			if atv.Type.Equal(oids["dc"]) || atv.Type.Equal(oids["emailaddress"]) {
				// IA5String
				atv.Value = asn1.RawValue{Tag: 22, Class: 0, Bytes: []byte(atv.Value.(string))}
				name.ExtraNames = append(name.ExtraNames, atv)
			}
		}
	}
	return nil
}
