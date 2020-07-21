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
			error = strings.Split(error, ".")[1]
		}
		err = fmt.Errorf("%w", &te.TextError{Msg: error})
	case strings.Contains(error, "Error 1406") || strings.Contains(error, "Error 1264"):
		key := strings.Split(strings.Split(error, "column '")[1], "'")[0]
		err = fmt.Errorf("%w", &te.TextError{Msg: fmt.Sprintf("字段过长[%s]", key)})
	default:
		err = fmt.Errorf("%v", errors.New(error))
	}

	return err
}
