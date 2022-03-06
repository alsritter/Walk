package walk

import (
	"io"
	"regexp"
)

type ILexer interface {
	Peek(i int) Token
	Read() Token
}

const regexPat = `\s*((//.*) | ([0-9]+) | ("(\\"|\\\\|\\n|[^"])*") | [A-Z_a-z][A-Z_a-z0-9]*|==|<=|>=|&&|\|\||[[:punct:]])?`

type lexer struct {
	queue   []Token        // list of tokens
	reg     *regexp.Regexp // regular expression
	hasMore bool
	reader  *LineNumberReader // Line number reader
}

func NewLexer(reader io.Reader) ILexer {
	lex := new(lexer)
	lex.queue = make([]Token, 0)
	lex.reg = regexp.MustCompile(regexPat)
	lex.hasMore = true
	lex.reader = NewLineNumberReader(reader)
	return lex
}

func (l *lexer) Read() Token {
	if l.fillQueue(0) {
		l.queue = l.queue[1:] // remove first Token.
		return l.queue[0]
	} else {
		return EOF()
	}
}

func (l *lexer) Peek(i int) Token {
	if l.fillQueue(0) {
		return l.queue[0]
	} else {
		return EOF()
	}
}

func (l *lexer) fillQueue(i int) bool {
	for i >= len(l.queue) {
		if l.hasMore {
			l.readLine()
		} else {
			return false
		}
	}
	return true
}

func (l *lexer) readLine() {
	line := ""
	if err := RecoverToError(func() {
		line = l.reader.ReadLine()
	}); err != nil {
		PanicError(NewParseError3("read a Line error", err))
	}

	if line == "" {
		l.hasMore = false
		return
	}

	lineNo := l.reader.GetLineNumber()
	pos := 0
	endPos := len(line)

	for pos < endPos {
		if s := l.reg.FindString(line[pos:]); s != "" {
			l.addToken(lineNo, s)
			pos = len(s)
		}
	}

}

func (l *lexer) addToken(lineNo int64, tokenStr string) {

}
