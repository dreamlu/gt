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

// 继承tag解析测试
func TestGetTags(t *testing.T) {
	type User struct {
		Name string `json:"name"`
	}
	type UserDe struct {
		User
		Other string `json:"other"`
	}

	type UserDeX struct {
		a []string
		UserDe
		OtherX string `json:"other_x"`
	}

	type UserMore struct {
		ShopName string `json:"shop_name"`
		UserDeX
	}
	// test tag
	t.Log(GetTags(UserMore{}))
}
