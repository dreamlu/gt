package mock

import (
	"github.com/dreamlu/gt/src/type/json"
	"github.com/dreamlu/gt/src/type/time"
	"testing"
)

type User struct {
	ID         uint       `json:"id"`
	Name       string     `json:"name"  faker:"name"`
	BirthDate  time.CDate `json:"birth_date" faker:"CDate"`
	CreateTime time.CTime `json:"create_time"`
	Info       json.CJSON `json:"info"`
	Str        []string   `json:"str"`
	Is         float64    `json:"is"`
}

func TestMock(t *testing.T) {
	var user User
	Mock(&user)
	t.Log(user)
}

func TestGetRand(t *testing.T) {
	t.Log(GetRand(Chinese, 10))
}
