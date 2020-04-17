package validator

import (
	"github.com/dreamlu/gt/tool/type/cmap"
	"testing"
)

// test validator
// 规则之外, 请额外处理
func TestValidator(t *testing.T) {

	type Test struct {
		ID   int64  `json:"id" valid:"required,min=0,max=5"`
		Name string `json:"name" valid:"required,len=2-5" trans:"用户名"`
	}

	// form data
	var maps = make(cmap.CMap)
	maps["name"] = append(maps["name"], "梦1")
	info := Valid(maps, Test{})
	t.Log(info)

	// json data
	var test = Test{
		ID:   6,
		Name: "梦1",
	}
	t.Log(Valid(test, Test{}))

}
