package sql

import (
	"errors"
	"fmt"
	"github.com/dreamlu/gt/tool/type/te"
	"testing"
)

func TestGetSQLError(t *testing.T) {
	msg := "Duplicate entry for key 'user.openid 唯一'"
	err := GetSQLError(msg)
	t.Log(err)
	///fmt.Println(errors.Unwrap(err))
	fmt.Println(errors.As(err, &te.TextErr))
}
