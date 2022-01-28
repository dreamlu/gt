package errors

import (
	"fmt"
	"github.com/pkg/errors"
)

// TextError customize cn text error
type TextError struct {
	Msg string
}

func (s *TextError) Error() string {
	return s.Msg
}

var TextErr *TextError

// Text return TextError type
func Text(msg string) error {
	return fmt.Errorf("%w", &TextError{Msg: msg})
}

func New(text string) error {
	return errors.New(text)
}

func Wrap(err error, text string) error {
	return errors.Wrap(err, text)
}

type QueryError struct {
	Query string
	Err   error
}

func (e *QueryError) Error() string {
	return e.Query + ": " + e.Err.Error()
}

func (e *QueryError) Unwrap() error {
	return e.Err
}
