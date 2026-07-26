/*
 *  Licensed to the Apache Software Foundation (ASF) under one or more
 *  contributor license agreements.  See the NOTICE file distributed with
 *  this work for additional information regarding copyright ownership.
 *  The ASF licenses this file to You under the Apache License, Version 2.0
 *  (the "License"); you may not use this file except in compliance with
 *  the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 * Original Code from github.com/apache/harmony/.../StreamTokenizer.java
 *
 * ported to golang — Copyright (c) 2026. terefang@gmail.com
 */

package ctokener

import (
	"bufio"
	"encoding/hex"
	"io"
	"strconv"
	"strings"
	"unicode"
)

type Token int

func (t *Token) ToString() string {

	if t == nil {
		return "NIL"
	}
	if *t == TOKEN_TYPE_UNKNOWN {
		return "UNKNOWN"
	}
	if *t == TOKEN_TYPE_EOF {
		return "EOF"
	}
	if *t == TOKEN_TYPE_EOL {
		return "EOL"
	}
	if *t == TOKEN_TYPE_NUMBER {
		return "NUMBER"
	}
	if *t == TOKEN_TYPE_CARDINAL {
		return "CARDINAL"
	}
	if *t == TOKEN_TYPE_IDENTIFIER {
		return "IDENTIFIER"
	}
	if *t == TOKEN_TYPE_STRING {
		return "STRING"
	}
	if *t == TOKEN_TYPE_STRING_MULTILINE {
		return "STRING_MULTILINE"
	}
	if *t == TOKEN_TYPE_BYTES {
		return "BYTES"
	}
	if *t == TOKEN_TYPE_BOOL {
		return "BOOL"
	}
	if *t == TOKEN_TYPE_COMMENT_SINGLE_LINE {
		return "COMMENT_SINGLE_LINE"
	}
	if *t == TOKEN_TYPE_COMMENT_SLASH_STAR {
		return "COMMENT_SLASH_STAR"
	}
	if *t == TOKEN_TYPE_COMMENT_SLASH_SLASH {
		return "COMMENT_SLASH_SLASH"
	}
	return string(rune(*t))
}

const (
	// token constants
	TOKEN_TYPE_UNKNOWN             = 0
	TOKEN_TYPE_EOF                 = -1
	TOKEN_TYPE_EOL                 = 10
	TOKEN_TYPE_NUMBER              = -2
	TOKEN_TYPE_CARDINAL            = -3
	TOKEN_TYPE_IDENTIFIER          = -4
	TOKEN_TYPE_STRING              = -5
	TOKEN_TYPE_BYTES               = -6
	TOKEN_TYPE_BOOL                = -7
	TOKEN_TYPE_COMMENT_SINGLE_LINE = -8
	TOKEN_TYPE_COMMENT_SLASH_SLASH = -9
	TOKEN_TYPE_COMMENT_SLASH_STAR  = -10
	TOKEN_TYPE_STRING_MULTILINE    = -11

	// lexical constants
	LEXICAL_UNKNOWN         = 0
	LEXICAL_COMMENT         = 1
	LEXICAL_QUOTE           = 2
	LEXICAL_WHITE           = 4
	LEXICAL_WORD            = 8
	LEXICAL_DIGIT           = 16
	LEXICAL_BYTES           = 32
	LEXICAL_QUOTE_MULTILINE = 64

	// other constants
	HEX_SET = "0123456789abcdefABCDEF"
)

type Position struct {
	line   int
	column int
}

func (p *Position) Line() int { return p.line }
func (p *Position) Col() int  { return p.column }

/**
 * Parses a stream into a set of defined tokens, one at a time. The different
 * types of tokens that can be found are numbers, identifiers, quoted strings,
 * and different comment styles. The class can be used for limited processing
 * of source code of programming languages like Java, although it is nowhere
 * near a full parser.
 */

