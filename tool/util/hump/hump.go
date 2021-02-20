package hump

import (
	"bytes"
	"strings"
	"unicode"
)

// Hump to underline, like json tag
func HumpToLine(str string) string {
	var buffer bytes.Buffer
	for i, v := range str {
		if unicode.IsUpper(v) {
			if i != 0 && unicode.IsLower(rune(str[i-1])) {
				buffer.WriteString("_")
			}
			buffer.WriteRune(unicode.ToLower(v))
		} else {
			buffer.WriteRune(v)
		}
	}
	return buffer.String()
}

// Underscore to hump
func LineToHump(str string) string {
	str = strings.Replace(str, "_", " ", -1)
	str = strings.Title(str)
	return strings.Replace(str, " ", "", -1)
}
