package sql

import (
	"errors"
	"fmt"
	"github.com/dreamlu/gt/tool/type/te"
	"testing"
)

func TestGetSQLError(t *testing.T) {
	msg := "record not found"
	err := GetSQLError(msg)
	///fmt.Println(errors.Unwrap(err))
	fmt.Println(errors.As(err, &te.TextErr))
}
