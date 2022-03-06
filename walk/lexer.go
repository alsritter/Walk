package walk

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

type ILexer interface {
	Peek(i int) Token
	Read() Token
}

const regexPat = `\s*((//.*)|([0-9]+)|("(\\"|\\\\|\\n|[^"])*")|[A-Z_a-z][A-Z_a-z0-9]*|==|<=|>=|&&|\|\||[[:punct:]])?`

type lexer struct {
	queue   []Token        // list of tokens
	regex   *regexp.Regexp // regular expression
	hasMore bool
	reader  *LineNumberReader // Line number reader
}

func NewLexer(reader io.Reader) ILexer {
	lex := new(lexer)
	lex.queue = make([]Token, 0)
	lex.regex = regexp.MustCompile(regexPat)
	lex.hasMore = true
	lex.reader = NewLineNumberReader(reader)
	return lex
}

func (l *lexer) Read() Token {
	if l.fillQueue(0) {
		ele := l.queue[0]
		l.queue = l.queue[1:] // remove first Token.
		return ele
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
		if s := l.regex.FindString(line[pos:]); s != "" {
			if err := RecoverToError(func() {
				l.addToken(lineNo, s)
			}); err != nil {
				PanicError(NewParseError3(fmt.Sprintf("bad token at line %d", lineNo), err))
			}

			pos = pos + len(s)
		} else {
			PanicError(NewParseError3(fmt.Sprintf("bad token at line %d", lineNo), nil))
		}
	}
	l.queue = append(l.queue, NewIdToken(lineNo, EOL))
}

func (l *lexer) addToken(lineNo int64, str string) {
	res := l.regex.FindAllStringSubmatch(str, -1)
	if len(res) == 0 {
		return
	}

	matcher := res[0]
	m := matcher[1]
	if m != "" { // if not a space
		if comment := matcher[2]; comment == "" { // if not a comment
			var token Token
			if number := matcher[3]; number != "" {
				num, err := strconv.Atoi(number)
				if err != nil {
					PanicError(NewParseError3(fmt.Sprintf("error parse number word: %s", number), nil))
				}

				token = NewNumToken(lineNo, int32(num))
			} else if str := matcher[4]; str != "" {
				token = NewStrToken(lineNo, l.toStringLiteral(str))
			} else {
				token = NewIdToken(lineNo, m)
			}

			l.queue = append(l.queue, token)
		}
	}
}

// toStringLiteral converts a special symbol inside a string
func (l *lexer) toStringLiteral(s string) string {
	sb := strings.Builder{}
	len := len(s) - 1
	for i := 0; i < len; i++ {
		c := s[i]
		if c == '\\' && i+1 < len {
			c2 := s[i+1]
			if c == '"' || c2 == '\\' {
				i++
				c = s[i]
			} else if c2 == 'n' {
				i++
				c = '\n'
			}
		}
		sb.WriteByte(c)
	}
	return sb.String()
}
