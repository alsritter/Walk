package walk

const EOL = "\\n" // EOL is end of line

// TOF is end of file. (end of the program)
func TOF() *Token {
	return NewToken(-1)
}

type Token struct {
	lineNumber int32
}

func NewToken(lineNumber int32) *Token {
	return &Token{lineNumber}
}

func (t *Token) GetLineNumber() int32 { return t.lineNumber }

func (t *Token) IsIdentifier() bool { return false }

func (t *Token) IsNumber() bool { return false }

func (t *Token) IsString() bool { return false }

// TODO: Implement exception handling of WalkException

func (t *Token) GetNumber() int32 { panic("not number token") }

func (t *Token) GetString() string { return "" }
