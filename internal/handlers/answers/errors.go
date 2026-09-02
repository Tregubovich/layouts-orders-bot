package answers

import "errors"

var (
	ErrWrongOption   = errors.New("wrong option")
	ErrFinishSession = errors.New("should finish session")

	ErrNoSession = errors.New("no session")
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

type NoSessionError struct {
	Message string
}

func (e *NoSessionError) Error() string {
	return e.Message
}

func (e *NoSessionError) Unwrap() error {
	return ErrNoSession
}
