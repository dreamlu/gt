package util

import (
	"bytes"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
	"unicode"
)

// HumpToLine Hump to underline, like json tag
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

// LineToHump Underscore to hump
func LineToHump(str string) string {
	str = strings.Replace(str, "_", " ", -1)
	str = cases.Title(language.English, cases.NoLower).String(str)
	return strings.Replace(str, " ", "", -1)
}