type CustomTokener struct {
	// internal
	reader *bufio.Reader

	// parse options
	LexicalTypes               []int
	ForceIndentifiersLowercase bool
	IsEolSignificant           bool
	IsCommentSignificant       bool
	IsSlashStarComments        bool
	IsSlashSlashComments       bool
	IsBacktickQuote            bool
	IsHexLiterals              bool
	IsByteLiterals             bool
	IsAutoUnicodeMode          bool
	IsPushBackToken            bool

	// state
	peekChar rune
	lastCR   bool
	Pos      Position
	// contains the TokenType as returned by the NextToken method
	TokenType Token
	// this value contains the single character for unknown token
	CharacterValue rune
	NumericValue   float64
	CardinalValue  int64
	StringValue    string
	ByteValue      []byte
	BooleanValue   bool
}

// create a new tokener/lexer from reader
func NewTokener(reader io.Reader) *CustomTokener {
	return &CustomTokener{
		Pos:          Position{line: 1, column: 0},
		reader:       bufio.NewReader(reader),
		peekChar:     -2,
		LexicalTypes: make([]int, 256),
	}
}

func (ct *CustomTokener) ConfigureSimpleDefaults() *CustomTokener {
	ct.SetWordChars('A', 'Z')
	ct.SetWordChars('a', 'z')
	ct.SetWordChars(160, 255)
	ct.SetWsChars(0, 32)
	ct.SetCommentChar('/')
	ct.SetQuoteChar('"')
	ct.SetQuoteChar('\'')
	ct.SetParseNumbers()
	return ct
}

func (ct *CustomTokener) ConfigureLdataSyntax() *CustomTokener {
	ct.ResetSyntax()
	ct.SetQuoteChar('"')
	ct.SetBackTickQuotes(true)
	ct.SetWsChars(0, 32)
	ct.SetWsChar(',')
	ct.SetWsChar(';')
	ct.SetWsChar(':')
	ct.SetWsChar('=')
	ct.SetWordChars('a', 'z')
	ct.SetWordChars('A', 'Z')
	ct.SetWordChars('0', '9')
	ct.SetWordChar('_')
	ct.SetWordChar('-')
	ct.SetWordChar('@')
	ct.SetParseNumbers()
	ct.SetHexLiterals(true)
	ct.SetSlashStarComments(true)
	ct.SetSlashSlashComments(true)
	ct.SetCommentChar('#')
	ct.SetCommentChar('!')
	ct.SetCommentChar('%')
	ct.SetAutoUnicodeMode(true)
	ct.SetCommentSignificant(false)
	ct.SetEolSignificant(false)
	return ct
}

// Specifies that all characters shall be treated as ordinary
func (ct *CustomTokener) ResetSyntax() *CustomTokener {
	ct.ResetChars(0, 255)
	return ct
}

// Specifies that character range shall be treated as ordinary
func (ct *CustomTokener) ResetChars(startChar rune, endChar rune) *CustomTokener {
	for n := startChar; n <= endChar; n++ {
		ct.ResetChar(n)
	}
	return ct
}

// Specifies that character shall be treated as ordinary.
func (ct *CustomTokener) ResetChar(n rune) *CustomTokener {
	ct.LexicalTypes[int(n)] = LEXICAL_UNKNOWN
	return ct
}

// Specifies that character range shall be treated as custom
func (ct *CustomTokener) SetCustomChars(x int, startChar rune, endChar rune) *CustomTokener {
	for n := startChar; n <= endChar; n++ {
		ct.SetCustomChar(x, n)
	}
	return ct
}

// Specifies that character shall be treated as custom
func (ct *CustomTokener) SetCustomChar(x int, n rune) *CustomTokener {
	ct.LexicalTypes[int(n)] = -x
	return ct
}

// Specifies that character range shall be treated as word
func (ct *CustomTokener) SetWordChars(startChar rune, endChar rune) *CustomTokener {
	for n := startChar; n <= endChar; n++ {
		ct.SetWordChar(n)
	}
	return ct
}

