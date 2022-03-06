package walk

const EOL = "\\n"          // EOL is end of line
var eof = NewBaseToken(-1) // singleton instance

// TOF is end of file. (end of the program)
func EOF() Token {
	return eof
}

type Token interface {
	GetLineNumber() int32
	IsIdentifier() bool
	IsNumber() bool
	IsString() bool
	GetNumber() int32
	GetString() string
	GetText() string
}

type BaseToken struct {
	lineNumber int32
}

func NewBaseToken(lineNumber int32) Token {
	return &BaseToken{lineNumber}
}

func (t *BaseToken) GetLineNumber() int32 { return t.lineNumber }

func (t *BaseToken) IsIdentifier() bool { return false }

func (t *BaseToken) IsNumber() bool { return false }

func (t *BaseToken) IsString() bool { return false }

func (t *BaseToken) GetNumber() int32 { panic(NewWalkError("not number token", nil)) }

func (t *BaseToken) GetString() string { return "" }

func (t *BaseToken) GetText() string { return "" }
