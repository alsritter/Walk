package walk

import "fmt"

const EOL = "\\n"          // EOL is end of line
var eof = NewBaseToken(-1) // singleton instance

// TOF is end of file. (end of the program)
func EOF() Token {
	return eof
}

type Token interface {
	GetLineNumber() int64
	IsIdentifier() bool
	IsNumber() bool
	IsString() bool
	GetNumber() int32
	GetString() string
	GetText() string
}

// ========================BaseToken==========================

type BaseToken struct {
	lineNumber int64
}

func NewBaseToken(lineNumber int64) Token {
	return &BaseToken{lineNumber}
}

func (t *BaseToken) GetLineNumber() int64 { return t.lineNumber }
func (t *BaseToken) IsIdentifier() bool   { return false }
func (t *BaseToken) IsNumber() bool       { return false }
func (t *BaseToken) IsString() bool       { return false }
func (t *BaseToken) GetNumber() int32     { panic(NewWalkError("not number token", nil)) }
func (t *BaseToken) GetString() string    { return "" }
func (t *BaseToken) GetText() string      { return "" }

// ==========================NumToken============================
type NumToken struct {
	value int32
	*BaseToken
}

func NewNumToken(line int64, value int32) *NumToken {
	return &NumToken{value: value, BaseToken: &BaseToken{lineNumber: line}}
}

func (t *NumToken) IsNumber() bool   { return true }
func (t *NumToken) GetText() string  { return fmt.Sprintf("%d", t.value) }
func (t *NumToken) GetNumber() int32 { return t.value }

// =========================StrToken============================
type StrToken struct {
	text string
	*BaseToken
}

func NewStrToken(line int64, str string) *StrToken {
	return &StrToken{text: str, BaseToken: &BaseToken{lineNumber: line}}
}

func (t *StrToken) IsString() bool  { return true }
func (t *StrToken) GetText() string { return t.text }

// =========================IdToken============================
type IdToken struct {
	text string
	*BaseToken
}

func NewIdToken(line int64, id string) *IdToken {
	return &IdToken{text: id, BaseToken: &BaseToken{lineNumber: line}}
}

func (t *IdToken) IsIdentifier() bool { return true }
func (t *IdToken) GetText() string    { return t.text }