// Specifies that character shall be treated as word
func (ct *CustomTokener) SetWordChar(n rune) *CustomTokener {
	ct.LexicalTypes[int(n)] = LEXICAL_WORD
	return ct
}

// Specifies that character range shall be treated as white-space
func (ct *CustomTokener) SetWsChars(startChar rune, endChar rune) *CustomTokener {
	for n := startChar; n <= endChar; n++ {
		ct.SetWsChar(n)
	}
	return ct
}

// Specifies that character shall be treated as white-space
func (ct *CustomTokener) SetWsChar(n rune) *CustomTokener {
	ct.LexicalTypes[int(n)] = LEXICAL_WHITE
	return ct
}

// Specifies that character shall be treated as single-line comment starter
func (ct *CustomTokener) SetCommentChar(n rune) *CustomTokener {
	ct.LexicalTypes[int(n)] = LEXICAL_COMMENT
	return ct
}

// Specifies that character shall be treated as string quote
func (ct *CustomTokener) SetQuoteChar(n rune) *CustomTokener {
	ct.LexicalTypes[int(n)] = LEXICAL_QUOTE
	return ct
}

// Specifies if c-style slash-slash comments should be recognized
func (ct *CustomTokener) SetSlashSlashComments(n bool) *CustomTokener {
	ct.IsSlashSlashComments = n
	return ct
}

// Specifies if c-style slash-star comments should be recognized
func (ct *CustomTokener) SetSlashStarComments(n bool) *CustomTokener {
	ct.IsSlashStarComments = n
	return ct
}

// Specifies if backticks be used for multi-line quotes
func (ct *CustomTokener) SetBackTickQuotes(n bool) *CustomTokener {
	ct.IsBacktickQuote = n
	return ct
}

// Specifies if hex/binary literals (0x../0b..) should be recognized
func (ct *CustomTokener) SetHexLiterals(n bool) *CustomTokener {
	ct.IsHexLiterals = n
	return ct
}

// Specifies if byte literals (<HH>) should be recognized
func (ct *CustomTokener) SetByteLiterals(n bool) *CustomTokener {
	ct.IsByteLiterals = n
	return ct
}

// Specifies that this tokenizer shall parse numbers.
func (ct *CustomTokener) SetParseNumbers() *CustomTokener {
	for n := '0'; n <= '9'; n++ {
		ct.LexicalTypes[int(n)] |= LEXICAL_DIGIT
	}
	ct.LexicalTypes[int('.')] |= LEXICAL_DIGIT
	ct.LexicalTypes[int('-')] |= LEXICAL_DIGIT
	return ct
}

// specifies that EOL is returned as Token
func (ct *CustomTokener) SetEolSignificant(n bool) *CustomTokener {
	ct.IsEolSignificant = n
	return ct
}

// specifies that COMMENT is returned as Token
func (ct *CustomTokener) SetCommentSignificant(n bool) *CustomTokener {
	ct.IsCommentSignificant = n
	return ct
}

// specifies that returned IDENTIFIERS is returned in locercase
func (ct *CustomTokener) SetLowercaseMode(n bool) *CustomTokener {
	ct.ForceIndentifiersLowercase = n
	return ct
}

// specifies that unicode runes are automatically recognized (IsLetter/IsDigit/IsSpace/IsControl)
func (ct *CustomTokener) SetAutoUnicodeMode(n bool) *CustomTokener {
	ct.IsAutoUnicodeMode = n
	return ct
}

func (ct *CustomTokener) checkCharType(_c rune, _start bool) int {
	if _start && ct.IsByteLiterals && (_c == '<') {
		return LEXICAL_BYTES
	}

	if _start && ct.IsBacktickQuote && (_c == '`') {
		return LEXICAL_QUOTE_MULTILINE
	}

	if _c < 256 {
		return ct.LexicalTypes[int(_c)]
	}

	if ct.IsAutoUnicodeMode {
		if unicode.IsLetter(_c) {
			return LEXICAL_WORD
		}
		if unicode.IsDigit(_c) {
			return LEXICAL_DIGIT
		}
		if unicode.IsSpace(_c) {
			return LEXICAL_WHITE
		}
		if unicode.IsControl(_c) {
			return LEXICAL_WHITE
		}
	}

	return LEXICAL_UNKNOWN
}

