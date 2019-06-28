// @author  dreamlu
package lib

import (
	"github.com/dreamlu/go-tool/util/result"
	"strings"
)

// 数据库错误过滤、转换(友好提示)
func GetSqlError(error string) (info result.MapData) {
	switch {
	case error == "record not found":
		info = result.MapNoResult
	case strings.Contains(error, "Duplicate entry"):
		//error = strings.Replace(error, "Error 1062: Duplicate entry", "", -1)
		errors := strings.Split(error, "for key ")
		//error = "已存在相同数据:" + errors[0]
		error = strings.Trim(errors[1], "'") //自定义数据库唯一约束名
		info = result.GetMapData(result.CodeText, error)
	case strings.Contains(error, "Data too long"):
		error = "存在字段范围过长"
		info = result.GetMapData(result.CodeText, error)
	default:
		info = result.GetMapData(result.CodeText, error)
	}

	return info
}
