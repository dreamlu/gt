package errors

import (
	"github.com/pkg/errors"
)

func New(text string) error {
	return errors.New(text)
}

func Wrap(err error, text string) error {
	return errors.Wrap(err, text)
}