// reads from stream also counting lines and columns
func (ct *CustomTokener) read() (rune, error) {
	_c, _sz, _err := ct.reader.ReadRune()
	ct.reportNextChar(_sz) // TODO
	if _c == '\n' {
		ct.reportNextLine()
	}
	if _c == '\r' {
		ct.reportNextChar(-1)
	}
	return _c, _err
}

// Indicates that the current token should be pushed back and returned again
// the next time NextToken() is called.
func (ct *CustomTokener) PushBack() {
	ct.IsPushBackToken = true
}

func (ct *CustomTokener) readXQuoted(_qqq string, _qs *strings.Builder) (int, error) {
	for !strings.HasSuffix(_qs.String(), _qqq) {
		_c, err := ct.read()
		if err != nil {
			return 0, err
		}
		if _c == '\\' {
			_c, err = ct.readEscaped()
			if err != nil {
				return 0, err
			}
		}
		_qs.WriteRune(_c)
	}
	str := _qs.String()
	_qs.Reset()
	_qs.WriteString(str[:len(str)-len(_qqq)])

	// text blocks remove all the incidental indentations and keep only essential indentations.
	return TOKEN_TYPE_STRING_MULTILINE, nil
}

func (ct *CustomTokener) readNumberLiteral(_c rune) (Token, error) {
	_sb := strings.Builder{}
	haveDecimal := false
	checkJustNegative := _c == '-'
	isHex := false
	isBinary := false
	for true {
		if ct.IsHexLiterals && _sb.Len() == 0 && _c == '0' {
			// should we detect hex-literals starting with '0x' and switch parsing ?
			// actually we can discard the zero here
			_p, err := ct.read()
			if err != nil {
				return 0, err
			}

			if _p == 'x' {
				isHex = true
				_p, err = ct.read()
				if err != nil {
					return 0, err
				}
				if strings.Index(HEX_SET, string(_p)) == -1 {
					break
				}
			} else if _p == 'b' {
				isBinary = true
				_p, err = ct.read()
				if err != nil {
					return 0, err
				}
				if _p != '0' && _p != '1' {
					break
				}
			} else {
				// alternative would be to have a float-literal starting with '0.'.
				if _p == '.' {
					haveDecimal = true
				} else if _p < '0' || _p > '9' {
					_sb.WriteRune('0')
					_c = _p
					break
				}
			}
			_c = _p
		}

		if isBinary {
			if _c != '_' {
				_sb.WriteRune(_c)
			}
			_p, err := ct.read()
			if err != nil {
				return 0, err
			}
			_c = _p
			if _c != '0' && _c != '1' && _c != '_' {
				break
			}
		} else if isHex {
			if _c != '_' {
				_sb.WriteRune(_c)
			}
			_p, err := ct.read()
			if err != nil {
				return 0, err
			}
			_c = _p
			if !(strings.Index(HEX_SET, string(_c)) != -1 || _c == '_') {
				break
			}
		} else {
			if _c == '.' {
				haveDecimal = true
			}
			if _c != '_' {
				_sb.WriteRune(_c)
			}
			_p, err := ct.read()
			if err != nil {
				return 0, err
			}
			_c = _p
			if (_c < '0' || _c > '9') && (haveDecimal || _c != '.') && _c != '_' {
				break
			}
		}
	}
	ct.peekChar = _c
	if checkJustNegative && _sb.Len() == 1 {
		ct.TokenType = TOKEN_TYPE_UNKNOWN
		ct.CharacterValue = '-'
		return ct.TokenType, nil
	}

	ct.CardinalValue = 0
	ct.NumericValue = 0

	ct.StringValue = _sb.String()
	if isHex {
		_cv, err := strconv.ParseUint(ct.StringValue, 16, 64)
		if err != nil {
			return 0, err
		}
		ct.CardinalValue = int64(_cv)
	} else if isBinary {
		_cv, err := strconv.ParseUint(ct.StringValue, 2, 64)
		if err != nil {
			return 0, err
		}
		ct.CardinalValue = int64(_cv)
	} else if haveDecimal {
		_cv, err := strconv.ParseFloat(ct.StringValue, 64)
		if err != nil {
			return 0, err
		}
		ct.NumericValue = _cv
	} else {
		_cv, err := strconv.ParseUint(ct.StringValue, 10, 64)
		if err != nil {
			return 0, err
		}
		ct.CardinalValue = int64(_cv)
	}

	if haveDecimal {
		ct.TokenType = TOKEN_TYPE_NUMBER
	} else {
		ct.TokenType = TOKEN_TYPE_CARDINAL
	}
	return ct.TokenType, nil
}

