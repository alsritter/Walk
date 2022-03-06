package walk

import (
	"bufio"
	"io"
)

type LineNumberReader struct {
	line    int64
	pos     int64
	scanner *bufio.Scanner
}

func NewLineNumberReader(reader io.Reader) *LineNumberReader {
	lineNumberReader := new(LineNumberReader)
	lineNumberReader.scanner = bufio.NewScanner(reader)
	lineNumberReader.scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		advance, token, err = bufio.ScanLines(data, atEOF)
		lineNumberReader.pos += int64(advance)
		return
	})

	return lineNumberReader
}

// read a line
func (r *LineNumberReader) ReadLine() string {
	if r.scanner.Scan() {
		r.line++
		return r.scanner.Text()
	}

	if r.scanner.Err() != nil {
		panic(NewWalkError("", nil))
	}

	return ""
}

// read the line number
func (r *LineNumberReader) GetLineNumber() int64 { return r.line }
