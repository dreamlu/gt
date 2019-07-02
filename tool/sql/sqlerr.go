// @author  dreamlu
package sql

import (
	"github.com/dreamlu/go-tool/tool/result"
	"github.com/pkg/errors"
	"strings"
)

// 数据库错误过滤、转换(友好提示)
func GetSQLError(error string) (err error) {

	switch {
	case error == "record not found":
		err = errors.New(result.MsgNoResult)
	case strings.Contains(error, "Duplicate entry"):
		//error = strings.Replace(error, "Error 1062: Duplicate entry", "", -1)
		errs := strings.Split(error, "for key ")
		//error = "已存在相同数据:" + errors[0]
		error = strings.Trim(errs[1], "'") //自定义数据库唯一约束名
		err = errors.New(error)
	case strings.Contains(error, "Data too long"):
		err = errors.New("存在字段范围过长")
	default:
		err = errors.New(error)
	}

	return err
}