func (ct *CustomTokener) readWordLiteral(_c rune) (Token, error) {
	_sb := strings.Builder{}
	for true {
		_sb.WriteRune(_c)
		_peek, err := ct.read()
		if err != nil {
			return 0, err
		}
		_c = _peek
		if _peek < 256 && ((ct.LexicalTypes[int(_peek)] & (LEXICAL_WORD | LEXICAL_DIGIT)) == 0) {
			ct.peekChar = _peek
			break
		}
	}
	ct.StringValue = _sb.String()
	if ct.ForceIndentifiersLowercase {
		ct.StringValue = strings.ToLower(ct.StringValue)
	}
	ct.TokenType = TOKEN_TYPE_IDENTIFIER
	return ct.TokenType, nil
}

func (ct *CustomTokener) eatComments(_c rune) (Token, error) {
	_sb := strings.Builder{}
	for _c != '\r' && _c != '\n' {
		_sb.WriteRune(_c)
		_peek, err := ct.read()
		if err != nil {
			return 0, err
		}
		_c = _peek
	}
	ct.peekChar = _c
	if !ct.IsCommentSignificant {
		return ct.NextToken()
	}
	ct.StringValue = _sb.String()
	ct.TokenType = TOKEN_TYPE_COMMENT_SINGLE_LINE
	return ct.TokenType, nil
}

func (ct *CustomTokener) eatSlashSlashComments() (Token, error) {
	_sb := strings.Builder{}
	for true {
		peekOne, err := ct.read()
		if err != nil {
			return 0, err
		}
		_sb.WriteRune(peekOne)
		if peekOne == '\r' || peekOne == '\n' {
			break
		}
	}
	if !ct.IsCommentSignificant {
		return ct.NextToken()
	}
	ct.StringValue = _sb.String()
	ct.TokenType = TOKEN_TYPE_COMMENT_SLASH_SLASH
	return ct.TokenType, nil
}

func (ct *CustomTokener) eatSlashStarComments() (Token, error) {
	_sb := strings.Builder{}
	peekOne, err := ct.read()
	if err != nil {
		return 0, err
	}
	for true {
		_sb.WriteRune(peekOne)
		currentChar := peekOne
		peekOne, err = ct.read()
		if err != nil {
			return 0, err
		}

		if currentChar == '*' && peekOne == '/' {
			peekOne, err = ct.read()
			if err != nil {
				return 0, err
			}
			break
		}
	}
	if !ct.IsCommentSignificant {
		return ct.NextToken()
	}
	ct.StringValue = _sb.String()
	ct.TokenType = TOKEN_TYPE_COMMENT_SLASH_STAR
	return ct.TokenType, nil
}

