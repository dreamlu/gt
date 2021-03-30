package tag

import (
	"github.com/dreamlu/gt/tool/reflect"
	"testing"
)

func TestGetGtTags(t *testing.T) {
	type User struct {
		Name   string `json:"name" gt:"field:t1.name;excel:名称"`
		Age    int    `json:"age" gt:"excel:性别"`
		Gender int    `json:"gender"`
	}
	var u = User{
		Name:   "测试",
		Age:    18,
		Gender: 1,
	}
	m := GetGtTags(u)
	t.Log(m)
	for k := range m {
		for _, v := range m[k].GtTags {
			t.Log(v)
		}
		t.Log(reflect.GetDataByFieldName(u, k))
	}
}

func TestGetJsonTags(t *testing.T) {
	type User struct {
		Name string `json:"name"`
	}
	type UserDe struct {
		User
		Other string `json:"other"`
	}

	type UserDeX struct {
		a []string `gt:"ignore"`
		UserDe
		OtherX string `json:"other_x"`
	}

	type UserMore struct {
		ShopName string `json:"shop_name"`
		UserDeX
	}
	m := GetJsonTags(UserMore{})
	for _, v := range m {
		t.Log(v)
	}
	t.Log(m)
}
