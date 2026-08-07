package xradius

type Attribute struct {
    Type       int
    Vendor     int
    VendorType int
    Value      []byte
}

type AttributeDefinition struct {
    Type       int
    Vendor     int
    VendorType int
    Name       string
    Encoding   string
    Encrypted  bool
}

var RadiusAttributeUserName = AttributeDefinition{1, 0, 0, "UserName", "string", false}
var RadiusAttributeUserPassword = AttributeDefinition{2, 0, 0, "UserPassword", "string", true}
var RadiusAttributeCHAPPassword = AttributeDefinition{3, 0, 0, "CHAPPassword", "octets", false}
var RadiusAttributeNASIPAddress = AttributeDefinition{4, 0, 0, "NASIPAddress", "ipaddr", false}
var RadiusAttributeNASPort = AttributeDefinition{5, 0, 0, "NASPort", "uint32", false}
var RadiusAttributeServiceType = AttributeDefinition{6, 0, 0, "ServiceType", "uint32", false}
var RadiusAttributeFramedProtocol = AttributeDefinition{7, 0, 0, "FramedProtocol", "uint32", false}
var RadiusAttributeFramedIPAddress = AttributeDefinition{8, 0, 0, "FramedIPAddress", "ipaddr", false}
var RadiusAttributeFramedIPNetmask = AttributeDefinition{9, 0, 0, "FramedIPNetmask", "ipaddr", false}
var RadiusAttributeFramedRouting = AttributeDefinition{10, 0, 0, "FramedRouting", "uint32", false}
var RadiusAttributeFilterId = AttributeDefinition{11, 0, 0, "FilterId", "string", false}
var RadiusAttributeFramedMTU = AttributeDefinition{12, 0, 0, "FramedMTU", "uint32", false}
var RadiusAttributeFramedCompression = AttributeDefinition{13, 0, 0, "FramedCompression", "uint32", false}
var RadiusAttributeLoginIPHost = AttributeDefinition{14, 0, 0, "LoginIPHost", "ipaddr", false}
var RadiusAttributeLoginService = AttributeDefinition{15, 0, 0, "LoginService", "uint32", false}
var RadiusAttributeLoginTCPPort = AttributeDefinition{16, 0, 0, "LoginTCPPort", "uint32", false}
var RadiusAttributeReplyMessage = AttributeDefinition{18, 0, 0, "ReplyMessage", "string", false}
var RadiusAttributeCallbackNumber = AttributeDefinition{19, 0, 0, "CallbackNumber", "string", false}
var RadiusAttributeCallbackId = AttributeDefinition{20, 0, 0, "CallbackId", "string", false}
var RadiusAttributeFramedRoute = AttributeDefinition{22, 0, 0, "FramedRoute", "string", false}
var RadiusAttributeFramedIPXNetwork = AttributeDefinition{23, 0, 0, "FramedIPXNetwork", "ipaddr", false}
var RadiusAttributeState = AttributeDefinition{24, 0, 0, "State", "octets", false}
var RadiusAttributeClass = AttributeDefinition{25, 0, 0, "Class", "octets", false}
var RadiusAttributeVendorSpecific = AttributeDefinition{26, 0, 0, "VendorSpecific", "vsa", false}
var RadiusAttributeSessionTimeout = AttributeDefinition{27, 0, 0, "SessionTimeout", "uint32", false}
var RadiusAttributeIdleTimeout = AttributeDefinition{28, 0, 0, "IdleTimeout", "uint32", false}
var RadiusAttributeTerminationAction = AttributeDefinition{29, 0, 0, "TerminationAction", "uint32", false}
var RadiusAttributeCalledStationId = AttributeDefinition{30, 0, 0, "CalledStationId", "string", false}
var RadiusAttributeCallingStationId = AttributeDefinition{31, 0, 0, "CallingStationId", "string", false}
var RadiusAttributeNASIdentifier = AttributeDefinition{32, 0, 0, "NASIdentifier", "string", false}
var RadiusAttributeProxyState = AttributeDefinition{33, 0, 0, "ProxyState", "octets", false}
var RadiusAttributeLoginLATService = AttributeDefinition{34, 0, 0, "LoginLATService", "string", false}
var RadiusAttributeLoginLATNode = AttributeDefinition{35, 0, 0, "LoginLATNode", "string", false}
var RadiusAttributeLoginLATGroup = AttributeDefinition{36, 0, 0, "LoginLATGroup", "octets", false}
var RadiusAttributeFramedAppleTalkLink = AttributeDefinition{37, 0, 0, "FramedAppleTalkLink", "uint32", false}
var RadiusAttributeFramedAppleTalkNetwork = AttributeDefinition{38, 0, 0, "FramedAppleTalkNetwork", "uint32", false}
var RadiusAttributeFramedAppleTalkZone = AttributeDefinition{39, 0, 0, "FramedAppleTalkZone", "string", false}
var RadiusAttributeCHAPChallenge = AttributeDefinition{60, 0, 0, "CHAPChallenge", "octets", false}
var RadiusAttributeNASPortType = AttributeDefinition{61, 0, 0, "NASPortType", "uint32", false}
var RadiusAttributePortLimit = AttributeDefinition{62, 0, 0, "PortLimit", "uint32", false}
var RadiusAttributeLoginLATPort = AttributeDefinition{63, 0, 0, "LoginLATPort", "string", false}

var RadiusAttributes = []AttributeDefinition{
    RadiusAttributeUserName,
    RadiusAttributeUserPassword,
    RadiusAttributeCHAPPassword,
    RadiusAttributeNASIPAddress,
    RadiusAttributeNASPort,
    RadiusAttributeServiceType,
    RadiusAttributeFramedProtocol,
    RadiusAttributeFramedIPAddress,
    RadiusAttributeFramedIPNetmask,
    RadiusAttributeFramedRouting,
    RadiusAttributeFilterId,
    RadiusAttributeFramedMTU,
    RadiusAttributeFramedCompression,
    RadiusAttributeLoginIPHost,
    RadiusAttributeLoginService,
    RadiusAttributeLoginTCPPort,
    RadiusAttributeReplyMessage,
    RadiusAttributeCallbackNumber,
    RadiusAttributeCallbackId,
    RadiusAttributeFramedRoute,
    RadiusAttributeFramedIPXNetwork,
    RadiusAttributeState,
    RadiusAttributeClass,
    RadiusAttributeVendorSpecific,
    RadiusAttributeSessionTimeout,
    RadiusAttributeIdleTimeout,
    RadiusAttributeTerminationAction,
    RadiusAttributeCalledStationId,
    RadiusAttributeCallingStationId,
    RadiusAttributeNASIdentifier,
    RadiusAttributeProxyState,
    RadiusAttributeLoginLATService,
    RadiusAttributeLoginLATNode,
    RadiusAttributeLoginLATGroup,
    RadiusAttributeFramedAppleTalkLink,
    RadiusAttributeFramedAppleTalkNetwork,
    RadiusAttributeFramedAppleTalkZone,
    RadiusAttributeCHAPChallenge,
    RadiusAttributeNASPortType,
    RadiusAttributePortLimit,
    RadiusAttributeLoginLATPort,
}