func (ct *CustomTokener) checkQuotedCharacter(_q rune) (Token, error) {
	_sb := strings.Builder{}
	peekOne, err := ct.read()
	for err == nil && peekOne != _q && peekOne != '\r' && peekOne != '\n' {
		readPeek := true
		if peekOne == '\\' {
			peekOne, err = ct.readEscaped()
			if err != nil {
				return 0, err
			}
			readPeek = false
			_sb.WriteRune(peekOne)
			peekOne, err = ct.read()
			if err != nil {
				return 0, err
			}
		}

		if readPeek {
			_sb.WriteRune(peekOne)
			peekOne, err = ct.read()
			if err != nil {
				return 0, err
			}
		}
	}

	if peekOne == _q {
		peekOne, err = ct.read()
		if err != nil {
			return 0, err
		}
	}
	ct.peekChar = peekOne
	ct.TokenType = TOKEN_TYPE_STRING
	ct.StringValue = _sb.String()
	return ct.TokenType, nil
}

// allows to break out lexer for here-document like parsing
func (ct *CustomTokener) ReadHereDocument(_qqq string, _qs *strings.Builder) (int, error) {
	for !strings.HasSuffix(_qs.String(), _qqq) {
		_c, err := ct.read()
		if err != nil {
			return 0, err
		}
		_qs.WriteRune(_c)
	}
	str := _qs.String()
	_qs.Reset()
	_qs.WriteString(str[:len(str)-len(_qqq)])

	return TOKEN_TYPE_STRING_MULTILINE, nil
}

func (ct *CustomTokener) reportNextLine() {
	ct.Pos.line++
	ct.Pos.column = 0
}

// scans the input for the next token. It returns the position of the token,
// the token's type, and the literal value.
func (ct *CustomTokener) NextToken() (Token, error) {
	if ct.IsPushBackToken {
		ct.IsPushBackToken = false
		return ct.TokenType, nil
	}

	// initialize token output
	ct.TokenType = 0
	ct.CharacterValue = 0
	ct.NumericValue = 0
	ct.CardinalValue = 0
	ct.StringValue = ""
	ct.ByteValue = nil
	ct.BooleanValue = false

	_tk, err := ct.doNextToken()
	if err == io.EOF {
		ct.TokenType = TOKEN_TYPE_EOF
		return ct.TokenType, nil
	}
	return _tk, err
}

func (ct *CustomTokener) doNextToken() (Token, error) {

	currentChar := rune(0)
	if ct.peekChar == -2 {
		_p, err := ct.read()
		if err != nil {
			return 0, err
		}
		currentChar = _p
	} else {
		currentChar = ct.peekChar
	}

	currentType := ct.checkCharType(currentChar, true)
	for currentType > 0 && (currentType&LEXICAL_WHITE) != 0 {
		// Skip over white space until we hit a new line or a real token
		if currentChar == '\r' {
			if ct.IsEolSignificant {
				ct.lastCR = true
				ct.peekChar = -2
				ct.TokenType = TOKEN_TYPE_EOL
				return ct.TokenType, nil
			}

			_p, err := ct.read()
			if err != nil {
				return 0, err
			}
			if _p == '\n' {
				_p, err = ct.read()
				if err != nil {
					return 0, err
				}
			}
			currentChar = _p
		} else if currentChar == '\n' {
			if ct.IsEolSignificant {
				ct.peekChar = -2
				ct.TokenType = TOKEN_TYPE_EOL
				return ct.TokenType, nil
			}

			_p, err := ct.read()
			if err != nil {
				return 0, err
			}
			currentChar = _p
		} else {
			_p, err := ct.read()
			if err != nil {
				return 0, err
			}
			currentChar = _p
		}

		currentType = ct.checkCharType(currentChar, true)
	}

	if currentType < 0 {
		ct.peekChar = -2
		ct.StringValue = string(currentChar)
		ct.CharacterValue = currentChar
		ct.CardinalValue = int64(currentChar)
		ct.TokenType = Token(-currentType)
		return ct.TokenType, nil
	}
	/**
	 * Check for digits before checking for words since digits can be
	 * contained within words.
	 */
	if (currentType & LEXICAL_DIGIT) != 0 {
		return ct.readNumberLiteral(currentChar)
	}

	// check for literal words and identifiers
	if (currentType & LEXICAL_WORD) != 0 {
		return ct.readWordLiteral(currentChar)
	}

	// check for byte literals if any
	if ct.IsByteLiterals && currentType == LEXICAL_BYTES {
		return ct.readByteLiteral(currentChar)
	}

	// Check for multi-line quoted character
	if ct.IsBacktickQuote && currentType == LEXICAL_QUOTE_MULTILINE {
		return ct.readBackTicks(currentChar)
	}

	// Check for quoted character
	if currentType == LEXICAL_QUOTE {
		return ct.checkQuotedCharacter(currentChar)
	}

	// Do comments, both "//" and "/*stuff*/"
	if currentChar == '/' && (ct.IsSlashSlashComments || ct.IsSlashStarComments) {
		_p, err := ct.read()
		if err != nil {
			return 0, err
		}

		if _p == '*' && ct.IsSlashStarComments {
			return ct.eatSlashStarComments()
		} else if _p == '/' && ct.IsSlashSlashComments {
			return ct.eatSlashSlashComments()
		} else if currentChar != LEXICAL_COMMENT {
			ct.peekChar = _p
			ct.TokenType = '/'
			return ct.TokenType, nil
		}
	}

	// Check for comment character
	if currentType == LEXICAL_COMMENT {
		return ct.eatComments(currentChar)
	}

	_p, err := ct.read()
	if err != nil {
		return 0, err
	}

	ct.peekChar = _p

	ct.TokenType = TOKEN_TYPE_UNKNOWN
	ct.CharacterValue = currentChar
	return ct.TokenType, nil
}

