package valid

import (
	"github.com/dreamlu/gt/tool/type/cmap"
	"testing"
)

// test validator
// 规则之外, 请额外处理
func TestValidator(t *testing.T) {

	type Test struct {
		ID   int64  `json:"id" gt:"valid:required,min=0,max=5"`
		Name string `json:"name" gt:"valid:required,len=2-5;trans:用户名"`
	}

	// json data
	var test = Test{
		ID:   6,
		Name: "梦",
	}
	t.Log(Valid(test))

	var test2 = &Test{
		ID:   6,
		Name: "梦",
	}
	t.Log(Valid(&test2))

	// json data
	var tests = []Test{
		{
			ID:   6,
			Name: "梦",
		},
	}
	t.Log(Valid(tests))

	// form data
	var maps = cmap.NewCMap().Set("id", "1")
	maps["name"] = append(maps["name"], "梦1")
	info := ValidModel(maps, Test{})
	//t.Log(info == nil)
	t.Log(info)
}
