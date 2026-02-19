package cli

import (
	"errors"
	"fmt"
)

var ErrHelpRequested = errors.New("help requested")

type HelpError struct {
	Err  error
	Help string
}

func (e *HelpError) Error() string {
	return e.Err.Error()
}

func (e *HelpError) Unwrap() error {
	return e.Err
}

func NewHelpError(message string) *HelpError {
	return &HelpError{
		Err:  fmt.Errorf("%w", ErrHelpRequested),
		Help: message,
	}
}