func (ct *CustomTokener) reportNextChar(sz int) {
	ct.Pos.column += sz
}

func (ct *CustomTokener) readEscaped() (rune, error) {
	_c, _err := ct.read()
	if _err != nil {
		return 0, _err
	}
	return ct.readEscape(_c)
}

func (ct *CustomTokener) readEscape(_c rune) (rune, error) {
	if _c == 'a' {
		return 0x7, nil
	}
	if _c == 'b' {
		return 0x8, nil
	}
	if _c == 't' {
		return 0x9, nil
	}
	if _c == 'f' {
		return 0xc, nil
	}
	if _c == 'n' {
		return 0xa, nil
	}
	if _c == 'v' {
		return 0xb, nil
	}
	if _c == 'r' {
		return 0xd, nil
	}
	if _c >= '0' && _c <= '7' {
		// OCTAL \000
		_sb := strings.Builder{}
		_sb.WriteRune(_c)
		_c, err := ct.read()
		if err != nil {
			return 0, err
		}
		_sb.WriteRune(_c)
		_c, err = ct.read()
		if err != nil {
			return 0, err
		}
		_sb.WriteRune(_c)
		_r, err := strconv.ParseUint(_sb.String(), 8, 8)
		if err != nil {
			return 0, err
		}
		return rune(_r), nil
	}
	if _c == 'x' {
		// \xHH
		_sb := strings.Builder{}
		for i := 0; i < 2; i++ {
			_c, err := ct.read()
			if err != nil {
				return 0, err
			}
			_sb.WriteRune(_c)
		}
		_r, err := strconv.ParseUint(_sb.String(), 16, 8)
		if err != nil {
			return 0, err
		}
		return rune(_r), nil
	}
	if _c == 'u' {
		// \uHHHH
		_sb := strings.Builder{}
		for i := 0; i < 4; i++ {
			_c, err := ct.read()
			if err != nil {
				return 0, err
			}
			_sb.WriteRune(_c)
		}
		_r, err := strconv.ParseUint(_sb.String(), 16, 16)
		if err != nil {
			return 0, err
		}
		return rune(_r), nil
	}
	if _c == 'U' {
		// \uHHHHHHHH
		_sb := strings.Builder{}
		for i := 0; i < 8; i++ {
			_c, err := ct.read()
			if err != nil {
				return 0, err
			}
			_sb.WriteRune(_c)
		}
		_r, err := strconv.ParseUint(_sb.String(), 16, 32)
		if err != nil {
			return 0, err
		}
		return rune(_r), nil
	}
	if _c == '(' {
		// \(HHH)
		_sb := strings.Builder{}
		_, err := ct.readSimpleQuoted(')', &_sb)
		if err != nil {
			return 0, err
		}
		_name := _sb.String()
		if strings.HasPrefix(_name, "0x") {
			_r, err := strconv.ParseUint(_name[2:], 16, 32)
			if err != nil {
				return 0, err
			}
			return rune(_r), nil
		}
		return 0, nil
	}
	if _c == '{' {
		// \(HHH)
		_sb := strings.Builder{}
		_, err := ct.readSimpleQuoted('}', &_sb)
		if err != nil {
			return 0, err
		}
		_name := _sb.String()
		if strings.HasPrefix(_name, "0x") {
			_r, err := strconv.ParseUint(_name[2:], 16, 32)
			if err != nil {
				return 0, err
			}
			return rune(_r), nil
		}
		return 0, nil
	}

	return _c, nil
}

