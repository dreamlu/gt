package tag

import (
	"github.com/dreamlu/gt/src/reflect"
	mr "reflect"
	"testing"
)

type User struct {
	Name string `json:"name" gt:"field:user.name;like;soft_del"`
}

type UserD struct {
	User
	Other string `json:"other" gt:"ignore;soft_del"`
}

func (User) TableName() string {
	return "users"
}

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
		t.Log(reflect.Field(u, k.Field))
	}
}

func TestGetTags(t *testing.T) {
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
	m := GetTypeTag(mr.TypeOf(u), "gt")
	//a := GetTypeTags(mr.TypeOf(u), []string{"gt", "json"}, IsGtTagIgnore)
	t.Log(m)
	//t.Log(a)
}

func TestGetJsonTags(t *testing.T) {
	type User struct {
		Name   string `json:"name"`
		Gender int    `json:"-"`
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
		ShopName string `json:"shop_name,omitempty"`
		UserDeX
	}
	m := GetJsonTags(UserMore{})
	t.Log(m)
	var arr []string
	for _, v := range m {
		arr = append(arr, v)
	}
	t.Log(arr)

	type NoTag struct {
		Name string
	}
	n := GetJsonTags(NoTag{})
	t.Log(n)
}
