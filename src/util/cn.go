package util

import (
	"regexp"
	"strings"
)

// GetCn get cn string
// unicode chinese range 4e00-9fa5，16 to 10: 19968 - 40869
func GetCn(str string) (cn string) {
	r := []rune(str)
	var strSlice []string
	for i := 0; i < len(r); i++ {
		if r[i] <= 40869 && r[i] >= 19968 {
			cn = cn + string(r[i])
			strSlice = append(strSlice, cn)
		}
	}
	return
}

func GetCnByRegexp(str string) (cn string) {
	re := regexp.MustCompile(`[\p{Han}]+`)
	matches := re.FindAllString(str, -1)
	cn = strings.Join(matches, "")
	return
}
