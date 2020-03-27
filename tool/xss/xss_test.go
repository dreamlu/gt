package xss

import (
	"fmt"
	"github.com/dreamlu/gt/tool/type/cmap"
	"testing"
)

// test xss
func TestXss(t *testing.T) {
	var maps = make(cmap.CMap)
	maps["name"] = append(maps["name"], "梦 '< and 1=1 \" --")
	XssMap(maps)
	fmt.Println(maps)
}
