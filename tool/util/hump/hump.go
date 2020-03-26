package hump

import (
	"bytes"
	"strings"
	"unicode"
)

// 驼峰转下划线
func HumpToLine(str string) string {
	var buffer bytes.Buffer
	for i, v := range str {
		if unicode.IsUpper(v) {
			if i != 0 {
				buffer.WriteString("_")
			}
			buffer.WriteRune(unicode.ToLower(v))
		} else {
			buffer.WriteRune(v)
		}
	}
	return buffer.String()
}

// 下划线转驼峰
func LineToHump(str string) string {
	str = strings.Replace(str, "_", " ", -1)
	str = strings.Title(str)
	return strings.Replace(str, " ", "", -1)
}
