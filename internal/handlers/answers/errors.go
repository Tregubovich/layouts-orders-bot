package answers

import "errors"

var (
	ErrWrongOption   = errors.New("wrong option")
	ErrFinishSession = errors.New("should finish session")
)

type WrongOptionError struct {
	Message string
}

func (e *WrongOptionError) Error() string {
	return e.Message
}

func (e *WrongOptionError) Unwrap() error {
	return ErrWrongOption
}
