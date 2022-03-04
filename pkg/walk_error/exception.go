package walk_error

import "fmt"

type WalkException struct {
	Msg  string
	Nest error
}

func NewWalkError(msg string) *WalkException {
	return &WalkException{Msg: msg}
}

func (w *WalkException) Error() string {
	return fmt.Sprintf("msg %s", w.Msg)
}

func (w *WalkException) Unwrap() error {
	return w.Nest
}

type ParseException struct{}
