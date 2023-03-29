package errors

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
)

// TextError customize cn text error
type TextError struct {
	Msg string
}

func (s *TextError) Error() string {
	return s.Msg
}

func (s *TextError) Unwrap() error {
	return s
}

func (s *TextError) Is(err error) bool {
	return reflect.TypeOf(err).Name() == reflect.TypeOf(s).Name()
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
