package valid

import (
	"errors"
	"github.com/dreamlu/gt/tool/type/cmap"
	"strconv"
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
		//ID:   6,
		Name: "梦",
	}
	t.Log(Valid(test))

	var test2 = &Test{
		ID:   6,
		Name: "梦",
	}
	t.Log(Valid(&test2))

	// json data
	var tests = []*Test{
		{
			ID:   6,
			Name: "梦",
		},
	}
	t.Log(Valid(tests))

	var tests2 = []*Test{
		{
			ID:   6,
			Name: "梦",
		},
	}
	t.Log(Valid(&tests2))

	// form data
	var maps = cmap.NewCMap().Set("id", "1")
	maps["name"] = append(maps["name"], "梦1")
	info := ValidModel(maps, Test{})
	//t.Log(info == nil)
	t.Log(info)
}

// add your custom rule
func TestCustomizeValid(t *testing.T) {
	type Test struct {
		Name string `json:"name" gt:"valid:required,large=3;trans:用户名"`
	}

	// json data
	var test = Test{
		Name: "梦sss",
	}
	t.Log(Valid(test))

	AddRule("large", func(rule string, data interface{}) error {
		num, _ := strconv.Atoi(rule)
		if v, ok := data.(string); ok {
			if length(v) > num {
				return errors.New("最大" + rule)
			}
		}
		return nil
	})

	t.Log(Valid(test))
}
