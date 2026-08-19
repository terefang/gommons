package xstrings

// Common string constants used throughout the stringx package.
// These constants provide commonly used string literals for better code readability and maintainability.
const (
	Ampersand       = "&"
	And             = "and"
	Or              = "or"
	At              = "@"
	Asterisk        = "*"
	Star            = Asterisk
	Backslash       = "\\"
	Colon           = ":"
	Comma           = ","
	Dash            = "-"
	Dollar          = "$"
	Dot             = "."
	DotDot          = ".."
	Empty           = ""
	Equals          = "="
	False           = "false"
	Slash           = "/"
	Hash            = "#"
	Hat             = "^"
	LeftBrace       = "{"
	LeftBracket     = "("
	LeftChev        = "<"
	Newline         = "\n"
	No              = "no"
	Null            = "null"
	Off             = "off"
	On              = "on"
	Percent         = "%"
	Pipe            = "|"
	Plus            = "+"
	QuestionMark    = "?"
	ExclamationMark = "!"
	Quote           = "\""
	Return          = "\r"
	Tab             = "\t"
	RightBrace      = "}"
	RightBracket    = ")"
	RightChev       = ">"
	Semicolon       = ";"
	SingleQuote     = "'"
	Backtick        = "`"
	Space           = " "
	Tilda           = "~"
	LeftSqBracket   = "["
	RightSqBracket  = "]"
	True            = "true"
	Underscore      = "_"
	Yes             = "yes"
	In              = "in"
	Crlf            = "\r\n"

	// Additional punctuation and symbols
	Exclamation     = "!"
	AtSign          = "@"
	HashTag         = "#"
	DollarSign      = "$"
	PercentSign     = "%"
	Caret           = "^"
	AmpersandSign   = "&"
	StarSign        = "*"
	PlusSign        = "+"
	MinusSign       = "-"
	EqualsSign      = "="
	UnderscoreSign  = "_"
	PipeSign        = "|"
	BackslashSign   = "\\"
	ForwardSlash    = "/"
	ColonSign       = ":"
	SemicolonSign   = ";"
	CommaSign       = ","
	DotSign         = "."
	QuestionSign    = "?"
	ExclamationSign = "!"

	// Bracket pairs (additional)
	LeftParen  = "("
	RightParen = ")"
	LeftAngle  = "<"
	RightAngle = ">"

	// Quote types (additional)
	DoubleQuote   = "\""
	BacktickQuote = "`"

	// Whitespace characters
	SpaceChar      = " "
	TabChar        = "\t"
	NewlineChar    = "\n"
	CarriageReturn = "\r"
	FormFeed       = "\f"
	VerticalTab    = "\v"

	// Common separators
	CommaSpace     = ", "
	SemicolonSpace = "; "
	ColonSpace     = ": "
	PipeSpace      = " | "
	SlashSpace     = " / "
	BackslashSpace = " \\ "

	// Common boolean values
	BooleanTrue     = "true"
	BooleanFalse    = "false"
	BooleanYes      = "yes"
	BooleanNo       = "no"
	BooleanOn       = "on"
	BooleanOff      = "off"
	BooleanEnabled  = "enabled"
	BooleanDisabled = "disabled"

	// Common null/empty values
	NullValue      = "null"
	UndefinedValue = "undefined"
	EmptyString    = ""
	ZeroValue      = "0"
	OneValue       = "1"

	// Common file extensions
	ExtJSONL  = ".jsonl"
	ExtJSON   = ".json"
	ExtXML    = ".xml"
	ExtHTML   = ".html"
	ExtCSS    = ".css"
	ExtJS     = ".js"
	ExtGo     = ".go"
	ExtTxt    = ".txt"
	ExtLog    = ".log"
	ExtYAML   = ".yaml"
	ExtYML    = ".yml"
	ExtTOML   = ".toml"
	ExtINI    = ".ini"
	ExtPdata  = ".pdata"
	ExtConf   = ".conf"
	ExtConfig = ".config"

	// Common protocols
	ProtocolHTTP            = "http"
	ProtocolHTTPS           = "https"
	ProtocolFTP             = "ftp"
	ProtocolFTPS            = "ftps"
	ProtocolSFTP            = "sftp"
	ProtocolSSH             = "ssh"
	ProtocolTCP             = "tcp"
	ProtocolUDP             = "udp"
	ProtocolWebSocket       = "ws"
	ProtocolWebSocketSecure = "wss"

	// Common encoding types
	EncodingUTF8   = "utf-8"
	EncodingUTF16  = "utf-16"
	EncodingASCII  = "ascii"
	EncodingBase64 = "base64"
	EncodingBase32 = "base32"
	EncodingBase36 = "base36"
	EncodingHex    = "hex"
	EncodingURL    = "url"

	// Common time zones
	TimezoneUTC = "UTC"
	TimezoneGMT = "GMT"
	TimezoneEST = "EST"
	TimezonePST = "PST"
	TimezoneCST = "CST"
	TimezoneMST = "MST"

	// Common units
	UnitBytes = "bytes"
	UnitKB    = "KB"
	UnitMB    = "MB"
	UnitGB    = "GB"
	UnitTB    = "TB"
	UnitPB    = "PB"
)

const Numbers = "0123456789"
const LowercaseLetters = "abcdefghijklmnopqrstuvwxyz"
const UppercaseLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const Letters = LowercaseLetters + UppercaseLetters
const Alphanumeric = Letters + Numbers

const HexStr = Numbers + "abcdef"
const HexUpper = Numbers + "ABCDEF"

const Base36Str = Numbers + LowercaseLetters
const Base36Upper = Numbers + UppercaseLetters

const Symbols = SimpleSymbols + Quotes + MathSymbols + Brackets + Punctuation
const SimpleSymbols = `@#$&_|`
const Quotes = `"'`
const MathSymbols = `%^*+-=/`
const Punctuation = `.,!?;:`
const Brackets = `()[]{}<>`
const AlphanumericSymbols = Letters + Numbers + Symbols

const CommonFieldSeparators = ` ,;:/|`
