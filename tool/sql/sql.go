package sql

import (
	"fmt"
	"github.com/dreamlu/gt/tool/result"
	"github.com/dreamlu/gt/tool/type/te"
	"github.com/pkg/errors"
	"strings"
)

// 数据库错误过滤、转换(友好提示)
func GetSQLError(error string) (err error) {

	switch {
	case error == "record not found":
		err = fmt.Errorf("%w", &te.TextError{Msg: result.MsgNoResult})
	case strings.Contains(error, "PRIMARY"):
		err = fmt.Errorf("%w", &te.TextError{Msg: "主键重复"})
	case strings.Contains(error, "Duplicate entry"):
		//error = strings.Replace(error, "Error 1062: Duplicate entry", "", -1)
		errs := strings.Split(error, "for key ")
		error = strings.Trim(errs[1], "'")
		if strings.Contains(error, ".") {
			error = strings.Split(error,".")[1]
		}
		err = fmt.Errorf("%w", &te.TextError{Msg: error})
	case strings.Contains(error, "Data too long"):
		err = fmt.Errorf("%w", &te.TextError{Msg: "存在字段范围过长"})
	default:
		err = fmt.Errorf("%v", errors.New(error))
	}

	return err
}
