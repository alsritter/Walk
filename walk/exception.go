package walk

import (
	"fmt"
	"strconv"
)

type WalkException struct {
	Msg  string
	Next error
}

func NewWalkError(msg string, err error) StError {
	return &WalkException{Msg: msg, Next: err}
}

func (w *WalkException) Error() string {
	return fmt.Sprintf("msg %s", w.Msg)
}

func (w *WalkException) Unwrap() error {
	return w.Next
}

func (w *WalkException) Is(target error) bool {
	_, ok := target.(*WalkException)
	return ok
}

// =======================ParseException==========================

type ParseException struct {
	msg  string
	Next error
}

func NewParseError(msg string, t Token, err error) StError {
	return &ParseException{msg: "syntax error around " + location(t) + ". " + msg, Next: err}
}

func NewParseError2(t Token, err error) StError {
	return &ParseException{msg: "syntax error around " + location(t) + ". ", Next: err}
}

func NewParseError3(msg string, err error) StError {
	return &ParseException{msg: msg, Next: err}
}

func (p *ParseException) Error() string {
	return p.msg
}

func location(t Token) string {
	if t == EOF() {
		return "the last line"
	}
	return "\"" + t.GetText() + "\" at line " + strconv.Itoa(int(t.GetLineNumber()))
}

func (p *ParseException) Unwrap() error {
	return p.Next
}

func (w *ParseException) Is(target error) bool {
	_, ok := target.(*ParseException)
	return ok
}

// =======================IndexOutOfBoundsException===============
type IndexOutOfBoundsException struct {
	msg  string
	Next error
}

func NewIndexOutOfBoundsException(msg string, err error) StError {
	return &IndexOutOfBoundsException{msg: "Index out of " + msg, Next: err}
}

func (i *IndexOutOfBoundsException) Error() string {
	return i.msg
}

func (i *IndexOutOfBoundsException) Unwrap() error {
	return i.Next
}

func (i *IndexOutOfBoundsException) Is(target error) bool {
	_, ok := target.(*IndexOutOfBoundsException)
	return ok
}

// =======================Exception Tools==========================

type StError interface {
	Error() string
	Unwrap() error
	Is(target error) bool
}

func PanicError(err StError) {
	if err != nil {
		panic(err)
	}
}

func RecoverToError(f func()) (retErr StError) {
	defer func() {
		panicErr := recover()
		if panicErr != nil {
			if err, ok := panicErr.(StError); ok {
				retErr = err
			} else {
				retErr = NewWalkError(fmt.Sprint(panicErr), nil)
			}
		}
	}()
	f()
	return nil
}