func (ct *CustomTokener) readSimpleQuoted(_q rune, _sb *strings.Builder) (*strings.Builder, error) {
	_qq := string(_q)
	for !strings.HasSuffix(_sb.String(), _qq) {
		_c, _err := ct.read()
		if _err != nil {
			return _sb, _err
		}
		_sb.WriteRune(_c)
	}
	str := _sb.String()
	_sb.Reset()
	_sb.WriteString(str[:len(str)-1])
	return _sb, nil
}

func (ct *CustomTokener) readComplexQuoted(_q rune, _sb *strings.Builder) (*strings.Builder, error) {
	_qq := string(_q)
	for !strings.HasSuffix(_sb.String(), _qq) {
		_c, _err := ct.read()
		if _err != nil {
			return _sb, _err
		}
		if _c == '\\' {
			_c, _err = ct.readEscaped()
		}
		_sb.WriteRune(_c)
	}
	str := _sb.String()
	str = strings.TrimLeft(str, "\r\n")
	_sb.Reset()
	_sb.WriteString(str[:len(str)-1])
	str = _sb.String()
	_parts := strings.Split(str, "\n")
	_parts0 := strings.TrimLeft(_parts[0], " \r\t\f")
	_sz := len(_parts[0]) - len(_parts0)
	for i := 0; i < len(_parts); i++ {
		runes := []rune(_parts[i])
		j := checkWsPrefix(runes)
		j = min(_sz, j)
		_parts[i] = string(runes[j:])
	}
	_sb.Reset()
	_sb.WriteString(strings.Join(_parts, "\n"))
	return _sb, nil
}

func (ct *CustomTokener) readByteLiteral(char rune) (Token, error) {
	_sb := strings.Builder{}
	_, err := ct.readSimpleQuoted('>', &_sb)
	if err != nil {
		return 0, err
	}
	_bv, err := hex.DecodeString(ct.StringValue)
	if err != nil {
		return 0, err
	}
	ct.ByteValue = _bv
	ct.TokenType = TOKEN_TYPE_BYTES
	return ct.TokenType, nil
}

func (ct *CustomTokener) readBackTicks(char rune) (Token, error) {
	_sb := strings.Builder{}
	_, err := ct.readComplexQuoted('`', &_sb)
	if err != nil {
		return 0, err
	}
	ct.StringValue = _sb.String()
	ct.TokenType = TOKEN_TYPE_STRING_MULTILINE
	return ct.TokenType, nil
}

func checkWsPrefix(s []rune) int {
	_sz := len(s)
	if _sz == 0 {
		return 0
	}
	i := 0
	for unicode.IsSpace(s[i]) {
		i++
		if i >= _sz {
			break
		}
	}
	return i
}
