package walk

import (
	"io"
	"regexp"

	"github.com/alsritter/walk/common"
	"github.com/alsritter/walk/pkg/token"
)

type ILexer interface {
	Peek(i int) token.Token
	Read() token.Token
}

const regexPat = `\s*((//.*) | ([0-9]+) | ("(\\"|\\\\|\\n|[^"])*") | [A-Z_a-z][A-Z_a-z0-9]*|==|<=|>=|&&|\|\||[[:punct:]])?`

type lexer struct {
	queue   []*token.BaseToken // list of tokens
	reg     *regexp.Regexp     // regular expression
	hasMore bool
	reader  *common.LineNumberReader // Line number reader
}

func NewLexer(reader io.Reader) ILexer {
	lex := new(lexer)
	lex.queue = make([]*token.BaseToken, 0)
	lex.reg = regexp.MustCompile(regexPat)
	lex.hasMore = true
	lex.reader = common.NewLineNumberReader(reader)
	return lex
}

func (l *lexer) Read() token.Token {
	if l.fillQueue(0) {
		l.queue = l.queue[1:] // remove first Token.
		return l.queue[0]
	} else {
		return token.EOF()
	}
}

func (l *lexer) Peek(i int) token.Token {
	if l.fillQueue(0) {
		return l.queue[0]
	} else {
		return token.EOF()
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
	line := l.reader.ReadLine()
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
